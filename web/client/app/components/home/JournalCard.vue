<script setup lang="ts">
import type { PublicJournal } from '~/types'
import { fmtR } from '~/utils/format'

const props = defineProps<{ journal: PublicJournal }>()

const isWin = computed(() => props.journal.result === 'win')
const isLoss = computed(() => props.journal.result === 'loss')

// Decorative shape only. It is not real trade data, so open trades get a flat neutral line.
const line = computed(() => {
  if (isWin.value) return 'M0,34 L30,30 L60,32 L90,20 L120,10'
  if (isLoss.value) return 'M0,10 L30,14 L60,12 L90,26 L120,34'
  return 'M0,22 L120,22'
})
// Green and red mean results only. An open trade has no result, so it stays neutral.
const tone = computed(() => (isWin.value ? 'var(--profit)' : isLoss.value ? 'var(--loss)' : 'var(--muted)'))
</script>

<template>
  <article class="flex h-full flex-col overflow-hidden rounded-2xl border border-border bg-surface" :aria-label="`${journal.pair} journal`">
    <div v-if="journal.image_url" class="h-36 w-full overflow-hidden border-b border-border bg-bg">
      <img
        :src="journal.image_url"
        :alt="`Chart screenshot for ${journal.pair}`"
        width="640"
        height="288"
        class="h-full w-full object-cover"
        loading="lazy"
        decoding="async"
      >
    </div>
    <svg v-else viewBox="0 0 120 40" class="h-36 w-full border-b border-border bg-bg" aria-hidden="true" preserveAspectRatio="none">
      <path :d="`${line} L120,40 L0,40 Z`" :fill="tone" opacity="0.14" />
      <path :d="line" fill="none" :stroke="tone" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" vector-effect="non-scaling-stroke" />
    </svg>

    <div class="flex flex-1 flex-col gap-2 p-4">
      <div class="flex items-center justify-between gap-2">
        <p class="tnum text-base font-bold">
          {{ journal.pair }}
        </p>
        <span class="inline-flex shrink-0 items-center gap-1 rounded-full border border-border px-2 py-0.5 text-xs font-medium text-muted">
          <span aria-hidden="true">{{ journal.direction === 'long' ? '↑' : '↓' }}</span>
          {{ journal.direction === 'long' ? 'Long' : 'Short' }}
          <span v-if="journal.timeframe" class="tnum">{{ journal.timeframe }}</span>
        </span>
      </div>

      <p class="truncate text-sm text-muted">
        {{ journal.setup_name || 'No setup' }}
      </p>

      <div class="mt-auto flex items-center justify-between gap-2 pt-1">
        <!-- Sign and arrow always shown, so win or loss never relies on colour alone -->
        <span
          v-if="(isWin || isLoss) && journal.r_multiple !== null"
          class="tnum inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-semibold"
          :class="isWin ? 'bg-profit/10 text-profit-text' : 'bg-loss/10 text-loss'"
        >
          <span aria-hidden="true">{{ isWin ? '↑' : '↓' }}</span>
          {{ isWin ? 'Win' : 'Loss' }} {{ fmtR(journal.r_multiple) }}
        </span>
        <span v-else class="rounded-full border border-border px-2.5 py-0.5 text-xs font-medium text-muted">Open</span>
        <span class="truncate text-xs text-muted">{{ journal.display_name || 'Trader' }}</span>
      </div>
    </div>
  </article>
</template>