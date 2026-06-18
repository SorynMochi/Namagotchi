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
