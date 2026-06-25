# Exact Game Numbers

Namigotchi needs idle-game numbers that can grow far beyond PostgreSQL `bigint`, Go `int64`, and JavaScript `Number` while remaining exact.

This foundation uses three layers:

1. PostgreSQL stores growth values as exact `numeric(120, 0)` through the `game_integer` domain.
2. Go represents those values with `database.GameInt`, backed by `math/big.Int`.
3. Browser code receives those values as JSON strings and parses/displays them through `window.NamigotchiNumbers`.

## What should use `game_integer`

Use `game_integer` for values that can grow indefinitely:

- Player level
- Nami level
- Playdeck level
- Work/activity/resource levels
- XP and XP into level
- Credits/currency values
- Nibbles, NamiCoin, future currencies
- Fans, memes, receipts, patterns, confidence, lost items, ingredients, materials
- HP, max HP, attack, defense, damage, recovery values once combat scaling grows
- Wins, losses, streaks, kill counts, achievement counters

Use `signed_game_integer` only for values that may intentionally go negative, such as temporary deltas or preview comparisons.

## What should stay small

Do not use `game_integer` for bounded or fractional UI/control values:

- Database IDs and foreign keys
- Care meters that are intentionally 0-100
- Percentages and progress-bar percentages
- Decay/recovery remainders
- Durations, queue positions, sort orders, timestamps
- Probability/rate values

## JSON rule

Game integers must be sent to the browser as strings:

```json
{
  "totalXp": "123456789012345678901234567890",
  "level": "420000000000000000000"
}
```

Never send large game integers as JSON numbers. JavaScript `Number` loses precision above `9007199254740991`.

## Frontend helpers

`web/game-numbers.js` exposes:

```js
NamigotchiNumbers.parse("12345678901234567890")
NamigotchiNumbers.format("12345678901234567890")
NamigotchiNumbers.formatFull("12345678901234567890")
NamigotchiNumbers.formatRatio(current, maximum)
NamigotchiNumbers.formatCreditsFromCents(creditsCents)
NamigotchiNumbers.percent(current, maximum)
```

The compact formatter supports one-decimal idle notation through Decillion and beyond:

```text
1.2K
44.8M
9.7Qa
3.1Qi
8.4Sx
1.0Dc
```

## Loading the frontend helper

Add the helper before or after `app.js` in `web/index.html`:

```html
<script src="/theme-effects.js"></script>
<script src="/app.js"></script>
<script src="/game-numbers.js"></script>
```

The helper includes backward-compatible global hooks for existing display calls:

```js
formatCompactNumber(value)
formatWholeCredits(cents)
formatCredits(cents)
percent(current, maximum)
```

## Migration pattern for existing columns

When converting a column, use this pattern:

```sql
alter table players
alter column total_xp type game_integer using total_xp::numeric;
```

For grouped conversions:

```sql
alter table players
alter column level type game_integer using level::numeric,
alter column total_xp type game_integer using total_xp::numeric,
alter column xp_into_level type game_integer using xp_into_level::numeric,
alter column currency_cents type game_integer using currency_cents::numeric,
alter column nibbles type game_integer using nibbles::numeric,
alter column namicoin type game_integer using namicoin::numeric;
```

## Go conversion pattern

Before:

```go
type Player struct {
    Level int `json:"level"`
    TotalXP int64 `json:"totalXp"`
}
```

After:

```go
type Player struct {
    Level GameInt `json:"level"`
    TotalXP GameInt `json:"totalXp"`
}
```

`GameInt` marshals to JSON as a string, so the frontend stays precise.

## Rollout order

1. Add the DB domains, Go type, and JS helpers.
2. Convert low-risk display-only fields to `GameInt` first.
3. Convert XP/currency/resource calculations.
4. Convert level curves and combat formulas.
5. Convert Playdeck HP/damage and future large stat systems.
6. Remove unsafe `Number(value).toLocaleString()` calls from growth-value displays.

This keeps the game playable while replacing the number engine one room at a time instead of yanking the entire house into a math portal.
