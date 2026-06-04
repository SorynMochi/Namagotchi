create table if not exists nami_messages (
	id bigserial primary key,
	player_id bigint not null references players(id) on delete cascade,
	trigger_key text not null,
	mood_key text not null default '',
	need_key text not null default '',
	severity text not null default 'info',
	message text not null,
	metadata_json jsonb not null default '{}'::jsonb,
	created_at timestamptz not null default now(),
	seen_at timestamptz
);

create index if not exists idx_nami_messages_player_created
on nami_messages (player_id, created_at desc, id desc);

create index if not exists idx_nami_messages_player_trigger_created
on nami_messages (player_id, trigger_key, created_at desc);

create table if not exists player_nami_message_state (
	player_id bigint primary key references players(id) on delete cascade,
	last_online_message_at timestamptz,
	next_random_message_at timestamptz,
	last_low_stat_message_at jsonb not null default '{}'::jsonb,
	updated_at timestamptz not null default now()
);