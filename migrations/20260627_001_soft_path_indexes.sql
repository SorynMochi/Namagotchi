create index if not exists idx_player_inventory_items_player_container
on player_inventory_items (player_id, container);

create index if not exists idx_player_inventory_items_player_equipped_slot
on player_inventory_items (player_id, equipped_slot)
where equipped_slot is not null;

create index if not exists idx_nami_messages_player_recent
on nami_messages (player_id, created_at desc, id desc);

create index if not exists idx_companion_care_actions_player_status_queue
on companion_care_actions (
player_id,
status,
queue_position,
created_at,
id
)
where status in ('active', 'queued');

create index if not exists idx_companion_care_actions_due_active
on companion_care_actions (
player_id,
completes_at,
id
)
where status = 'active';

create index if not exists idx_playdeck_combat_log_player_recent
on playdeck_combat_log (player_id, created_at desc, id desc);