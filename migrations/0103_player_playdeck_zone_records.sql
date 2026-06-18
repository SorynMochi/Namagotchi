create table if not exists player_playdeck_zone_records (
  player_id bigint not null references players(id) on delete cascade,
  zone_id integer not null,
  max_streak bigint not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  primary key (player_id, zone_id)
);