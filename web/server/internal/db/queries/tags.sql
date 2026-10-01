-- tags

-- name: CreateTag :one
insert into tags (user_id, name, kind)
values ($1, $2, $3)
returning id, user_id, name, kind, created_at;

-- name: GetTag :one
select id, user_id, name, kind, created_at
from tags where id = $1 and user_id = $2;

-- name: ListTagsByUser :many
select id, user_id, name, kind, created_at
from tags where user_id = $1 order by name asc;

-- name: UpdateTag :one
update tags set
  name = coalesce(sqlc.narg(name), name),
  kind = coalesce(sqlc.narg(kind), kind)
where id = sqlc.arg(id) and user_id = sqlc.arg(user_id)
returning id, user_id, name, kind, created_at;

-- name: DeleteTag :exec
delete from tags where id = $1 and user_id = $2;
