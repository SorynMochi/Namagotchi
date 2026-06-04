package database

import (
	"context"
	"fmt"
	"hash/fnv"
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

func (s *Store) DevPlayerID(ctx context.Context) (int64, error) {
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

	if limit > 100 {
		limit = 100
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

	triggerKey := "care_" + rule.Key
	severity := "info"

	if levelUps > 0 {
		triggerKey = "nami_level_up"
		severity = "happy"
	}

	if after.MoodScore < 20 {
		severity = "urgent"
	} else if after.MoodScore < 40 {
		severity = "low"
	} else if after.MoodScore >= 75 {
		severity = "happy"
	}

	context := NamiProceduralContext{
		TriggerKey:   triggerKey,
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
		recentText[message.Message] = true
	}

	actionPool := namiActionMessagePool(context)
	moodPool := namiMoodMessagePool(context.MoodKey)
	needPool := namiNeedMessagePool(context.NeedKey)
	openingPool := namiOpeningMessagePool(context)
	closerPool := namiCloserMessagePool(context)

	if len(actionPool) == 0 {
		actionPool = []string{"I noticed that. I am placing it carefully in my little internal scrapbook."}
	}

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

	for attempt := 0; attempt < 160; attempt++ {
		seed := fmt.Sprintf("%s|attempt:%d", baseSeed, attempt)

		pieces := []string{
			pickNamiMessagePart(openingPool, seed+"|opening"),
			pickNamiMessagePart(actionPool, seed+"|action"),
		}

		if shouldUseNamiPart(seed+"|mood", 85) {
			pieces = append(pieces, pickNamiMessagePart(moodPool, seed+"|mood"))
		}

		if shouldUseNamiPart(seed+"|need", 65) {
			pieces = append(pieces, pickNamiMessagePart(needPool, seed+"|need"))
		}

		if shouldUseNamiPart(seed+"|closer", 35) {
			pieces = append(pieces, pickNamiMessagePart(closerPool, seed+"|closer"))
		}

		message := cleanNamiMessage(strings.Join(pieces, " "))
		if message != "" && !recentText[message] {
			return message
		}
	}

	return cleanNamiMessage(strings.Join([]string{
		pickNamiMessagePart(openingPool, baseSeed+"|fallback-opening"),
		pickNamiMessagePart(actionPool, baseSeed+"|fallback-action"),
		pickNamiMessagePart(closerPool, baseSeed+"|fallback-closer"),
	}, " "))
}

func namiOpeningMessagePool(context NamiProceduralContext) []string {
	if context.LevelUps > 0 || context.TriggerKey == "nami_level_up" {
		return []string{
			"Soryn!",
			"Soryn, look!",
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
			"Soryn...",
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
			"Soryn...",
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
			"Soryn!",
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
			"Soryn.",
			"Hey Soryn.",
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

func namiActionMessagePool(context NamiProceduralContext) []string {
	if context.LevelUps > 0 || context.TriggerKey == "nami_level_up" {
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
			"I’m going to sleep now. Keep the room cozy, okay?",
			"Bedtime accepted. I will be brave and extremely small.",
			"Tuck-in successful. I am entering blanket mode.",
			"I am going to sleep. Please keep the moon from being weird.",
			"Sleep mode activated. I expect dreams with snacks.",
			"I will rest now. Stay somewhere nearby in spirit.",
			"Blanket nest prepared. Tiny diva powering down.",
			"I am sleepy enough to stop pretending I am not.",
			"Goodnight, Soryn. I am keeping a tiny light on inside.",
			"I will recharge. Do not let the room become suspicious.",
		}
	case "wake_up":
		return []string{
			"I’m awake. Soft, sleepy, and accepting tribute.",
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
		return "Cozy LAN Café"
	case 3:
		return "Neon Mall Net"
	default:
		return "Unknown Zone"
	}
}
