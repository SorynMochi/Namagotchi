package database

import (
	"context"
	"fmt"
)

const (
	PlaydeckCombatLogStorageLimit = 50
	SecurityEventLogStorageLimit  = 100
	DevAuditLogStorageLimit       = 100
)

type SoftLogPruneResult struct {
	PlaydeckCombatLogs int64
	NamiMessages       int64
	SecurityEventLogs  int64
	DevAuditLogs       int64
}

func (r SoftLogPruneResult) TotalPruned() int64 {
	return r.PlaydeckCombatLogs + r.NamiMessages + r.SecurityEventLogs + r.DevAuditLogs
}

func (s *Store) PruneSoftLogs(ctx context.Context) (SoftLogPruneResult, error) {
	var result SoftLogPruneResult
	var err error

	result.PlaydeckCombatLogs, err = s.pruneNewestPerPlayer(ctx, "playdeck_combat_log", "player_id", PlaydeckCombatLogStorageLimit)
	if err != nil {
		return result, fmt.Errorf("prune playdeck combat log: %w", err)
	}

	result.NamiMessages, err = s.pruneNewestPerPlayer(ctx, "nami_messages", "player_id", NamiMessageStorageLimit)
	if err != nil {
		return result, fmt.Errorf("prune nami messages: %w", err)
	}

	result.SecurityEventLogs, err = s.pruneNewestGlobal(ctx, "security_event_logs", SecurityEventLogStorageLimit)
	if err != nil {
		return result, fmt.Errorf("prune security event logs: %w", err)
	}

	result.DevAuditLogs, err = s.pruneNewestGlobal(ctx, "dev_audit_logs", DevAuditLogStorageLimit)
	if err != nil {
		return result, fmt.Errorf("prune dev audit logs: %w", err)
	}

	return result, nil
}

func (s *Store) pruneNewestPerPlayer(ctx context.Context, table string, partitionColumn string, limit int) (int64, error) {
	if limit < 1 {
		return 0, nil
	}

	tableName, err := retentionIdentifier(table)
	if err != nil {
		return 0, err
	}

	partitionName, err := retentionIdentifier(partitionColumn)
	if err != nil {
		return 0, err
	}

	query := fmt.Sprintf(`
with ranked_rows as (
select
id,
row_number() over (
partition by %s
order by created_at desc, id desc
) as row_rank
from %s
)
delete from %s target
using ranked_rows ranked
where target.id = ranked.id
and ranked.row_rank > $1
`, partitionName, tableName, tableName)

	commandTag, err := s.Pool.Exec(ctx, query, limit)
	if err != nil {
		return 0, err
	}

	return commandTag.RowsAffected(), nil
}

func (s *Store) pruneNewestGlobal(ctx context.Context, table string, limit int) (int64, error) {
	if limit < 1 {
		return 0, nil
	}

	tableName, err := retentionIdentifier(table)
	if err != nil {
		return 0, err
	}

	query := fmt.Sprintf(`
with ranked_rows as (
select
id,
row_number() over (
order by created_at desc, id desc
) as row_rank
from %s
)
delete from %s target
using ranked_rows ranked
where target.id = ranked.id
and ranked.row_rank > $1
`, tableName, tableName)

	commandTag, err := s.Pool.Exec(ctx, query, limit)
	if err != nil {
		return 0, err
	}

	return commandTag.RowsAffected(), nil
}

func retentionIdentifier(value string) (string, error) {
	switch value {
	case "playdeck_combat_log",
		"nami_messages",
		"security_event_logs",
		"dev_audit_logs",
		"player_id",
		"created_at",
		"id":
		return value, nil
	default:
		return "", fmt.Errorf("unsupported retention identifier %q", value)
	}
}
