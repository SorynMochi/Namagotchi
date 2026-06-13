create table if not exists playdeck_zones (
    id integer primary key,
    zone_key text not null unique,
    name text not null,
    min_level integer not null default 1,
    softcap_level integer not null default 50,
    description text not null default '',
    is_unlocked_default boolean not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    constraint playdeck_zone_level_positive check (min_level >= 1 and softcap_level >= min_level)
);

create table if not exists playdeck_enemies (
    id bigserial primary key,
    zone_id integer not null references playdeck_zones(id) on delete cascade,
    enemy_key text not null,
    name text not null,
    min_level integer not null default 1,
    max_level integer not null default 50,
    base_hp integer not null default 25,
    hp_per_level double precision not null default 4,
    base_attack integer not null default 4,
    attack_per_level double precision not null default 1,
    base_xp bigint not null default 20,
    base_credits_cents bigint not null default 2500,
    base_nibbles bigint not null default 1,
    weight integer not null default 100,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    unique (zone_id, enemy_key),

    constraint playdeck_enemy_level_range check (min_level >= 1 and max_level >= min_level),
    constraint playdeck_enemy_stats_positive check (
        base_hp > 0
        and hp_per_level >= 0
        and base_attack > 0
        and attack_per_level >= 0
        and base_xp >= 0
        and base_credits_cents >= 0
        and base_nibbles >= 0
        and weight > 0
    )
);

create table if not exists item_definitions (
    id bigserial primary key,
    item_key text not null unique,
    name text not null,
    item_type text not null,
    rarity text not null default 'common',
    equipment_slot text not null default 'none',
    stackable boolean not null default false,
    max_stack integer not null default 1,
    power_level integer not null default 1,
    attack_bonus integer not null default 0,
    defense_bonus integer not null default 0,
    max_hp_bonus integer not null default 0,
    value_cents bigint not null default 0,
    description text not null default '',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    constraint item_type_allowed check (
        item_type in ('gear', 'cosmetic', 'ingredient', 'pattern', 'consumable', 'material')
    ),
    constraint item_rarity_allowed check (
        rarity in ('common', 'uncommon', 'rare', 'epic', 'legendary')
    ),
    constraint item_equipment_slot_allowed check (
        equipment_slot in (
            'none',
            'top',
            'bottom',
            'dress',
            'footwear',
            'outerwear',
            'necklace',
            'bag',
            'accessory'
        )
    ),
    constraint item_stack_positive check (max_stack >= 1),
    constraint item_power_positive check (power_level >= 1)
);

create table if not exists playdeck_equipment_slots (
    slot_key text primary key,
    display_name text not null,
    accepts_slot text not null,
    sort_order integer not null unique,
    created_at timestamptz not null default now(),

    constraint equipment_slot_accepts_allowed check (
        accepts_slot in (
            'top',
            'bottom',
            'dress',
            'footwear',
            'outerwear',
            'necklace',
            'bag',
            'accessory'
        )
    )
);

create table if not exists player_inventory_items (
    id bigserial primary key,
    player_id bigint not null references players(id) on delete cascade,
    item_definition_id bigint not null references item_definitions(id) on delete restrict,
    container text not null default 'wardrobe',
    quantity integer not null default 1,
    equipped_slot text references playdeck_equipment_slots(slot_key) on delete set null,
    metadata jsonb not null default '{}'::jsonb,
    acquired_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    constraint player_inventory_container_allowed check (
        container in ('wardrobe', 'materials', 'consumables')
    ),
    constraint player_inventory_quantity_positive check (quantity >= 1)
);

create unique index if not exists player_inventory_equipped_slot_unique_idx
on player_inventory_items(player_id, equipped_slot)
where equipped_slot is not null;

create index if not exists player_inventory_player_container_idx
on player_inventory_items(player_id, container, acquired_at desc);

create table if not exists playdeck_combat_state (
    player_id bigint primary key references players(id) on delete cascade,
    zone_id integer not null references playdeck_zones(id) on delete restrict,
    enemy_id bigint not null references playdeck_enemies(id) on delete restrict,
    enemy_level integer not null default 1,
    enemy_hp integer not null default 1,
    player_hp integer not null default 100,
    wins bigint not null default 0,
    losses bigint not null default 0,
    last_outcome text not null default 'ready',
    updated_at timestamptz not null default now(),

    constraint playdeck_combat_enemy_level_positive check (enemy_level >= 1),
    constraint playdeck_combat_hp_nonnegative check (enemy_hp >= 0 and player_hp >= 0),
    constraint playdeck_combat_counts_nonnegative check (wins >= 0 and losses >= 0),
    constraint playdeck_combat_outcome_allowed check (
        last_outcome in ('ready', 'fighting', 'win', 'loss')
    )
);

create table if not exists playdeck_combat_log (
    id bigserial primary key,
    player_id bigint not null references players(id) on delete cascade,
    zone_id integer not null references playdeck_zones(id) on delete restrict,
    enemy_name text not null,
    enemy_level integer not null,
    outcome text not null,
    player_damage integer not null default 0,
    enemy_damage integer not null default 0,
    xp_gained bigint not null default 0,
    credits_cents_gained bigint not null default 0,
    nibbles_gained bigint not null default 0,
    item_definition_id bigint references item_definitions(id) on delete set null,
    item_quantity integer not null default 0,
    created_at timestamptz not null default now(),

    constraint playdeck_log_outcome_allowed check (
        outcome in ('hit', 'win', 'loss')
    )
);

create index if not exists playdeck_combat_log_player_created_idx
on playdeck_combat_log(player_id, created_at desc, id desc);

create table if not exists playdeck_drop_tables (
    id bigserial primary key,
    zone_id integer not null references playdeck_zones(id) on delete cascade,
    enemy_id bigint references playdeck_enemies(id) on delete cascade,
    item_definition_id bigint not null references item_definitions(id) on delete cascade,
    drop_chance_basis_points integer not null,
    min_quantity integer not null default 1,
    max_quantity integer not null default 1,
    created_at timestamptz not null default now(),

    constraint playdeck_drop_chance_range check (
        drop_chance_basis_points between 1 and 10000
    ),
    constraint playdeck_drop_quantity_range check (
        min_quantity >= 1 and max_quantity >= min_quantity
    )
);

create index if not exists playdeck_drop_tables_zone_enemy_idx
on playdeck_drop_tables(zone_id, enemy_id);

insert into playdeck_zones (
    id,
    zone_key,
    name,
    min_level,
    softcap_level,
    description,
    is_unlocked_default
)
values
    (1, 'starter_deck', 'Starter Deck', 1, 50, 'Training bots, pop-up goblins, and suspiciously smug tutorials.', true),
    (2, 'cozy_lan_cafe', 'Cozy LAN Café', 50, 100, 'Friendly chaos, snack-powered rivals, and rogue leaderboard gremlins.', false),
    (3, 'neon_mall_net', 'Neon Mall Net', 250, 1500, 'Fashion bots, coupon phantoms, and glittery malware.', false)
on conflict (id) do update
set
    zone_key = excluded.zone_key,
    name = excluded.name,
    min_level = excluded.min_level,
    softcap_level = excluded.softcap_level,
    description = excluded.description,
    is_unlocked_default = excluded.is_unlocked_default,
    updated_at = now();

insert into playdeck_enemies (
    zone_id,
    enemy_key,
    name,
    min_level,
    max_level,
    base_hp,
    hp_per_level,
    base_attack,
    attack_per_level,
    base_xp,
    base_credits_cents,
    base_nibbles,
    weight
)
values
    (1, 'tutorial_popup_gremlin', 'Tutorial Pop-Up Gremlin', 1, 50, 22, 4.2, 4, 0.85, 24, 2600, 1, 45),
    (1, 'cozy_cache_sprite', 'Cozy Cache Sprite', 1, 50, 18, 3.7, 3, 0.75, 22, 2400, 1, 35),
    (1, 'idle_snack_daemon', 'Idle Snack Daemon', 1, 50, 28, 4.8, 5, 0.95, 30, 3000, 1, 20)
on conflict (zone_id, enemy_key) do update
set
    name = excluded.name,
    min_level = excluded.min_level,
    max_level = excluded.max_level,
    base_hp = excluded.base_hp,
    hp_per_level = excluded.hp_per_level,
    base_attack = excluded.base_attack,
    attack_per_level = excluded.attack_per_level,
    base_xp = excluded.base_xp,
    base_credits_cents = excluded.base_credits_cents,
    base_nibbles = excluded.base_nibbles,
    weight = excluded.weight,
    updated_at = now();

insert into playdeck_equipment_slots (
    slot_key,
    display_name,
    accepts_slot,
    sort_order
)
values
    ('top', 'Top', 'top', 10),
    ('bottom', 'Bottom', 'bottom', 20),
    ('dress', 'Dress / Outfit', 'dress', 30),
    ('footwear', 'Footwear', 'footwear', 40),
    ('outerwear', 'Outerwear', 'outerwear', 50),
    ('necklace', 'Necklace', 'necklace', 60),
    ('bag', 'Bag', 'bag', 70),
    ('accessory_1', 'Accessory 1', 'accessory', 80),
    ('accessory_2', 'Accessory 2', 'accessory', 90)
on conflict (slot_key) do update
set
    display_name = excluded.display_name,
    accepts_slot = excluded.accepts_slot,
    sort_order = excluded.sort_order;

insert into item_definitions (
    item_key,
    name,
    item_type,
    rarity,
    equipment_slot,
    stackable,
    max_stack,
    power_level,
    attack_bonus,
    defense_bonus,
    max_hp_bonus,
    value_cents,
    description
)
values
    ('tutorial_sticker_choker', 'Tutorial Sticker Choker', 'gear', 'common', 'necklace', false, 1, 1, 1, 0, 2, 1500, 'A glossy little starter necklace with suspiciously motivational sticker energy.'),
    ('pixel_petal_pin', 'Pixel Petal Pin', 'gear', 'common', 'accessory', false, 1, 1, 0, 1, 3, 1400, 'A tiny accessory pin that looks soft until the pixels bite back.'),
    ('snack_cache_satchel', 'Snack Cache Satchel', 'gear', 'common', 'bag', false, 1, 1, 0, 2, 4, 1700, 'A small bag with emergency snack logic and beginner-friendly padding.'),
    ('debug_boots', 'Debug Boots', 'gear', 'uncommon', 'footwear', false, 1, 3, 2, 1, 4, 3000, 'Boots designed for stepping over tiny bugs with maximum dignity.'),
    ('cache_thread_spool', 'Cache Thread Spool', 'material', 'common', 'none', true, 999, 1, 0, 0, 0, 250, 'A spool of strange digital thread. Useful later for crafting, tailoring, and tiny textile crimes.')
on conflict (item_key) do update
set
    name = excluded.name,
    item_type = excluded.item_type,
    rarity = excluded.rarity,
    equipment_slot = excluded.equipment_slot,
    stackable = excluded.stackable,
    max_stack = excluded.max_stack,
    power_level = excluded.power_level,
    attack_bonus = excluded.attack_bonus,
    defense_bonus = excluded.defense_bonus,
    max_hp_bonus = excluded.max_hp_bonus,
    value_cents = excluded.value_cents,
    description = excluded.description,
    updated_at = now();

insert into playdeck_drop_tables (
    zone_id,
    enemy_id,
    item_definition_id,
    drop_chance_basis_points,
    min_quantity,
    max_quantity
)
select
    z.id,
    null,
    i.id,
    drop_rule.drop_chance_basis_points,
    drop_rule.min_quantity,
    drop_rule.max_quantity
from (
    values
        ('cache_thread_spool', 3500, 1, 3),
        ('tutorial_sticker_choker', 400, 1, 1),
        ('pixel_petal_pin', 400, 1, 1),
        ('snack_cache_satchel', 300, 1, 1),
        ('debug_boots', 120, 1, 1)
) as drop_rule(item_key, drop_chance_basis_points, min_quantity, max_quantity)
join playdeck_zones z on z.zone_key = 'starter_deck'
join item_definitions i on i.item_key = drop_rule.item_key
where not exists (
    select 1
    from playdeck_drop_tables existing
    where existing.zone_id = z.id
        and existing.enemy_id is null
        and existing.item_definition_id = i.id
);

insert into playdeck_combat_state (
    player_id,
    zone_id,
    enemy_id,
    enemy_level,
    enemy_hp,
    player_hp
)
select
    p.id,
    z.id,
    e.id,
    greatest(1, least(p.level, e.max_level)),
    e.base_hp + round(e.hp_per_level * greatest(1, least(p.level, e.max_level)))::integer,
    100 + greatest(1, p.level) * 4
from players p
join playdeck_zones z on z.zone_key = 'starter_deck'
join playdeck_enemies e on e.zone_id = z.id and e.enemy_key = 'tutorial_popup_gremlin'
on conflict (player_id) do nothing;