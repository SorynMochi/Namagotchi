alter table players
  add column if not exists created_at timestamptz;

update players
set created_at = now()
where created_at is null;

alter table players
  alter column created_at set default now();

alter table players
  alter column created_at set not null;

alter table players
  add column if not exists online_seconds bigint;

update players
set online_seconds = 0
where online_seconds is null
  or online_seconds < 0;

alter table players
  alter column online_seconds set default 0;

alter table players
  alter column online_seconds set not null;

alter table players
  add column if not exists last_seen_at timestamptz;
