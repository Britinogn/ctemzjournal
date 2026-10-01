-- Trading Journal V1 init
-- Supabase Postgres. RLS on everywhere (second layer; Go service scopes by user_id too).
-- Times UTC; display/calendar grouping uses profiles.timezone (default Africa/Lagos) in app.

create extension if not exists "pgcrypto";

-- ---------- helpers ----------
create or replace function public.set_updated_at()
returns trigger language plpgsql as $$
begin
  new.updated_at = now();
  return new;
end $$;

-- auto-create profile on signup (auth.users insert)
create or replace function public.handle_new_user()
returns trigger language plpgsql security definer set search_path = public as $$
begin
  insert into public.profiles (id, display_name)
  values (new.id, coalesce(new.raw_user_meta_data ->> 'display_name', split_part(new.email, '@', 1)))
  on conflict (id) do nothing;
  return new;
end $$;

-- ---------- tables ----------
create table public.profiles (
  id uuid primary key,
  display_name text,
  role text not null default 'user' check (role in ('user', 'admin')),
  status text not null default 'active' check (status in ('active', 'suspended')),
  timezone text not null default 'Africa/Lagos',
  avatar_path text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table public.accounts (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references public.profiles (id) on delete cascade,
  name text not null,
  type text not null check (type in ('demo', 'live')),
  currency text not null default 'USD',
  starting_balance numeric(18, 2) not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table public.setups (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references public.profiles (id) on delete cascade,
  name text not null,
  rules text,
  invalidation text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (user_id, name)
);

create table public.tags (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references public.profiles (id) on delete cascade,
  name text not null,
  kind text not null default 'general' check (kind in ('mistake', 'general')),
  created_at timestamptz not null default now(),
  unique (user_id, name)
);

create table public.trades (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references public.profiles (id) on delete cascade,
  account_id uuid not null references public.accounts (id) on delete cascade,
  setup_id uuid references public.setups (id) on delete set null,
  pair text not null,
  direction text not null check (direction in ('long', 'short')),
  timeframe text,
  opened_at timestamptz,
  closed_at timestamptz,
  entry numeric(18, 6),
  stop_loss numeric(18, 6),
  take_profit numeric(18, 6),
  exit_price numeric(18, 6),
  lot_size numeric(18, 2),
  commission numeric(18, 2) not null default 0,
  swap numeric(18, 2) not null default 0,
  risk_amount numeric(18, 2),
  pnl numeric(18, 2),
  r_multiple numeric(18, 4),
  followed_rules boolean,
  emotion text,
  notes text,
  status text not null default 'open' check (status in ('open', 'closed')),
  is_public boolean not null default false,
  hidden_by_admin boolean not null default false,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table public.trade_tags (
  trade_id uuid not null references public.trades (id) on delete cascade,
  tag_id uuid not null references public.tags (id) on delete cascade,
  primary key (trade_id, tag_id)
);

create table public.trade_images (
  id uuid primary key default gen_random_uuid(),
  trade_id uuid not null references public.trades (id) on delete cascade,
  public_id text not null unique,
  kind text check (kind in ('entry', 'exit')),
  position int not null default 0,
  created_at timestamptz not null default now()
);

-- single-row site settings (id = 1)
create table public.site_settings (
  id int primary key check (id = 1),
  site_name text not null default 'CTemz Journal',
  tagline text,
  logo_path text,
  favicon_path text,
  contact_email text,
  footer_text text,
  risk_disclaimer text,
  social_links jsonb not null default '{}',
  allow_signups boolean not null default true,
  maintenance_mode boolean not null default false,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table public.audit_log (
  id uuid primary key default gen_random_uuid(),
  admin_id uuid references public.profiles (id) on delete set null,
  action text not null,
  target_type text not null,
  target_id text not null,
  meta jsonb not null default '{}',
  created_at timestamptz not null default now()
);

-- ---------- indexes ----------
create index accounts_user_id_idx on public.accounts (user_id);
create index setups_user_id_idx on public.setups (user_id);
create index tags_user_id_idx on public.tags (user_id);
create index trades_user_opened_idx on public.trades (user_id, opened_at desc);
create index trades_account_idx on public.trades (account_id);
create index trades_pair_idx on public.trades (pair);
create index trades_status_idx on public.trades (status);
create index trades_public_idx on public.trades (is_public, hidden_by_admin) where is_public = true;
create index trade_images_trade_idx on public.trade_images (trade_id);
create index audit_log_admin_idx on public.audit_log (admin_id);
create index audit_log_created_idx on public.audit_log (created_at desc);

-- ---------- updated_at triggers ----------
create trigger trg_profiles_updated before update on public.profiles
  for each row execute function public.set_updated_at();
create trigger trg_accounts_updated before update on public.accounts
  for each row execute function public.set_updated_at();
create trigger trg_setups_updated before update on public.setups
  for each row execute function public.set_updated_at();
create trigger trg_trades_updated before update on public.trades
  for each row execute function public.set_updated_at();
create trigger trg_site_settings_updated before update on public.site_settings
  for each row execute function public.set_updated_at();

-- ---------- new-user trigger ----------
drop trigger if exists on_auth_user_created on auth.users;
create trigger on_auth_user_created
  after insert on auth.users
  for each row execute function public.handle_new_user();

-- ---------- seed ----------
insert into public.site_settings (id) values (1) on conflict (id) do nothing;

-- public logo/favicon bucket (Supabase Storage)
insert into storage.buckets (id, name, public)
values ('site-assets', 'site-assets', true)
on conflict (id) do nothing;

-- ---------- RLS (second layer) ----------
alter table public.profiles enable row level security;
alter table public.accounts enable row level security;
alter table public.setups enable row level security;
alter table public.tags enable row level security;
alter table public.trades enable row level security;
alter table public.trade_tags enable row level security;
alter table public.trade_images enable row level security;
alter table public.site_settings enable row level security;
alter table public.audit_log enable row level security;

-- profiles: owner read/update own only (admin writes go via service_role, bypasses RLS)
create policy profiles_owner_select on public.profiles
  for select to authenticated using (auth.uid() = id);
create policy profiles_owner_update on public.profiles
  for update to authenticated using (auth.uid() = id) with check (auth.uid() = id);

-- owner tables: full access to own rows only
create policy accounts_owner_all on public.accounts
  for all to authenticated using (auth.uid() = user_id) with check (auth.uid() = user_id);
create policy setups_owner_all on public.setups
  for all to authenticated using (auth.uid() = user_id) with check (auth.uid() = user_id);
create policy tags_owner_all on public.tags
  for all to authenticated using (auth.uid() = user_id) with check (auth.uid() = user_id);
create policy trades_owner_all on public.trades
  for all to authenticated using (auth.uid() = user_id) with check (auth.uid() = user_id);

-- trade_tags: only via trades you own
create policy trade_tags_owner_all on public.trade_tags
  for all to authenticated
  using (exists (select 1 from public.trades t where t.id = trade_tags.trade_id and t.user_id = auth.uid()))
  with check (exists (select 1 from public.trades t where t.id = trade_tags.trade_id and t.user_id = auth.uid()));

-- trade_images: only via trades you own
create policy trade_images_owner_all on public.trade_images
  for all to authenticated
  using (exists (select 1 from public.trades t where t.id = trade_images.trade_id and t.user_id = auth.uid()))
  with check (exists (select 1 from public.trades t where t.id = trade_images.trade_id and t.user_id = auth.uid()));

-- site_settings: public read, no write policies (admin writes via service_role)
create policy site_settings_public_read on public.site_settings
  for select to anon, authenticated using (true);

-- public journals read: anyone can read public, non-hidden trades
-- (app public endpoint adds column filtering + signed image URLs; RLS only gates rows)
create policy trades_public_read on public.trades
  for select to anon, authenticated using (is_public = true and hidden_by_admin = false);

-- audit_log: no policies (service_role only)
