package server

import (
	"log"
	"net/http"

	"github.com/SorynMochi/Namagotchi/internal/database"
)

type devAuditResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *devAuditResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (s *Server) withDevAudit(command string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recorder := &devAuditResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next(recorder, r)

		accountID, _ := database.AuthAccountIDFromContext(r.Context())

		if err := s.Store.RecordDevAuditLog(r.Context(), database.DevAuditLog{
			AccountID:  accountID,
			Command:    command,
			Method:     r.Method,
			Path:       r.URL.Path,
			StatusCode: recorder.statusCode,
			RemoteAddr: r.RemoteAddr,
			UserAgent:  r.UserAgent(),
		}); err != nil {
			log.Printf("record dev audit log failed: %v", err)
		}
	}
}

func (s *Server) HandleDevAuditLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	logs, err := s.Store.RecentDevAuditLogs(r.Context(), 100)
	if err != nil {
		log.Printf("get dev audit logs failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to load dev audit logs")
		return
	}

	writeJSON(w, http.StatusOK, logs)
}
