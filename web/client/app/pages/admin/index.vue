<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { Tick02Icon } from '~/utils/icons'
import type { AdminOverview, AuditEntry, RatesResponse } from '~/types'
import { adminAuditKey, adminOverviewKey, ratesKey } from '~/types'
import { timeAgo } from '~/utils/format'

definePageMeta({ middleware: 'admin', layout: 'admin' })

const api = useApi()

const { data: overview, isPending, isError, refetch } = useQuery({
  queryKey: adminOverviewKey(),
  queryFn: () => api.get<AdminOverview>('/admin/overview'),
})

const { data: recentAudit } = useQuery({
  queryKey: adminAuditKey(0),
  queryFn: () => api.get<AuditEntry[]>('/admin/audit?limit=5'),
  staleTime: 60_000,
})

const { data: rates } = useQuery({
  queryKey: ratesKey(),
  queryFn: () => api.get<RatesResponse>('/public/rates'),
  staleTime: 60_000,
})

const weekDelta = computed(() => {
  const cur = overview.value?.new_users_this_week ?? 0
  const prev = overview.value?.new_users_prev_week ?? 0
  if (prev <= 0)
    return cur > 0 ? `+${cur} new` : 'No signups yet'
  const pct = Math.round(((cur - prev) / prev) * 100)
  return `${pct >= 0 ? 'Up' : 'Down'} ${Math.abs(pct)}% on last week`
})

/** Last 7 calendar days filled with signup counts (zeros included). */
const signupBars = computed(() => {
  const byDate = new Map((overview.value?.signups_last_7d ?? []).map(d => [d.date, d.count]))
  const bars: Array<{ label: string; count: number; peak: boolean }> = []
  const today = new Date()
  const counts: number[] = []
  for (let i = 6; i >= 0; i--) {
    const d = new Date(Date.UTC(today.getUTCFullYear(), today.getUTCMonth(), today.getUTCDate() - i))
    const key = d.toISOString().slice(0, 10)
    const count = byDate.get(key) ?? 0
    counts.push(count)
    bars.push({
      label: d.toLocaleDateString('en-GB', { weekday: 'short', timeZone: 'UTC' }),
      count,
      peak: false,
    })
  }
  const max = Math.max(0, ...counts)
  for (const b of bars) {
    if (b.count > 0 && b.count === max)
      b.peak = true
  }
  const peakDay = bars.find(b => b.peak)?.label
  return { bars, max: Math.max(1, max), peakDay }
})

const ratesFresh = computed(() => {
  if (!rates.value?.updated_at)
    return { label: 'No data', fresh: false }
  if (rates.value.stale)
    return { label: `Stale, ${timeAgo(rates.value.updated_at)}`, fresh: false }
  return { label: `Fresh, ${timeAgo(rates.value.updated_at)}`, fresh: true }
})

function actionLabel(action: string): string {
  return action
    .split('.')
    .map(part => part.replace(/_/g, ' '))
    .map(part => part.charAt(0).toUpperCase() + part.slice(1))
    .join(' · ')
}
</script>

<template>
  <div>
    <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
      <h1 class="text-xl font-bold tracking-tight md:text-2xl">
        Admin overview
      </h1>
      <ThemeToggle />
    </div>

    <div v-if="isPending" class="grid grid-cols-1 gap-3 sm:grid-cols-3">
      <div v-for="i in 3" :key="i" class="h-28 animate-pulse rounded-2xl bg-surface" />
    </div>
    <div v-else-if="isError" class="rounded-2xl border border-border bg-surface p-8 text-center">
      <p class="text-sm text-muted">
        Couldn't load overview.
      </p>
      <button
        type="button"
        class="mt-3 rounded-xl bg-primary px-5 py-2 text-sm font-semibold text-on-primary transition hover:opacity-90"
        @click="() => refetch()"
      >
        Retry
      </button>
    </div>

    <template v-else-if="overview">
      <div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
        <div class="rounded-2xl border border-border bg-surface p-4">
          <p class="text-xs font-medium text-muted">Users</p>
          <p class="tnum mt-1 text-2xl font-bold tracking-tight">
            {{ overview.user_count.toLocaleString('en-US') }}
          </p>
          <p class="mt-0.5 text-xs text-muted">
            {{ overview.active_users.toLocaleString('en-US') }} active, {{ overview.suspended_users.toLocaleString('en-US') }} suspended
          </p>
        </div>
        <div class="rounded-2xl border border-border bg-surface p-4">
          <p class="text-xs font-medium text-muted">Trades logged</p>
          <p class="tnum mt-1 text-2xl font-bold tracking-tight">
            {{ overview.trade_count.toLocaleString('en-US') }}
          </p>
          <p class="mt-0.5 text-xs text-muted">
            {{ overview.public_trades.toLocaleString('en-US') }} public
          </p>
        </div>
        <div class="rounded-2xl border border-border bg-surface p-4">
          <p class="text-xs font-medium text-muted">New users this week</p>
          <p class="tnum mt-1 text-2xl font-bold tracking-tight text-profit-text">
            +{{ overview.new_users_this_week.toLocaleString('en-US') }}
          </p>
          <p class="mt-0.5 text-xs text-muted">
            {{ weekDelta }}
          </p>
        </div>
      </div>

      <div class="mt-4 grid gap-4 xl:grid-cols-5">
        <section class="rounded-2xl border border-border bg-surface p-4 xl:col-span-3" aria-label="Signups, last 7 days">
          <div class="mb-3 flex items-baseline justify-between">
            <h2 class="text-sm font-semibold">
              Signups, last 7 days
            </h2>
            <p v-if="signupBars.peakDay" class="text-xs text-muted">
              Peak on {{ signupBars.peakDay }}
            </p>
          </div>
          <div class="flex h-40 items-end gap-2">
            <div v-for="b in signupBars.bars" :key="b.label" class="flex min-w-0 flex-1 flex-col items-center gap-1">
              <span class="tnum text-xs font-bold">{{ b.count }}</span>
              <div
                class="w-full max-w-10 rounded-t-lg"
                :class="b.peak ? 'bg-primary' : 'bg-primary/40'"
                :style="{ height: `${Math.max(4, (b.count / signupBars.max) * 110)}px` }"
              />
              <span class="text-[11px] text-muted">{{ b.label }}</span>
            </div>
          </div>
        </section>

        <section class="rounded-2xl border border-border bg-surface p-4 xl:col-span-2" aria-label="Needs a look">
          <h2 class="mb-1 text-sm font-semibold">
            Needs a look
          </h2>
          <ul class="divide-y divide-border text-sm">
            <li class="flex items-center justify-between gap-2 py-2.5">
              <span>Exchange rates</span>
              <span
                class="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-semibold"
                :class="ratesFresh.fresh ? 'bg-profit/10 text-profit-text' : 'bg-warning/10 text-warning-text'"
              >
                <UiAppIcon v-if="ratesFresh.fresh" :icon="Tick02Icon" :size="12" />
                {{ ratesFresh.label }}
              </span>
            </li>
            <li class="flex items-center justify-between gap-2 py-2.5">
              <span>Suspended accounts</span>
              <NuxtLink
                to="/admin/users"
                class="rounded-full px-2.5 py-0.5 text-xs font-semibold"
                :class="overview.suspended_users > 0 ? 'bg-warning/10 text-warning-text' : 'bg-bg text-muted'"
              >
                {{ overview.suspended_users > 0 ? `⚠ ${overview.suspended_users}` : overview.suspended_users }}
              </NuxtLink>
            </li>
            <li class="flex items-center justify-between gap-2 py-2.5">
              <span>Hidden journals</span>
              <NuxtLink to="/admin/journals" class="rounded-full bg-bg px-2.5 py-0.5 text-xs font-semibold text-muted hover:text-text">
                {{ overview.hidden_journals }}
              </NuxtLink>
            </li>
            <li class="flex items-center justify-between gap-2 py-2.5">
              <span>Signups</span>
              <NuxtLink to="/admin/settings" class="rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-semibold text-primary">
                Open
              </NuxtLink>
            </li>
          </ul>
        </section>
      </div>

      <section class="mt-4 rounded-2xl border border-border bg-surface p-4" aria-label="Latest admin actions">
        <div class="mb-2 flex items-center justify-between">
          <h2 class="text-sm font-semibold">
            Latest admin actions
          </h2>
          <NuxtLink to="/admin/audit" class="text-sm font-medium text-primary hover:underline">
            View audit log
          </NuxtLink>
        </div>
        <ul v-if="(recentAudit ?? []).length > 0" class="divide-y divide-border">
          <li v-for="row in recentAudit" :key="row.ID" class="flex items-center justify-between gap-2 py-2 text-sm">
            <p class="min-w-0 truncate">
              <span class="font-medium">{{ actionLabel(row.Action) }}</span>
              <span class="tnum ml-2 text-xs text-muted">{{ row.TargetID.slice(0, 8) }}</span>
            </p>
          </li>
        </ul>
        <p v-else class="py-4 text-center text-sm text-muted">
          No admin actions yet.
        </p>
      </section>
    </template>
  </div>
</template>
