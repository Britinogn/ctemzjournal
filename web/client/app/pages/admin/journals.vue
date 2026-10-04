<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { toast } from 'vue-sonner'
import { adminJournalsKey } from '~/types'
import { fmtR } from '~/utils/format'

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
const LIMIT = 12

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

function resultOf(j: AdminJournal): { win: boolean; label: string } | null {
  if (j.RMultiple === null || j.Pnl === null)
    return null
  const win = j.Pnl > 0
  return { win, label: `${win ? 'Win' : 'Loss'} ${fmtR(j.RMultiple)}` }
}

/** Decorative trend accent (no per-trade history exists): rising for wins. */
function sparkPath(win: boolean | null): string {
  if (win === null)
    return 'M0,30 L60,30 L120,30'
  return win
    ? 'M0,34 L30,30 L60,32 L90,20 L120,10'
    : 'M0,10 L30,14 L60,12 L90,26 L120,34'
}
</script>

<template>
  <div>
    <h1 class="mb-1 text-xl font-bold tracking-tight md:text-2xl">
      Public journals
    </h1>
    <p class="mb-4 text-sm text-muted">
      Public journals show pair, setup and result in R only. Hide one to remove it from the home page.
    </p>

    <div v-if="isPending" class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
      <div v-for="i in 6" :key="i" class="h-64 animate-pulse rounded-2xl bg-surface" />
    </div>
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

    <template v-else>
      <div v-if="(journals ?? []).length > 0" class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
        <article v-for="j in journals" :key="j.ID" class="overflow-hidden rounded-2xl border border-border bg-surface" :aria-label="`${j.Pair} journal`">
          <svg viewBox="0 0 120 40" class="h-20 w-full" aria-hidden="true" preserveAspectRatio="none">
            <path
              :d="sparkPath(resultOf(j)?.win ?? null)"
              fill="none"
              :stroke="(resultOf(j)?.win ?? false) ? 'var(--profit)' : 'var(--loss)'"
              stroke-width="2.5"
              stroke-linecap="round"
            />
            <path
              :d="`${sparkPath(resultOf(j)?.win ?? null)} L120,40 L0,40 Z`"
              :fill="(resultOf(j)?.win ?? false) ? 'var(--profit)' : 'var(--loss)'"
              opacity="0.12"
              stroke="none"
            />
          </svg>
          <div class="p-4">
            <div class="flex items-center justify-between gap-2">
              <p class="tnum text-base font-bold">
                {{ j.Pair }}
              </p>
              <span class="rounded-full border border-border px-2 py-0.5 text-xs capitalize text-muted">
                {{ j.Direction === 'long' ? '↑ Long' : '↓ Short' }}{{ j.Timeframe ? ` · ${j.Timeframe}` : '' }}
              </span>
            </div>
            <p class="mt-0.5 truncate text-sm text-muted">
              {{ j.SetupName || 'No setup' }}{{ j.Timeframe ? ` ${j.Timeframe}` : '' }}
            </p>
            <div class="mt-1.5 flex items-center justify-between gap-2">
              <span
                v-if="resultOf(j)"
                class="rounded-full px-2 py-0.5 text-xs font-semibold"
                :class="resultOf(j)!.win ? 'bg-profit/10 text-profit-text' : 'bg-loss/10 text-loss'"
              >
                {{ resultOf(j)!.win ? '↑' : '↓' }} {{ resultOf(j)!.label }}
              </span>
              <span v-else class="text-xs text-muted">Open</span>
              <span class="truncate text-xs text-muted">{{ j.DisplayName || 'Trader' }}</span>
            </div>
            <p v-if="j.HiddenByAdmin" class="mt-1.5 inline-block rounded-full bg-warning/10 px-2 py-0.5 text-xs font-semibold text-warning-text">
              Hidden from home
            </p>
            <button
              type="button"
              :disabled="acting === j.ID"
              class="mt-3 w-full rounded-xl border border-border py-2 text-sm font-semibold transition hover:border-primary disabled:opacity-50"
              @click="setHidden(j.ID, !j.HiddenByAdmin)"
            >
              {{ acting === j.ID ? 'Working…' : j.HiddenByAdmin ? 'Restore' : 'Hide from home' }}
            </button>
          </div>
        </article>
      </div>
      <p v-else class="rounded-2xl border border-border bg-surface p-8 text-center text-sm text-muted">
        No public journals right now.
      </p>
      <div class="mt-4 flex items-center justify-end gap-2">
        <button
          type="button" :disabled="page === 0"
          class="rounded-xl border border-border bg-surface px-4 py-1.5 text-sm font-medium transition hover:border-primary disabled:opacity-40"
          @click="page--"
        >
          Previous
        </button>
        <button
          type="button" :disabled="(journals ?? []).length < LIMIT"
          class="rounded-xl border border-border bg-surface px-4 py-1.5 text-sm font-medium transition hover:border-primary disabled:opacity-40"
          @click="page++"
        >
          Next
        </button>
      </div>
    </template>
  </div>
</template>
