<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { meKey, type Profile } from '~/types'

/** Admin shell: sidebar (tablet/desktop, no Trader toggle),
 *  top bar + bottom tabs (mobile). */
const api = useApi()
const { data: profile } = useQuery({
  queryKey: meKey(),
  queryFn: () => api.get<Profile>('/me'),
  staleTime: 5 * 60_000,
})
</script>

<template>
  <div class="min-h-screen bg-bg text-text">
    <LayoutAdminlayoutAdminSidebar :profile="profile ?? null" />
    <div class="pb-24 md:pb-10 md:pl-20 xl:pl-64">
      <!-- Mobile top bar -->
      <header class="sticky top-0 z-30 border-b border-border bg-surface/95 backdrop-blur md:hidden">
        <div class="flex h-14 items-center justify-between px-4">
          <NuxtLink to="/admin" class="flex items-center">
            <BrandLogo :size="28" />
          </NuxtLink>
          <div class="flex items-center gap-2">
            <span class="rounded-full bg-warning/10 px-2.5 py-1 text-xs font-semibold text-warning-text">Admin</span>
            <ThemeToggle />
          <LayoutUserlayoutUserMenu :profile="profile ?? null" placement="down" />
          </div>
        </div>
      </header>
      <main class="mx-auto w-full max-w-6xl px-4 pt-4 md:px-6 md:pt-6">
        <slot />
      </main>
    </div>
    <LayoutAdminlayoutAdminBottomNav />
  </div>
</template>
