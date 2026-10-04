<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { ratesKey, type RatesResponse } from '~/types'
import { timeAgo } from '~/utils/format'

const api = useApi()

const { data: rates, isPending, isError } = useQuery({
  queryKey: ratesKey(),
  queryFn: () => api.get<RatesResponse>('/public/rates'),
  refetchInterval: 60_000,
})

const items = computed(() => (rates.value?.rates ?? []).slice(0, 8))
</script>

<template>
  <section class="mx-auto w-full max-w-6xl px-4 pt-16 md:px-6" aria-labelledby="prices-title">
    <div class="mb-4 flex flex-wrap items-end justify-between gap-x-4 gap-y-1">
      <div>
        <p class="text-xs font-semibold uppercase tracking-wide text-primary">
          Market data
        </p>
        <h2 id="prices-title" class="text-2xl font-bold tracking-tight">
          Live prices
        </h2>
      </div>
      <div class="flex flex-wrap items-center gap-x-3 gap-y-1 text-sm text-muted">
        <!-- Warning colour plus an icon and words, so staleness is not colour alone -->
        <span
          v-if="rates?.stale"
          class="inline-flex items-center gap-1 rounded-full bg-warning/10 px-2.5 py-0.5 font-semibold text-warning-text"
        >
          <svg class="h-3.5 w-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M12 4l9 16H3L12 4zM12 10v4M12 17h.01" />
          </svg>
          Prices may be out of date
        </span>
        <span v-if="rates?.updated_at" class="inline-flex items-center gap-1.5">
          <span class="h-1.5 w-1.5 rounded-full" :class="rates.stale ? 'bg-warning' : 'bg-primary'" aria-hidden="true" />
          Updated {{ timeAgo(rates.updated_at) }}
        </span>
        <NuxtLink to="/prices" class="font-medium text-primary hover:underline">
          All prices
        </NuxtLink>
      </div>
    </div>

    <!-- Loading: same shape as the real grid, so nothing jumps when data arrives -->
    <div v-if="isPending" class="grid grid-cols-2 gap-3 sm:grid-cols-4" aria-hidden="true">
      <div v-for="i in 8" :key="i" class="h-[74px] animate-pulse rounded-2xl border border-border bg-surface" />
    </div>

    <p v-else-if="isError || items.length === 0" class="rounded-2xl border border-border bg-surface p-6 text-center text-sm text-muted">
      Prices are not available right now. They refresh every minute, so check back shortly.
    </p>

    <dl v-else class="grid grid-cols-2 gap-3 sm:grid-cols-4">
      <div v-for="r in items" :key="r.pair" class="rounded-2xl border border-border bg-surface px-4 py-3">
        <dt class="tnum text-xs font-medium text-muted">
          {{ r.pair }}
        </dt>
        <dd class="tnum mt-0.5 text-xl font-semibold" :class="r.stale ? 'text-warning-text' : ''">
          {{ r.price }}
        </dd>
      </div>
    </dl>
  </section>
</template>