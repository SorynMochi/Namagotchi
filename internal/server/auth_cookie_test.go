package server

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthCookieSecureCanBeForcedOn(t *testing.T) {
	t.Setenv("AUTH_SECURE_COOKIE", "true")
	t.Setenv("TRUST_PROXY_HEADERS", "")

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	if !authCookieShouldBeSecure(request) {
		t.Fatal("expected secure cookie to be forced on")
	}
}

func TestAuthCookieSecureCanBeForcedOff(t *testing.T) {
	t.Setenv("AUTH_SECURE_COOKIE", "false")
	t.Setenv("TRUST_PROXY_HEADERS", "1")

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("X-Forwarded-Proto", "https")

	if authCookieShouldBeSecure(request) {
		t.Fatal("expected secure cookie to be forced off")
	}
}

func TestAuthCookieSecureUsesDirectTLS(t *testing.T) {
	t.Setenv("AUTH_SECURE_COOKIE", "")
	t.Setenv("TRUST_PROXY_HEADERS", "")

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.TLS = &tls.ConnectionState{}

	if !authCookieShouldBeSecure(request) {
		t.Fatal("expected TLS request to use secure cookie")
	}
}

func TestAuthCookieSecureIgnoresUntrustedForwardedProto(t *testing.T) {
	t.Setenv("AUTH_SECURE_COOKIE", "")
	t.Setenv("TRUST_PROXY_HEADERS", "")

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("X-Forwarded-Proto", "https")

	if authCookieShouldBeSecure(request) {
		t.Fatal("expected forwarded proto to be ignored unless proxy headers are trusted")
	}
}

func TestAuthCookieSecureUsesTrustedForwardedProto(t *testing.T) {
	t.Setenv("AUTH_SECURE_COOKIE", "")
	t.Setenv("TRUST_PROXY_HEADERS", "1")

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("X-Forwarded-Proto", "https")

	if !authCookieShouldBeSecure(request) {
		t.Fatal("expected trusted forwarded proto to use secure cookie")
	}
}

func TestPublicRequestSchemeUsesTrustedForwardedProtoFirstValue(t *testing.T) {
	t.Setenv("TRUST_PROXY_HEADERS", "1")

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("X-Forwarded-Proto", "https, http")

	got := publicRequestScheme(request)
	if got != "https" {
		t.Fatalf("expected https scheme, got %q", got)
	}
}
