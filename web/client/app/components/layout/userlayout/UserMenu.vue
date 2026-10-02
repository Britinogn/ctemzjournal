<script setup lang="ts">
import { onClickOutside } from '@vueuse/core'
import { Logout02Icon, Moon02Icon, Settings01Icon, Sun03Icon } from '~/utils/icons'
import type { Profile } from '~/types'

withDefaults(defineProps<{
  profile: Profile | null;
  /** Sidebar opens upward, top bars drop downward. */
  placement?: 'up' | 'down';
}>(), {
  placement: 'down',
})

const { logout } = useAuth()
const colorMode = useColorMode()
const isDark = computed(() => colorMode.value === 'dark')

function toggleTheme(): void {
  colorMode.preference = isDark.value ? 'light' : 'dark'
  open.value = false
}
const open = ref(false)
const busy = ref(false)
const menu = ref<HTMLElement | null>(null)

onClickOutside(menu, () => {
  open.value = false
})

async function onLogout(): Promise<void> {
  if (busy.value)
    return
  busy.value = true
  try {
    await logout()
  }
  finally {
    busy.value = false
    open.value = false
  }
}
</script>

<template>
  <div ref="menu" class="relative">
    <button
      type="button"
      aria-haspopup="menu"
      :aria-expanded="open"
      aria-label="Account menu"
      class="rounded-full outline-none transition focus-visible:ring-2 focus-visible:ring-primary"
      @click="open = !open"
      @keydown.escape="open = false"
    >
      <LayoutUserlayoutUserAvatar :name="profile?.DisplayName" :size="36" />
    </button>

    <div
      v-if="open"
      role="menu"
      :class="[
        'absolute right-0 z-50 w-56 overflow-hidden rounded-2xl border border-border bg-surface shadow-xl',
        placement === 'up' ? 'bottom-full mb-2' : 'top-full mt-2',
      ]"
    >
      <div class="border-b border-border px-4 py-3">
        <p class="truncate text-sm font-semibold">
          {{ profile?.DisplayName || 'Trader' }}
        </p>
        <p class="text-xs capitalize text-muted">
          {{ profile?.Role || 'user' }}
        </p>
      </div>
      <div class="p-1.5">
        <NuxtLink
          to="/dashboard/settings"
          role="menuitem"
          class="flex items-center gap-2.5 rounded-xl px-3 py-2 text-sm transition hover:bg-bg"
          @click="open = false"
        >
          <UiAppIcon :icon="Settings01Icon" :size="18" class="text-muted" />
          Settings
        </NuxtLink>
        <button
          type="button"
          role="menuitem"
          class="flex w-full items-center gap-2.5 rounded-xl px-3 py-2 text-left text-sm transition hover:bg-bg"
          @click="toggleTheme"
        >
          <UiAppIcon :icon="isDark ? Sun03Icon : Moon02Icon" :size="18" class="text-muted" />
          {{ isDark ? 'Light mode' : 'Dark mode' }}
        </button>
        <button
          type="button"
          role="menuitem"
          :disabled="busy"
          class="flex w-full items-center gap-2.5 rounded-xl px-3 py-2 text-left text-sm text-loss transition hover:bg-loss/10 disabled:opacity-50"
          @click="onLogout"
        >
          <UiAppIcon :icon="Logout02Icon" :size="18" />
          {{ busy ? 'Signing out…' : 'Log out' }}
        </button>
      </div>
    </div>
  </div>
</template>
