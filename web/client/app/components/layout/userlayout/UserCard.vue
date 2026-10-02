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
  <div class="flex items-center gap-3 rounded-2xl border border-border bg-surface p-3">
    <LayoutUserlayoutUserAvatar
      :name="profile?.DisplayName"
      :email="undefined"
      :size="40"
    />
    <div class="min-w-0 flex-1 leading-tight">
      <p class="truncate text-sm font-semibold">
        {{ profile?.DisplayName || 'Trader' }}
      </p>
      <p class="text-xs capitalize text-muted">
        {{ profile?.Role || 'user' }}
      </p>
    </div>
    <button
      type="button"
      aria-label="Log out"
      :disabled="busy"
      class="inline-flex h-9 w-9 shrink-0 items-center justify-center rounded-xl border border-border text-muted transition hover:border-loss hover:text-loss disabled:opacity-50 xl:h-10 xl:w-10"
      @click="onLogout"
    >
      <UiAppIcon :icon="Logout02Icon" :size="18" />
    </button>
  </div>
</template>
