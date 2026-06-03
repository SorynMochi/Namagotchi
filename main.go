package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatal("Could not connect to PostgreSQL:", err)
	}
	defer conn.Close(ctx)

	var now string
	err = conn.QueryRow(ctx, "select now()::text").Scan(&now)
	if err != nil {
		log.Fatal("Database query failed:", err)
	}

	fmt.Println("Connected to PostgreSQL!")
	fmt.Println("Database time:", now)
}
