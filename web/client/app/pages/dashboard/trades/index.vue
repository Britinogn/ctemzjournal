<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { toast } from 'vue-sonner'
import type { SupabaseClient } from '@supabase/supabase-js'
import { tradesKey, type Trade, type TradeFilters } from '~/types'

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
</script>

<template>
  <div>
    <h1 class="mb-4 text-xl font-bold tracking-tight md:text-2xl">
      Trades
    </h1>

    <!-- Filters -->
    <div class="mb-4 grid grid-cols-2 gap-2 md:flex md:flex-wrap md:items-center">
      <label class="relative col-span-2 md:w-64">
        <span class="sr-only">Search pair</span>
        <input
          v-model="search"
          type="search"
          placeholder="Search pair (e.g. EUR/USD)"
          class="w-full rounded-xl border border-border bg-surface py-2.5 pl-4 pr-4 text-sm outline-none transition placeholder:text-muted focus:border-primary"
        >
      </label>
      <select v-model="status" aria-label="Status" class="rounded-xl border border-border bg-surface px-3 py-2.5 text-sm outline-none transition focus:border-primary">
        <option value="">All status</option>
        <option value="open">Open</option>
        <option value="closed">Closed</option>
      </select>
      <select v-model="result" aria-label="Result" class="rounded-xl border border-border bg-surface px-3 py-2.5 text-sm outline-none transition focus:border-primary">
        <option value="">Any result</option>
        <option value="win">Win</option>
        <option value="loss">Loss</option>
      </select>
      <select v-model="setupId" aria-label="Setup" class="rounded-xl border border-border bg-surface px-3 py-2.5 text-sm outline-none transition focus:border-primary">
        <option value="">All setups</option>
        <option v-for="s in setups ?? []" :key="s.ID" :value="s.ID">
          {{ s.Name }}
        </option>
      </select>
      <select v-model="datePreset" aria-label="Date range" class="rounded-xl border border-border bg-surface px-3 py-2.5 text-sm outline-none transition focus:border-primary">
        <option value="7">Last 7 days</option>
        <option value="30">Last 30 days</option>
        <option value="90">Last 90 days</option>
        <option value="all">All time</option>
      </select>
      <button
        type="button"
        :disabled="exporting"
        class="col-span-2 inline-flex items-center justify-center gap-2 rounded-xl border border-border bg-surface px-4 py-2.5 text-sm font-semibold transition hover:border-primary disabled:opacity-50 md:col-span-1"
        @click="onExport"
      >
        {{ exporting ? 'Exporting…' : 'Export CSV' }}
      </button>
    </div>

    <!-- Loading / error -->
    <div v-if="isPending" class="h-96 animate-pulse rounded-2xl bg-surface" />
    <div v-else-if="isError" class="rounded-2xl border border-border bg-surface p-8 text-center">
      <p class="text-sm text-muted">
        Couldn't load trades.
      </p>
      <button
        type="button"
        class="mt-3 rounded-xl bg-primary px-5 py-2 text-sm font-semibold text-on-primary transition hover:opacity-90"
        @click="() => refetch()"
      >
        Retry
      </button>
    </div>

    <!-- Table card -->
    <section v-else class="rounded-2xl border border-border bg-surface p-4" aria-label="Trades">
      <TradesRecentTrades :trades="trades ?? []" :setup-names="setupNames" />
      <div class="mt-3 flex items-center justify-between border-t border-border pt-3 text-sm">
        <p class="text-muted">
          Showing {{ (trades ?? []).length === 0 ? 0 : page * LIMIT + 1 }} to {{ page * LIMIT + (trades ?? []).length }}
        </p>
        <div class="flex gap-2">
          <button
            type="button"
            :disabled="page === 0"
            class="rounded-xl border border-border px-4 py-1.5 text-sm font-medium transition hover:border-primary disabled:opacity-40"
            @click="page--"
          >
            Previous
          </button>
          <button
            type="button"
            :disabled="(trades ?? []).length < LIMIT"
            class="rounded-xl border border-border px-4 py-1.5 text-sm font-medium transition hover:border-primary disabled:opacity-40"
            @click="page++"
          >
            Next
          </button>
        </div>
      </div>
    </section>
  </div>
</template>
