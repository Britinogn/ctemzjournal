<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'

defineProps<{
  /**
   * Hugeicons icon data (e.g. Moon02Icon). Deliberately `unknown`: the core
   * icon union includes readonly variants that fight template type-checking.
   * The runtime only needs an Array — cast once, here, nowhere else.
   */
  icon: unknown;
  size?: number | string;
}>()
</script>

<template>
  <HugeiconsIcon :icon="icon as unknown[]" :size="size ?? 20" />
</template>
