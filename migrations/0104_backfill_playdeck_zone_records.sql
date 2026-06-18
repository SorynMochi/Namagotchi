insert into player_playdeck_zone_records (player_id, zone_id, max_streak)
select
  player_id,
  coalesce(nullif(playdeck_zone_id, 0), 1),
  greatest(coalesce(playdeck_streak, 0), 0)
from player_tick_state
where coalesce(playdeck_streak, 0) > 0
on conflict (player_id, zone_id) do update
set max_streak = greatest(player_playdeck_zone_records.max_streak, excluded.max_streak),
    updated_at = now();