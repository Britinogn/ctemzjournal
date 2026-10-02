<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { toast } from 'vue-sonner'
import { adminJournalsKey } from '~/types'
import { fmtDay, fmtMoney, fmtR } from '~/utils/format'

definePageMeta({ middleware: 'admin', layout: 'admin' })

interface AdminJournal {
  ID: string;
  UserID: string;
  Pair: string;
  Direction: string;
  Timeframe: string | null;
  Pnl: number | null;
  RMultiple: number | null;
  IsPublic: boolean;
  HiddenByAdmin: boolean;
  CreatedAt: string;
  DisplayName: string | null;
  SetupName: string | null;
}

const api = useApi()
const queryClient = useQueryClient()
const page = ref(0)
const LIMIT = 20

const { data: journals, isPending, isError, refetch } = useQuery({
  queryKey: computed(() => adminJournalsKey(page.value)),
  queryFn: () => api.get<AdminJournal[]>(`/admin/journals?limit=${LIMIT}&offset=${page.value * LIMIT}`),
})

const acting = ref<string | null>(null)

async function setHidden(id: string, hidden: boolean): Promise<void> {
  if (acting.value)
    return
  acting.value = id
  try {
    await api.patch(`/admin/journals/${id}/hide`, { hidden })
    queryClient.invalidateQueries({ queryKey: ['admin', 'journals'] })
    await refetch()
    toast.success(hidden ? 'Journal hidden from public' : 'Journal restored to public')
  }
  catch {
    toast.error('Could not update journal')
  }
  finally {
    acting.value = null
  }
}
</script>

<template>
  <div>
    <h1 class="mb-1 text-xl font-bold tracking-tight md:text-2xl">
      Public journals
    </h1>
    <p class="mb-4 text-sm text-muted">
      Hide journals that shouldn't be on the home page. Hidden ones never appear publicly.
    </p>

    <div v-if="isPending" class="h-96 animate-pulse rounded-2xl bg-surface" />
    <div v-else-if="isError" class="rounded-2xl border border-border bg-surface p-8 text-center">
      <p class="text-sm text-muted">
        Couldn't load journals.
      </p>
      <button
        type="button"
        class="mt-3 rounded-xl bg-primary px-5 py-2 text-sm font-semibold text-on-primary transition hover:opacity-90"
        @click="() => refetch()"
      >
        Retry
      </button>
    </div>

    <section v-else class="overflow-hidden rounded-2xl border border-border bg-surface" aria-label="Journals">
      <ul v-if="(journals ?? []).length > 0" class="divide-y divide-border">
        <li v-for="j in journals" :key="j.ID" class="flex flex-wrap items-center gap-2 px-4 py-3">
          <div class="min-w-0 flex-1">
            <p class="tnum text-sm font-bold">
              {{ j.Pair }}
              <span class="ml-1 text-xs font-medium capitalize" :class="j.Direction === 'long' ? 'text-profit-text' : 'text-loss'">
                {{ j.Direction }}
              </span>
            </p>
            <p class="truncate text-xs text-muted">
              {{ j.DisplayName || 'Trader' }} · {{ j.SetupName || 'No setup' }} · {{ fmtDay(j.CreatedAt) }}
              <span v-if="j.RMultiple !== null" class="tnum font-semibold"> · {{ fmtR(j.RMultiple) }}</span>
              <span v-if="j.Pnl !== null" class="tnum"> · {{ fmtMoney(j.Pnl) }}</span>
            </p>
          </div>
          <span
            v-if="j.HiddenByAdmin"
            class="rounded-full bg-warning/10 px-2.5 py-0.5 text-xs font-semibold text-warning-text"
          >
            Hidden
          </span>
          <button
            type="button"
            :disabled="acting === j.ID"
            class="shrink-0 rounded-lg border border-border px-3 py-1.5 text-xs font-semibold transition hover:border-primary disabled:opacity-50"
            @click="setHidden(j.ID, !j.HiddenByAdmin)"
          >
            {{ acting === j.ID ? 'Working…' : j.HiddenByAdmin ? 'Restore' : 'Hide' }}
          </button>
        </li>
      </ul>
      <p v-else class="p-8 text-center text-sm text-muted">
        No public journals right now.
      </p>
      <div class="flex items-center justify-end gap-2 border-t border-border px-4 py-3">
        <button
          type="button" :disabled="page === 0"
          class="rounded-xl border border-border px-4 py-1.5 text-sm font-medium transition hover:border-primary disabled:opacity-40"
          @click="page--"
        >
          Previous
        </button>
        <button
          type="button" :disabled="(journals ?? []).length < LIMIT"
          class="rounded-xl border border-border px-4 py-1.5 text-sm font-medium transition hover:border-primary disabled:opacity-40"
          @click="page++"
        >
          Next
        </button>
      </div>
    </section>
  </div>
</template>
