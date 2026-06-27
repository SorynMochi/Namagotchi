package database

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const devAdminMaxInt64 = int64(^uint64(0) >> 1)

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

type DevCurrencyPlayer struct {
	PlayerID       int64            `json:"playerId"`
	PlayerName     string           `json:"playerName"`
	Before         map[string]int64 `json:"before"`
	ChangedBy      map[string]int64 `json:"changedBy"`
	After          map[string]int64 `json:"after"`
	CreditsBefore  int64            `json:"creditsBefore"`
	CreditsAfter   int64            `json:"creditsAfter"`
	NibblesBefore  int64            `json:"nibblesBefore"`
	NibblesAfter   int64            `json:"nibblesAfter"`
	NamiCoinBefore int64            `json:"namiCoinBefore"`
	NamiCoinAfter  int64            `json:"namiCoinAfter"`
}

type DevCurrencyResult struct {
	OK           bool                `json:"ok"`
	Message      string              `json:"message"`
	Operation    string              `json:"operation"`
	Currency     string              `json:"currency"`
	AmountInput  string              `json:"amountInput"`
	Amount       int64               `json:"amount"`
	AmountWasAll bool                `json:"amountWasAll"`
	Target       DevTargetSummary    `json:"target"`
	Players      []DevCurrencyPlayer `json:"players"`
}

type DevLevelResetPlayer struct {
	PlayerID   int64                `json:"playerId"`
	PlayerName string               `json:"playerName"`
	Resets     []DevLevelResetEntry `json:"resets"`
}

type DevLevelResetEntry struct {
	Key               string `json:"key"`
	Label             string `json:"label"`
	BeforeLevel       int    `json:"beforeLevel"`
	BeforeTotalXP     int64  `json:"beforeTotalXp"`
	BeforeXPIntoLevel int64  `json:"beforeXpIntoLevel"`
	AfterLevel        int    `json:"afterLevel"`
	AfterTotalXP      int64  `json:"afterTotalXp"`
	AfterXPIntoLevel  int64  `json:"afterXpIntoLevel"`
}

type DevResetLevelsResult struct {
	OK            bool                  `json:"ok"`
	Message       string                `json:"message"`
	ActivityInput string                `json:"activityInput"`
	Target        DevTargetSummary      `json:"target"`
	Players       []DevLevelResetPlayer `json:"players"`
}

type devTargetPlayer struct {
	ID   int64
	Name string
}

type devLevelResetSpec struct {
	Key   string
	Label string
	Kind  string
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

func parseDevAdminAmount(raw string) (int64, bool, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return 0, false, fmt.Errorf("amount is required")
	}

	if strings.EqualFold(value, "ALL") {
		return 0, true, nil
	}

	value = strings.ReplaceAll(value, ",", "")
	value = strings.ReplaceAll(value, "_", "")

	multiplier := float64(1)
	last := strings.ToLower(value[len(value)-1:])

	switch last {
	case "k":
		multiplier = 1_000
		value = strings.TrimSpace(value[:len(value)-1])
	case "m":
		multiplier = 1_000_000
		value = strings.TrimSpace(value[:len(value)-1])
	case "b":
		multiplier = 1_000_000_000
		value = strings.TrimSpace(value[:len(value)-1])
	case "t":
		multiplier = 1_000_000_000_000
		value = strings.TrimSpace(value[:len(value)-1])
	case "q":
		multiplier = 1_000_000_000_000_000
		value = strings.TrimSpace(value[:len(value)-1])
	}

	if value == "" {
		return 0, false, fmt.Errorf("amount is invalid")
	}

	number, err := strconv.ParseFloat(value, 64)
	if err != nil || math.IsNaN(number) || math.IsInf(number, 0) || number < 0 {
		return 0, false, fmt.Errorf("amount must be a non-negative number")
	}

	expanded := number * multiplier
	if expanded > float64(devAdminMaxInt64) {
		return 0, false, fmt.Errorf("amount is too large")
	}

	return int64(math.Round(expanded)), false, nil
}

type devCurrencyAsset struct {
	Key     string
	JSONKey string
	Label   string
	Kind    string
}

func devCurrencyAssets() []devCurrencyAsset {
	return []devCurrencyAsset{
		{Key: "credits", JSONKey: "credits", Label: "Credits", Kind: "player"},
		{Key: "nibbles", JSONKey: "nibbles", Label: "Nibbles", Kind: "player"},
		{Key: "namicoin", JSONKey: "namiCoins", Label: "NamiCoins", Kind: "player"},
		{Key: "fans", JSONKey: "fans", Label: "Fans", Kind: "resource"},
		{Key: "memes", JSONKey: "memes", Label: "Memes", Kind: "resource"},
		{Key: "lost_items", JSONKey: "lostItems", Label: "Lost Items", Kind: "resource"},
		{Key: "confidence", JSONKey: "confidence", Label: "Confidence", Kind: "resource"},
		{Key: "receipts", JSONKey: "receipts", Label: "Receipts", Kind: "resource"},
		{Key: "patterns", JSONKey: "patterns", Label: "Patterns", Kind: "resource"},
	}
}

func normalizeDevCurrencyType(raw string, allowAll bool) ([]devCurrencyAsset, string, error) {
	value := strings.TrimSpace(strings.ToLower(raw))
	value = strings.ReplaceAll(value, "_", "")
	value = strings.ReplaceAll(value, "-", "")
	value = strings.ReplaceAll(value, " ", "")

	allAssets := devCurrencyAssets()

	if allowAll && value == "all" {
		return allAssets, "ALL", nil
	}

	for _, asset := range allAssets {
		switch asset.Key {
		case "credits":
			if value == "credit" || value == "credits" {
				return []devCurrencyAsset{asset}, asset.Label, nil
			}
		case "nibbles":
			if value == "nibble" || value == "nibbles" {
				return []devCurrencyAsset{asset}, asset.Label, nil
			}
		case "namicoin":
			if value == "namicoin" || value == "namicoins" || value == "nami" || value == "namic" {
				return []devCurrencyAsset{asset}, asset.Label, nil
			}
		case "fans":
			if value == "fan" || value == "fans" {
				return []devCurrencyAsset{asset}, asset.Label, nil
			}
		case "memes":
			if value == "meme" || value == "memes" {
				return []devCurrencyAsset{asset}, asset.Label, nil
			}
		case "lost_items":
			if value == "lostitem" || value == "lostitems" {
				return []devCurrencyAsset{asset}, asset.Label, nil
			}
		case "confidence":
			if value == "confidence" {
				return []devCurrencyAsset{asset}, asset.Label, nil
			}
		case "receipts":
			if value == "receipt" || value == "receipts" {
				return []devCurrencyAsset{asset}, asset.Label, nil
			}
		case "patterns":
			if value == "pattern" || value == "patterns" {
				return []devCurrencyAsset{asset}, asset.Label, nil
			}
		}
	}

	if allowAll {
		return nil, "", fmt.Errorf("currencyType must be Credits, Nibbles, NamiCoins, Fans, Memes, Lost Items, Confidence, Receipts, Patterns, or ALL")
	}

	return nil, "", fmt.Errorf("currencyType must be Credits, Nibbles, NamiCoins, Fans, Memes, Lost Items, Confidence, Receipts, or Patterns")
}

func currencyMap(
	creditsCents int64,
	nibbles int64,
	namiCoin int64,
	fans int64,
	memes int64,
	lostItems int64,
	confidence int64,
	receipts int64,
	patterns int64,
) map[string]int64 {
	return map[string]int64{
		"credits":    creditsCents / 100,
		"nibbles":    nibbles,
		"namiCoins":  namiCoin,
		"fans":       fans,
		"memes":      memes,
		"lostItems":  lostItems,
		"confidence": confidence,
		"receipts":   receipts,
		"patterns":   patterns,
	}
}

func changedCurrencyMap() map[string]int64 {
	changed := map[string]int64{}

	for _, asset := range devCurrencyAssets() {
		changed[asset.JSONKey] = 0
	}

	return changed
}

func addSafeInt64(before int64, delta int64) int64 {
	if delta > 0 && before > devAdminMaxInt64-delta {
		return devAdminMaxInt64
	}

	return before + delta
}

func (s *Store) DevAddCurrency(ctx context.Context, playerName string, currencyType string, amountInput string) (DevCurrencyResult, error) {
	target, players, err := s.resolveDevTargetPlayers(ctx, playerName, "GIVEALL")
	if err != nil {
		return DevCurrencyResult{}, err
	}

	assets, currencyLabel, err := normalizeDevCurrencyType(currencyType, true)
	if err != nil {
		return DevCurrencyResult{}, err
	}

	amount, amountWasAll, err := parseDevAdminAmount(amountInput)
	if err != nil {
		return DevCurrencyResult{}, err
	}

	if amountWasAll {
		return DevCurrencyResult{}, fmt.Errorf("ALL amount is only supported for Remove Currency")
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return DevCurrencyResult{}, fmt.Errorf("begin add currency: %w", err)
	}
	defer tx.Rollback(ctx)

	result := DevCurrencyResult{
		OK:          true,
		Operation:   "add",
		Currency:    currencyLabel,
		AmountInput: strings.TrimSpace(amountInput),
		Amount:      amount,
		Target:      target,
		Players:     make([]DevCurrencyPlayer, 0, len(players)),
	}

	for _, player := range players {
		entry := DevCurrencyPlayer{
			PlayerID:   player.ID,
			PlayerName: player.Name,
			ChangedBy:  changedCurrencyMap(),
		}

		if _, err := tx.Exec(ctx, `
insert into player_resources (player_id)
values ($1)
on conflict (player_id) do nothing
`, player.ID); err != nil {
			return DevCurrencyResult{}, fmt.Errorf("ensure resources for %s: %w", player.Name, err)
		}

		if err := tx.QueryRow(ctx, `
select currency_cents, nibbles, namicoin
from players
where id = $1
for update
`, player.ID).Scan(&entry.CreditsBefore, &entry.NibblesBefore, &entry.NamiCoinBefore); err != nil {
			return DevCurrencyResult{}, fmt.Errorf("load currency for %s: %w", player.Name, err)
		}

		var fansBefore, memesBefore, lostItemsBefore, confidenceBefore, receiptsBefore, patternsBefore int64
		if err := tx.QueryRow(ctx, `
select fans, memes, lost_items, confidence, receipts, patterns
from player_resources
where player_id = $1
for update
`, player.ID).Scan(
			&fansBefore,
			&memesBefore,
			&lostItemsBefore,
			&confidenceBefore,
			&receiptsBefore,
			&patternsBefore,
		); err != nil {
			return DevCurrencyResult{}, fmt.Errorf("load resources for %s: %w", player.Name, err)
		}

		creditsAfter := entry.CreditsBefore
		nibblesAfter := entry.NibblesBefore
		namiCoinAfter := entry.NamiCoinBefore
		fansAfter := fansBefore
		memesAfter := memesBefore
		lostItemsAfter := lostItemsBefore
		confidenceAfter := confidenceBefore
		receiptsAfter := receiptsBefore
		patternsAfter := patternsBefore

		for _, asset := range assets {
			switch asset.Key {
			case "credits":
				if amount > devAdminMaxInt64/100 {
					return DevCurrencyResult{}, fmt.Errorf("credits amount is too large")
				}
				creditsAfter = addSafeInt64(creditsAfter, amount*100)
			case "nibbles":
				nibblesAfter = addSafeInt64(nibblesAfter, amount)
			case "namicoin":
				namiCoinAfter = addSafeInt64(namiCoinAfter, amount)
			case "fans":
				fansAfter = addSafeInt64(fansAfter, amount)
			case "memes":
				memesAfter = addSafeInt64(memesAfter, amount)
			case "lost_items":
				lostItemsAfter = addSafeInt64(lostItemsAfter, amount)
			case "confidence":
				confidenceAfter = addSafeInt64(confidenceAfter, amount)
			case "receipts":
				receiptsAfter = addSafeInt64(receiptsAfter, amount)
			case "patterns":
				patternsAfter = addSafeInt64(patternsAfter, amount)
			}

			entry.ChangedBy[asset.JSONKey] = amount
		}

		if err := tx.QueryRow(ctx, `
update players
set currency_cents = $2,
nibbles = $3,
namicoin = $4,
updated_at = now()
where id = $1
returning currency_cents, nibbles, namicoin
`,
			player.ID,
			creditsAfter,
			nibblesAfter,
			namiCoinAfter,
		).Scan(&entry.CreditsAfter, &entry.NibblesAfter, &entry.NamiCoinAfter); err != nil {
			return DevCurrencyResult{}, fmt.Errorf("add currency for %s: %w", player.Name, err)
		}

		if err := tx.QueryRow(ctx, `
update player_resources
set fans = $2,
memes = $3,
lost_items = $4,
confidence = $5,
receipts = $6,
patterns = $7,
updated_at = now()
where player_id = $1
returning fans, memes, lost_items, confidence, receipts, patterns
`,
			player.ID,
			fansAfter,
			memesAfter,
			lostItemsAfter,
			confidenceAfter,
			receiptsAfter,
			patternsAfter,
		).Scan(
			&fansAfter,
			&memesAfter,
			&lostItemsAfter,
			&confidenceAfter,
			&receiptsAfter,
			&patternsAfter,
		); err != nil {
			return DevCurrencyResult{}, fmt.Errorf("add resources for %s: %w", player.Name, err)
		}

		entry.Before = currencyMap(entry.CreditsBefore, entry.NibblesBefore, entry.NamiCoinBefore, fansBefore, memesBefore, lostItemsBefore, confidenceBefore, receiptsBefore, patternsBefore)
		entry.After = currencyMap(entry.CreditsAfter, entry.NibblesAfter, entry.NamiCoinAfter, fansAfter, memesAfter, lostItemsAfter, confidenceAfter, receiptsAfter, patternsAfter)
		result.Players = append(result.Players, entry)
	}

	if err := tx.Commit(ctx); err != nil {
		return DevCurrencyResult{}, fmt.Errorf("commit add currency: %w", err)
	}

	result.Message = fmt.Sprintf("Added %d to %s for %d player(s).", amount, currencyLabel, len(players))
	return result, nil
}

func (s *Store) DevRemoveCurrency(ctx context.Context, playerName string, currencyType string, amountInput string) (DevCurrencyResult, error) {
	target, players, err := s.resolveDevTargetPlayers(ctx, playerName, "RESETALL")
	if err != nil {
		return DevCurrencyResult{}, err
	}

	assets, currencyLabel, err := normalizeDevCurrencyType(currencyType, true)
	if err != nil {
		return DevCurrencyResult{}, err
	}

	amount, amountWasAll, err := parseDevAdminAmount(amountInput)
	if err != nil {
		return DevCurrencyResult{}, err
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return DevCurrencyResult{}, fmt.Errorf("begin remove currency: %w", err)
	}
	defer tx.Rollback(ctx)

	result := DevCurrencyResult{
		OK:           true,
		Operation:    "remove",
		Currency:     currencyLabel,
		AmountInput:  strings.TrimSpace(amountInput),
		Amount:       amount,
		AmountWasAll: amountWasAll,
		Target:       target,
		Players:      make([]DevCurrencyPlayer, 0, len(players)),
	}

	for _, player := range players {
		entry := DevCurrencyPlayer{
			PlayerID:   player.ID,
			PlayerName: player.Name,
			ChangedBy:  changedCurrencyMap(),
		}

		if _, err := tx.Exec(ctx, `
insert into player_resources (player_id)
values ($1)
on conflict (player_id) do nothing
`, player.ID); err != nil {
			return DevCurrencyResult{}, fmt.Errorf("ensure resources for %s: %w", player.Name, err)
		}

		if err := tx.QueryRow(ctx, `
select currency_cents, nibbles, namicoin
from players
where id = $1
for update
`, player.ID).Scan(&entry.CreditsBefore, &entry.NibblesBefore, &entry.NamiCoinBefore); err != nil {
			return DevCurrencyResult{}, fmt.Errorf("load currency for %s: %w", player.Name, err)
		}

		var fansBefore, memesBefore, lostItemsBefore, confidenceBefore, receiptsBefore, patternsBefore int64
		if err := tx.QueryRow(ctx, `
select fans, memes, lost_items, confidence, receipts, patterns
from player_resources
where player_id = $1
for update
`, player.ID).Scan(
			&fansBefore,
			&memesBefore,
			&lostItemsBefore,
			&confidenceBefore,
			&receiptsBefore,
			&patternsBefore,
		); err != nil {
			return DevCurrencyResult{}, fmt.Errorf("load resources for %s: %w", player.Name, err)
		}

		creditsAfter := entry.CreditsBefore
		nibblesAfter := entry.NibblesBefore
		namiCoinAfter := entry.NamiCoinBefore
		fansAfter := fansBefore
		memesAfter := memesBefore
		lostItemsAfter := lostItemsBefore
		confidenceAfter := confidenceBefore
		receiptsAfter := receiptsBefore
		patternsAfter := patternsBefore

		for _, asset := range assets {
			switch asset.Key {
			case "credits":
				removeCreditsCents := int64(0)
				if amountWasAll {
					removeCreditsCents = entry.CreditsBefore
				} else {
					if amount > devAdminMaxInt64/100 {
						return DevCurrencyResult{}, fmt.Errorf("credits amount is too large")
					}
					removeCreditsCents = amount * 100
					if removeCreditsCents > entry.CreditsBefore {
						removeCreditsCents = entry.CreditsBefore
					}
				}
				creditsAfter = entry.CreditsBefore - removeCreditsCents
				entry.ChangedBy[asset.JSONKey] = removeCreditsCents / 100
			case "nibbles":
				removeNibbles := amount
				if amountWasAll || removeNibbles > entry.NibblesBefore {
					removeNibbles = entry.NibblesBefore
				}
				nibblesAfter = entry.NibblesBefore - removeNibbles
				entry.ChangedBy[asset.JSONKey] = removeNibbles
			case "namicoin":
				removeNamiCoin := amount
				if amountWasAll || removeNamiCoin > entry.NamiCoinBefore {
					removeNamiCoin = entry.NamiCoinBefore
				}
				namiCoinAfter = entry.NamiCoinBefore - removeNamiCoin
				entry.ChangedBy[asset.JSONKey] = removeNamiCoin
			case "fans":
				removeFans := amount
				if amountWasAll || removeFans > fansBefore {
					removeFans = fansBefore
				}
				fansAfter = fansBefore - removeFans
				entry.ChangedBy[asset.JSONKey] = removeFans
			case "memes":
				removeMemes := amount
				if amountWasAll || removeMemes > memesBefore {
					removeMemes = memesBefore
				}
				memesAfter = memesBefore - removeMemes
				entry.ChangedBy[asset.JSONKey] = removeMemes
			case "lost_items":
				removeLostItems := amount
				if amountWasAll || removeLostItems > lostItemsBefore {
					removeLostItems = lostItemsBefore
				}
				lostItemsAfter = lostItemsBefore - removeLostItems
				entry.ChangedBy[asset.JSONKey] = removeLostItems
			case "confidence":
				removeConfidence := amount
				if amountWasAll || removeConfidence > confidenceBefore {
					removeConfidence = confidenceBefore
				}
				confidenceAfter = confidenceBefore - removeConfidence
				entry.ChangedBy[asset.JSONKey] = removeConfidence
			case "receipts":
				removeReceipts := amount
				if amountWasAll || removeReceipts > receiptsBefore {
					removeReceipts = receiptsBefore
				}
				receiptsAfter = receiptsBefore - removeReceipts
				entry.ChangedBy[asset.JSONKey] = removeReceipts
			case "patterns":
				removePatterns := amount
				if amountWasAll || removePatterns > patternsBefore {
					removePatterns = patternsBefore
				}
				patternsAfter = patternsBefore - removePatterns
				entry.ChangedBy[asset.JSONKey] = removePatterns
			}
		}

		if err := tx.QueryRow(ctx, `
update players
set currency_cents = $2,
nibbles = $3,
namicoin = $4,
updated_at = now()
where id = $1
returning currency_cents, nibbles, namicoin
`,
			player.ID,
			creditsAfter,
			nibblesAfter,
			namiCoinAfter,
		).Scan(&entry.CreditsAfter, &entry.NibblesAfter, &entry.NamiCoinAfter); err != nil {
			return DevCurrencyResult{}, fmt.Errorf("remove currency for %s: %w", player.Name, err)
		}

		if err := tx.QueryRow(ctx, `
update player_resources
set fans = $2,
memes = $3,
lost_items = $4,
confidence = $5,
receipts = $6,
patterns = $7,
updated_at = now()
where player_id = $1
returning fans, memes, lost_items, confidence, receipts, patterns
`,
			player.ID,
			fansAfter,
			memesAfter,
			lostItemsAfter,
			confidenceAfter,
			receiptsAfter,
			patternsAfter,
		).Scan(
			&fansAfter,
			&memesAfter,
			&lostItemsAfter,
			&confidenceAfter,
			&receiptsAfter,
			&patternsAfter,
		); err != nil {
			return DevCurrencyResult{}, fmt.Errorf("remove resources for %s: %w", player.Name, err)
		}

		entry.Before = currencyMap(entry.CreditsBefore, entry.NibblesBefore, entry.NamiCoinBefore, fansBefore, memesBefore, lostItemsBefore, confidenceBefore, receiptsBefore, patternsBefore)
		entry.After = currencyMap(entry.CreditsAfter, entry.NibblesAfter, entry.NamiCoinAfter, fansAfter, memesAfter, lostItemsAfter, confidenceAfter, receiptsAfter, patternsAfter)
		result.Players = append(result.Players, entry)
	}

	if err := tx.Commit(ctx); err != nil {
		return DevCurrencyResult{}, fmt.Errorf("commit remove currency: %w", err)
	}

	if amountWasAll {
		result.Message = fmt.Sprintf("Removed all selected assets from %d player(s).", len(players))
	} else {
		result.Message = fmt.Sprintf("Removed up to %d from %s for %d player(s).", amount, currencyLabel, len(players))
	}

	return result, nil
}
func normalizeDevLevelResetSpecs(raw string) ([]devLevelResetSpec, string, error) {
	value := strings.TrimSpace(strings.ToLower(raw))
	value = strings.ReplaceAll(value, "_", " ")
	value = strings.ReplaceAll(value, "-", " ")
	value = strings.Join(strings.Fields(value), " ")

	allSpecs := []devLevelResetSpec{
		{Key: "playdeck", Label: "Playdeck Level", Kind: "player"},
		{Key: "nami", Label: "Nami Level", Kind: "companion"},
		{Key: "streaming", Label: "Streaming Work Level", Kind: "activity"},
		{Key: "doom_scrolling", Label: "Doom Scrolling Work Level", Kind: "activity"},
		{Key: "cleaning", Label: "Cleaning Work Level", Kind: "activity"},
		{Key: "exercising", Label: "Exercise Work Level", Kind: "activity"},
		{Key: "shopping", Label: "Shopping Work Level", Kind: "activity"},
		{Key: "designing", Label: "Designing Work Level", Kind: "activity"},
	}

	if value == "all" {
		return allSpecs, "ALL", nil
	}

	switch value {
	case "playdeck", "playdeck level", "player level":
		return []devLevelResetSpec{allSpecs[0]}, allSpecs[0].Label, nil
	case "nami", "nami level", "companion", "companion level":
		return []devLevelResetSpec{allSpecs[1]}, allSpecs[1].Label, nil
	case "streaming", "streaming level", "streaming work", "streaming work level":
		return []devLevelResetSpec{allSpecs[2]}, allSpecs[2].Label, nil
	case "doom scrolling", "doom scrolling level", "doom scrolling work", "doom scrolling work level", "scrolling", "scrolling level":
		return []devLevelResetSpec{allSpecs[3]}, allSpecs[3].Label, nil
	case "cleaning", "cleaning level", "cleaning work", "cleaning work level":
		return []devLevelResetSpec{allSpecs[4]}, allSpecs[4].Label, nil
	case "exercising", "exercise", "exercise level", "exercising level", "exercise work", "exercise work level", "exercising work level":
		return []devLevelResetSpec{allSpecs[5]}, allSpecs[5].Label, nil
	case "shopping", "shopping level", "shopping work", "shopping work level":
		return []devLevelResetSpec{allSpecs[6]}, allSpecs[6].Label, nil
	case "designing", "designing level", "designing work", "designing work level":
		return []devLevelResetSpec{allSpecs[7]}, allSpecs[7].Label, nil
	default:
		return nil, "", fmt.Errorf("activityName must be Playdeck Level, Nami Level, one of the six work levels, or ALL")
	}
}

func (s *Store) DevResetLevels(ctx context.Context, playerName string, activityName string) (DevResetLevelsResult, error) {
	target, players, err := s.resolveDevTargetPlayers(ctx, playerName, "RESETALL")
	if err != nil {
		return DevResetLevelsResult{}, err
	}

	specs, activityLabel, err := normalizeDevLevelResetSpecs(activityName)
	if err != nil {
		return DevResetLevelsResult{}, err
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return DevResetLevelsResult{}, fmt.Errorf("begin reset levels: %w", err)
	}
	defer tx.Rollback(ctx)

	result := DevResetLevelsResult{
		OK:            true,
		ActivityInput: activityLabel,
		Target:        target,
		Players:       make([]DevLevelResetPlayer, 0, len(players)),
	}

	for _, player := range players {
		playerResult := DevLevelResetPlayer{
			PlayerID:   player.ID,
			PlayerName: player.Name,
			Resets:     make([]DevLevelResetEntry, 0, len(specs)),
		}

		for _, spec := range specs {
			entry := DevLevelResetEntry{
				Key:   spec.Key,
				Label: spec.Label,
			}

			switch spec.Kind {
			case "player":
				if err := tx.QueryRow(ctx, `
select level, total_xp, xp_into_level
from players
where id = $1
for update
`, player.ID).Scan(&entry.BeforeLevel, &entry.BeforeTotalXP, &entry.BeforeXPIntoLevel); err != nil {
					return DevResetLevelsResult{}, fmt.Errorf("load playdeck level for %s: %w", player.Name, err)
				}

				if _, err := tx.Exec(ctx, `
update players
set level = 1,
total_xp = 0,
xp_into_level = 0,
updated_at = now()
where id = $1
`, player.ID); err != nil {
					return DevResetLevelsResult{}, fmt.Errorf("reset playdeck level for %s: %w", player.Name, err)
				}

			case "companion":
				if err := tx.QueryRow(ctx, `
select level, total_xp, xp_into_level
from companion_states
where player_id = $1
for update
`, player.ID).Scan(&entry.BeforeLevel, &entry.BeforeTotalXP, &entry.BeforeXPIntoLevel); err != nil {
					return DevResetLevelsResult{}, fmt.Errorf("load Nami level for %s: %w", player.Name, err)
				}

				if _, err := tx.Exec(ctx, `
update companion_states
set level = 1,
total_xp = 0,
xp_into_level = 0,
last_xp_gained = 0,
updated_at = now()
where player_id = $1
`, player.ID); err != nil {
					return DevResetLevelsResult{}, fmt.Errorf("reset Nami level for %s: %w", player.Name, err)
				}

			case "activity":
				if _, err := tx.Exec(ctx, `
insert into player_activity_skills (player_id, activity_key, level, total_xp, xp_into_level)
values ($1, $2, 1, 0, 0)
on conflict (player_id, activity_key) do nothing
`, player.ID, spec.Key); err != nil {
					return DevResetLevelsResult{}, fmt.Errorf("ensure activity level for %s / %s: %w", player.Name, spec.Label, err)
				}

				if err := tx.QueryRow(ctx, `
select level, total_xp, xp_into_level
from player_activity_skills
where player_id = $1
and activity_key = $2
for update
`, player.ID, spec.Key).Scan(&entry.BeforeLevel, &entry.BeforeTotalXP, &entry.BeforeXPIntoLevel); err != nil {
					return DevResetLevelsResult{}, fmt.Errorf("load activity level for %s / %s: %w", player.Name, spec.Label, err)
				}

				if _, err := tx.Exec(ctx, `
update player_activity_skills
set level = 1,
total_xp = 0,
xp_into_level = 0,
updated_at = now()
where player_id = $1
and activity_key = $2
`, player.ID, spec.Key); err != nil {
					return DevResetLevelsResult{}, fmt.Errorf("reset activity level for %s / %s: %w", player.Name, spec.Label, err)
				}
			}

			entry.AfterLevel = 1
			entry.AfterTotalXP = 0
			entry.AfterXPIntoLevel = 0
			playerResult.Resets = append(playerResult.Resets, entry)
		}

		result.Players = append(result.Players, playerResult)
	}

	if err := tx.Commit(ctx); err != nil {
		return DevResetLevelsResult{}, fmt.Errorf("commit reset levels: %w", err)
	}

	result.Message = fmt.Sprintf("Reset %s for %d player(s).", activityLabel, len(players))
	return result, nil
}

type DevResetServerResult struct {
	OK                          bool    `json:"ok"`
	Message                     string  `json:"message"`
	PreservedDisplayName        string  `json:"preservedDisplayName"`
	PreservedPlayerIDs          []int64 `json:"preservedPlayerIds"`
	PreservedAccountIDs         []int64 `json:"preservedAccountIds"`
	DeletedPlayers              int64   `json:"deletedPlayers"`
	DeletedAuthAccounts         int64   `json:"deletedAuthAccounts"`
	DeletedAuthSessions         int64   `json:"deletedAuthSessions"`
	DeletedAuthCredentials      int64   `json:"deletedAuthCredentials"`
	DeletedAuthIdentities       int64   `json:"deletedAuthIdentities"`
	DeletedWardrobeItems        int64   `json:"deletedWardrobeItems"`
	DeletedCareActions          int64   `json:"deletedCareActions"`
	DeletedNamiMessages         int64   `json:"deletedNamiMessages"`
	DeletedActivityLogs         int64   `json:"deletedActivityLogs"`
	DeletedPlaydeckCombatLogs   int64   `json:"deletedPlaydeckCombatLogs"`
	ResetPlayers                int64   `json:"resetPlayers"`
	ResetResourceRows           int64   `json:"resetResourceRows"`
	ResetActivitySkillRows      int64   `json:"resetActivitySkillRows"`
	ResetCompanionRows          int64   `json:"resetCompanionRows"`
	ResetTickStateRows          int64   `json:"resetTickStateRows"`
	ResetPlaydeckStateRows      int64   `json:"resetPlaydeckStateRows"`
	ResetPlaydeckZoneRecordRows int64   `json:"resetPlaydeckZoneRecordRows"`
}

func devInt64ListSQL(ids []int64) string {
	seen := map[int64]bool{}
	parts := []string{}

	for _, id := range ids {
		if id <= 0 || seen[id] {
			continue
		}

		seen[id] = true
		parts = append(parts, strconv.FormatInt(id, 10))
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, ", ")
}

func (s *Store) DevResetServer(ctx context.Context) (DevResetServerResult, error) {
	const preservedDisplayName = "Soryn"

	if err := s.ensurePlaydeckZoneRecordsTable(ctx); err != nil {
		return DevResetServerResult{}, err
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return DevResetServerResult{}, fmt.Errorf("begin reset server: %w", err)
	}
	defer tx.Rollback(ctx)

	tableExists := func(tableName string) (bool, error) {
		var exists bool
		err := tx.QueryRow(ctx, `
select to_regclass($1) is not null
`, "public."+tableName).Scan(&exists)
		if err != nil {
			return false, fmt.Errorf("check table %s: %w", tableName, err)
		}

		return exists, nil
	}

	columnExists := func(tableName string, columnName string) (bool, error) {
		var exists bool
		err := tx.QueryRow(ctx, `
select exists (
select 1
from information_schema.columns
where table_schema = 'public'
and table_name = $1
and column_name = $2
)
`, tableName, columnName).Scan(&exists)
		if err != nil {
			return false, fmt.Errorf("check column %s.%s: %w", tableName, columnName, err)
		}

		return exists, nil
	}

	deleteAllIfExists := func(tableName string) (int64, error) {
		exists, err := tableExists(tableName)
		if err != nil {
			return 0, err
		}

		if !exists {
			return 0, nil
		}

		commandTag, err := tx.Exec(ctx, fmt.Sprintf(`delete from %s`, tableName))
		if err != nil {
			return 0, fmt.Errorf("delete from %s: %w", tableName, err)
		}

		return commandTag.RowsAffected(), nil
	}

	deleteWhereAccountIDIfExists := func(tableName string, accountIDFilter string) (int64, error) {
		exists, err := tableExists(tableName)
		if err != nil {
			return 0, err
		}

		if !exists {
			return 0, nil
		}

		commandTag, err := tx.Exec(ctx, fmt.Sprintf(`delete from %s where %s`, tableName, accountIDFilter))
		if err != nil {
			return 0, fmt.Errorf("delete from %s: %w", tableName, err)
		}

		return commandTag.RowsAffected(), nil
	}

	result := DevResetServerResult{
		OK:                   true,
		PreservedDisplayName: preservedDisplayName,
	}

	preservedAccountIDs := map[int64]bool{}

	accountRows, err := tx.Query(ctx, `
select id
from auth_accounts
where lower(display_name) = lower($1)
`, preservedDisplayName)
	if err != nil {
		return DevResetServerResult{}, fmt.Errorf("load preserved auth accounts: %w", err)
	}

	for accountRows.Next() {
		var accountID int64
		if err := accountRows.Scan(&accountID); err != nil {
			accountRows.Close()
			return DevResetServerResult{}, fmt.Errorf("scan preserved auth account: %w", err)
		}

		preservedAccountIDs[accountID] = true
	}
	if err := accountRows.Err(); err != nil {
		accountRows.Close()
		return DevResetServerResult{}, fmt.Errorf("iterate preserved auth accounts: %w", err)
	}
	accountRows.Close()

	hasPlayerAccountID, err := columnExists("players", "account_id")
	if err != nil {
		return DevResetServerResult{}, err
	}

	if hasPlayerAccountID {
		rows, err := tx.Query(ctx, `
select coalesce(account_id, 0)
from players
where lower(display_name) = lower($1)
`, preservedDisplayName)
		if err != nil {
			return DevResetServerResult{}, fmt.Errorf("load preserved player account ids: %w", err)
		}

		for rows.Next() {
			var accountID int64
			if err := rows.Scan(&accountID); err != nil {
				rows.Close()
				return DevResetServerResult{}, fmt.Errorf("scan preserved player account id: %w", err)
			}

			if accountID > 0 {
				preservedAccountIDs[accountID] = true
			}
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return DevResetServerResult{}, fmt.Errorf("iterate preserved player account ids: %w", err)
		}
		rows.Close()
	}

	for accountID := range preservedAccountIDs {
		result.PreservedAccountIDs = append(result.PreservedAccountIDs, accountID)
	}

	commandTag, err := tx.Exec(ctx, `
delete from players
where lower(display_name) <> lower($1)
`, preservedDisplayName)
	if err != nil {
		return DevResetServerResult{}, fmt.Errorf("delete non-Soryn players: %w", err)
	}
	result.DeletedPlayers = commandTag.RowsAffected()

	accountIDFilter := "true"
	accountIDList := devInt64ListSQL(result.PreservedAccountIDs)
	if accountIDList != "" {
		accountIDFilter = "account_id not in (" + accountIDList + ")"
	}

	result.DeletedAuthSessions, err = deleteWhereAccountIDIfExists("auth_sessions", accountIDFilter)
	if err != nil {
		return DevResetServerResult{}, err
	}

	result.DeletedAuthCredentials, err = deleteWhereAccountIDIfExists("auth_credentials", accountIDFilter)
	if err != nil {
		return DevResetServerResult{}, err
	}

	result.DeletedAuthIdentities, err = deleteWhereAccountIDIfExists("auth_identities", accountIDFilter)
	if err != nil {
		return DevResetServerResult{}, err
	}

	authAccountFilter := "true"
	if accountIDList != "" {
		authAccountFilter = "id not in (" + accountIDList + ")"
	}

	commandTag, err = tx.Exec(ctx, "delete from auth_accounts where "+authAccountFilter)
	if err != nil {
		return DevResetServerResult{}, fmt.Errorf("delete non-Soryn auth accounts: %w", err)
	}
	result.DeletedAuthAccounts = commandTag.RowsAffected()

	rows, err := tx.Query(ctx, `
select id
from players
where lower(display_name) = lower($1)
order by id
`, preservedDisplayName)
	if err != nil {
		return DevResetServerResult{}, fmt.Errorf("load preserved players: %w", err)
	}

	for rows.Next() {
		var playerID int64
		if err := rows.Scan(&playerID); err != nil {
			rows.Close()
			return DevResetServerResult{}, fmt.Errorf("scan preserved player: %w", err)
		}

		result.PreservedPlayerIDs = append(result.PreservedPlayerIDs, playerID)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return DevResetServerResult{}, fmt.Errorf("iterate preserved players: %w", err)
	}
	rows.Close()

	if len(result.PreservedPlayerIDs) == 0 {
		return DevResetServerResult{}, fmt.Errorf("cannot reset server: preserved player %q was not found", preservedDisplayName)
	}

	if err := tx.QueryRow(ctx, `
select count(*)::bigint
from player_inventory_items
`).Scan(&result.DeletedWardrobeItems); err != nil {
		return DevResetServerResult{}, fmt.Errorf("count wardrobe items before reset: %w", err)
	}

	if _, err := tx.Exec(ctx, `
truncate table player_inventory_items restart identity cascade
`); err != nil {
		return DevResetServerResult{}, fmt.Errorf("truncate wardrobe inventory: %w", err)
	}

	result.DeletedCareActions, err = deleteAllIfExists("companion_care_actions")
	if err != nil {
		return DevResetServerResult{}, err
	}

	result.DeletedNamiMessages, err = deleteAllIfExists("nami_messages")
	if err != nil {
		return DevResetServerResult{}, err
	}

	result.DeletedActivityLogs, err = deleteAllIfExists("activity_log")
	if err != nil {
		return DevResetServerResult{}, err
	}

	combatLogCountA, err := deleteAllIfExists("playdeck_combat_log")
	if err != nil {
		return DevResetServerResult{}, err
	}

	combatLogCountB, err := deleteAllIfExists("playdeck_combat_logs")
	if err != nil {
		return DevResetServerResult{}, err
	}

	result.DeletedPlaydeckCombatLogs = combatLogCountA + combatLogCountB

	commandTag, err = tx.Exec(ctx, `
update players
set level = 1,
total_xp = 0,
xp_into_level = 0,
currency_cents = 0,
nibbles = 0,
namicoin = 0,
updated_at = now()
`)
	if err != nil {
		return DevResetServerResult{}, fmt.Errorf("reset remaining players: %w", err)
	}
	result.ResetPlayers = commandTag.RowsAffected()

	if _, err := tx.Exec(ctx, `delete from player_resources`); err != nil {
		return DevResetServerResult{}, fmt.Errorf("clear player resources: %w", err)
	}

	commandTag, err = tx.Exec(ctx, `
insert into player_resources (player_id)
select id
from players
`)
	if err != nil {
		return DevResetServerResult{}, fmt.Errorf("recreate player resources: %w", err)
	}
	result.ResetResourceRows = commandTag.RowsAffected()

	if _, err := tx.Exec(ctx, `delete from player_activity_skills`); err != nil {
		return DevResetServerResult{}, fmt.Errorf("clear activity skills: %w", err)
	}

	commandTag, err = tx.Exec(ctx, `
insert into player_activity_skills (player_id, activity_key)
select p.id, task.activity_key
from players p
cross join (
values
('streaming'),
('doom_scrolling'),
('cleaning'),
('exercising'),
('shopping'),
('designing')
) as task(activity_key)
`)
	if err != nil {
		return DevResetServerResult{}, fmt.Errorf("recreate activity skills: %w", err)
	}
	result.ResetActivitySkillRows = commandTag.RowsAffected()

	if _, err := tx.Exec(ctx, `delete from companion_states`); err != nil {
		return DevResetServerResult{}, fmt.Errorf("clear companion states: %w", err)
	}

	commandTag, err = tx.Exec(ctx, `
insert into companion_states (
player_id,
companion_name,
mood_score,
satiety,
connection,
energy,
comfort,
playfulness,
inspiration,
cleanliness,
status,
last_interaction_at,
last_decay_at
)
select
id,
'Nami-chan',
80.00,
85,
95,
75,
90,
80,
80,
80,
'awake',
now(),
now()
from players
`)
	if err != nil {
		return DevResetServerResult{}, fmt.Errorf("recreate companion states: %w", err)
	}
	result.ResetCompanionRows = commandTag.RowsAffected()

	if _, err := tx.Exec(ctx, `delete from player_tick_state`); err != nil {
		return DevResetServerResult{}, fmt.Errorf("clear player tick state: %w", err)
	}

	commandTag, err = tx.Exec(ctx, `
insert into player_tick_state (player_id)
select id
from players
`)
	if err != nil {
		return DevResetServerResult{}, fmt.Errorf("recreate player tick state: %w", err)
	}
	result.ResetTickStateRows = commandTag.RowsAffected()

	commandTag, err = tx.Exec(ctx, `delete from player_playdeck_zone_records`)
	if err != nil {
		return DevResetServerResult{}, fmt.Errorf("clear playdeck zone records: %w", err)
	}
	result.ResetPlaydeckZoneRecordRows = commandTag.RowsAffected()

	if _, err := tx.Exec(ctx, `delete from playdeck_combat_state`); err != nil {
		return DevResetServerResult{}, fmt.Errorf("clear playdeck combat state: %w", err)
	}

	for _, playerID := range result.PreservedPlayerIDs {
		if err := ensurePlaydeckStateTx(ctx, tx, playerID); err != nil {
			return DevResetServerResult{}, fmt.Errorf("recreate playdeck combat state for player %d: %w", playerID, err)
		}
		result.ResetPlaydeckStateRows++
	}

	if err := tx.Commit(ctx); err != nil {
		return DevResetServerResult{}, fmt.Errorf("commit reset server: %w", err)
	}

	result.Message = "Server reset complete. Soryn was preserved; other player accounts were deleted and core player state was reset."

	return result, nil
}
