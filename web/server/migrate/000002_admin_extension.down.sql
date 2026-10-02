-- rollback for 000002 (email column only; trigger keeps newest version)
alter table if exists public.profiles drop column if exists email;
