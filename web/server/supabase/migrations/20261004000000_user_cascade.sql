-- User cascade delete: removing an auth user atomically removes their
-- profile, accounts, trades, tags and images (via existing cascades).
-- Suspend stays the reversible in-app action; hard delete is deliberate.

-- Probe leftover from diagnostics (the only orphan); nothing else is deleted.
delete from public.profiles where id = '00000000-0000-0000-0000-000000000000';

alter table public.profiles
  add constraint profiles_id_auth_fkey
  foreign key (id) references auth.users (id)
  on delete cascade;
