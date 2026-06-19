package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *Store) SetAuthDisplayNameForAccountAllowReserved(ctx context.Context, accountID int64, displayName string, allowReserved bool) (AuthAccount, error) {
	var account AuthAccount

	if accountID < 1 {
		return account, fmt.Errorf("account id must be positive")
	}

	displayName = cleanAuthDisplayName(displayName)
	if err := ValidateAuthDisplayNameSyntax(displayName); err != nil {
		return account, err
	}

	if IsReservedAuthDisplayName(displayName) && !allowReserved {
		return account, ErrAuthDisplayNameReserved
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return account, err
	}
	defer tx.Rollback(ctx)

	var existingID int64
	err = tx.QueryRow(ctx, `
select id
from auth_accounts
where lower(display_name) = lower($1)
and id <> $2
`, displayName, accountID).Scan(&existingID)
	if err == nil {
		return account, ErrAuthDisplayNameTaken
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return account, err
	}

	err = tx.QueryRow(ctx, `
select id
from players
where account_id is not null
and account_id <> $2
and lower(display_name) = lower($1)
`, displayName, accountID).Scan(&existingID)
	if err == nil {
		return account, ErrAuthDisplayNameTaken
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return account, err
	}

	var playerID int64
	err = tx.QueryRow(ctx, `
select id
from players
where account_id = $1
`, accountID).Scan(&playerID)
	if errors.Is(err, pgx.ErrNoRows) {
		err = tx.QueryRow(ctx, `
update players
set account_id = $1,
display_name = $2
where id = (
select id
from players
where account_id is null
order by id
limit 1
)
returning id
`, accountID, displayName).Scan(&playerID)
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return account, fmt.Errorf("no unclaimed gameplay player is available for this account")
		}

		return account, err
	}

	_, err = tx.Exec(ctx, `
update players
set display_name = $2
where id = $1
`, playerID, displayName)
	if err != nil {
		return account, err
	}

	err = tx.QueryRow(ctx, `
update auth_accounts
set display_name = $2,
updated_at = now()
where id = $1
returning id, display_name, email, avatar_url, created_at, updated_at, last_login_at
`, accountID, displayName).Scan(
		&account.ID,
		&account.DisplayName,
		&account.Email,
		&account.AvatarURL,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.LastLoginAt,
	)
	if err != nil {
		return account, err
	}

	_, err = tx.Exec(ctx, `
update auth_identities
set display_name = $2
where account_id = $1
`, accountID, displayName)
	if err != nil {
		return account, err
	}

	if err := tx.Commit(ctx); err != nil {
		return account, err
	}

	return account, nil
}
