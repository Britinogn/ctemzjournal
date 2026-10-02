<script setup lang="ts">
withDefaults(defineProps<{
  label: string;
  value: string;
  sub?: string;
  tone?: 'neutral' | 'profit' | 'loss' | 'warning';
  icon: unknown;
}>(), {
  tone: 'neutral',
})

const tones: Record<string, string> = {
  neutral: 'text-text',
  profit: 'text-profit',
  loss: 'text-loss',
  warning: 'text-warning',
}
</script>

<template>
  <div class="rounded-2xl border border-border bg-surface p-4">
    <div class="flex items-center gap-1.5 text-xs font-medium text-muted">
      <UiAppIcon :icon="icon" :size="14" />
      {{ label }}
    </div>
    <p class="tnum mt-1 text-2xl font-bold tracking-tight" :class="tones[tone]">
      {{ value }}
    </p>
    <p v-if="sub" class="mt-0.5 text-xs text-muted">
      {{ sub }}
    </p>
  </div>
</template>
