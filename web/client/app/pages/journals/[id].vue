<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import type { PublicJournal } from '~/types'
import { fmtR } from '~/utils/format'

definePageMeta({ layout: 'public' })

// The notes and the full picture list are not on your PublicJournal type yet. The API has to send them (see the reply).
interface JournalImage {
  url: string
  kind?: string | null
}
type JournalDetail = PublicJournal & {
  notes?: string | null
  images?: JournalImage[] | null
}

const route = useRoute()
const id = computed(() => String(route.params.id ?? ''))
const api = useApi()

const { data: journal, isPending, isError } = useQuery({
  queryKey: computed(() => ['journals', 'detail', id.value]),
  queryFn: async (): Promise<JournalDetail | null> => {
    try {
      return await api.get<JournalDetail>(`/public/journals/${id.value}`)
    }
    catch {
      // Falls back to the list, which is how this page worked before. It only finds the 100 newest journals.
      const list = await api.get<JournalDetail[]>('/public/journals?limit=100')
      return list.find(j => j.id === id.value) ?? null
    }
  },
  staleTime: 60_000,
})

// Every picture there is. If the API only sends the first one, that single picture is shown.
const pictures = computed(() => {
  const j = journal.value
  if (!j)
    return []
  const list = (j.images ?? []).filter(i => i?.url)
  const items = list.length > 0 ? list : j.image_url ? [{ url: j.image_url, kind: null }] : []
  return items.map((img, i) => ({ ID: String(i), URL: img.url, Kind: img.kind ?? null }))
})

useSeoMeta({
  title: computed(() => journal.value ? `${journal.value.pair} journal — Ctemz Journal` : 'Journal — Ctemz Journal'),
})
</script>

<template>
  <div class="mx-auto w-full max-w-3xl flex-1 px-4 py-8 md:px-6">
    <NuxtLink
      to="/journals"
      class="-ml-1 inline-flex min-h-11 items-center rounded-lg px-1 text-sm font-medium text-primary hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
    >
      ← All journals
    </NuxtLink>

    <div v-if="isPending" class="mt-4 h-96 animate-pulse rounded-2xl border border-border bg-surface" aria-hidden="true" />
    <div v-else-if="isError" class="mt-4 rounded-2xl border border-border bg-surface p-8 text-center text-sm text-muted">
      Couldn't load this journal.
    </div>
    <div v-else-if="!journal" class="mt-4 rounded-2xl border border-border bg-surface p-8 text-center">
      <p class="text-sm font-semibold">
        Journal not found
      </p>
      <p class="mt-1 text-sm text-muted">
        It may have been made private or hidden.
      </p>
      <NuxtLink
        to="/journals"
        class="mt-3 inline-flex min-h-11 items-center text-sm font-medium text-primary hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
      >
        Back to journals
      </NuxtLink>
    </div>

    <article v-else class="mt-4 overflow-hidden rounded-2xl border border-border bg-surface">
      <!--
        All the pictures. Each opens in the viewer, which has arrows, swipe and the arrow keys.
        "contain" keeps the whole chart in view instead of cropping it.
      -->
      <div v-if="pictures.length > 0" class="border-b border-border bg-bg p-3 md:p-4">
        <TradesImageGallery
          :images="pictures"
          fit="contain"
          :grid-class="pictures.length === 1 ? 'grid-cols-1' : 'grid-cols-2'"
        />
      </div>

      <div class="p-4 md:p-6">
        <div class="flex flex-wrap items-center gap-2">
          <h1 class="tnum text-2xl font-bold tracking-tight">
            {{ journal.pair }}
          </h1>
          <span class="rounded-full border border-border px-2.5 py-0.5 text-xs capitalize text-muted">
            {{ journal.direction === 'long' ? '↑ Long' : '↓ Short' }}{{ journal.timeframe ? ` · ${journal.timeframe}` : '' }}
          </span>
          <span
            v-if="journal.result && journal.r_multiple !== null"
            class="tnum rounded-full px-2.5 py-0.5 text-xs font-semibold"
            :class="journal.result === 'win' ? 'bg-profit/10 text-profit-text' : 'bg-loss/10 text-loss'"
          >
            {{ journal.result === 'win' ? '↑ Win' : '↓ Loss' }} {{ fmtR(journal.r_multiple) }}
          </span>
        </div>

        <dl class="mt-4 grid grid-cols-2 gap-3 text-sm sm:grid-cols-4">
          <div>
            <dt class="text-xs text-muted">
              Setup
            </dt>
            <dd class="font-medium">
              {{ journal.setup_name || 'No setup' }}
            </dd>
          </div>
          <div>
            <dt class="text-xs text-muted">
              Trader
            </dt>
            <dd class="font-medium">
              {{ journal.display_name || 'Trader' }}
            </dd>
          </div>
          <div>
            <dt class="text-xs text-muted">
              Date
            </dt>
            <dd class="tnum font-medium">
              {{ journal.date.slice(0, 10) }}
            </dd>
          </div>
          <div>
            <dt class="text-xs text-muted">
              Result
            </dt>
            <dd class="font-medium capitalize">
              {{ journal.result ?? 'Open' }}
            </dd>
          </div>
        </dl>

        <!-- The trader's notes, only when there are some -->
        <section v-if="journal.notes?.trim()" class="mt-5 border-t border-border pt-4" aria-label="Notes">
          <h2 class="text-sm font-semibold">
            Notes
          </h2>
          <p class="mt-1.5 whitespace-pre-wrap wrap-break-word text-sm leading-relaxed">
            {{ journal.notes }}
          </p>
        </section>

        <p class="mt-5 border-t border-border pt-3 text-xs text-muted">
          Shared publicly. Money and lots are never shown.
        </p>
      </div>
    </article>
  </div>
</template>