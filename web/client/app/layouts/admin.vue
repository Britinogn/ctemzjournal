<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { meKey, type Profile } from '~/types'

/** Minimal admin shell (full admin area comes later): brand + avatar menu. */
const api = useApi()
const { data: profile } = useQuery({
  queryKey: meKey(),
  queryFn: () => api.get<Profile>('/me'),
  staleTime: 5 * 60_000,
})
</script>

<template>
  <div class="min-h-screen bg-bg text-text">
    <header class="sticky top-0 z-30 border-b border-border bg-surface/95 backdrop-blur">
      <div class="mx-auto flex h-14 w-full max-w-6xl items-center justify-between px-4 md:px-6">
        <NuxtLink to="/admin/settings" class="flex items-center">
          <BrandLogo :size="28" />
        </NuxtLink>
        <div class="flex items-center gap-2">
          <span class="rounded-full bg-warning/10 px-2.5 py-1 text-xs font-semibold text-warning-text">Admin</span>
          <LayoutUserlayoutUserMenu :profile="profile ?? null" placement="down" />
        </div>
      </div>
    </header>
    <main class="mx-auto w-full max-w-6xl px-4 pb-10 pt-4 md:px-6 md:pt-6">
      <slot />
    </main>
  </div>
</template>
