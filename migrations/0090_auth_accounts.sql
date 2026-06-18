create table if not exists auth_accounts (
  id bigserial primary key,
  display_name text not null,
  email text not null default '',
  avatar_url text not null default '',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  last_login_at timestamptz
);

create unique index if not exists auth_accounts_display_name_lower_key
on auth_accounts (lower(display_name));

create table if not exists auth_credentials (
  account_id bigint primary key references auth_accounts(id) on delete cascade,
  email_normalized text not null unique,
  password_hash text not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists auth_identities (
  id bigserial primary key,
  account_id bigint not null references auth_accounts(id) on delete cascade,
  provider text not null,
  provider_user_id text not null,
  email text not null default '',
  display_name text not null default '',
  avatar_url text not null default '',
  created_at timestamptz not null default now(),
  last_login_at timestamptz,
  unique(provider, provider_user_id)
);

create index if not exists auth_identities_account_id_idx
on auth_identities(account_id);

create table if not exists auth_sessions (
  session_hash text primary key,
  account_id bigint not null references auth_accounts(id) on delete cascade,
  created_at timestamptz not null default now(),
  last_seen_at timestamptz not null default now(),
  expires_at timestamptz not null
);

create index if not exists auth_sessions_account_id_idx
on auth_sessions(account_id);

create index if not exists auth_sessions_expires_at_idx
on auth_sessions(expires_at);

create table if not exists auth_oauth_states (
  state_hash text primary key,
  provider text not null,
  redirect_path text not null default '/',
  created_at timestamptz not null default now(),
  expires_at timestamptz not null
);

create index if not exists auth_oauth_states_expires_at_idx
on auth_oauth_states(expires_at);
