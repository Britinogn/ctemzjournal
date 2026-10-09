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
usePageSeo('Log in', 'Log in to your Ctemz Journal account.')

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

/* ---------- Shared look ---------- */
const field
  = 'h-12 w-full rounded-xl border border-border bg-bg pl-11 text-sm outline-none transition-colors placeholder:text-muted focus:border-primary focus:ring-2 focus:ring-primary/20'
const focusRing
  = 'focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary'
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
              inputmode="email"
              autocapitalize="none"
              spellcheck="false"
              placeholder="you@example.com"
              :class="[field, 'pr-4']"
            >
          </div>
        </div>

        <!-- Password -->
        <div>
          <div class="mb-1.5 flex items-center justify-between">
            <label for="login-password" class="block text-sm font-medium">
              Password
            </label>
            <!-- The padding makes the tap area taller without moving the row -->
            <NuxtLink
              to="/auth/forgot-password"
              :class="['-my-2.5 rounded-md py-2.5 text-sm font-medium text-primary hover:underline', focusRing]"
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
              :class="[field, 'pr-12']"
            >
            <!-- Was an unnamed icon. Now it has a name, a state, and a full 48px tap area. -->
            <button
              type="button"
              :aria-label="showPassword ? 'Hide password' : 'Show password'"
              :aria-pressed="showPassword"
              class="absolute inset-y-0 right-0 flex w-12 items-center justify-center rounded-r-xl text-muted transition-colors hover:text-text focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-primary"
              @click="showPassword = !showPassword"
            >
              <UiAppIcon
                :icon="showPassword ? EyeOffIcon : EyeIcon"
                :size="18"
              />
            </button>
          </div>
        </div>

        <!-- Error: an icon next to the words, so it does not rely on colour alone -->
        <p
          v-if="error"
          role="alert"
          class="flex items-start gap-2 rounded-xl bg-loss/10 px-4 py-3 text-sm text-loss"
        >
          <svg class="mt-0.5 h-4 w-4 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M12 4l9 16H3L12 4zM12 10v4M12 17h.01" />
          </svg>
          {{ error }}
        </p>

        <!-- Submit -->
        <button
          type="submit"
          :disabled="loading"
          :aria-busy="loading"
          :class="['relative flex h-12 w-full items-center justify-center gap-2 rounded-xl bg-primary text-sm font-semibold text-on-primary transition-colors hover:opacity-90 disabled:cursor-wait disabled:opacity-60', focusRing]"
        >
          <UiAppIcon
            v-if="loading"
            :icon="Loading03Icon"
            :size="18"
            class="absolute animate-spin motion-reduce:animate-none"
          />
          <span :class="{ 'opacity-0': loading }">
            Log in
          </span>
        </button>
      </form>

      <!-- Footer links -->
      <div class="mt-7 space-y-1 text-center text-sm">
        <p class="text-muted">
          No account?
          <NuxtLink to="/auth/signup" :class="['rounded-md px-1 py-2 font-medium text-primary hover:underline', focusRing]">
            Sign up
          </NuxtLink>
        </p>

        <NuxtLink
          to="/"
          :class="['inline-flex min-h-11 items-center gap-1.5 rounded-lg px-2 font-medium text-muted transition-colors hover:text-primary', focusRing]"
        >
          <UiAppIcon :icon="ArrowLeft01Icon" :size="16" />
          Back to home
        </NuxtLink>
      </div>
    </div>
  </div>
</template>