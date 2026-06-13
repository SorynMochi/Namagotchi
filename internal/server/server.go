package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SorynMochi/Namagotchi/internal/database"
)

type Server struct {
	Store     *database.Store
	StartedAt time.Time
}

type StatusResponse struct {
	Server      string `json:"server"`
	Database    string `json:"database"`
	Timestamp   string `json:"timestamp"`
	Uptime      string `json:"uptime"`
	Version     string `json:"version"`
	OnlineUsers int    `json:"onlineUsers"`
}

type MessageResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

type GatheringRequest struct {
	Task string `json:"task"`
}

type CareActionRequest struct {
	Action string `json:"action"`
}

func New(store *database.Store, startedAt time.Time) *Server {
	return &Server{
		Store:     store,
		StartedAt: startedAt,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/status", s.HandleStatus)
	mux.HandleFunc("/health", s.HandleStatus)
	mux.HandleFunc("/api/dev/seed-player", s.HandleSeedDevPlayer)
	mux.HandleFunc("/api/dev/force-tick", s.HandleForceTick)
	mux.HandleFunc("/api/dev/rewind-care-decay", s.HandleRewindCareDecay)
	mux.HandleFunc("/api/player/status", s.HandlePlayerStatus)
	mux.HandleFunc("/api/player/settle-ticks", s.HandleSettleTicks)
	mux.HandleFunc("/api/player/gathering", s.HandleGatheringTask)
	mux.HandleFunc("/api/player/care", s.HandleCareAction)
	mux.HandleFunc("/api/nami/messages", s.HandleNamiMessages)

	mux.Handle("/", http.FileServer(http.Dir("web")))

	return mux
}

func (s *Server) HandleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	response := StatusResponse{
		Server:      "online",
		Database:    s.databaseStatus(),
		Timestamp:   time.Now().Format(time.RFC3339),
		Uptime:      time.Since(s.StartedAt).Round(time.Second).String(),
		Version:     "phase-3b",
		OnlineUsers: 1,
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) HandleSeedDevPlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := s.Store.SeedDevPlayer(r.Context()); err != nil {
		log.Printf("seed dev player failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to seed dev player")
		return
	}

	writeJSON(w, http.StatusOK, MessageResponse{
		OK:      true,
		Message: "Dev player Soryn and Nami-chan are ready.",
	})
}

func (s *Server) HandlePlayerStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	_, _ = s.Store.SettleDevTicks(r.Context(), 0)

	if err := s.Store.SettleDevCareActions(r.Context()); err != nil {
		log.Printf("settle dev care actions failed: %v", err)
	}

	if err := s.Store.SettleDevCareDecay(r.Context()); err != nil {
		log.Printf("settle dev care decay failed: %v", err)
	}

	if err := s.Store.GenerateDevPassiveNamiMessages(r.Context()); err != nil {
		log.Printf("generate passive nami messages failed: %v", err)
	}

	status, err := s.Store.GetDevPlayerStatus(r.Context())
	if err != nil {
		log.Printf("get player status failed: %v", err)
		writeError(w, http.StatusNotFound, "player status not found; visit /api/dev/seed-player first")
		return
	}

	writeJSON(w, http.StatusOK, status)
}

func (s *Server) HandleSettleTicks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	result, err := s.Store.SettleDevTicks(r.Context(), 0)
	if err != nil {
		log.Printf("settle ticks failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to settle ticks")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) HandleForceTick(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	ticks := int64(1)

	if rawTicks := r.URL.Query().Get("ticks"); rawTicks != "" {
		parsedTicks, err := strconv.ParseInt(rawTicks, 10, 64)
		if err != nil || parsedTicks < 1 {
			writeError(w, http.StatusBadRequest, "ticks must be a positive whole number")
			return
		}

		ticks = parsedTicks
	}

	result, err := s.Store.SettleDevTicks(r.Context(), ticks)
	if err != nil {
		log.Printf("force tick failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to force tick")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) HandleRewindCareDecay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	rawHours := strings.TrimSpace(r.URL.Query().Get("hours"))
	if rawHours == "" {
		rawHours = "2"
	}

	hours, err := strconv.ParseFloat(rawHours, 64)
	if err != nil || hours <= 0 {
		writeError(w, http.StatusBadRequest, "hours must be a positive number")
		return
	}

	duration := time.Duration(hours * float64(time.Hour))

	if err := s.Store.RewindDevCareDecay(r.Context(), duration); err != nil {
		log.Printf("rewind care decay failed: %v", err)
		writeError(w, http.StatusInternalServerError, "rewind care decay failed")
		return
	}

	if err := s.Store.SettleDevCareDecay(r.Context()); err != nil {
		log.Printf("settle dev care decay failed: %v", err)
		writeError(w, http.StatusInternalServerError, "settle dev care decay failed")
		return
	}

	status, err := s.Store.GetDevPlayerStatus(r.Context())
	if err != nil {
		log.Printf("get player status after care decay rewind failed: %v", err)
		writeError(w, http.StatusInternalServerError, "player status failed")
		return
	}

	writeJSON(w, http.StatusOK, status)
}

func (s *Server) HandleGatheringTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var request GatheringRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid gathering request")
		return
	}

	request.Task = strings.TrimSpace(strings.ToLower(request.Task))
	if err := s.Store.SetDevGatheringTask(r.Context(), request.Task); err != nil {
		log.Printf("set gathering task failed: %v", err)
		writeError(w, http.StatusBadRequest, "invalid gathering task")
		return
	}

	writeJSON(w, http.StatusOK, MessageResponse{
		OK:      true,
		Message: "Gathering task updated.",
	})
}

func (s *Server) HandleCareAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var request CareActionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid care action request")
		return
	}

	request.Action = strings.TrimSpace(strings.ToLower(request.Action))
	result, err := s.Store.StartOrQueueDevCareAction(r.Context(), request.Action)
	if err != nil {
		log.Printf("care action failed: %v", err)
		writeError(w, http.StatusBadRequest, "invalid care action")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) HandleNamiMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	messages, err := s.Store.GetRecentDevNamiMessages(r.Context(), 100)
	if err != nil {
		log.Printf("get nami messages failed: %v", err)
		writeError(w, http.StatusNotFound, "nami messages not found; visit /api/dev/seed-player first")
		return
	}

	writeJSON(w, http.StatusOK, messages)
}

func (s *Server) databaseStatus() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := s.Store.Pool.Ping(ctx); err != nil {
		log.Printf("database status check failed: %v", err)
		return "offline"
	}

	return "online"
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to write JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, MessageResponse{
		OK:      false,
		Message: message,
	})
}
