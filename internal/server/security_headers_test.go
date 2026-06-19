package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSecurityHeadersAreSet(t *testing.T) {
	server := &Server{}

	handler := server.withSecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	checks := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
		"Permissions-Policy":     "camera=(), microphone=(), geolocation=(), payment=()",
	}

	for header, expected := range checks {
		if got := response.Header.Get(header); got != expected {
			t.Fatalf("expected %s=%q, got %q", header, expected, got)
		}
	}
}

func TestSecurityHeadersSetNoStoreForAPI(t *testing.T) {
	server := &Server{}

	handler := server.withSecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if got := response.Header.Get("Cache-Control"); got != "no-store" {
		t.Fatalf("expected Cache-Control no-store for API path, got %q", got)
	}
}

func TestSecurityHeadersDoNotSetNoStoreForStaticPath(t *testing.T) {
	server := &Server{}

	handler := server.withSecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if got := response.Header.Get("Cache-Control"); got != "" {
		t.Fatalf("expected no Cache-Control for static path, got %q", got)
	}
}

func TestSecurityHeadersSetHSTSForHTTPS(t *testing.T) {
	t.Setenv("TRUST_PROXY_HEADERS", "1")

	server := &Server{}

	handler := server.withSecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("X-Forwarded-Proto", "https")

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if got := response.Header.Get("Strict-Transport-Security"); got != "max-age=31536000; includeSubDomains" {
		t.Fatalf("expected HSTS header for HTTPS request, got %q", got)
	}
}

func TestSecurityHeadersSetCSPReportOnly(t *testing.T) {
	server := &Server{}

	handler := server.withSecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	got := response.Header.Get("Content-Security-Policy-Report-Only")
	if got == "" {
		t.Fatal("expected CSP report-only header")
	}

	requiredParts := []string{
		"default-src 'self'",
		"script-src 'self'",
		"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
		"font-src 'self' https://fonts.gstatic.com data:",
		"object-src 'none'",
		"frame-ancestors 'none'",
	}

	for _, part := range requiredParts {
		if !strings.Contains(got, part) {
			t.Fatalf("expected CSP report-only header to contain %q, got %q", part, got)
		}
	}
}
