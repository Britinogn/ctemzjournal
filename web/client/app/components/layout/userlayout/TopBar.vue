<script setup lang="ts">
import type { Profile } from '~/types'

defineProps<{ profile: Profile | null }>()
</script>

<template>
  <!-- Mobile top bar (tablet and desktop use the sidebar rail). -->
  <header class="sticky top-0 z-30 border-b border-border bg-surface/95 backdrop-blur md:hidden">
    <div class="flex h-14 items-center justify-between px-4">
      <NuxtLink to="/dashboard" class="flex items-center">
        <BrandLogo :size="28" />
      </NuxtLink>
      <LayoutUserlayoutUserAvatar :name="profile?.DisplayName" :size="34" />
    </div>
  </header>
</template>
