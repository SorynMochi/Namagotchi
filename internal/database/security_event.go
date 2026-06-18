package database

import (
	"context"
	"fmt"
	"time"
)

type SecurityEventLog struct {
	ID         int64     `json:"id"`
	AccountID  int64     `json:"accountId"`
	EventType  string    `json:"eventType"`
	Reason     string    `json:"reason"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	StatusCode int       `json:"statusCode"`
	RemoteAddr string    `json:"remoteAddr"`
	UserAgent  string    `json:"userAgent"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (s *Store) RecordSecurityEventLog(ctx context.Context, entry SecurityEventLog) error {
	var accountID any
	if entry.AccountID > 0 {
		accountID = entry.AccountID
	}

	_, err := s.Pool.Exec(ctx, `
insert into security_event_logs (
account_id,
event_type,
reason,
method,
path,
status_code,
remote_addr,
user_agent
)
values ($1, $2, $3, $4, $5, $6, $7, $8)
`,
		accountID,
		entry.EventType,
		entry.Reason,
		entry.Method,
		entry.Path,
		entry.StatusCode,
		entry.RemoteAddr,
		entry.UserAgent,
	)

	if err != nil {
		return fmt.Errorf("record security event log: %w", err)
	}

	return nil
}

func (s *Store) RecentSecurityEventLogs(ctx context.Context, limit int) ([]SecurityEventLog, error) {
	if limit < 1 {
		limit = 1
	}

	if limit > 200 {
		limit = 200
	}

	rows, err := s.Pool.Query(ctx, `
select
id,
coalesce(account_id, 0),
event_type,
reason,
method,
path,
status_code,
remote_addr,
user_agent,
created_at
from security_event_logs
order by created_at desc, id desc
limit $1
`, limit)
	if err != nil {
		return nil, fmt.Errorf("get recent security event logs: %w", err)
	}
	defer rows.Close()

	logs := make([]SecurityEventLog, 0, limit)

	for rows.Next() {
		var entry SecurityEventLog

		if err := rows.Scan(
			&entry.ID,
			&entry.AccountID,
			&entry.EventType,
			&entry.Reason,
			&entry.Method,
			&entry.Path,
			&entry.StatusCode,
			&entry.RemoteAddr,
			&entry.UserAgent,
			&entry.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan security event log: %w", err)
		}

		logs = append(logs, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate security event logs: %w", err)
	}

	return logs, nil
}
