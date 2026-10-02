<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { meKey, type Profile } from '~/types'

/** Admin shell: brand, Admin badge, avatar menu, section tabs. */
const api = useApi()
const { data: profile } = useQuery({
  queryKey: meKey(),
  queryFn: () => api.get<Profile>('/me'),
  staleTime: 5 * 60_000,
})

const route = useRoute()
const tabs = [
  { label: 'Overview', to: '/admin' },
  { label: 'Users', to: '/admin/users' },
  { label: 'Journals', to: '/admin/journals' },
  { label: 'Settings', to: '/admin/settings' },
  { label: 'Audit', to: '/admin/audit' },
]

function tabActive(to: string): boolean {
  if (to === '/admin')
    return route.path === '/admin'
  return route.path === to || route.path.startsWith(`${to}/`)
}
</script>

<template>
  <div class="min-h-screen bg-bg text-text">
    <header class="sticky top-0 z-30 border-b border-border bg-surface/95 backdrop-blur">
      <div class="mx-auto flex h-14 w-full max-w-6xl items-center justify-between px-4 md:px-6">
        <NuxtLink to="/admin" class="flex items-center">
          <BrandLogo :size="28" />
        </NuxtLink>
        <div class="flex items-center gap-2">
          <span class="rounded-full bg-warning/10 px-2.5 py-1 text-xs font-semibold text-warning-text">Admin</span>
          <LayoutUserlayoutUserMenu :profile="profile ?? null" placement="down" />
        </div>
      </div>
      <nav class="mx-auto w-full max-w-6xl overflow-x-auto px-4 md:px-6" aria-label="Admin sections">
        <div class="flex gap-1 pb-2">
          <NuxtLink
            v-for="tab in tabs"
            :key="tab.to"
            :to="tab.to"
            :aria-current="tabActive(tab.to) ? 'page' : undefined"
            :class="[
              'whitespace-nowrap rounded-lg px-3 py-1.5 text-sm font-medium transition',
              tabActive(tab.to) ? 'bg-primary/10 text-primary' : 'text-muted hover:text-text',
            ]"
          >
            {{ tab.label }}
          </NuxtLink>
        </div>
      </nav>
    </header>
    <main class="mx-auto w-full max-w-6xl px-4 pb-10 pt-4 md:px-6 md:pt-6">
      <slot />
    </main>
  </div>
</template>
