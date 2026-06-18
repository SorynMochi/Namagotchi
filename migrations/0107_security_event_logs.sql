create table if not exists security_event_logs (
  id bigserial primary key,
  account_id bigint references auth_accounts(id) on delete set null,
  event_type text not null,
  reason text not null,
  method text not null,
  path text not null,
  status_code integer not null default 0,
  remote_addr text not null default '',
  user_agent text not null default '',
  created_at timestamptz not null default now()
);

create index if not exists security_event_logs_created_at_idx
on security_event_logs (created_at desc);

create index if not exists security_event_logs_account_id_idx
on security_event_logs (account_id);

create index if not exists security_event_logs_event_type_idx
on security_event_logs (event_type);