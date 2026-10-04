<script setup lang="ts">
import {
  Mail01Icon,
  LockKeyIcon,
  EyeIcon,
  EyeOffIcon,
  Loading03Icon,
  ArrowLeft01Icon,
  UserIcon,
} from '~/utils/icons'

definePageMeta({ middleware: 'guest' })

const email = ref('')
const password = ref('')
const showPassword = ref(false)
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
  <div class="w-full max-w-md">
    <div class="rounded-2xl border border-border bg-surface p-7 shadow-sm sm:p-8">

      <!-- Header -->
      <div class="mb-7">
        <div class="mb-4 flex h-11 w-11 items-center justify-center rounded-xl bg-primary/10 text-primary">
          <UiAppIcon :icon="UserIcon" :size="18" />
        </div>
        <h1 class="text-2xl font-bold tracking-tight">
          Welcome back
        </h1>
        <p class="mt-1.5 text-sm text-muted">
          Log in to your trading journal
        </p>
      </div>

      <!-- Form -->
      <form class="space-y-5" @submit.prevent="onSubmit">
        <!-- Email -->
        <div>
          <label for="login-email" class="mb-1.5 block text-sm font-medium">
            Email
          </label>
          <div class="relative">
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5 text-muted">
              <UiAppIcon :icon="Mail01Icon" :size="18" />
            </div>
            <input
              id="login-email"
              v-model="email"
              type="email"
              required
              autocomplete="email"
              placeholder="you@example.com"
              class="w-full rounded-xl border border-border bg-bg py-3 pl-11 pr-4 text-sm outline-none transition
                placeholder:text-muted
                focus:border-primary focus:ring-2 focus:ring-primary/20"
            >
          </div>
        </div>

        <!-- Password -->
        <div>
          <div class="mb-1.5 flex items-center justify-between">
            <label for="login-password" class="block text-sm font-medium">
              Password
            </label>
            <NuxtLink
              to="/auth/forgot-password"
              class="text-sm font-medium text-primary hover:underline"
            >
              Forgot password?
            </NuxtLink>
          </div>
          <div class="relative">
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5 text-muted">
              <UiAppIcon :icon="LockKeyIcon" :size="18" />
            </div>
            <input
              id="login-password"
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              required
              autocomplete="current-password"
              placeholder="••••••••"
              class="w-full rounded-xl border border-border bg-bg py-3 pl-11 pr-12 text-sm outline-none transition
                placeholder:text-muted
                focus:border-primary focus:ring-2 focus:ring-primary/20"
            >
            <button
              type="button"
              class="absolute inset-y-0 right-0 flex items-center pr-3.5 text-muted transition hover:text-foreground"
              @click="showPassword = !showPassword"
            >
              <UiAppIcon
                :icon="showPassword ? EyeOffIcon : EyeIcon"
                :size="18"
              />
            </button>
          </div>
        </div>

        <!-- Error -->
        <p
          v-if="error"
          role="alert"
          class="rounded-xl bg-loss/10 px-4 py-2.5 text-sm text-loss"
        >
          {{ error }}
        </p>

        <!-- Submit -->
        <button
          type="submit"
          :disabled="loading"
          class="relative flex w-full items-center justify-center gap-2 rounded-xl bg-primary py-3 text-sm font-semibold text-on-primary
            transition hover:opacity-90
            disabled:cursor-not-allowed disabled:opacity-50"
        >
          <UiAppIcon
            v-if="loading"
            :icon="Loading03Icon"
            :size="18"
            class="absolute animate-spin"
          />
          <span :class="{ 'opacity-0': loading }">
            Log in
          </span>
        </button>
      </form>

      <!-- Footer links -->
      <div class="mt-7 space-y-3 text-center text-sm">
        <p class="text-muted">
          No account?
          <NuxtLink to="/auth/signup" class="font-medium text-primary hover:underline">
            Sign up
          </NuxtLink>
        </p>

        <NuxtLink
          to="/"
          class="inline-flex items-center gap-1.5 font-medium text-muted transition hover:text-primary"
        >
          <UiAppIcon :icon="ArrowLeft01Icon" :size="16" />
          Back to home
        </NuxtLink>
      </div>
    </div>
  </div>
</template>