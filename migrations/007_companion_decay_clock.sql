alter table companion_states
add column if not exists last_decay_at timestamptz;

update companion_states
set last_decay_at = coalesce(last_decay_at, last_interaction_at, updated_at, now())
where last_decay_at is null;

alter table companion_states
alter column last_decay_at set default now();

alter table companion_states
alter column last_decay_at set not null;