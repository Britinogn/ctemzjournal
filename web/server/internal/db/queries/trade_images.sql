-- trade_images (only Cloudinary public_id is stored)

-- name: CreateTradeImage :one
insert into trade_images (trade_id, public_id, kind, position)
values ($1, $2, $3, $4)
returning id, trade_id, public_id, kind, position, created_at;

-- name: ListImagesByTrade :many
select id, trade_id, public_id, kind, position, created_at
from trade_images where trade_id = $1 order by position asc;

-- name: CountImagesByTrade :one
select count(*) from trade_images where trade_id = $1;

-- name: GetTradeImage :one
select id, trade_id, public_id, kind, position, created_at
from trade_images where id = $1;

-- name: DeleteTradeImage :exec
delete from trade_images where id = $1;
