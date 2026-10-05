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

// The keys are computed, so moving from one trade to another refetches instead of showing the old one.
const { data: trade, isPending, isError } = useQuery({
  queryKey: computed(() => ['trade', tradeId.value]),
  queryFn: () => api.get<Trade>(`/trades/${tradeId.value}`),
})

const { data: images } = useQuery({
  queryKey: computed(() => ['trade-images', tradeId.value]),
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

/* ---------- Numbers: three blocks, like a receipt ---------- */
function row(label: string, v: number | null, money = false) {
  return {
    label,
    display: v === null ? '—' : money ? fmtMoney(v, currency.value) : String(v),
    tone: v === null ? 'muted' : v < 0 && money ? 'loss' : 'text',
  }
}

// Prices, then size and costs, then risk. Net P&L has its own highlighted row underneath.
const numberBlocks = computed(() => {
  const t = trade.value
  if (!t)
    return []
  return [
    [row('Entry', t.Entry), row('Stop loss', t.StopLoss), row('Take profit', t.TakeProfit), row('Exit price', t.ExitPrice)],
    [row('Lot size', t.LotSize), row('Commission', t.Commission, true), row('Swap', t.Swap, true)],
    // With no P&L yet, the dash for Net P&L sits in this block instead of the highlighted row.
    t.Pnl === null ? [row('Risk', t.RiskAmount, true), row('Net P&L', null, true)] : [row('Risk', t.RiskAmount, true)],
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
const skel = 'animate-pulse rounded-2xl border border-border bg-surface'
</script>

<template>
  <div>
    <!-- Top actions -->
    <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
      <NuxtLink
        to="/dashboard/trades"
        class="-ml-1 inline-flex min-h-11 items-center gap-1.5 rounded-lg px-1 text-sm font-medium text-muted transition-colors hover:text-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
      >
        <UiAppIcon :icon="ArrowLeft01Icon" :size="16" />
        Trades
      </NuxtLink>

      <div v-if="trade" class="flex gap-2">
        <NuxtLink
          :to="`/dashboard/trades/${trade.ID}/edit`"
          class="inline-flex h-11 items-center gap-1.5 rounded-xl border border-border bg-bg px-4 text-sm font-semibold transition-colors hover:border-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
        >
          <UiAppIcon :icon="Edit02Icon" :size="16" />
          Edit
        </NuxtLink>
        <button
          type="button"
          class="inline-flex h-11 items-center gap-1.5 rounded-xl border border-loss/40 px-4 text-sm font-semibold text-loss transition-colors hover:bg-loss/10 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-loss"
          @click="confirmDelete = true"
        >
          <UiAppIcon :icon="Delete02Icon" :size="16" />
          Delete
        </button>
      </div>
    </div>

    <!-- Loading: the same shape as the real page -->
    <div v-if="isPending" class="space-y-4" aria-hidden="true">
      <div :class="[skel, 'h-32']" />
      <div class="grid gap-4 xl:grid-cols-[minmax(0,1.1fr)_minmax(0,1fr)]">
        <div :class="[skel, 'h-80']" />
        <div :class="[skel, 'h-80']" />
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="isError || !trade" :class="[panel, 'text-center']">
      <h1 class="sr-only">
        Trade
      </h1>
      <p class="text-sm text-muted">
        Trade not found. It may have been deleted.
      </p>
      <NuxtLink
        to="/dashboard/trades"
        class="mt-4 inline-flex h-11 items-center rounded-xl border border-border bg-bg px-5 text-sm font-semibold transition-colors hover:border-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
      >
        Back to trades
      </NuxtLink>
    </div>

    <template v-else>
      <!-- Header card -->
      <section :class="panel" aria-label="Trade summary">
        <div class="flex flex-wrap items-center gap-2">
          <h1 class="tnum text-2xl font-bold tracking-tight md:text-3xl">
            {{ trade.Pair }}
          </h1>
          <!-- Long and short are not results, so they get an arrow and a word, not green or red -->
          <span class="inline-flex items-center gap-1 rounded-full border border-border px-2.5 py-0.5 text-xs font-semibold">
            <span aria-hidden="true">{{ trade.Direction === 'long' ? '↑' : '↓' }}</span>
            {{ trade.Direction === 'long' ? 'Long' : 'Short' }}
          </span>
          <span
            class="rounded-full px-2.5 py-0.5 text-xs font-medium"
            :class="trade.Status === 'open' ? 'border border-border font-semibold' : 'bg-bg text-muted'"
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

        <!-- The result: signed amount with an arrow, then the R as a win or loss chip -->
        <div v-if="trade.RMultiple !== null" class="mt-5 flex flex-wrap items-center gap-x-4 gap-y-2">
          <p
            class="tnum inline-flex items-center gap-1.5 text-3xl font-bold"
            :class="(trade.Pnl ?? 0) >= 0 ? 'text-profit-text' : 'text-loss'"
          >
            <span aria-hidden="true">{{ (trade.Pnl ?? 0) >= 0 ? '↑' : '↓' }}</span>
            {{ fmtMoney(trade.Pnl ?? 0, currency) }}
          </p>
          <span
            class="tnum inline-flex items-center gap-1 rounded-full px-3 py-1 text-sm font-semibold"
            :class="trade.RMultiple >= 0 ? 'bg-profit/10 text-profit-text' : 'bg-loss/10 text-loss'"
          >
            <span aria-hidden="true">{{ trade.RMultiple >= 0 ? '↑' : '↓' }}</span>
            {{ trade.RMultiple >= 0 ? 'Win' : 'Loss' }} {{ fmtR(trade.RMultiple) }}
          </span>
        </div>
      </section>

      <!-- Two columns on a wide screen: the numbers, then the screenshots and sharing -->
      <div class="mt-4 grid items-start gap-4 xl:grid-cols-[minmax(0,1.1fr)_minmax(0,1fr)]">
        <!-- Numbers -->
        <section :class="panel" aria-label="Numbers">
          <h2 class="text-sm font-semibold">
            Numbers
          </h2>

          <!-- Three blocks of rows: label on the left, value on the right, a line between blocks -->
          <dl
            v-for="(block, bi) in numberBlocks"
            :key="bi"
            :class="bi === 0 ? 'mt-3' : 'mt-3 border-t border-border pt-3'"
          >
            <div v-for="r in block" :key="r.label" class="flex min-h-10 items-center justify-between gap-4">
              <dt class="text-sm text-muted">
                {{ r.label }}
              </dt>
              <dd
                class="tnum text-sm font-semibold"
                :class="{ 'text-loss': r.tone === 'loss', 'text-muted': r.tone === 'muted' }"
              >
                {{ r.display }}
              </dd>
            </div>
          </dl>

          <!-- The bottom line gets its own tinted row, with an arrow and a sign -->
          <dl v-if="trade.Pnl !== null" class="mt-3">
            <div
              class="flex items-center justify-between gap-4 rounded-xl px-4 py-3"
              :class="trade.Pnl >= 0 ? 'bg-profit/10' : 'bg-loss/10'"
            >
              <dt class="text-sm font-semibold">
                Net P&amp;L
              </dt>
              <dd
                class="tnum inline-flex items-center gap-1.5 text-lg font-bold"
                :class="trade.Pnl >= 0 ? 'text-profit-text' : 'text-loss'"
              >
                <span aria-hidden="true">{{ trade.Pnl >= 0 ? '↑' : '↓' }}</span>
                {{ fmtMoney(trade.Pnl, currency) }}
              </dd>
            </div>
          </dl>

          <!-- Emotion, rules and open date as three small tiles -->
          <dl class="mt-4 grid grid-cols-2 gap-2 sm:grid-cols-3">
            <div class="rounded-xl bg-bg px-3 py-2.5">
              <dt class="text-xs text-muted">
                Emotion
              </dt>
              <dd class="mt-0.5 truncate text-sm font-medium capitalize">
                {{ trade.Emotion ?? '—' }}
              </dd>
            </div>
            <div class="rounded-xl bg-bg px-3 py-2.5">
              <dt class="text-xs text-muted">
                Rules
              </dt>
              <!-- Broken rules use the warning colour, with an icon and the word -->
              <dd class="mt-0.5 flex items-center gap-1.5 text-sm font-medium" :class="trade.FollowedRules === false ? 'text-warning-text' : ''">
                <svg v-if="trade.FollowedRules === true" class="h-4 w-4 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                  <path d="M5 12l5 5 9-10" />
                </svg>
                <svg v-else-if="trade.FollowedRules === false" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                  <path d="M12 4l9 16H3L12 4zM12 10v4M12 17h.01" />
                </svg>
                {{ trade.FollowedRules === null ? '—' : trade.FollowedRules ? 'Followed' : 'Broken' }}
              </dd>
            </div>
            <div class="col-span-2 rounded-xl bg-bg px-3 py-2.5 sm:col-span-1">
              <dt class="text-xs text-muted">
                Opened
              </dt>
              <dd class="tnum mt-0.5 text-sm font-medium">
                {{ trade.OpenedAt ? fmtDay(trade.OpenedAt) : '—' }}
              </dd>
            </div>
          </dl>

          <div v-if="trade.Notes" class="mt-4 border-t border-border pt-4">
            <p class="text-xs text-muted">
              Notes
            </p>
            <p class="mt-1.5 whitespace-pre-wrap wrap-break-word text-sm leading-relaxed">
              {{ trade.Notes }}
            </p>
          </div>
        </section>

        <div class="space-y-4">
          <!-- Screenshots -->
          <section :class="panel" aria-label="Screenshots">
            <div class="mb-3 flex items-center justify-between">
              <h2 class="text-sm font-semibold">
                Screenshots
              </h2>
              <span class="tnum text-xs text-muted">
                {{ (images ?? []).length }} total
              </span>
            </div>

            <TradesImageGallery v-if="(images ?? []).length > 0" :images="images ?? []" />

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
          <section :class="panel" aria-label="Visibility">
            <label class="flex min-h-12 cursor-pointer items-center justify-between gap-4" :class="toggling ? 'cursor-wait opacity-70' : ''">
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
                role="switch"
                class="peer sr-only"
                :checked="trade.IsPublic"
                :disabled="toggling"
                @change="onToggleVisibility"
              >
              <span
                class="relative h-7 w-12 shrink-0 rounded-full bg-border transition-colors
                peer-checked:bg-primary
                peer-focus-visible:outline-2 peer-focus-visible:outline-offset-2 peer-focus-visible:outline-primary
                after:absolute after:left-1 after:top-1 after:h-5 after:w-5
                after:rounded-full after:bg-white after:shadow-sm after:transition-transform
                peer-checked:after:translate-x-5"
              />
            </label>
          </section>
        </div>
      </div>
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