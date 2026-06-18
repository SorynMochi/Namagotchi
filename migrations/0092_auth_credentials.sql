create table if not exists auth_credentials (
  account_id bigint primary key references auth_accounts(id) on delete cascade,
  email_normalized text not null unique,
  password_hash text not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
