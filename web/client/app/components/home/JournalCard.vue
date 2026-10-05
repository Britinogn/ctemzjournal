<script setup lang="ts">
import type { PublicJournal } from '~/types'
import { fmtR } from '~/utils/format'

const props = defineProps<{ journal: PublicJournal }>()

const isWin = computed(() => props.journal.result === 'win')
const isLoss = computed(() => props.journal.result === 'loss')

// Decorative shape only. An open trade has no result, so it gets a flat neutral line, not a red one.
const line = computed(() => {
  if (isWin.value)
    return 'M0,34 L30,30 L60,32 L90,20 L120,10'
  if (isLoss.value)
    return 'M0,10 L30,14 L60,12 L90,26 L120,34'
  return 'M0,22 L120,22'
})
const tone = computed(() => (isWin.value ? 'var(--profit)' : isLoss.value ? 'var(--loss)' : 'var(--muted)'))
</script>

<template>
  <NuxtLink
    :to="`/journals/${journal.id}`"
    class="flex h-full flex-col overflow-hidden rounded-2xl border border-border bg-surface transition-colors hover:border-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
    :aria-label="`${journal.pair} journal by ${journal.display_name || 'a trader'}`"
  >
    <div v-if="journal.image_url" class="h-36 w-full overflow-hidden border-b border-border">
      <img :src="journal.image_url" :alt="`${journal.pair} chart`" class="h-full w-full object-cover" loading="lazy" decoding="async">
    </div>
    <svg v-else viewBox="0 0 120 40" class="h-36 w-full border-b border-border bg-bg" aria-hidden="true" preserveAspectRatio="none">
      <path :d="`${line} L120,40 L0,40 Z`" :fill="tone" opacity="0.14" />
      <path :d="line" fill="none" :stroke="tone" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" vector-effect="non-scaling-stroke" />
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
          v-if="(isWin || isLoss) && journal.r_multiple !== null"
          class="tnum rounded-full px-2 py-0.5 text-xs font-semibold"
          :class="isWin ? 'bg-profit/10 text-profit-text' : 'bg-loss/10 text-loss'"
        >
          {{ isWin ? '↑' : '↓' }} {{ isWin ? 'Win' : 'Loss' }} {{ fmtR(journal.r_multiple) }}
        </span>
        <span v-else class="text-xs text-muted">Open</span>
        <span class="truncate text-xs text-muted">{{ journal.display_name || 'Trader' }}</span>
      </div>
    </div>
  </NuxtLink>
</template>