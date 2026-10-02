<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { meKey, type Profile } from '~/types'

/** Trader shell: sidebar rail (tablet/desktop), top bar + bottom nav (mobile). */
const api = useApi()
const { data: profile } = useQuery({
  queryKey: meKey(),
  queryFn: () => api.get<Profile>('/me'),
  staleTime: 5 * 60_000,
})
</script>

<template>
  <div class="min-h-screen bg-bg text-text">
    <LayoutUserlayoutSidebar :profile="profile ?? null" />
    <LayoutUserlayoutTopBar :profile="profile ?? null" />
    <div class="pb-24 md:pb-10 md:pl-20 xl:pl-64">
      <main class="mx-auto w-full max-w-6xl px-4 pt-4 md:px-6 md:pt-6">
        <slot />
      </main>
    </div>
    <LayoutUserlayoutBottomNav />
  </div>
</template>
