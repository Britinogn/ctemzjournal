<script setup lang="ts">
import { keepPreviousData, useQuery } from '@tanstack/vue-query'
import { toast } from 'vue-sonner'
import type { SupabaseClient } from '@supabase/supabase-js'
import { tradesKey, type Trade, type TradeFilters } from '~/types'
import {
  Search01Icon,
  Download01Icon,
  Loading03Icon,
} from '~/utils/icons'

definePageMeta({ middleware: 'auth', layout: 'dashboard' })

const api = useApi()
const { data: setups } = useSetups()

const search = ref('')
const status = ref('')
const result = ref('')
const setupId = ref('')
const datePreset = ref('30')
const page = ref(0)
const LIMIT = 20

const setupNames = computed<Record<string, string>>(() =>
  Object.fromEntries((setups.value ?? []).map(s => [s.ID, s.Name])),
)

function range(): { from?: string, to?: string } {
  if (datePreset.value === 'all')
    return {}
  const days = Number(datePreset.value)
  const to = new Date()
  const from = new Date(to.getTime() - days * 86_400_000)
  return { from: from.toISOString(), to: to.toISOString() }
}

const filters = computed<TradeFilters>(() => ({
  ...(search.value.includes('/') ? { pair: search.value.trim().toUpperCase() } : {}),
  ...(status.value ? { status: status.value as TradeFilters['status'] } : {}),
  ...(result.value ? { result: result.value as TradeFilters['result'] } : {}),
  ...(setupId.value ? { setup: setupId.value } : {}),
  ...range(),
  limit: LIMIT,
  offset: page.value * LIMIT,
}))

// One extra row is requested, only to learn whether a next page exists.
const { data: raw, isPending, isError, isPlaceholderData, refetch } = useQuery({
  queryKey: computed(() => tradesKey(filters.value)),
  queryFn: () => {
    const q = new URLSearchParams()
    for (const [k, v] of Object.entries({ ...filters.value, limit: LIMIT + 1 })) {
      if (v !== undefined && v !== '')
        q.set(k, String(v))
    }
    const qs = q.toString()
    return api.get<Trade[]>(`/trades${qs ? `?${qs}` : ''}`)
  },
  placeholderData: keepPreviousData, // keeps the old rows on screen while the next page loads
})

const pageTrades = computed(() => (raw.value ?? []).slice(0, LIMIT))
const hasMore = computed(() => (raw.value?.length ?? 0) > LIMIT)

// The server filters on a full pair ("EUR/USD"). Typing part of one ("eur") now narrows the rows on screen instead of doing nothing.
const pairTerm = computed(() => {
  const s = search.value.trim().toUpperCase().replace(/\s+/g, '')
  return s && !s.includes('/') ? s : ''
})
const trades = computed(() =>
  pairTerm.value
    ? pageTrades.value.filter(t => t.Pair.replace('/', '').includes(pairTerm.value))
    : pageTrades.value,
)

const filtersOn = computed(() =>
  !!(search.value.trim() || status.value || result.value || setupId.value || datePreset.value !== '30'),
)
function clearFilters(): void {
  search.value = ''
  status.value = ''
  result.value = ''
  setupId.value = ''
  datePreset.value = '30'
}

function resetPage(): void {
  page.value = 0
}
watch([search, status, result, setupId, datePreset], resetPage)

// Back to the top of the list when the page changes.
const topRef = ref<HTMLElement | null>(null)
watch(page, () => {
  nextTick(() => {
    const calm = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    topRef.value?.scrollIntoView({ behavior: calm ? 'auto' : 'smooth', block: 'start' })
  })
})

const exporting = ref(false)

async function onExport(): Promise<void> {
  if (exporting.value)
    return
  exporting.value = true
  try {
    const { $supabase } = useNuxtApp() as unknown as { $supabase: SupabaseClient }
    const { data: { session } } = await $supabase.auth.getSession()
    const q = new URLSearchParams()
    for (const [k, v] of Object.entries(filters.value)) {
      if (v !== undefined && v !== '' && k !== 'limit' && k !== 'offset')
        q.set(k, String(v))
    }
    const qs = q.toString()
    const blob = await $fetch<Blob>(`/trades/export.csv${qs ? `?${qs}` : ''}`, {
      baseURL: useRuntimeConfig().public.apiUrl as string,
      headers: session?.access_token ? { Authorization: `Bearer ${session.access_token}` } : {},
      responseType: 'blob',
    })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'trades-export.csv'
    a.click()
    URL.revokeObjectURL(url)
    toast.success('Export downloaded')
  }
  catch {
    toast.error('Export failed')
  }
  finally {
    exporting.value = false
  }
}

/* ---------- Shared look ---------- */
const panel = 'rounded-2xl border border-border bg-surface p-4 md:p-5'
const input = 'h-11 w-full rounded-xl border border-border bg-bg px-4 text-sm outline-none transition-colors placeholder:text-muted focus:border-primary focus-visible:ring-2 focus-visible:ring-primary/30'
// appearance-none plus an arrow drawn below, so all four dropdowns look the same in every browser
const select = 'h-11 w-full appearance-none rounded-xl border border-border bg-bg pl-3.5 pr-9 text-sm outline-none transition-colors focus:border-primary focus-visible:ring-2 focus-visible:ring-primary/30'
const pagerBtn = 'inline-flex h-11 items-center rounded-xl border border-border bg-surface px-5 text-sm font-semibold transition-colors hover:border-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary disabled:cursor-not-allowed disabled:opacity-40'
</script>

<template>
  <div>
    <!-- Mobile heading -->
    <div class="mb-4 flex items-center justify-between gap-3 md:justify-end">
      <h1 class="text-xl font-bold tracking-tight md:sr-only">
        Trades
      </h1>
      <NuxtLink
        to="/dashboard/trades/new"
        class="inline-flex h-11 items-center gap-2 rounded-xl bg-primary px-5 text-sm font-semibold text-on-primary transition-colors hover:opacity-90 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
      >
        <UiAppIcon :icon="Add01Icon" :size="18" />
        Log trade
      </NuxtLink>
    </div>

    <!-- Filters -->
    <section :class="[panel, 'mb-4']" aria-label="Filters">
      <div class="grid grid-cols-2 gap-3 md:flex md:flex-wrap md:items-center">
        <!-- Search -->
        <label class="relative col-span-2 md:w-64">
          <span class="sr-only">Search pair</span>
          <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5 text-muted">
            <UiAppIcon :icon="Search01Icon" :size="16" />
          </div>
          <input
            v-model="search"
            type="search"
            autocomplete="off"
            placeholder="Search pair (e.g. EUR/USD)"
            :class="[input, 'pl-10']"
          >
        </label>

        <div class="relative md:w-40">
          <select v-model="status" aria-label="Status" :class="select">
            <option value="">All status</option>
            <option value="open">Open</option>
            <option value="closed">Closed</option>
          </select>
          <svg class="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M6 9l6 6 6-6" /></svg>
        </div>

        <div class="relative md:w-40">
          <select v-model="result" aria-label="Result" :class="select">
            <option value="">Any result</option>
            <option value="win">Win</option>
            <option value="loss">Loss</option>
          </select>
          <svg class="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M6 9l6 6 6-6" /></svg>
        </div>

        <div class="relative md:w-44">
          <select v-model="setupId" aria-label="Setup" :class="select">
            <option value="">All setups</option>
            <option v-for="s in setups ?? []" :key="s.ID" :value="s.ID">
              {{ s.Name }}
            </option>
          </select>
          <svg class="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M6 9l6 6 6-6" /></svg>
        </div>

        <div class="relative md:w-44">
          <select v-model="datePreset" aria-label="Date range" :class="select">
            <option value="7">Last 7 days</option>
            <option value="30">Last 30 days</option>
            <option value="90">Last 90 days</option>
            <option value="all">All time</option>
          </select>
          <svg class="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M6 9l6 6 6-6" /></svg>
        </div>

        <!-- Only shown when something is filtered -->
        <button
          v-if="filtersOn"
          type="button"
          class="inline-flex h-11 items-center justify-center rounded-xl px-3 text-sm font-semibold text-primary transition-colors hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
          @click="clearFilters"
        >
          Clear filters
        </button>

        <button
          type="button"
          :disabled="exporting"
          :aria-busy="exporting"
          class="col-span-2 inline-flex h-11 items-center justify-center gap-2 rounded-xl border border-border bg-bg px-4 text-sm font-semibold transition-colors hover:border-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary disabled:cursor-wait disabled:opacity-60 md:col-span-1 md:ml-auto"
          @click="onExport"
        >
          <UiAppIcon
            v-if="exporting"
            :icon="Loading03Icon"
            :size="16"
            class="animate-spin motion-reduce:animate-none"
          />
          <UiAppIcon
            v-else
            :icon="Download01Icon"
            :size="16"
          />
          {{ exporting ? 'Exporting…' : 'Export CSV' }}
        </button>
      </div>
    </section>

    <!-- Loading -->
    <div v-if="isPending" class="space-y-2" aria-hidden="true">
      <div v-for="i in 5" :key="i" class="h-16 animate-pulse rounded-xl border border-border bg-surface" />
    </div>

    <!-- Error -->
    <div
      v-else-if="isError"
      :class="[panel, 'text-center']"
    >
      <p class="text-sm text-muted">
        Couldn't load trades.
      </p>
      <button
        type="button"
        class="mt-4 inline-flex h-11 items-center rounded-xl bg-primary px-6 text-sm font-semibold text-on-primary transition-colors hover:opacity-90 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
        @click="() => refetch()"
      >
        Retry
      </button>
    </div>

    <!-- Table. While the next page loads, the old rows stay and fade. -->
    <section
      v-else
      ref="topRef"
      :class="[panel, 'scroll-mt-20 transition-opacity', isPlaceholderData ? 'opacity-60' : '']"
      aria-label="Trades"
      :aria-busy="isPlaceholderData"
    >
      <!-- Filters that match nothing get their own message and a way out -->
      <div v-if="trades.length === 0 && filtersOn" class="flex flex-col items-center px-4 py-10 text-center">
        <p class="text-sm text-muted">
          No trades match these filters.
        </p>
        <button
          type="button"
          class="mt-4 inline-flex h-11 items-center rounded-xl border border-border bg-surface px-5 text-sm font-semibold transition-colors hover:border-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
          @click="clearFilters"
        >
          Clear filters
        </button>
      </div>
      <TradesRecentTrades
        v-else
        :trades="trades"
        :setup-names="setupNames"
      />

      <!-- Pagination -->
      <div class="mt-4 flex flex-wrap items-center justify-between gap-3 border-t border-border pt-4 text-sm">
        <p class="text-muted" aria-live="polite">
          Showing
          <span class="tnum">
            {{ pageTrades.length === 0 ? 0 : page * LIMIT + 1 }}
          </span>
          to
          <span class="tnum">
            {{ page * LIMIT + pageTrades.length }}
          </span>
        </p>
        <div class="flex items-center gap-2">
          <button
            type="button"
            :disabled="page === 0"
            :class="pagerBtn"
            @click="page--"
          >
            Previous
          </button>
          <span class="tnum px-1 text-muted">Page {{ page + 1 }}</span>
          <button
            type="button"
            :disabled="!hasMore"
            :class="pagerBtn"
            @click="page++"
          >
            Next
          </button>
        </div>
      </div>
    </section>
  </div>
</template>