create table if not exists auth_oauth_states (
  state_hash text primary key,
  provider text not null,
  redirect_path text not null default '/',
  created_at timestamptz not null default now(),
  expires_at timestamptz not null
);
