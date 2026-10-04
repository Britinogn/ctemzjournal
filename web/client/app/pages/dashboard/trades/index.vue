<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
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

function range(): { from?: string; to?: string } {
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

const { data: trades, isPending, isError, refetch } = useQuery({
  queryKey: computed(() => tradesKey(filters.value)),
  queryFn: () => {
    const q = new URLSearchParams()
    for (const [k, v] of Object.entries(filters.value)) {
      if (v !== undefined && v !== '')
        q.set(k, String(v))
    }
    const qs = q.toString()
    return api.get<Trade[]>(`/trades${qs ? `?${qs}` : ''}`)
  },
})

function resetPage(): void {
  page.value = 0
}
watch([search, status, result, setupId, datePreset], resetPage)

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
const select = 'h-11 rounded-xl border border-border bg-bg px-3 text-sm outline-none transition-colors focus:border-primary focus-visible:ring-2 focus-visible:ring-primary/30'
</script>

<template>
  <div>
    <!-- Mobile heading -->
    <h1 class="mb-4 text-xl font-bold tracking-tight md:sr-only">
      Trades
    </h1>

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
            placeholder="Search pair (e.g. EUR/USD)"
            :class="[input, 'pl-10']"
          >
        </label>

        <select v-model="status" aria-label="Status" :class="select">
          <option value="">All status</option>
          <option value="open">Open</option>
          <option value="closed">Closed</option>
        </select>

        <select v-model="result" aria-label="Result" :class="select">
          <option value="">Any result</option>
          <option value="win">Win</option>
          <option value="loss">Loss</option>
        </select>

        <select v-model="setupId" aria-label="Setup" :class="select">
          <option value="">All setups</option>
          <option v-for="s in setups ?? []" :key="s.ID" :value="s.ID">
            {{ s.Name }}
          </option>
        </select>

        <select v-model="datePreset" aria-label="Date range" :class="select">
          <option value="7">Last 7 days</option>
          <option value="30">Last 30 days</option>
          <option value="90">Last 90 days</option>
          <option value="all">All time</option>
        </select>

        <button
          type="button"
          :disabled="exporting"
          class="col-span-2 inline-flex h-11 items-center justify-center gap-2 rounded-xl border border-border bg-bg px-4 text-sm font-semibold transition hover:border-primary disabled:opacity-50 md:col-span-1"
          @click="onExport"
        >
          <UiAppIcon
            v-if="exporting"
            :icon="Loading03Icon"
            :size="16"
            class="animate-spin"
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
        class="mt-3 inline-flex h-11 items-center rounded-xl bg-primary px-5 text-sm font-semibold text-on-primary transition hover:opacity-90"
        @click="() => refetch()"
      >
        Retry
      </button>
    </div>

    <!-- Table -->
    <section
      v-else
      :class="panel"
      aria-label="Trades"
    >
      <TradesRecentTrades
        :trades="trades ?? []"
        :setup-names="setupNames"
      />

      <!-- Pagination -->
      <div class="mt-4 flex items-center justify-between border-t border-border pt-4 text-sm">
        <p class="text-muted">
          Showing
          <span class="tnum">
            {{ (trades ?? []).length === 0 ? 0 : page * LIMIT + 1 }}
          </span>
          to
          <span class="tnum">
            {{ page * LIMIT + (trades ?? []).length }}
          </span>
        </p>
        <div class="flex gap-2">
          <button
            type="button"
            :disabled="page === 0"
            class="inline-flex h-9 items-center rounded-xl border border-border px-4 text-sm font-medium transition hover:border-primary disabled:opacity-40"
            @click="page--"
          >
            Previous
          </button>
          <button
            type="button"
            :disabled="(trades ?? []).length < LIMIT"
            class="inline-flex h-9 items-center rounded-xl border border-border px-4 text-sm font-medium transition hover:border-primary disabled:opacity-40"
            @click="page++"
          >
            Next
          </button>
        </div>
      </div>
    </section>
  </div>
</template>