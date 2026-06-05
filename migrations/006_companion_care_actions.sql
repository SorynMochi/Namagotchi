create table if not exists companion_care_actions (
	id bigserial primary key,
	player_id bigint not null references players(id) on delete cascade,
	action_key text not null,
	action_name text not null,
	status text not null default 'queued',
	queue_position integer,
	duration_seconds integer not null,
	started_at timestamptz,
	completes_at timestamptz,
	completed_at timestamptz,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now(),

	constraint companion_care_action_status_allowed check (
		status in ('queued', 'active', 'completed', 'cancelled')
	),

	constraint companion_care_action_key_allowed check (
		action_key in (
			'meal',
			'snack',
			'drink',
			'cuddle',
			'play',
			'write_together',
			'read_together',
			'boop',
			'nap',
			'bath',
			'freshen_up',
			'put_to_bed',
			'wake_up'
		)
	),

	constraint companion_care_action_queue_position_valid check (
		queue_position is null or queue_position between 1 and 3
	),

	constraint companion_care_action_duration_positive check (
		duration_seconds > 0
	)
);

create unique index if not exists companion_care_one_active_idx
on companion_care_actions (player_id)
where status = 'active';

create unique index if not exists companion_care_one_queued_action_idx
on companion_care_actions (player_id, action_key)
where status = 'queued';

create unique index if not exists companion_care_queue_position_idx
on companion_care_actions (player_id, queue_position)
where status = 'queued';

create index if not exists companion_care_due_active_idx
on companion_care_actions (completes_at)
where status = 'active';

create index if not exists companion_care_player_status_idx
on companion_care_actions (player_id, status, queue_position, created_at);

create index if not exists companion_care_player_history_idx
on companion_care_actions (player_id, created_at desc, id desc);