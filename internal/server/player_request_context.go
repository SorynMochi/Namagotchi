package server

import (
	"context"
	"log"

	"github.com/SorynMochi/Namagotchi/internal/database"
)

func (s *Server) withRequestPlayerID(ctx context.Context, accountID int64) context.Context {
	if accountID < 1 {
		return ctx
	}

	if _, ok := database.AuthPlayerIDFromContext(ctx); ok {
		return ctx
	}

	playerID, err := s.Store.PlayerIDForAccount(ctx, accountID)
	if err != nil {
		log.Printf("get account player id for request context failed: %v", err)
		return ctx
	}

	return database.WithAuthPlayerID(ctx, playerID)
}
