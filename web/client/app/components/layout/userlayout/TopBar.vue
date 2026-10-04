<script setup lang="ts">
import { useWindowScroll } from '@vueuse/core'
import type { Profile } from '~/types'

defineProps<{ profile: Profile | null }>()

// A soft shadow appears only once content has scrolled under the bar.
const { y } = useWindowScroll()
const scrolled = computed(() => y.value > 4)
</script>

<template>
  <!-- Mobile top bar (tablet and desktop use the sidebar rail). -->
  <header
    class="sticky top-0 z-30 border-b border-border bg-surface/95 pt-[env(safe-area-inset-top)] backdrop-blur transition-shadow md:hidden"
    :class="scrolled ? 'shadow-sm' : ''"
  >
    <div class="flex h-14 items-center justify-between gap-3 px-4">
      <!-- 44px tall tap area, nudged left so the logo still lines up with the page content -->
      <NuxtLink
        to="/dashboard"
        class="-ml-1 flex min-h-11 min-w-0 items-center rounded-lg px-1 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
      >
        <BrandLogo :size="28" />
      </NuxtLink>
      <div class="flex shrink-0 items-center">
        <LayoutUserlayoutUserMenu :profile="profile" placement="down" />
      </div>
    </div>
  </header>
</template>