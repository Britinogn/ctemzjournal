-- audit_log (every admin write)

-- name: CreateAuditLog :one
insert into audit_log (admin_id, action, target_type, target_id, meta)
values ($1, $2, $3, $4, $5)
returning id, admin_id, action, target_type, target_id, meta, created_at;

-- name: ListAuditLog :many
select id, admin_id, action, target_type, target_id, meta, created_at
from audit_log order by created_at desc limit $1 offset $2;
