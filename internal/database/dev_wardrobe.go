package database

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

type DevWardrobeSpawnResult struct {
	ItemID int64              `json:"itemId"`
	Detail WardrobeItemDetail `json:"detail"`
}

type devWardrobeSlot struct {
	Key   string
	Label string
	Nouns []string
}

type devRolledWardrobeStatLine struct {
	Source    string
	AffixKey  string
	StatKey   string
	Value     float64
	SortOrder int
}

var devWardrobeSlots = []devWardrobeSlot{
	{Key: "top", Label: "Top", Nouns: []string{"Crop Top", "Blouse", "Cami", "Sweater Top"}},
	{Key: "bottom", Label: "Bottom", Nouns: []string{"Skirt", "Shorts", "Leggings", "Pants"}},
	{Key: "dress", Label: "Dress", Nouns: []string{"Dress", "Mini Dress", "Outfit", "Jumper Dress"}},
	{Key: "footwear", Label: "Footwear", Nouns: []string{"Boots", "Sneakers", "Heels", "Ankle Boots"}},
	{Key: "outerwear", Label: "Outerwear", Nouns: []string{"Jacket", "Cardigan", "Coat", "Bolero"}},
	{Key: "necklace", Label: "Necklace", Nouns: []string{"Choker", "Pendant", "Necklace", "Charm Chain"}},
	{Key: "accessory", Label: "Accessory", Nouns: []string{"Hair Pin", "Bracelet", "Ribbon", "Charm"}},
}

var devWardrobeRarities = []struct {
	Key    string
	Weight int
}{
	{Key: "basic", Weight: 42},
	{Key: "cute", Weight: 25},
	{Key: "chic", Weight: 15},
	{Key: "trendy", Weight: 9},
	{Key: "glam", Weight: 5},
	{Key: "iconic", Weight: 3},
	{Key: "devastating", Weight: 1},
}

var devWardrobeNameBits = []string{
	"Pixel",
	"Velvet",
	"Bubblegum",
	"Moonlit",
	"Cozy",
	"Neon",
	"Ribbon",
	"Starlit",
	"Cafe",
	"Chaos",
	"Glitter",
	"Dreamy",
}

var devWardrobePrefixPool = []string{
	"playdeck_xp_percent",
	"work_xp_percent",
	"global_xp_percent",
	"work_resources_percent",
	"drop_rate_percent",
	"credit_rate_percent",
	"ingredient_quality_percent",
}

var devWardrobeImplicitPool = []string{
	"max_health_percent",
	"attack_percent",
	"attack_speed_percent",
	"beauty",
	"glamor",
	"crit_rate_percent",
	"crit_damage_percent",
	"charm",
	"humor",
	"targeting_percent",
	"dodge_percent",
	"recovery",
}

var devWardrobeSuffixPoolsBySlot = map[string][]string{
	"top": {
		"max_health_percent",
		"attack_percent",
		"beauty",
		"glamor",
		"charm",
		"recovery",
	},
	"bottom": {
		"max_health_percent",
		"beauty",
		"glamor",
		"dodge_percent",
		"humor",
		"recovery",
	},
	"dress": {
		"max_health_percent",
		"attack_percent",
		"beauty",
		"glamor",
		"charm",
		"humor",
		"recovery",
	},
	"footwear": {
		"attack_speed_percent",
		"targeting_percent",
		"dodge_percent",
		"humor",
		"recovery",
	},
	"outerwear": {
		"max_health_percent",
		"beauty",
		"glamor",
		"dodge_percent",
		"recovery",
	},
	"necklace": {
		"attack_percent",
		"crit_rate_percent",
		"crit_damage_percent",
		"charm",
		"targeting_percent",
	},
	"accessory": {
		"attack_percent",
		"attack_speed_percent",
		"crit_rate_percent",
		"crit_damage_percent",
		"charm",
		"humor",
		"targeting_percent",
		"dodge_percent",
		"recovery",
	},
}

func (s *Store) SpawnDevRandomWardrobeItem(ctx context.Context) (DevWardrobeSpawnResult, error) {
	playerID, err := s.DevPlayerID(ctx)
	if err != nil {
		return DevWardrobeSpawnResult{}, err
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return DevWardrobeSpawnResult{}, fmt.Errorf("begin spawn dev wardrobe item: %w", err)
	}
	defer tx.Rollback(ctx)

	var playerLevel int
	if err := tx.QueryRow(ctx, `
		select level
		from players
		where id = $1
	`, playerID).Scan(&playerLevel); err != nil {
		return DevWardrobeSpawnResult{}, fmt.Errorf("load player level for dev wardrobe item: %w", err)
	}

	var wardrobeUsed int
	if err := tx.QueryRow(ctx, `
		select count(*)::int
		from player_inventory_items
		where player_id = $1
			and container = 'wardrobe'
	`, playerID).Scan(&wardrobeUsed); err != nil {
		return DevWardrobeSpawnResult{}, fmt.Errorf("count wardrobe items: %w", err)
	}

	if wardrobeUsed >= WardrobeCapacity {
		return DevWardrobeSpawnResult{}, fmt.Errorf("wardrobe is full")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	slot := devRandomWardrobeSlot(rng)
	rarity := devRandomWardrobeRarity(rng)
	powerLevel := devRandomWardrobePowerLevel(rng, playerLevel, rarity)
	itemName := devRandomWardrobeItemName(rng, slot)
	itemKey := fmt.Sprintf(
		"dev_%s_%s_%d_%05d",
		slot.Key,
		rarity,
		time.Now().UnixNano(),
		rng.Intn(100000),
	)

	var itemDefinitionID int64
	if err := tx.QueryRow(ctx, `
		insert into item_definitions (
			item_key,
			name,
			item_type,
			rarity,
			equipment_slot,
			stackable,
			max_stack,
			power_level,
			attack_bonus,
			defense_bonus,
			max_hp_bonus,
			value_cents,
			description
		)
		values ($1, $2, 'gear', $3, $4, false, 1, $5, 0, 0, 0, 0, $6)
		returning id
	`,
		itemKey,
		itemName,
		rarity,
		slot.Key,
		powerLevel,
		fmt.Sprintf("Dev-spawned %s %s for testing wardrobe stat lines.", rarity, slot.Label),
	).Scan(&itemDefinitionID); err != nil {
		return DevWardrobeSpawnResult{}, fmt.Errorf("insert dev wardrobe item definition: %w", err)
	}

	tailoringMax := devWardrobeTailoringMax(rng, rarity, powerLevel)

	var itemID int64
	if err := tx.QueryRow(ctx, `
		insert into player_inventory_items (
			player_id,
			item_definition_id,
			container,
			quantity,
			tailoring_current,
			tailoring_max
		)
		values ($1, $2, 'wardrobe', 1, 0, $3)
		returning id
	`, playerID, itemDefinitionID, tailoringMax).Scan(&itemID); err != nil {
		return DevWardrobeSpawnResult{}, fmt.Errorf("insert dev wardrobe item: %w", err)
	}

	statLines := devRollWardrobeStatLines(rng, rarity, slot.Key, powerLevel)
	for _, line := range statLines {
		if _, err := tx.Exec(ctx, `
			insert into player_inventory_item_stat_lines (
				player_inventory_item_id,
				stat_source,
				affix_key,
				stat_key,
				value,
				sort_order
			)
			values ($1, $2, $3, $4, $5, $6)
		`,
			itemID,
			line.Source,
			line.AffixKey,
			line.StatKey,
			line.Value,
			line.SortOrder,
		); err != nil {
			return DevWardrobeSpawnResult{}, fmt.Errorf("insert dev wardrobe stat line: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return DevWardrobeSpawnResult{}, fmt.Errorf("commit spawn dev wardrobe item: %w", err)
	}

	detail, err := s.GetWardrobeItemDetail(ctx, playerID, itemID, "")
	if err != nil {
		return DevWardrobeSpawnResult{}, err
	}

	return DevWardrobeSpawnResult{
		ItemID: itemID,
		Detail: detail,
	}, nil
}

func devRandomWardrobeSlot(rng *rand.Rand) devWardrobeSlot {
	return devWardrobeSlots[rng.Intn(len(devWardrobeSlots))]
}

func devRandomWardrobeRarity(rng *rand.Rand) string {
	totalWeight := 0
	for _, rarity := range devWardrobeRarities {
		totalWeight += rarity.Weight
	}

	roll := rng.Intn(totalWeight)
	for _, rarity := range devWardrobeRarities {
		if roll < rarity.Weight {
			return rarity.Key
		}

		roll -= rarity.Weight
	}

	return "basic"
}

func devRandomWardrobePowerLevel(rng *rand.Rand, playerLevel int, rarity string) int {
	if playerLevel < 1 {
		playerLevel = 1
	}

	rarityBonus := map[string]int{
		"basic":       0,
		"cute":        1,
		"chic":        2,
		"trendy":      4,
		"glam":        7,
		"iconic":      10,
		"devastating": 15,
	}[rarity]

	level := playerLevel + rng.Intn(8) - 2 + rarityBonus
	if level < 1 {
		return 1
	}

	return level
}

func devRandomWardrobeItemName(rng *rand.Rand, slot devWardrobeSlot) string {
	prefix := devWardrobeNameBits[rng.Intn(len(devWardrobeNameBits))]
	noun := slot.Nouns[rng.Intn(len(slot.Nouns))]

	return fmt.Sprintf("%s %s", prefix, noun)
}

func devWardrobeTailoringMax(rng *rand.Rand, rarity string, powerLevel int) int {
	base := 8 + powerLevel*2
	rarityBonus := map[string]int{
		"basic":       0,
		"cute":        8,
		"chic":        16,
		"trendy":      26,
		"glam":        38,
		"iconic":      52,
		"devastating": 70,
	}[rarity]

	return base + rarityBonus + rng.Intn(8)
}

func devRollWardrobeStatLines(rng *rand.Rand, rarity string, slotKey string, powerLevel int) []devRolledWardrobeStatLine {
	lines := make([]devRolledWardrobeStatLine, 0, 8)

	implicitStat := devWardrobeImplicitPool[rng.Intn(len(devWardrobeImplicitPool))]
	lines = append(lines, devRolledWardrobeStatLine{
		Source:    "implicit",
		AffixKey:  fmt.Sprintf("implicit:%s", implicitStat),
		StatKey:   implicitStat,
		Value:     devRollWardrobeStatValue(rng, implicitStat, rarity, powerLevel),
		SortOrder: 10,
	})

	affixCount := devWardrobeAffixCount(rarity)
	maxPrefixCount := minInt(3, affixCount)
	minPrefixCount := maxInt(0, affixCount-3)
	prefixCount := minPrefixCount + rng.Intn(maxPrefixCount-minPrefixCount+1)
	suffixCount := affixCount - prefixCount

	usedAffixStats := make(map[string]bool, affixCount)

	for index := 0; index < prefixCount; index++ {
		statKey, ok := devPickUnusedStatKey(rng, devWardrobePrefixPool, usedAffixStats)
		if !ok {
			break
		}

		usedAffixStats[statKey] = true
		lines = append(lines, devRolledWardrobeStatLine{
			Source:    "prefix",
			AffixKey:  fmt.Sprintf("prefix:%s", statKey),
			StatKey:   statKey,
			Value:     devRollWardrobeStatValue(rng, statKey, rarity, powerLevel),
			SortOrder: 100 + index,
		})
	}

	suffixPool := devWardrobeSuffixPoolsBySlot[strings.ToLower(slotKey)]
	if len(suffixPool) == 0 {
		suffixPool = devWardrobeImplicitPool
	}

	for index := 0; index < suffixCount; index++ {
		statKey, ok := devPickUnusedStatKey(rng, suffixPool, usedAffixStats)
		if !ok {
			break
		}

		usedAffixStats[statKey] = true
		lines = append(lines, devRolledWardrobeStatLine{
			Source:    "suffix",
			AffixKey:  fmt.Sprintf("suffix:%s", statKey),
			StatKey:   statKey,
			Value:     devRollWardrobeStatValue(rng, statKey, rarity, powerLevel),
			SortOrder: 200 + index,
		})
	}

	return lines
}

func devWardrobeAffixCount(rarity string) int {
	switch rarity {
	case "basic":
		return 1
	case "cute":
		return 2
	case "chic":
		return 3
	case "trendy":
		return 4
	case "glam":
		return 5
	case "iconic":
		return 6
	case "devastating":
		return 6
	default:
		return 1
	}
}

func devPickUnusedStatKey(rng *rand.Rand, pool []string, used map[string]bool) (string, bool) {
	candidates := make([]string, 0, len(pool))

	for _, statKey := range pool {
		if !used[statKey] {
			candidates = append(candidates, statKey)
		}
	}

	if len(candidates) == 0 {
		return "", false
	}

	return candidates[rng.Intn(len(candidates))], true
}

func devRollWardrobeStatValue(rng *rand.Rand, statKey string, rarity string, powerLevel int) float64 {
	if powerLevel < 1 {
		powerLevel = 1
	}

	rarityMultiplier := devWardrobeRarityMultiplier(rarity)
	power := math.Sqrt(float64(powerLevel))

	switch statKey {
	case "global_xp_percent":
		return devRoundStat(rollFloat(rng, 1.5, 4.5) * rarityMultiplier)
	case "playdeck_xp_percent", "work_xp_percent":
		return devRoundStat((rollFloat(rng, 2.0, 7.0) + power*0.2) * rarityMultiplier)
	case "work_resources_percent", "credit_rate_percent":
		return devRoundStat((rollFloat(rng, 2.5, 8.5) + power*0.25) * rarityMultiplier)
	case "drop_rate_percent":
		return devRoundStat((rollFloat(rng, 3.0, 10.0) + power*0.3) * rarityMultiplier)
	case "ingredient_quality_percent":
		return devRoundStat((rollFloat(rng, 2.0, 8.0) + power*0.25) * rarityMultiplier)
	case "max_health_percent", "attack_percent", "targeting_percent", "dodge_percent":
		return devRoundStat((rollFloat(rng, 2.0, 7.5) + power*0.18) * rarityMultiplier)
	case "attack_speed_percent", "crit_rate_percent":
		return devRoundStat((rollFloat(rng, 1.0, 4.0) + power*0.08) * rarityMultiplier)
	case "crit_damage_percent":
		return devRoundStat((rollFloat(rng, 4.0, 12.0) + power*0.25) * rarityMultiplier)
	case "beauty":
		return devRoundStat((rollFloat(rng, 2.0, 8.0) + power*0.75) * rarityMultiplier)
	case "glamor":
		return devRoundStat((rollFloat(rng, 8.0, 24.0) + power*4.0) * rarityMultiplier)
	case "charm", "humor", "recovery":
		return devRoundStat((rollFloat(rng, 1.0, 5.0) + power*0.6) * rarityMultiplier)
	default:
		return devRoundStat((rollFloat(rng, 1.0, 5.0) + power*0.2) * rarityMultiplier)
	}
}

func devWardrobeRarityMultiplier(rarity string) float64 {
	switch rarity {
	case "basic":
		return 1.0
	case "cute":
		return 1.25
	case "chic":
		return 1.55
	case "trendy":
		return 1.95
	case "glam":
		return 2.45
	case "iconic":
		return 3.1
	case "devastating":
		return 4.0
	default:
		return 1.0
	}
}

func rollFloat(rng *rand.Rand, minValue float64, maxValue float64) float64 {
	return minValue + rng.Float64()*(maxValue-minValue)
}

func devRoundStat(value float64) float64 {
	return math.Round(value*10) / 10
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}

	return b
}
