<script setup lang="ts">
import { fmtR } from '~/utils/format'

export interface BreakdownItem {
  label: string;
  value: number;
}

const props = defineProps<{
  title: string;
  items: BreakdownItem[];
}>();

const maxAbs = computed(() => Math.max(0.0001, ...props.items.map(i => Math.abs(i.value))))
</script>

<template>
  <div>
    <div class="mb-3 flex items-center justify-between">
      <h3 class="text-sm font-semibold">
        {{ title }}
      </h3>
      <span class="text-xs text-muted">Avg R</span>
    </div>
    <ul class="space-y-2.5">
      <li v-for="item in items" :key="item.label">
        <div class="mb-1 flex items-center justify-between gap-2 text-sm">
          <span class="truncate">{{ item.label }}</span>
          <span class="tnum shrink-0 font-bold" :class="item.value >= 0 ? 'text-profit-text' : 'text-loss'">
            {{ fmtR(item.value) }}
          </span>
        </div>
        <div class="h-1.5 overflow-hidden rounded-full bg-bg">
          <div
            class="h-full rounded-full transition-all"
            :class="item.value >= 0 ? 'bg-profit' : 'bg-loss'"
            :style="{ width: `${Math.max(2, (Math.abs(item.value) / maxAbs) * 100)}%` }"
          />
        </div>
      </li>
    </ul>
    <p v-if="items.length === 0" class="py-4 text-center text-sm text-muted">
      No closed trades yet
    </p>
  </div>
</template>
