package database

import (
	"context"
	"fmt"
	"strings"
)

func (s *Store) SetGatheringTaskForAccount(ctx context.Context, accountID int64, task string) error {
	if accountID < 1 {
		return fmt.Errorf("account id must be positive")
	}

	task = strings.TrimSpace(strings.ToLower(task))
	if !ValidGatheringTask(task) {
		return fmt.Errorf("invalid gathering task: %s", task)
	}

	playerID, err := s.PlayerIDForAccount(ctx, accountID)
	if err != nil {
		return err
	}

	if _, err := s.SettleDevTicks(ctx, 0); err != nil {
		return fmt.Errorf("settle ticks before gathering switch: %w", err)
	}

	commandTag, err := s.Pool.Exec(ctx, `
update player_tick_state
set active_gathering_task = $1,
updated_at = now()
where player_id = $2
`, task, playerID)
	if err != nil {
		return fmt.Errorf("set gathering task for account: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("player tick state not found")
	}

	return nil
}
