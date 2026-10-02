<script setup lang="ts">
import { initials } from '~/utils/format';

const props = withDefaults(defineProps<{
  name?: string | null;
  email?: string;
  role?: string;
  size?: number;
}>(), {
  name: null,
  email: '',
  role: 'user',
  size: 40,
});

const label = computed(() => initials(props.name, props.email));
</script>

<template>
  <span
    class="inline-flex shrink-0 items-center justify-center rounded-full bg-primary/15 font-semibold text-primary"
    :style="{ width: `${size}px`, height: `${size}px`, fontSize: `${Math.round(size * 0.36)}px` }"
    :title="name || email || 'Account'"
  >
    {{ label }}
  </span>
</template>
