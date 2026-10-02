<script setup lang="ts">
import {
  DashboardSquare02Icon,
  GlobeIcon,
  Settings01Icon,
  ShieldCheckIcon,
  UsersIcon,
} from '~/utils/icons'

const route = useRoute()

const tabs = [
  { label: 'Overview', to: '/admin', icon: DashboardSquare02Icon, exact: true },
  { label: 'Users', to: '/admin/users', icon: UsersIcon, exact: false },
  { label: 'Journals', to: '/admin/journals', icon: GlobeIcon, exact: false },
  { label: 'Audit', to: '/admin/audit', icon: ShieldCheckIcon, exact: false },
  { label: 'Site', to: '/admin/settings', icon: Settings01Icon, exact: false },
]

function active(tab: { to: string; exact: boolean }): boolean {
  if (tab.exact)
    return route.path === tab.to
  return route.path === tab.to || route.path.startsWith(`${tab.to}/`)
}
</script>

<template>
  <!-- Mobile bottom tabs (tablet and desktop use the admin sidebar). -->
  <nav class="fixed inset-x-0 bottom-0 z-30 border-t border-border bg-surface/95 backdrop-blur md:hidden" aria-label="Admin">
    <div class="grid grid-cols-5 px-2 pb-[env(safe-area-inset-bottom)]">
      <NuxtLink
        v-for="tab in tabs"
        :key="tab.to"
        :to="tab.to"
        :class="[
          'flex flex-col items-center gap-0.5 py-2 text-[11px] font-medium transition',
          active(tab) ? 'text-primary' : 'text-muted',
        ]"
      >
        <UiAppIcon :icon="tab.icon" :size="22" />
        {{ tab.label }}
      </NuxtLink>
    </div>
  </nav>
</template>
