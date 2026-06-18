begin;

alter table player_inventory_items
add column if not exists tailoring_current integer not null default 0;

alter table player_inventory_items
add column if not exists tailoring_max integer not null default 0;

alter table player_inventory_items
drop constraint if exists player_inventory_tailoring_nonnegative;

alter table player_inventory_items
add constraint player_inventory_tailoring_nonnegative check (
    tailoring_current >= 0
    and tailoring_max >= 0
    and tailoring_current <= tailoring_max
);

create table if not exists wardrobe_stat_definitions (
    stat_key text primary key,
    display_name text not null,
    value_kind text not null,
    stat_family text not null,
    applies_to text not null,
    allowed_as_implicit boolean not null default false,
    bypasses_beauty boolean not null default false,
    bypasses_glamor boolean not null default false,
    tooltip text not null default '',
    sort_order integer not null unique,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    constraint wardrobe_stat_value_kind_allowed check (
        value_kind in ('percent', 'flat')
    ),
    constraint wardrobe_stat_family_allowed check (
        stat_family in ('prefix', 'suffix')
    ),
    constraint wardrobe_stat_applies_to_allowed check (
        applies_to in ('playdeck', 'work', 'global')
    )
);

create table if not exists player_inventory_item_stat_lines (
    id bigserial primary key,
    player_inventory_item_id bigint not null references player_inventory_items(id) on delete cascade,
    stat_source text not null,
    affix_key text not null default '',
    stat_key text not null references wardrobe_stat_definitions(stat_key) on delete restrict,
    value numeric(18, 4) not null,
    sort_order integer not null default 0,
    created_at timestamptz not null default now(),

    constraint wardrobe_item_stat_source_allowed check (
        stat_source in ('implicit', 'prefix', 'suffix', 'devastating_implicit')
    )
);

create index if not exists player_inventory_item_stat_lines_item_idx
on player_inventory_item_stat_lines(player_inventory_item_id, sort_order, id);

create index if not exists player_inventory_item_stat_lines_stat_idx
on player_inventory_item_stat_lines(stat_key);

insert into wardrobe_stat_definitions (
    stat_key,
    display_name,
    value_kind,
    stat_family,
    applies_to,
    allowed_as_implicit,
    bypasses_beauty,
    bypasses_glamor,
    tooltip,
    sort_order
)
values
    (
        'playdeck_xp_percent',
        'Playdeck XP',
        'percent',
        'prefix',
        'playdeck',
        false,
        false,
        false,
        'Increases XP gained from Playdeck combat.',
        100
    ),
    (
        'work_xp_percent',
        'Work XP',
        'percent',
        'prefix',
        'work',
        false,
        false,
        false,
        'Increases XP gained from Work activities.',
        110
    ),
    (
        'global_xp_percent',
        'Global XP',
        'percent',
        'prefix',
        'global',
        false,
        false,
        false,
        'Increases both Playdeck XP and Work XP gains.',
        120
    ),
    (
        'work_resources_percent',
        'Work Yield',
        'percent',
        'prefix',
        'work',
        false,
        false,
        false,
        'Increases resources gained from Work activities.',
        130
    ),
    (
        'drop_rate_percent',
        'Drop Rate',
        'percent',
        'prefix',
        'playdeck',
        false,
        false,
        false,
        'Increases how often drops happen, including items, Nibbles, and ingredient drops. Over 100% final chance can create extra drop rolls.',
        140
    ),
    (
        'credit_rate_percent',
        'Credits',
        'percent',
        'prefix',
        'playdeck',
        false,
        false,
        false,
        'Increases Credits gained.',
        150
    ),
    (
        'ingredient_quality_percent',
        'Ingredients',
        'percent',
        'prefix',
        'playdeck',
        false,
        false,
        false,
        'Improves the chance for higher-quality ingredients when ingredients drop.',
        160
    ),
    (
        'max_health_percent',
        'Sparkles',
        'percent',
        'suffix',
        'playdeck',
        true,
        false,
        false,
        'Increases maximum health. More Sparkles means Nami-Chan can take more damage before losing.',
        200
    ),
    (
        'attack_percent',
        'Chaos',
        'percent',
        'suffix',
        'playdeck',
        true,
        false,
        false,
        'Increases attack damage dealt in Playdeck combat.',
        210
    ),
    (
        'attack_speed_percent',
        'Zoomies',
        'percent',
        'suffix',
        'playdeck',
        true,
        false,
        false,
        'Increases attacks per combat round. Zoomies affects Playdeck attacks only, not game tick speed.',
        220
    ),
    (
        'beauty',
        'Beauty',
        'flat',
        'suffix',
        'playdeck',
        true,
        false,
        false,
        'Reduces incoming damage before Glamor Shield is hit. Higher Beauty gives less extra reduction the more you already have.',
        230
    ),
    (
        'glamor',
        'Glamor Shield',
        'flat',
        'suffix',
        'playdeck',
        true,
        false,
        false,
        'Adds a shield pool at the start of each fight. Damage hits Glamor Shield before health. Refreshes between fights, not between rounds.',
        240
    ),
    (
        'crit_rate_percent',
        'Bonk Chance',
        'percent',
        'suffix',
        'playdeck',
        true,
        false,
        false,
        'Increases the chance for each attack to crit. Over 100% can create multiple crits.',
        250
    ),
    (
        'crit_damage_percent',
        'Bonk Damage',
        'percent',
        'suffix',
        'playdeck',
        true,
        false,
        false,
        'Increases crit damage. Multiple crits stack additively with reduced value per extra crit.',
        260
    ),
    (
        'charm',
        'Charm',
        'flat',
        'suffix',
        'playdeck',
        true,
        true,
        false,
        'Applies stacking damage over time when an attack hits. Charm cannot crit and bypasses Beauty, but it hits Glamor Shield before health.',
        270
    ),
    (
        'humor',
        'Sass',
        'flat',
        'suffix',
        'playdeck',
        true,
        true,
        false,
        'Reflects flat damage back when hit. Sass cannot crit and bypasses Beauty, but it hits Glamor Shield before health.',
        280
    ),
    (
        'targeting_percent',
        'Focus',
        'percent',
        'suffix',
        'playdeck',
        true,
        false,
        false,
        'Increases hit reliability in Playdeck combat. Focus works against the enemy’s Dodge.',
        290
    ),
    (
        'dodge_percent',
        'Dodge',
        'percent',
        'suffix',
        'playdeck',
        true,
        false,
        false,
        'Increases the chance to evade incoming attacks. Final evade chance is capped at 70%.',
        300
    ),
    (
        'recovery',
        'Patch-Up',
        'flat',
        'suffix',
        'playdeck',
        true,
        false,
        false,
        'Restores health at the beginning of each combat turn before Charm damage happens.',
        310
    )
on conflict (stat_key) do update
set
    display_name = excluded.display_name,
    value_kind = excluded.value_kind,
    stat_family = excluded.stat_family,
    applies_to = excluded.applies_to,
    allowed_as_implicit = excluded.allowed_as_implicit,
    bypasses_beauty = excluded.bypasses_beauty,
    bypasses_glamor = excluded.bypasses_glamor,
    tooltip = excluded.tooltip,
    sort_order = excluded.sort_order,
    updated_at = now();

commit;