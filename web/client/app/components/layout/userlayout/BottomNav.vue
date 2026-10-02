<script setup lang="ts">
import {
  Add01Icon,
  Home01Icon,
  Menu01Icon,
  TradeUpIcon,
  Wallet01Icon,
} from '~/utils/icons'

const route = useRoute()

const tabs = [
  { label: 'Home', to: '/dashboard', icon: Home01Icon, exact: true },
  { label: 'Trades', to: '/dashboard/trades', icon: TradeUpIcon, exact: false },
  { label: 'Log', to: '/dashboard/trades/new', icon: Add01Icon, fab: true, exact: true },
  { label: 'Accounts', to: '/dashboard/accounts', icon: Wallet01Icon, exact: false },
  { label: 'More', to: '/dashboard/settings', icon: Menu01Icon, exact: false },
]

function active(tab: { to: string; exact: boolean }): boolean {
  if (tab.exact)
    return route.path === tab.to
  return route.path === tab.to || route.path.startsWith(`${tab.to}/`)
}
</script>

<template>
  <!-- Mobile bottom nav (tablet and desktop use the sidebar rail). -->
  <nav class="fixed inset-x-0 bottom-0 z-30 border-t border-border bg-surface/95 backdrop-blur md:hidden" aria-label="Primary">
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
        <span
          v-if="tab.fab"
          class="-mt-6 inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-primary text-on-primary shadow-lg"
        >
          <UiAppIcon :icon="tab.icon" :size="22" />
        </span>
        <UiAppIcon v-else :icon="tab.icon" :size="22" />
        {{ tab.label }}
      </NuxtLink>
    </div>
  </nav>
</template>
