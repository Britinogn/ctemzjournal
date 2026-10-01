-- profiles (id = supabase auth uid)

-- name: GetProfileByID :one
select id, display_name, role, status, timezone, avatar_path, created_at, updated_at
from profiles where id = $1;

-- name: UpdateProfile :one
update profiles set
  display_name = coalesce(sqlc.narg(display_name), display_name),
  timezone = coalesce(sqlc.narg(timezone), timezone),
  avatar_path = coalesce(sqlc.narg(avatar_path), avatar_path)
where id = $1
returning id, display_name, role, status, timezone, avatar_path, created_at, updated_at;

-- name: UpdateProfileStatus :one
update profiles set status = $2 where id = $1
returning id, display_name, role, status, timezone, avatar_path, created_at, updated_at;

-- name: ListUsers :many
select id, display_name, role, status, timezone, avatar_path, created_at, updated_at
from profiles
where ($1 = '' or display_name ilike '%' || $1 || '%')
order by created_at desc
limit $2 offset $3;

-- name: CountUsers :one
select count(*) from profiles;

-- name: CountNewUsersSince :one
select count(*) from profiles where created_at >= $1;
