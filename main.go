package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SorynMochi/Namagotchi/internal/database"
	"github.com/SorynMochi/Namagotchi/internal/server"
)

const listenAddress = ":8080"

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	ctx := context.Background()

	store, err := database.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatal("database connection failed:", err)
	}
	defer store.Close()

	if err := store.RunMigrations(ctx, "migrations"); err != nil {
		log.Fatal("database migrations failed:", err)
	}

	app := server.New(store, time.Now())

	log.Printf("Namagotchi Phase 2 server running at http://localhost%s", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, app.Routes()))
}
