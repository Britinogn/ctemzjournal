<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { meKey, type Profile } from '~/types'

/** Admin shell: sidebar (tablet/desktop), one top bar for every size,
 *  and bottom tabs (phone). */
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
    <!-- Bottom padding clears the phone tabs, including the home-indicator area -->
    <div class="pb-[calc(6rem+env(safe-area-inset-bottom))] md:pb-10 md:pl-20 xl:pl-64">
      <LayoutAdminlayoutAdminTopbar :profile="profile ?? null" />
      <main class="mx-auto w-full max-w-6xl px-4 pt-4 md:px-6 md:pt-6">
        <slot />
      </main>
    </div>
    <LayoutAdminlayoutAdminBottomNav />
  </div>
</template>