package server

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/SorynMochi/Namagotchi/internal/database"
)

func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account, ok := s.AuthAccountFromRequest(r)
		if !ok {
			s.recordSecurityEvent(r, http.StatusUnauthorized, "auth", "login required")
			writeError(w, http.StatusUnauthorized, "login required")
			return
		}

		ctx := database.WithAuthAccountID(r.Context(), account.ID)
		next(w, r.WithContext(ctx))
	}
}

func (s *Server) requireDev(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !devCommandsEnabled() {
			s.recordSecurityEvent(r, http.StatusNotFound, "dev", "dev commands disabled")
			http.NotFound(w, r)
			return
		}

		account, ok := s.AuthAccountFromRequest(r)
		if !ok {
			s.recordSecurityEvent(r, http.StatusUnauthorized, "dev", "login required")
			writeError(w, http.StatusUnauthorized, "login required")
			return
		}

		allowed, err := s.authAccountCanUseDevCommands(r, account)
		if err != nil {
			s.recordSecurityEvent(r, http.StatusInternalServerError, "dev", "dev authorization failed")
			writeError(w, http.StatusInternalServerError, "dev authorization failed")
			return
		}

		if !allowed {
			s.recordSecurityEvent(r, http.StatusForbidden, "dev", "dev command access denied")
			writeError(w, http.StatusForbidden, "dev command access denied")
			return
		}

		ctx := database.WithAuthAccountID(r.Context(), account.ID)
		next(w, r.WithContext(ctx))
	}
}

func (s *Server) authAccountCanUseDevCommands(r *http.Request, account database.AuthAccount) (bool, error) {
	if account.ID < 1 {
		return false, nil
	}

	for _, rawID := range splitEnvList(os.Getenv("NAMIGOTCHI_DEV_ACCOUNT_IDS")) {
		allowedID, err := strconv.ParseInt(rawID, 10, 64)
		if err == nil && allowedID == account.ID {
			return true, nil
		}
	}

	for _, email := range splitEnvList(os.Getenv("NAMIGOTCHI_DEV_GOOGLE_EMAILS")) {
		allowed, err := s.Store.AuthAccountHasVerifiedProviderEmail(r.Context(), account.ID, "google", email)
		if err != nil {
			return false, err
		}

		if allowed {
			return true, nil
		}
	}

	return false, nil
}

func devCommandsEnabled() bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv("NAMIGOTCHI_DEV_COMMANDS")))

	return value == "1" ||
		value == "true" ||
		value == "yes" ||
		value == "on"
}

func splitEnvList(value string) []string {
	fields := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ';' || r == '\n' || r == '\r' || r == '\t' || r == ' '
	})

	result := make([]string, 0, len(fields))
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field != "" {
			result = append(result, field)
		}
	}

	return result
}
