with ranked_messages as (
  select
    id,
    row_number() over (
      partition by player_id
      order by created_at desc, id desc
    ) as message_rank
  from nami_messages
)
delete from nami_messages message
using ranked_messages ranked
where message.id = ranked.id
  and ranked.message_rank > 50;