-- rollback for 000003 (drops the auth cascade only)
alter table if exists public.profiles drop constraint if exists profiles_id_auth_fkey;
