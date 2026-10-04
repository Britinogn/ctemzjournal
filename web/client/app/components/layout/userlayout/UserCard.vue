<script setup lang="ts">
import { Logout02Icon } from '~/utils/icons'
import type { Profile } from '~/types'

defineProps<{ profile: Profile | null }>()

const { logout } = useAuth()
const busy = ref(false)

async function onLogout(): Promise<void> {
  if (busy.value)
    return
  busy.value = true
  try {
    await logout()
  }
  finally {
    busy.value = false
  }
}
</script>

<template>
  <!--
    No popover any more. With Settings gone (the sidebar already has a Settings tab) the menu held one item,
    so Log out is its own button.
    Tablet (icon rail): avatar above an icon-only logout button.
    Desktop: one card with the name and role, and a Log out button under it.
  -->
  <div
    class="flex flex-col items-center gap-2 xl:items-stretch xl:gap-3 xl:rounded-2xl xl:border xl:border-border xl:bg-bg xl:p-3"
    role="group"
    aria-label="Account"
  >
    <div class="flex items-center gap-3" :title="profile?.DisplayName || 'Trader'">
      <LayoutUserlayoutUserAvatar
        :name="profile?.DisplayName"
        :email="undefined"
        :size="40"
        class="shrink-0"
      />
      <div class="hidden min-w-0 flex-1 leading-tight xl:block">
        <p class="truncate text-sm font-semibold">
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

    <!-- The text stays for screen readers on tablet, where only the icon shows -->
    <button
      type="button"
      title="Log out"
      :disabled="busy"
      :aria-busy="busy"
      class="flex h-11 w-full items-center justify-center gap-2 rounded-xl border border-border text-muted transition-colors hover:bg-bg hover:text-text focus-visible:outline-2-2 focus-visible:outline-2-offset-2 focus-visible:outline-2-primary disabled:cursor-wait disabled:opacity-60 xl:h-10 xl:bg-surface xl:text-sm xl:font-semibold"
      @click="onLogout"
    >
      <UiAppIcon
        :icon="Logout02Icon"
        :size="20"
        :class="busy ? 'animate-pulse motion-reduce:animate-none' : ''"
      />
      <span class="max-xl:sr-only">{{ busy ? 'Signing out…' : 'Log out' }}</span>
    </button>
  </div>
</template>