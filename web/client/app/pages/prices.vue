<script setup lang="ts">
import { useIntervalFn } from '@vueuse/core'
import { useQuery } from '@tanstack/vue-query'
import { ratesKey, type RatesResponse } from '~/types'
import { timeAgo } from '~/utils/format'

definePageMeta({ layout: 'public' })

const title = 'Live prices | Ctemz Journal'
const description = 'Indicative forex, gold and naira prices for journal context. Refreshed every minute.'
useSeoMeta({ title, description, ogTitle: title, ogDescription: description })

const api = useApi()

const { data: rates, isPending, isError, isFetching, refetch } = useQuery({
  queryKey: ratesKey(),
  queryFn: () => api.get<RatesResponse>('/public/rates'),
  refetchInterval: 60_000,
})

/* "Updated 3 min ago" should keep counting even when no new data arrives */
const tick = ref(0)
useIntervalFn(() => { tick.value++ }, 30_000)
function ago(date?: string | null): string {
  void tick.value // re-runs the label every 30 seconds
  return date ? timeAgo(date) : ''
}

/* ---------- Search and groups ---------- */
type Cat = 'all' | 'forex' | 'metals' | 'naira'
const q = ref('')
const cat = ref<Cat>('all')

function categoryOf(pair: string): Exclude<Cat, 'all'> {
  const p = pair.toUpperCase()
  if (/^XA[UG]/.test(p)) return 'metals'
  if (p.includes('NGN')) return 'naira'
  return 'forex'
}

const all = computed(() => rates.value?.rates ?? [])

const chips = computed(() => {
  const present = new Set(all.value.map(r => categoryOf(r.pair)))
  const list: { value: Cat, label: string }[] = [
    { value: 'all', label: 'All' },
    { value: 'forex', label: 'Forex' },
    { value: 'metals', label: 'Gold and metals' },
    { value: 'naira', label: 'Naira' },
  ]
  return list.filter(c => c.value === 'all' || present.has(c.value))
})

const visible = computed(() => {
  const term = q.value.trim().toUpperCase().replace(/\s+/g, '')
  return all.value.filter((r) => {
    const okCat = cat.value === 'all' || categoryOf(r.pair) === cat.value
    const okText = !term || r.pair.replace('/', '').includes(term.replace('/', ''))
    return okCat && okText
  })
})

function clear() {
  q.value = ''
  cat.value = 'all'
}
</script>

<template>
  <div class="mx-auto w-full max-w-6xl flex-1 px-4 pb-16 pt-10 md:px-6 md:pt-16">
    <header class="flex flex-wrap items-end justify-between gap-x-6 gap-y-3">
      <div>
        <p class="text-xs font-semibold uppercase tracking-wide text-primary">
          Market data
        </p>
        <h1 class="mt-1 text-3xl font-extrabold tracking-tight md:text-4xl">
          Live prices
        </h1>
        <p class="mt-2 max-w-xl text-pretty text-base text-muted">
          Indicative prices for journal context, not trading advice.
        </p>
      </div>

      <div class="flex items-center gap-3 text-sm text-muted">
        <span v-if="rates?.updated_at" class="inline-flex items-center gap-1.5">
          <span class="h-1.5 w-1.5 rounded-full" :class="rates.stale ? 'bg-warning' : 'bg-primary'" aria-hidden="true" />
          Updated {{ ago(rates.updated_at) }}
        </span>
        <button
          type="button"
          class="inline-flex h-10 items-center gap-2 rounded-xl border border-border bg-surface px-4 font-semibold text-text transition-colors hover:border-primary disabled:opacity-60"
          :disabled="isFetching"
          :aria-busy="isFetching"
          @click="() => refetch()"
        >
          <svg class="h-4 w-4" :class="isFetching ? 'animate-spin motion-reduce:animate-none' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M20 12a8 8 0 11-2.3-5.6M20 4v5h-5" />
          </svg>
          {{ isFetching ? 'Refreshing' : 'Refresh' }}
        </button>
      </div>
    </header>

    <!-- A page-level warning, in words and an icon, when the provider failed -->
    <div
      v-if="rates?.stale"
      class="mt-5 flex items-start gap-3 rounded-xl border border-warning/40 bg-warning/10 px-4 py-3 text-sm text-warning-text"
      role="status"
    >
      <svg class="mt-0.5 h-4 w-4 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
        <path d="M12 4l9 16H3L12 4zM12 10v4M12 17h.01" />
      </svg>
      <p>
        <span class="font-semibold">These prices may be out of date.</span>
        The last update from the price provider failed, so you are seeing the most recent prices we have.
      </p>
    </div>

    <!-- Search and groups -->
    <div class="mt-6 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div class="relative md:w-72">
        <label for="prices-filter" class="sr-only">Search pairs</label>
        <svg class="pointer-events-none absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" aria-hidden="true">
          <circle cx="11" cy="11" r="7" /><path d="M20 20l-4-4" />
        </svg>
        <input
          id="prices-filter"
          v-model="q"
          type="search"
          inputmode="search"
          autocomplete="off"
          placeholder="Search pairs, for example EUR"
          class="h-11 w-full rounded-xl border border-border bg-surface pl-10 pr-4 text-sm outline-none transition-colors placeholder:text-muted focus:border-primary focus-visible:ring-2 focus-visible:ring-primary/30"
          @keydown.esc="q = ''"
        >
      </div>

      <div v-if="chips.length > 2" class="-mx-1 flex gap-1.5 overflow-x-auto px-1 pb-0.5" role="group" aria-label="Filter by group">
        <button
          v-for="c in chips"
          :key="c.value"
          type="button"
          :aria-pressed="cat === c.value"
          class="inline-flex h-10 shrink-0 items-center rounded-full px-4 text-sm font-semibold transition-colors focus-visible:outline-2-2 focus-visible:outline-2-offset-2 focus-visible:outline-2-primary"
          :class="cat === c.value ? 'bg-primary text-on-primary' : 'border border-border bg-surface text-muted hover:text-text'"
          @click="cat = c.value"
        >
          {{ c.label }}
        </button>
      </div>
    </div>

    <!-- Loading: same shape as the real grid -->
    <div v-if="isPending" class="mt-5 grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4" aria-hidden="true">
      <div v-for="i in 8" :key="i" class="h-23 animate-pulse rounded-2xl border border-border bg-surface" />
    </div>

    <div v-else-if="isError" class="mt-5 rounded-2xl border border-border bg-surface p-8 text-center">
      <p class="text-sm text-muted">
        Prices did not load. Check your connection and try again.
      </p>
      <button
        type="button"
        class="mt-4 inline-flex h-11 items-center rounded-xl bg-primary px-6 text-sm font-semibold text-on-primary transition-colors hover:opacity-90"
        @click="() => refetch()"
      >
        Try again
      </button>
    </div>

    <template v-else>
      <p class="mt-4 text-sm text-muted" aria-live="polite">
        {{ visible.length }} {{ visible.length === 1 ? 'pair' : 'pairs' }}
      </p>

      <dl v-if="visible.length > 0" class="mt-2 grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
        <div v-for="r in visible" :key="r.pair" class="rounded-2xl border border-border bg-surface px-4 py-3.5">
          <dt class="tnum text-xs font-medium text-muted">
            {{ r.pair }}
          </dt>
          <dd class="tnum mt-0.5 text-2xl font-semibold" :class="r.stale ? 'text-warning-text' : ''">
            {{ r.price }}
          </dd>
          <dd class="mt-1 flex items-center gap-1 text-xs" :class="r.stale ? 'font-medium text-warning-text' : 'text-muted'">
            <svg v-if="r.stale" class="h-3 w-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
              <path d="M12 4l9 16H3L12 4zM12 10v4M12 17h.01" />
            </svg>
            {{ r.stale ? 'Out of date, ' : '' }}{{ ago(r.updated_at) }}
          </dd>
        </div>
      </dl>

      <div v-else class="mt-2 rounded-2xl border border-border bg-surface p-10 text-center">
        <p class="text-sm text-muted">
          <template v-if="all.length === 0">
            No prices are available yet. They refresh every minute, so check back shortly.
          </template>
          <template v-else>
            No pairs match "{{ q }}" in this group.
          </template>
        </p>
        <button
          v-if="all.length > 0"
          type="button"
          class="mt-4 inline-flex h-11 items-center rounded-xl border border-border px-5 text-sm font-semibold transition-colors hover:border-primary"
          @click="clear"
        >
          Clear search
        </button>
      </div>
    </template>

    <p class="mt-8 text-xs text-muted">
      Prices refresh every minute while this page is open.
    </p>
  </div>
</template>