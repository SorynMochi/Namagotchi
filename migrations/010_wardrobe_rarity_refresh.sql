alter table item_definitions
drop constraint if exists item_rarity_allowed;

update item_definitions
set rarity = case lower(rarity)
    when 'common' then 'basic'
    when 'uncommon' then 'cute'
    when 'rare' then 'trendy'
    when 'epic' then 'iconic'
    when 'legendary' then 'devastating'
    else 'basic'
end;

update item_definitions
set rarity = case item_key
    when 'tutorial_sticker_choker' then 'cute'
    when 'pixel_petal_pin' then 'chic'
    when 'snack_cache_satchel' then 'cute'
    when 'debug_boots' then 'trendy'
    when 'cache_thread_spool' then 'basic'
    else rarity
end;

alter table item_definitions
add constraint item_rarity_allowed check (
    rarity in (
        'basic',
        'cute',
        'chic',
        'trendy',
        'glam',
        'iconic',
        'devastating'
    )
);