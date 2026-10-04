<script setup lang="ts">
import { keepPreviousData, useQuery } from '@tanstack/vue-query'
import type { PublicJournal } from '~/types'

definePageMeta({ layout: 'public' })

const title = 'Public journals | Ctemz Journal'
const description = 'Trades that traders chose to share. Results are shown in R, never in money.'
useSeoMeta({ title, description, ogTitle: title, ogDescription: description })

const api = useApi()
const route = useRoute()
const router = useRouter()
const LIMIT = 12

/* ---------- State lives in the URL, so filters and pages can be shared and Back works ---------- */
type ResultFilter = '' | 'win' | 'loss' | 'open'
type DirFilter = '' | 'long' | 'short'

const page = computed(() => {
  const n = Number.parseInt(String(route.query.page ?? '1'), 10)
  return Number.isFinite(n) && n > 1 ? n - 1 : 0
})
const resultFilter = computed<ResultFilter>(() => {
  const v = String(route.query.result ?? '')
  return v === 'win' || v === 'loss' || v === 'open' ? v : ''
})
const dirFilter = computed<DirFilter>(() => {
  const v = String(route.query.dir ?? '')
  return v === 'long' || v === 'short' ? v : ''
})

function setQuery(patch: Record<string, string | number | undefined>, push = false) {
  const q: Record<string, any> = { ...route.query, ...patch }
  for (const k of Object.keys(q)) {
    if (q[k] === undefined || q[k] === '' || (k === 'page' && Number(q[k]) <= 1))
      delete q[k]
  }
  return push ? router.push({ query: q }) : router.replace({ query: q })
}

/* ---------- Data. One extra row tells us whether a next page exists. ---------- */
const { data, isPending, isError, isPlaceholderData, refetch } = useQuery({
  queryKey: computed(() => ['journals', 'list', LIMIT, page.value]),
  queryFn: () => api.get<PublicJournal[]>(`/public/journals?limit=${LIMIT + 1}&offset=${page.value * LIMIT}`),
  staleTime: 60_000,
  placeholderData: keepPreviousData, // keeps the old cards on screen while the next page loads
})

const journals = computed(() => (data.value ?? []).slice(0, LIMIT))
const hasMore = computed(() => (data.value?.length ?? 0) > LIMIT)

/* ---------- Filters work on the page that is loaded ---------- */
const isClosed = (j: PublicJournal) => j.result === 'win' || j.result === 'loss'

const counts = computed(() => {
  const all = journals.value
  return {
    '': all.length,
    'win': all.filter(j => j.result === 'win').length,
    'loss': all.filter(j => j.result === 'loss').length,
    'open': all.filter(j => !isClosed(j)).length,
  }
})

const visible = computed(() =>
  journals.value.filter((j) => {
    const okResult = resultFilter.value === ''
      || (resultFilter.value === 'open' ? !isClosed(j) : j.result === resultFilter.value)
    const okDir = dirFilter.value === '' || j.direction === dirFilter.value
    return okResult && okDir
  }),
)

const filtersOn = computed(() => resultFilter.value !== '' || dirFilter.value !== '')

const resultChips: { value: ResultFilter, label: string, arrow?: string }[] = [
  { value: '', label: 'All' },
  { value: 'win', label: 'Wins', arrow: '↑' },
  { value: 'loss', label: 'Losses', arrow: '↓' },
  { value: 'open', label: 'Open' },
]
const dirChips: { value: DirFilter, label: string, arrow?: string }[] = [
  { value: '', label: 'Any side' },
  { value: 'long', label: 'Long', arrow: '↑' },
  { value: 'short', label: 'Short', arrow: '↓' },
]

/* ---------- Move to the top of the list when the page changes ---------- */
const topRef = ref<HTMLElement | null>(null)
watch(page, () => {
  nextTick(() => {
    const calm = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    topRef.value?.scrollIntoView({ behavior: calm ? 'auto' : 'smooth', block: 'start' })
  })
})

const summary = computed(() =>
  `Showing ${visible.value.length} of ${journals.value.length} journals on this page.`,
)
</script>

<template>
  <div class="mx-auto w-full max-w-6xl flex-1 px-4 pb-16 pt-10 md:px-6 md:pt-16">
    <header>
      <p class="text-xs font-semibold uppercase tracking-wide text-primary">
        Community
      </p>
      <h1 class="mt-1 text-balance text-3xl font-extrabold tracking-tight md:text-4xl">
        Public journals
      </h1>
      <p class="mt-2 max-w-xl text-pretty text-base text-muted">
        Trades that traders chose to share. Results are shown in R, never in money.
      </p>
    </header>

    <!-- Filters stay reachable while you scroll -->
    <div ref="topRef" class="sticky top-16 z-30 -mx-4 mt-6 scroll-mt-16 border-b border-border bg-bg/90 px-4 py-3 backdrop-blur md:-mx-6 md:px-6">
      <div class="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
        <div class="-mx-1 flex gap-1.5 overflow-x-auto px-1 pb-0.5" role="group" aria-label="Filter by result">
          <button
            v-for="c in resultChips"
            :key="c.value"
            type="button"
            :aria-pressed="resultFilter === c.value"
            class="inline-flex h-10 shrink-0 items-center gap-1.5 rounded-full px-4 text-sm font-semibold transition-colors  focus-visible:outline-2-2 focus-visible:outline-2-offset-2 focus-visible:outline-2-primary"
            :class="resultFilter === c.value ? 'bg-primary text-on-primary' : 'border border-border bg-surface text-muted hover:text-text'"
            @click="setQuery({ result: c.value })"
          >
            <span v-if="c.arrow" aria-hidden="true">{{ c.arrow }}</span>
            {{ c.label }}
            <span class="tnum text-xs opacity-75">{{ counts[c.value] }}</span>
          </button>
        </div>

        <div class="-mx-1 flex gap-1.5 overflow-x-auto px-1 pb-0.5" role="group" aria-label="Filter by side">
          <button
            v-for="c in dirChips"
            :key="c.value"
            type="button"
            :aria-pressed="dirFilter === c.value"
            class="inline-flex h-10 shrink-0 items-center gap-1.5 rounded-full px-4 text-sm font-medium transition-colors focus-visible:outline-2 focus-visible:outline-2-2 focus-visible:outline-2-offset-2 focus-visible:outline-2-primary"
            :class="dirFilter === c.value ? 'border border-primary bg-primary/10 text-primary' : 'border border-border bg-surface text-muted hover:text-text'"
            @click="setQuery({ dir: c.value })"
          >
            <span v-if="c.arrow" aria-hidden="true">{{ c.arrow }}</span>
            {{ c.label }}
          </button>
        </div>
      </div>
    </div>

    <p class="mt-3 text-sm text-muted" aria-live="polite">
      <template v-if="!isPending && !isError">
        {{ summary }}
        <span v-if="filtersOn">Filters apply to this page only.</span>
      </template>
      <template v-else>
        &nbsp;
      </template>
    </p>

    <!-- Loading: same grid shape, with borders so skeletons show on the light background -->
    <div v-if="isPending" class="mt-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-3" aria-hidden="true">
      <div v-for="i in 6" :key="i" class="h-64 animate-pulse rounded-2xl border border-border bg-surface" />
    </div>

    <div v-else-if="isError" class="mt-4 rounded-2xl border border-border bg-surface p-8 text-center">
      <p class="text-sm text-muted">
        Journals did not load. Check your connection and try again.
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
      <div
        v-if="visible.length > 0"
        class="mt-4 grid gap-4 transition-opacity sm:grid-cols-2 lg:grid-cols-3"
        :class="isPlaceholderData ? 'opacity-60' : ''"
      >
        <HomeJournalCard v-for="j in visible" :key="j.id" :journal="j" />
      </div>

      <!-- Empty states say what happened and what to do next -->
      <div v-else class="mt-4 rounded-2xl border border-border bg-surface p-10 text-center">
        <template v-if="filtersOn && journals.length > 0">
          <p class="text-sm text-muted">
            Nothing on this page matches those filters.
          </p>
          <button
            type="button"
            class="mt-4 inline-flex h-11 items-center rounded-xl border border-border px-5 text-sm font-semibold transition-colors hover:border-primary"
            @click="setQuery({ result: undefined, dir: undefined })"
          >
            Clear filters
          </button>
        </template>
        <template v-else-if="page > 0">
          <p class="text-sm text-muted">
            There are no more journals after this point.
          </p>
          <button
            type="button"
            class="mt-4 inline-flex h-11 items-center rounded-xl border border-border px-5 text-sm font-semibold transition-colors hover:border-primary"
            @click="setQuery({ page: undefined }, true)"
          >
            Back to the first page
          </button>
        </template>
        <p v-else class="text-sm text-muted">
          No public journals yet. Yours could be the first.
        </p>
      </div>

      <nav v-if="page > 0 || hasMore" class="mt-8 flex items-center justify-between gap-3" aria-label="Pagination">
        <button
          type="button"
          :disabled="page === 0"
          class="inline-flex h-11 items-center rounded-xl border border-border bg-surface px-5 text-sm font-semibold transition-colors hover:border-primary disabled:cursor-not-allowed disabled:opacity-40"
          @click="setQuery({ page: page }, true)"
        >
          Previous
        </button>
        <span class="tnum text-sm text-muted">Page {{ page + 1 }}</span>
        <button
          type="button"
          :disabled="!hasMore"
          class="inline-flex h-11 items-center rounded-xl border border-border bg-surface px-5 text-sm font-semibold transition-colors hover:border-primary disabled:cursor-not-allowed disabled:opacity-40"
          @click="setQuery({ page: page + 2 }, true)"
        >
          Next
        </button>
      </nav>
    </template>
  </div>
</template>