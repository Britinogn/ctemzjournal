<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import type { PublicJournal } from '~/types'
import { fmtR } from '~/utils/format'

definePageMeta({ layout: 'public' })

const route = useRoute()
const id = computed(() => String(route.params.id ?? ''))
const api = useApi()

const { data: journals, isPending, isError } = useQuery({
  queryKey: ['journals', 'detail-lookup'],
  queryFn: () => api.get<PublicJournal[]>('/public/journals?limit=100'),
  staleTime: 60_000,
})

const journal = computed(() => (journals.value ?? []).find(j => j.id === id.value) ?? null)

const previewOpen = ref(false)

function openPreview(): void {
  if (journal.value?.image_url)
    previewOpen.value = true
}

function closePreview(): void {
  previewOpen.value = false
}

useSeoMeta({
  title: computed(() => journal.value ? `${journal.value.pair} journal — Ctemz Journal` : 'Journal — Ctemz Journal'),
})
</script>

<template>
  <div class="mx-auto w-full max-w-3xl flex-1 px-4 py-8 md:px-6">
    <NuxtLink to="/journals" class="text-sm font-medium text-primary hover:underline">
      ← All journals
    </NuxtLink>

    <div v-if="isPending" class="mt-4 h-96 animate-pulse rounded-2xl bg-surface" />
    <div v-else-if="isError" class="mt-4 rounded-2xl border border-border bg-surface p-8 text-center text-sm text-muted">
      Couldn't load this journal.
    </div>
    <div v-else-if="!journal" class="mt-4 rounded-2xl border border-border bg-surface p-8 text-center">
      <p class="text-sm font-semibold">Journal not found</p>
      <p class="mt-1 text-sm text-muted">It may have been made private or hidden.</p>
      <NuxtLink to="/journals" class="mt-3 inline-block text-sm font-medium text-primary hover:underline">
        Back to journals
      </NuxtLink>
    </div>

    <article v-else class="mt-4 overflow-hidden rounded-2xl border border-border bg-surface">
      <button
        v-if="journal.image_url"
        type="button"
        class="block w-full cursor-zoom-in"
        aria-label="View photo full size"
        @click="openPreview"
      >
        <img :src="journal.image_url" :alt="`${journal.pair} chart`" class="max-h-[480px] w-full object-contain bg-bg">
      </button>
      <div class="p-4 md:p-6">
        <div class="flex flex-wrap items-center gap-2">
          <h1 class="tnum text-2xl font-bold tracking-tight">{{ journal.pair }}</h1>
          <span class="rounded-full border border-border px-2.5 py-0.5 text-xs capitalize text-muted">
            {{ journal.direction === 'long' ? '↑ Long' : '↓ Short' }}{{ journal.timeframe ? ` · ${journal.timeframe}` : '' }}
          </span>
          <span
            v-if="journal.result && journal.r_multiple !== null"
            class="rounded-full px-2.5 py-0.5 text-xs font-semibold"
            :class="journal.result === 'win' ? 'bg-profit/10 text-profit-text' : 'bg-loss/10 text-loss'"
          >
            {{ journal.result === 'win' ? '↑ Win' : '↓ Loss' }} {{ fmtR(journal.r_multiple) }}
          </span>
        </div>
        <dl class="mt-4 grid grid-cols-2 gap-3 text-sm sm:grid-cols-4">
          <div>
            <dt class="text-xs text-muted">Setup</dt>
            <dd class="font-medium">{{ journal.setup_name || 'No setup' }}</dd>
          </div>
          <div>
            <dt class="text-xs text-muted">Trader</dt>
            <dd class="font-medium">{{ journal.display_name || 'Trader' }}</dd>
          </div>
          <div>
            <dt class="text-xs text-muted">Date</dt>
            <dd class="tnum font-medium">{{ journal.date.slice(0, 10) }}</dd>
          </div>
          <div>
            <dt class="text-xs text-muted">Result</dt>
            <dd class="font-medium capitalize">{{ journal.result ?? 'Open' }}</dd>
          </div>
        </dl>
        <p class="mt-4 border-t border-border pt-3 text-xs text-muted">
          Shared publicly. Money, lots and notes are never shown.
        </p>
      </div>
    </article>

    <!-- Photo preview -->
    <div
      v-if="previewOpen && journal?.image_url"
      class="fixed inset-0 z-50 flex items-center justify-center bg-text/70 p-4"
      role="dialog"
      aria-modal="true"
      aria-label="Photo preview"
      @click.self="closePreview"
      @keydown.escape="closePreview"
    >
      <div class="relative max-h-full max-w-4xl overflow-hidden rounded-2xl bg-surface">
        <img :src="journal.image_url" :alt="`${journal.pair} chart full size`" class="max-h-[85vh] w-full object-contain">
        <button
          type="button"
          aria-label="Close preview"
          class="absolute right-3 top-3 rounded-xl bg-text/70 px-3 py-1.5 text-sm font-semibold text-bg transition hover:opacity-90"
          @click="closePreview"
        >
          Close
        </button>
      </div>
    </div>
  </div>
</template>