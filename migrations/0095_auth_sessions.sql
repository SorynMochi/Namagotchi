create table if not exists auth_sessions (
  session_hash text primary key,
  account_id bigint not null references auth_accounts(id) on delete cascade,
  created_at timestamptz not null default now(),
  last_seen_at timestamptz not null default now(),
  expires_at timestamptz not null
);
