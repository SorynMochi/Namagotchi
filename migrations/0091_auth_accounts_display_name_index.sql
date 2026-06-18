create unique index if not exists auth_accounts_display_name_lower_key
on auth_accounts (lower(display_name));
