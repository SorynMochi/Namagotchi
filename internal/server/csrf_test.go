package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequireCSRFAllowsSafeMethods(t *testing.T) {
	s := &Server{}
	called := false

	handler := s.requireCSRF(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	recorder := httptest.NewRecorder()

	handler(recorder, request)

	if !called {
		t.Fatal("expected safe request to reach next handler")
	}

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}
}

func TestRequireCSRFRejectsPostWithoutToken(t *testing.T) {
	s := &Server{}
	called := false

	handler := s.requireCSRF(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodPost, "/test", nil)
	recorder := httptest.NewRecorder()

	handler(recorder, request)

	if called {
		t.Fatal("expected request without CSRF token to be blocked")
	}

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), "csrf token required") {
		t.Fatalf("expected csrf token required message, got %q", recorder.Body.String())
	}
}

func TestRequireCSRFRejectsPostWithMismatchedToken(t *testing.T) {
	s := &Server{}
	called := false

	handler := s.requireCSRF(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodPost, "/test", nil)
	request.AddCookie(&http.Cookie{
		Name:  csrfCookieName,
		Value: "cookie-token",
	})
	request.Header.Set(csrfHeaderName, "header-token")

	recorder := httptest.NewRecorder()

	handler(recorder, request)

	if called {
		t.Fatal("expected request with mismatched CSRF token to be blocked")
	}

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), "csrf token invalid") {
		t.Fatalf("expected csrf token invalid message, got %q", recorder.Body.String())
	}
}

func TestRequireCSRFAllowsPostWithMatchingToken(t *testing.T) {
	s := &Server{}
	called := false

	handler := s.requireCSRF(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodPost, "/test", nil)
	request.AddCookie(&http.Cookie{
		Name:  csrfCookieName,
		Value: "matching-token",
	})
	request.Header.Set(csrfHeaderName, "matching-token")

	recorder := httptest.NewRecorder()

	handler(recorder, request)

	if !called {
		t.Fatal("expected request with matching CSRF token to reach next handler")
	}

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}
}

func TestHandleCSRFTokenIssuesCookieAndResponse(t *testing.T) {
	s := &Server{}

	request := httptest.NewRequest(http.MethodGet, "/api/auth/csrf", nil)
	recorder := httptest.NewRecorder()

	s.HandleCSRFToken(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	response := recorder.Result()
	defer response.Body.Close()

	var csrfCookie *http.Cookie
	for _, cookie := range response.Cookies() {
		if cookie.Name == csrfCookieName {
			csrfCookie = cookie
			break
		}
	}

	if csrfCookie == nil {
		t.Fatal("expected CSRF cookie to be set")
	}

	if csrfCookie.Value == "" {
		t.Fatal("expected CSRF cookie value to be non-empty")
	}

	body := recorder.Body.String()

	if !strings.Contains(body, `"ok":true`) {
		t.Fatalf("expected ok response, got %q", body)
	}

	if !strings.Contains(body, csrfCookie.Value) {
		t.Fatalf("expected response body to contain CSRF token")
	}
}
