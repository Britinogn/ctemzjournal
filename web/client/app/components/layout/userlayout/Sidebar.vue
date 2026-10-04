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

// "group" only controls where the dividers go: daily use, the lists you manage, then settings.
const nav = [
  { label: 'Dashboard', to: '/dashboard', icon: DashboardSquare02Icon, exact: true, group: 0 },
  { label: 'Trades', to: '/dashboard/trades', icon: TradeUpIcon, exact: false, group: 0 },
  { label: 'Accounts', to: '/dashboard/accounts', icon: Wallet01Icon, exact: false, group: 1 },
  { label: 'Setups', to: '/dashboard/setups', icon: Target01Icon, exact: false, group: 1 },
  { label: 'Tags', to: '/dashboard/tags', icon: Tag01Icon, exact: false, group: 1 },
  { label: 'Settings', to: '/dashboard/settings', icon: Settings01Icon, exact: false, group: 2 },
]

function isActive(item: { to: string, exact: boolean }): boolean {
  if (item.exact)
    return route.path === item.to
  return route.path === item.to || route.path.startsWith(`${item.to}/`)
}

const isAdminRoute = computed(() => route.path.startsWith('/admin'))

const focusRing
  = 'focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary'
</script>

<template>
  <aside class="fixed inset-y-0 left-0 z-30 hidden w-20 flex-col border-r border-border bg-surface md:flex xl:w-64" aria-label="Sidebar">
    <div class="flex h-16 items-center justify-center px-3 xl:justify-start xl:px-5">
      <NuxtLink to="/dashboard" class="flex items-center rounded-lg" :class="focusRing">
        <BrandLogo :show-name="false" :size="34" />
        <!-- max-xl:sr-only keeps the name for screen readers when the sidebar is the icon rail -->
        <span class="ml-2 text-lg font-bold tracking-tight max-xl:sr-only">Ctemz Journal</span>
      </NuxtLink>
    </div>

    <!-- Role badge: static "Trader" for traders, Trader/Admin switcher for admins -->
    <div class="px-3 xl:px-4">
      <div class="grid grid-cols-1 gap-1 rounded-2xl border border-border bg-bg p-1">
        <NuxtLink
          to="/dashboard"
          :class="[
            'rounded-xl px-3 py-2 text-center text-sm font-semibold transition-colors',
            focusRing,
            !isAdminRoute ? 'bg-surface text-text shadow-sm ring-1 ring-border' : 'text-muted hover:text-text',
          ]"
        >
          <span class="xl:hidden" aria-hidden="true">T</span>
          <span class="max-xl:sr-only">Trader</span>
        </NuxtLink>
      </div>
    </div>

    <nav class="mt-4 flex-1 space-y-1 overflow-y-auto px-3 xl:px-4" aria-label="Primary">
      <template v-for="(item, i) in nav" :key="item.to">
        <div v-if="i > 0 && nav[i - 1]!.group !== item.group" class="!my-3 border-t border-border" aria-hidden="true" />
        <NuxtLink
          :to="item.to"
          :title="item.label"
          :aria-current="isActive(item) ? 'page' : undefined"
          :class="[
            'relative flex h-11 items-center justify-center gap-3 rounded-xl text-sm font-medium transition-colors xl:justify-start xl:px-3',
            focusRing,
            isActive(item)
              ? 'bg-primary/10 font-semibold text-primary'
              : 'text-muted hover:bg-bg hover:text-text',
          ]"
        >
          <!-- Edge marker, so the current page is not shown by colour alone -->
          <span
            v-if="isActive(item)"
            class="absolute inset-y-2 -left-3 w-1 rounded-r-full bg-primary xl:-left-4"
            aria-hidden="true"
          />
          <UiAppIcon :icon="item.icon" :size="20" class="shrink-0" />
          <span class="max-xl:sr-only">{{ item.label }}</span>
        </NuxtLink>
      </template>
    </nav>

    <div class="space-y-3 p-3 xl:p-4">
      <NuxtLink
        to="/dashboard/trades/new"
        title="Log trade"
        :class="['flex h-11 items-center justify-center gap-2 rounded-xl bg-primary text-sm font-semibold text-on-primary transition-colors hover:opacity-90', focusRing]"
      >
        <UiAppIcon :icon="Add01Icon" :size="18" />
        <span class="max-xl:sr-only">Log trade</span>
      </NuxtLink>
      <LayoutUserlayoutUserCard :profile="profile" />
    </div>
  </aside>
</template>