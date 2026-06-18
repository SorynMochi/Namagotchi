package server

import (
"net/http"
"os"
"strings"

"github.com/SorynMochi/Namagotchi/internal/database"
)

func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
return func(w http.ResponseWriter, r *http.Request) {
account, ok := s.AuthAccountFromRequest(r)
if !ok {
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
http.NotFound(w, r)
return
}

next(w, r)
}
}

func devCommandsEnabled() bool {
value := strings.TrimSpace(strings.ToLower(os.Getenv("NAMIGOTCHI_DEV_COMMANDS")))

return value == "1" ||
value == "true" ||
value == "yes" ||
value == "on"
}
