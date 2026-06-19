create unique index if not exists auth_accounts_display_name_lower_unique_idx
on auth_accounts (lower(display_name));

create unique index if not exists players_claimed_display_name_lower_unique_idx
on players (lower(display_name))
where account_id is not null;