package database

import (
	"context"
	"fmt"
	"strings"
)

type DevTargetSummary struct {
	Mode        string   `json:"mode"`
	Requested   string   `json:"requested"`
	PlayerCount int      `json:"playerCount"`
	PlayerNames []string `json:"playerNames"`
}

type DevChainResetPlayer struct {
	PlayerID           int64  `json:"playerId"`
	PlayerName         string `json:"playerName"`
	BeforeCurrentChain int64  `json:"beforeCurrentChain"`
	BeforeMaxChain     int64  `json:"beforeMaxChain"`
	AfterCurrentChain  int64  `json:"afterCurrentChain"`
	AfterMaxChain      int64  `json:"afterMaxChain"`
}

type DevResetChainResult struct {
	OK       bool                  `json:"ok"`
	Message  string                `json:"message"`
	ResetMax bool                  `json:"resetMax"`
	Target   DevTargetSummary      `json:"target"`
	Players  []DevChainResetPlayer `json:"players"`
}

type DevClearWardrobePlayer struct {
	PlayerID     int64  `json:"playerId"`
	PlayerName   string `json:"playerName"`
	DeletedItems int64  `json:"deletedItems"`
}

type DevClearWardrobeResult struct {
	OK           bool                     `json:"ok"`
	Message      string                   `json:"message"`
	Target       DevTargetSummary         `json:"target"`
	DeletedItems int64                    `json:"deletedItems"`
	Players      []DevClearWardrobePlayer `json:"players"`
}

type devTargetPlayer struct {
	ID   int64
	Name string
}

func (s *Store) resolveDevTargetPlayers(ctx context.Context, rawPlayerName string, allToken string) (DevTargetSummary, []devTargetPlayer, error) {
	requested := strings.TrimSpace(rawPlayerName)
	if requested == "" {
		return DevTargetSummary{}, nil, fmt.Errorf("playerName is required")
	}

	allToken = strings.TrimSpace(strings.ToUpper(allToken))
	allMode := strings.EqualFold(requested, allToken)

	query := `
select id, display_name
from players
`
	args := []any{}

	if !allMode {
		query += ` where display_name = $1`
		args = append(args, requested)
	}

	query += ` order by display_name, id`

	rows, err := s.Pool.Query(ctx, query, args...)
	if err != nil {
		return DevTargetSummary{}, nil, fmt.Errorf("resolve target players: %w", err)
	}
	defer rows.Close()

	players := []devTargetPlayer{}
	names := []string{}

	for rows.Next() {
		var player devTargetPlayer
		if err := rows.Scan(&player.ID, &player.Name); err != nil {
			return DevTargetSummary{}, nil, fmt.Errorf("scan target player: %w", err)
		}

		players = append(players, player)
		names = append(names, player.Name)
	}

	if err := rows.Err(); err != nil {
		return DevTargetSummary{}, nil, fmt.Errorf("iterate target players: %w", err)
	}

	if len(players) == 0 {
		if allMode {
			return DevTargetSummary{}, nil, fmt.Errorf("no players found")
		}

		return DevTargetSummary{}, nil, fmt.Errorf("player %q not found; exact display_name match required", requested)
	}

	mode := "single"
	if allMode {
		mode = "all"
	}

	return DevTargetSummary{
		Mode:        mode,
		Requested:   requested,
		PlayerCount: len(players),
		PlayerNames: names,
	}, players, nil
}

func (s *Store) DevResetChain(ctx context.Context, playerName string, resetMax bool) (DevResetChainResult, error) {
	target, players, err := s.resolveDevTargetPlayers(ctx, playerName, "RESETALL")
	if err != nil {
		return DevResetChainResult{}, err
	}

	if err := s.ensurePlaydeckZoneRecordsTable(ctx); err != nil {
		return DevResetChainResult{}, err
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return DevResetChainResult{}, fmt.Errorf("begin reset chain: %w", err)
	}
	defer tx.Rollback(ctx)

	result := DevResetChainResult{
		OK:       true,
		ResetMax: resetMax,
		Target:   target,
		Players:  make([]DevChainResetPlayer, 0, len(players)),
	}

	for _, player := range players {
		entry := DevChainResetPlayer{
			PlayerID:   player.ID,
			PlayerName: player.Name,
		}

		if _, err := tx.Exec(ctx, `
insert into player_tick_state (player_id)
values ($1)
on conflict (player_id) do nothing
`, player.ID); err != nil {
			return DevResetChainResult{}, fmt.Errorf("ensure tick state for %s: %w", player.Name, err)
		}

		if err := tx.QueryRow(ctx, `
select coalesce(playdeck_streak, 0)
from player_tick_state
where player_id = $1
for update
`, player.ID).Scan(&entry.BeforeCurrentChain); err != nil {
			return DevResetChainResult{}, fmt.Errorf("load current chain for %s: %w", player.Name, err)
		}

		if err := tx.QueryRow(ctx, `
select coalesce(max(max_streak), 0)
from player_playdeck_zone_records
where player_id = $1
`, player.ID).Scan(&entry.BeforeMaxChain); err != nil {
			return DevResetChainResult{}, fmt.Errorf("load max chain for %s: %w", player.Name, err)
		}

		if _, err := tx.Exec(ctx, `
update player_tick_state
set playdeck_streak = 0,
last_tick_at = now(),
updated_at = now()
where player_id = $1
`, player.ID); err != nil {
			return DevResetChainResult{}, fmt.Errorf("reset current chain for %s: %w", player.Name, err)
		}

		entry.AfterCurrentChain = 0
		entry.AfterMaxChain = entry.BeforeMaxChain

		if resetMax {
			if _, err := tx.Exec(ctx, `
update player_playdeck_zone_records
set max_streak = 0,
updated_at = now()
where player_id = $1
`, player.ID); err != nil {
				return DevResetChainResult{}, fmt.Errorf("reset max chain for %s: %w", player.Name, err)
			}

			entry.AfterMaxChain = 0
		}

		result.Players = append(result.Players, entry)
	}

	if err := tx.Commit(ctx); err != nil {
		return DevResetChainResult{}, fmt.Errorf("commit reset chain: %w", err)
	}

	if resetMax {
		result.Message = fmt.Sprintf("Reset current and max chain for %d player(s).", len(players))
	} else {
		result.Message = fmt.Sprintf("Reset current chain for %d player(s).", len(players))
	}

	return result, nil
}

func (s *Store) DevClearWardrobe(ctx context.Context, playerName string) (DevClearWardrobeResult, error) {
	target, players, err := s.resolveDevTargetPlayers(ctx, playerName, "RESETALL")
	if err != nil {
		return DevClearWardrobeResult{}, err
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return DevClearWardrobeResult{}, fmt.Errorf("begin clear wardrobe: %w", err)
	}
	defer tx.Rollback(ctx)

	result := DevClearWardrobeResult{
		OK:      true,
		Target:  target,
		Players: make([]DevClearWardrobePlayer, 0, len(players)),
	}

	for _, player := range players {
		commandTag, err := tx.Exec(ctx, `
delete from player_inventory_items
where player_id = $1
and container = 'wardrobe'
`, player.ID)
		if err != nil {
			return DevClearWardrobeResult{}, fmt.Errorf("clear wardrobe for %s: %w", player.Name, err)
		}

		deleted := commandTag.RowsAffected()
		result.DeletedItems += deleted
		result.Players = append(result.Players, DevClearWardrobePlayer{
			PlayerID:     player.ID,
			PlayerName:   player.Name,
			DeletedItems: deleted,
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return DevClearWardrobeResult{}, fmt.Errorf("commit clear wardrobe: %w", err)
	}

	result.Message = fmt.Sprintf("Deleted %d wardrobe item(s) for %d player(s).", result.DeletedItems, len(players))

	return result, nil
}
