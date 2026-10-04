<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { journalsKey, type PublicJournal } from '~/types'

const api = useApi()

const { data: journals, isPending, isError, refetch } = useQuery({
  queryKey: journalsKey(9),
  queryFn: () => api.get<PublicJournal[]>('/public/journals?limit=9'),
  staleTime: 60_000,
})
</script>

<template>
  <section class="mx-auto w-full max-w-6xl px-4 pt-16 md:px-6" aria-labelledby="journals-title">
    <div class="mb-5 flex items-end justify-between gap-4">
      <div>
        <p class="text-xs font-semibold uppercase tracking-wide text-primary">
          Community
        </p>
        <h2 id="journals-title" class="text-2xl font-bold tracking-tight">
          Recent journals
        </h2>
        <p class="mt-1 text-sm text-muted">
          Shared by traders. Results are shown in R, never in money.
        </p>
      </div>
      <NuxtLink to="/journals" class="shrink-0 text-sm font-medium text-primary hover:underline">
        View all
      </NuxtLink>
    </div>

    <!-- The border keeps skeletons visible on the light background -->
    <div v-if="isPending" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3" aria-hidden="true">
      <div v-for="i in 6" :key="i" class="h-64 animate-pulse rounded-2xl border border-border bg-surface" />
    </div>

    <div v-else-if="isError" class="rounded-2xl border border-border bg-surface p-8 text-center">
      <p class="text-sm text-muted">
        Journals did not load. Check your connection and try again.
      </p>
      <button
        type="button"
        class="mt-3 inline-flex h-10 items-center rounded-xl border border-border px-4 text-sm font-semibold transition-colors hover:border-primary"
        @click="refetch()"
      >
        Try again
      </button>
    </div>

    <div v-else-if="(journals ?? []).length > 0" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <HomeJournalCard v-for="j in journals" :key="j.id" :journal="j" />
    </div>

    <p v-else class="rounded-2xl border border-border bg-surface p-8 text-center text-sm text-muted">
      No public journals yet. Yours could be the first.
    </p>
  </section>
</template>