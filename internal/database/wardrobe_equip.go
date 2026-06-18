package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (s *Store) EquipWardrobeItem(ctx context.Context, playerID int64, itemID int64, requestedSlotKey string) (WardrobeItemDetail, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return WardrobeItemDetail{}, fmt.Errorf("begin equip wardrobe item: %w", err)
	}
	defer tx.Rollback(ctx)

	item, err := loadWardrobeInventoryItemByIDTx(ctx, tx, playerID, itemID)
	if err != nil {
		return WardrobeItemDetail{}, err
	}

	if strings.ToLower(item.ItemType) != "gear" {
		return WardrobeItemDetail{}, fmt.Errorf("item is not gear")
	}

	targetSlotKey, err := resolveWardrobeEquipSlotTx(ctx, tx, item.EquipmentSlot, requestedSlotKey)
	if err != nil {
		return WardrobeItemDetail{}, err
	}

	if targetSlotKey == "" {
		return WardrobeItemDetail{}, fmt.Errorf("item has no wearable slot")
	}

	if _, err := tx.Exec(ctx, `
update player_inventory_items
set equipped_slot = null,
updated_at = now()
where player_id = $1
and container = 'wardrobe'
and (
id = $2
or equipped_slot = $3
)
`, playerID, itemID, targetSlotKey); err != nil {
		return WardrobeItemDetail{}, fmt.Errorf("clear wardrobe slot before equip: %w", err)
	}

	commandTag, err := tx.Exec(ctx, `
update player_inventory_items
set equipped_slot = $3,
updated_at = now()
where player_id = $1
and id = $2
and container = 'wardrobe'
`, playerID, itemID, targetSlotKey)
	if err != nil {
		return WardrobeItemDetail{}, fmt.Errorf("equip wardrobe item: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return WardrobeItemDetail{}, fmt.Errorf("wardrobe item not found")
	}

	if err := normalizePlaydeckHPAfterWardrobeChangeTx(ctx, tx, playerID); err != nil {
		return WardrobeItemDetail{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return WardrobeItemDetail{}, fmt.Errorf("commit equip wardrobe item: %w", err)
	}

	return s.GetWardrobeItemDetail(ctx, playerID, itemID, targetSlotKey)
}

func (s *Store) UnequipWardrobeItem(ctx context.Context, playerID int64, itemID int64, requestedSlotKey string) (WardrobeItemDetail, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return WardrobeItemDetail{}, fmt.Errorf("begin unequip wardrobe item: %w", err)
	}
	defer tx.Rollback(ctx)

	item, err := loadWardrobeInventoryItemByIDTx(ctx, tx, playerID, itemID)
	if err != nil {
		return WardrobeItemDetail{}, err
	}

	compareSlot := strings.TrimSpace(strings.ToLower(requestedSlotKey))
	if compareSlot == "" {
		compareSlot = strings.TrimSpace(strings.ToLower(item.EquippedSlot))
	}

	commandTag, err := tx.Exec(ctx, `
update player_inventory_items
set equipped_slot = null,
updated_at = now()
where player_id = $1
and id = $2
and container = 'wardrobe'
`, playerID, itemID)
	if err != nil {
		return WardrobeItemDetail{}, fmt.Errorf("unequip wardrobe item: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return WardrobeItemDetail{}, fmt.Errorf("wardrobe item not found")
	}

	if err := normalizePlaydeckHPAfterWardrobeChangeTx(ctx, tx, playerID); err != nil {
		return WardrobeItemDetail{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return WardrobeItemDetail{}, fmt.Errorf("commit unequip wardrobe item: %w", err)
	}

	return s.GetWardrobeItemDetail(ctx, playerID, itemID, compareSlot)
}

func resolveWardrobeEquipSlotTx(ctx context.Context, tx pgx.Tx, equipmentSlot string, requestedSlotKey string) (string, error) {
	equipmentSlot = strings.TrimSpace(strings.ToLower(equipmentSlot))
	requestedSlotKey = strings.TrimSpace(strings.ToLower(requestedSlotKey))

	if equipmentSlot == "" || equipmentSlot == "none" {
		return "", nil
	}

	if equipmentSlot == "accessory" && requestedSlotKey == "" {
		requestedSlotKey = "accessory_1"
	}

	if requestedSlotKey != "" {
		var count int
		if err := tx.QueryRow(ctx, `
select count(*)::int
from playdeck_equipment_slots
where slot_key = $1
and accepts_slot = $2
`, requestedSlotKey, equipmentSlot).Scan(&count); err != nil {
			return "", fmt.Errorf("validate wardrobe equip slot: %w", err)
		}

		if count == 0 {
			return "", fmt.Errorf("item cannot be worn in %s", requestedSlotKey)
		}

		return requestedSlotKey, nil
	}

	var resolvedSlotKey string
	if err := tx.QueryRow(ctx, `
select slot_key
from playdeck_equipment_slots
where accepts_slot = $1
order by sort_order
limit 1
`, equipmentSlot).Scan(&resolvedSlotKey); err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}

		return "", fmt.Errorf("resolve wardrobe equip slot: %w", err)
	}

	return resolvedSlotKey, nil
}

func normalizePlaydeckHPAfterWardrobeChangeTx(ctx context.Context, tx pgx.Tx, playerID int64) error {
	if err := ensurePlaydeckStateTx(ctx, tx, playerID); err != nil {
		return err
	}

	var level int
	if err := tx.QueryRow(ctx, `
select level
from players
where id = $1
`, playerID).Scan(&level); err != nil {
		return fmt.Errorf("load level after wardrobe change: %w", err)
	}

	equipmentStats, err := loadPlaydeckEquipmentStatsTx(ctx, tx, playerID)
	if err != nil {
		return err
	}

	maxHP := PlaydeckMaxHP(level, equipmentStats)

	if _, err := tx.Exec(ctx, `
update playdeck_combat_state
set player_hp = least(greatest(player_hp, 1), $2),
updated_at = now()
where player_id = $1
`, playerID, maxHP); err != nil {
		return fmt.Errorf("normalize playdeck hp after wardrobe change: %w", err)
	}

	return nil
}
