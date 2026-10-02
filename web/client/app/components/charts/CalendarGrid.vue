<script setup lang="ts">
import { ArrowLeft01Icon, ArrowRight01Icon } from '~/utils/icons'
import type { CalendarDay } from '~/types'
import { fmtSigned } from '~/utils/format'

const props = defineProps<{ days: CalendarDay[] }>()

const weekdays = ['M', 'T', 'W', 'T', 'F', 'S', 'S']

/** Months present in the data, newest first (YYYY-MM). */
const months = computed(() => {
  const set = new Set(props.days.map(d => d.date.slice(0, 7)))
  return [...set].sort().reverse()
})

const cursor = ref('')
watchEffect(() => {
  if (!cursor.value && months.value.length > 0)
    cursor.value = months.value[0]
})

function shift(dir: 1 | -1): void {
  const i = months.value.indexOf(cursor.value)
  const next = months.value[i + dir]
  if (next)
    cursor.value = next
}

const monthDays = computed(() => props.days.filter(d => d.date.startsWith(cursor.value)))
const monthPnl = computed(() => monthDays.value.reduce((s, d) => s + d.pnl, 0))

const monthLabel = computed(() => {
  if (!cursor.value)
    return ''
  const [y, m] = cursor.value.split('-').map(Number)
  return new Date(Date.UTC(y, m - 1, 1)).toLocaleDateString('en-GB', { month: 'long', timeZone: 'UTC' })
})

/** Monday-first offset for the 1st of the cursor month. */
const leadBlanks = computed(() => {
  if (!cursor.value)
    return 0
  const [y, m] = cursor.value.split('-').map(Number)
  return (new Date(Date.UTC(y, m - 1, 1)).getUTCDay() + 6) % 7
})

function cellTone(pnl: number): string {
  if (pnl > 0)
    return 'bg-profit/10'
  if (pnl < 0)
    return 'bg-loss/10'
  return 'bg-bg'
}

function textTone(pnl: number): string {
  if (pnl > 0)
    return 'text-profit-text'
  if (pnl < 0)
    return 'text-loss'
  return 'text-muted'
}
</script>

<template>
  <div>
    <div class="mb-3 flex items-center justify-between">
      <div class="flex items-center gap-1">
        <button
          type="button"
          aria-label="Previous month"
          class="inline-flex h-8 w-8 items-center justify-center rounded-lg text-muted transition hover:bg-bg hover:text-text"
          @click="shift(1)"
        >
          <UiAppIcon :icon="ArrowLeft01Icon" :size="18" />
        </button>
        <h3 class="min-w-28 text-center text-sm font-semibold">
          {{ monthLabel }}
        </h3>
        <button
          type="button"
          aria-label="Next month"
          class="inline-flex h-8 w-8 items-center justify-center rounded-lg text-muted transition hover:bg-bg hover:text-text"
          @click="shift(-1)"
        >
          <UiAppIcon :icon="ArrowRight01Icon" :size="18" />
        </button>
      </div>
      <p class="tnum text-sm font-bold" :class="textTone(monthPnl)">
        {{ fmtSigned(monthPnl, 0) }}
      </p>
    </div>

    <div class="grid grid-cols-7 gap-1 text-center text-[11px] font-medium text-muted">
      <span v-for="w in weekdays" :key="w" class="py-1">{{ w }}</span>
    </div>
    <div class="grid grid-cols-7 gap-1">
      <span v-for="i in leadBlanks" :key="`b-${i}`" />
      <div
        v-for="d in monthDays"
        :key="d.date"
        :class="['rounded-lg px-1 py-1.5 text-center', cellTone(d.pnl)]"
      >
        <p class="text-[11px] text-muted">
          {{ Number(d.date.slice(8, 10)) }}
        </p>
        <p class="tnum text-[11px] font-bold leading-tight" :class="textTone(d.pnl)">
          {{ fmtSigned(d.pnl, 0) }}
        </p>
      </div>
    </div>
    <p class="mt-3 text-[11px] leading-relaxed text-muted">
      Days are grouped in Lagos time. Each cell shows the day's result with a sign.
    </p>
  </div>
</template>
