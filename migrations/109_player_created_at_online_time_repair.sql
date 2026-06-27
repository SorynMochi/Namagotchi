update players p
set created_at = coalesce(
  (
    select a.created_at
    from auth_accounts a
    where lower(a.display_name) = lower(p.display_name)
      and a.created_at >= '2020-01-01'::timestamptz
      and a.created_at <= now() + interval '1 day'
    order by a.created_at
    limit 1
  ),
  now()
)
where p.created_at < '2020-01-01'::timestamptz
   or p.created_at > now() + interval '1 day';

update players
set online_seconds = 0
where online_seconds is null
   or online_seconds < 0;
