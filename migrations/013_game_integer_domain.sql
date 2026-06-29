do $$
begin
    if not exists (
        select 1
        from pg_type
        where typname = 'game_integer'
    ) then
        create domain game_integer as numeric(120, 0)
        check (value >= 0);
    end if;

    if not exists (
        select 1
        from pg_type
        where typname = 'signed_game_integer'
    ) then
        create domain signed_game_integer as numeric(120, 0);
    end if;
end $$;

comment on domain game_integer is 'Exact non-negative idle-game integer. Use for XP, levels, currencies, resources, HP, damage totals, streaks, counters, and other values that may grow beyond int64.';
comment on domain signed_game_integer is 'Exact signed idle-game integer. Use only for deltas or values that may intentionally go below zero.';
