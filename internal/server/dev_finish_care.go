package server

import (
	"log"
	"net/http"

	"github.com/SorynMochi/Namagotchi/internal/database"
)

type DevFinishCareResponse struct {
	OK      bool                    `json:"ok"`
	Message string                  `json:"message"`
	Care    database.CareQueueState `json:"care"`
}

func (s *Server) HandleFinishCareAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	careState, message, err := s.Store.FinishActiveDevCareAction(r.Context())
	if err != nil {
		log.Printf("finish active care action failed: %v", err)
		writeError(w, http.StatusInternalServerError, "finish active care action failed")
		return
	}

	writeJSON(w, http.StatusOK, DevFinishCareResponse{
		OK:      true,
		Message: message,
		Care:    careState,
	})
}
