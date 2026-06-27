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

type DevCurrencyRequest struct {
	PlayerName     string `json:"playerName"`
	CurrencyType   string `json:"currencyType"`
	CurrencyAmount string `json:"currencyAmount"`
}

type DevResetLevelsRequest struct {
	PlayerName   string `json:"playerName"`
	ActivityName string `json:"activityName"`
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

func decodeDevCurrencyRequest(r *http.Request) (DevCurrencyRequest, error) {
	var request DevCurrencyRequest

	if r.Body != nil {
		defer r.Body.Close()
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, err
	}

	request.PlayerName = strings.TrimSpace(request.PlayerName)
	request.CurrencyType = strings.TrimSpace(request.CurrencyType)
	request.CurrencyAmount = strings.TrimSpace(request.CurrencyAmount)

	return request, nil
}

func decodeDevResetLevelsRequest(r *http.Request) (DevResetLevelsRequest, error) {
	var request DevResetLevelsRequest

	if r.Body != nil {
		defer r.Body.Close()
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, err
	}

	request.PlayerName = strings.TrimSpace(request.PlayerName)
	request.ActivityName = strings.TrimSpace(request.ActivityName)

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

func (s *Server) HandleDevAddCurrency(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	request, err := decodeDevCurrencyRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid add currency request")
		return
	}

	result, err := s.Store.DevAddCurrency(r.Context(), request.PlayerName, request.CurrencyType, request.CurrencyAmount)
	if err != nil {
		log.Printf("add currency failed: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) HandleDevRemoveCurrency(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	request, err := decodeDevCurrencyRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid remove currency request")
		return
	}

	result, err := s.Store.DevRemoveCurrency(r.Context(), request.PlayerName, request.CurrencyType, request.CurrencyAmount)
	if err != nil {
		log.Printf("remove currency failed: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) HandleDevResetLevels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	request, err := decodeDevResetLevelsRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid reset levels request")
		return
	}

	result, err := s.Store.DevResetLevels(r.Context(), request.PlayerName, request.ActivityName)
	if err != nil {
		log.Printf("reset levels failed: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}
