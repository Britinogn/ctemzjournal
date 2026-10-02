<script setup lang="ts">
import { Alert02Icon, ArrowDown01Icon, ArrowUp01Icon, Tick02Icon } from '~/utils/icons'
import type { Trade } from '~/types'
import { fmtDay, fmtMoney, fmtR } from '~/utils/format'

defineProps<{
  trades: Trade[];
  /** Setup ID → name map (dashboard resolves via /setups). */
  setupNames?: Record<string, string>;
}>()

function setupLabel(t: Trade, names?: Record<string, string>): string {
  const name = (t.SetupID && names?.[t.SetupID]) || null
  const base = name ?? 'No setup'
  return t.Timeframe ? `${base} · ${t.Timeframe}` : base
}

function sideTone(dir: string): string {
  return dir === 'long' ? 'text-profit-text' : 'text-loss'
}

function resultOf(t: Trade): { label: string; cls: string } | null {
  if (t.RMultiple === null || t.Pnl === null)
    return null
  const win = t.Pnl > 0
  return {
    label: `${win ? 'Win' : 'Loss'} ${fmtR(t.RMultiple)}`,
    cls: win ? 'bg-profit/10 text-profit-text' : 'bg-loss/10 text-loss',
  }
}
</script>

<template>
  <div>
    <!-- Desktop table -->
    <table class="hidden w-full text-left text-sm md:table">
      <thead>
        <tr class="text-xs text-muted">
          <th class="pb-2 pr-3 font-medium">Pair</th>
          <th class="pb-2 pr-3 font-medium">Side</th>
          <th class="pb-2 pr-3 font-medium">Setup</th>
          <th class="pb-2 pr-3 font-medium">Rules</th>
          <th class="pb-2 pr-3 font-medium">Date</th>
          <th class="pb-2 pr-3 text-right font-medium">P&amp;L</th>
          <th class="pb-2 text-right font-medium">Result</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-border">
        <tr v-for="t in trades" :key="t.ID">
          <td class="tnum py-2.5 pr-3 font-semibold">{{ t.Pair }}</td>
          <td class="py-2.5 pr-3">
            <span class="inline-flex items-center gap-1 text-xs" :class="sideTone(t.Direction)">
              <UiAppIcon :icon="t.Direction === 'long' ? ArrowUp01Icon : ArrowDown01Icon" :size="14" />
              {{ t.Direction === 'long' ? 'Long' : 'Short' }}
            </span>
          </td>
          <td class="py-2.5 pr-3 text-muted">{{ setupLabel(t, setupNames) }}</td>
          <td class="py-2.5 pr-3">
            <span
              v-if="t.FollowedRules === true"
              class="inline-flex items-center gap-1 rounded-full bg-primary/10 px-2 py-0.5 text-xs font-medium text-primary"
            >
              <UiAppIcon :icon="Tick02Icon" :size="12" /> Followed
            </span>
            <span
              v-else-if="t.FollowedRules === false"
              class="inline-flex items-center gap-1 rounded-full bg-warning/10 px-2 py-0.5 text-xs font-medium text-warning-text"
            >
              <UiAppIcon :icon="Alert02Icon" :size="12" /> Broken
            </span>
            <span v-else class="text-xs text-muted">—</span>
          </td>
          <td class="py-2.5 pr-3 text-muted">{{ t.ClosedAt ? fmtDay(t.ClosedAt) : 'Open' }}</td>
          <td class="tnum py-2.5 pr-3 text-right font-semibold" :class="t.Pnl !== null && t.Pnl < 0 ? 'text-loss' : 'text-profit-text'">
            {{ t.Pnl === null ? '—' : fmtMoney(t.Pnl) }}
          </td>
          <td class="py-2.5 text-right">
            <span
              v-if="resultOf(t)"
              class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-semibold"
              :class="resultOf(t)!.cls"
            >
              <UiAppIcon :icon="t.Pnl! > 0 ? ArrowUp01Icon : ArrowDown01Icon" :size="12" />
              {{ resultOf(t)!.label }}
            </span>
            <span v-else class="text-xs text-muted">—</span>
          </td>
        </tr>
      </tbody>
    </table>

    <!-- Mobile / small-tablet cards -->
    <ul class="space-y-2 md:hidden">
      <li v-for="t in trades" :key="t.ID" class="rounded-xl border border-border p-3">
        <div class="flex items-center justify-between">
          <p class="tnum font-bold">{{ t.Pair }}</p>
          <p class="tnum text-sm font-bold" :class="t.Pnl !== null && t.Pnl < 0 ? 'text-loss' : 'text-profit-text'">
            {{ t.Pnl === null ? '—' : fmtMoney(t.Pnl) }}
          </p>
        </div>
        <div class="mt-1.5 flex items-center justify-between text-xs text-muted">
          <span class="inline-flex items-center gap-1" :class="sideTone(t.Direction)">
            <UiAppIcon :icon="t.Direction === 'long' ? ArrowUp01Icon : ArrowDown01Icon" :size="12" />
            {{ t.Direction === 'long' ? 'Long' : 'Short' }} · {{ t.ClosedAt ? fmtDay(t.ClosedAt) : 'Open' }}
          </span>
          <span v-if="resultOf(t)" class="font-semibold" :class="resultOf(t)!.cls.split(' ')[1]">
            {{ resultOf(t)!.label }}
          </span>
        </div>
      </li>
    </ul>

    <p v-if="trades.length === 0" class="py-6 text-center text-sm text-muted">
      No trades yet — log your first one.
    </p>
  </div>
</template>
