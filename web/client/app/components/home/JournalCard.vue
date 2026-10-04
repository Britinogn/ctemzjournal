<script setup lang="ts">
import type { PublicJournal } from '~/types'
import { fmtR } from '~/utils/format'

defineProps<{ journal: PublicJournal }>()

function sparkPath(win: boolean | null): string {
  if (win === null)
    return 'M0,30 L60,30 L120,30'
  return win
    ? 'M0,34 L30,30 L60,32 L90,20 L120,10'
    : 'M0,10 L30,14 L60,12 L90,26 L120,34'
}
</script>

<template>
  <NuxtLink
    :to="`/journals/${journal.id}`"
    class="flex h-full flex-col overflow-hidden rounded-2xl border border-border bg-surface transition hover:border-primary"
    :aria-label="`${journal.pair} journal by ${journal.display_name || 'a trader'}`"
  >
    <div v-if="journal.image_url" class="h-36 w-full overflow-hidden border-b border-border">
      <img :src="journal.image_url" :alt="`${journal.pair} chart`" class="h-full w-full object-cover" loading="lazy">
    </div>
    <svg v-else viewBox="0 0 120 40" class="h-36 w-full" aria-hidden="true" preserveAspectRatio="none">
      <path
        :d="sparkPath(journal.result === null ? null : journal.result === 'win')"
        fill="none"
        :stroke="(journal.result ?? '') === 'win' ? 'var(--profit)' : 'var(--loss)'"
        stroke-width="2.5"
        stroke-linecap="round"
      />
      <path
        :d="`${sparkPath(journal.result === null ? null : journal.result === 'win')} L120,40 L0,40 Z`"
        :fill="(journal.result ?? '') === 'win' ? 'var(--profit)' : 'var(--loss)'"
        opacity="0.12"
        stroke="none"
      />
    </svg>
    <div class="flex flex-1 flex-col p-4">
      <div class="flex items-center justify-between gap-2">
        <p class="tnum text-base font-bold">
          {{ journal.pair }}
        </p>
        <span class="rounded-full border border-border px-2 py-0.5 text-xs capitalize text-muted">
          {{ journal.direction === 'long' ? '↑ Long' : '↓ Short' }}
        </span>
      </div>
      <p class="mt-0.5 truncate text-sm text-muted">
        {{ journal.setup_name || 'No setup' }}{{ journal.timeframe ? ` ${journal.timeframe}` : '' }}
      </p>
      <div class="mt-1.5 flex items-center justify-between gap-2">
        <span
          v-if="journal.result && journal.r_multiple !== null"
          class="rounded-full px-2 py-0.5 text-xs font-semibold"
          :class="journal.result === 'win' ? 'bg-profit/10 text-profit-text' : 'bg-loss/10 text-loss'"
        >
          {{ journal.result === 'win' ? '↑' : '↓' }} {{ journal.result === 'win' ? 'Win' : 'Loss' }} {{ fmtR(journal.r_multiple) }}
        </span>
        <span v-else class="text-xs text-muted">Open</span>
        <span class="truncate text-xs text-muted">{{ journal.display_name || 'Trader' }}</span>
      </div>
    </div>
  </NuxtLink>
</template>