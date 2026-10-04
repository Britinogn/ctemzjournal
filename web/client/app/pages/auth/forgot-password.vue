<script setup lang="ts">
import { toast } from 'vue-sonner'
import type { SupabaseClient } from '@supabase/supabase-js'
import {
  LockKeyIcon,
  ArrowLeft01Icon,
  Loading03Icon,
  Tick01Icon,
} from '~/utils/icons'

// or however you re-export them

const email = ref('')
const loading = ref(false)
const sent = ref(false)

async function onSubmit(): Promise<void> {
  if (loading.value || sent.value)
    return
  loading.value = true
  try {
    const { $supabase } = useNuxtApp() as unknown as { $supabase: SupabaseClient }
    const { error } = await $supabase.auth.resetPasswordForEmail(email.value.trim(), {
      redirectTo: `${window.location.origin}/auth/reset-password`,
    })
    if (error)
      throw new Error(error.message)
    sent.value = true
    toast.success('Reset link sent — check your inbox')
  }
  catch (err) {
    toast.error(err instanceof Error ? err.message : 'Request failed')
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
          <UiAppIcon :icon="LockKeyIcon" :size="20" />
        </div>
        <h1 class="text-2xl font-bold tracking-tight">
          Forgot password?
        </h1>
        <p class="mt-1.5 text-sm text-muted">
          No worries — enter your email and we’ll send you a reset link.
        </p>
      </div>

      <!-- Form -->
      <form v-if="!sent" class="space-y-5" @submit.prevent="onSubmit">
        <div>
          <label for="forgot-email" class="mb-1.5 block text-sm font-medium">
            Email address
          </label>
          <input
            id="forgot-email"
            v-model="email"
            type="email"
            required
            autocomplete="email"
            placeholder="you@example.com"
            class="w-full rounded-xl border border-border bg-bg px-4 py-3 text-sm outline-none transition
              placeholder:text-muted
              focus:border-primary focus:ring-2 focus:ring-primary/20"
          >
        </div>

        <button
          type="submit"
          :disabled="loading || !email"
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
            Send reset link
          </span>
        </button>
      </form>

      <!-- Success state -->
      <div v-else class="rounded-xl border border-profit/20 bg-profit/10 px-5 py-4">
        <div class="flex items-start gap-3">
          <div class="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-profit text-white">
            <UiAppIcon :icon="Tick01Icon" :size="12" />
          </div>
          <div>
            <p class="text-sm font-medium text-profit-text">
              Reset link sent
            </p>
            <p class="mt-1 text-sm text-muted">
              Check your inbox and follow the link to set a new password.
            </p>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="mt-7 text-center">
        <NuxtLink
          to="/auth/login"
          class="inline-flex items-center gap-1.5 text-sm font-medium text-muted transition hover:text-primary"
        >
          <UiAppIcon :icon="ArrowLeft01Icon" :size="16" />
          Back to login
        </NuxtLink>
      </div>
    </div>
  </div>
</template>