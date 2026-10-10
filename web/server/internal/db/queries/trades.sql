-- trades

-- name: CreateTrade :one
insert into trades (
  user_id, account_id, setup_id, pair, direction, timeframe,
  opened_at, closed_at, entry, stop_loss, take_profit, exit_price,
  lot_size, commission, swap, risk_amount, pnl, r_multiple,
  followed_rules, emotion, notes, status, is_public
) values (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
  $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23
) returning *;

-- name: GetTrade :one
select * from trades where id = $1 and user_id = $2;

-- name: ListTrades :many
select * from trades
where user_id = sqlc.arg(user_id)
  and (sqlc.arg(pair)::text = '' or pair = sqlc.arg(pair)::text)
  and (sqlc.narg(setup_id)::uuid is null or setup_id = sqlc.narg(setup_id)::uuid)
  and (sqlc.arg(status)::text = '' or status = sqlc.arg(status)::text)
  and (sqlc.narg(account_id)::uuid is null or account_id = sqlc.narg(account_id)::uuid)
  and (sqlc.arg(direction)::text = '' or direction = sqlc.arg(direction)::text)
  and (sqlc.narg(opened_from)::timestamptz is null or opened_at >= sqlc.narg(opened_from)::timestamptz)
  and (sqlc.narg(opened_to)::timestamptz is null or opened_at <= sqlc.narg(opened_to)::timestamptz)
order by opened_at desc nulls last, created_at desc
limit sqlc.arg(page_limit) offset sqlc.arg(page_offset);

-- name: CountTrades :one
select count(*) from trades where user_id = $1;

-- name: CountOpenTrades :one
select count(*) from trades where user_id = $1 and status = 'open';

-- name: CountOpenTradesByAccount :one
select count(*) from trades where user_id = $1 and account_id = $2 and status = 'open';

-- name: CountAllTrades :one
select count(*) from trades;

-- name: UpdateTrade :one
update trades set
  account_id = coalesce(sqlc.narg(account_id), account_id),
  setup_id = coalesce(sqlc.narg(setup_id), setup_id),
  pair = coalesce(sqlc.narg(pair), pair),
  direction = coalesce(sqlc.narg(direction), direction),
  timeframe = coalesce(sqlc.narg(timeframe), timeframe),
  opened_at = coalesce(sqlc.narg(opened_at), opened_at),
  closed_at = coalesce(sqlc.narg(closed_at), closed_at),
  entry = coalesce(sqlc.narg(entry), entry),
  stop_loss = coalesce(sqlc.narg(stop_loss), stop_loss),
  take_profit = coalesce(sqlc.narg(take_profit), take_profit),
  exit_price = coalesce(sqlc.narg(exit_price), exit_price),
  lot_size = coalesce(sqlc.narg(lot_size), lot_size),
  commission = coalesce(sqlc.narg(commission), commission),
  swap = coalesce(sqlc.narg(swap), swap),
  risk_amount = coalesce(sqlc.narg(risk_amount), risk_amount),
  pnl = coalesce(sqlc.narg(pnl), pnl),
  r_multiple = coalesce(sqlc.narg(r_multiple), r_multiple),
  followed_rules = coalesce(sqlc.narg(followed_rules), followed_rules),
  emotion = coalesce(sqlc.narg(emotion), emotion),
  notes = coalesce(sqlc.narg(notes), notes),
  status = coalesce(sqlc.narg(status), status)
where id = sqlc.arg(id) and user_id = sqlc.arg(user_id)
returning *;

-- name: UpdateTradeVisibility :one
update trades set is_public = $3 where id = $1 and user_id = $2
returning *;

-- name: ReopenTrade :one
update trades set status = 'open', exit_price = null, closed_at = null, pnl = null, r_multiple = null
where id = $1 and user_id = $2
returning *;

-- name: HideJournal :one
update trades set hidden_by_admin = $2 where id = $1
returning *;

-- name: DeleteTrade :exec
delete from trades where id = $1 and user_id = $2;

-- name: ListPublicJournals :many
select t.*, p.display_name, s.name as setup_name
from trades t
join profiles p on p.id = t.user_id
left join setups s on s.id = t.setup_id
where t.is_public = true and t.hidden_by_admin = false
order by t.created_at desc
limit $1 offset $2;

-- name: GetPublicJournal :one
select t.*, p.display_name, s.name as setup_name
from trades t
join profiles p on p.id = t.user_id
left join setups s on s.id = t.setup_id
where t.id = $1 and t.is_public = true and t.hidden_by_admin = false;

-- name: ListAdminJournals :many
select t.*, p.display_name, s.name as setup_name
from trades t
join profiles p on p.id = t.user_id
left join setups s on s.id = t.setup_id
where t.is_public = true
order by t.created_at desc
limit $1 offset $2;

-- name: ExportTradesByUser :many
select * from trades where user_id = $1 order by opened_at asc nulls last;
