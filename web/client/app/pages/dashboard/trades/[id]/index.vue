<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { toast } from 'vue-sonner'
import type { Trade, TradeImageWithUrl } from '~/types'
import { fmtDay, fmtMoney, fmtR } from '~/utils/format'
import {
  ArrowLeft01Icon,
  Edit02Icon,
  Delete02Icon,
  ChartIcon,
} from '~/utils/icons'

definePageMeta({ middleware: 'auth', layout: 'dashboard' })

const route = useRoute()
const tradeId = computed(() => String(route.params.id ?? ''))
const api = useApi()
const queryClient = useQueryClient()

const { data: trade, isPending, isError } = useQuery({
  queryKey: ['trade', tradeId.value],
  queryFn: () => api.get<Trade>(`/trades/${tradeId.value}`),
})

const { data: images } = useQuery({
  queryKey: ['trade-images', tradeId.value],
  queryFn: () => api.get<TradeImageWithUrl[]>(`/trades/${tradeId.value}/images`),
})

const { data: accounts } = useAccounts()
const { data: setups } = useSetups()

const account = computed(() =>
  accounts.value?.find(a => a.ID === trade.value?.AccountID),
)

const currency = computed(() =>
  account.value?.Currency === 'NGN' ? '₦' : '$',
)

const setupName = computed(() => {
  if (!trade.value?.SetupID)
    return 'No setup'
  return setups.value?.find(s => s.ID === trade.value?.SetupID)?.Name ?? 'Setup'
})

const moneyRows = computed(() => {
  const t = trade.value
  if (!t)
    return []
  const m = (label: string, v: number | null, money = true) => ({
    label,
    display: v === null ? '—' : money ? fmtMoney(v, currency.value) : String(v),
    tone: v === null ? 'muted' : v < 0 && money ? 'loss' : 'text',
  })
  return [
    m('Entry', t.Entry, false),
    m('Stop loss', t.StopLoss, false),
    m('Take profit', t.TakeProfit, false),
    m('Exit price', t.ExitPrice, false),
    m('Lot size', t.LotSize, false),
    m('Commission', t.Commission),
    m('Swap', t.Swap),
    m('Risk', t.RiskAmount),
    m('Net P&L', t.Pnl),
  ]
})

const toggling = ref(false)
async function onToggleVisibility(): Promise<void> {
  if (!trade.value || toggling.value)
    return
  toggling.value = true
  try {
    const updated = await api.patch<Trade>(`/trades/${trade.value.ID}/visibility`, {
      is_public: !trade.value.IsPublic,
    })
    queryClient.setQueryData(['trade', tradeId.value], updated)
    toast.success(updated.IsPublic ? 'Trade is now public' : 'Trade is now private')
  }
  catch {
    toast.error('Could not change visibility')
  }
  finally {
    toggling.value = false
  }
}

const confirmDelete = ref(false)
const deleting = ref(false)
async function onDelete(): Promise<void> {
  if (!trade.value || deleting.value)
    return
  deleting.value = true
  try {
    await api.del(`/trades/${trade.value.ID}`)
    queryClient.invalidateQueries({ queryKey: ['trades'] })
    queryClient.invalidateQueries({ queryKey: ['dashboard'] })
    toast.success('Trade deleted')
    await navigateTo('/dashboard/trades')
  }
  catch {
    toast.error('Could not delete trade')
  }
  finally {
    deleting.value = false
    confirmDelete.value = false
  }
}

/* ---------- Shared look ---------- */
const panel = 'rounded-2xl border border-border bg-surface p-4 md:p-5'
</script>

<template>
  <div>
    <!-- Mobile heading -->
    <h1 class="mb-4 text-xl font-bold tracking-tight md:sr-only">
      Trade
    </h1>

    <!-- Top actions -->
    <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
      <NuxtLink
        to="/dashboard/trades"
        class="inline-flex items-center gap-1.5 text-sm font-medium text-muted transition hover:text-primary"
      >
        <UiAppIcon :icon="ArrowLeft01Icon" :size="16" />
        Trades
      </NuxtLink>

      <div v-if="trade" class="flex gap-2">
        <NuxtLink
          :to="`/dashboard/trades/${trade.ID}/edit`"
          class="inline-flex h-10 items-center gap-1.5 rounded-xl border border-border bg-bg px-4 text-sm font-semibold transition hover:border-primary"
        >
          <UiAppIcon :icon="Edit02Icon" :size="16" />
          Edit
        </NuxtLink>
        <button
          type="button"
          class="inline-flex h-10 items-center gap-1.5 rounded-xl border border-loss/40 px-4 text-sm font-semibold text-loss transition hover:bg-loss/10"
          @click="confirmDelete = true"
        >
          <UiAppIcon :icon="Delete02Icon" :size="16" />
          Delete
        </button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="isPending" class="space-y-3" aria-hidden="true">
      <div class="h-32 animate-pulse rounded-2xl border border-border bg-surface" />
      <div class="h-48 animate-pulse rounded-2xl border border-border bg-surface" />
      <div class="h-40 animate-pulse rounded-2xl border border-border bg-surface" />
    </div>

    <!-- Error -->
    <div
      v-else-if="isError || !trade"
      :class="[panel, 'text-center']"
    >
      <p class="text-sm text-muted">
        Trade not found. It may have been deleted.
      </p>
      <NuxtLink
        to="/dashboard/trades"
        class="mt-3 inline-flex h-11 items-center rounded-xl border border-border bg-bg px-5 text-sm font-semibold transition hover:border-primary"
      >
        Back to trades
      </NuxtLink>
    </div>

    <template v-else>
      <!-- Header card -->
      <section :class="panel" aria-label="Trade summary">
        <div class="flex flex-wrap items-center gap-2">
          <h1 class="tnum text-2xl font-bold tracking-tight">
            {{ trade.Pair }}
          </h1>
          <span
            class="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-semibold"
            :class="trade.Direction === 'long'
              ? 'bg-profit/10 text-profit-text'
              : 'bg-loss/10 text-loss'"
          >
            {{ trade.Direction === 'long' ? '↑ Long' : '↓ Short' }}
          </span>
          <span
            class="rounded-full px-2.5 py-0.5 text-xs font-medium"
            :class="trade.Status === 'open'
              ? 'bg-warning/10 text-warning-text'
              : 'bg-bg text-muted'"
          >
            {{ trade.Status === 'open' ? 'Open' : 'Closed' }}
          </span>
          <span
            v-if="trade.IsPublic"
            class="rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium text-primary"
          >
            Public
          </span>
        </div>

        <p class="mt-1.5 text-sm text-muted">
          {{ setupName }}
          <template v-if="trade.Timeframe"> · {{ trade.Timeframe }}</template>
          · {{ account?.Name ?? 'Account' }}
          · {{ trade.ClosedAt ? fmtDay(trade.ClosedAt) : 'Open' }}
        </p>

        <div
          v-if="trade.RMultiple !== null"
          class="mt-4 flex items-end gap-4"
        >
          <p
            class="tnum text-3xl font-bold"
            :class="(trade.Pnl ?? 0) >= 0 ? 'text-profit-text' : 'text-loss'"
          >
            {{ fmtMoney(trade.Pnl ?? 0, currency) }}
          </p>
          <p
            class="tnum text-lg font-bold"
            :class="trade.RMultiple >= 0 ? 'text-profit-text' : 'text-loss'"
          >
            {{ fmtR(trade.RMultiple) }}
          </p>
        </div>
      </section>

      <!-- Numbers -->
      <section :class="[panel, 'mt-4']" aria-label="Numbers">
        <h2 class="text-sm font-semibold">
          Numbers
        </h2>

        <dl class="mt-3 grid grid-cols-2 gap-x-4 gap-y-3 sm:grid-cols-3">
          <div v-for="row in moneyRows" :key="row.label">
            <dt class="text-xs text-muted">
              {{ row.label }}
            </dt>
            <dd
              class="tnum text-sm font-semibold"
              :class="{
                'text-loss': row.tone === 'loss',
                'text-muted': row.tone === 'muted',
              }"
            >
              {{ row.display }}
            </dd>
          </div>
        </dl>

        <dl class="mt-4 grid grid-cols-2 gap-x-4 gap-y-3 border-t border-border pt-4 sm:grid-cols-3">
          <div>
            <dt class="text-xs text-muted">Emotion</dt>
            <dd class="text-sm font-medium">
              {{ trade.Emotion ?? '—' }}
            </dd>
          </div>
          <div>
            <dt class="text-xs text-muted">Rules</dt>
            <dd class="text-sm font-medium">
              {{ trade.FollowedRules === null ? '—' : trade.FollowedRules ? 'Followed' : 'Broken' }}
            </dd>
          </div>
          <div>
            <dt class="text-xs text-muted">Opened</dt>
            <dd class="tnum text-sm font-medium">
              {{ trade.OpenedAt ? fmtDay(trade.OpenedAt) : '—' }}
            </dd>
          </div>
        </dl>

        <div
          v-if="trade.Notes"
          class="mt-4 border-t border-border pt-4"
        >
          <p class="text-xs text-muted">
            Notes
          </p>
          <p class="mt-1 whitespace-pre-wrap text-sm">
            {{ trade.Notes }}
          </p>
        </div>
      </section>

      <!-- Screenshots -->
      <section :class="[panel, 'mt-4']" aria-label="Screenshots">
        <div class="mb-3 flex items-center justify-between">
          <h2 class="text-sm font-semibold">
            Screenshots
          </h2>
          <span class="tnum text-xs text-muted">
            {{ (images ?? []).length }} total
          </span>
        </div>

        <div
          v-if="(images ?? []).length > 0"
          class="grid grid-cols-2 gap-2 sm:grid-cols-3"
        >
          <figure
            v-for="img in images"
            :key="img.ID"
            class="overflow-hidden rounded-xl border border-border bg-bg"
          >
            <img
              :src="img.URL"
              :alt="`${img.Kind ?? 'Trade'} screenshot`"
              class="aspect-video w-full object-cover"
              loading="lazy"
            >
            <figcaption class="px-2 py-1.5 text-center text-[11px] capitalize text-muted">
              {{ img.Kind ?? 'chart' }}
            </figcaption>
          </figure>
        </div>

        <div
          v-else
          class="flex flex-col items-center rounded-xl border border-dashed border-border px-4 py-8 text-center"
        >
          <span class="inline-flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 text-primary">
            <UiAppIcon :icon="ChartIcon" :size="20" />
          </span>
          <p class="mt-2 text-sm text-muted">
            No screenshots on this trade.
          </p>
        </div>
      </section>

      <!-- Visibility -->
      <section :class="[panel, 'mt-4']" aria-label="Visibility">
        <label class="flex cursor-pointer items-center justify-between gap-3">
          <span>
            <span class="block text-sm font-medium">
              Make public
            </span>
            <span class="block text-xs text-muted">
              Shows pair, setup and R only. Never money or lots.
            </span>
          </span>
          <input
            type="checkbox"
            class="peer sr-only"
            :checked="trade.IsPublic"
            :disabled="toggling"
            @change="onToggleVisibility"
          >
          <span
            class="relative h-6 w-11 shrink-0 rounded-full bg-border transition
            peer-checked:bg-primary
            peer-focus-visible:ring-2 peer-focus-visible:ring-primary
            after:absolute after:left-0.5 after:top-0.5 after:h-5 after:w-5
            after:rounded-full after:bg-white after:transition
            peer-checked:after:translate-x-5"
          />
        </label>
      </section>
    </template>

    <UiConfirmDialog
      :open="confirmDelete"
      title="Delete trade?"
      message="The trade and its screenshots will be gone for good."
      :busy="deleting"
      @confirm="onDelete"
      @close="confirmDelete = false"
    />
  </div>
</template>