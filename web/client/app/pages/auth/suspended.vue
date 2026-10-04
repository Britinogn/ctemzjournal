<script setup lang="ts">
import {
  Alert02Icon,
  Loading03Icon,
  Logout02Icon,
} from '~/utils/icons'

const { logout } = useAuth()
const loading = ref(false)

async function onLogout(): Promise<void> {
  loading.value = true
  try {
    await logout()
  }
  finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="w-full max-w-md">
    <div class="rounded-2xl border border-warning/30 bg-surface p-7 text-center shadow-sm sm:p-8">

      <!-- Header -->
      <div class="mb-6">
        <div class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-warning/10 text-warning">
          <UiAppIcon :icon="Alert02Icon" :size="22" />
        </div>
        <h1 class="text-2xl font-bold tracking-tight">
          Account suspended
        </h1>
        <p class="mt-2 text-sm text-muted">
          This account has been suspended.
          Contact support if you think this is a mistake.
        </p>
      </div>

      <!-- Action -->
      <button
        type="button"
        :disabled="loading"
        class="relative flex w-full items-center justify-center gap-2 rounded-xl border border-border py-3 text-sm font-semibold transition
          hover:border-primary hover:text-primary
          disabled:cursor-not-allowed disabled:opacity-50"
        @click="onLogout"
      >
        <UiAppIcon
          v-if="loading"
          :icon="Loading03Icon"
          :size="18"
          class="absolute animate-spin"
        />
        <span :class="{ 'opacity-0': loading }" class="inline-flex cursor-pointer items-center gap-2">
          <UiAppIcon :icon="Logout02Icon" :size="16" />
          Sign out
        </span>
      </button>
    </div>
  </div>
</template>