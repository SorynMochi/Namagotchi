package server

import (
	"fmt"
	"net/http"

	"github.com/SorynMochi/Namagotchi/internal/database"
)

func (s *Server) playerStatusForRequest(r *http.Request) (*database.PlayerStatus, error) {
	accountID, ok := database.AuthAccountIDFromContext(r.Context())
	if !ok || accountID < 1 {
		return nil, fmt.Errorf("authenticated account id missing from request context")
	}

	return s.Store.GetPlayerStatusForAccount(r.Context(), accountID)
}
