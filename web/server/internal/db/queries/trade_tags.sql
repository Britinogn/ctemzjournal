-- trade_tags

-- name: AddTradeTag :exec
insert into trade_tags (trade_id, tag_id) values ($1, $2)
on conflict do nothing;

-- name: RemoveTradeTag :exec
delete from trade_tags where trade_id = $1 and tag_id = $2;

-- name: ListTagsByTrade :many
select g.id, g.user_id, g.name, g.kind, g.created_at
from tags g join trade_tags tt on tt.tag_id = g.id
where tt.trade_id = $1 order by g.name asc;

-- name: ClearTradeTags :exec
delete from trade_tags where trade_id = $1;

-- name: ListTradeIDsByTag :many
select trade_id from trade_tags where tag_id = $1;
