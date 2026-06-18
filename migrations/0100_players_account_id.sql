alter table players
add column if not exists account_id bigint references auth_accounts(id) on delete set null;
