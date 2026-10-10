-- Per-user timezone on signup: trigger reads raw_user_meta_data ->> 'timezone'.
-- Empty/invalid values fall back to the column default ('Africa/Lagos').
-- Existing rows untouched.

create or replace function public.handle_new_user()
returns trigger language plpgsql security definer set search_path = public as $$
begin
  insert into public.profiles (id, display_name, email, timezone)
  values (
    new.id,
    coalesce(new.raw_user_meta_data ->> 'display_name', split_part(new.email, '@', 1)),
    new.email,
    coalesce(nullif(new.raw_user_meta_data ->> 'timezone', ''), 'Africa/Lagos')
  )
  on conflict (id) do nothing;
  return new;
end $$;
