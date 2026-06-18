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
