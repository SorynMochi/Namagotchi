alter table companion_states
add column if not exists satiety_decay_remainder double precision not null default 0,
add column if not exists connection_decay_remainder double precision not null default 0,
add column if not exists energy_decay_remainder double precision not null default 0,
add column if not exists comfort_decay_remainder double precision not null default 0,
add column if not exists playfulness_decay_remainder double precision not null default 0,
add column if not exists inspiration_decay_remainder double precision not null default 0,
add column if not exists cleanliness_decay_remainder double precision not null default 0,
add column if not exists sleep_energy_recovery_remainder double precision not null default 0;