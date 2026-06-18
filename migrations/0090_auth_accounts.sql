create table if not exists auth_accounts (
  id bigserial primary key,
  display_name text not null,
  email text not null default '',
  avatar_url text not null default '',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  last_login_at timestamptz
);
