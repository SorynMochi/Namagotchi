create table if not exists player_activity_skills (
    player_id bigint not null references players(id) on delete cascade,
    activity_key text not null,
    level integer not null default 1,
    total_xp bigint not null default 0,
    xp_into_level bigint not null default 0,
    updated_at timestamptz not null default now(),

    primary key (player_id, activity_key),

    constraint player_activity_skill_key_allowed check (
        activity_key in (
            'streaming',
            'doom_scrolling',
            'cleaning',
            'exercising',
            'shopping',
            'designing'
        )
    ),

    constraint player_activity_skill_level_positive check (level >= 1),
    constraint player_activity_skill_xp_nonnegative check (total_xp >= 0 and xp_into_level >= 0)
);

insert into player_activity_skills (player_id, activity_key)
select p.id, activity.activity_key
from players p
cross join (
    values
        ('streaming'),
        ('doom_scrolling'),
        ('cleaning'),
        ('exercising'),
        ('shopping'),
        ('designing')
) as activity(activity_key)
on conflict (player_id, activity_key) do nothing;