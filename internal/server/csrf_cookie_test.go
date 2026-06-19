package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCSRFCookieIsReadableByJavaScript(t *testing.T) {
	t.Setenv("AUTH_SECURE_COOKIE", "")
	t.Setenv("TRUST_PROXY_HEADERS", "")

	request := httptest.NewRequest(http.MethodGet, "/api/auth/csrf", nil)
	recorder := httptest.NewRecorder()

	setCSRFCookie(recorder, request, "test-token")

	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected one csrf cookie, got %d", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != csrfCookieName {
		t.Fatalf("expected csrf cookie name %q, got %q", csrfCookieName, cookie.Name)
	}

	if cookie.HttpOnly {
		t.Fatal("expected csrf cookie to be readable by JavaScript")
	}

	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("expected SameSite=Lax, got %v", cookie.SameSite)
	}
}

func TestCSRFCookieUsesSecureWhenForced(t *testing.T) {
	t.Setenv("AUTH_SECURE_COOKIE", "1")
	t.Setenv("TRUST_PROXY_HEADERS", "")

	request := httptest.NewRequest(http.MethodGet, "/api/auth/csrf", nil)
	recorder := httptest.NewRecorder()

	setCSRFCookie(recorder, request, "test-token")

	cookie := recorder.Result().Cookies()[0]
	if !cookie.Secure {
		t.Fatal("expected csrf cookie to be secure when AUTH_SECURE_COOKIE=1")
	}
}

func TestCSRFCookieUsesSecureForTrustedHTTPSProxy(t *testing.T) {
	t.Setenv("AUTH_SECURE_COOKIE", "")
	t.Setenv("TRUST_PROXY_HEADERS", "1")

	request := httptest.NewRequest(http.MethodGet, "/api/auth/csrf", nil)
	request.Header.Set("X-Forwarded-Proto", "https")
	recorder := httptest.NewRecorder()

	setCSRFCookie(recorder, request, "test-token")

	cookie := recorder.Result().Cookies()[0]
	if !cookie.Secure {
		t.Fatal("expected csrf cookie to be secure for trusted HTTPS proxy request")
	}
}

func TestCSRFCookieIgnoresUntrustedHTTPSProxyHeader(t *testing.T) {
	t.Setenv("AUTH_SECURE_COOKIE", "")
	t.Setenv("TRUST_PROXY_HEADERS", "")

	request := httptest.NewRequest(http.MethodGet, "/api/auth/csrf", nil)
	request.Header.Set("X-Forwarded-Proto", "https")
	recorder := httptest.NewRecorder()

	setCSRFCookie(recorder, request, "test-token")

	cookie := recorder.Result().Cookies()[0]
	if cookie.Secure {
		t.Fatal("expected csrf cookie to ignore X-Forwarded-Proto unless proxy headers are trusted")
	}
}

func TestEnsureCSRFTokenReusesExistingCookie(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/auth/csrf", nil)
	request.AddCookie(&http.Cookie{
		Name:  csrfCookieName,
		Value: "existing-token",
	})

	recorder := httptest.NewRecorder()

	got := ensureCSRFToken(recorder, request)
	if got != "existing-token" {
		t.Fatalf("expected existing token, got %q", got)
	}

	if cookies := recorder.Result().Cookies(); len(cookies) != 0 {
		t.Fatalf("expected no new cookie when csrf cookie already exists, got %d", len(cookies))
	}
}

func TestEnsureCSRFTokenCreatesCookieWhenMissing(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/auth/csrf", nil)
	recorder := httptest.NewRecorder()

	got := ensureCSRFToken(recorder, request)
	if got == "" {
		t.Fatal("expected generated csrf token")
	}

	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected generated csrf cookie, got %d", len(cookies))
	}

	if cookies[0].Value != got {
		t.Fatalf("expected cookie value to match generated token")
	}
}
