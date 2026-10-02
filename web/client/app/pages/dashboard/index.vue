<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import {
  ChartLineData01Icon,
  Target01Icon,
  Tick02Icon,
  TradeDownIcon,
  TradeUpIcon,
  Wallet01Icon,
} from '~/utils/icons'
import { dashboardKey, setupsKey, type DashboardOverview, type Setup } from '~/types'
import { fmtMoney, fmtPct, fmtR } from '~/utils/format'

definePageMeta({ middleware: 'auth', layout: 'dashboard' })

const api = useApi()
const selectedAccount = ref<string>('all')

const { data, isPending, isError, refetch } = useQuery({
  queryKey: computed(() => dashboardKey(selectedAccount.value === 'all' ? undefined : selectedAccount.value)),
  queryFn: () => {
    const q = selectedAccount.value === 'all' ? '' : `?account=${selectedAccount.value}`
    return api.get<DashboardOverview>(`/dashboard${q}`)
  },
})

const currency = computed(() => {
  if (selectedAccount.value === 'all' || !data.value)
    return '$'
  return data.value.Accounts.find(a => a.ID === selectedAccount.value)?.Currency === 'NGN' ? '₦' : '$'
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

const equityInfo = computed(() => {
  const points = data.value?.Equity ?? []
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
  const s = data.value?.Summary
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
  (data.value?.BySetup ?? []).map(g => ({ label: g.setup_name, value: g.avg_r })),
)
const pairItems = computed(() =>
  (data.value?.ByPair ?? []).map(g => ({ label: g.pair, value: g.avg_r })),
)
</script>

<template>
  <div>
    <!-- Header: title + account switcher -->
    <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
      <h1 class="text-xl font-bold tracking-tight md:text-2xl">
        Dashboard
      </h1>
      <label class="flex items-center gap-2 text-sm">
        <span class="hidden text-muted sm:inline">Account</span>
        <select
          v-model="selectedAccount"
          class="rounded-xl border border-border bg-surface px-3 py-2 text-sm outline-none transition focus:border-primary"
        >
          <option value="all">
            All accounts
          </option>
          <option v-for="a in data?.Accounts ?? []" :key="a.ID" :value="a.ID">
            {{ a.Name }} ({{ a.Currency }})
          </option>
        </select>
      </label>
    </div>

    <!-- Loading skeletons -->
    <div v-if="isPending" class="grid grid-cols-2 gap-3 md:grid-cols-3 xl:grid-cols-6">
      <div v-for="i in 6" :key="i" class="h-28 animate-pulse rounded-2xl bg-surface" />
    </div>

    <!-- Error -->
    <div v-else-if="isError" class="rounded-2xl border border-border bg-surface p-8 text-center">
      <p class="text-sm text-muted">
        Couldn't load your dashboard.
      </p>
      <button
        type="button"
        class="mt-3 rounded-xl bg-primary px-5 py-2 text-sm font-semibold text-on-primary transition hover:opacity-90"
        @click="() => refetch()"
      >
        Retry
      </button>
    </div>

    <template v-else-if="data">
      <!-- Stat cards -->
      <div class="grid grid-cols-2 gap-3 md:grid-cols-3 xl:grid-cols-6">
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
        <section class="rounded-2xl border border-border bg-surface p-4 xl:col-span-3" aria-label="Equity curve">
          <div class="mb-2 flex flex-wrap items-end justify-between gap-2">
            <div>
              <h2 class="text-sm font-semibold text-muted">
                Equity curve
              </h2>
              <p class="tnum mt-0.5 text-2xl font-bold tracking-tight">
                {{ fmtMoney(equityInfo.current, currency) }}
                <span class="text-sm font-semibold" :class="equityInfo.change >= 0 ? 'text-profit-text' : 'text-loss'">
                  {{ fmtMoney(equityInfo.change, currency) }}
                  in {{ range === 0 ? 'all' : `${range} days` }}
                </span>
              </p>
            </div>
            <div class="flex gap-1 text-xs font-medium" role="tablist" aria-label="Equity range">
              <button
                v-for="r in ([30, 60, 0] as const)"
                :key="r"
                type="button"
                :aria-pressed="range === r"
                :class="[
                  'rounded-lg px-3 py-1.5 transition',
                  range === r ? 'bg-primary/10 text-primary' : 'text-muted hover:text-text',
                ]"
                @click="range = r"
              >
                {{ r === 0 ? 'All' : `${r}D` }}
              </button>
            </div>
          </div>
          <ChartsEquityChart :points="data.Equity" :range-days="range" />
        </section>

        <section class="rounded-2xl border border-border bg-surface p-4 xl:col-span-2" aria-label="Performance calendar">
          <ChartsCalendarGrid :days="data.Calendar" />
        </section>
      </div>

      <!-- Breakdowns -->
      <div class="mt-4 grid gap-4 md:grid-cols-2">
        <section class="rounded-2xl border border-border bg-surface p-4" aria-label="Breakdown by setup">
          <ChartsBreakdownBars title="By setup" :items="setupItems" />
        </section>
        <section class="rounded-2xl border border-border bg-surface p-4" aria-label="Breakdown by pair">
          <ChartsBreakdownBars title="By pair" :items="pairItems" />
        </section>
      </div>

      <!-- Recent trades -->
      <section class="mt-4 rounded-2xl border border-border bg-surface p-4" aria-label="Recent trades">
        <div class="mb-3 flex items-center justify-between">
          <h2 class="text-sm font-semibold">
            Recent trades
          </h2>
          <NuxtLink to="/dashboard/trades" class="text-sm font-medium text-primary hover:underline">
            View all
          </NuxtLink>
        </div>
        <TradesRecentTrades :trades="data.Recent" :setup-names="setupNames" />
      </section>
    </template>
  </div>
</template>
