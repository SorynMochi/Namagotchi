alter table players
add column if not exists xp_into_level bigint not null default 0;

alter table players
add column if not exists nibbles bigint not null default 0;

alter table players
add column if not exists namicoin bigint not null default 0;

create table if not exists player_tick_state (
    player_id bigint primary key references players(id) on delete cascade,
    playdeck_enabled boolean not null default true,
    playdeck_zone_id integer not null default 1,
    playdeck_streak bigint not null default 0,
    playdeck_timeout_ticks integer not null default 0,
    active_gathering_task text not null default 'streaming',
    gathering_remainder double precision not null default 0,
    last_tick_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint player_tick_gathering_task_allowed check (
        active_gathering_task in (
            'streaming',
            'doom_scrolling',
            'cleaning',
            'exercising',
            'shopping',
            'designing'
        )
    )
);

insert into player_tick_state (player_id)
select id
from players
on conflict (player_id) do nothing;