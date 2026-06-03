package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
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
	mux.HandleFunc("/api/player/status", s.HandlePlayerStatus)

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
		Version:     "phase-3a",
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

	status, err := s.Store.GetDevPlayerStatus(r.Context())
	if err != nil {
		log.Printf("get player status failed: %v", err)
		writeError(w, http.StatusNotFound, "player status not found; visit /api/dev/seed-player first")
		return
	}

	writeJSON(w, http.StatusOK, status)
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
