<script setup lang="ts">
import { useIntervalFn, useLocalStorage } from '@vueuse/core'
import { useQuery } from '@tanstack/vue-query'
import { ratesKey, type RatesResponse } from '~/types'
import { timeAgo } from '~/utils/format'

definePageMeta({ layout: 'public' })

usePageSeo(
  'Live Forex Prices',
  'Check live prices for the pairs you trade, then log your trades and track your performance in Ctemz Journal.',
)
// const description = 'Indicative forex, gold and naira prices for journal context. Refreshed every minute.'
// useSeoMeta({ title, description, ogTitle: title, ogDescription: description })

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

/* ---------- Pinned pairs: the ones you trade, kept on this device and shown first ---------- */
// initOnMounted keeps the server and the browser in agreement on the first paint.
const pinned = useLocalStorage<string[]>('ctemz:pinned-pairs', [], { initOnMounted: true })
const isPinned = (pair: string): boolean => pinned.value.includes(pair)
function togglePin(pair: string): void {
  pinned.value = isPinned(pair) ? pinned.value.filter(p => p !== pair) : [...pinned.value, pair]
}

/* ---------- Movement since the last update ---------- */
// The API sends one price per pair, so the change is worked out here, between one refresh and the next.
const moved = ref<Record<string, number>>({})
watch(rates, (now, before) => {
  if (!now || !before)
    return
  const next = { ...moved.value }
  for (const r of now.rates) {
    const old = before.rates.find(x => x.pair === r.pair)
    if (!old)
      continue
    const a = Number(old.price)
    const b = Number(r.price)
    if (Number.isFinite(a) && Number.isFinite(b) && a !== b)
      next[r.pair] = b - a
  }
  moved.value = next
})
const decimalsOf = (price: string | number): number => String(price).split('.')[1]?.length ?? 0

/* ---------- Copy a price, ready to paste into a trade ---------- */
const copied = ref('')
async function copyPrice(pair: string, price: string | number): Promise<void> {
  try {
    await navigator.clipboard.writeText(String(price))
    copied.value = pair
    setTimeout(() => {
      if (copied.value === pair)
        copied.value = ''
    }, 1500)
  }
  catch {
    // Clipboard is not available (an old browser or an insecure page). The price is still on screen.
  }
}

/* ---------- Search and groups ---------- */
type Cat = 'all' | 'pinned' | 'forex' | 'metals' | 'naira'
const q = ref('')
const cat = ref<Cat>('all')

function categoryOf(pair: string): 'forex' | 'metals' | 'naira' {
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
    ...(pinned.value.length > 0 ? [{ value: 'pinned' as const, label: 'Pinned' }] : []),
    { value: 'forex', label: 'Forex' },
    { value: 'metals', label: 'Gold and metals' },
    { value: 'naira', label: 'Naira' },
  ]
  return list.filter(c => c.value === 'all' || c.value === 'pinned' || present.has(c.value))
})

// Removing the last pin while looking at Pinned would leave an empty page, so go back to All.
watch(() => pinned.value.length, (n) => {
  if (n === 0 && cat.value === 'pinned')
    cat.value = 'all'
})

const visible = computed(() => {
  const term = q.value.trim().toUpperCase().replace(/\s+/g, '')
  const list = all.value.filter((r) => {
    const okCat = cat.value === 'all'
      || (cat.value === 'pinned' ? isPinned(r.pair) : categoryOf(r.pair) === cat.value)
    const okText = !term || r.pair.replace('/', '').includes(term.replace('/', ''))
    return okCat && okText
  })
  // Pinned pairs first. The sort is stable, so the rest keep the order the API sent.
  return [...list].sort((a, b) => Number(isPinned(b.pair)) - Number(isPinned(a.pair)))
})

function clear() {
  q.value = ''
  cat.value = 'all'
}

const focusRing
  = 'focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary'
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
          :class="['inline-flex h-11 items-center gap-2 rounded-xl border border-border bg-surface px-4 font-semibold text-text transition-colors hover:border-primary disabled:opacity-60', focusRing]"
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
          :class="[
            'inline-flex h-10 shrink-0 items-center gap-1.5 rounded-full px-4 text-sm font-semibold transition-colors',
            focusRing,
            cat === c.value ? 'bg-primary text-on-primary' : 'border border-border bg-surface text-muted hover:text-text',
          ]"
          @click="cat = c.value"
        >
          <svg v-if="c.value === 'pinned'" class="h-3.5 w-3.5" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
            <path d="M12 3l2.6 5.6 6.1.7-4.5 4.2 1.2 6L12 16.5 6.6 19.5l1.2-6L3.3 9.3l6.1-.7L12 3z" />
          </svg>
          {{ c.label }}
        </button>
      </div>
    </div>

    <!-- Loading: same shape as the real grid -->
    <div v-if="isPending" class="mt-5 grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4" aria-hidden="true">
      <div v-for="i in 8" :key="i" class="h-28 animate-pulse rounded-2xl border border-border bg-surface" />
    </div>

    <div v-else-if="isError" class="mt-5 rounded-2xl border border-border bg-surface p-8 text-center">
      <p class="text-sm text-muted">
        Prices did not load. Check your connection and try again.
      </p>
      <button
        type="button"
        :class="['mt-4 inline-flex h-11 items-center rounded-xl bg-primary px-6 text-sm font-semibold text-on-primary transition-colors hover:opacity-90', focusRing]"
        @click="() => refetch()"
      >
        Try again
      </button>
    </div>

    <template v-else>
      <p class="mt-4 text-sm text-muted" aria-live="polite">
        {{ visible.length }} {{ visible.length === 1 ? 'pair' : 'pairs' }}
      </p>
      <!-- Says "Copied" to screen readers, since the visible change is only a small icon -->
      <p class="sr-only" role="status">
        {{ copied ? `Copied ${copied} price` : '' }}
      </p>

      <dl v-if="visible.length > 0" class="mt-2 grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
        <div v-for="r in visible" :key="r.pair" class="relative rounded-2xl border border-border bg-surface px-4 py-3.5">
          <dt class="tnum pr-8 text-xs font-medium text-muted">
            {{ r.pair }}
          </dt>

          <!-- Tap the price to copy it, for example to paste into the entry field of a trade -->
          <dd class="mt-0.5">
            <button
              type="button"
              :class="['-mx-1.5 inline-flex min-h-11 items-center gap-2 rounded-lg px-1.5 text-left', focusRing]"
              :title="`Copy ${r.pair} price`"
              :aria-label="`Copy ${r.pair} price, ${r.price}`"
              @click="copyPrice(r.pair, r.price)"
            >
              <span class="tnum text-2xl font-semibold" :class="r.stale ? 'text-warning-text' : ''">{{ r.price }}</span>
              <svg v-if="copied === r.pair" class="h-4 w-4 shrink-0 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                <path d="M5 12l5 5 9-10" />
              </svg>
              <svg v-else class="h-4 w-4 shrink-0 text-muted/70" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                <rect x="9" y="9" width="11" height="11" rx="2" /><path d="M5 15V6a2 2 0 012-2h9" />
              </svg>
            </button>
          </dd>

          <dd class="flex flex-wrap items-center gap-x-2 text-xs" :class="r.stale ? 'font-medium text-warning-text' : 'text-muted'">
            <span class="inline-flex items-center gap-1">
              <svg v-if="r.stale" class="h-3 w-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                <path d="M12 4l9 16H3L12 4zM12 10v4M12 17h.01" />
              </svg>
              {{ r.stale ? 'Out of date, ' : '' }}{{ ago(r.updated_at) }}
            </span>
            <!-- How far it moved at the last refresh. An arrow and a number, with no green or red. -->
            <span v-if="moved[r.pair] !== undefined" class="tnum inline-flex items-center gap-0.5 font-medium text-text" title="Change at the last update">
              <span aria-hidden="true">{{ moved[r.pair]! > 0 ? '↑' : '↓' }}</span>
              <span class="sr-only">{{ moved[r.pair]! > 0 ? 'Up' : 'Down' }}</span>
              {{ Math.abs(moved[r.pair]!).toFixed(decimalsOf(r.price)) }}
            </span>
          </dd>

          <!-- Pin: the pairs you care about go to the top and stay there on this device -->
          <dd class="absolute right-1 top-1">
            <button
              type="button"
              :aria-pressed="isPinned(r.pair)"
              :aria-label="`Pin ${r.pair}`"
              :title="isPinned(r.pair) ? 'Unpin' : 'Pin to the top'"
              :class="['inline-flex h-11 w-11 items-center justify-center rounded-xl transition-colors', focusRing, isPinned(r.pair) ? 'text-primary' : 'text-muted/60 hover:text-text']"
              @click="togglePin(r.pair)"
            >
              <svg class="h-4.5 w-4.5" viewBox="0 0 24 24" :fill="isPinned(r.pair) ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round" aria-hidden="true">
                <path d="M12 3l2.6 5.6 6.1.7-4.5 4.2 1.2 6L12 16.5 6.6 19.5l1.2-6L3.3 9.3l6.1-.7L12 3z" />
              </svg>
            </button>
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
          :class="['mt-4 inline-flex h-11 items-center rounded-xl border border-border px-5 text-sm font-semibold transition-colors hover:border-primary', focusRing]"
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