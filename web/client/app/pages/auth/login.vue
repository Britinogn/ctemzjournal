<script setup lang="ts">
definePageMeta({ middleware: 'guest' })

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref<string | null>(null)
const { login } = useAuth()

async function onSubmit(): Promise<void> {
  if (loading.value)
    return
  loading.value = true
  error.value = null
  try {
    await login(email.value.trim(), password.value)
  }
  catch (err) {
    error.value = err instanceof Error ? err.message : 'Login failed'
  }
  finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="rounded-2xl border border-border bg-surface p-6 shadow-sm sm:p-8">
    <h1 class="text-2xl font-bold tracking-tight">
      Welcome back
    </h1>
    <p class="mt-1 text-sm text-muted">
      Log in to your trading journal
    </p>

    <form class="mt-6 space-y-4" @submit.prevent="onSubmit">
      <div>
        <label for="login-email" class="mb-1.5 block text-sm font-medium">Email</label>
        <input
          id="login-email"
          v-model="email"
          type="email"
          required
          autocomplete="email"
          placeholder="you@example.com"
          class="w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary"
        >
      </div>
      <div>
        <div class="mb-1.5 flex items-center justify-between">
          <label for="login-password" class="block text-sm font-medium">Password</label>
          <NuxtLink to="/auth/forgot-password" class="text-sm text-primary hover:underline">
            Forgot password?
          </NuxtLink>
        </div>
        <input
          id="login-password"
          v-model="password"
          type="password"
          required
          autocomplete="current-password"
          placeholder="••••••••"
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
        {{ loading ? 'Logging in…' : 'Log in' }}
      </button>
    </form>

    <p class="mt-6 text-center text-sm text-muted">
      No account?
      <NuxtLink to="/auth/signup" class="font-medium text-primary hover:underline">
        Sign up
      </NuxtLink>
    </p>
  </div>
</template>
