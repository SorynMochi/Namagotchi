package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

const listenAddress = ":8080"

type statusResponse struct {
	Server    string `json:"server"`
	Database  string `json:"database"`
	Timestamp string `json:"timestamp"`
	Uptime    string `json:"uptime"`
	Version   string `json:"version"`
}

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	startedAt := time.Now()
	mux := http.NewServeMux()

	mux.HandleFunc("/api/status", statusHandler(databaseURL, startedAt))
	mux.HandleFunc("/health", statusHandler(databaseURL, startedAt))
	mux.Handle("/", http.FileServer(http.Dir("web")))

	log.Printf("Namagotchi Phase 1 server running at http://localhost%s", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, mux))
}

func statusHandler(databaseURL string, startedAt time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		response := statusResponse{
			Server:    "online",
			Database:  databaseStatus(databaseURL),
			Timestamp: time.Now().Format(time.RFC3339),
			Uptime:    time.Since(startedAt).Round(time.Second).String(),
			Version:   "phase-1",
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("failed to write status response: %v", err)
		}
	}
}

func databaseStatus(databaseURL string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Printf("database status check failed to connect: %v", err)
		return "offline"
	}
	defer conn.Close(context.Background())

	var ok int
	if err := conn.QueryRow(ctx, "select 1").Scan(&ok); err != nil {
		log.Printf("database status check failed: %v", err)
		return "offline"
	}

	return "online"
}
