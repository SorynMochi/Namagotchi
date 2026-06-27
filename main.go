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

	go runSoftCleanupJob(ctx, store)

	app := server.New(store, time.Now())

	log.Printf("Namigotchi Idle server running at http://localhost%s", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, app.Routes()))
}

const softCleanupInterval = time.Hour

func runSoftCleanupJob(ctx context.Context, store *database.Store) {
	if store == nil {
		return
	}

	runCleanup := func() {
		cleanupCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		result, err := store.PruneSoftLogs(cleanupCtx)
		if err != nil {
			log.Printf("soft cleanup failed: %v", err)
			return
		}

		if result.TotalPruned() > 0 {
			log.Printf(
				"soft cleanup pruned playdeck=%d nami_messages=%d security=%d dev_audit=%d",
				result.PlaydeckCombatLogs,
				result.NamiMessages,
				result.SecurityEventLogs,
				result.DevAuditLogs,
			)
		}
	}

	runCleanup()

	ticker := time.NewTicker(softCleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			runCleanup()
		}
	}
}
