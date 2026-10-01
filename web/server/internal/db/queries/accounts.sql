-- accounts

-- name: CreateAccount :one
insert into accounts (user_id, name, type, currency, starting_balance)
values ($1, $2, $3, $4, $5)
returning id, user_id, name, type, currency, starting_balance, created_at, updated_at;

-- name: GetAccount :one
select id, user_id, name, type, currency, starting_balance, created_at, updated_at
from accounts where id = $1 and user_id = $2;

-- name: ListAccountsByUser :many
select id, user_id, name, type, currency, starting_balance, created_at, updated_at
from accounts where user_id = $1 order by created_at asc;

-- name: UpdateAccount :one
update accounts set
  name = coalesce(sqlc.narg(name), name),
  type = coalesce(sqlc.narg(type), type),
  currency = coalesce(sqlc.narg(currency), currency),
  starting_balance = coalesce(sqlc.narg(starting_balance), starting_balance)
where id = sqlc.arg(id) and user_id = sqlc.arg(user_id)
returning id, user_id, name, type, currency, starting_balance, created_at, updated_at;

-- name: DeleteAccount :exec
delete from accounts where id = $1 and user_id = $2;
