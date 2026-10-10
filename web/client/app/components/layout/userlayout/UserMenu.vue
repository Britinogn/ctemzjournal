<script setup lang="ts">
import { onClickOutside } from '@vueuse/core'
import { Logout02Icon, Moon02Icon, Sun03Icon } from '~/utils/icons'
import type { Profile } from '~/types'

withDefaults(defineProps<{
  profile: Profile | null
  /** Sidebar opens upward, top bars drop downward. */
  placement?: 'up' | 'down'
}>(), {
  placement: 'down',
})

const { logout } = useAuth()
const colorMode = useColorMode()
const isDark = computed(() => colorMode.value === 'dark')

const open = ref(false)
const busy = ref(false)
const menu = ref<HTMLElement | null>(null)
const trigger = ref<HTMLButtonElement | null>(null)

function toggleTheme(): void {
  colorMode.preference = isDark.value ? 'light' : 'dark'
  open.value = false
}

function close(returnFocus = false): void {
  open.value = false
  if (returnFocus)
    trigger.value?.focus()
}

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

const itemCls
  = 'flex h-11 w-full items-center gap-2.5 rounded-xl px-3 text-left text-sm font-medium transition-colors hover:bg-bg focus-visible:outline-2-2 focus-visible:-outline-offset-2 focus-visible:outline-2-primary disabled:cursor-wait disabled:opacity-60'
</script>

<template>
  <div ref="menu" class="relative" @keydown.escape="close(true)">
    <!-- 44px tap area around the 36px avatar -->
    <button
      ref="trigger"
      type="button"
      aria-haspopup="menu"
      :aria-expanded="open"
      aria-label="Account menu"
      class="inline-flex h-11 w-11 items-center justify-center rounded-full outline-none transition-shadow focus-visible:ring-2 focus-visible:ring-primary"
      :class="open ? 'ring-2 ring-primary/40' : ''"
      @click="open = !open"
    >
      <LayoutUserlayoutUserAvatar :name="profile?.DisplayName" :src="profile?.AvatarPath" :rev="profile?.UpdatedAt" :size="36" />
    </button>

    <!-- Settings is not here: the sidebar tab and the bottom "More" tab already lead to it. -->
    <div
      v-if="open"
      role="menu"
      :class="[
        'absolute right-0 z-50 w-60 overflow-hidden rounded-2xl border border-border bg-surface shadow-xl',
        placement === 'up' ? 'bottom-full mb-2' : 'top-full mt-2',
      ]"
    >
      <div class="flex items-center gap-3 border-b border-border px-4 py-3">
        <LayoutUserlayoutUserAvatar :name="profile?.DisplayName" :src="profile?.AvatarPath" :rev="profile?.UpdatedAt" :size="36" />
        <div class="min-w-0">
          <p class="truncate text-sm font-semibold leading-tight">
            {{ profile?.DisplayName || 'Trader' }}
          </p>
          <p
            class="text-xs font-medium capitalize"
            :class="profile?.Role === 'admin' ? 'text-warning-text' : 'text-muted'"
          >
            {{ profile?.Role || 'user' }}
          </p>
        </div>
      </div>
      <div class="p-1.5">
        <button type="button" role="menuitem" :class="itemCls" @click="toggleTheme">
          <UiAppIcon :icon="isDark ? Sun03Icon : Moon02Icon" :size="18" class="text-muted" />
          {{ isDark ? 'Light mode' : 'Dark mode' }}
        </button>
        <button type="button" role="menuitem" :disabled="busy" :aria-busy="busy" :class="itemCls" @click="onLogout">
          <UiAppIcon :icon="Logout02Icon" :size="18" class="text-muted" />
          {{ busy ? 'Signing out…' : 'Log out' }}
        </button>
      </div>
    </div>
  </div>
</template>