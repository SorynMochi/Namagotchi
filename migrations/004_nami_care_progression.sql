alter table companion_states
add column if not exists level integer not null default 1;

alter table companion_states
add column if not exists total_xp bigint not null default 0;

alter table companion_states
add column if not exists xp_into_level bigint not null default 0;

alter table companion_states
add column if not exists last_xp_gained bigint not null default 0;

alter table companion_states
add column if not exists last_action text not null default 'Created';

alter table companion_states
add column if not exists sleep_started_at timestamptz;

alter table companion_states
add column if not exists energy_at_sleep_start integer;