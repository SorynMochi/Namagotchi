package database

import (
"context"
"errors"
"fmt"

"github.com/jackc/pgx/v5"
)

func (s *Store) PlayerIDForAccount(ctx context.Context, accountID int64) (int64, error) {
if accountID < 1 {
return 0, fmt.Errorf("account id must be positive")
}

var playerID int64

err := s.Pool.QueryRow(ctx, `
select id
from players
where account_id = $1
`, accountID).Scan(&playerID)
if err == nil {
return playerID, nil
}
if err != nil && !errors.Is(err, pgx.ErrNoRows) {
return 0, fmt.Errorf("get player for account: %w", err)
}

err = s.Pool.QueryRow(ctx, `
update players
set account_id = $1
where id = (
select id
from players
where account_id is null
order by
case when display_name = 'Soryn' then 0 else 1 end,
id
limit 1
)
returning id
`, accountID).Scan(&playerID)
if err == nil {
return playerID, nil
}
if errors.Is(err, pgx.ErrNoRows) {
return 0, fmt.Errorf("no unclaimed gameplay player is available for this account")
}

return 0, fmt.Errorf("claim player for account: %w", err)
}
