package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

const WardrobeCapacity = 100

type WardrobeStatus struct {
	Used     int `json:"used"`
	Capacity int `json:"capacity"`
}

type PlaydeckStatus struct {
	PlayerHP         int                       `json:"playerHp"`
	PlayerMaxHP      int                       `json:"playerMaxHp"`
	Attack           int                       `json:"attack"`
	Defense          int                       `json:"defense"`
	EquipmentPower   int                       `json:"equipmentPower"`
	Zone             PlaydeckZoneStatus        `json:"zone"`
	Enemy            PlaydeckEnemyStatus       `json:"enemy"`
	Wins             int64                     `json:"wins"`
	Losses           int64                     `json:"losses"`
	LastOutcome      string                    `json:"lastOutcome"`
	TimeoutTicks     int                       `json:"timeoutTicks"`
	Equipment        []EquipmentSlotStatus     `json:"equipment"`
	InventoryPreview []InventoryItemStatus     `json:"inventoryPreview"`
	CombatLog        []PlaydeckCombatLogStatus `json:"combatLog"`
}

type PlaydeckZoneStatus struct {
	ID           int    `json:"id"`
	Key          string `json:"key"`
	Name         string `json:"name"`
	MinLevel     int    `json:"minLevel"`
	SoftcapLevel int    `json:"softcapLevel"`
	Description  string `json:"description"`
}

type PlaydeckEnemyStatus struct {
	ID       int64  `json:"id"`
	Key      string `json:"key"`
	Name     string `json:"name"`
	Level    int    `json:"level"`
	HP       int    `json:"hp"`
	MaxHP    int    `json:"maxHp"`
	Attack   int    `json:"attack"`
	IsActive bool   `json:"isActive"`
}

type EquipmentSlotStatus struct {
	SlotKey      string `json:"slotKey"`
	DisplayName  string `json:"displayName"`
	AcceptsSlot  string `json:"acceptsSlot"`
	ItemID       int64  `json:"itemId"`
	ItemName     string `json:"itemName"`
	Rarity       string `json:"rarity"`
	PowerLevel   int    `json:"powerLevel"`
	AttackBonus  int    `json:"attackBonus"`
	DefenseBonus int    `json:"defenseBonus"`
	MaxHPBonus   int    `json:"maxHpBonus"`
}

type InventoryItemStatus struct {
	ID            int64  `json:"id"`
	ItemKey       string `json:"itemKey"`
	Name          string `json:"name"`
	ItemType      string `json:"itemType"`
	Rarity        string `json:"rarity"`
	EquipmentSlot string `json:"equipmentSlot"`
	Quantity      int    `json:"quantity"`
	EquippedSlot  string `json:"equippedSlot"`
	PowerLevel    int    `json:"powerLevel"`
	AttackBonus   int    `json:"attackBonus"`
	DefenseBonus  int    `json:"defenseBonus"`
	MaxHPBonus    int    `json:"maxHpBonus"`
}

type PlaydeckCombatLogStatus struct {
	ID                 int64     `json:"id"`
	EnemyName          string    `json:"enemyName"`
	EnemyLevel         int       `json:"enemyLevel"`
	Outcome            string    `json:"outcome"`
	PlayerDamage       int       `json:"playerDamage"`
	EnemyDamage        int       `json:"enemyDamage"`
	XPGained           int64     `json:"xpGained"`
	CreditsCentsGained int64     `json:"creditsCentsGained"`
	NibblesGained      int64     `json:"nibblesGained"`
	ItemName           string    `json:"itemName"`
	ItemQuantity       int       `json:"itemQuantity"`
	CreatedAt          time.Time `json:"createdAt"`
}

type PlaydeckEquipmentStats struct {
	AttackBonus    int
	DefenseBonus   int
	MaxHPBonus     int
	EquipmentPower int
}

func (s *Store) GetWardrobeStatus(ctx context.Context, playerID int64) (WardrobeStatus, error) {
	status := WardrobeStatus{
		Capacity: WardrobeCapacity,
	}

	if err := s.Pool.QueryRow(ctx, `
		select count(*)::int
		from player_inventory_items
		where player_id = $1
			and container = 'wardrobe'
	`, playerID).Scan(&status.Used); err != nil {
		return WardrobeStatus{}, fmt.Errorf("get wardrobe status: %w", err)
	}

	return status, nil
}

func (s *Store) GetPlaydeckStatus(ctx context.Context, playerID int64) (PlaydeckStatus, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return PlaydeckStatus{}, fmt.Errorf("begin get playdeck status: %w", err)
	}
	defer tx.Rollback(ctx)

	var level int
	var timeoutTicks int
	if err := tx.QueryRow(ctx, `
		select p.level,
			t.playdeck_timeout_ticks
		from players p
		join player_tick_state t on t.player_id = p.id
		where p.id = $1
	`, playerID).Scan(&level, &timeoutTicks); err != nil {
		return PlaydeckStatus{}, fmt.Errorf("load playdeck player basics: %w", err)
	}

	equipmentStats, err := loadPlaydeckEquipmentStatsTx(ctx, tx, playerID)
	if err != nil {
		return PlaydeckStatus{}, err
	}

	playerMaxHP := PlaydeckMaxHP(level, equipmentStats)

	if err := ensurePlaydeckStateTx(ctx, tx, playerID); err != nil {
		return PlaydeckStatus{}, err
	}

	var status PlaydeckStatus
	status.PlayerMaxHP = playerMaxHP
	status.Attack = PlaydeckAttack(level, equipmentStats)
	status.Defense = PlaydeckDefense(level, equipmentStats)
	status.EquipmentPower = equipmentStats.EquipmentPower
	status.TimeoutTicks = timeoutTicks

	err = tx.QueryRow(ctx, `
		select
			z.id,
			z.zone_key,
			z.name,
			z.min_level,
			z.softcap_level,
			z.description,
			e.id,
			e.enemy_key,
			e.name,
			s.enemy_level,
			s.enemy_hp,
			e.base_hp + round(e.hp_per_level * s.enemy_level)::integer,
			e.base_attack + round(e.attack_per_level * s.enemy_level)::integer,
			s.player_hp,
			s.wins,
			s.losses,
			s.last_outcome
		from playdeck_combat_state s
		join playdeck_zones z on z.id = s.zone_id
		join playdeck_enemies e on e.id = s.enemy_id
		where s.player_id = $1
	`, playerID).Scan(
		&status.Zone.ID,
		&status.Zone.Key,
		&status.Zone.Name,
		&status.Zone.MinLevel,
		&status.Zone.SoftcapLevel,
		&status.Zone.Description,
		&status.Enemy.ID,
		&status.Enemy.Key,
		&status.Enemy.Name,
		&status.Enemy.Level,
		&status.Enemy.HP,
		&status.Enemy.MaxHP,
		&status.Enemy.Attack,
		&status.PlayerHP,
		&status.Wins,
		&status.Losses,
		&status.LastOutcome,
	)
	if err != nil {
		return PlaydeckStatus{}, fmt.Errorf("load playdeck combat state: %w", err)
	}

	if status.PlayerHP <= 0 || status.PlayerHP > status.PlayerMaxHP {
		status.PlayerHP = status.PlayerMaxHP
	}

	status.Enemy.IsActive = status.Enemy.HP > 0

	status.Equipment, err = loadEquipmentSlotsTx(ctx, tx, playerID)
	if err != nil {
		return PlaydeckStatus{}, err
	}

	status.InventoryPreview, err = loadWardrobePreviewTx(ctx, tx, playerID, WardrobeCapacity)
	if err != nil {
		return PlaydeckStatus{}, err
	}

	status.CombatLog, err = loadPlaydeckCombatLogTx(ctx, tx, playerID, 10)
	if err != nil {
		return PlaydeckStatus{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return PlaydeckStatus{}, fmt.Errorf("commit get playdeck status: %w", err)
	}

	return status, nil
}

func ensurePlaydeckStateTx(ctx context.Context, tx pgx.Tx, playerID int64) error {
	var exists bool
	if err := tx.QueryRow(ctx, `
		select exists (
			select 1
			from playdeck_combat_state
			where player_id = $1
		)
	`, playerID).Scan(&exists); err != nil {
		return fmt.Errorf("check playdeck state: %w", err)
	}

	if exists {
		return nil
	}

	var level int
	if err := tx.QueryRow(ctx, `
		select level
		from players
		where id = $1
	`, playerID).Scan(&level); err != nil {
		return fmt.Errorf("load player level for playdeck state: %w", err)
	}

	if level < 1 {
		level = 1
	}

	equipmentStats, err := loadPlaydeckEquipmentStatsTx(ctx, tx, playerID)
	if err != nil {
		return err
	}

	var zoneID int
	var enemyID int64
	var enemyLevel int
	var enemyMaxHP int

	err = tx.QueryRow(ctx, `
		select
			z.id,
			e.id,
			greatest(e.min_level, least($1::int, e.max_level)),
			e.base_hp + round(e.hp_per_level * greatest(e.min_level, least($1::int, e.max_level)))::integer
		from playdeck_zones z
		join playdeck_enemies e on e.zone_id = z.id
		where z.zone_key = 'starter_deck'
		order by e.weight desc, e.id
		limit 1
	`, level).Scan(&zoneID, &enemyID, &enemyLevel, &enemyMaxHP)
	if err != nil {
		return fmt.Errorf("choose starter playdeck enemy: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		insert into playdeck_combat_state (
			player_id,
			zone_id,
			enemy_id,
			enemy_level,
			enemy_hp,
			player_hp
		)
		values ($1, $2, $3, $4, $5, $6)
		on conflict (player_id) do nothing
	`, playerID, zoneID, enemyID, enemyLevel, enemyMaxHP, PlaydeckMaxHP(level, equipmentStats)); err != nil {
		return fmt.Errorf("insert playdeck combat state: %w", err)
	}

	return nil
}

func loadPlaydeckEquipmentStatsTx(ctx context.Context, tx pgx.Tx, playerID int64) (PlaydeckEquipmentStats, error) {
	var stats PlaydeckEquipmentStats

	if err := tx.QueryRow(ctx, `
		select
			coalesce(sum(d.attack_bonus), 0)::int,
			coalesce(sum(d.defense_bonus), 0)::int,
			coalesce(sum(d.max_hp_bonus), 0)::int,
			coalesce(sum(d.power_level), 0)::int
		from player_inventory_items i
		join item_definitions d on d.id = i.item_definition_id
		where i.player_id = $1
			and i.equipped_slot is not null
	`, playerID).Scan(
		&stats.AttackBonus,
		&stats.DefenseBonus,
		&stats.MaxHPBonus,
		&stats.EquipmentPower,
	); err != nil {
		return PlaydeckEquipmentStats{}, fmt.Errorf("load playdeck equipment stats: %w", err)
	}

	return stats, nil
}

func loadEquipmentSlotsTx(ctx context.Context, tx pgx.Tx, playerID int64) ([]EquipmentSlotStatus, error) {
	rows, err := tx.Query(ctx, `
		select
			s.slot_key,
			s.display_name,
			s.accepts_slot,
			coalesce(i.id, 0),
			coalesce(d.name, ''),
			coalesce(d.rarity, ''),
			coalesce(d.power_level, 0),
			coalesce(d.attack_bonus, 0),
			coalesce(d.defense_bonus, 0),
			coalesce(d.max_hp_bonus, 0)
		from playdeck_equipment_slots s
		left join player_inventory_items i
			on i.player_id = $1
			and i.equipped_slot = s.slot_key
		left join item_definitions d on d.id = i.item_definition_id
		order by s.sort_order
	`, playerID)
	if err != nil {
		return nil, fmt.Errorf("query equipment slots: %w", err)
	}
	defer rows.Close()

	var slots []EquipmentSlotStatus
	for rows.Next() {
		var slot EquipmentSlotStatus
		if err := rows.Scan(
			&slot.SlotKey,
			&slot.DisplayName,
			&slot.AcceptsSlot,
			&slot.ItemID,
			&slot.ItemName,
			&slot.Rarity,
			&slot.PowerLevel,
			&slot.AttackBonus,
			&slot.DefenseBonus,
			&slot.MaxHPBonus,
		); err != nil {
			return nil, fmt.Errorf("scan equipment slot: %w", err)
		}

		slots = append(slots, slot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate equipment slots: %w", err)
	}

	return slots, nil
}

func loadWardrobePreviewTx(ctx context.Context, tx pgx.Tx, playerID int64, limit int) ([]InventoryItemStatus, error) {
	if limit < 1 {
		limit = 1
	}

	rows, err := tx.Query(ctx, `
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
			d.max_hp_bonus
		from player_inventory_items i
		join item_definitions d on d.id = i.item_definition_id
		where i.player_id = $1
			and i.container = 'wardrobe'
		order by
			case lower(d.rarity)
				when 'devastating' then 7
				when 'iconic' then 6
				when 'glam' then 5
				when 'trendy' then 4
				when 'chic' then 3
				when 'cute' then 2
				when 'basic' then 1
				else 0
			end desc,
			case when i.equipped_slot is not null then 1 else 0 end desc,
			d.name asc,
			i.acquired_at desc,
			i.id desc
		limit $2
	`, playerID, limit)
	if err != nil {
		return nil, fmt.Errorf("query wardrobe preview: %w", err)
	}
	defer rows.Close()

	var items []InventoryItemStatus
	for rows.Next() {
		var item InventoryItemStatus
		if err := rows.Scan(
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
		); err != nil {
			return nil, fmt.Errorf("scan wardrobe preview: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate wardrobe preview: %w", err)
	}

	return items, nil
}

func loadPlaydeckCombatLogTx(ctx context.Context, tx pgx.Tx, playerID int64, limit int) ([]PlaydeckCombatLogStatus, error) {
	if limit < 1 {
		limit = 1
	}

	rows, err := tx.Query(ctx, `
		select
			l.id,
			l.enemy_name,
			l.enemy_level,
			l.outcome,
			l.player_damage,
			l.enemy_damage,
			l.xp_gained,
			l.credits_cents_gained,
			l.nibbles_gained,
			coalesce(d.name, ''),
			l.item_quantity,
			l.created_at
		from playdeck_combat_log l
		left join item_definitions d on d.id = l.item_definition_id
		where l.player_id = $1
		order by l.created_at desc, l.id desc
		limit $2
	`, playerID, limit)
	if err != nil {
		return nil, fmt.Errorf("query playdeck combat log: %w", err)
	}
	defer rows.Close()

	var logs []PlaydeckCombatLogStatus
	for rows.Next() {
		var logEntry PlaydeckCombatLogStatus
		if err := rows.Scan(
			&logEntry.ID,
			&logEntry.EnemyName,
			&logEntry.EnemyLevel,
			&logEntry.Outcome,
			&logEntry.PlayerDamage,
			&logEntry.EnemyDamage,
			&logEntry.XPGained,
			&logEntry.CreditsCentsGained,
			&logEntry.NibblesGained,
			&logEntry.ItemName,
			&logEntry.ItemQuantity,
			&logEntry.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan playdeck combat log: %w", err)
		}

		logs = append(logs, logEntry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate playdeck combat log: %w", err)
	}

	return logs, nil
}

func PlaydeckMaxHP(level int, equipment PlaydeckEquipmentStats) int {
	if level < 1 {
		level = 1
	}

	return 100 + level*4 + equipment.MaxHPBonus
}

func PlaydeckAttack(level int, equipment PlaydeckEquipmentStats) int {
	if level < 1 {
		level = 1
	}

	return 8 + level*2 + equipment.AttackBonus
}

func PlaydeckDefense(level int, equipment PlaydeckEquipmentStats) int {
	if level < 1 {
		level = 1
	}

	return level/2 + equipment.DefenseBonus
}
