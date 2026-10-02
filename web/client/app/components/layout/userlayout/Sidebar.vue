<script setup lang="ts">
import {
  Add01Icon,
  DashboardSquare02Icon,
  Settings01Icon,
  Tag01Icon,
  Target01Icon,
  TradeUpIcon,
  Wallet01Icon,
} from '~/utils/icons'
import type { Profile } from '~/types'

defineProps<{ profile: Profile | null }>()

const route = useRoute()

const nav = [
  { label: 'Dashboard', to: '/dashboard', icon: DashboardSquare02Icon, exact: true },
  { label: 'Trades', to: '/dashboard/trades', icon: TradeUpIcon, exact: false },
  { label: 'Accounts', to: '/dashboard/accounts', icon: Wallet01Icon, exact: false },
  { label: 'Setups', to: '/dashboard/setups', icon: Target01Icon, exact: false },
  { label: 'Tags', to: '/dashboard/tags', icon: Tag01Icon, exact: false },
  { label: 'Settings', to: '/dashboard/settings', icon: Settings01Icon, exact: false },
]

function isActive(item: { to: string; exact: boolean }): boolean {
  if (item.exact)
    return route.path === item.to
  return route.path === item.to || route.path.startsWith(`${item.to}/`)
}

const isAdminRoute = computed(() => route.path.startsWith('/admin'))
</script>

<template>
  <aside class="fixed inset-y-0 left-0 z-30 hidden w-20 flex-col border-r border-border bg-surface md:flex xl:w-64">
    <div class="flex h-16 items-center px-4 xl:px-5">
      <NuxtLink to="/dashboard" class="flex items-center">
        <BrandLogo :show-name="false" :size="34" />
        <span class="ml-2 hidden text-lg font-bold tracking-tight xl:inline">Ctemz Journal</span>
      </NuxtLink>
    </div>

    <!-- Role badge: static "Trader" for traders, Trader/Admin switcher for admins -->
    <div class="px-3 xl:px-4">
      <div class="grid grid-cols-1 gap-1 rounded-2xl border border-border bg-bg p-1 xl:grid-cols-2">
        <NuxtLink
          to="/dashboard"
          :class="[
            'rounded-xl px-3 py-2 text-center text-sm font-medium transition',
            !isAdminRoute ? 'bg-surface text-text shadow-sm' : 'text-muted hover:text-text',
          ]"
        >
          <span class="xl:hidden">T</span><span class="hidden xl:inline">Trader</span>
        </NuxtLink>
        <NuxtLink
          v-if="profile?.Role === 'admin'"
          to="/admin"
          :class="[
            'rounded-xl px-3 py-2 text-center text-sm font-medium transition',
            isAdminRoute ? 'bg-surface text-text shadow-sm' : 'text-muted hover:text-text',
          ]"
        >
          <span class="xl:hidden">A</span><span class="hidden xl:inline">Admin</span>
        </NuxtLink>
      </div>
    </div>

    <nav class="mt-4 flex-1 space-y-1 overflow-y-auto px-3 xl:px-4" aria-label="Primary">
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

    <div class="space-y-3 p-3 xl:p-4">
      <NuxtLink
        to="/dashboard/trades/new"
        class="flex items-center justify-center gap-2 rounded-xl bg-primary py-2.5 text-sm font-semibold text-on-primary transition hover:opacity-90"
      >
        <UiAppIcon :icon="Add01Icon" :size="18" />
        <span class="hidden xl:inline">Log trade</span>
      </NuxtLink>
      <LayoutUserlayoutUserCard :profile="profile" />
    </div>
  </aside>
</template>
