package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/SorynMochi/Namagotchi/internal/database"
)

type AuthDisplayNameRequest struct {
	DisplayName string `json:"displayName"`
}

func (s *Server) HandleAuthDisplayName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	currentAccount, ok := s.AuthAccountFromRequest(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "login required")
		return
	}

	var request AuthDisplayNameRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid display name request")
		return
	}

	if !displayNameChangeAllowedDuringOnboarding(currentAccount, request.DisplayName) {
		writeError(w, http.StatusForbidden, "display name is already set. Name changes will be available later.")
		return
	}

	allowReservedName := reservedDisplayNameAllowedForAccount(currentAccount)

	account, err := s.Store.SetAuthDisplayNameForAccountAllowReserved(
		r.Context(),
		currentAccount.ID,
		request.DisplayName,
		allowReservedName,
	)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrAuthDisplayNameInvalid):
			writeError(w, http.StatusBadRequest, "display name must be 3 to 15 characters and may only use letters, numbers, underscores, or hyphens")
		case errors.Is(err, database.ErrAuthDisplayNameReserved):
			writeError(w, http.StatusBadRequest, "display name is reserved")
		case errors.Is(err, database.ErrAuthDisplayNameTaken):
			writeError(w, http.StatusConflict, "display name is already taken")
		default:
			log.Printf("set auth display name failed: %v", err)
			writeError(w, http.StatusInternalServerError, "display name update failed")
		}
		return
	}

	writeJSON(w, http.StatusOK, AuthResponse{
		OK:       true,
		LoggedIn: true,
		Account:  &account,
		Message:  "Display name updated.",
	})
}

func displayNameChangeAllowedDuringOnboarding(account database.AuthAccount, requestedDisplayName string) bool {
	if isOnboardingPlaceholderDisplayName(account.DisplayName) {
		return true
	}

	return strings.EqualFold(
		strings.TrimSpace(account.DisplayName),
		strings.TrimSpace(requestedDisplayName),
	)
}

func isOnboardingPlaceholderDisplayName(displayName string) bool {
	name := strings.TrimSpace(displayName)
	if name == "" {
		return true
	}

	lowerName := strings.ToLower(name)
	if lowerName == "player" {
		return true
	}

	if strings.HasPrefix(lowerName, "player_") {
		return isASCIIDigits(lowerName[len("player_"):])
	}

	if strings.HasPrefix(lowerName, "namifan") {
		return isASCIIDigits(lowerName[len("namifan"):])
	}

	return false
}

func isASCIIDigits(value string) bool {
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

func reservedDisplayNameAllowedForAccount(account database.AuthAccount) bool {
	accountEmail := strings.TrimSpace(strings.ToLower(account.Email))
	if accountEmail == "" {
		return false
	}

	allowedEmails := strings.Split(os.Getenv("NAMIGOTCHI_RESERVED_DISPLAY_NAME_EMAILS"), ",")
	for _, allowedEmail := range allowedEmails {
		if accountEmail == strings.TrimSpace(strings.ToLower(allowedEmail)) {
			return true
		}
	}

	return false
}
