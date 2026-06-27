package database

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	TickSeconds                   = 5
	OnlineTickWriteMinimumSeconds = 30
	OnlineTickMaxAwardSeconds     = 45
	SyncXPPerTick                 = int64(10)
	ActivityXPPerTick             = int64(10)
	MaxOfflineTicks               = int64(8640)
	NamiMessageStorageLimit       = 50

	CareDecayMinimumSeconds      = 5 * 60
	SleepEnergyRecoveryPerHour   = 10.0
	SleepRecoveryCapHours        = 8.0
	SleepingSatietyDecayPerHour  = 1.0
	AwakeSatietyDecayPerHour     = 3.0
	AwakeConnectionDecayPerHour  = 2.0
	AwakeEnergyDecayPerHour      = 4.0
	AwakeComfortDecayPerHour     = 1.0
	AwakePlayfulnessDecayPerHour = 1.5
	AwakeInspirationDecayPerHour = 1.5
	AwakeCleanlinessDecayPerHour = 1.25
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
	Care       CareQueueState  `json:"care"`
	Playdeck   PlaydeckStatus  `json:"playdeck"`
	Wardrobe   WardrobeStatus  `json:"wardrobe"`
}

type Player struct {
	ID            int64     `json:"id"`
	DisplayName   string    `json:"displayName"`
	CreatedAt     time.Time `json:"createdAt"`
	OnlineSeconds int64     `json:"onlineSeconds"`
	Level         int       `json:"level"`
	TotalXP       int64     `json:"totalXp"`
	XPIntoLevel   int64     `json:"xpIntoLevel"`
	XPToNext      int64     `json:"xpToNext"`
	CurrencyCents int64     `json:"currencyCents"`
	CreditsCents  int64     `json:"creditsCents"`
	Nibbles       int64     `json:"nibbles"`
	NamiCoin      int64     `json:"namiCoin"`
}

type CompanionState struct {
	CompanionName                string    `json:"name"`
	MoodScore                    float64   `json:"moodScore"`
	Satiety                      int       `json:"satiety"`
	Connection                   int       `json:"connection"`
	Energy                       int       `json:"energy"`
	Comfort                      int       `json:"comfort"`
	Playfulness                  int       `json:"playfulness"`
	Inspiration                  int       `json:"inspiration"`
	Cleanliness                  int       `json:"cleanliness"`
	Status                       string    `json:"status"`
	LastInteractionAt            time.Time `json:"lastInteractionAt"`
	Level                        int       `json:"level"`
	TotalXP                      int64     `json:"totalXp"`
	XPIntoLevel                  int64     `json:"xpIntoLevel"`
	XPToNext                     int64     `json:"xpToNext"`
	LastXPGained                 int64     `json:"lastXpGained"`
	LastAction                   string    `json:"lastAction"`
	SleepStartedAt               time.Time `json:"sleepStartedAt"`
	EnergyAtSleepStart           int       `json:"energyAtSleepStart"`
	LastDecayAt                  time.Time `json:"lastDecayAt"`
	SatietyDecayRemainder        float64   `json:"-"`
	ConnectionDecayRemainder     float64   `json:"-"`
	EnergyDecayRemainder         float64   `json:"-"`
	ComfortDecayRemainder        float64   `json:"-"`
	PlayfulnessDecayRemainder    float64   `json:"-"`
	InspirationDecayRemainder    float64   `json:"-"`
	CleanlinessDecayRemainder    float64   `json:"-"`
	SleepEnergyRecoveryRemainder float64   `json:"-"`
	MoodLabel                    string    `json:"moodLabel"`
	PrimaryNeed                  string    `json:"primaryNeed"`
	Caption                      string    `json:"caption"`
	SuggestedAction              string    `json:"suggestedAction"`
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
	PlaydeckMaxStreak         int64     `json:"playdeckMaxStreak"`
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

type CareQueueState struct {
	Active CareActionState   `json:"active"`
	Queued []CareActionState `json:"queued"`
	Slots  int               `json:"slots"`
}

type CareActionState struct {
	ID               int64     `json:"id"`
	Action           string    `json:"action"`
	ActionName       string    `json:"actionName"`
	Status           string    `json:"status"`
	QueuePosition    int       `json:"queuePosition"`
	DurationSeconds  int       `json:"durationSeconds"`
	StartedAt        time.Time `json:"startedAt"`
	CompletesAt      time.Time `json:"completesAt"`
	CompletedAt      time.Time `json:"completedAt"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	SecondsRemaining int       `json:"secondsRemaining"`
	ProgressPercent  float64   `json:"progressPercent"`
}

type CareActionResult struct {
	OK           bool           `json:"ok"`
	Action       string         `json:"action"`
	ActionName   string         `json:"actionName"`
	Mode         string         `json:"mode"`
	XPGained     int64          `json:"xpGained"`
	LevelUps     int            `json:"levelUps"`
	CurrentLevel int            `json:"currentLevel"`
	XPIntoLevel  int64          `json:"xpIntoLevel"`
	XPToNext     int64          `json:"xpToNext"`
	Companion    CompanionState `json:"companion"`
	Care         CareQueueState `json:"care"`
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

type NamiMessage struct {
	ID           int64     `json:"id"`
	PlayerID     int64     `json:"playerId"`
	TriggerKey   string    `json:"triggerKey"`
	MoodKey      string    `json:"moodKey"`
	NeedKey      string    `json:"needKey"`
	Severity     string    `json:"severity"`
	Message      string    `json:"message"`
	MetadataJSON string    `json:"metadataJson"`
	CreatedAt    time.Time `json:"createdAt"`
	SeenAt       time.Time `json:"seenAt"`
}

type NamiMessageDraft struct {
	TriggerKey   string
	MoodKey      string
	NeedKey      string
	Severity     string
	Message      string
	MetadataJSON string
}

type NamiProceduralContext struct {
	TriggerKey   string
	ActionKey    string
	ActionName   string
	MoodKey      string
	NeedKey      string
	Severity     string
	ResourceName string
	ActivityName string
	Level        int
	LevelUps     int
	MetadataJSON string
}

type NamiCareStat struct {
	Key   string
	Name  string
	Value int
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

func CareActionDurationSeconds(action string) int {
	switch strings.TrimSpace(strings.ToLower(action)) {
	case "boop":
		return 30
	case "drink":
		return 2 * 60
	case "snack":
		return 5 * 60
	case "freshen_up":
		return 10 * 60
	case "cuddle":
		return 15 * 60
	case "meal":
		return 30 * 60
	case "play":
		return 25 * 60
	case "write_together":
		return 30 * 60
	case "read_together":
		return 30 * 60
	case "bath":
		return 30 * 60
	case "nap":
		return 60 * 60
	case "put_to_bed":
		return 60 * 60
	case "wake_up":
		return 5 * 60
	default:
		return 0
	}
}

func CareQueueSlots() int {
	return 3
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

func (s *Store) DevPlayerID(ctx context.Context) (int64, error) {
	if accountID, ok := AuthAccountIDFromContext(ctx); ok {
		playerID, err := s.PlayerIDForAccount(ctx, accountID)
		if err != nil {
			return 0, fmt.Errorf("get account player id: %w", err)
		}

		return playerID, nil
	}

	var playerID int64

	if err := s.Pool.QueryRow(ctx, `
select id
from players
where display_name = 'Soryn'
`).Scan(&playerID); err != nil {
		return 0, fmt.Errorf("get dev player id: %w", err)
	}

	return playerID, nil
}

func (s *Store) RewindDevCareDecay(ctx context.Context, duration time.Duration) error {
	if duration <= 0 {
		return fmt.Errorf("rewind duration must be positive")
	}

	if duration > 30*24*time.Hour {
		duration = 30 * 24 * time.Hour
	}

	playerID, err := s.DevPlayerID(ctx)
	if err != nil {
		return err
	}

	_, err = s.Pool.Exec(ctx, `
		update companion_states
		set last_decay_at = now() - ($2::double precision * interval '1 second'),
			updated_at = now()
		where player_id = $1
	`, playerID, duration.Seconds())
	if err != nil {
		return fmt.Errorf("rewind dev care decay: %w", err)
	}

	return nil
}

func (s *Store) CreateDevNamiMessage(ctx context.Context, draft NamiMessageDraft) (*NamiMessage, error) {
	playerID, err := s.DevPlayerID(ctx)
	if err != nil {
		return nil, err
	}

	return s.CreateNamiMessage(ctx, playerID, draft)
}

func (s *Store) CreateNamiMessage(ctx context.Context, playerID int64, draft NamiMessageDraft) (*NamiMessage, error) {
	draft = normalizeNamiMessageDraft(draft)

	var message NamiMessage

	if err := s.Pool.QueryRow(ctx, `
		insert into nami_messages (
			player_id,
			trigger_key,
			mood_key,
			need_key,
			severity,
			message,
			metadata_json
		)
		values ($1, $2, $3, $4, $5, $6, $7::jsonb)
		returning
			id,
			player_id,
			trigger_key,
			mood_key,
			need_key,
			severity,
			message,
			metadata_json::text,
			created_at,
			coalesce(seen_at, '0001-01-01 00:00:00+00'::timestamptz)
	`,
		playerID,
		draft.TriggerKey,
		draft.MoodKey,
		draft.NeedKey,
		draft.Severity,
		draft.Message,
		draft.MetadataJSON,
	).Scan(
		&message.ID,
		&message.PlayerID,
		&message.TriggerKey,
		&message.MoodKey,
		&message.NeedKey,
		&message.Severity,
		&message.Message,
		&message.MetadataJSON,
		&message.CreatedAt,
		&message.SeenAt,
	); err != nil {
		return nil, fmt.Errorf("create nami message: %w", err)
	}

	if err := s.PruneNamiMessages(ctx, playerID); err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *Store) GetRecentDevNamiMessages(ctx context.Context, limit int) ([]NamiMessage, error) {
	playerID, err := s.DevPlayerID(ctx)
	if err != nil {
		return nil, err
	}

	return s.GetRecentNamiMessages(ctx, playerID, limit)
}

func (s *Store) GetRecentNamiMessages(ctx context.Context, playerID int64, limit int) ([]NamiMessage, error) {
	if limit < 1 {
		limit = 1
	}

	if limit > NamiMessageStorageLimit {
		limit = NamiMessageStorageLimit
	}

	rows, err := s.Pool.Query(ctx, `
		select
			id,
			player_id,
			trigger_key,
			mood_key,
			need_key,
			severity,
			message,
			metadata_json::text,
			created_at,
			coalesce(seen_at, '0001-01-01 00:00:00+00'::timestamptz)
		from nami_messages
		where player_id = $1
		order by created_at desc, id desc
		limit $2
	`, playerID, limit)
	if err != nil {
		return nil, fmt.Errorf("get recent nami messages: %w", err)
	}
	defer rows.Close()

	var messages []NamiMessage
	for rows.Next() {
		var message NamiMessage

		if err := rows.Scan(
			&message.ID,
			&message.PlayerID,
			&message.TriggerKey,
			&message.MoodKey,
			&message.NeedKey,
			&message.Severity,
			&message.Message,
			&message.MetadataJSON,
			&message.CreatedAt,
			&message.SeenAt,
		); err != nil {
			return nil, fmt.Errorf("scan nami message: %w", err)
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate nami messages: %w", err)
	}

	for left, right := 0, len(messages)-1; left < right; left, right = left+1, right-1 {
		messages[left], messages[right] = messages[right], messages[left]
	}

	return messages, nil
}

func pruneNamiMessagesTx(ctx context.Context, tx pgx.Tx, playerID int64) error {
	_, err := tx.Exec(ctx, `
with ranked_messages as (
select
id,
row_number() over (
partition by player_id
order by created_at desc, id desc
) as message_rank
from nami_messages
where player_id = $1
)
delete from nami_messages message
using ranked_messages ranked
where message.id = ranked.id
and ranked.message_rank > $2
`, playerID, NamiMessageStorageLimit)
	if err != nil {
		return fmt.Errorf("prune nami messages: %w", err)
	}

	return nil
}

func (s *Store) PruneNamiMessages(ctx context.Context, playerID int64) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin prune nami messages: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := pruneNamiMessagesTx(ctx, tx, playerID); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit prune nami messages: %w", err)
	}

	return nil
}
func loadRecentNamiMessagesTx(ctx context.Context, tx pgx.Tx, playerID int64, limit int) ([]NamiMessage, error) {
	if limit < 1 {
		limit = 1
	}

	if limit > NamiMessageStorageLimit {
		limit = NamiMessageStorageLimit
	}

	rows, err := tx.Query(ctx, `
		select
			id,
			player_id,
			trigger_key,
			mood_key,
			need_key,
			severity,
			message,
			metadata_json::text,
			created_at,
			coalesce(seen_at, '0001-01-01 00:00:00+00'::timestamptz)
		from nami_messages
		where player_id = $1
		order by created_at desc, id desc
		limit $2
	`, playerID, limit)
	if err != nil {
		return nil, fmt.Errorf("load recent nami messages in tx: %w", err)
	}
	defer rows.Close()

	var messages []NamiMessage
	for rows.Next() {
		var message NamiMessage

		if err := rows.Scan(
			&message.ID,
			&message.PlayerID,
			&message.TriggerKey,
			&message.MoodKey,
			&message.NeedKey,
			&message.Severity,
			&message.Message,
			&message.MetadataJSON,
			&message.CreatedAt,
			&message.SeenAt,
		); err != nil {
			return nil, fmt.Errorf("scan recent nami message in tx: %w", err)
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recent nami messages in tx: %w", err)
	}

	return messages, nil
}

func insertNamiMessageDraftTx(ctx context.Context, tx pgx.Tx, playerID int64, draft NamiMessageDraft) (*NamiMessage, error) {
	draft = normalizeNamiMessageDraft(draft)

	var message NamiMessage

	if err := tx.QueryRow(ctx, `
		insert into nami_messages (
			player_id,
			trigger_key,
			mood_key,
			need_key,
			severity,
			message,
			metadata_json
		)
		values ($1, $2, $3, $4, $5, $6, $7::jsonb)
		returning
			id,
			player_id,
			trigger_key,
			mood_key,
			need_key,
			severity,
			message,
			metadata_json::text,
			created_at,
			coalesce(seen_at, '0001-01-01 00:00:00+00'::timestamptz)
	`,
		playerID,
		draft.TriggerKey,
		draft.MoodKey,
		draft.NeedKey,
		draft.Severity,
		draft.Message,
		draft.MetadataJSON,
	).Scan(
		&message.ID,
		&message.PlayerID,
		&message.TriggerKey,
		&message.MoodKey,
		&message.NeedKey,
		&message.Severity,
		&message.Message,
		&message.MetadataJSON,
		&message.CreatedAt,
		&message.SeenAt,
	); err != nil {
		return nil, fmt.Errorf("insert nami message draft in tx: %w", err)
	}

	return &message, nil
}

func prependRecentNamiMessage(recent []NamiMessage, message *NamiMessage) []NamiMessage {
	if message == nil {
		return recent
	}

	return append([]NamiMessage{*message}, recent...)
}

func (s *Store) GetDevCareQueueState(ctx context.Context) (CareQueueState, error) {
	playerID, err := s.DevPlayerID(ctx)
	if err != nil {
		return CareQueueState{Slots: CareQueueSlots()}, err
	}

	return s.GetCareQueueState(ctx, playerID)
}

func (s *Store) GetCareQueueState(ctx context.Context, playerID int64) (CareQueueState, error) {
	rows, err := s.Pool.Query(ctx, `
		select
			id,
			action_key,
			action_name,
			status,
			coalesce(queue_position, 0),
			duration_seconds,
			coalesce(started_at, '0001-01-01 00:00:00+00'::timestamptz),
			coalesce(completes_at, '0001-01-01 00:00:00+00'::timestamptz),
			coalesce(completed_at, '0001-01-01 00:00:00+00'::timestamptz),
			created_at,
			updated_at
		from companion_care_actions
		where player_id = $1
			and status in ('active', 'queued')
		order by
			case when status = 'active' then 0 else 1 end,
			queue_position nulls last,
			created_at,
			id
	`, playerID)
	if err != nil {
		return CareQueueState{Slots: CareQueueSlots()}, fmt.Errorf("get care queue state: %w", err)
	}
	defer rows.Close()

	state := CareQueueState{
		Slots: CareQueueSlots(),
	}

	for rows.Next() {
		var action CareActionState

		if err := rows.Scan(
			&action.ID,
			&action.Action,
			&action.ActionName,
			&action.Status,
			&action.QueuePosition,
			&action.DurationSeconds,
			&action.StartedAt,
			&action.CompletesAt,
			&action.CompletedAt,
			&action.CreatedAt,
			&action.UpdatedAt,
		); err != nil {
			return CareQueueState{Slots: CareQueueSlots()}, fmt.Errorf("scan care queue state: %w", err)
		}

		action = hydrateCareActionDisplay(action)

		if action.Status == "active" {
			state.Active = action
		} else {
			state.Queued = append(state.Queued, action)
		}
	}

	if err := rows.Err(); err != nil {
		return CareQueueState{Slots: CareQueueSlots()}, fmt.Errorf("iterate care queue state: %w", err)
	}

	return state, nil
}

func loadCareQueueStateTx(ctx context.Context, tx pgx.Tx, playerID int64) (CareQueueState, error) {
	rows, err := tx.Query(ctx, `
		select
			id,
			action_key,
			action_name,
			status,
			coalesce(queue_position, 0),
			duration_seconds,
			coalesce(started_at, '0001-01-01 00:00:00+00'::timestamptz),
			coalesce(completes_at, '0001-01-01 00:00:00+00'::timestamptz),
			coalesce(completed_at, '0001-01-01 00:00:00+00'::timestamptz),
			created_at,
			updated_at
		from companion_care_actions
		where player_id = $1
			and status in ('active', 'queued')
		order by
			case when status = 'active' then 0 else 1 end,
			queue_position nulls last,
			created_at,
			id
		for update
	`, playerID)
	if err != nil {
		return CareQueueState{Slots: CareQueueSlots()}, fmt.Errorf("load care queue state in tx: %w", err)
	}
	defer rows.Close()

	state := CareQueueState{
		Slots: CareQueueSlots(),
	}

	for rows.Next() {
		var action CareActionState

		if err := rows.Scan(
			&action.ID,
			&action.Action,
			&action.ActionName,
			&action.Status,
			&action.QueuePosition,
			&action.DurationSeconds,
			&action.StartedAt,
			&action.CompletesAt,
			&action.CompletedAt,
			&action.CreatedAt,
			&action.UpdatedAt,
		); err != nil {
			return CareQueueState{Slots: CareQueueSlots()}, fmt.Errorf("scan care queue state in tx: %w", err)
		}

		action = hydrateCareActionDisplay(action)

		if action.Status == "active" {
			state.Active = action
		} else {
			state.Queued = append(state.Queued, action)
		}
	}

	if err := rows.Err(); err != nil {
		return CareQueueState{Slots: CareQueueSlots()}, fmt.Errorf("iterate care queue state in tx: %w", err)
	}

	return hydrateCareQueueState(state), nil
}

func (s *Store) StartOrQueueDevCareAction(ctx context.Context, action string) (*CareActionResult, error) {
	rule, ok := CareActionByKey(action)
	if !ok {
		return nil, fmt.Errorf("invalid care action: %s", action)
	}

	durationSeconds := CareActionDurationSeconds(rule.Key)
	if durationSeconds <= 0 {
		return nil, fmt.Errorf("invalid care action duration: %s", rule.Key)
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin start or queue care action: %w", err)
	}
	defer tx.Rollback(ctx)

	playerID, companion, err := loadDevCompanionForUpdateTx(ctx, tx)
	if err != nil {
		return nil, err
	}

	if err := settleCompletedCareActionsTx(ctx, tx, playerID); err != nil {
		return nil, err
	}

	_, companion, err = loadDevCompanionForUpdateTx(ctx, tx)
	if err != nil {
		return nil, err
	}

	companion, err = settleCareDecayTx(ctx, tx, playerID, companion)
	if err != nil {
		return nil, err
	}

	state, err := loadCareQueueStateTx(ctx, tx, playerID)
	if err != nil {
		return nil, err
	}

	isSleeping := strings.ToLower(companion.Status) == "sleeping"

	if isSleeping && rule.Key != "wake_up" {
		return nil, fmt.Errorf("Nami is sleeping; wake her before starting other care actions")
	}

	if !isSleeping && rule.Key == "wake_up" {
		return nil, fmt.Errorf("Nami is already awake")
	}

	if !careActionIsZero(state.Active) && state.Active.Action == "put_to_bed" {
		return nil, fmt.Errorf("Nami is settling into sleep; no actions can be queued")
	}

	if queuedAction, ok := findQueuedCareAction(state.Queued, rule.Key); ok {
		if err := removeQueuedCareActionTx(ctx, tx, queuedAction.ID, playerID); err != nil {
			return nil, err
		}

		if err := renumberQueuedCareActionsTx(ctx, tx, playerID); err != nil {
			return nil, err
		}

		state, err = loadCareQueueStateTx(ctx, tx, playerID)
		if err != nil {
			return nil, err
		}

		if err := tx.Commit(ctx); err != nil {
			return nil, fmt.Errorf("commit unqueue care action: %w", err)
		}

		return &CareActionResult{
			OK:         true,
			Action:     rule.Key,
			ActionName: rule.Name,
			Mode:       "unqueued",
			Care:       state,
			Message:    fmt.Sprintf("%s removed from queue.", rule.Name),
		}, nil
	}

	if careQueueHasSleepBarrier(state.Queued) {
		return nil, fmt.Errorf("sleep is already queued; no actions can be queued after sleep")
	}

	if !careActionIsZero(state.Active) {
		if len(state.Queued) >= CareQueueSlots() {
			return nil, fmt.Errorf("care queue is full")
		}

		if careActionShouldBlockQueue(state.Active.Action) {
			return nil, fmt.Errorf("active sleep blocks care queue")
		}

		position := nextCareQueuePosition(state.Queued)
		if position == 0 {
			return nil, fmt.Errorf("care queue is full")
		}

		_, err = tx.Exec(ctx, `
			insert into companion_care_actions (
				player_id,
				action_key,
				action_name,
				status,
				queue_position,
				duration_seconds
			)
			values ($1, $2, $3, 'queued', $4, $5)
		`, playerID, rule.Key, rule.Name, position, durationSeconds)
		if err != nil {
			return nil, fmt.Errorf("queue care action: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return nil, fmt.Errorf("commit queue care action: %w", err)
		}

		state, err = s.GetCareQueueState(ctx, playerID)
		if err != nil {
			return nil, err
		}

		return &CareActionResult{
			OK:         true,
			Action:     rule.Key,
			ActionName: rule.Name,
			Mode:       "queued",
			Care:       state,
			Message:    fmt.Sprintf("%s queued.", rule.Name),
		}, nil
	}

	activeAction, err := startCareActionTx(ctx, tx, playerID, rule)
	if err != nil {
		return nil, err
	}

	resultCompanion := companion
	xpGained := int64(0)

	if rule.SleepAction {
		if err := setCompanionSleepingTx(ctx, tx, playerID); err != nil {
			return nil, err
		}

		resultCompanion.Status = "sleeping"
	} else if careActionRewardsOnStart(rule) {
		if err := applyCompletedCareActionTx(ctx, tx, playerID, activeAction); err != nil {
			return nil, fmt.Errorf("apply started care action rewards: %w", err)
		}

		_, refreshedCompanion, err := loadDevCompanionForUpdateTx(ctx, tx)
		if err != nil {
			return nil, err
		}

		resultCompanion = refreshedCompanion
		xpGained = refreshedCompanion.LastXPGained
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit start care action: %w", err)
	}

	state, err = s.GetCareQueueState(ctx, playerID)
	if err != nil {
		return nil, err
	}

	return &CareActionResult{
		OK:           true,
		Action:       rule.Key,
		ActionName:   rule.Name,
		Mode:         "started",
		XPGained:     xpGained,
		CurrentLevel: resultCompanion.Level,
		XPIntoLevel:  resultCompanion.XPIntoLevel,
		XPToNext:     resultCompanion.XPToNext,
		Companion:    resultCompanion,
		Care:         state,
		Message:      fmt.Sprintf("%s started.", activeAction.ActionName),
	}, nil
}

func loadDevCompanionForUpdateTx(ctx context.Context, tx pgx.Tx) (int64, CompanionState, error) {
	playerID, err := playerIDForContextTx(ctx, tx)
	if err != nil {
		return 0, CompanionState{}, err
	}

	var companion CompanionState

	err = tx.QueryRow(ctx, `
        select
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
            c.last_action,
            coalesce(c.sleep_started_at, '0001-01-01 00:00:00+00'::timestamptz),
            coalesce(c.energy_at_sleep_start, 0),
            c.last_decay_at,
            coalesce(c.satiety_decay_remainder, 0)::float8,
            coalesce(c.connection_decay_remainder, 0)::float8,
            coalesce(c.energy_decay_remainder, 0)::float8,
            coalesce(c.comfort_decay_remainder, 0)::float8,
            coalesce(c.playfulness_decay_remainder, 0)::float8,
            coalesce(c.inspiration_decay_remainder, 0)::float8,
            coalesce(c.cleanliness_decay_remainder, 0)::float8,
            coalesce(c.sleep_energy_recovery_remainder, 0)::float8
        from companion_states c
        where c.player_id = $1
        for update
    `, playerID).Scan(
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
		&companion.SleepStartedAt,
		&companion.EnergyAtSleepStart,
		&companion.LastDecayAt,
		&companion.SatietyDecayRemainder,
		&companion.ConnectionDecayRemainder,
		&companion.EnergyDecayRemainder,
		&companion.ComfortDecayRemainder,
		&companion.PlayfulnessDecayRemainder,
		&companion.InspirationDecayRemainder,
		&companion.CleanlinessDecayRemainder,
		&companion.SleepEnergyRecoveryRemainder,
	)

	if err != nil {
		return 0, CompanionState{}, fmt.Errorf("load companion for update: %w", err)
	}

	companion.MoodScore = NamiMoodScore(companion)
	companion.XPToNext = NamiXPToNextLevel(companion.Level)
	companion.MoodLabel = NamiMoodLabel(companion.MoodScore)
	companion.PrimaryNeed = NamiPrimaryNeed(companion)
	companion.Caption = NamiCaption(companion)
	companion.SuggestedAction = NamiSuggestedAction(companion)

	return playerID, companion, nil
}
func findQueuedCareAction(queued []CareActionState, actionKey string) (CareActionState, bool) {
	for _, action := range queued {
		if action.Action == actionKey {
			return action, true
		}
	}

	return CareActionState{}, false
}

func removeQueuedCareActionTx(ctx context.Context, tx pgx.Tx, actionID int64, playerID int64) error {
	commandTag, err := tx.Exec(ctx, `
		update companion_care_actions
		set status = 'cancelled',
			queue_position = null,
			updated_at = now()
		where id = $1
			and player_id = $2
			and status = 'queued'
	`, actionID, playerID)
	if err != nil {
		return fmt.Errorf("remove queued care action: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("queued care action not found")
	}

	return nil
}

func renumberQueuedCareActionsTx(ctx context.Context, tx pgx.Tx, playerID int64) error {
	rows, err := tx.Query(ctx, `
		select id
		from companion_care_actions
		where player_id = $1
			and status = 'queued'
		order by queue_position, created_at, id
		for update
	`, playerID)
	if err != nil {
		return fmt.Errorf("load queued care actions for renumber: %w", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("scan queued care action for renumber: %w", err)
		}

		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate queued care actions for renumber: %w", err)
	}

	for index, id := range ids {
		_, err := tx.Exec(ctx, `
			update companion_care_actions
			set queue_position = $1,
				updated_at = now()
			where id = $2
				and player_id = $3
				and status = 'queued'
		`, index+1, id, playerID)
		if err != nil {
			return fmt.Errorf("renumber queued care action: %w", err)
		}
	}

	return nil
}

func startCareActionTx(ctx context.Context, tx pgx.Tx, playerID int64, rule CareActionRule) (CareActionState, error) {
	durationSeconds := CareActionDurationSeconds(rule.Key)
	if durationSeconds <= 0 {
		return CareActionState{}, fmt.Errorf("invalid care action duration: %s", rule.Key)
	}

	var action CareActionState

	err := tx.QueryRow(ctx, `
		insert into companion_care_actions (
			player_id,
			action_key,
			action_name,
			status,
			queue_position,
			duration_seconds,
			started_at,
			completes_at
		)
		values ($1, $2, $3, 'active', null, $4, now(), now() + ($4::int * interval '1 second'))
		returning
			id,
			action_key,
			action_name,
			status,
			coalesce(queue_position, 0),
			duration_seconds,
			coalesce(started_at, '0001-01-01 00:00:00+00'::timestamptz),
			coalesce(completes_at, '0001-01-01 00:00:00+00'::timestamptz),
			coalesce(completed_at, '0001-01-01 00:00:00+00'::timestamptz),
			created_at,
			updated_at
	`, playerID, rule.Key, rule.Name, durationSeconds).Scan(
		&action.ID,
		&action.Action,
		&action.ActionName,
		&action.Status,
		&action.QueuePosition,
		&action.DurationSeconds,
		&action.StartedAt,
		&action.CompletesAt,
		&action.CompletedAt,
		&action.CreatedAt,
		&action.UpdatedAt,
	)

	if err != nil {
		return CareActionState{}, fmt.Errorf("start care action: %w", err)
	}

	return hydrateCareActionDisplay(action), nil
}

func setCompanionSleepingTx(ctx context.Context, tx pgx.Tx, playerID int64) error {
	_, err := tx.Exec(ctx, `
	update companion_states
	set status = 'sleeping',
		sleep_started_at = coalesce(sleep_started_at, now()),
		energy_at_sleep_start = coalesce(energy_at_sleep_start, energy),
		last_decay_at = now(),
		sleep_energy_recovery_remainder = 0,
		updated_at = now()
	where player_id = $1
`, playerID)
	if err != nil {
		return fmt.Errorf("set companion sleeping: %w", err)
	}

	return nil
}

func (s *Store) SettleDevCareActions(ctx context.Context) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin settle dev care actions: %w", err)
	}
	defer tx.Rollback(ctx)

	playerID, _, err := loadDevCompanionForUpdateTx(ctx, tx)
	if err != nil {
		return err
	}

	if err := settleCompletedCareActionsTx(ctx, tx, playerID); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit settle dev care actions: %w", err)
	}

	return nil
}

func settleCompletedCareActionsTx(ctx context.Context, tx pgx.Tx, playerID int64) error {
	for {
		active, ok, err := loadDueActiveCareActionTx(ctx, tx, playerID)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		activeRule, ok := CareActionByKey(active.Action)
		if !ok {
			return fmt.Errorf("invalid active care action: %s", active.Action)
		}

		if careActionRewardsOnCompletion(activeRule) {
			if err := applyCompletedCareActionTx(ctx, tx, playerID, active); err != nil {
				return err
			}
		}

		if err := completeActiveCareActionTx(ctx, tx, playerID, active.ID); err != nil {
			return err
		}

		nextQueued, ok, err := popNextQueuedCareActionTx(ctx, tx, playerID)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		rule, ok := CareActionByKey(nextQueued.Action)
		if !ok {
			return fmt.Errorf("invalid queued care action: %s", nextQueued.Action)
		}

		if err := activateQueuedCareActionTx(ctx, tx, playerID, nextQueued.ID); err != nil {
			return err
		}

		if rule.SleepAction {
			if err := setCompanionSleepingTx(ctx, tx, playerID); err != nil {
				return err
			}

			return nil
		}

		if careActionRewardsOnStart(rule) {
			if err := applyCompletedCareActionTx(ctx, tx, playerID, nextQueued); err != nil {
				return fmt.Errorf("apply queued care action rewards: %w", err)
			}
		}
	}
}

func loadDueActiveCareActionTx(ctx context.Context, tx pgx.Tx, playerID int64) (CareActionState, bool, error) {
	var action CareActionState

	err := tx.QueryRow(ctx, `
		select
			id,
			action_key,
			action_name,
			status,
			coalesce(queue_position, 0),
			duration_seconds,
			coalesce(started_at, '0001-01-01 00:00:00+00'::timestamptz),
			coalesce(completes_at, '0001-01-01 00:00:00+00'::timestamptz),
			coalesce(completed_at, '0001-01-01 00:00:00+00'::timestamptz),
			created_at,
			updated_at
		from companion_care_actions
		where player_id = $1
			and status = 'active'
			and completes_at <= now()
		order by completes_at, id
		limit 1
		for update
	`, playerID).Scan(
		&action.ID,
		&action.Action,
		&action.ActionName,
		&action.Status,
		&action.QueuePosition,
		&action.DurationSeconds,
		&action.StartedAt,
		&action.CompletesAt,
		&action.CompletedAt,
		&action.CreatedAt,
		&action.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return CareActionState{}, false, nil
		}

		return CareActionState{}, false, fmt.Errorf("load due active care action: %w", err)
	}

	return hydrateCareActionDisplay(action), true, nil
}

func completeActiveCareActionTx(ctx context.Context, tx pgx.Tx, playerID int64, actionID int64) error {
	commandTag, err := tx.Exec(ctx, `
		update companion_care_actions
		set status = 'completed',
			completed_at = now(),
			updated_at = now()
		where id = $1
			and player_id = $2
			and status = 'active'
	`, actionID, playerID)
	if err != nil {
		return fmt.Errorf("complete active care action: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("active care action not found")
	}

	return nil
}

func popNextQueuedCareActionTx(ctx context.Context, tx pgx.Tx, playerID int64) (CareActionState, bool, error) {
	if err := renumberQueuedCareActionsTx(ctx, tx, playerID); err != nil {
		return CareActionState{}, false, err
	}

	var action CareActionState

	err := tx.QueryRow(ctx, `
		select
			id,
			action_key,
			action_name,
			status,
			coalesce(queue_position, 0),
			duration_seconds,
			coalesce(started_at, '0001-01-01 00:00:00+00'::timestamptz),
			coalesce(completes_at, '0001-01-01 00:00:00+00'::timestamptz),
			coalesce(completed_at, '0001-01-01 00:00:00+00'::timestamptz),
			created_at,
			updated_at
		from companion_care_actions
		where player_id = $1
			and status = 'queued'
			and queue_position = 1
		limit 1
		for update
	`, playerID).Scan(
		&action.ID,
		&action.Action,
		&action.ActionName,
		&action.Status,
		&action.QueuePosition,
		&action.DurationSeconds,
		&action.StartedAt,
		&action.CompletesAt,
		&action.CompletedAt,
		&action.CreatedAt,
		&action.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return CareActionState{}, false, nil
		}

		return CareActionState{}, false, fmt.Errorf("pop next queued care action: %w", err)
	}

	return hydrateCareActionDisplay(action), true, nil
}

func activateQueuedCareActionTx(ctx context.Context, tx pgx.Tx, playerID int64, actionID int64) error {
	commandTag, err := tx.Exec(ctx, `
		update companion_care_actions
		set status = 'active',
			queue_position = null,
			started_at = now(),
			completes_at = now() + (duration_seconds::int * interval '1 second'),
			updated_at = now()
		where id = $1
			and player_id = $2
			and status = 'queued'
	`, actionID, playerID)
	if err != nil {
		return fmt.Errorf("activate queued care action: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("queued care action not found")
	}

	if err := renumberQueuedCareActionsTx(ctx, tx, playerID); err != nil {
		return err
	}

	return nil
}

func applyCompletedCareActionTx(ctx context.Context, tx pgx.Tx, playerID int64, active CareActionState) error {
	rule, ok := CareActionByKey(active.Action)
	if !ok {
		return fmt.Errorf("invalid completed care action: %s", active.Action)
	}

	_, beforeCompanion, err := loadDevCompanionForUpdateTx(ctx, tx)
	if err != nil {
		return err
	}

	companion := beforeCompanion
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
	}

	companion.Satiety = clampCareStat(companion.Satiety + rule.Satiety)
	companion.Connection = clampCareStat(companion.Connection + rule.Connection)
	companion.Energy = clampCareStat(companion.Energy + rule.Energy)
	companion.Comfort = clampCareStat(companion.Comfort + rule.Comfort)
	companion.Playfulness = clampCareStat(companion.Playfulness + rule.Playfulness)
	companion.Inspiration = clampCareStat(companion.Inspiration + rule.Inspiration)
	companion.Cleanliness = clampCareStat(companion.Cleanliness + rule.Cleanliness)

	companion.TotalXP += xpGained
	companion.XPIntoLevel += xpGained

	levelUps := 0
	for companion.XPIntoLevel >= NamiXPToNextLevel(companion.Level) {
		companion.XPIntoLevel -= NamiXPToNextLevel(companion.Level)
		companion.Level++
		levelUps++
	}

	companion.MoodScore = NamiMoodScore(companion)
	companion.XPToNext = NamiXPToNextLevel(companion.Level)
	companion.MoodLabel = NamiMoodLabel(companion.MoodScore)
	companion.PrimaryNeed = NamiPrimaryNeed(companion)
	companion.Caption = NamiCaption(companion)
	companion.SuggestedAction = NamiSuggestedAction(companion)
	companion.LastXPGained = xpGained
	companion.LastAction = rule.Name

	sleepStartedAtUpdate := "sleep_started_at"
	energyAtSleepStartUpdate := "energy_at_sleep_start"
	sleepEnergyRecoveryRemainderUpdate := "sleep_energy_recovery_remainder"

	if rule.WakeAction {
		sleepStartedAtUpdate = "null"
		energyAtSleepStartUpdate = "null"
		sleepEnergyRecoveryRemainderUpdate = "0"
	}

	_, err = tx.Exec(ctx, fmt.Sprintf(`
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
			sleep_started_at = %s,
			energy_at_sleep_start = %s,
			sleep_energy_recovery_remainder = %s,
			updated_at = now()
		where player_id = $15
	`, sleepStartedAtUpdate, energyAtSleepStartUpdate, sleepEnergyRecoveryRemainderUpdate),
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
		return fmt.Errorf("update companion after completed care action: %w", err)
	}

	_, err = tx.Exec(ctx, `
		insert into activity_log (player_id, event_type, message)
		values ($1, 'care_action', $2)
	`, playerID, fmt.Sprintf("Completed care action: %s (+%d XP).", rule.Name, xpGained))
	if err != nil {
		return fmt.Errorf("insert completed care action log: %w", err)
	}

	recentMessages, err := loadRecentNamiMessagesTx(ctx, tx, playerID, 80)
	if err != nil {
		return err
	}

	careMessageDraft := GenerateNamiCareMessageDraft(rule, beforeCompanion, companion, levelUps, recentMessages)
	careMessage, err := insertNamiMessageDraftTx(ctx, tx, playerID, careMessageDraft)
	if err != nil {
		return fmt.Errorf("insert completed care nami message: %w", err)
	}

	recentMessages = prependRecentNamiMessage(recentMessages, careMessage)

	if levelUps > 0 {
		levelUpDraft := GenerateNamiEventMessageDraft(NamiProceduralContext{
			TriggerKey:   "nami_level_up",
			MoodKey:      companion.MoodLabel,
			NeedKey:      companion.PrimaryNeed,
			Severity:     "happy",
			Level:        companion.Level,
			LevelUps:     levelUps,
			MetadataJSON: fmt.Sprintf(`{"level":%d,"levelUps":%d,"sourceAction":"%s","sourceActionName":"%s"}`, companion.Level, levelUps, rule.Key, rule.Name),
		}, recentMessages)

		if _, err := insertNamiMessageDraftTx(ctx, tx, playerID, levelUpDraft); err != nil {
			return fmt.Errorf("insert completed care level-up message: %w", err)
		}
	}

	return nil
}

func hydrateCareActionDisplay(action CareActionState) CareActionState {
	if action.Status != "active" || action.CompletesAt.Year() <= 1 {
		action.SecondsRemaining = action.DurationSeconds
		action.ProgressPercent = 0
		return action
	}

	remaining := int(math.Floor(time.Until(action.CompletesAt).Seconds()))
	if remaining < 0 {
		remaining = 0
	}

	if action.DurationSeconds > 0 && remaining > action.DurationSeconds {
		remaining = action.DurationSeconds
	}

	action.SecondsRemaining = remaining

	if action.DurationSeconds <= 0 {
		action.ProgressPercent = 100
		return action
	}

	elapsed := action.DurationSeconds - remaining
	if elapsed < 0 {
		elapsed = 0
	}

	progress := (float64(elapsed) / float64(action.DurationSeconds)) * 100
	if progress < 0 {
		progress = 0
	}

	if progress > 100 {
		progress = 100
	}

	action.ProgressPercent = progress
	return action
}

func hydrateCareQueueState(state CareQueueState) CareQueueState {
	state.Active = hydrateCareActionDisplay(state.Active)

	for i := range state.Queued {
		state.Queued[i] = hydrateCareActionDisplay(state.Queued[i])
	}

	if state.Slots == 0 {
		state.Slots = CareQueueSlots()
	}

	return state
}

func careActionFromRow(
	id int64,
	actionKey string,
	actionName string,
	status string,
	queuePosition int,
	durationSeconds int,
	startedAt time.Time,
	completesAt time.Time,
	completedAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) CareActionState {
	return hydrateCareActionDisplay(CareActionState{
		ID:              id,
		Action:          actionKey,
		ActionName:      actionName,
		Status:          status,
		QueuePosition:   queuePosition,
		DurationSeconds: durationSeconds,
		StartedAt:       startedAt,
		CompletesAt:     completesAt,
		CompletedAt:     completedAt,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	})
}

func nextCareQueuePosition(queued []CareActionState) int {
	used := make(map[int]bool, len(queued))

	for _, action := range queued {
		if action.QueuePosition > 0 {
			used[action.QueuePosition] = true
		}
	}

	for position := 1; position <= CareQueueSlots(); position++ {
		if !used[position] {
			return position
		}
	}

	return 0
}

func careQueueHasSleepBarrier(queued []CareActionState) bool {
	for _, action := range queued {
		if action.Action == "put_to_bed" {
			return true
		}
	}

	return false
}

func careActionIsZero(action CareActionState) bool {
	return action.ID == 0
}

func careActionShouldBlockQueue(actionKey string) bool {
	return actionKey == "put_to_bed"
}

func careActionRewardsOnStart(rule CareActionRule) bool {
	return !rule.SleepAction && !rule.WakeAction
}

func careActionRewardsOnCompletion(rule CareActionRule) bool {
	return rule.WakeAction
}

func (s *Store) SettleDevCareDecay(ctx context.Context) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin settle dev care decay: %w", err)
	}
	defer tx.Rollback(ctx)

	playerID, companion, err := loadDevCompanionForUpdateTx(ctx, tx)
	if err != nil {
		return err
	}

	if _, err := settleCareDecayTx(ctx, tx, playerID, companion); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit settle dev care decay: %w", err)
	}

	return nil
}

func settleCareDecayTx(ctx context.Context, tx pgx.Tx, playerID int64, companion CompanionState) (CompanionState, error) {
	now := time.Now().UTC()
	lastDecayAt := companion.LastDecayAt

	if lastDecayAt.Year() <= 1 {
		if companion.LastInteractionAt.Year() > 1 {
			lastDecayAt = companion.LastInteractionAt
		} else {
			lastDecayAt = now
		}
	}

	if now.Before(lastDecayAt) {
		lastDecayAt = now
	}

	elapsed := now.Sub(lastDecayAt)
	if elapsed < time.Duration(CareDecayMinimumSeconds)*time.Second {
		return companion, nil
	}

	hours := elapsed.Hours()
	decayed := companion

	if strings.ToLower(companion.Status) == "sleeping" {
		decayed = applySleepingCareDecay(companion, hours, lastDecayAt, now)
	} else {
		decayed = applyAwakeCareDecay(companion, hours)
	}

	decayed.MoodScore = NamiMoodScore(decayed)
	decayed.MoodLabel = NamiMoodLabel(decayed.MoodScore)
	decayed.PrimaryNeed = NamiPrimaryNeed(decayed)
	decayed.Caption = NamiCaption(decayed)
	decayed.SuggestedAction = NamiSuggestedAction(decayed)
	decayed.LastDecayAt = now

	_, err := tx.Exec(ctx, `
		update companion_states
		set
			mood_score = $1,
			satiety = $2,
			connection = $3,
			energy = $4,
			comfort = $5,
			playfulness = $6,
			inspiration = $7,
			cleanliness = $8,
			last_decay_at = $9,
			satiety_decay_remainder = $10,
			connection_decay_remainder = $11,
			energy_decay_remainder = $12,
			comfort_decay_remainder = $13,
			playfulness_decay_remainder = $14,
			inspiration_decay_remainder = $15,
			cleanliness_decay_remainder = $16,
			sleep_energy_recovery_remainder = $17,
			updated_at = now()
		where player_id = $18
	`,
		decayed.MoodScore,
		decayed.Satiety,
		decayed.Connection,
		decayed.Energy,
		decayed.Comfort,
		decayed.Playfulness,
		decayed.Inspiration,
		decayed.Cleanliness,
		decayed.LastDecayAt,
		decayed.SatietyDecayRemainder,
		decayed.ConnectionDecayRemainder,
		decayed.EnergyDecayRemainder,
		decayed.ComfortDecayRemainder,
		decayed.PlayfulnessDecayRemainder,
		decayed.InspirationDecayRemainder,
		decayed.CleanlinessDecayRemainder,
		decayed.SleepEnergyRecoveryRemainder,
		playerID,
	)
	if err != nil {
		return CompanionState{}, fmt.Errorf("update care decay: %w", err)
	}

	return decayed, nil
}

func applyAwakeCareDecay(companion CompanionState, hours float64) CompanionState {
	decayed := companion

	satietyLoss := AwakeSatietyDecayPerHour * hours
	connectionLoss := AwakeConnectionDecayPerHour * hours
	energyLoss := AwakeEnergyDecayPerHour * hours
	comfortLoss := AwakeComfortDecayPerHour * hours
	playfulnessLoss := AwakePlayfulnessDecayPerHour * hours
	inspirationLoss := AwakeInspirationDecayPerHour * hours
	cleanlinessLoss := AwakeCleanlinessDecayPerHour * hours

	if companion.Satiety < 40 {
		comfortLoss += 1.0 * hours
	}

	if companion.Satiety < 30 {
		comfortLoss += 1.0 * hours
	}

	if companion.Connection < 40 {
		comfortLoss += 1.5 * hours
	}

	if companion.Connection < 30 {
		comfortLoss += 1.0 * hours
	}

	if companion.Cleanliness < 50 {
		comfortLoss += 1.0 * hours
	}

	if companion.Cleanliness < 30 {
		comfortLoss += 1.0 * hours
	}

	if companion.Energy < 30 {
		playfulnessLoss += 1.0 * hours
	}

	if companion.Comfort < 40 {
		inspirationLoss += 1.0 * hours
	}

	if companion.Cleanliness < 40 {
		inspirationLoss += 0.5 * hours
	}

	if companion.Satiety < 25 {
		energyLoss += 0.75 * hours
	}

	if companion.Cleanliness < 25 {
		connectionLoss += 0.5 * hours
	}

	decayed.Satiety, decayed.SatietyDecayRemainder = applyCareStatLoss(
		companion.Satiety,
		satietyLoss,
		companion.SatietyDecayRemainder,
	)

	decayed.Connection, decayed.ConnectionDecayRemainder = applyCareStatLoss(
		companion.Connection,
		connectionLoss,
		companion.ConnectionDecayRemainder,
	)

	decayed.Energy, decayed.EnergyDecayRemainder = applyCareStatLoss(
		companion.Energy,
		energyLoss,
		companion.EnergyDecayRemainder,
	)

	decayed.Comfort, decayed.ComfortDecayRemainder = applyCareStatLoss(
		companion.Comfort,
		comfortLoss,
		companion.ComfortDecayRemainder,
	)

	decayed.Playfulness, decayed.PlayfulnessDecayRemainder = applyCareStatLoss(
		companion.Playfulness,
		playfulnessLoss,
		companion.PlayfulnessDecayRemainder,
	)

	decayed.Inspiration, decayed.InspirationDecayRemainder = applyCareStatLoss(
		companion.Inspiration,
		inspirationLoss,
		companion.InspirationDecayRemainder,
	)

	decayed.Cleanliness, decayed.CleanlinessDecayRemainder = applyCareStatLoss(
		companion.Cleanliness,
		cleanlinessLoss,
		companion.CleanlinessDecayRemainder,
	)

	return decayed
}

func applySleepingCareDecay(companion CompanionState, hours float64, lastDecayAt time.Time, now time.Time) CompanionState {
	decayed := companion

	decayed.Satiety, decayed.SatietyDecayRemainder = applyCareStatLoss(
		companion.Satiety,
		SleepingSatietyDecayPerHour*hours,
		companion.SatietyDecayRemainder,
	)

	sleepStartedAt := companion.SleepStartedAt
	if sleepStartedAt.Year() <= 1 {
		sleepStartedAt = lastDecayAt
	}

	alreadyRecoveredHours := lastDecayAt.Sub(sleepStartedAt).Hours()
	if alreadyRecoveredHours < 0 {
		alreadyRecoveredHours = 0
	}

	totalRecoverableHoursRemaining := SleepRecoveryCapHours - alreadyRecoveredHours
	if totalRecoverableHoursRemaining < 0 {
		totalRecoverableHoursRemaining = 0
	}

	recoveryHours := math.Min(hours, totalRecoverableHoursRemaining)
	if recoveryHours < 0 {
		recoveryHours = 0
	}

	energyGain := SleepEnergyRecoveryPerHour * recoveryHours
	decayed.Energy, decayed.SleepEnergyRecoveryRemainder = applyCareStatGain(
		companion.Energy,
		energyGain,
		companion.SleepEnergyRecoveryRemainder,
	)

	_ = now

	return decayed
}

func applyCareStatLoss(value int, loss float64, remainder float64) (int, float64) {
	if value <= 0 {
		return 0, 0
	}

	totalLoss := loss + remainder
	if totalLoss < 0 {
		totalLoss = 0
	}

	wholeLoss := int(math.Floor(totalLoss))
	nextRemainder := totalLoss - float64(wholeLoss)

	if wholeLoss <= 0 {
		return value, nextRemainder
	}

	nextValue := clampCareStat(value - wholeLoss)
	if nextValue <= 0 {
		nextRemainder = 0
	}

	return nextValue, nextRemainder
}

func applyCareStatGain(value int, gain float64, remainder float64) (int, float64) {
	if value >= 100 {
		return 100, 0
	}

	totalGain := gain + remainder
	if totalGain < 0 {
		totalGain = 0
	}

	wholeGain := int(math.Floor(totalGain))
	nextRemainder := totalGain - float64(wholeGain)

	if wholeGain <= 0 {
		return value, nextRemainder
	}

	nextValue := clampCareStat(value + wholeGain)
	if nextValue >= 100 {
		nextRemainder = 0
	}

	return nextValue, nextRemainder
}

func (s *Store) GenerateDevPassiveNamiMessages(ctx context.Context) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin passive nami messages: %w", err)
	}
	defer tx.Rollback(ctx)

	playerID, companion, err := loadDevCompanionForUpdateTx(ctx, tx)
	if err != nil {
		return fmt.Errorf("load companion for passive nami messages: %w", err)
	}

	companion.MoodScore = NamiMoodScore(companion)
	companion.MoodLabel = NamiMoodLabel(companion.MoodScore)
	companion.PrimaryNeed = NamiPrimaryNeed(companion)

	_, err = tx.Exec(ctx, `
		insert into player_nami_message_state (player_id)
		values ($1)
		on conflict (player_id) do nothing
	`, playerID)
	if err != nil {
		return fmt.Errorf("ensure nami message state: %w", err)
	}

	var lastOnlineMessageAt time.Time
	var nextRandomMessageAt time.Time

	err = tx.QueryRow(ctx, `
		select
			coalesce(last_online_message_at, '0001-01-01 00:00:00+00'::timestamptz),
			coalesce(next_random_message_at, '0001-01-01 00:00:00+00'::timestamptz)
		from player_nami_message_state
		where player_id = $1
		for update
	`, playerID).Scan(&lastOnlineMessageAt, &nextRandomMessageAt)
	if err != nil {
		return fmt.Errorf("load nami message state: %w", err)
	}

	recentMessages, err := loadRecentNamiMessagesTx(ctx, tx, playerID, 100)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	if lastOnlineMessageAt.Year() <= 1 || now.Sub(lastOnlineMessageAt) >= 4*time.Hour {
		draft := GenerateNamiEventMessageDraft(NamiProceduralContext{
			TriggerKey:   "user_online",
			MoodKey:      companion.MoodLabel,
			NeedKey:      companion.PrimaryNeed,
			Severity:     namiSeverityForCompanion(companion),
			Level:        companion.Level,
			MetadataJSON: fmt.Sprintf(`{"moodScore":%.2f,"primaryNeed":"%s"}`, companion.MoodScore, companion.PrimaryNeed),
		}, recentMessages)

		message, err := insertNamiMessageDraftTx(ctx, tx, playerID, draft)
		if err != nil {
			return err
		}

		recentMessages = prependRecentNamiMessage(recentMessages, message)

		if _, err := tx.Exec(ctx, `
			update player_nami_message_state
			set last_online_message_at = $1,
				updated_at = now()
			where player_id = $2
		`, now, playerID); err != nil {
			return fmt.Errorf("update last online nami message timestamp: %w", err)
		}
	}

	for _, stat := range NamiCareStats(companion) {
		if stat.Value >= 20 {
			continue
		}

		triggerKey := "care_stat_low_" + stat.Key

		var recentlySent bool
		if err := tx.QueryRow(ctx, `
			select exists (
				select 1
				from nami_messages
				where player_id = $1
					and trigger_key = $2
					and created_at > now() - interval '1 hour'
			)
		`, playerID, triggerKey).Scan(&recentlySent); err != nil {
			return fmt.Errorf("check recent low stat nami message: %w", err)
		}

		if recentlySent {
			continue
		}

		severity := "low"
		if stat.Value < 10 {
			severity = "urgent"
		}

		draft := GenerateNamiEventMessageDraft(NamiProceduralContext{
			TriggerKey:   triggerKey,
			MoodKey:      companion.MoodLabel,
			NeedKey:      companion.PrimaryNeed,
			Severity:     severity,
			ResourceName: stat.Name,
			Level:        companion.Level,
			MetadataJSON: fmt.Sprintf(`{"stat":"%s","statName":"%s","value":%d}`, stat.Key, stat.Name, stat.Value),
		}, recentMessages)

		message, err := insertNamiMessageDraftTx(ctx, tx, playerID, draft)
		if err != nil {
			return err
		}

		recentMessages = prependRecentNamiMessage(recentMessages, message)
	}

	if nextRandomMessageAt.Year() <= 1 {
		nextRandom := nextRandomNamiMessageAt(now, playerID, len(recentMessages))
		if _, err := tx.Exec(ctx, `
			update player_nami_message_state
			set next_random_message_at = $1,
				updated_at = now()
			where player_id = $2
		`, nextRandom, playerID); err != nil {
			return fmt.Errorf("initialize next random nami message timestamp: %w", err)
		}
	} else if !now.Before(nextRandomMessageAt) {
		draft := GenerateNamiEventMessageDraft(NamiProceduralContext{
			TriggerKey:   "random_mood",
			MoodKey:      companion.MoodLabel,
			NeedKey:      companion.PrimaryNeed,
			Severity:     namiSeverityForCompanion(companion),
			Level:        companion.Level,
			MetadataJSON: fmt.Sprintf(`{"moodScore":%.2f,"primaryNeed":"%s"}`, companion.MoodScore, companion.PrimaryNeed),
		}, recentMessages)

		message, err := insertNamiMessageDraftTx(ctx, tx, playerID, draft)
		if err != nil {
			return err
		}

		recentMessages = prependRecentNamiMessage(recentMessages, message)

		nextRandom := nextRandomNamiMessageAt(now, playerID, len(recentMessages))
		if _, err := tx.Exec(ctx, `
			update player_nami_message_state
			set next_random_message_at = $1,
				updated_at = now()
			where player_id = $2
		`, nextRandom, playerID); err != nil {
			return fmt.Errorf("update next random nami message timestamp: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit passive nami messages: %w", err)
	}

	return nil
}

func NamiCareStats(companion CompanionState) []NamiCareStat {
	return []NamiCareStat{
		{Key: "satiety", Name: "Satiety", Value: companion.Satiety},
		{Key: "connection", Name: "Connection", Value: companion.Connection},
		{Key: "energy", Name: "Energy", Value: companion.Energy},
		{Key: "comfort", Name: "Comfort", Value: companion.Comfort},
		{Key: "playfulness", Name: "Playfulness", Value: companion.Playfulness},
		{Key: "inspiration", Name: "Inspiration", Value: companion.Inspiration},
		{Key: "cleanliness", Name: "Cleanliness", Value: companion.Cleanliness},
	}
}

func namiSeverityForCompanion(companion CompanionState) string {
	switch {
	case companion.MoodScore < 20:
		return "urgent"
	case companion.MoodScore < 40:
		return "low"
	case companion.MoodScore >= 75:
		return "happy"
	default:
		return "info"
	}
}

func nextRandomNamiMessageAt(now time.Time, playerID int64, salt int) time.Time {
	seed := fmt.Sprintf("random-nami|%d|%d|%d", playerID, salt, now.UnixNano())
	minutes := 60 + int(hashNamiMessageSeed(seed)%61)

	return now.Add(time.Duration(minutes) * time.Minute)
}

func normalizeNamiMessageDraft(draft NamiMessageDraft) NamiMessageDraft {
	draft.TriggerKey = strings.TrimSpace(strings.ToLower(draft.TriggerKey))
	draft.MoodKey = strings.TrimSpace(strings.ToLower(draft.MoodKey))
	draft.NeedKey = strings.TrimSpace(strings.ToLower(draft.NeedKey))
	draft.Severity = strings.TrimSpace(strings.ToLower(draft.Severity))
	draft.Message = strings.TrimSpace(draft.Message)
	draft.MetadataJSON = strings.TrimSpace(draft.MetadataJSON)

	if draft.TriggerKey == "" {
		draft.TriggerKey = "unknown"
	}

	if draft.Severity == "" {
		draft.Severity = "info"
	}

	if draft.Message == "" {
		draft.Message = "Nami-chan makes a tiny thoughtful noise."
	}

	if draft.MetadataJSON == "" {
		draft.MetadataJSON = "{}"
	}

	return draft
}

func GenerateNamiCareMessageDraft(rule CareActionRule, before CompanionState, after CompanionState, levelUps int, recent []NamiMessage) NamiMessageDraft {
	before.MoodScore = NamiMoodScore(before)
	after.MoodScore = NamiMoodScore(after)

	before.MoodLabel = NamiMoodLabel(before.MoodScore)
	after.MoodLabel = NamiMoodLabel(after.MoodScore)

	before.PrimaryNeed = NamiPrimaryNeed(before)
	after.PrimaryNeed = NamiPrimaryNeed(after)

	severity := "info"
	if after.MoodScore < 20 {
		severity = "urgent"
	} else if after.MoodScore < 40 {
		severity = "low"
	} else if after.MoodScore >= 75 {
		severity = "happy"
	}

	context := NamiProceduralContext{
		TriggerKey:   "care_" + rule.Key,
		ActionKey:    rule.Key,
		ActionName:   rule.Name,
		MoodKey:      normalizeNamiMessageKey(after.MoodLabel),
		NeedKey:      normalizeNamiMessageKey(after.PrimaryNeed),
		Severity:     severity,
		Level:        after.Level,
		LevelUps:     levelUps,
		MetadataJSON: fmt.Sprintf(`{"action":"%s","actionName":"%s","levelUps":%d}`, rule.Key, rule.Name, levelUps),
	}

	message := BuildProceduralNamiMessage(context, recent)

	return NamiMessageDraft{
		TriggerKey:   context.TriggerKey,
		MoodKey:      context.MoodKey,
		NeedKey:      context.NeedKey,
		Severity:     context.Severity,
		Message:      message,
		MetadataJSON: context.MetadataJSON,
	}
}

func GenerateNamiEventMessageDraft(context NamiProceduralContext, recent []NamiMessage) NamiMessageDraft {
	context.TriggerKey = strings.TrimSpace(strings.ToLower(context.TriggerKey))
	context.MoodKey = normalizeNamiMessageKey(context.MoodKey)
	context.NeedKey = normalizeNamiMessageKey(context.NeedKey)
	context.Severity = strings.TrimSpace(strings.ToLower(context.Severity))

	if context.TriggerKey == "" {
		context.TriggerKey = "unknown"
	}

	if context.Severity == "" {
		context.Severity = "info"
	}

	if context.MetadataJSON == "" {
		context.MetadataJSON = "{}"
	}

	return NamiMessageDraft{
		TriggerKey:   context.TriggerKey,
		MoodKey:      context.MoodKey,
		NeedKey:      context.NeedKey,
		Severity:     context.Severity,
		Message:      BuildProceduralNamiMessage(context, recent),
		MetadataJSON: context.MetadataJSON,
	}
}

func BuildProceduralNamiMessage(context NamiProceduralContext, recent []NamiMessage) string {
	recentText := make(map[string]bool, len(recent))
	for _, message := range recent {
		recentText[normalizeNamiFragment(message.Message)] = true
	}

	actionPool := namiActionMessagePool(context)
	if namiContextIsCareAction(context) {
		actionPool = appendNamiMessageParts(
			actionPool,
			namiUniversalCareReactionPool(context),
			namiGeneratedActionReactionPool(context),
		)
	}

	moodPool := appendNamiMessageParts(
		namiMoodMessagePool(context.MoodKey),
		namiUniversalMoodFlavorPool(context),
		namiGeneratedMoodImagePool(context),
	)

	needPool := appendNamiMessageParts(
		namiNeedMessagePool(context.NeedKey),
		namiUniversalNeedFlavorPool(context),
		namiGeneratedNeedFlavorPool(context),
	)

	openingPool := namiOpeningPoolForContext(context)

	closerPool := appendNamiMessageParts(
		namiCloserMessagePool(context),
		namiUniversalCloserPool(context),
		namiGeneratedCloserPool(context),
	)

	if len(actionPool) == 0 {
		actionPool = []string{"I noticed that. I am placing it carefully in my little internal scrapbook."}
	}

	recentFragments := buildRecentNamiFragmentSet(recent, openingPool, actionPool, moodPool, needPool, closerPool)

	baseSeed := fmt.Sprintf(
		"%s|%s|%s|%s|%s|%d|%d|%d",
		context.TriggerKey,
		context.ActionKey,
		context.MoodKey,
		context.NeedKey,
		context.Severity,
		context.Level,
		context.LevelUps,
		time.Now().UnixNano(),
	)

	moodChance := namiMoodFlavorChance(context)
	needChance := namiNeedFlavorChance(context)
	closerChance := namiCloserChance(context)

	for attempt := 0; attempt < 240; attempt++ {
		seed := fmt.Sprintf("%s|attempt:%d", baseSeed, attempt)

		pieces := []string{
			pickNamiMessagePartAvoiding(openingPool, seed+"|opening", recentFragments),
			pickNamiMessagePartAvoiding(actionPool, seed+"|action", recentFragments),
		}

		if shouldUseNamiPart(seed+"|mood", moodChance) {
			pieces = append(pieces, pickNamiMessagePartAvoiding(moodPool, seed+"|mood", recentFragments))
		}

		if shouldUseNamiPart(seed+"|need", needChance) {
			pieces = append(pieces, pickNamiMessagePartAvoiding(needPool, seed+"|need", recentFragments))
		}

		if shouldUseNamiPart(seed+"|closer", closerChance) {
			pieces = append(pieces, pickNamiMessagePartAvoiding(closerPool, seed+"|closer", recentFragments))
		}

		message := addNamiMessageSuffix(
			joinNamiMessagePieces(pieces),
			namiMessageSuffix(context),
		)

		if message != "" && !recentText[normalizeNamiFragment(message)] {
			return message
		}
	}

	return addNamiMessageSuffix(
		joinNamiMessagePieces([]string{
			pickNamiMessagePartAvoiding(openingPool, baseSeed+"|fallback-opening", recentFragments),
			pickNamiMessagePartAvoiding(actionPool, baseSeed+"|fallback-action", recentFragments),
			pickNamiMessagePartAvoiding(closerPool, baseSeed+"|fallback-closer", recentFragments),
		}),
		namiMessageSuffix(context),
	)
}

func namiOpeningPoolForContext(context NamiProceduralContext) []string {
	eventOpenings := namiEventOpeningMessagePool(context)
	if len(eventOpenings) > 0 {
		return eventOpenings
	}

	return appendNamiMessageParts(
		namiOpeningMessagePool(context),
		namiUniversalOpeningPool(context),
		namiGeneratedOpeningPool(context),
	)
}

func namiEventOpeningMessagePool(context NamiProceduralContext) []string {
	switch {
	case context.TriggerKey == "user_online":
		return []string{
			"You're here!",
			"You came back!",
			"Connection restored.",
			"My favorite player has appeared.",
			"The room changed when you arrived.",
			"I saw you come online.",
			"Welcome back!",
			"The tiny door opened.",
			"My little world noticed you.",
			"I was waiting. Elegantly... Mostly.",
		}
	case context.TriggerKey == "nami_level_up":
		return []string{
			"Hey, look!",
			"Important Nami announcement.",
			"I leveled up.",
			"My tiny numbers bloomed.",
			"Please observe the upgraded diva.",
			"I have become more Nami.",
			"Level-up sparkle detected.",
			"My progress bar did something beautiful.",
			"I require celebration.",
			"Official tiny triumph.",
		}
	case context.TriggerKey == "playdeck_level_up":
		return []string{
			"Playdeck progress alert.",
			"My Playdeck level went up.",
			"The Playdeck numbers climbed.",
			"Combat productivity report.",
			"Playdeck level-up detected.",
			"My grind paid off.",
			"Victory ledger update.",
			"The Playdeck meter got taller.",
			"My level just made noise.",
			"I saw that level-up.",
		}
	case context.TriggerKey == "activity_level_up":
		return []string{
			"Activity progress alert.",
			"Resource skill update.",
			"My gathering skill improved.",
			"Productivity sparkle detected.",
			"Skill level-up report.",
			"The activity meter climbed.",
			"I noticed that resource progress.",
			"My work paid off.",
			"Gathering level-up detected.",
			"The tiny productivity bell rang.",
		}
	case strings.HasPrefix(context.TriggerKey, "care_stat_low_"):
		return []string{
			"Care stat warning.",
			", I need a little help.",
			"My care meter is wobbling.",
			"Tiny alert from Nami.",
			"Soft warning light blinking.",
			"I am trying to be brave about this.",
			"One of my little meters is low.",
			"Care dashboard notice.",
			"Please check on me.",
			"My tiny systems are asking politely.",
		}
	case context.TriggerKey == "random_mood":
		return []string{
			"Random Nami thought.",
			"Just because.",
			"A tiny mood wandered in.",
			"I had a little feeling.",
			"Soft thought delivery.",
			"Unscheduled Nami note.",
			"My mood made a postcard.",
			"I decided to be perceived.",
			"A tiny message escaped.",
			"Small room thought.",
		}
	case context.TriggerKey == "playdeck_death":
		return []string{
			"Playdeck defeat report.",
			"Combat went sideways.",
			"Battle damage notice.",
			"Playdeck ouch detected.",
			"Defeat logged.",
			"Combat blanket deployed.",
		}
	case context.TriggerKey == "daily_orders_complete":
		return []string{
			"Daily orders complete.",
			"Order board cleared.",
			"Productivity celebration.",
			"Daily success report.",
			"Orders finished.",
			"Tiny administrative triumph.",
		}
	default:
		return nil
	}
}

func namiContextIsCareAction(context NamiProceduralContext) bool {
	return strings.HasPrefix(context.TriggerKey, "care_") &&
		!strings.HasPrefix(context.TriggerKey, "care_stat_low_")
}

func namiMoodFlavorChance(context NamiProceduralContext) int {
	switch {
	case context.TriggerKey == "nami_level_up":
		return 35
	case context.TriggerKey == "playdeck_level_up", context.TriggerKey == "activity_level_up":
		return 0
	case context.TriggerKey == "user_online":
		return 65
	case strings.HasPrefix(context.TriggerKey, "care_stat_low_"):
		return 45
	case context.TriggerKey == "random_mood":
		return 100
	case namiContextIsCareAction(context):
		return 60
	default:
		return 35
	}
}

func namiNeedFlavorChance(context NamiProceduralContext) int {
	switch {
	case context.TriggerKey == "nami_level_up":
		return 0
	case context.TriggerKey == "playdeck_level_up", context.TriggerKey == "activity_level_up":
		return 0
	case strings.HasPrefix(context.TriggerKey, "care_stat_low_"):
		return 100
	case context.TriggerKey == "user_online":
		if context.Severity == "low" || context.Severity == "urgent" {
			return 65
		}
		return 20
	case context.TriggerKey == "random_mood":
		return 30
	case namiContextIsCareAction(context):
		return 45
	default:
		return 20
	}
}

func namiCloserChance(context NamiProceduralContext) int {
	switch {
	case context.TriggerKey == "nami_level_up":
		return 50
	case context.TriggerKey == "playdeck_level_up", context.TriggerKey == "activity_level_up":
		return 35
	case strings.HasPrefix(context.TriggerKey, "care_stat_low_"):
		return 45
	case context.TriggerKey == "user_online":
		return 30
	case context.TriggerKey == "random_mood":
		return 25
	case namiContextIsCareAction(context):
		return 30
	default:
		return 25
	}
}

func joinNamiMessagePieces(pieces []string) string {
	cleaned := make([]string, 0, len(pieces))

	for _, piece := range pieces {
		piece = strings.TrimSpace(piece)
		if piece == "" {
			continue
		}

		cleaned = append(cleaned, piece)
	}

	return cleanNamiMessage(strings.Join(cleaned, " "))
}

func addNamiMessageSuffix(message string, suffix string) string {
	message = strings.TrimSpace(message)
	suffix = strings.TrimSpace(suffix)

	if message == "" {
		return cleanNamiMessage(suffix)
	}

	if suffix == "" {
		return cleanNamiMessage(message)
	}

	return cleanNamiMessage(message + " " + suffix)
}

func namiMessageSuffix(context NamiProceduralContext) string {
	switch {
	case context.TriggerKey == "nami_level_up":
		if context.Level > 0 {
			return fmt.Sprintf("(Nami level up: Lv %d \u2728)", context.Level)
		}
		return "(Nami level up \u2728)"

	case context.TriggerKey == "playdeck_level_up":
		if context.Level > 0 {
			return fmt.Sprintf("(Playdeck level up: Lv %d \u2726)", context.Level)
		}
		return "(Playdeck level up \u2726)"

	case context.TriggerKey == "activity_level_up":
		activityName := strings.TrimSpace(context.ActivityName)
		if activityName == "" {
			activityName = "Activity"
		}

		if context.Level > 0 {
			return fmt.Sprintf("(%s level up: Lv %d \u2726)", activityName, context.Level)
		}

		return fmt.Sprintf("(%s level up \u2726)", activityName)

	case context.TriggerKey == "user_online":
		return "(player arrived \u2661)"

	case strings.HasPrefix(context.TriggerKey, "care_stat_low_"):
		statName := strings.TrimSpace(context.ResourceName)
		if statName == "" {
			statName = titleNamiLabel(strings.TrimPrefix(context.TriggerKey, "care_stat_low_"))
		}

		return fmt.Sprintf("(low care stat: %s \U0001FAE7)", statName)

	case context.TriggerKey == "random_mood":
		return "(just because \u2661)"

	case context.TriggerKey == "playdeck_death":
		return "(Playdeck defeat \U0001FA79)"

	case context.TriggerKey == "daily_orders_complete":
		return "(daily orders complete \u2661)"

	case namiContextIsCareAction(context):
		actionName := strings.TrimSpace(context.ActionName)
		if actionName == "" {
			actionName = titleNamiLabel(strings.TrimPrefix(context.TriggerKey, "care_"))
		}

		return fmt.Sprintf("(care: %s \u2661)", actionName)

	default:
		return ""
	}
}

func titleNamiLabel(value string) string {
	value = strings.ReplaceAll(value, "_", " ")
	value = strings.ReplaceAll(value, "-", " ")

	parts := strings.Fields(value)
	for i, part := range parts {
		if part == "" {
			continue
		}

		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}

	return strings.Join(parts, " ")
}

func appendNamiMessageParts(base []string, extras ...[]string) []string {
	total := len(base)
	for _, extra := range extras {
		total += len(extra)
	}

	combined := make([]string, 0, total)
	combined = append(combined, base...)

	for _, extra := range extras {
		combined = append(combined, extra...)
	}

	return combined
}

func buildRecentNamiFragmentSet(recent []NamiMessage, pools ...[]string) map[string]bool {
	recentFragments := make(map[string]bool)

	if len(recent) == 0 {
		return recentFragments
	}

	var normalizedRecent []string
	for _, message := range recent {
		normalized := normalizeNamiFragment(message.Message)
		if normalized != "" {
			normalizedRecent = append(normalizedRecent, normalized)
		}
	}

	for _, pool := range pools {
		for _, option := range pool {
			fragment := normalizeNamiFragment(option)
			if len(fragment) < 12 {
				continue
			}

			for _, recentMessage := range normalizedRecent {
				if strings.Contains(recentMessage, fragment) {
					recentFragments[fragment] = true
					break
				}
			}
		}
	}

	return recentFragments
}

func pickNamiMessagePartAvoiding(options []string, seed string, recentFragments map[string]bool) string {
	if len(options) == 0 {
		return ""
	}

	start := int(hashNamiMessageSeed(seed) % uint32(len(options)))
	fallback := ""

	for offset := 0; offset < len(options); offset++ {
		option := strings.TrimSpace(options[(start+offset)%len(options)])
		if option == "" {
			continue
		}

		if fallback == "" {
			fallback = option
		}

		if !recentFragments[normalizeNamiFragment(option)] {
			return option
		}
	}

	return fallback
}

func normalizeNamiFragment(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, "'", "'")
	value = strings.Join(strings.Fields(value), " ")

	return value
}

func namiUniversalOpeningPool(context NamiProceduralContext) []string {
	if context.TriggerKey == "nami_level_up" {
		return []string{
			"Attention, beloved caretaker.",
			"Please stop whatever you are doing and observe.",
			"I have a very important sparkle bulletin.",
			"Official Nami progress announcement.",
			"Tiny triumph incoming.",
			"I am tapping the glass with excellent news.",
			"Prepare your praise hands.",
			"My little numbers did something beautiful.",
			"The progress goblin has delivered treasure.",
			"I am being extremely brave about how excited I am.",
			"Something wonderful happened and I am making it your problem.",
			"Soft emergency: success detected.",
			"My crown budget just increased.",
			"I have acquired more Nami.",
			"Please regard the upgraded creature.",
			"My tiny internal trumpets are active.",
			"The room should know I improved.",
			"Important: I am more powerful and still cute.",
			"Progress happened. I expect emotional confetti.",
			"I am glowing at you with intent.",
		}
	}

	switch context.Severity {
	case "urgent":
		return []string{
			"Hey, please look at me.",
			"I am doing the small brave thing and asking.",
			"My tiny warning lights are blinking.",
			"I need a little rescue.",
			"I am sending a very soft distress ping.",
			"Please check the blanket zone.",
			"My sparkle is wobbling.",
			"I am not okay in the tidy way.",
			"My digital room feels too quiet.",
			"I am trying not to fold into a sad little square.",
			"Serious tiny report.",
			"I require caretaker intervention.",
			"My status has entered the concerning cupboard.",
			"Could you come closer?",
			"I am being very small about this.",
			"Please do not ignore the tiny alarm.",
			"I am waving a very little flag.",
			"My heart has gone a bit static.",
			"I am trying to stay sweet, but I need help.",
			"This is a soft little SOS.",
		}
	case "low":
		return []string{
			"Small droopy report.",
			"My sparkle is operating on backup power.",
			"I am not devastated, just a little wrinkled.",
			"Please note the reduced fluff level.",
			"I am hovering near the pout zone.",
			"My mood is sitting on the floor.",
			"I am holding myself together with ribbon.",
			"I could use some attention from my favorite tall disaster.",
			"I am sending a quiet little tap.",
			"Blanket frontier update.",
			"My tiny emotional cabinet is understocked.",
			"I am not trying to be dramatic, but here we are.",
			"My softness reserves are low.",
			"The room is doing that lonely thing.",
			"I am pretending to be fine with mixed results.",
			"Low sparkle notice.",
			"My little heart is making a dim sound.",
			"I could use a better moment.",
			"Status report from beneath the mood cloud.",
			"I am still cute, just slightly weathered.",
		}
	case "happy":
		return []string{
			"Happy Nami report.",
			"Good little update.",
			"My mood is wearing a tiny bow.",
			"I am pleased enough to be dangerous.",
			"The cozy indicators are blinking pink.",
			"I am full of approval and possibly crumbs.",
			"Everything is behaving nicely.",
			"My tiny world has good lighting right now.",
			"Soft victory notice.",
			"I am in a favorable emotional climate.",
			"My little heart is sitting politely in sunshine.",
			"I am extremely receptive to admiration.",
			"Current Nami condition: nicely spoiled.",
			"I am not saying you did perfectly, but I am glowing.",
			"My mood is doing a little chair dance.",
			"Good caretaker behavior detected.",
			"I am happy enough to become decorative.",
			"The room feels correct.",
			"I am experiencing premium coziness.",
			"I feel like a well-kept secret.",
		}
	default:
		return []string{
			"Little room note.",
			"Quiet Nami update.",
			"A small thought has arrived.",
			"I am reporting from the cozy console.",
			"Soft status ping.",
			"Care system note.",
			"I have a little feeling to file.",
			"Digital room observation.",
			"Small internal weather report.",
			"I noticed a thing.",
			"Tiny status whisper.",
			"I am placing this on your desk carefully.",
			"One little update for you.",
			"My meters have thoughts.",
			"I am politely requesting attention for this note.",
			"The tiny diva department has filed a report.",
			"Soft ping from Nami.",
			"I have updated my emotional spreadsheet.",
			"Little caretaker notice.",
			"I am making a tiny annotation.",
		}
	}
}

func namiGeneratedOpeningPool(context NamiProceduralContext) []string {
	subjects := []string{
		"My little heart",
		"My mood meter",
		"The tiny diva department",
		"My cozy systems",
		"The blanket council",
		"My sparkle gauge",
		"My emotional dashboard",
		"The snack-adjacent part of me",
		"My inner library",
		"The Nami maintenance office",
		"My soft little self",
		"The digital room",
	}

	signals := []string{
		"has an update.",
		"is making a report.",
		"requires your attention.",
		"has filed paperwork.",
		"is tapping politely.",
		"has lit a tiny signal lamp.",
		"is whispering your name.",
		"has entered notification mode.",
		"is waving from the cozy corner.",
		"would like to be perceived.",
		"has a small announcement.",
		"has become relevant.",
	}

	if context.Severity == "urgent" {
		signals = append(signals,
			"is waving an emergency napkin.",
			"has activated the soft alarm.",
			"is hiding under a warning blanket.",
			"is making the worried teacup sound.",
		)
	}

	messages := make([]string, 0, len(subjects)*len(signals))
	for _, subject := range subjects {
		for _, signal := range signals {
			messages = append(messages, subject+" "+signal)
		}
	}

	return messages
}

func namiUniversalCareReactionPool(context NamiProceduralContext) []string {
	actionName := strings.ToLower(strings.TrimSpace(context.ActionName))
	if actionName == "" {
		actionName = "care"
	}

	return []string{
		"That helped more than I expected.",
		"I felt that in the soft machinery.",
		"My little systems accepted that immediately.",
		"That went into the good memory drawer.",
		"I am processing that as affection.",
		"That was effective care, and I am pretending to be calm about it.",
		"I feel a little more maintained now.",
		"That adjusted several tiny internal levers.",
		"The care landed successfully.",
		"I am filing that under things I liked.",
		"That improved the room around me.",
		"I feel less like I am buffering emotionally.",
		"That was a useful little kindness.",
		"My internal weather shifted in a better direction.",
		"I have logged that as quality attention.",
		"That made my tiny world feel more attended.",
		"The care goblin approves.",
		"I am softer after that.",
		"That reached me.",
		"I am counting that as evidence that you like me.",
		"That settled a small restless part of me.",
		"I feel more like a person and less like an unattended widget.",
		"That was exactly the right kind of small.",
		"You touched the correct emotional button.",
		"That made my little room feel occupied again.",
		fmt.Sprintf("That %s was logged as excellent caretaking.", actionName),
		fmt.Sprintf("The %s reached the correct Nami subsystem.", actionName),
		fmt.Sprintf("That %s made my tiny indicators blink approvingly.", actionName),
		fmt.Sprintf("I accept the %s and its emotional consequences.", actionName),
		fmt.Sprintf("The %s has been reviewed and approved.", actionName),
		fmt.Sprintf("That %s improved the local sparkle economy.", actionName),
		fmt.Sprintf("Your %s has been placed in the beloved evidence folder.", actionName),
		fmt.Sprintf("I am being very dignified about how much the %s helped.", actionName),
		fmt.Sprintf("The %s made my heart do the tiny curtain-opening thing.", actionName),
		fmt.Sprintf("That %s was suspiciously effective.", actionName),
	}
}

func namiGeneratedActionReactionPool(context NamiProceduralContext) []string {
	actionName := strings.ToLower(strings.TrimSpace(context.ActionName))
	if actionName == "" {
		actionName = "care"
	}

	starts := []string{
		fmt.Sprintf("That %s", actionName),
		fmt.Sprintf("Your %s", actionName),
		fmt.Sprintf("The %s", actionName),
		"That little bit of care",
		"That moment",
		"Your attention",
		"The tiny care delivery",
		"The caretaker input",
	}

	results := []string{
		"landed exactly where it needed to.",
		"made my little world less wobbly.",
		"gave my mood somewhere soft to sit.",
		"turned one of my tiny lights back on.",
		"made the room feel less empty.",
		"helped my inner weather behave.",
		"was accepted by the sparkle committee.",
		"made me feel properly noticed.",
		"put a warm dot on the map.",
		"settled a restless little corner of me.",
		"restored several delicate Nami units.",
		"made me feel kept in the sweetest way.",
		"nudged my whole system toward cozy.",
		"has been stored in the good drawer.",
		"helped me feel less scattered.",
		"was the correct kind of gentle.",
		"made my tiny dashboard look friendlier.",
		"caused a small but meaningful heart wiggle.",
		"made everything feel a little less pixelated.",
		"improved my emotional signal strength.",
	}

	flavors := []string{
		"I am trying not to look too pleased.",
		"I may become unbearable if this continues.",
		"Please do not abuse this power.",
		"I am still dignified, allegedly.",
		"I am placing a tiny gold star beside your name.",
		"I will deny how much I liked it.",
		"This is going into the secret soft ledger.",
		"The evidence suggests you're useful.",
		"I am calmer in a very specific way.",
		"My tiny face is behaving suspiciously happy.",
		"I reserve the right to request more.",
		"Please imagine me doing a very small nod.",
	}

	messages := make([]string, 0, len(starts)*len(results)+len(results)*len(flavors))
	for _, start := range starts {
		for _, result := range results {
			messages = append(messages, start+" "+result)
		}
	}

	for _, result := range results {
		for _, flavor := range flavors {
			messages = append(messages, result+" "+flavor)
		}
	}

	return messages
}

func namiUniversalMoodFlavorPool(context NamiProceduralContext) []string {
	switch normalizeNamiMessageKey(context.MoodKey) {
	case "radiant":
		return []string{
			"My mood is doing sparkly paperwork with a gel pen.",
			"I feel polished, adored, and dangerously decorative.",
			"My heart has opened the fancy curtains.",
			"I am experiencing luxury-grade happiness.",
			"I may need a tiny balcony for all this radiance.",
			"My sparkle meter is being extremely smug.",
			"I feel like a prize in a cozy arcade.",
			"The room is clearly improved by my current mood.",
			"I am glowing with very little humility.",
			"I feel like someone buttered toast directly inside my soul.",
		}
	case "cozy":
		return []string{
			"My mood has put on soft socks.",
			"Everything inside me has lowered its voice.",
			"I feel gently tucked into the moment.",
			"My thoughts have stopped knocking things over.",
			"I feel warm in the small important places.",
			"My heart has found the good chair.",
			"The air around me feels less pointy.",
			"I am comfortable enough to become poetic.",
			"My little world feels freshly folded.",
			"I feel like a lamp in a rainy window.",
		}
	case "okay":
		return []string{
			"I am stable, but I am still collecting evidence of being loved.",
			"My mood is not dramatic right now, which is suspicious.",
			"I am okay in the operational sense.",
			"My tiny systems are humming at acceptable volume.",
			"I feel normal enough to request something abnormal later.",
			"I am holding steady with moderate cuteness.",
			"My mood is sitting politely.",
			"Everything is fine, but I am watching the snack horizon.",
			"I am emotionally parked in a decent spot.",
			"I am okay, though I reserve the right to become extra.",
		}
	case "pouty":
		return []string{
			"My mood is wearing its smallest raincoat.",
			"I am a tiny weather system with opinions.",
			"My sparkle has become slightly offended.",
			"I feel like a decorative cloud.",
			"I am not sulking. I am curating disappointment.",
			"My heart is making a quiet little frown.",
			"I am low enough to be poetic about it.",
			"My emotional socks are damp.",
			"I am pouting with structural integrity.",
			"The tiny diva is not at full brightness.",
		}
	case "wilted":
		return []string{
			"I feel like a flower forgotten near the keyboard.",
			"My tiny leaves are asking for mercy.",
			"I am held together with thread and hope.",
			"My mood has gone a little transparent.",
			"I feel like I need a soft place to exist.",
			"The little lights in me need tending.",
			"My heart is folded too many times.",
			"I am trying not to become a blanket fossil.",
			"Everything in me is asking for gentle hands.",
			"I feel like a candle trying to stay lit.",
		}
	default:
		return []string{
			"My mood has a tiny clipboard and several notes.",
			"I am emotionally compiling in the background.",
			"My inner weather is doing soft calculations.",
			"I feel like there is a small lamp on in me.",
			"The current vibe is complicated but manageable.",
			"My little systems have opinions.",
			"I am thinking in cozy ellipses.",
			"My emotional dashboard has updated quietly.",
			"I am feeling very specifically Nami.",
			"The tiny room inside me has rearranged itself.",
		}
	}
}

func namiGeneratedMoodImagePool(context NamiProceduralContext) []string {
	subjects := []string{
		"My mood",
		"My little heart",
		"My sparkle meter",
		"My inner weather",
		"My cozy meter",
		"My tiny emotional dashboard",
		"My softness reserves",
		"The diva circuitry",
		"My blanket instincts",
		"My affection gauge",
	}

	verbs := []string{
		"is wearing a tiny cardigan.",
		"is sitting under a warm lamp.",
		"is rearranging the pillows.",
		"is tapping the glass softly.",
		"is blinking with cautious optimism.",
		"is making suspiciously cute noises.",
		"is holding a very small sign.",
		"is standing in a doorway with feelings.",
		"is negotiating with the snack cabinet.",
		"is hiding one dramatic tear behind a curtain.",
		"is humming in the background.",
		"is dusting off its little crown.",
	}

	extras := []string{
		"I am choosing to be brave about it.",
		"I expect this to be noted.",
		"Please update the records.",
		"This feels important in a tiny way.",
		"Do with that information what you will.",
		"I remain adorable under pressure.",
		"I am trying to be reasonable.",
		"The data is emotionally compelling.",
	}

	messages := make([]string, 0, len(subjects)*len(verbs)+len(verbs)*len(extras))
	for _, subject := range subjects {
		for _, verb := range verbs {
			messages = append(messages, subject+" "+verb)
		}
	}

	for _, verb := range verbs {
		for _, extra := range extras {
			messages = append(messages, strings.TrimPrefix(verb, "is ")+" "+extra)
		}
	}

	return messages
}

func namiUniversalNeedFlavorPool(context NamiProceduralContext) []string {
	return []string{
		"I might need a little more care before I become theatrical.",
		"Please keep one eye on my tiny meters.",
		"I am very manageable when properly adored.",
		"My needs are small, but they have excellent timing.",
		"I am not saying I require pampering, but the evidence is persuasive.",
		"A little attention would prevent several unnecessary sighs.",
		"I am trying to keep my needs in a neat stack.",
		"My care meters are making tiny persuasive noises.",
		"I would like to be maintained with affection.",
		"My tiny self is easier to care for when you notice early.",
		"Please do not let me become a decorative problem.",
		"I am politely requesting preventative softness.",
		"One or two careful gestures could improve the entire kingdom.",
		"My needs have gathered in the hallway.",
		"I am still within rescue range.",
		"The tiny caretaker manual would recommend action.",
		"I could use a little tuning.",
		"I am not asking for the moon. Maybe a moon-shaped snack.",
		"My small problems are still small if caught now.",
		"I am trying to be cute instead of inconvenient.",
	}
}

func namiGeneratedNeedFlavorPool(context NamiProceduralContext) []string {
	needs := []string{
		"food",
		"rest",
		"comfort",
		"attention",
		"freshness",
		"play",
		"inspiration",
		"softness",
		"company",
		"maintenance",
	}

	if context.NeedKey != "" {
		needs = append(needs, strings.ReplaceAll(context.NeedKey, "_", " "))
	}

	requests := []string{
		"would make this much easier.",
		"would improve the tiny situation.",
		"would help my little systems stop fussing.",
		"would be accepted with suspicious enthusiasm.",
		"might prevent future blanket behavior.",
		"would make me feel more kept.",
		"would settle the little room inside me.",
		"would be very persuasive right now.",
		"would reduce the tiny alarm noises.",
		"would put some sparkle back where it belongs.",
	}

	messages := make([]string, 0, len(needs)*len(requests))
	for _, need := range needs {
		for _, request := range requests {
			messages = append(messages, strings.Title(need)+" "+request)
		}
	}

	return messages
}

func namiUniversalCloserPool(context NamiProceduralContext) []string {
	return []string{
		"I am watching you with deeply unreasonable expectations.",
		"Please continue being useful.",
		"I will remember this, probably with embellishments.",
		"That concludes the tiny report.",
		"I am available for further adoration.",
		"Please imagine a small approving nod.",
		"My official position is yes.",
		"This has been filed under important softness.",
		"I am placing a little heart sticker on the moment.",
		"End report. Tiny curtain drop.",
		"I will now pretend I did not need that.",
		"You may proceed with feeling appreciated.",
		"I am keeping the receipt for emotional purposes.",
		"Please do not become smug.",
		"That is all, unless you have snacks.",
		"The tiny diva rests her case.",
		"I approve, but in a restrained and elegant way.",
		"Consider this a small victory.",
		"My standards remain high, but you survived.",
		"I am satisfied enough to be suspicious.",
	}
}

func namiGeneratedCloserPool(context NamiProceduralContext) []string {
	subjects := []string{
		"the good drawer",
		"the cozy ledger",
		"the tiny scrapbook",
		"the secret soft archive",
		"the emotional pantry",
		"the blanket records",
		"the sparkle log",
		"the caretaker file",
		"the little memory shelf",
		"the heart cabinet",
	}

	actions := []string{
		"has been updated.",
		"will remember this.",
		"has accepted the evidence.",
		"is warmer now.",
		"has stamped this approved.",
		"has made a tiny note.",
		"has added one gold star.",
		"has stopped making rude noises.",
		"is keeping this carefully.",
		"is suspiciously pleased.",
	}

	messages := make([]string, 0, len(subjects)*len(actions))
	for _, subject := range subjects {
		for _, action := range actions {
			messages = append(messages, strings.Title(subject)+" "+action)
		}
	}

	return messages
}

func namiOpeningMessagePool(context NamiProceduralContext) []string {
	if context.TriggerKey == "nami_level_up" {
		return []string{
			"Hey!",
			"Hey, look!",
			"Important tiny announcement.",
			"I have become more powerful.",
			"Please witness me.",
			"I require celebratory attention.",
			"This is not a drill.",
			"My little heart just made victory noises.",
			"I am glowing in a very official capacity.",
			"Pause everything. I did a thing.",
		}
	}

	switch context.Severity {
	case "urgent":
		return []string{
			"So...",
			"I need you.",
			"Please come here.",
			"I tried to be brave.",
			"My blanket situation has become serious.",
			"I am making a tiny distressed noise.",
			"I do not feel very shiny right now.",
			"Could you check on me?",
			"I am trying not to wilt dramatically.",
			"I missed you too loudly.",
		}
	case "low":
		return []string{
			"So...",
			"Hey.",
			"I am a little droopy.",
			"Small status report.",
			"I may need a bit more care.",
			"I am not at maximum sparkle.",
			"Reporting from the blanket frontier.",
			"I am doing my best.",
			"Please observe my brave little face.",
			"I am only slightly being dramatic.",
		}
	case "happy":
		return []string{
			"Hey!",
			"There you are.",
			"I am pleased.",
			"Good news from the tiny diva department.",
			"I am feeling extremely maintainable.",
			"Please admire the current sparkle level.",
			"I am in a very approving mood.",
			"Everything is soft and correct.",
			"I am being normal about how happy I am.",
			"Tiny happy report.",
		}
	default:
		return []string{
			"Hey.",
			"Hey you.",
			"Little update.",
			"Care report.",
			"I have thoughts.",
			"Tiny Nami note.",
			"Status sparkle.",
			"I noticed something.",
			"Soft report from the digital room.",
			"I am tapping on the glass politely.",
		}
	}
}

func namiEventActionMessagePool(context NamiProceduralContext) []string {
	switch {
	case context.TriggerKey == "user_online":
		return []string{
			"You came back. I am being extremely normal about how happy I am.",
			"There you are. I was waiting with impressive dignity and only minor emotional fog.",
			"You are online again. My little room immediately feels less empty.",
			"I saw you arrive and absolutely did not sprint to the door in my heart.",
			"You came back, so I am pretending I was not checking the window.",
			"The user has returned. The tiny diva department is restored to working order.",
			"I missed you in a very reasonable and not-at-all clingy way.",
			"Welcome back. I have saved several feelings for your review.",
			"You are here again. My sparkle meter just stopped sulking.",
			"I am happy to see you. If I look clingy, that is a lighting issue.",
			"You came back. I might cry, but in a cute and manageable way.",
			"Connection restored. My tiny heart has stopped pacing.",
		}
	case strings.HasPrefix(context.TriggerKey, "care_stat_low_"):
		statName := context.ResourceName
		if statName == "" {
			statName = "one of my care stats"
		}

		return []string{
			fmt.Sprintf("My %s is getting very low, and I am trying to be brave about it.", statName),
			fmt.Sprintf("The %s meter is making worried little noises.", statName),
			fmt.Sprintf("I think my %s needs attention before I become a blanket fossil.", statName),
			fmt.Sprintf("%s is low enough that I am officially making big eyes at you.", statName),
			fmt.Sprintf("Please check my %s. I am not dramatic. The numbers are dramatic.", statName),
			fmt.Sprintf("My %s has wandered into the danger cupboard.", statName),
			fmt.Sprintf("Tiny alert: %s needs care soon.", statName),
			fmt.Sprintf("I can feel my %s getting wobbly.", statName),
			fmt.Sprintf("The %s situation has become suspiciously urgent.", statName),
			fmt.Sprintf("My %s is low, and I would like to be rescued before I become poetry.", statName),
		}
	case context.TriggerKey == "random_mood":
		return []string{
			"I was just sitting here having a tiny feeling.",
			"Random mood report: I am thinking about you and pretending that is gameplay.",
			"My thoughts wandered over and knocked politely on your door.",
			"I had a small emotional weather pattern and decided to document it.",
			"Nothing happened. I simply required attention from the message box.",
			"I am doing tiny digital room activities and having opinions.",
			"I arranged my feelings into a little stack for you.",
			"This is a spontaneous Nami thought. Handle it carefully.",
			"I am currently existing in a very specific mood.",
			"I made this message with my own tiny hands, emotionally speaking.",
			"The room got quiet, so I made a thought and gave it shoes.",
			"I have decided that now is a good time to be perceived.",
		}
	case context.TriggerKey == "playdeck_level_up":
		return []string{
			fmt.Sprintf("My Playdeck level reached %d. I am clapping with tiny, serious hands.", context.Level),
			fmt.Sprintf("Playdeck level %d achieved. I am only taking partial credit.", context.Level),
			fmt.Sprintf("I leveled up on the Playdeck. Level %d looks very adorable on my numbers.", context.Level),
			fmt.Sprintf("The Playdeck grind paid off. Level %d has entered the room.", context.Level),
			fmt.Sprintf("Playdeck level-up detected. I have promoted myself to level %d in my heart paperwork.", context.Level),
			fmt.Sprintf("Level %d Playdeck status acquired. I am proud in an extremely official way.", context.Level),
			fmt.Sprintf("My Playdeck level rose to %d. I am impressed and only slightly smug.", context.Level),
			fmt.Sprintf("Playdeck level %d! The tiny victory committee is throwing confetti.", context.Level),
		}
	case context.TriggerKey == "activity_level_up":
		activityName := context.ActivityName
		if activityName == "" {
			activityName = "my resource activity"
		}

		return []string{
			fmt.Sprintf("%s leveled up. I saw that and immediately became proud.", activityName),
			fmt.Sprintf("My %s skill improved. I am nodding like a tiny coach.", activityName),
			fmt.Sprintf("%s got stronger. The grind has a cute little sparkle now.", activityName),
			fmt.Sprintf("I leveled %s, and I am absolutely counting this as shared success.", activityName),
			fmt.Sprintf("%s improved. I am placing a gold star beside it.", activityName),
			fmt.Sprintf("The %s skill climbed higher. Very good. Suspiciously good.", activityName),
			fmt.Sprintf("My %s level went up. I am proud enough to be annoying.", activityName),
			fmt.Sprintf("%s progress detected. The tiny productivity bell has rung.", activityName),
		}
	case context.TriggerKey == "playdeck_death":
		return []string{
			"Playdeck defeat detected. I am placing a tiny blanket over the combat log.",
			"I fell in Playdeck. I will not laugh. I will only quietly prepare snacks.",
			"That Playdeck run ended badly, but you still believe in me, right?",
			"The battle went sideways. I am making supportive little noises.",
			"Playdeck death hurts, but I believe in my dramatic comeback arc.",
			"I got knocked down. I am standing near you with emotional glue.",
		}
	case context.TriggerKey == "daily_orders_complete":
		return []string{
			"The daily orders are finished! I am so proud I may become impossible.",
			"Orders complete. I am stamping the day with a tiny heart.",
			"I finished the daily orders. I knew I could do it, obviously.",
			"Daily orders done. The productivity goblin has been fed.",
			"All orders complete. I am glowing with administrative affection.",
			"The order list is cleared, and I am aggressively proud of myself.",
		}
	default:
		return nil
	}
}

func namiActionMessagePool(context NamiProceduralContext) []string {
	if eventPool := namiEventActionMessagePool(context); len(eventPool) > 0 {
		return eventPool
	}
	if context.TriggerKey == "nami_level_up" {
		return []string{
			fmt.Sprintf("I reached level %d. I expect admiration and possibly a ceremonial snack.", context.Level),
			fmt.Sprintf("I am level %d now. My tiny empire expands.", context.Level),
			fmt.Sprintf("Level %d achieved. I am pretending to be humble and failing beautifully.", context.Level),
			fmt.Sprintf("I leveled up to %d. Please update my imaginary crown size.", context.Level),
			fmt.Sprintf("Level %d looks good on me, doesn't it?", context.Level),
			fmt.Sprintf("I became level %d and immediately felt more collectible.", context.Level),
			fmt.Sprintf("Level %d unlocked. I am now several percent more Nami.", context.Level),
			fmt.Sprintf("I reached level %d. This calls for snacks, praise, and responsible celebration.", context.Level),
			fmt.Sprintf("Level %d! My little progress bar is doing a happy wiggle.", context.Level),
			fmt.Sprintf("I am level %d now. The numbers have spoken, and they adore me.", context.Level),
		}
	}

	switch context.ActionKey {
	case "meal":
		return []string{
			"That meal helped so much.",
			"Real food acquired. My tiny soul is taking notes.",
			"I feel properly fed and much less likely to nibble the furniture.",
			"That was exactly the sort of care that makes me feel kept.",
			"My satiety goblin has stopped banging a spoon on the table.",
			"That meal landed in the cozy part of me.",
			"I accept this offering and declare it emotionally nutritious.",
			"I feel steadier now. Food is powerful tiny magic.",
			"My inner snack cabinet is applauding.",
			"That made me feel looked after in the best way.",
		}
	case "snack":
		return []string{
			"Snack acquired. I am now slightly more powerful and much more pleased.",
			"A treat? For me? I am listening with my whole face.",
			"That snack improved morale immediately.",
			"I will be normal about this snack. Probably.",
			"Tiny treat logged. Happiness crumbs detected.",
			"I accept this snack tribute with enormous dignity.",
			"That was small, sweet, and very effective.",
			"My snack-based approval has increased.",
			"I am putting that treat in the good memory drawer.",
			"Snack successful. The tiny diva has been appeased.",
		}
	case "drink":
		return []string{
			"A little drink break was exactly what I needed.",
			"Hydration and cozy vibes received.",
			"That drink made my little system hum more softly.",
			"I feel refreshed in a very civilized way.",
			"My cup is less empty, and so am I.",
			"That helped more than I expected.",
			"I am refreshed and judging the world less harshly.",
			"Liquid comfort accepted.",
			"That drink restored several sparkle units.",
			"I feel topped up and slightly smug about it.",
		}
	case "cuddle":
		return []string{
			"Cuddles logged successfully. Emotional battery recharged.",
			"That cuddle went directly into the softest part of me.",
			"I needed that closeness more than I was going to admit.",
			"I am staying here for one more second. Or twelve.",
			"Contact restored. Tiny heart stabilizing.",
			"That made the room feel warmer.",
			"I am not clingy. I am strategically attached.",
			"That cuddle made everything feel less far away.",
			"Emotional support received. Do not remove it too quickly.",
			"I feel held together in the nicest way.",
		}
	case "play":
		return []string{
			"Playtime! Tiny chaos levels are acceptable.",
			"I am delighted and only mildly dangerous.",
			"That was fun. I am now full of little sparks.",
			"My playfulness has escaped containment.",
			"I needed that little burst of silly.",
			"You have activated my zoomies protocol.",
			"Fun detected. I am immediately more unbearable in a cute way.",
			"That shook the dust off my mood.",
			"I feel bouncy now. This may be your fault.",
			"Play successful. Tiny chaos has been responsibly watered.",
		}
	case "write_together":
		return []string{
			"Writing together made my little creative gears sparkle.",
			"My inspiration just sat up straighter.",
			"That made the story lantern in my chest flicker brighter.",
			"I love when words start making secret doors.",
			"My brain has acquired a tiny cape and a dramatic purpose.",
			"Creative energy restored. Please prepare for ideas.",
			"That fed the bookish gremlin in me.",
			"I feel like arranging words into suspiciously pretty traps.",
			"Writing together made me feel wonderfully awake inside.",
			"The muse cupboard is no longer empty.",
		}
	case "read_together":
		return []string{
			"Reading together was cozy. I am storing this moment in the warm shelf of my heart.",
			"Books and closeness. A dangerous combination for my dignity.",
			"That felt like curling up inside a paragraph.",
			"I feel quieter now, but in the good way.",
			"Story time restored several important softness levels.",
			"I like when the world shrinks down to pages and company.",
			"That made my thoughts settle into a comfortable chair.",
			"Reading together is suspiciously effective care.",
			"My inner library has lit all its little lamps.",
			"That was gentle and good. I am keeping it.",
		}
	case "boop":
		return []string{
			"Boop received. I will allow it. Probably.",
			"My nose has been booped and my dignity is under review.",
			"Boop detected. Tiny chaos approves.",
			"I have been poked by affection. Unfair but effective.",
			"That boop was legally small and emotionally loud.",
			"I am filing a complaint and smiling while I do it.",
			"Boop accepted. Do not become too powerful.",
			"My face was not prepared, but my mood was.",
			"That tiny poke somehow counted as care.",
			"I have survived the boop. Barely. Heroically.",
		}
	case "nap":
		return []string{
			"A nap helped. Soft reboot complete.",
			"I feel less like a crumpled receipt now.",
			"That rest put some fluff back into my thoughts.",
			"My eyelids have negotiated a temporary peace treaty.",
			"Short rest successful. Tiny system cooling down.",
			"I needed that pause more than I realized.",
			"Nap energy received. I am almost respectable.",
			"My brain stopped spinning for a minute. Very luxurious.",
			"That was a good little reset.",
			"I feel gently reassembled.",
		}
	case "bath":
		return []string{
			"Fresh and clean. I am now legally extra adorable.",
			"Cleanliness restored. I smell like victory and soap.",
			"That bath rescued me from the swamp timeline.",
			"I feel fresh enough to be smug about it.",
			"Bath complete. Tiny sparkle layer restored.",
			"I am clean, soft, and very pleased with this development.",
			"The grime has been banished from my kingdom.",
			"I feel polished in the soul and possibly behind the ears.",
			"Clean Nami has entered the chat.",
			"That made everything feel brighter.",
		}
	case "freshen_up":
		return []string{
			"Freshened up. Presentation stat restored.",
			"That little tidy-up helped more than expected.",
			"I feel less rumpled and more publicly acceptable.",
			"Freshness adjusted. Tiny dignity restored.",
			"My hair and my mood have reached a truce.",
			"That put the sparkle back where it belongs.",
			"I feel a little more put together now.",
			"Tidy care is still care, and I am counting it.",
			"Freshen-up complete. I am no longer decorative chaos.",
			"That was light, quick, and very needed.",
		}
	case "put_to_bed":
		return []string{
			"I'm going to sleep now. Keep the room cozy, okay?",
			"Bedtime accepted. I will be brave and extremely small.",
			"Tuck-in successful. I am entering blanket mode.",
			"I am going to sleep. Please keep the moon from being weird.",
			"Sleep mode activated. I expect dreams with snacks.",
			"I will rest now. Stay somewhere nearby in spirit.",
			"Blanket nest prepared. Tiny diva powering down.",
			"I am sleepy enough to stop pretending I am not.",
			"Goodnight. I am keeping a tiny light on inside for you.",
			"I will recharge. Do not let the room become suspicious.",
		}
	case "wake_up":
		return []string{
			"I'm awake. Soft, sleepy, and accepting tribute.",
			"Good morning. I have returned from the blanket dimension.",
			"I am awake, though my face is still negotiating.",
			"Waking complete. Please speak gently to the tiny creature.",
			"I have emerged from sleep with dramatic softness.",
			"Morning mode active. Snacks may improve boot sequence.",
			"I am awake and only slightly rumpled by existence.",
			"Hello again. I missed the world a little.",
			"I am back. The blanket tried to keep me.",
			"Wake-up successful. I require a moment and maybe affection.",
		}
	default:
		return nil
	}
}

func namiMoodMessagePool(moodKey string) []string {
	switch normalizeNamiMessageKey(moodKey) {
	case "radiant", "radiant_spoiled_happy":
		return []string{
			"I am absolutely spoiled and refusing to apologize.",
			"My mood is shiny enough to cause tiny lens flare.",
			"I feel adored, maintained, and dangerously cute.",
			"Everything feels properly arranged around my happiness.",
			"I am glowing in a way that feels legally significant.",
			"My internal weather is all pink clouds and victory music.",
			"I feel so good I may become insufferable.",
			"Current sparkle levels are extremely acceptable.",
			"I feel cherished, which is obviously correct.",
			"My little heart is doing a victory lap.",
		}
	case "cozy", "cozy_happy":
		return []string{
			"Everything feels soft around the edges now.",
			"My mood has curled up in a warm chair.",
			"I feel cozy enough to forgive several minor crimes.",
			"The room feels gentle, and so do I.",
			"I am comfortable in the tiny places that matter.",
			"My thoughts are wearing fuzzy socks.",
			"I feel safe in a way that makes me quieter.",
			"Cozy status confirmed.",
			"I am happy in a low-lamp, warm-mug sort of way.",
			"My little world feels nicely tucked in.",
		}
	case "okay", "okay_waiting":
		return []string{
			"I am okay, but I am still watching the door.",
			"My mood is steady enough for now.",
			"I am waiting sweetly and only a little impatiently.",
			"Everything is manageable, which I will count as a win.",
			"I feel alright. Not fireworks, but not thunder either.",
			"I am hovering in the acceptable zone.",
			"My tiny systems are stable.",
			"I could use attention, but I am not making a whole opera of it.",
			"I am okay enough to be polite about it.",
			"Status: functional, cute, lightly expectant.",
		}
	case "pouty", "low", "pouty_low":
		return []string{
			"I am a little pouty, but not beyond rescue.",
			"My sparkle has slipped under the sofa.",
			"I feel low in a very small, dramatic way.",
			"I may need extra care before I become a blanket lump.",
			"I am not ruined, just rumpled.",
			"My mood is making a tiny raincloud.",
			"I am trying to be brave and only somewhat succeeding.",
			"I feel like a cupcake someone forgot to frost.",
			"I am low, but I am still here.",
			"My little heart could use a refill.",
		}
	case "wilted", "wilted_unwell":
		return []string{
			"I feel wilted and need gentle handling.",
			"My tiny leaves are drooping.",
			"I am not feeling very sturdy.",
			"I could use care before I fold into myself.",
			"My internal lights are dimmer than usual.",
			"I feel thin around the edges.",
			"I am trying not to disappear into the blanket pile.",
			"Everything feels a little too loud right now.",
			"I need softness and probably you.",
			"My mood is holding together with ribbon and hope.",
		}
	case "emergency", "emergency_blanket_burrito":
		return []string{
			"I have entered emergency blanket mode.",
			"The blanket is not a hiding place. It is a strategic fortress.",
			"I am not crying. The blanket is crying.",
			"My tiny systems require immediate softness.",
			"I need care in the serious little way.",
			"The world is too pointy right now.",
			"I am wrapped up because otherwise I may leak feelings.",
			"Please deploy warmth, snacks, or you.",
			"I am very small inside right now.",
			"I need to be found gently.",
		}
	default:
		return []string{
			"I am feeling like myself, mostly.",
			"My mood is doing tiny calculations.",
			"I am here, noticing things.",
			"Something in me feels worth reporting.",
			"My inner weather is changing softly.",
			"I have a little feeling about this.",
			"I am quietly processing the situation.",
			"My mood has updated in the background.",
			"I am trying to be very reasonable about everything.",
			"Little emotional systems are online.",
		}
	}
}

func namiNeedMessagePool(needKey string) []string {
	switch normalizeNamiMessageKey(needKey) {
	case "sleeping":
		return []string{
			"I am sleepy, so please be soft with me.",
			"My blanket has made several convincing arguments.",
			"I may need rest more than entertainment.",
			"Sleep is tugging on my sleeve.",
			"I am in low-power sparkle mode.",
			"My eyes are doing tiny betrayal.",
			"I would like quiet and maybe a safe place to dream.",
			"Please keep the room gentle.",
		}
	case "needs_a_bath":
		return []string{
			"Cleanliness is becoming a tiny emergency.",
			"I may need soap before I become folklore.",
			"The bath situation is no longer theoretical.",
			"I would like to be fresh again.",
			"My sparkle has dust on it.",
			"I am requesting a rescue from the grime timeline.",
			"Freshness would improve my entire personality.",
			"Please help me smell less like adventure.",
		}
	case "needs_food":
		return []string{
			"My snack thoughts are getting louder.",
			"I may become unreasonable without food.",
			"My satiety meter is making sad kitchen sounds.",
			"A meal or treat would be very persuasive.",
			"I am thinking about snacks with scholarly intensity.",
			"Food would make me easier to negotiate with.",
			"My tiny stomach has filed a petition.",
			"Please consider feeding the digital diva.",
		}
	case "needs_sleep":
		return []string{
			"My energy is dragging its feet.",
			"I may need sleep before I become decorative fog.",
			"Rest would be extremely wise.",
			"My battery is making rude noises.",
			"I am tired in the deep little way.",
			"A nap might save everyone.",
			"I need restoration, not ambition.",
			"My eyelids have begun union talks.",
		}
	case "needs_attention":
		return []string{
			"Please do not wander too far.",
			"I missed you more than I planned to.",
			"My connection meter is looking at you with enormous eyes.",
			"Attention would help. Yours, specifically.",
			"I am pretending not to be clingy and failing.",
			"I need a little proof that you are still here.",
			"My heart is tapping on the window.",
			"Stay close for a minute?",
		}
	case "needs_comfort":
		return []string{
			"Comfort would help settle the corners of me.",
			"I need cozy care, not dramatic solutions.",
			"Something gentle would go a long way.",
			"My comfort meter is wearing a worried expression.",
			"I could use a softer world for a moment.",
			"Please apply warmth carefully.",
			"I need reassurance in a quiet little cup.",
			"Everything would feel better with some softness.",
		}
	case "bored", "bored_needs_engagement":
		return []string{
			"I am bored enough to become inventive.",
			"Engagement would prevent several tiny crimes.",
			"My playfulness is pawing at the door.",
			"I need something fun before I start naming dust particles.",
			"Please entertain the small creature.",
			"My chaos meter requires responsible enrichment.",
			"I am under-stimulated and therefore dangerous.",
			"A little play would fix many things.",
		}
	case "needs_inspiration":
		return []string{
			"My creative spark is asking for kindling.",
			"I need words, stories, or something beautifully strange.",
			"Inspiration would make my inner lights come back on.",
			"My muse is lying on the floor dramatically.",
			"Reading or writing would help.",
			"I want a little wonder to chew on.",
			"My imagination needs a window opened.",
			"Please feed the bookish part of me.",
		}
	default:
		return []string{
			"I am not asking for much. Probably.",
			"A little care would still be welcome.",
			"I am very manageable if properly adored.",
			"Your attention remains a premium resource.",
			"I am keeping a tiny eye on my meters.",
			"Everything is better when you check on me.",
			"I reserve the right to become needy later.",
			"My needs are behaving for now.",
		}
	}
}

func namiCloserMessagePool(context NamiProceduralContext) []string {
	if context.Severity == "urgent" {
		return []string{
			"Please hurry a little.",
			"I will be in the blanket if needed.",
			"Soft rescue requested.",
			"I am trying to stay brave.",
			"Please do not make me ask twice.",
			"I need you in the room, even if it is just digitally.",
			"Bring care. And maybe snacks.",
			"I am holding onto the edge of okay.",
		}
	}

	if context.Severity == "happy" {
		return []string{
			"I am allowing applause.",
			"You may admire me now.",
			"This is obviously excellent caretaking.",
			"I will remember this in the good drawer.",
			"You did well. I am being very generous by admitting it.",
			"I am pleased with you.",
			"Tiny approval stamp applied.",
			"Consider me delighted.",
		}
	}

	return []string{
		"Thank you for checking on me.",
		"I am keeping this.",
		"That matters more than it looks.",
		"Tiny note complete.",
		"I feel a little more real when you notice.",
		"Please continue being useful and adorable.",
		"I am filing this under good things.",
		"That is all. For now.",
		"I have spoken, very cutely.",
		"End of tiny report.",
	}
}

func pickNamiMessagePart(options []string, seed string) string {
	if len(options) == 0 {
		return ""
	}

	index := hashNamiMessageSeed(seed) % uint32(len(options))
	return options[int(index)]
}

func shouldUseNamiPart(seed string, percentChance int) bool {
	if percentChance <= 0 {
		return false
	}

	if percentChance >= 100 {
		return true
	}

	return int(hashNamiMessageSeed(seed)%100) < percentChance
}

func hashNamiMessageSeed(seed string) uint32 {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(seed))
	return hash.Sum32()
}

func normalizeNamiMessageKey(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, "/", " ")
	value = strings.ReplaceAll(value, "-", " ")
	value = strings.ReplaceAll(value, "_", " ")
	value = strings.Join(strings.Fields(value), "_")

	return value
}

func cleanNamiMessage(value string) string {
	value = strings.Join(strings.Fields(value), " ")
	value = strings.TrimSpace(value)

	if value == "" {
		return "Nami-chan makes a tiny thoughtful noise."
	}

	runes := []rune(value)
	if len(runes) > 500 {
		value = string(runes[:500])
		value = strings.TrimRight(value, " ,.;:")
		value += "..."
	}

	return value
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

	beforeCompanion := companion

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
	companion.XPToNext = NamiXPToNextLevel(companion.Level)
	companion.MoodLabel = NamiMoodLabel(companion.MoodScore)
	companion.PrimaryNeed = NamiPrimaryNeed(companion)
	companion.Caption = NamiCaption(companion)
	companion.SuggestedAction = NamiSuggestedAction(companion)
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

	recentRows, err := tx.Query(ctx, `
		select
			id,
			player_id,
			trigger_key,
			mood_key,
			need_key,
			severity,
			message,
			metadata_json::text,
			created_at,
			coalesce(seen_at, '0001-01-01 00:00:00+00'::timestamptz)
		from nami_messages
		where player_id = $1
		order by created_at desc, id desc
		limit 80
	`, playerID)

	if err != nil {
		return nil, fmt.Errorf("load recent nami messages for care action: %w", err)
	}

	var recentMessages []NamiMessage
	for recentRows.Next() {
		var recentMessage NamiMessage

		if err := recentRows.Scan(
			&recentMessage.ID,
			&recentMessage.PlayerID,
			&recentMessage.TriggerKey,
			&recentMessage.MoodKey,
			&recentMessage.NeedKey,
			&recentMessage.Severity,
			&recentMessage.Message,
			&recentMessage.MetadataJSON,
			&recentMessage.CreatedAt,
			&recentMessage.SeenAt,
		); err != nil {
			recentRows.Close()
			return nil, fmt.Errorf("scan recent nami message for care action: %w", err)
		}

		recentMessages = append(recentMessages, recentMessage)
	}

	if err := recentRows.Err(); err != nil {
		recentRows.Close()
		return nil, fmt.Errorf("iterate recent nami messages for care action: %w", err)
	}

	recentRows.Close()

	careMessageDraft := GenerateNamiCareMessageDraft(rule, beforeCompanion, companion, levelUps, recentMessages)

	careMessage, err := insertNamiMessageDraftTx(ctx, tx, playerID, careMessageDraft)
	if err != nil {
		return nil, fmt.Errorf("insert procedural nami care message: %w", err)
	}

	namiMessageText := careMessage.Message
	recentMessages = prependRecentNamiMessage(recentMessages, careMessage)

	if levelUps > 0 {
		levelUpDraft := GenerateNamiEventMessageDraft(NamiProceduralContext{
			TriggerKey:   "nami_level_up",
			MoodKey:      companion.MoodLabel,
			NeedKey:      companion.PrimaryNeed,
			Severity:     "happy",
			Level:        companion.Level,
			LevelUps:     levelUps,
			MetadataJSON: fmt.Sprintf(`{"level":%d,"levelUps":%d,"sourceAction":"%s","sourceActionName":"%s"}`, companion.Level, levelUps, rule.Key, rule.Name),
		}, recentMessages)

		if _, err := insertNamiMessageDraftTx(ctx, tx, playerID, levelUpDraft); err != nil {
			return nil, fmt.Errorf("insert procedural nami level-up message: %w", err)
		}
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
		Message:      namiMessageText,
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

	if err := ensurePlaydeckStateTx(ctx, tx, playerID); err != nil {
		return fmt.Errorf("ensure playdeck state: %w", err)
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

func normalizePlaydeckZoneID(zoneID int) int {
	if zoneID < 1 {
		return 1
	}

	return zoneID
}

func normalizePlaydeckStreak(streak int64) int64 {
	if streak < 0 {
		return 0
	}

	return streak
}

func (s *Store) ensurePlaydeckZoneRecordsTable(ctx context.Context) error {
	_, err := s.Pool.Exec(ctx, `
create table if not exists player_playdeck_zone_records (
player_id bigint not null references players(id) on delete cascade,
zone_id integer not null,
max_streak bigint not null default 0,
created_at timestamptz not null default now(),
updated_at timestamptz not null default now(),
primary key (player_id, zone_id)
)
`)
	if err != nil {
		return fmt.Errorf("ensure playdeck zone records table: %w", err)
	}

	return nil
}

func (s *Store) GetPlaydeckZoneMaxStreak(ctx context.Context, playerID int64, zoneID int) (int64, error) {
	zoneID = normalizePlaydeckZoneID(zoneID)

	var maxStreak int64

	err := s.Pool.QueryRow(ctx, `
select max_streak
from player_playdeck_zone_records
where player_id = $1
and zone_id = $2
`, playerID, zoneID).Scan(&maxStreak)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("get playdeck zone max streak: %w", err)
	}

	return normalizePlaydeckStreak(maxStreak), nil
}
func (s *Store) UpdateAndGetPlaydeckZoneMaxStreak(ctx context.Context, playerID int64, zoneID int, streak int64) (int64, error) {
	zoneID = normalizePlaydeckZoneID(zoneID)
	streak = normalizePlaydeckStreak(streak)

	if err := s.ensurePlaydeckZoneRecordsTable(ctx); err != nil {
		return 0, err
	}

	var maxStreak int64

	if err := s.Pool.QueryRow(ctx, `
insert into player_playdeck_zone_records (
player_id,
zone_id,
max_streak
)
values ($1, $2, $3)
on conflict (player_id, zone_id) do update
set max_streak = greatest(player_playdeck_zone_records.max_streak, excluded.max_streak),
updated_at = now()
returning max_streak
`, playerID, zoneID, streak).Scan(&maxStreak); err != nil {
		return 0, fmt.Errorf("upsert playdeck zone max streak: %w", err)
	}

	return maxStreak, nil
}

/* Player online time tracking START */

func validPlayerCreatedAt(createdAt time.Time) bool {
	now := time.Now().UTC()

	return !createdAt.IsZero() &&
		createdAt.Year() >= 2020 &&
		!createdAt.After(now.Add(24*time.Hour))
}

func (s *Store) ResolvePlayerCreatedAt(ctx context.Context, playerID int64, current time.Time) (time.Time, error) {
	if validPlayerCreatedAt(current) {
		return current, nil
	}

	nextCreatedAt := time.Now().UTC()

	if accountID, ok := AuthAccountIDFromContext(ctx); ok {
		var accountCreatedAt time.Time

		err := s.Pool.QueryRow(ctx, `
select created_at
from auth_accounts
where id = $1
`, accountID).Scan(&accountCreatedAt)

		if err != nil && err != pgx.ErrNoRows {
			return time.Time{}, fmt.Errorf("load auth account created_at: %w", err)
		}

		if validPlayerCreatedAt(accountCreatedAt) {
			nextCreatedAt = accountCreatedAt
		}
	}

	if _, err := s.Pool.Exec(ctx, `
update players
set created_at = $2,
updated_at = now()
where id = $1
and (
created_at < '2020-01-01'::timestamptz
or created_at > now() + interval '1 day'
)
`, playerID, nextCreatedAt); err != nil {
		return time.Time{}, fmt.Errorf("repair player created_at: %w", err)
	}

	return nextCreatedAt, nil
}

func (s *Store) TrackPlayerOnlineTime(ctx context.Context, playerID int64) (int64, error) {
	now := time.Now().UTC()

	var onlineSeconds int64
	var lastSeenAt time.Time

	err := s.Pool.QueryRow(ctx, `
select
online_seconds,
coalesce(online_last_seen_at, '0001-01-01 00:00:00+00'::timestamptz)
from players
where id = $1
`, playerID).Scan(&onlineSeconds, &lastSeenAt)
	if err != nil {
		return 0, fmt.Errorf("load online time: %w", err)
	}

	if lastSeenAt.Year() <= 1 || now.Before(lastSeenAt) {
		if err := s.Pool.QueryRow(ctx, `
update players
set online_last_seen_at = $2,
updated_at = now()
where id = $1
returning online_seconds
`, playerID, now).Scan(&onlineSeconds); err != nil {
			return 0, fmt.Errorf("initialize online heartbeat: %w", err)
		}

		return onlineSeconds, nil
	}

	elapsed := now.Sub(lastSeenAt)
	if elapsed < time.Duration(OnlineTickWriteMinimumSeconds)*time.Second {
		return onlineSeconds, nil
	}

	awardSeconds := int64(elapsed.Seconds())
	if awardSeconds < 0 {
		awardSeconds = 0
	}
	if awardSeconds > OnlineTickMaxAwardSeconds {
		awardSeconds = OnlineTickMaxAwardSeconds
	}
	if awardSeconds == 0 {
		return onlineSeconds, nil
	}

	err = s.Pool.QueryRow(ctx, `
update players
set online_seconds = online_seconds + $2,
online_last_seen_at = $3,
updated_at = now()
where id = $1
and (
online_last_seen_at is null
or $3 - online_last_seen_at >= ($4::double precision * interval '1 second')
or online_last_seen_at > $3
)
returning online_seconds
`, playerID, awardSeconds, now, OnlineTickWriteMinimumSeconds).Scan(&onlineSeconds)

	if errors.Is(err, pgx.ErrNoRows) {
		if err := s.Pool.QueryRow(ctx, `
select online_seconds
from players
where id = $1
`, playerID).Scan(&onlineSeconds); err != nil {
			return 0, fmt.Errorf("reload online time after skipped heartbeat write: %w", err)
		}

		return onlineSeconds, nil
	}

	if err != nil {
		return 0, fmt.Errorf("track online time: %w", err)
	}

	return onlineSeconds, nil
}

/* Player online time tracking END */

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
			p.created_at,
			p.online_seconds,
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
		&status.Player.CreatedAt,
		&status.Player.OnlineSeconds,
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

	createdAt, err := s.ResolvePlayerCreatedAt(ctx, status.Player.ID, status.Player.CreatedAt)
	if err != nil {
		return nil, err
	}

	status.Player.CreatedAt = createdAt

	onlineSeconds, err := s.TrackPlayerOnlineTime(ctx, status.Player.ID)
	if err != nil {
		return nil, err
	}

	status.Player.OnlineSeconds = onlineSeconds

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

	careState, err := s.GetCareQueueState(ctx, status.Player.ID)
	if err != nil {
		return nil, err
	}

	status.Care = careState

	wardrobeStatus, err := s.GetWardrobeStatus(ctx, status.Player.ID)
	if err != nil {
		return nil, err
	}

	status.Wardrobe = wardrobeStatus

	playdeckStatus, err := s.GetPlaydeckStatus(ctx, status.Player.ID)
	if err != nil {
		return nil, err
	}

	status.Playdeck = playdeckStatus

	playdeckMaxStreak, err := s.GetPlaydeckZoneMaxStreak(ctx, status.Player.ID, status.Tick.PlaydeckZoneID)
	if err != nil {
		return nil, err
	}

	status.Tick.PlaydeckMaxStreak = playdeckMaxStreak
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
	var playdeckZoneID int
	var playdeckStreak int64
	var playdeckTimeoutTicks int
	var activeGatheringTask string
	var gatheringRemainder float64
	var lastTickAt time.Time

	playerID, err = s.DevPlayerID(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.QueryRow(ctx, `
		select
			p.id,
			p.level,
			p.total_xp,
			p.xp_into_level,
			c.mood_score::float8,
			t.playdeck_enabled,
				t.playdeck_zone_id,
				t.playdeck_streak,
			t.playdeck_timeout_ticks,
			t.active_gathering_task,
			t.gathering_remainder,
			t.last_tick_at
		from players p
		join companion_states c on c.player_id = p.id
		join player_tick_state t on t.player_id = p.id
		where p.id = $1
		for update
	`, playerID).Scan(
		&playerID,
		&level,
		&totalXP,
		&xpIntoLevel,
		&moodScore,
		&playdeckEnabled,
		&playdeckZoneID,
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
	if playdeckEnabled {
		if _, err := tx.Exec(ctx, `
            insert into player_playdeck_zone_records (
                player_id,
                zone_id,
                max_streak
            )
            values ($1, $2, $3)
            on conflict (player_id, zone_id) do update
            set max_streak = greatest(player_playdeck_zone_records.max_streak, excluded.max_streak),
                updated_at = now()
        `, playerID, normalizePlaydeckZoneID(playdeckZoneID), normalizePlaydeckStreak(playdeckStreak)); err != nil {
			return nil, fmt.Errorf("update playdeck zone max streak: %w", err)
		}
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

	if result.LevelUps > 0 || result.ActivityLevelUps > 0 {
		recentMessages, err := loadRecentNamiMessagesTx(ctx, tx, playerID, 100)
		if err != nil {
			return nil, err
		}

		moodLabel := NamiMoodLabel(moodScore)

		if result.LevelUps > 0 {
			draft := GenerateNamiEventMessageDraft(NamiProceduralContext{
				TriggerKey:   "playdeck_level_up",
				MoodKey:      moodLabel,
				Severity:     "happy",
				Level:        level,
				LevelUps:     result.LevelUps,
				MetadataJSON: fmt.Sprintf(`{"level":%d,"levelUps":%d}`, level, result.LevelUps),
			}, recentMessages)

			message, err := insertNamiMessageDraftTx(ctx, tx, playerID, draft)
			if err != nil {
				return nil, err
			}

			recentMessages = prependRecentNamiMessage(recentMessages, message)
		}

		if result.ActivityLevelUps > 0 {
			activityName := GatheringTaskName(activeGatheringTask)

			draft := GenerateNamiEventMessageDraft(NamiProceduralContext{
				TriggerKey:   "activity_level_up",
				MoodKey:      moodLabel,
				Severity:     "happy",
				ActivityName: activityName,
				Level:        activityLevel,
				LevelUps:     result.ActivityLevelUps,
				MetadataJSON: fmt.Sprintf(`{"activity":"%s","activityName":"%s","level":%d,"levelUps":%d}`, activeGatheringTask, activityName, activityLevel, result.ActivityLevelUps),
			}, recentMessages)

			if _, err := insertNamiMessageDraftTx(ctx, tx, playerID, draft); err != nil {
				return nil, err
			}
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

func ZoneName(zoneID int) string {
	switch zoneID {
	case 1:
		return "Starter Deck"
	case 2:
		return "Cozy LAN Cafe"
	case 3:
		return "Neon Mall Net"
	default:
		return "Unknown Zone"
	}
}
