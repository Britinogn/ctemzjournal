<script setup lang="ts">
import { onClickOutside } from '@vueuse/core'
import { Logout02Icon, Settings01Icon } from '~/utils/icons'
import type { Profile } from '~/types'

withDefaults(defineProps<{
  profile: Profile | null;
  /** Admin shell points this at /admin/settings instead. */
  settingsTo?: string;
}>(), {
  settingsTo: '/dashboard/settings',
})

const { logout } = useAuth()
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
      class="flex w-full items-center gap-3 rounded-2xl border border-border bg-surface p-3 text-left transition hover:border-primary"
      @click="open = !open"
      @keydown.escape="open = false"
    >
      <LayoutUserlayoutUserAvatar
        :name="profile?.DisplayName"
        :email="undefined"
        :size="40"
      />
      <div class="hidden min-w-0 flex-1 leading-tight xl:block">
        <p class="truncate text-sm font-semibold">
          {{ profile?.DisplayName || 'Trader' }}
        </p>
        <p class="text-xs capitalize text-muted">
          {{ profile?.Role || 'user' }}
        </p>
      </div>
    </button>

    <div
      v-if="open"
      role="menu"
      class="absolute inset-x-0 bottom-full z-50 mb-2 overflow-hidden rounded-2xl border border-border bg-surface shadow-xl"
    >
      <div class="p-1.5">
        <NuxtLink
          :to="settingsTo"
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
