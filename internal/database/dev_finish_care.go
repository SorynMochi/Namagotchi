package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s *Store) FinishActiveDevCareAction(ctx context.Context) (CareQueueState, string, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return CareQueueState{}, "", fmt.Errorf("begin finish active care action: %w", err)
	}
	defer tx.Rollback(ctx)

	playerID, _, err := loadDevCompanionForUpdateTx(ctx, tx)
	if err != nil {
		return CareQueueState{}, "", err
	}

	active, ok, err := loadActiveCareActionForFinishTx(ctx, tx, playerID)
	if err != nil {
		return CareQueueState{}, "", err
	}

	if !ok {
		state, err := loadCareQueueStateTx(ctx, tx, playerID)
		if err != nil {
			return CareQueueState{}, "", err
		}

		return state, "No active care action to finish.", nil
	}

	commandTag, err := tx.Exec(ctx, `
        update companion_care_actions
        set completes_at = now(),
            updated_at = now()
        where id = $1
            and player_id = $2
            and status = 'active'
    `, active.ID, playerID)
	if err != nil {
		return CareQueueState{}, "", fmt.Errorf("mark active care action complete: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return CareQueueState{}, "", fmt.Errorf("active care action not found")
	}

	if err := settleCompletedCareActionsTx(ctx, tx, playerID); err != nil {
		return CareQueueState{}, "", err
	}

	state, err := loadCareQueueStateTx(ctx, tx, playerID)
	if err != nil {
		return CareQueueState{}, "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return CareQueueState{}, "", fmt.Errorf("commit finish active care action: %w", err)
	}

	return state, fmt.Sprintf("Finished active care action: %s.", active.ActionName), nil
}

func loadActiveCareActionForFinishTx(ctx context.Context, tx pgx.Tx, playerID int64) (CareActionState, bool, error) {
	var (
		id              int64
		actionKey       string
		actionName      string
		status          string
		queuePosition   int
		durationSeconds int
		startedAt       time.Time
		completesAt     time.Time
		completedAt     time.Time
		createdAt       time.Time
		updatedAt       time.Time
	)

	err := tx.QueryRow(ctx, `
        select
            id,
            action_key,
            action_name,
            status,
            coalesce(queue_position, 0),
            duration_seconds,
            coalesce(started_at, '0001-01-01 00:00:00+00'::timestamptz),
            coalesce(completes_at, '0001-01-01 00:00:00+00'::timestamptz),
            coalesce(completed_at, '0001-01-01 00:00:00+00'::timestamptz),
            created_at,
            updated_at
        from companion_care_actions
        where player_id = $1
            and status = 'active'
        order by started_at, id
        limit 1
        for update
    `, playerID).Scan(
		&id,
		&actionKey,
		&actionName,
		&status,
		&queuePosition,
		&durationSeconds,
		&startedAt,
		&completesAt,
		&completedAt,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return CareActionState{}, false, nil
		}

		return CareActionState{}, false, fmt.Errorf("load active care action for finish: %w", err)
	}

	return careActionFromRow(
		id,
		actionKey,
		actionName,
		status,
		queuePosition,
		durationSeconds,
		startedAt,
		completesAt,
		completedAt,
		createdAt,
		updatedAt,
	), true, nil
}
