package server

import (
"encoding/json"
"errors"
"log"
"net/http"

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

account, err := s.Store.SetAuthDisplayNameForAccount(r.Context(), currentAccount.ID, request.DisplayName)
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