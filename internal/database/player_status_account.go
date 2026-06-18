package database

import (
	"context"
	"fmt"
	"math"
	"time"
)

func (s *Store) GetPlayerStatusForAccount(ctx context.Context, accountID int64) (*PlayerStatus, error) {
	if accountID < 1 {
		return nil, fmt.Errorf("account id must be positive")
	}

	playerID, err := s.PlayerIDForAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return s.GetPlayerStatusByID(ctx, playerID)
}

func (s *Store) GetPlayerStatusByID(ctx context.Context, playerID int64) (*PlayerStatus, error) {
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
where p.id = $1
`, playerID).Scan(
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
		return nil, fmt.Errorf("get player status by id: %w", err)
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
