create unique index if not exists players_account_id_key
on players(account_id)
where account_id is not null;
