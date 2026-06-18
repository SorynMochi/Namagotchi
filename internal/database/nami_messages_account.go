package database

import (
	"context"
	"fmt"
)

func (s *Store) GetRecentNamiMessagesForAccount(ctx context.Context, accountID int64, limit int) ([]NamiMessage, error) {
	if accountID < 1 {
		return nil, fmt.Errorf("account id must be positive")
	}

	playerID, err := s.PlayerIDForAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return s.GetRecentNamiMessages(ctx, playerID, limit)
}
