-- setups

-- name: CreateSetup :one
insert into setups (user_id, name, rules, invalidation)
values ($1, $2, $3, $4)
returning id, user_id, name, rules, invalidation, created_at, updated_at;

-- name: GetSetup :one
select id, user_id, name, rules, invalidation, created_at, updated_at
from setups where id = $1 and user_id = $2;

-- name: ListSetupsByUser :many
select id, user_id, name, rules, invalidation, created_at, updated_at
from setups where user_id = $1 order by name asc;

-- name: UpdateSetup :one
update setups set
  name = coalesce(sqlc.narg(name), name),
  rules = coalesce(sqlc.narg(rules), rules),
  invalidation = coalesce(sqlc.narg(invalidation), invalidation)
where id = sqlc.arg(id) and user_id = sqlc.arg(user_id)
returning id, user_id, name, rules, invalidation, created_at, updated_at;

-- name: DeleteSetup :exec
delete from setups where id = $1 and user_id = $2;
