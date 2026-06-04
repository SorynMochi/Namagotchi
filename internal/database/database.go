package database

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	TickSeconds       = 5
	SyncXPPerTick     = int64(10)
	ActivityXPPerTick = int64(10)
	MaxOfflineTicks   = int64(8640)
)

type Store struct {
	Pool *pgxpool.Pool
}

type PlayerStatus struct {
	Player     Player          `json:"player"`
	Companion  CompanionState  `json:"companion"`
	Resources  PlayerResources `json:"resources"`
	Activities ActivitySkills  `json:"activities"`
	Tick       TickState       `json:"tick"`
}

type Player struct {
	ID            int64  `json:"id"`
	DisplayName   string `json:"displayName"`
	Level         int    `json:"level"`
	TotalXP       int64  `json:"totalXp"`
	XPIntoLevel   int64  `json:"xpIntoLevel"`
	XPToNext      int64  `json:"xpToNext"`
	CurrencyCents int64  `json:"currencyCents"`
	CreditsCents  int64  `json:"creditsCents"`
	Nibbles       int64  `json:"nibbles"`
	NamiCoin      int64  `json:"namiCoin"`
}

type CompanionState struct {
	CompanionName     string    `json:"name"`
	MoodScore         float64   `json:"moodScore"`
	Satiety           int       `json:"satiety"`
	Connection        int       `json:"connection"`
	Energy            int       `json:"energy"`
	Comfort           int       `json:"comfort"`
	Playfulness       int       `json:"playfulness"`
	Inspiration       int       `json:"inspiration"`
	Cleanliness       int       `json:"cleanliness"`
	Status            string    `json:"status"`
	LastInteractionAt time.Time `json:"lastInteractionAt"`
	Level             int       `json:"level"`
	TotalXP           int64     `json:"totalXp"`
	XPIntoLevel       int64     `json:"xpIntoLevel"`
	XPToNext          int64     `json:"xpToNext"`
	LastXPGained      int64     `json:"lastXpGained"`
	LastAction        string    `json:"lastAction"`
	MoodLabel         string    `json:"moodLabel"`
	PrimaryNeed       string    `json:"primaryNeed"`
	Caption           string    `json:"caption"`
	SuggestedAction   string    `json:"suggestedAction"`
}

type PlayerResources struct {
	Fans        int64 `json:"fans"`
	Memes       int64 `json:"memes"`
	LostItems   int64 `json:"lostItems"`
	Confidence  int64 `json:"confidence"`
	Receipts    int64 `json:"receipts"`
	Patterns    int64 `json:"patterns"`
	GlitchDrops int64 `json:"glitchDrops"`
}

type ActivitySkill struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	TotalXP     int64  `json:"totalXp"`
	XPIntoLevel int64  `json:"xpIntoLevel"`
	XPToNext    int64  `json:"xpToNext"`
}

type ActivitySkills struct {
	Streaming     ActivitySkill `json:"streaming"`
	DoomScrolling ActivitySkill `json:"doomScrolling"`
	Cleaning      ActivitySkill `json:"cleaning"`
	Exercising    ActivitySkill `json:"exercising"`
	Shopping      ActivitySkill `json:"shopping"`
	Designing     ActivitySkill `json:"designing"`
}

type TickState struct {
	PlaydeckEnabled           bool      `json:"playdeckEnabled"`
	PlaydeckZoneID            int       `json:"playdeckZoneId"`
	PlaydeckZoneName          string    `json:"playdeckZoneName"`
	PlaydeckStreak            int64     `json:"playdeckStreak"`
	PlaydeckTimeoutTicks      int       `json:"playdeckTimeoutTicks"`
	ActiveGatheringTask       string    `json:"activeGatheringTask"`
	ActiveGatheringName       string    `json:"activeGatheringName"`
	ActiveGatheringOutput     string    `json:"activeGatheringOutput"`
	ActiveActivityLevel       int       `json:"activeActivityLevel"`
	ActiveActivityXPIntoLevel int64     `json:"activeActivityXpIntoLevel"`
	ActiveActivityXPToNext    int64     `json:"activeActivityXpToNext"`
	GatheringRemainder        float64   `json:"gatheringRemainder"`
	ResourcePerTick           float64   `json:"resourcePerTick"`
	ResourcePerTickDisplay    int64     `json:"resourcePerTickDisplay"`
	LastTickAt                time.Time `json:"lastTickAt"`
	NextTickAt                time.Time `json:"nextTickAt"`
	SecondsUntilNextTick      int       `json:"secondsUntilNextTick"`
}

type TickResult struct {
	OK                   bool   `json:"ok"`
	TicksProcessed       int64  `json:"ticksProcessed"`
	SyncXPGained         int64  `json:"syncXpGained"`
	CreditsCentsGained   int64  `json:"creditsCentsGained"`
	NibblesGained        int64  `json:"nibblesGained"`
	ResourceName         string `json:"resourceName"`
	ResourceAmountGained int64  `json:"resourceAmountGained"`
	ActivityName         string `json:"activityName"`
	ActivityXPGained     int64  `json:"activityXpGained"`
	ActivityLevelUps     int    `json:"activityLevelUps"`
	ActivityCurrentLevel int    `json:"activityCurrentLevel"`
	ActivityXPIntoLevel  int64  `json:"activityXpIntoLevel"`
	ActivityXPToNext     int64  `json:"activityXpToNext"`
	LevelUps             int    `json:"levelUps"`
	CurrentLevel         int    `json:"currentLevel"`
	XPIntoLevel          int64  `json:"xpIntoLevel"`
	XPToNext             int64  `json:"xpToNext"`
	Message              string `json:"message"`
}

type CareActionResult struct {
	OK           bool           `json:"ok"`
	Action       string         `json:"action"`
	ActionName   string         `json:"actionName"`
	XPGained     int64          `json:"xpGained"`
	LevelUps     int            `json:"levelUps"`
	CurrentLevel int            `json:"currentLevel"`
	XPIntoLevel  int64          `json:"xpIntoLevel"`
	XPToNext     int64          `json:"xpToNext"`
	Companion    CompanionState `json:"companion"`
	Message      string         `json:"message"`
}

type CareActionRule struct {
	Key         string
	Name        string
	Satiety     int
	Connection  int
	Energy      int
	Comfort     int
	Playfulness int
	Inspiration int
	Cleanliness int
	SleepAction bool
	WakeAction  bool
}

func CareActionByKey(action string) (CareActionRule, bool) {
	switch strings.TrimSpace(strings.ToLower(action)) {
	case "meal":
		return CareActionRule{Key: "meal", Name: "Meal", Satiety: 30, Comfort: 6, Cleanliness: -2}, true
	case "snack":
		return CareActionRule{Key: "snack", Name: "Snack", Satiety: 12, Comfort: 3, Playfulness: 4, Cleanliness: -1}, true
	case "drink":
		return CareActionRule{Key: "drink", Name: "Drink", Satiety: 3, Energy: 8, Comfort: 4}, true
	case "cuddle":
		return CareActionRule{Key: "cuddle", Name: "Cuddle", Connection: 18, Energy: -3, Comfort: 12}, true
	case "play":
		return CareActionRule{Key: "play", Name: "Play", Connection: 8, Energy: -8, Playfulness: 20, Cleanliness: -3}, true
	case "write_together":
		return CareActionRule{Key: "write_together", Name: "Write Together", Connection: 6, Energy: -5, Inspiration: 20}, true
	case "read_together":
		return CareActionRule{Key: "read_together", Name: "Read Together", Connection: 4, Energy: -3, Comfort: 12, Inspiration: 8}, true
	case "boop":
		return CareActionRule{Key: "boop", Name: "Boop", Connection: 4, Playfulness: 8}, true
	case "nap":
		return CareActionRule{Key: "nap", Name: "Nap", Energy: 15, Comfort: 5}, true
	case "bath":
		return CareActionRule{Key: "bath", Name: "Bath", Connection: 8, Energy: -6, Comfort: 10, Cleanliness: 45}, true
	case "freshen_up":
		return CareActionRule{Key: "freshen_up", Name: "Freshen Up", Comfort: 8, Cleanliness: 20}, true
	case "put_to_bed":
		return CareActionRule{Key: "put_to_bed", Name: "Put To Bed", SleepAction: true}, true
	case "wake_up":
		return CareActionRule{Key: "wake_up", Name: "Wake Up", WakeAction: true}, true
	default:
		return CareActionRule{}, false
	}
}

func clampCareStat(value int) int {
	if value < 0 {
		return 0
	}

	if value > 100 {
		return 100
	}

	return value
}

func usefulCareXP(current int, delta int) int64 {
	if delta == 0 {
		return 0
	}

	next := clampCareStat(current + delta)

	if delta > 0 {
		return int64(next - current)
	}

	return int64(-(current - next))
}

func (s *Store) ApplyDevCareAction(ctx context.Context, action string) (*CareActionResult, error) {
	rule, ok := CareActionByKey(action)
	if !ok {
		return nil, fmt.Errorf("invalid care action: %s", action)
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin care action: %w", err)
	}
	defer tx.Rollback(ctx)

	var playerID int64
	var companion CompanionState

	err = tx.QueryRow(ctx, `
		select
			p.id,
			c.companion_name,
			c.level,
			c.total_xp,
			c.xp_into_level,
			c.mood_score::float8,
			c.satiety,
			c.connection,
			c.energy,
			c.comfort,
			c.playfulness,
			c.inspiration,
			c.cleanliness,
			c.status,
			c.last_interaction_at,
			c.last_xp_gained,
			c.last_action
		from players p
		join companion_states c on c.player_id = p.id
		where p.display_name = 'Soryn'
		for update
	`).Scan(
		&playerID,
		&companion.CompanionName,
		&companion.Level,
		&companion.TotalXP,
		&companion.XPIntoLevel,
		&companion.MoodScore,
		&companion.Satiety,
		&companion.Connection,
		&companion.Energy,
		&companion.Comfort,
		&companion.Playfulness,
		&companion.Inspiration,
		&companion.Cleanliness,
		&companion.Status,
		&companion.LastInteractionAt,
		&companion.LastXPGained,
		&companion.LastAction,
	)

	if err != nil {
		return nil, fmt.Errorf("load companion for care action: %w", err)
	}

	xpGained := int64(0)

	if rule.SleepAction {
		companion.Status = "sleeping"
	} else if rule.WakeAction {
		companion.Status = "awake"
	} else {
		xpGained += usefulCareXP(companion.Satiety, rule.Satiety)
		xpGained += usefulCareXP(companion.Connection, rule.Connection)
		xpGained += usefulCareXP(companion.Energy, rule.Energy)
		xpGained += usefulCareXP(companion.Comfort, rule.Comfort)
		xpGained += usefulCareXP(companion.Playfulness, rule.Playfulness)
		xpGained += usefulCareXP(companion.Inspiration, rule.Inspiration)
		xpGained += usefulCareXP(companion.Cleanliness, rule.Cleanliness)

		if xpGained < 0 {
			xpGained = 0
		}

		companion.Satiety = clampCareStat(companion.Satiety + rule.Satiety)
		companion.Connection = clampCareStat(companion.Connection + rule.Connection)
		companion.Energy = clampCareStat(companion.Energy + rule.Energy)
		companion.Comfort = clampCareStat(companion.Comfort + rule.Comfort)
		companion.Playfulness = clampCareStat(companion.Playfulness + rule.Playfulness)
		companion.Inspiration = clampCareStat(companion.Inspiration + rule.Inspiration)
		companion.Cleanliness = clampCareStat(companion.Cleanliness + rule.Cleanliness)
		companion.Status = "awake"
	}

	companion.TotalXP += xpGained
	companion.XPIntoLevel += xpGained

	levelUps := 0
	for companion.XPIntoLevel >= NamiXPToNextLevel(companion.Level) {
		companion.XPIntoLevel -= NamiXPToNextLevel(companion.Level)
		companion.Level++
		levelUps++
	}

	companion.MoodScore = NamiMoodScore(companion)
	companion.LastXPGained = xpGained
	companion.LastAction = rule.Name

	_, err = tx.Exec(ctx, `
		update companion_states
		set
			level = $1,
			total_xp = $2,
			xp_into_level = $3,
			last_xp_gained = $4,
			last_action = $5,
			mood_score = $6,
			satiety = $7,
			connection = $8,
			energy = $9,
			comfort = $10,
			playfulness = $11,
			inspiration = $12,
			cleanliness = $13,
			status = $14,
			last_interaction_at = now(),
			updated_at = now()
		where player_id = $15
	`,
		companion.Level,
		companion.TotalXP,
		companion.XPIntoLevel,
		companion.LastXPGained,
		companion.LastAction,
		companion.MoodScore,
		companion.Satiety,
		companion.Connection,
		companion.Energy,
		companion.Comfort,
		companion.Playfulness,
		companion.Inspiration,
		companion.Cleanliness,
		companion.Status,
		playerID,
	)

	if err != nil {
		return nil, fmt.Errorf("update companion after care action: %w", err)
	}

	_, err = tx.Exec(ctx, `
		insert into activity_log (player_id, event_type, message)
		values ($1, 'care_action', $2)
	`, playerID, fmt.Sprintf("Nami-chan received care action: %s (+%d XP).", rule.Name, xpGained))

	if err != nil {
		return nil, fmt.Errorf("insert care action log: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit care action: %w", err)
	}

	companion.XPToNext = NamiXPToNextLevel(companion.Level)
	companion.MoodLabel = NamiMoodLabel(companion.MoodScore)
	companion.PrimaryNeed = NamiPrimaryNeed(companion)
	companion.Caption = NamiCaption(companion)
	companion.SuggestedAction = NamiSuggestedAction(companion)

	return &CareActionResult{
		OK:           true,
		Action:       rule.Key,
		ActionName:   rule.Name,
		XPGained:     xpGained,
		LevelUps:     levelUps,
		CurrentLevel: companion.Level,
		XPIntoLevel:  companion.XPIntoLevel,
		XPToNext:     companion.XPToNext,
		Companion:    companion,
		Message:      fmt.Sprintf("%s complete. Nami-chan gained %d XP.", rule.Name, xpGained),
	}, nil
}

func Connect(ctx context.Context, databaseURL string) (*Store, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create PostgreSQL pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping PostgreSQL: %w", err)
	}

	return &Store{Pool: pool}, nil
}

func (s *Store) Close() {
	s.Pool.Close()
}

func (s *Store) RunMigrations(ctx context.Context, migrationsDir string) error {
	if _, err := s.Pool.Exec(ctx, `
		create table if not exists schema_migrations (
			version text primary key,
			applied_at timestamptz not null default now()
		);
	`); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".sql") {
			migrationFiles = append(migrationFiles, name)
		}
	}

	sort.Strings(migrationFiles)

	for _, fileName := range migrationFiles {
		alreadyApplied, err := s.migrationApplied(ctx, fileName)
		if err != nil {
			return err
		}

		if alreadyApplied {
			continue
		}

		path := filepath.Join(migrationsDir, fileName)
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", fileName, err)
		}

		tx, err := s.Pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin migration %s: %w", fileName, err)
		}

		if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("execute migration %s: %w", fileName, err)
		}

		if _, err := tx.Exec(ctx, `
			insert into schema_migrations (version)
			values ($1)
		`, fileName); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("record migration %s: %w", fileName, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit migration %s: %w", fileName, err)
		}
	}

	return nil
}

func (s *Store) migrationApplied(ctx context.Context, version string) (bool, error) {
	var exists bool
	err := s.Pool.QueryRow(ctx, `
		select exists (
			select 1
			from schema_migrations
			where version = $1
		)
	`, version).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check migration %s: %w", version, err)
	}

	return exists, nil
}

func (s *Store) SeedDevPlayer(ctx context.Context) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin dev seed: %w", err)
	}

	defer tx.Rollback(ctx)

	var playerID int64
	if err := tx.QueryRow(ctx, `
		insert into players (display_name)
		values ('Soryn')
		on conflict (display_name) do update
		set updated_at = now()
		returning id
	`).Scan(&playerID); err != nil {
		return fmt.Errorf("upsert dev player: %w", err)
	}

	for _, activityKey := range GatheringTaskKeys() {
		if _, err := tx.Exec(ctx, `
			insert into player_activity_skills (player_id, activity_key)
			values ($1, $2)
			on conflict (player_id, activity_key) do update
			set updated_at = now()
		`, playerID, activityKey); err != nil {
			return fmt.Errorf("upsert activity skill %s: %w", activityKey, err)
		}
	}

	if _, err := tx.Exec(ctx, `
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
			last_interaction_at
		)
		values (
			$1,
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
			now()
		)
		on conflict (player_id) do update
		set
			companion_name = excluded.companion_name,
			mood_score = excluded.mood_score,
			satiety = excluded.satiety,
			connection = excluded.connection,
			energy = excluded.energy,
			comfort = excluded.comfort,
			playfulness = excluded.playfulness,
			inspiration = excluded.inspiration,
			cleanliness = excluded.cleanliness,
			status = excluded.status,
			last_interaction_at = excluded.last_interaction_at,
			updated_at = now()
	`, playerID); err != nil {
		return fmt.Errorf("upsert companion state: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		insert into player_resources (player_id)
		values ($1)
		on conflict (player_id) do update
		set updated_at = now()
	`, playerID); err != nil {
		return fmt.Errorf("upsert player resources: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		insert into player_tick_state (player_id)
		values ($1)
		on conflict (player_id) do update
		set updated_at = now()
	`, playerID); err != nil {
		return fmt.Errorf("upsert player tick state: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		insert into activity_log (player_id, event_type, message)
		values ($1, 'dev_seed', 'Dev player Soryn and Nami-chan were seeded.')
	`, playerID); err != nil {
		return fmt.Errorf("insert activity log: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit dev seed: %w", err)
	}

	return nil
}

func (s *Store) GetDevPlayerStatus(ctx context.Context) (*PlayerStatus, error) {
	var status PlayerStatus

	err := s.Pool.QueryRow(ctx, `
		select
			p.id,
			p.display_name,
			p.level,
			p.total_xp,
			p.xp_into_level,
			p.currency_cents,
			p.nibbles,
			p.namicoin,
			c.companion_name,
			c.level,
			c.total_xp,
			c.xp_into_level,
			c.last_xp_gained,
			c.last_action,
			c.mood_score::float8,
			c.satiety,
			c.connection,
			c.energy,
			c.comfort,
			c.playfulness,
			c.inspiration,
			c.cleanliness,
			c.status,
			c.last_interaction_at,
			r.fans,
			r.memes,
			r.lost_items,
			r.confidence,
			r.receipts,
			r.patterns,
			r.glitch_drops,
			t.playdeck_enabled,
			t.playdeck_zone_id,
			t.playdeck_streak,
			t.playdeck_timeout_ticks,
			t.active_gathering_task,
			t.gathering_remainder,
			t.last_tick_at
		from players p
		join companion_states c on c.player_id = p.id
		join player_resources r on r.player_id = p.id
		join player_tick_state t on t.player_id = p.id
		where p.display_name = 'Soryn'
	`).Scan(
		&status.Player.ID,
		&status.Player.DisplayName,
		&status.Player.Level,
		&status.Player.TotalXP,
		&status.Player.XPIntoLevel,
		&status.Player.CurrencyCents,
		&status.Player.Nibbles,
		&status.Player.NamiCoin,
		&status.Companion.CompanionName,
		&status.Companion.Level,
		&status.Companion.TotalXP,
		&status.Companion.XPIntoLevel,
		&status.Companion.LastXPGained,
		&status.Companion.LastAction,
		&status.Companion.MoodScore,
		&status.Companion.Satiety,
		&status.Companion.Connection,
		&status.Companion.Energy,
		&status.Companion.Comfort,
		&status.Companion.Playfulness,
		&status.Companion.Inspiration,
		&status.Companion.Cleanliness,
		&status.Companion.Status,
		&status.Companion.LastInteractionAt,
		&status.Resources.Fans,
		&status.Resources.Memes,
		&status.Resources.LostItems,
		&status.Resources.Confidence,
		&status.Resources.Receipts,
		&status.Resources.Patterns,
		&status.Resources.GlitchDrops,
		&status.Tick.PlaydeckEnabled,
		&status.Tick.PlaydeckZoneID,
		&status.Tick.PlaydeckStreak,
		&status.Tick.PlaydeckTimeoutTicks,
		&status.Tick.ActiveGatheringTask,
		&status.Tick.GatheringRemainder,
		&status.Tick.LastTickAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get dev player status: %w", err)
	}

	activities, err := s.GetPlayerActivitySkills(ctx, status.Player.ID)
	if err != nil {
		return nil, err
	}

	status.Activities = activities
	activeActivity := status.Activities.ByKey(status.Tick.ActiveGatheringTask)

	status.Player.CreditsCents = status.Player.CurrencyCents
	status.Player.XPToNext = XPToNextLevel(status.Player.Level)
	status.Companion.XPToNext = NamiXPToNextLevel(status.Companion.Level)
	status.Companion.MoodScore = NamiMoodScore(status.Companion)
	status.Companion.MoodLabel = NamiMoodLabel(status.Companion.MoodScore)
	status.Companion.PrimaryNeed = NamiPrimaryNeed(status.Companion)
	status.Companion.Caption = NamiCaption(status.Companion)
	status.Companion.SuggestedAction = NamiSuggestedAction(status.Companion)
	status.Tick.PlaydeckZoneName = ZoneName(status.Tick.PlaydeckZoneID)
	status.Tick.ActiveGatheringName = GatheringTaskName(status.Tick.ActiveGatheringTask)
	status.Tick.ActiveGatheringOutput = GatheringResourceName(status.Tick.ActiveGatheringTask)
	status.Tick.ActiveActivityLevel = activeActivity.Level
	status.Tick.ActiveActivityXPIntoLevel = activeActivity.XPIntoLevel
	status.Tick.ActiveActivityXPToNext = activeActivity.XPToNext
	status.Tick.ResourcePerTick = ResourcePerTick(activeActivity.Level, status.Companion.MoodScore)
	status.Tick.ResourcePerTickDisplay = int64(math.Round(status.Tick.ResourcePerTick))
	status.Tick.NextTickAt = status.Tick.LastTickAt.Add(TickSeconds * time.Second)
	status.Tick.SecondsUntilNextTick = SecondsUntil(status.Tick.NextTickAt)

	return &status, nil
}

func (s *Store) SetDevGatheringTask(ctx context.Context, task string) error {
	task = strings.TrimSpace(strings.ToLower(task))
	if !ValidGatheringTask(task) {
		return fmt.Errorf("invalid gathering task: %s", task)
	}

	if _, err := s.SettleDevTicks(ctx, 0); err != nil {
		return fmt.Errorf("settle ticks before gathering switch: %w", err)
	}

	commandTag, err := s.Pool.Exec(ctx, `
		update player_tick_state
		set active_gathering_task = $1,
			updated_at = now()
		where player_id = (
			select id
			from players
			where display_name = 'Soryn'
		)
	`, task)

	if err != nil {
		return fmt.Errorf("set dev gathering task: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("dev player not found")
	}

	return nil
}

func (s *Store) SettleDevTicks(ctx context.Context, forcedTicks int64) (*TickResult, error) {
	if forcedTicks < 0 {
		forcedTicks = 0
	}

	if forcedTicks > 100 {
		forcedTicks = 100
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tick settlement: %w", err)
	}

	defer tx.Rollback(ctx)

	var playerID int64
	var level int
	var totalXP int64
	var xpIntoLevel int64
	var moodScore float64
	var playdeckEnabled bool
	var playdeckStreak int64
	var playdeckTimeoutTicks int
	var activeGatheringTask string
	var gatheringRemainder float64
	var lastTickAt time.Time

	if err := tx.QueryRow(ctx, `
		select
			p.id,
			p.level,
			p.total_xp,
			p.xp_into_level,
			c.mood_score::float8,
			t.playdeck_enabled,
			t.playdeck_streak,
			t.playdeck_timeout_ticks,
			t.active_gathering_task,
			t.gathering_remainder,
			t.last_tick_at
		from players p
		join companion_states c on c.player_id = p.id
		join player_tick_state t on t.player_id = p.id
		where p.display_name = 'Soryn'
		for update
	`).Scan(
		&playerID,
		&level,
		&totalXP,
		&xpIntoLevel,
		&moodScore,
		&playdeckEnabled,
		&playdeckStreak,
		&playdeckTimeoutTicks,
		&activeGatheringTask,
		&gatheringRemainder,
		&lastTickAt,
	); err != nil {
		return nil, fmt.Errorf("load tick state: %w", err)
	}

	now := time.Now().UTC()
	ticksToProcess := forcedTicks
	newLastTickAt := lastTickAt

	if ticksToProcess == 0 {
		elapsedTicks := int64(now.Sub(lastTickAt) / (TickSeconds * time.Second))
		if elapsedTicks <= 0 {
			return &TickResult{
				OK:           true,
				CurrentLevel: level,
				XPIntoLevel:  xpIntoLevel,
				XPToNext:     XPToNextLevel(level),
				Message:      "No ticks ready yet.",
			}, nil
		}

		ticksToProcess = elapsedTicks
		if ticksToProcess > MaxOfflineTicks {
			ticksToProcess = MaxOfflineTicks
			newLastTickAt = now
		} else {
			newLastTickAt = lastTickAt.Add(time.Duration(ticksToProcess*TickSeconds) * time.Second)
		}
	} else {
		newLastTickAt = now
	}

	resourceColumn, resourceName := GatheringResourceColumn(activeGatheringTask)

	if _, err := tx.Exec(ctx, `
	insert into player_activity_skills (player_id, activity_key)
	values ($1, $2)
	on conflict (player_id, activity_key) do nothing
`, playerID, activeGatheringTask); err != nil {
		return nil, fmt.Errorf("ensure active activity skill: %w", err)
	}

	var activityLevel int
	var activityTotalXP int64
	var activityXPIntoLevel int64

	if err := tx.QueryRow(ctx, `
	select level,
		total_xp,
		xp_into_level
	from player_activity_skills
	where player_id = $1
		and activity_key = $2
	for update
`, playerID, activeGatheringTask).Scan(
		&activityLevel,
		&activityTotalXP,
		&activityXPIntoLevel,
	); err != nil {
		return nil, fmt.Errorf("load active activity skill: %w", err)
	}

	result := &TickResult{
		OK:                   true,
		ResourceName:         resourceName,
		ActivityName:         GatheringTaskName(activeGatheringTask),
		ActivityCurrentLevel: activityLevel,
	}

	for i := int64(0); i < ticksToProcess; i++ {
		result.TicksProcessed++

		if playdeckEnabled {
			if playdeckTimeoutTicks > 0 {
				playdeckTimeoutTicks--
			} else {
				playdeckStreak++

				result.SyncXPGained += SyncXPPerTick
				totalXP += SyncXPPerTick
				xpIntoLevel += SyncXPPerTick

				for xpIntoLevel >= XPToNextLevel(level) {
					xpIntoLevel -= XPToNextLevel(level)
					level++
					result.LevelUps++
				}

				result.CreditsCentsGained += CreditsCentsPerPlaydeckWin(level)
				result.NibblesGained += NibblesPerPlaydeckWin(level)
			}
		}

		resourceRaw := ResourcePerTick(activityLevel, moodScore) + gatheringRemainder
		resourceWhole := int64(math.Floor(resourceRaw + 0.000000001))
		gatheringRemainder = resourceRaw - float64(resourceWhole)
		result.ResourceAmountGained += resourceWhole
	}

	if result.TicksProcessed > 0 {
		result.ActivityXPGained = result.TicksProcessed * ActivityXPPerTick
		activityTotalXP += result.ActivityXPGained
		activityXPIntoLevel += result.ActivityXPGained

		for activityXPIntoLevel >= ActivityXPToNextLevel(activityLevel) {
			activityXPIntoLevel -= ActivityXPToNextLevel(activityLevel)
			activityLevel++
			result.ActivityLevelUps++
		}
	}

	if _, err := tx.Exec(ctx, `
	update players
		set level = $1,
			total_xp = $2,
			xp_into_level = $3,
			currency_cents = currency_cents + $4,
			nibbles = nibbles + $5,
			updated_at = now()
		where id = $6
	`, level, totalXP, xpIntoLevel, result.CreditsCentsGained, result.NibblesGained, playerID); err != nil {
		return nil, fmt.Errorf("update player after ticks: %w", err)
	}

	if result.ResourceAmountGained > 0 {
		resourceSQL := fmt.Sprintf(`
			update player_resources
			set %s = %s + $1,
				updated_at = now()
			where player_id = $2
		`, resourceColumn, resourceColumn)

		if _, err := tx.Exec(ctx, resourceSQL, result.ResourceAmountGained, playerID); err != nil {
			return nil, fmt.Errorf("update gathering resource: %w", err)
		}
	}

	if result.TicksProcessed > 0 {
		if _, err := tx.Exec(ctx, `
		update player_activity_skills
		set level = $1,
			total_xp = $2,
			xp_into_level = $3,
			updated_at = now()
		where player_id = $4
			and activity_key = $5
	`, activityLevel, activityTotalXP, activityXPIntoLevel, playerID, activeGatheringTask); err != nil {
			return nil, fmt.Errorf("update activity skill after ticks: %w", err)
		}
	}

	if _, err := tx.Exec(ctx, `
		update player_tick_state
		set playdeck_streak = $1,
			playdeck_timeout_ticks = $2,
			gathering_remainder = $3,
			last_tick_at = $4,
			updated_at = now()
		where player_id = $5
	`, playdeckStreak, playdeckTimeoutTicks, gatheringRemainder, newLastTickAt, playerID); err != nil {
		return nil, fmt.Errorf("update tick state: %w", err)
	}

	if result.TicksProcessed > 0 {
		if _, err := tx.Exec(ctx, `
			insert into activity_log (player_id, event_type, message)
			values ($1, 'tick_settlement', $2)
		`, playerID, fmt.Sprintf(
			"Processed %d tick(s): +%d Sync XP, +%d %s.",
			result.TicksProcessed,
			result.SyncXPGained,
			result.ResourceAmountGained,
			resourceName,
		)); err != nil {
			return nil, fmt.Errorf("insert tick activity log: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tick settlement: %w", err)
	}

	result.CurrentLevel = level
	result.XPIntoLevel = xpIntoLevel
	result.XPToNext = XPToNextLevel(level)
	result.ActivityCurrentLevel = activityLevel
	result.ActivityXPIntoLevel = activityXPIntoLevel
	result.ActivityXPToNext = ActivityXPToNextLevel(activityLevel)
	result.Message = fmt.Sprintf("Processed %d tick(s).", result.TicksProcessed)

	return result, nil
}

func (s *Store) GetPlayerActivitySkills(ctx context.Context, playerID int64) (ActivitySkills, error) {
	skills := DefaultActivitySkills()

	rows, err := s.Pool.Query(ctx, `
		select activity_key,
			level,
			total_xp,
			xp_into_level
		from player_activity_skills
		where player_id = $1
	`, playerID)
	if err != nil {
		return skills, fmt.Errorf("query player activity skills: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var skill ActivitySkill

		if err := rows.Scan(
			&skill.Key,
			&skill.Level,
			&skill.TotalXP,
			&skill.XPIntoLevel,
		); err != nil {
			return skills, fmt.Errorf("scan player activity skill: %w", err)
		}

		skill.Name = GatheringTaskName(skill.Key)
		skill.XPToNext = ActivityXPToNextLevel(skill.Level)
		skills.Set(skill)
	}

	if err := rows.Err(); err != nil {
		return skills, fmt.Errorf("iterate player activity skills: %w", err)
	}

	return skills, nil
}

func DefaultActivitySkills() ActivitySkills {
	return ActivitySkills{
		Streaming:     NewActivitySkill("streaming"),
		DoomScrolling: NewActivitySkill("doom_scrolling"),
		Cleaning:      NewActivitySkill("cleaning"),
		Exercising:    NewActivitySkill("exercising"),
		Shopping:      NewActivitySkill("shopping"),
		Designing:     NewActivitySkill("designing"),
	}
}

func NewActivitySkill(key string) ActivitySkill {
	return ActivitySkill{
		Key:         key,
		Name:        GatheringTaskName(key),
		Level:       1,
		TotalXP:     0,
		XPIntoLevel: 0,
		XPToNext:    ActivityXPToNextLevel(1),
	}
}

func (skills *ActivitySkills) Set(skill ActivitySkill) {
	if skill.Level < 1 {
		skill.Level = 1
	}

	skill.Name = GatheringTaskName(skill.Key)
	skill.XPToNext = ActivityXPToNextLevel(skill.Level)

	switch skill.Key {
	case "streaming":
		skills.Streaming = skill
	case "doom_scrolling":
		skills.DoomScrolling = skill
	case "cleaning":
		skills.Cleaning = skill
	case "exercising":
		skills.Exercising = skill
	case "shopping":
		skills.Shopping = skill
	case "designing":
		skills.Designing = skill
	}
}

func (skills ActivitySkills) ByKey(key string) ActivitySkill {
	switch key {
	case "streaming":
		return skills.Streaming
	case "doom_scrolling":
		return skills.DoomScrolling
	case "cleaning":
		return skills.Cleaning
	case "exercising":
		return skills.Exercising
	case "shopping":
		return skills.Shopping
	case "designing":
		return skills.Designing
	default:
		return NewActivitySkill("streaming")
	}
}

func ActivityXPToNextLevel(level int) int64 {
	if level < 1 {
		level = 1
	}

	ticksToNext := math.Round(60 + 12*math.Sqrt(float64(level)))
	return int64(ticksToNext) * ActivityXPPerTick
}

func XPToNextLevel(level int) int64 {
	if level < 1 {
		level = 1
	}

	ticksToNext := math.Round(60 + 12*math.Sqrt(float64(level)))
	return int64(ticksToNext) * SyncXPPerTick
}

func ResourcePerTick(level int, moodScore float64) float64 {
	if level < 1 {
		level = 1
	}

	if moodScore < 0 {
		moodScore = 0
	}

	if moodScore > 100 {
		moodScore = 100
	}

	base := math.Pow(float64(level), 1.1) + 100
	moodMultiplier := (moodScore / 200) + 1

	return base * moodMultiplier
}

func CreditsCentsPerPlaydeckWin(level int) int64 {
	if level < 1 {
		level = 1
	}

	credits := math.Pow(float64(level), 1.05) + 25
	return int64(math.Round(credits * 100))
}

func NibblesPerPlaydeckWin(level int) int64 {
	if level < 1 {
		level = 1
	}

	nibbles := math.Round(math.Pow(float64(level), 0.45))
	if nibbles < 1 {
		nibbles = 1
	}

	return int64(nibbles)
}

func SecondsUntil(next time.Time) int {
	seconds := int(math.Ceil(time.Until(next).Seconds()))

	if seconds < 0 {
		return 0
	}

	if seconds > TickSeconds {
		return TickSeconds
	}

	return seconds
}

func GatheringTaskKeys() []string {
	return []string{
		"streaming",
		"doom_scrolling",
		"cleaning",
		"exercising",
		"shopping",
		"designing",
	}
}

func ValidGatheringTask(task string) bool {
	switch task {
	case "streaming", "doom_scrolling", "cleaning", "exercising", "shopping", "designing":
		return true
	default:
		return false
	}
}

func GatheringTaskName(task string) string {
	switch task {
	case "streaming":
		return "Streaming"
	case "doom_scrolling":
		return "Scrolling"
	case "cleaning":
		return "Cleaning"
	case "exercising":
		return "Exercise"
	case "shopping":
		return "Shopping"
	case "designing":
		return "Designing"
	default:
		return "Unknown"
	}
}

func GatheringResourceName(task string) string {
	switch task {
	case "streaming":
		return "Fans"
	case "doom_scrolling":
		return "Memes"
	case "cleaning":
		return "Lost Items"
	case "exercising":
		return "Confidence"
	case "shopping":
		return "Receipts"
	case "designing":
		return "Patterns"
	default:
		return "Resources"
	}
}

func GatheringResourceColumn(task string) (string, string) {
	switch task {
	case "streaming":
		return "fans", "Fans"
	case "doom_scrolling":
		return "memes", "Memes"
	case "cleaning":
		return "lost_items", "Lost Items"
	case "exercising":
		return "confidence", "Confidence"
	case "shopping":
		return "receipts", "Receipts"
	case "designing":
		return "patterns", "Patterns"
	default:
		return "fans", "Fans"
	}
}

func NamiXPToNextLevel(level int) int64 {
	if level < 1 {
		level = 1
	}

	nextLevel := level + 1
	raw := 100 + nextLevel*(10+level)

	return int64(math.Floor(float64(raw)/10) * 10)
}

func NamiMoodScore(companion CompanionState) float64 {
	score :=
		float64(companion.Comfort)*0.25 +
			float64(companion.Connection)*0.20 +
			float64(companion.Energy)*0.15 +
			float64(companion.Satiety)*0.15 +
			float64(companion.Cleanliness)*0.10 +
			float64(companion.Playfulness)*0.10 +
			float64(companion.Inspiration)*0.05

	if score < 0 {
		return 0
	}

	if score > 100 {
		return 100
	}

	return score
}

func NamiMoodLabel(score float64) string {
	switch {
	case score >= 90:
		return "Radiant"
	case score >= 75:
		return "Cozy"
	case score >= 60:
		return "Okay"
	case score >= 40:
		return "Pouty"
	case score >= 20:
		return "Wilted"
	default:
		return "Emergency Blanket Burrito"
	}
}

func NamiPrimaryNeed(companion CompanionState) string {
	if companion.Status == "sleeping" {
		return "Sleeping"
	}

	switch {
	case companion.Cleanliness < 30:
		return "Needs a Bath"
	case companion.Satiety < 30:
		return "Needs Food"
	case companion.Energy < 25:
		return "Needs Sleep"
	case companion.Connection < 30:
		return "Needs Attention"
	case companion.Comfort < 30:
		return "Needs Comfort"
	case companion.Playfulness < 30:
		return "Bored"
	case companion.Inspiration < 30:
		return "Needs Inspiration"
	default:
		return NamiMoodLabel(companion.MoodScore)
	}
}

func NamiCaption(companion CompanionState) string {
	switch companion.PrimaryNeed {
	case "Sleeping":
		return "Nami-chan is asleep and recovering energy."
	case "Needs a Bath":
		return "Nami-chan is dirty and needs a bath."
	case "Needs Food":
		return "Nami-chan needs food, a snack, or a drink."
	case "Needs Sleep":
		return "Nami-chan is sleepy and needs rest."
	case "Needs Attention":
		return "Nami-chan misses you and needs attention."
	case "Needs Comfort":
		return "Nami-chan needs cozy care."
	case "Bored":
		return "Nami-chan is bored and needs engagement."
	case "Needs Inspiration":
		return "Nami-chan needs creativity, reading, or writing time."
	case "Radiant":
		return "Nami-chan is radiant, spoiled, and very pleased."
	case "Cozy":
		return "Nami-chan is happy, cozy, and content."
	case "Okay":
		return "Nami-chan is okay and waiting sweetly."
	case "Pouty":
		return "Nami-chan is pouty and a bit low."
	case "Wilted":
		return "Nami-chan is wilted and needs care."
	default:
		return "Nami-chan has retreated into emergency blanket mode."
	}
}

func NamiSuggestedAction(companion CompanionState) string {
	switch companion.PrimaryNeed {
	case "Sleeping":
		return "Let her sleep or wake her up."
	case "Needs a Bath":
		return "Bath or freshen up."
	case "Needs Food":
		return "Meal, snack, or drink."
	case "Needs Sleep":
		return "Nap or put her to bed."
	case "Needs Attention":
		return "Cuddle, read together, or boop."
	case "Needs Comfort":
		return "Cuddle, read together, or freshen up."
	case "Bored":
		return "Play or boop."
	case "Needs Inspiration":
		return "Write together or read together."
	default:
		return "Any care action would make her happy."
	}
}

func ZoneName(zoneID int) string {
	switch zoneID {
	case 1:
		return "Starter Deck"
	case 2:
		return "Cozy LAN Café"
	case 3:
		return "Neon Mall Net"
	default:
		return "Unknown Zone"
	}
}
