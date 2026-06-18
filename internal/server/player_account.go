package server

import (
	"fmt"
	"net/http"

	"github.com/SorynMochi/Namagotchi/internal/database"
)

func accountIDForRequest(r *http.Request) (int64, error) {
	accountID, ok := database.AuthAccountIDFromContext(r.Context())
	if !ok || accountID < 1 {
		return 0, fmt.Errorf("authenticated account id missing from request context")
	}

	return accountID, nil
}

func (s *Server) playerIDForRequest(r *http.Request) (int64, error) {
	accountID, err := accountIDForRequest(r)
	if err != nil {
		return 0, err
	}

	return s.Store.PlayerIDForAccount(r.Context(), accountID)
}
