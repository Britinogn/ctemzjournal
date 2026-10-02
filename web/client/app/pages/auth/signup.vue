<script setup lang="ts">
definePageMeta({ middleware: 'guest' })

const displayName = ref('')
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref<string | null>(null)
const { signup } = useAuth()

async function onSubmit(): Promise<void> {
  if (loading.value)
    return
  loading.value = true
  error.value = null
  try {
    await signup(email.value.trim(), password.value, displayName.value.trim() || undefined)
  }
  catch (err) {
    error.value = err instanceof Error ? err.message : 'Signup failed'
  }
  finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="rounded-2xl border border-border bg-surface p-6 shadow-sm sm:p-8">
    <h1 class="text-2xl font-bold tracking-tight">
      Create your journal
    </h1>
    <p class="mt-1 text-sm text-muted">
      Start tracking every trade
    </p>

    <form class="mt-6 space-y-4" @submit.prevent="onSubmit">
      <div>
        <label for="signup-name" class="mb-1.5 block text-sm font-medium">Display name</label>
        <input
          id="signup-name"
          v-model="displayName"
          type="text"
          autocomplete="nickname"
          placeholder="Ada Obi"
          class="w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary"
        >
      </div>
      <div>
        <label for="signup-email" class="mb-1.5 block text-sm font-medium">Email</label>
        <input
          id="signup-email"
          v-model="email"
          type="email"
          required
          autocomplete="email"
          placeholder="you@example.com"
          class="w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary"
        >
      </div>
      <div>
        <label for="signup-password" class="mb-1.5 block text-sm font-medium">Password</label>
        <input
          id="signup-password"
          v-model="password"
          type="password"
          required
          minlength="6"
          autocomplete="new-password"
          placeholder="Minimum 6 characters"
          class="w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary"
        >
      </div>

      <p v-if="error" role="alert" class="rounded-xl bg-loss/10 px-4 py-2.5 text-sm text-loss">
        {{ error }}
      </p>

      <button
        type="submit"
        :disabled="loading"
        class="w-full rounded-xl bg-primary py-2.5 text-sm font-semibold text-on-primary transition hover:opacity-90 disabled:opacity-50"
      >
        {{ loading ? 'Creating account…' : 'Sign up' }}
      </button>
    </form>

    <p class="mt-6 text-center text-sm text-muted">
      Have an account?
      <NuxtLink to="/auth/login" class="font-medium text-primary hover:underline">
        Log in
      </NuxtLink>
    </p>
  </div>
</template>
