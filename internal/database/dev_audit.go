package database

import (
	"context"
	"fmt"
	"time"
)

type DevAuditLog struct {
	ID         int64     `json:"id"`
	AccountID  int64     `json:"accountId"`
	Command    string    `json:"command"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	StatusCode int       `json:"statusCode"`
	RemoteAddr string    `json:"remoteAddr"`
	UserAgent  string    `json:"userAgent"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (s *Store) RecordDevAuditLog(ctx context.Context, entry DevAuditLog) error {
	_, err := s.Pool.Exec(ctx, `
insert into dev_audit_logs (
account_id,
command,
method,
path,
status_code,
remote_addr,
user_agent
)
values ($1, $2, $3, $4, $5, $6, $7)
`,
		entry.AccountID,
		entry.Command,
		entry.Method,
		entry.Path,
		entry.StatusCode,
		entry.RemoteAddr,
		entry.UserAgent,
	)

	if err != nil {
		return fmt.Errorf("record dev audit log: %w", err)
	}

	return nil
}

func (s *Store) RecentDevAuditLogs(ctx context.Context, limit int) ([]DevAuditLog, error) {
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
command,
method,
path,
status_code,
remote_addr,
user_agent,
created_at
from dev_audit_logs
where command <> 'audit-logs'
order by created_at desc, id desc
limit $1
`, limit)
	if err != nil {
		return nil, fmt.Errorf("get recent dev audit logs: %w", err)
	}
	defer rows.Close()

	logs := make([]DevAuditLog, 0, limit)

	for rows.Next() {
		var entry DevAuditLog

		if err := rows.Scan(
			&entry.ID,
			&entry.AccountID,
			&entry.Command,
			&entry.Method,
			&entry.Path,
			&entry.StatusCode,
			&entry.RemoteAddr,
			&entry.UserAgent,
			&entry.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan dev audit log: %w", err)
		}

		logs = append(logs, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate dev audit logs: %w", err)
	}

	return logs, nil
}
