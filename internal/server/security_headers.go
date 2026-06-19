package server

import (
	"net/http"
	"strings"
)

func (s *Server) withSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setHeaderIfEmpty(w, "X-Content-Type-Options", "nosniff")
		setHeaderIfEmpty(w, "X-Frame-Options", "DENY")
		setHeaderIfEmpty(w, "Referrer-Policy", "strict-origin-when-cross-origin")
		setHeaderIfEmpty(w, "Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()")
		setHeaderIfEmpty(w, "Content-Security-Policy-Report-Only", contentSecurityPolicyReportOnlyValue())

		if strings.HasPrefix(r.URL.Path, "/api/") {
			setHeaderIfEmpty(w, "Cache-Control", "no-store")
		}

		if publicRequestScheme(r) == "https" {
			setHeaderIfEmpty(w, "Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		next.ServeHTTP(w, r)
	})
}

func contentSecurityPolicyReportOnlyValue() string {
	return strings.Join([]string{
		"default-src 'self'",
		"script-src 'self'",
		"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
		"font-src 'self' https://fonts.gstatic.com data:",
		"img-src 'self' data: blob:",
		"media-src 'self' blob:",
		"connect-src 'self'",
		"object-src 'none'",
		"base-uri 'self'",
		"frame-ancestors 'none'",
		"form-action 'self'",
	}, "; ")
}

func setHeaderIfEmpty(w http.ResponseWriter, name string, value string) {
	if w.Header().Get(name) == "" {
		w.Header().Set(name, value)
	}
}
