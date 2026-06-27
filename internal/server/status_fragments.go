package server

import (
	"log"
	"net/http"

	"github.com/SorynMochi/Namagotchi/internal/database"
)

type PlayerCoreStatusResponse struct {
	Player     database.Player          `json:"player"`
	Companion  database.CompanionState  `json:"companion"`
	Resources  database.PlayerResources `json:"resources"`
	Activities database.ActivitySkills  `json:"activities"`
	Tick       database.TickState       `json:"tick"`
	Wardrobe   database.WardrobeStatus  `json:"wardrobe"`
}

type PlayerCareStatusResponse struct {
	Care database.CareQueueState `json:"care"`
}

type PlayerPlaydeckStatusResponse struct {
	Playdeck database.PlaydeckStatus `json:"playdeck"`
}

type PlayerWardrobeStatusResponse struct {
	Wardrobe database.WardrobeStatus `json:"wardrobe"`
}

type PlayerResourcesStatusResponse struct {
	Resources  database.PlayerResources `json:"resources"`
	Activities database.ActivitySkills  `json:"activities"`
	Tick       database.TickState       `json:"tick"`
}

func (s *Server) HandlePlayerCoreStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	status, err := s.playerStatusForRequest(r)
	if err != nil {
		log.Printf("get player core status failed: %v", err)
		writeError(w, http.StatusNotFound, "player status not found")
		return
	}

	s.applyAccountCreatedAtToStatus(r, status)

	writeJSON(w, http.StatusOK, PlayerCoreStatusResponse{
		Player:     status.Player,
		Companion:  status.Companion,
		Resources:  status.Resources,
		Activities: status.Activities,
		Tick:       status.Tick,
		Wardrobe:   status.Wardrobe,
	})
}

func (s *Server) HandlePlayerCareStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	playerID, err := s.playerIDForRequest(r)
	if err != nil {
		log.Printf("get player for care status failed: %v", err)
		writeError(w, http.StatusNotFound, "player not found")
		return
	}

	care, err := s.Store.GetCareQueueState(r.Context(), playerID)
	if err != nil {
		log.Printf("get care status failed: %v", err)
		writeError(w, http.StatusInternalServerError, "care status failed")
		return
	}

	writeJSON(w, http.StatusOK, PlayerCareStatusResponse{
		Care: care,
	})
}

func (s *Server) HandlePlayerPlaydeckStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	playerID, err := s.playerIDForRequest(r)
	if err != nil {
		log.Printf("get player for playdeck status failed: %v", err)
		writeError(w, http.StatusNotFound, "player not found")
		return
	}

	playdeck, err := s.Store.GetPlaydeckStatus(r.Context(), playerID)
	if err != nil {
		log.Printf("get playdeck status failed: %v", err)
		writeError(w, http.StatusInternalServerError, "playdeck status failed")
		return
	}

	writeJSON(w, http.StatusOK, PlayerPlaydeckStatusResponse{
		Playdeck: playdeck,
	})
}

func (s *Server) HandlePlayerWardrobeStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	playerID, err := s.playerIDForRequest(r)
	if err != nil {
		log.Printf("get player for wardrobe status failed: %v", err)
		writeError(w, http.StatusNotFound, "player not found")
		return
	}

	wardrobe, err := s.Store.GetWardrobeStatus(r.Context(), playerID)
	if err != nil {
		log.Printf("get wardrobe status failed: %v", err)
		writeError(w, http.StatusInternalServerError, "wardrobe status failed")
		return
	}

	writeJSON(w, http.StatusOK, PlayerWardrobeStatusResponse{
		Wardrobe: wardrobe,
	})
}

func (s *Server) HandlePlayerResourcesStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	status, err := s.playerStatusForRequest(r)
	if err != nil {
		log.Printf("get player resources status failed: %v", err)
		writeError(w, http.StatusNotFound, "player status not found")
		return
	}

	writeJSON(w, http.StatusOK, PlayerResourcesStatusResponse{
		Resources:  status.Resources,
		Activities: status.Activities,
		Tick:       status.Tick,
	})
}
