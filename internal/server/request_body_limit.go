package server

import (
	"net/http"
	"strings"
)

const apiRequestBodyLimitBytes int64 = 64 << 10

func (s *Server) withRequestBodyLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if shouldLimitRequestBody(r) {
			if r.ContentLength > apiRequestBodyLimitBytes {
				s.recordSecurityEvent(r, http.StatusRequestEntityTooLarge, "request_body", "request body too large")
				writeError(w, http.StatusRequestEntityTooLarge, "request body too large")
				return
			}

			r.Body = http.MaxBytesReader(w, r.Body, apiRequestBodyLimitBytes)
		}

		next.ServeHTTP(w, r)
	})
}

func shouldLimitRequestBody(r *http.Request) bool {
	if r == nil || r.Body == nil {
		return false
	}

	if !strings.HasPrefix(r.URL.Path, "/api/") {
		return false
	}

	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return true
	default:
		return false
	}
}
