-- site_settings (single row id = 1)

-- name: GetSiteSettings :one
select * from site_settings where id = 1;

-- name: UpdateSiteSettings :one
update site_settings set
  site_name = coalesce(sqlc.narg(site_name), site_name),
  tagline = coalesce(sqlc.narg(tagline), tagline),
  logo_path = coalesce(sqlc.narg(logo_path), logo_path),
  favicon_path = coalesce(sqlc.narg(favicon_path), favicon_path),
  contact_email = coalesce(sqlc.narg(contact_email), contact_email),
  footer_text = coalesce(sqlc.narg(footer_text), footer_text),
  risk_disclaimer = coalesce(sqlc.narg(risk_disclaimer), risk_disclaimer),
  social_links = coalesce(sqlc.narg(social_links), social_links),
  allow_signups = coalesce(sqlc.narg(allow_signups), allow_signups),
  maintenance_mode = coalesce(sqlc.narg(maintenance_mode), maintenance_mode)
where id = 1
returning *;
