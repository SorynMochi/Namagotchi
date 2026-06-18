create table if not exists dev_audit_logs (
  id bigserial primary key,
  account_id bigint references auth_accounts(id) on delete set null,
  command text not null,
  method text not null,
  path text not null,
  status_code integer not null default 0,
  remote_addr text not null default '',
  user_agent text not null default '',
  created_at timestamptz not null default now()
);

create index if not exists dev_audit_logs_created_at_idx
on dev_audit_logs (created_at desc);

create index if not exists dev_audit_logs_account_id_idx
on dev_audit_logs (account_id);