<script setup lang="ts">
import {
  CategoryScale,
  Chart as ChartJS,
  Filler,
  LinearScale,
  LineElement,
  PointElement,
  Tooltip,
} from 'chart.js'
import type { TooltipItem } from 'chart.js'
import { Line } from 'vue-chartjs'
import type { EquityPoint } from '~/types'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Tooltip, Filler)

const props = withDefaults(defineProps<{
  points: EquityPoint[];
  /** Days back from the last point (0 = all). */
  rangeDays?: number;
}>(), {
  rangeDays: 0,
})

const colorMode = useColorMode()

function cssVar(name: string, fallback: string): string {
  if (import.meta.server || typeof document === 'undefined')
    return fallback
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim() || fallback
}

const filtered = computed(() => {
  if (!props.rangeDays || props.points.length === 0)
    return props.points
  const lastPoint = props.points.at(-1)
  if (!lastPoint)
    return props.points
  const last = new Date(lastPoint.date).getTime()
  const cutoff = last - props.rangeDays * 86_400_000
  return props.points.filter(p => new Date(p.date).getTime() >= cutoff)
})

const chartData = computed(() => {
  const primary = cssVar('--primary', '#0E7C86')
  return {
    labels: filtered.value.map(p => p.date),
    datasets: [
      {
        data: filtered.value.map(p => p.equity),
        borderColor: primary,
        backgroundColor: `${primary}1A`,
        fill: true,
        tension: 0.3,
        pointRadius: 0,
        pointHoverRadius: 4,
        borderWidth: 2,
      },
    ],
  }
})

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      callbacks: {
        title: (items: TooltipItem<'line'>[]) => items[0]?.label ?? '',
        label: (item: TooltipItem<'line'>) => `Equity: ${item.parsed.y ?? 0}`,
      },
    },
  },
  scales: {
    x: {
      grid: { color: cssVar('--border', '#E2E8F0') },
      ticks: { color: cssVar('--muted', '#64748B'), maxTicksLimit: 6 },
    },
    y: {
      grid: { color: cssVar('--border', '#E2E8F0') },
      ticks: { color: cssVar('--muted', '#64748B'), maxTicksLimit: 6 },
    },
  },
}))
</script>

<template>
  <!-- remount on theme change so Chart.js re-reads the CSS variables -->
  <div :key="colorMode.value" class="h-64 w-full">
    <Line :data="chartData" :options="chartOptions" />
  </div>
</template>