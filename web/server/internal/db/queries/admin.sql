-- admin aggregates (overview cards, signup chart, users-with-counts)

-- name: CountUsersByStatus :one
select count(*) from profiles where status = $1;

-- name: CountPublicTrades :one
select count(*) from trades where is_public = true and hidden_by_admin = false;

-- name: CountHiddenJournals :one
select count(*) from trades where hidden_by_admin = true;

-- name: SignupsByDay :many
select date_trunc('day', created_at)::date as day, count(*) as signups
from profiles
where created_at >= sqlc.arg(since)::timestamptz
group by 1
order by 1;

-- name: ListUsersAdmin :many
select p.id, p.display_name, p.email, p.role, p.status, p.timezone,
  p.avatar_path, p.created_at, p.updated_at,
  (select count(*) from trades t where t.user_id = p.id) as trade_count
from profiles p
where (sqlc.arg(search)::text = '' or p.display_name ilike '%' || sqlc.arg(search)::text || '%' or p.email ilike '%' || sqlc.arg(search)::text || '%')
order by p.created_at desc
limit sqlc.arg(page_limit) offset sqlc.arg(page_offset);
