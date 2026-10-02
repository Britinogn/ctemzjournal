<script setup lang="ts">
import {
  DashboardSquare02Icon,
  GlobeIcon,
  Settings01Icon,
  ShieldCheckIcon,
  UsersIcon,
} from '~/utils/icons'
import type { Profile } from '~/types'

defineProps<{ profile: Profile | null }>()

const route = useRoute()

const nav = [
  { label: 'Overview', to: '/admin', icon: DashboardSquare02Icon, exact: true },
  { label: 'Users', to: '/admin/users', icon: UsersIcon, exact: false },
  { label: 'Journals', to: '/admin/journals', icon: GlobeIcon, exact: false },
  { label: 'Site settings', to: '/admin/settings', icon: Settings01Icon, exact: false },
  { label: 'Audit log', to: '/admin/audit', icon: ShieldCheckIcon, exact: false },
]

function isActive(item: { to: string; exact: boolean }): boolean {
  if (item.exact)
    return route.path === item.to
  return route.path === item.to || route.path.startsWith(`${item.to}/`)
}
</script>

<template>
  <aside class="fixed inset-y-0 left-0 z-30 hidden w-20 flex-col border-r border-border bg-surface md:flex xl:w-64">
    <div class="flex h-16 items-center px-4 xl:px-5">
      <NuxtLink to="/admin" class="flex items-center">
        <BrandLogo :show-name="false" :size="34" />
        <span class="ml-2 hidden text-lg font-bold tracking-tight xl:inline">Ctemz Journal</span>
      </NuxtLink>
    </div>

    <div class="px-3 xl:px-4">
      <span class="block rounded-xl bg-warning/10 px-3 py-2 text-center text-sm font-semibold text-warning-text">
        <span class="xl:hidden">A</span><span class="hidden xl:inline">Admin</span>
      </span>
    </div>

    <p class="mt-4 hidden px-5 text-xs font-medium uppercase tracking-wide text-muted xl:block">
      Admin area
    </p>
    <nav class="mt-1 flex-1 space-y-1 overflow-y-auto px-3 xl:px-4" aria-label="Admin">
      <NuxtLink
        v-for="item in nav"
        :key="item.to"
        :to="item.to"
        :aria-current="isActive(item) ? 'page' : undefined"
        :class="[
          'flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition',
          isActive(item)
            ? 'bg-primary/10 text-primary'
            : 'text-muted hover:bg-bg hover:text-text',
        ]"
      >
        <UiAppIcon :icon="item.icon" :size="20" class="shrink-0" />
        <span class="hidden xl:inline">{{ item.label }}</span>
      </NuxtLink>
    </nav>

    <div class="p-3 xl:p-4">
      <LayoutUserlayoutUserCard :profile="profile" settings-to="/admin/settings" />
    </div>
  </aside>
</template>
