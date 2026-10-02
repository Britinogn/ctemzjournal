<script setup lang="ts">
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
  <div class="rounded-2xl border border-warning/40 bg-surface p-6 text-center shadow-sm sm:p-8">
    <h1 class="text-2xl font-bold tracking-tight">
      Account suspended
    </h1>
    <p class="mt-2 text-sm text-muted">
      This account has been suspended. Contact support if you think this is a mistake.
    </p>
    <button
      type="button"
      :disabled="loading"
      class="mt-6 w-full rounded-xl border border-border py-2.5 text-sm font-semibold transition hover:border-primary disabled:opacity-50"
      @click="onLogout"
    >
      {{ loading ? 'Signing out…' : 'Sign out' }}
    </button>
  </div>
</template>
