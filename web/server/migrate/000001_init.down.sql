-- rollback for 000001_init
drop trigger if exists on_auth_user_created on auth.users;
drop table if exists public.audit_log cascade;
drop table if exists public.trade_images cascade;
drop table if exists public.trade_tags cascade;
drop table if exists public.trades cascade;
drop table if exists public.tags cascade;
drop table if exists public.setups cascade;
drop table if exists public.accounts cascade;
drop table if exists public.site_settings cascade;
drop table if exists public.profiles cascade;
drop function if exists public.handle_new_user() cascade;
drop function if exists public.set_updated_at() cascade;
