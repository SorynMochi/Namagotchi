package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiterAllowsWithinLimit(t *testing.T) {
	limiter := newRateLimiter()
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	if !limiter.allow("test-key", 2, time.Minute, now) {
		t.Fatal("expected first request to be allowed")
	}

	if !limiter.allow("test-key", 2, time.Minute, now.Add(time.Second)) {
		t.Fatal("expected second request to be allowed")
	}
}

func TestRateLimiterBlocksOverLimit(t *testing.T) {
	limiter := newRateLimiter()
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	if !limiter.allow("test-key", 2, time.Minute, now) {
		t.Fatal("expected first request to be allowed")
	}

	if !limiter.allow("test-key", 2, time.Minute, now.Add(time.Second)) {
		t.Fatal("expected second request to be allowed")
	}

	if limiter.allow("test-key", 2, time.Minute, now.Add(2*time.Second)) {
		t.Fatal("expected third request to be blocked")
	}
}

func TestRateLimiterResetsAfterWindow(t *testing.T) {
	limiter := newRateLimiter()
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	if !limiter.allow("test-key", 1, time.Minute, now) {
		t.Fatal("expected first request to be allowed")
	}

	if limiter.allow("test-key", 1, time.Minute, now.Add(10*time.Second)) {
		t.Fatal("expected second request inside window to be blocked")
	}

	if !limiter.allow("test-key", 1, time.Minute, now.Add(2*time.Minute)) {
		t.Fatal("expected request after reset window to be allowed")
	}
}

func TestClientIPForRateLimitUsesForwardedFor(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.RemoteAddr = "10.0.0.1:12345"
	request.Header.Set("X-Forwarded-For", "203.0.113.7, 10.0.0.2")

	got := clientIPForRateLimit(request)
	if got != "203.0.113.7" {
		t.Fatalf("expected forwarded IP, got %q", got)
	}
}
