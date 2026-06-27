package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type DevPlayerNameRequest struct {
	PlayerName string `json:"playerName"`
}

func decodeDevPlayerNameRequest(r *http.Request) (DevPlayerNameRequest, error) {
	var request DevPlayerNameRequest

	if r.Body != nil {
		defer r.Body.Close()
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, err
	}

	request.PlayerName = strings.TrimSpace(request.PlayerName)
	return request, nil
}

func (s *Server) HandleDevResetChain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	request, err := decodeDevPlayerNameRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid reset chain request")
		return
	}

	result, err := s.Store.DevResetChain(r.Context(), request.PlayerName, false)
	if err != nil {
		log.Printf("reset chain failed: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) HandleDevResetMaxChain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	request, err := decodeDevPlayerNameRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid reset max chain request")
		return
	}

	result, err := s.Store.DevResetChain(r.Context(), request.PlayerName, true)
	if err != nil {
		log.Printf("reset max chain failed: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) HandleDevClearWardrobe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	request, err := decodeDevPlayerNameRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid clear wardrobe request")
		return
	}

	result, err := s.Store.DevClearWardrobe(r.Context(), request.PlayerName)
	if err != nil {
		log.Printf("clear wardrobe failed: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}
