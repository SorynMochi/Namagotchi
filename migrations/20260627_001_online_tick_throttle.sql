alter table players
add column if not exists online_last_seen_at timestamptz;