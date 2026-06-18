create unique index if not exists auth_accounts_email_lower_key
on auth_accounts (lower(email))
where email <> '';
