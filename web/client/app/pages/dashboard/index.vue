<script setup lang="ts">
import { keepPreviousData, useQuery } from '@tanstack/vue-query'
import {
  ChartLineData01Icon,
  Target01Icon,
  Tick02Icon,
  TradeDownIcon,
  TradeUpIcon,
  Wallet01Icon,
} from '~/utils/icons'
import { dashboardKey, meKey, setupsKey, type DashboardOverview, type Profile, type Setup } from '~/types'
import { fmtMoney, fmtPct, fmtR } from '~/utils/format'

definePageMeta({ middleware: 'auth', layout: 'dashboard' })

const api = useApi()
const selectedAccount = ref<string>('all')

const { data, isPending, isError, isPlaceholderData, refetch } = useQuery({
  queryKey: computed(() => dashboardKey(selectedAccount.value === 'all' ? undefined : selectedAccount.value)),
  queryFn: () => {
    const q = selectedAccount.value === 'all' ? '' : `?account=${selectedAccount.value}`
    return api.get<DashboardOverview>(`/dashboard${q}`)
  },
  // Keeps the old numbers on screen while another account loads, so the page does not flash empty.
  placeholderData: keepPreviousData,
})

const currency = computed(() => {
  if (selectedAccount.value === 'all' || !data.value)
    return '$'
  return data.value.accounts.find(a => a.ID === selectedAccount.value)?.Currency === 'NGN' ? '₦' : '$'
})

const range = ref<0 | 30 | 60>(0)

const { data: setups } = useQuery({
  queryKey: setupsKey(),
  queryFn: () => api.get<Setup[]>('/setups'),
  staleTime: 5 * 60_000,
})

const setupNames = computed<Record<string, string>>(() =>
  Object.fromEntries((setups.value ?? []).map(s => [s.ID, s.Name])),
)

/* ---------- Greeting ---------- */
// Same key and stale time as the layout, so this reads the cached /me and makes no extra request.
const { data: me } = useQuery({
  queryKey: meKey(),
  queryFn: () => api.get<Profile>('/me'),
  staleTime: 5 * 60_000,
})
const firstName = computed(() => (me.value?.DisplayName ?? '').trim().split(/\s+/)[0] ?? '')

// The time of day is read in the browser only, so the server and browser can never disagree.
const greeting = ref('Welcome back')
onMounted(() => {
  const h = new Date().getHours()
  greeting.value = h < 12 ? 'Good morning' : h < 17 ? 'Good afternoon' : 'Good evening'
})

// One line that changes with where the trader is, with no random picks, so it stays steady between visits.
const tagline = computed(() => {
  const s = data.value?.summary
  if (!s)
    return ''
  if (s.total_trades === 0)
    return 'Log your first trade and the story begins.'
  if (s.open_trades > 0)
    return `${s.open_trades} ${s.open_trades === 1 ? 'trade is' : 'trades are'} still open. Let ${s.open_trades === 1 ? 'it' : 'them'} play out.`
  return 'Know why you win, and why you lose.'
})

const equityInfo = computed(() => {
  const points = data.value?.equity_curve ?? []
  const lastPoint = points[points.length - 1]
  const firstPoint = points[0]
  if (points.length === 0 || !lastPoint || !firstPoint)
    return { current: 0, change: 0 }
  const current = lastPoint.equity
  const last = new Date(lastPoint.date).getTime()
  const cutoff = range.value === 0 ? 0 : last - range.value * 86_400_000
  const base = points.find(p => new Date(p.date).getTime() >= cutoff)?.equity ?? firstPoint.equity
  const start = range.value === 0 ? firstPoint.equity : base
  return { current, change: current - start }
})

const statCards = computed(() => {
  const s = data.value?.summary
  if (!s)
    return []
  return [
    {
      label: 'Win rate',
      value: fmtPct(s.win_rate),
      sub: `${s.wins} of ${s.total} closed`,
      tone: 'neutral' as const,
      icon: Tick02Icon,
    },
    {
      label: 'Average R',
      value: fmtR(s.avg_r),
      sub: 'Per closed trade',
      tone: s.avg_r >= 0 ? 'profit' as const : 'loss' as const,
      icon: TradeUpIcon,
    },
    {
      label: 'Expectancy',
      value: fmtMoney(s.expectancy, currency.value),
      sub: 'Per trade, after costs',
      tone: s.expectancy >= 0 ? 'profit' as const : 'loss' as const,
      icon: ChartLineData01Icon,
    },
    {
      label: 'Max drawdown',
      value: fmtMoney(-s.max_drawdown, currency.value),
      sub: 'Peak-to-trough',
      tone: s.max_drawdown > 0 ? 'loss' as const : ('neutral' as const),
      icon: TradeDownIcon,
    },
    {
      label: 'Rules followed',
      value: fmtPct(s.rule_rate),
      sub: `On ${s.total} trades`,
      tone: 'neutral' as const,
      icon: Target01Icon,
    },
    {
      label: 'Open trades',
      value: String(s.open_trades),
      sub: `of ${s.total_trades} total`,
      tone: 'neutral' as const,
      icon: Wallet01Icon,
    },
  ]
})

const setupItems = computed(() =>
  (data.value?.by_setup ?? []).map(g => ({ label: g.setup_name, value: g.avg_r })),
)
const pairItems = computed(() =>
  (data.value?.by_pair ?? []).map(g => ({ label: g.pair, value: g.avg_r })),
)

/* ---------- Shared look, defined once ---------- */
const panel = 'rounded-2xl border border-border bg-surface p-4 md:p-5'
// Six stat cards: 2 across on phones, 3 on tablet and laptop, 6 only when there is room (the sidebar takes 256px).
const statsGrid = 'grid grid-cols-2 gap-3 sm:grid-cols-3 2xl:grid-cols-6'
const skel = 'animate-pulse rounded-2xl border border-border bg-surface'
</script>

<template>
  <div>
    <!--
      A greeting leads the page. The real h1 stays for screen readers, because the top bar
      already shows "Dashboard" from tablet up and the greeting is the visible heading on phones.
    -->
    <div class="mb-5 flex flex-wrap items-end justify-between gap-x-4 gap-y-3">
      <div class="min-w-0">
        <h1 class="sr-only">
          Dashboard
        </h1>
        <p class="truncate text-2xl font-extrabold tracking-tight md:text-3xl">
          {{ greeting }}<template v-if="firstName">
            , {{ firstName }}
          </template>
        </p>
        <!-- min-h keeps the space while the numbers load, so the page does not shift -->
        <p class="mt-1 min-h-6 text-sm text-muted md:text-base">
          {{ tagline }}
        </p>
      </div>
      <div class="flex items-center gap-2">
        <label for="dashboard-account" class="text-sm text-muted max-sm:sr-only">Account</label>
        <div class="relative">
          <select
            id="dashboard-account"
            v-model="selectedAccount"
            class="h-11 w-full min-w-44 appearance-none rounded-xl border border-border bg-surface pl-3.5 pr-10 text-sm font-medium outline-none transition-colors focus:border-primary focus-visible:ring-2 focus-visible:ring-primary/30"
          >
            <option value="all">
              All accounts
            </option>
            <option v-for="a in data?.accounts ?? []" :key="a.ID" :value="a.ID">
              {{ a.Name }} ({{ a.Currency }})
            </option>
          </select>
          <svg class="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M6 9l6 6 6-6" />
          </svg>
        </div>
      </div>
    </div>

    <!-- Loading: the same shape as the real page, so nothing jumps when data arrives -->
    <div v-if="isPending" class="space-y-4" aria-hidden="true">
      <div :class="statsGrid">
        <div v-for="i in 6" :key="i" :class="[skel, 'h-28']" />
      </div>
      <div class="grid gap-4 xl:grid-cols-5">
        <div :class="[skel, 'h-80 xl:col-span-3']" />
        <div :class="[skel, 'h-80 xl:col-span-2']" />
      </div>
      <div class="grid gap-4 md:grid-cols-2">
        <div :class="[skel, 'h-56']" />
        <div :class="[skel, 'h-56']" />
      </div>
      <div :class="[skel, 'h-64']" />
    </div>

    <!-- Error -->
    <div v-else-if="isError" class="rounded-2xl border border-border bg-surface p-8 text-center">
      <p class="text-sm text-muted">
        Couldn't load your dashboard.
      </p>
      <button
        type="button"
        class="mt-4 inline-flex h-11 items-center rounded-xl bg-primary px-6 text-sm font-semibold text-on-primary transition-colors hover:opacity-90 focus-visible:outline-2-2 focus-visible:outline-2-offset-2 focus-visible:outline-2-primary"
        @click="() => refetch()"
      >
        Retry
      </button>
    </div>

    <!-- While another account loads, the old numbers stay but fade, so it is clear they are about to change -->
    <div
      v-else-if="data"
      class="transition-opacity"
      :class="isPlaceholderData ? 'opacity-60' : ''"
      :aria-busy="isPlaceholderData"
    >
      <!-- Stat cards -->
      <div :class="statsGrid">
        <UiStatCard
          v-for="card in statCards"
          :key="card.label"
          :label="card.label"
          :value="card.value"
          :sub="card.sub"
          :tone="card.tone"
          :icon="card.icon"
        />
      </div>

      <!-- Equity + calendar -->
      <div class="mt-4 grid gap-4 xl:grid-cols-5">
        <section :class="[panel, 'xl:col-span-3']" aria-label="Equity curve">
          <div class="mb-3 flex flex-wrap items-end justify-between gap-3">
            <div>
              <h2 class="text-sm font-semibold text-muted">
                Equity curve
              </h2>
              <p class="tnum mt-0.5 flex flex-wrap items-baseline gap-x-3 text-2xl font-bold tracking-tight">
                {{ fmtMoney(equityInfo.current, currency) }}
                <!-- The arrow means the direction never relies on colour alone -->
                <span class="inline-flex items-center gap-1 text-sm font-semibold" :class="equityInfo.change >= 0 ? 'text-profit-text' : 'text-loss'">
                  <span aria-hidden="true">{{ equityInfo.change >= 0 ? '↑' : '↓' }}</span>
                  {{ fmtMoney(equityInfo.change, currency) }}
                  in {{ range === 0 ? 'all' : `${range} days` }}
                </span>
              </p>
            </div>
            <!-- A group of toggle buttons, not a tablist: nothing here switches panels -->
            <div class="inline-flex gap-0.5 rounded-xl border border-border bg-bg p-0.5 text-xs font-semibold" role="group" aria-label="Equity range">
              <button
                v-for="r in ([30, 60, 0] as const)"
                :key="r"
                type="button"
                :aria-pressed="range === r"
                :class="[
                  'h-9 min-w-11 rounded-lg px-3 transition-colors focus-visible:outline-2-2 focus-visible:outline-2-offset-2 focus-visible:outline-2-primary',
                  range === r ? 'bg-surface text-primary shadow-sm ring-1 ring-border' : 'text-muted hover:text-text',
                ]"
                @click="range = r"
              >
                {{ r === 0 ? 'All' : `${r}D` }}
              </button>
            </div>
          </div>
          <ChartsEquityChart :points="data.equity_curve" :range-days="range" />
        </section>

        <section :class="[panel, 'xl:col-span-2']" aria-label="Performance calendar">
          <ChartsCalendarGrid :days="data.calendar" />
        </section>
      </div>

      <!-- Breakdowns -->
      <div class="mt-4 grid gap-4 md:grid-cols-2">
        <section :class="panel" aria-label="Breakdown by setup">
          <ChartsBreakdownBars title="By setup" :items="setupItems" />
        </section>
        <section :class="panel" aria-label="Breakdown by pair">
          <ChartsBreakdownBars title="By pair" :items="pairItems" />
        </section>
      </div>

      <!-- Recent trades -->
      <section :class="['mt-4', panel]" aria-label="Recent trades">
        <div class="mb-3 flex items-center justify-between">
          <h2 class="text-sm font-semibold">
            Recent trades
          </h2>
          <NuxtLink
            to="/dashboard/trades"
            class="inline-flex min-h-11 items-center rounded-lg px-2 text-sm font-medium text-primary hover:underline focus-visible:outline-2-2 focus-visible:outline-2-offset-2 focus-visible:outline-2-primary"
          >
            View all
          </NuxtLink>
        </div>
        <TradesRecentTrades :trades="data.recent_trades" :setup-names="setupNames" />
      </section>
    </div>
  </div>
</template>