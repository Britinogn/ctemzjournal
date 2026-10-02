<script setup lang="ts">
import { Moon02Icon, Sun03Icon } from '~/utils/icons'

const colorMode = useColorMode()
const isDark = computed(() => colorMode.value === 'dark')

function toggle(): void {
  colorMode.preference = isDark.value ? 'light' : 'dark'
}
</script>

<template>
  <button
    type="button"
    :aria-label="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
    class="inline-flex h-10 w-10 items-center justify-center rounded-xl border border-border bg-surface text-muted transition hover:text-text"
    @click="toggle"
  >
    <UiAppIcon :icon="isDark ? Sun03Icon : Moon02Icon" :size="20" />
  </button>
</template>
