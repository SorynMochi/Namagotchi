begin;

update player_inventory_items
set equipped_slot = null
where equipped_slot = 'bag';

update item_definitions
set equipment_slot = 'accessory'
where equipment_slot = 'bag';

delete from playdeck_equipment_slots
where slot_key = 'bag';

alter table playdeck_equipment_slots
drop constraint if exists equipment_slot_accepts_allowed;

alter table playdeck_equipment_slots
add constraint equipment_slot_accepts_allowed check (
    accepts_slot in (
        'top',
        'bottom',
        'dress',
        'footwear',
        'outerwear',
        'necklace',
        'accessory'
    )
);

alter table item_definitions
drop constraint if exists item_equipment_slot_allowed;

alter table item_definitions
add constraint item_equipment_slot_allowed check (
    equipment_slot in (
        'none',
        'top',
        'bottom',
        'dress',
        'footwear',
        'outerwear',
        'necklace',
        'accessory'
    )
);

commit;