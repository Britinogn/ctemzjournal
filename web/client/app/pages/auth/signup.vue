<!-- pages/auth/signup.vue -->
<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { siteSettingsKey, type PublicSettings } from '~/types'
import {
  UserIcon,
  Mail01Icon,
  LockKeyIcon,
  EyeIcon,
  EyeOffIcon,
  Loading03Icon,
  ArrowLeft01Icon,
} from '~/utils/icons'

definePageMeta({ middleware: 'guest' })

const api = useApi()

const { data: settings } = useQuery({
  queryKey: siteSettingsKey(),
  queryFn: () => api.get<PublicSettings>('/public/site-settings'),
  staleTime: 30 * 60_000,
})

const signupsOpen = computed(() => settings.value?.allow_signups ?? true)

const displayName = ref('')
const email = ref('')
const password = ref('')
const showPassword = ref(false)
const agreed = ref(false)
const loading = ref(false)
const error = ref<string | null>(null)
const { signup } = useAuth()

async function onSubmit(): Promise<void> {
  if (loading.value)
    return
  if (!agreed.value) {
    error.value = 'Please accept the Terms of Service and Privacy Policy to continue.'
    return
  }
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
  <div class="w-full max-w-md">
    <div class="rounded-2xl border border-border bg-surface p-7 shadow-sm sm:p-8">
      <!-- Header -->
      <div class="mb-7">
        <div class="mb-4 flex h-11 w-11 items-center justify-center rounded-xl bg-primary/10 text-primary">
          <UiAppIcon :icon="UserIcon" :size="18" />
        </div>
        <h1 class="text-2xl font-bold tracking-tight">
          Create your journal
        </h1>
        <p class="mt-1.5 text-sm text-muted">
          Start tracking every trade and find your edge
        </p>
      </div>

      <!-- Form -->
      <form v-if="signupsOpen" class="space-y-5" @submit.prevent="onSubmit">
        <!-- Display name -->
        <div>
          <label for="signup-name" class="mb-1.5 block text-sm font-medium">
            Display name
          </label>
          <div class="relative">
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5 text-muted">
              <UiAppIcon :icon="UserIcon" :size="18" />
            </div>
            <input
              id="signup-name"
              v-model="displayName"
              type="text"
              autocomplete="nickname"
              placeholder="Ada Obi"
              class="w-full rounded-xl border border-border bg-bg py-3 pl-11 pr-4 text-sm outline-none transition
                placeholder:text-muted
                focus:border-primary focus:ring-2 focus:ring-primary/20"
            >
          </div>
        </div>

        <!-- Email -->
        <div>
          <label for="signup-email" class="mb-1.5 block text-sm font-medium">
            Email
          </label>
          <div class="relative">
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5 text-muted">
              <UiAppIcon :icon="Mail01Icon" :size="18" />
            </div>
            <input
              id="signup-email"
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
          <label for="signup-password" class="mb-1.5 block text-sm font-medium">
            Password
          </label>
          <div class="relative">
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5 text-muted">
              <UiAppIcon :icon="LockKeyIcon" :size="18" />
            </div>
            <input
              id="signup-password"
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              required
              minlength="6"
              autocomplete="new-password"
              placeholder="Minimum 6 characters"
              class="w-full rounded-xl border border-border bg-bg py-3 pl-11 pr-12 text-sm outline-none transition
                placeholder:text-muted
                focus:border-primary focus:ring-2 focus:ring-primary/20"
            >
            <button
              type="button"
              :aria-label="showPassword ? 'Hide password' : 'Show password'"
              class="absolute inset-y-0 right-0 flex items-center pr-3.5 text-muted transition hover:text-text"
              @click="showPassword = !showPassword"
            >
              <UiAppIcon
                :icon="showPassword ? EyeOffIcon : EyeIcon"
                :size="18"
              />
            </button>
          </div>
        </div>

        <!-- Terms & privacy -->
        <label class="flex cursor-pointer items-start gap-3 text-sm leading-relaxed text-muted">
          <input
            v-model="agreed"
            type="checkbox"
            required
            class="mt-1 h-4 w-4 shrink-0 rounded border-border accent-primary"
          >
          <span>
            I am 18 or older and agree to the
            <NuxtLink
              to="/terms"
              target="_blank"
              class="font-medium text-primary hover:underline"
            >Terms of Service</NuxtLink>
            and
            <NuxtLink
              to="/privacy"
              target="_blank"
              class="font-medium text-primary hover:underline"
            >Privacy Policy</NuxtLink>.
          </span>
        </label>

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
          :disabled="loading || !agreed"
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
            Create account
          </span>
        </button>
      </form>

      <!-- Signups closed -->
      <div
        v-else
        class="mt-2 rounded-xl border border-warning/20 bg-warning/10 px-5 py-4 text-center"
      >
        <p class="text-sm font-medium text-warning-text">
          Signups are currently closed
        </p>
        <p class="mt-1 text-sm text-muted">
          Ask an admin for an account, or check back later.
        </p>
      </div>

      <!-- Footer links -->
      <div class="mt-7 space-y-3 text-center text-sm">
        <p class="text-muted">
          Already have an account?
          <NuxtLink to="/auth/login" class="font-medium text-primary hover:underline">
            Log in
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