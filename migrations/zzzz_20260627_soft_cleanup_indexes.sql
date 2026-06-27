create index if not exists idx_player_inventory_items_player_container
on player_inventory_items (player_id, container);

create index if not exists idx_player_inventory_items_player_equipped_slot
on player_inventory_items (player_id, equipped_slot)
where equipped_slot is not null;

create index if not exists idx_player_activity_skills_player_activity
on player_activity_skills (player_id, activity_key);

create index if not exists idx_nami_messages_player_created_at
on nami_messages (player_id, created_at desc, id desc);

create index if not exists idx_playdeck_combat_log_player_created_at
on playdeck_combat_log (player_id, created_at desc, id desc);

create index if not exists idx_companion_care_actions_player_status_queue
on companion_care_actions (player_id, status, queue_position);

create index if not exists idx_auth_sessions_session_hash_expires_at
on auth_sessions (session_hash, expires_at);

create index if not exists idx_security_event_logs_created_at
on security_event_logs (created_at desc, id desc);

create index if not exists idx_dev_audit_logs_created_at
on dev_audit_logs (created_at desc, id desc);