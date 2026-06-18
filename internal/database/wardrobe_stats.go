package database

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
)

type WardrobeStatLine struct {
	InventoryItemID int64   `json:"-"`
	Source          string  `json:"source"`
	AffixKey        string  `json:"affixKey"`
	StatKey         string  `json:"statKey"`
	DisplayName     string  `json:"displayName"`
	ValueKind       string  `json:"valueKind"`
	Value           float64 `json:"value"`
	Tooltip         string  `json:"tooltip"`
	SortOrder       int     `json:"sortOrder"`
}

type WardrobeStatComparison struct {
	StatKey       string  `json:"statKey"`
	DisplayName   string  `json:"displayName"`
	ValueKind     string  `json:"valueKind"`
	Tooltip       string  `json:"tooltip"`
	ItemValue     float64 `json:"itemValue"`
	EquippedValue float64 `json:"equippedValue"`
	Delta         float64 `json:"delta"`
}

type WardrobeAccessoryCompareSlot struct {
	SlotKey     string `json:"slotKey"`
	DisplayName string `json:"displayName"`
	Selected    bool   `json:"selected"`
}

type WardrobeItemDetail struct {
	Item                  InventoryItemStatus            `json:"item"`
	CompareSlot           string                         `json:"compareSlot"`
	CompareItem           *InventoryItemStatus           `json:"compareItem,omitempty"`
	Comparisons           []WardrobeStatComparison       `json:"comparisons"`
	AccessoryCompareSlots []WardrobeAccessoryCompareSlot `json:"accessoryCompareSlots,omitempty"`
}

func (s *Store) GetWardrobeItemDetail(ctx context.Context, playerID int64, itemID int64, compareSlot string) (WardrobeItemDetail, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return WardrobeItemDetail{}, fmt.Errorf("begin wardrobe item detail: %w", err)
	}
	defer tx.Rollback(ctx)

	item, err := loadWardrobeInventoryItemByIDTx(ctx, tx, playerID, itemID)
	if err != nil {
		return WardrobeItemDetail{}, err
	}

	itemStatLines, err := loadWardrobeStatLinesForItemsTx(ctx, tx, []int64{item.ID})
	if err != nil {
		return WardrobeItemDetail{}, err
	}

	item.StatLines = itemStatLines[item.ID]

	resolvedCompareSlot, err := resolveWardrobeCompareSlotTx(ctx, tx, item, compareSlot)
	if err != nil {
		return WardrobeItemDetail{}, err
	}

	detail := WardrobeItemDetail{
		Item:        item,
		CompareSlot: resolvedCompareSlot,
	}

	if item.EquipmentSlot == "accessory" {
		detail.AccessoryCompareSlots = []WardrobeAccessoryCompareSlot{
			{SlotKey: "accessory_1", DisplayName: "Accessory 1", Selected: resolvedCompareSlot == "accessory_1"},
			{SlotKey: "accessory_2", DisplayName: "Accessory 2", Selected: resolvedCompareSlot == "accessory_2"},
		}
	}

	var equippedStatLines []WardrobeStatLine
	if resolvedCompareSlot != "" {
		equippedItem, found, err := loadEquippedWardrobeItemInSlotTx(ctx, tx, playerID, resolvedCompareSlot)
		if err != nil {
			return WardrobeItemDetail{}, err
		}

		if found {
			equippedStatLineMap, err := loadWardrobeStatLinesForItemsTx(ctx, tx, []int64{equippedItem.ID})
			if err != nil {
				return WardrobeItemDetail{}, err
			}

			equippedItem.StatLines = equippedStatLineMap[equippedItem.ID]
			equippedStatLines = equippedItem.StatLines
			detail.CompareItem = &equippedItem
		}
	}

	detail.Comparisons = buildWardrobeStatComparisons(item.StatLines, equippedStatLines)

	if err := tx.Commit(ctx); err != nil {
		return WardrobeItemDetail{}, fmt.Errorf("commit wardrobe item detail: %w", err)
	}

	return detail, nil
}

func loadWardrobeInventoryItemByIDTx(ctx context.Context, tx pgx.Tx, playerID int64, itemID int64) (InventoryItemStatus, error) {
	var item InventoryItemStatus

	if err := tx.QueryRow(ctx, `
		select
			i.id,
			d.item_key,
			d.name,
			d.item_type,
			d.rarity,
			d.equipment_slot,
			i.quantity,
			coalesce(i.equipped_slot, ''),
			d.power_level,
			d.attack_bonus,
			d.defense_bonus,
			d.max_hp_bonus,
			i.tailoring_current,
			i.tailoring_max
		from player_inventory_items i
		join item_definitions d on d.id = i.item_definition_id
		where i.player_id = $1
			and i.id = $2
			and i.container = 'wardrobe'
	`, playerID, itemID).Scan(
		&item.ID,
		&item.ItemKey,
		&item.Name,
		&item.ItemType,
		&item.Rarity,
		&item.EquipmentSlot,
		&item.Quantity,
		&item.EquippedSlot,
		&item.PowerLevel,
		&item.AttackBonus,
		&item.DefenseBonus,
		&item.MaxHPBonus,
		&item.TailoringCurrent,
		&item.TailoringMax,
	); err != nil {
		return InventoryItemStatus{}, fmt.Errorf("load wardrobe item detail: %w", err)
	}

	return item, nil
}

func loadEquippedWardrobeItemInSlotTx(ctx context.Context, tx pgx.Tx, playerID int64, slotKey string) (InventoryItemStatus, bool, error) {
	var item InventoryItemStatus

	if err := tx.QueryRow(ctx, `
		select
			i.id,
			d.item_key,
			d.name,
			d.item_type,
			d.rarity,
			d.equipment_slot,
			i.quantity,
			coalesce(i.equipped_slot, ''),
			d.power_level,
			d.attack_bonus,
			d.defense_bonus,
			d.max_hp_bonus,
			i.tailoring_current,
			i.tailoring_max
		from player_inventory_items i
		join item_definitions d on d.id = i.item_definition_id
		where i.player_id = $1
			and i.equipped_slot = $2
			and i.container = 'wardrobe'
		limit 1
	`, playerID, slotKey).Scan(
		&item.ID,
		&item.ItemKey,
		&item.Name,
		&item.ItemType,
		&item.Rarity,
		&item.EquipmentSlot,
		&item.Quantity,
		&item.EquippedSlot,
		&item.PowerLevel,
		&item.AttackBonus,
		&item.DefenseBonus,
		&item.MaxHPBonus,
		&item.TailoringCurrent,
		&item.TailoringMax,
	); err != nil {
		if err == pgx.ErrNoRows {
			return InventoryItemStatus{}, false, nil
		}

		return InventoryItemStatus{}, false, fmt.Errorf("load equipped wardrobe item in slot: %w", err)
	}

	return item, true, nil
}

func resolveWardrobeCompareSlotTx(ctx context.Context, tx pgx.Tx, item InventoryItemStatus, compareSlot string) (string, error) {
	equipmentSlot := strings.TrimSpace(strings.ToLower(item.EquipmentSlot))
	compareSlot = strings.TrimSpace(strings.ToLower(compareSlot))

	if equipmentSlot == "" || equipmentSlot == "none" {
		return "", nil
	}

	if equipmentSlot == "accessory" {
		if compareSlot == "accessory_2" {
			return "accessory_2", nil
		}

		return "accessory_1", nil
	}

	if compareSlot != "" {
		var count int
		if err := tx.QueryRow(ctx, `
			select count(*)::int
			from playdeck_equipment_slots
			where slot_key = $1
				and accepts_slot = $2
		`, compareSlot, equipmentSlot).Scan(&count); err != nil {
			return "", fmt.Errorf("validate wardrobe compare slot: %w", err)
		}

		if count > 0 {
			return compareSlot, nil
		}
	}

	var resolvedSlot string
	if err := tx.QueryRow(ctx, `
		select slot_key
		from playdeck_equipment_slots
		where accepts_slot = $1
		order by sort_order
		limit 1
	`, equipmentSlot).Scan(&resolvedSlot); err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}

		return "", fmt.Errorf("resolve wardrobe compare slot: %w", err)
	}

	return resolvedSlot, nil
}

func loadWardrobeStatLinesForItemsTx(ctx context.Context, tx pgx.Tx, itemIDs []int64) (map[int64][]WardrobeStatLine, error) {
	result := make(map[int64][]WardrobeStatLine)

	uniqueIDs := uniquePositiveInt64s(itemIDs)
	if len(uniqueIDs) == 0 {
		return result, nil
	}

	rows, err := tx.Query(ctx, `
		select
			l.player_inventory_item_id,
			l.stat_source,
			l.affix_key,
			l.stat_key,
			d.display_name,
			d.value_kind,
			l.value::double precision,
			d.tooltip,
			l.sort_order
		from player_inventory_item_stat_lines l
		join wardrobe_stat_definitions d on d.stat_key = l.stat_key
		where l.player_inventory_item_id = any($1)
		order by
			l.player_inventory_item_id,
			case l.stat_source
				when 'implicit' then 10
				when 'prefix' then 20
				when 'suffix' then 30
				when 'devastating_implicit' then 40
				else 90
			end,
			l.sort_order,
			d.sort_order,
			l.id
	`, uniqueIDs)
	if err != nil {
		return nil, fmt.Errorf("query wardrobe stat lines: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var line WardrobeStatLine

		if err := rows.Scan(
			&line.InventoryItemID,
			&line.Source,
			&line.AffixKey,
			&line.StatKey,
			&line.DisplayName,
			&line.ValueKind,
			&line.Value,
			&line.Tooltip,
			&line.SortOrder,
		); err != nil {
			return nil, fmt.Errorf("scan wardrobe stat line: %w", err)
		}

		result[line.InventoryItemID] = append(result[line.InventoryItemID], line)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate wardrobe stat lines: %w", err)
	}

	return result, nil
}

func buildWardrobeStatComparisons(itemLines []WardrobeStatLine, equippedLines []WardrobeStatLine) []WardrobeStatComparison {
	itemTotals := sumWardrobeStatLines(itemLines)
	equippedTotals := sumWardrobeStatLines(equippedLines)

	keys := make(map[string]bool, len(itemTotals)+len(equippedTotals))
	for key := range itemTotals {
		keys[key] = true
	}
	for key := range equippedTotals {
		keys[key] = true
	}

	comparisons := make([]WardrobeStatComparison, 0, len(keys))
	for key := range keys {
		itemTotal := itemTotals[key]
		equippedTotal := equippedTotals[key]

		displayName := itemTotal.DisplayName
		valueKind := itemTotal.ValueKind
		tooltip := itemTotal.Tooltip

		if displayName == "" {
			displayName = equippedTotal.DisplayName
			valueKind = equippedTotal.ValueKind
			tooltip = equippedTotal.Tooltip
		}

		comparisons = append(comparisons, WardrobeStatComparison{
			StatKey:       key,
			DisplayName:   displayName,
			ValueKind:     valueKind,
			Tooltip:       tooltip,
			ItemValue:     itemTotal.Value,
			EquippedValue: equippedTotal.Value,
			Delta:         itemTotal.Value - equippedTotal.Value,
		})
	}

	sort.SliceStable(comparisons, func(i, j int) bool {
		return comparisons[i].DisplayName < comparisons[j].DisplayName
	})

	return comparisons
}

func sumWardrobeStatLines(lines []WardrobeStatLine) map[string]WardrobeStatLine {
	totals := make(map[string]WardrobeStatLine)

	for _, line := range lines {
		total := totals[line.StatKey]
		if total.StatKey == "" {
			total = WardrobeStatLine{
				StatKey:     line.StatKey,
				DisplayName: line.DisplayName,
				ValueKind:   line.ValueKind,
				Tooltip:     line.Tooltip,
			}
		}

		total.Value += line.Value
		totals[line.StatKey] = total
	}

	return totals
}

func uniquePositiveInt64s(values []int64) []int64 {
	seen := make(map[int64]bool, len(values))
	result := make([]int64, 0, len(values))

	for _, value := range values {
		if value <= 0 || seen[value] {
			continue
		}

		seen[value] = true
		result = append(result, value)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})

	return result
}
