create table if not exists players (
    id bigserial primary key,
    display_name text not null unique,
    level integer not null default 1,
    total_xp bigint not null default 0,
    currency_cents bigint not null default 0,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists companion_states (
    player_id bigint primary key references players(id) on delete cascade,
    companion_name text not null default 'Nami-chan',
    mood_score numeric(5, 2) not null default 80.00,
    satiety integer not null default 85,
    connection integer not null default 95,
    energy integer not null default 75,
    comfort integer not null default 90,
    playfulness integer not null default 80,
    inspiration integer not null default 80,
    cleanliness integer not null default 80,
    status text not null default 'awake',
    last_interaction_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint companion_satiety_range check (satiety between 0 and 100),
    constraint companion_connection_range check (connection between 0 and 100),
    constraint companion_energy_range check (energy between 0 and 100),
    constraint companion_comfort_range check (comfort between 0 and 100),
    constraint companion_playfulness_range check (playfulness between 0 and 100),
    constraint companion_inspiration_range check (inspiration between 0 and 100),
    constraint companion_cleanliness_range check (cleanliness between 0 and 100),
    constraint companion_status_allowed check (status in ('awake', 'sleeping'))
);

create table if not exists player_resources (
    player_id bigint primary key references players(id) on delete cascade,
    fans bigint not null default 0,
    memes bigint not null default 0,
    lost_items bigint not null default 0,
    confidence bigint not null default 0,
    receipts bigint not null default 0,
    patterns bigint not null default 0,
    glitch_drops bigint not null default 0,
    updated_at timestamptz not null default now()
);

create table if not exists activity_log (
    id bigserial primary key,
    player_id bigint references players(id) on delete cascade,
    event_type text not null,
    message text not null,
    created_at timestamptz not null default now()
);

create index if not exists activity_log_player_created_idx
on activity_log(player_id, created_at desc);