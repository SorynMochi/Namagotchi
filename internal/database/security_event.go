package database

import (
	"context"
	"fmt"
)

type SecurityEventLog struct {
	AccountID  int64
	EventType  string
	Reason     string
	Method     string
	Path       string
	StatusCode int
	RemoteAddr string
	UserAgent  string
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
