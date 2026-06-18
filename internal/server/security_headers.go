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

		if strings.HasPrefix(r.URL.Path, "/api/") {
			setHeaderIfEmpty(w, "Cache-Control", "no-store")
		}

		if publicRequestScheme(r) == "https" {
			setHeaderIfEmpty(w, "Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		next.ServeHTTP(w, r)
	})
}

func setHeaderIfEmpty(w http.ResponseWriter, name string, value string) {
	if w.Header().Get(name) == "" {
		w.Header().Set(name, value)
	}
}
