package server

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type rateLimitBucket struct {
	Count     int
	ResetAt   time.Time
	UpdatedAt time.Time
}

type rateLimiter struct {
	mu      sync.Mutex
	buckets map[string]rateLimitBucket
}

func newRateLimiter() *rateLimiter {
	return &rateLimiter{
		buckets: make(map[string]rateLimitBucket),
	}
}

func (l *rateLimiter) allow(key string, limit int, window time.Duration, now time.Time) bool {
	if l == nil {
		return true
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.cleanupExpired(now)

	bucket := l.buckets[key]
	if bucket.ResetAt.IsZero() || now.After(bucket.ResetAt) {
		bucket = rateLimitBucket{
			Count:     0,
			ResetAt:   now.Add(window),
			UpdatedAt: now,
		}
	}

	bucket.Count++
	bucket.UpdatedAt = now
	l.buckets[key] = bucket

	return bucket.Count <= limit
}

func (l *rateLimiter) cleanupExpired(now time.Time) {
	for key, bucket := range l.buckets {
		if now.Sub(bucket.UpdatedAt) > 2*time.Hour || now.After(bucket.ResetAt.Add(time.Hour)) {
			delete(l.buckets, key)
		}
	}
}

var globalRateLimiter = newRateLimiter()

func (s *Server) requireRateLimit(name string, limit int, window time.Duration, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := name + ":" + clientIPForRateLimit(r)

		if !globalRateLimiter.allow(key, limit, window, time.Now().UTC()) {
			s.recordSecurityEvent(r, http.StatusTooManyRequests, "rate_limit", name+" rate limit exceeded")
			writeError(w, http.StatusTooManyRequests, "too many attempts; please wait and try again")
			return
		}

		next(w, r)
	}
}

func clientIPForRateLimit(r *http.Request) string {
	forwardedFor := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if forwardedFor != "" {
		first := strings.TrimSpace(strings.Split(forwardedFor, ",")[0])
		if first != "" {
			return first
		}
	}

	realIP := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if realIP != "" {
		return realIP
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}
