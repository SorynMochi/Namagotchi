package server

import (
	"log"
	"net/http"

	"github.com/SorynMochi/Namagotchi/internal/database"
)

func (s *Server) recordSecurityEvent(r *http.Request, statusCode int, eventType string, reason string) {
	if s == nil || s.Store == nil {
		log.Printf("security event: type=%s reason=%s method=%s path=%s status=%d", eventType, reason, r.Method, r.URL.Path, statusCode)
		return
	}

	accountID, _ := database.AuthAccountIDFromContext(r.Context())

	if accountID < 1 {
		if account, ok := s.AuthAccountFromRequest(r); ok {
			accountID = account.ID
		}
	}

	if err := s.Store.RecordSecurityEventLog(r.Context(), database.SecurityEventLog{
		AccountID:  accountID,
		EventType:  eventType,
		Reason:     reason,
		Method:     r.Method,
		Path:       r.URL.Path,
		StatusCode: statusCode,
		RemoteAddr: r.RemoteAddr,
		UserAgent:  r.UserAgent(),
	}); err != nil {
		log.Printf("record security event failed: %v", err)
	}
}
