<script setup lang="ts">
import { toast } from 'vue-sonner'
import type { SupabaseClient } from '@supabase/supabase-js'
import {
  LockKeyIcon,
  EyeIcon,
  EyeOffIcon,
  Loading03Icon,
  ArrowLeft01Icon,
} from '~/utils/icons'

const password = ref('')
const showPassword = ref(false)
const loading = ref(false)

async function onSubmit(): Promise<void> {
  if (loading.value)
    return
  loading.value = true
  try {
    const { $supabase } = useNuxtApp() as unknown as { $supabase: SupabaseClient }
    const { error } = await $supabase.auth.updateUser({ password: password.value })
    if (error)
      throw new Error(error.message)
    toast.success('Password updated — log in with the new one')
    await navigateTo('/auth/login')
  }
  catch (err) {
    toast.error(err instanceof Error ? err.message : 'Update failed')
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
          <UiAppIcon :icon="LockKeyIcon" :size="18" />
        </div>
        <h1 class="text-2xl font-bold tracking-tight">
          Set new password
        </h1>
        <p class="mt-1.5 text-sm text-muted">
          Choose something strong and memorable
        </p>
      </div>

      <!-- Form -->
      <form class="space-y-5" @submit.prevent="onSubmit">
        <!-- New password -->
        <div>
          <label for="reset-password" class="mb-1.5 block text-sm font-medium">
            New password
          </label>
          <div class="relative">
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5 text-muted">
              <UiAppIcon :icon="LockKeyIcon" :size="18" />
            </div>
            <input
              id="reset-password"
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

        <!-- Submit -->
        <button
          type="submit"
          :disabled="loading || password.length < 6"
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
            Save new password
          </span>
        </button>
      </form>

      <!-- Footer -->
      <div class="mt-7 text-center">
        <NuxtLink
          to="/auth/login"
          class="inline-flex items-center gap-1.5 text-sm font-medium text-muted transition hover:text-primary"
        >
          <UiAppIcon :icon="ArrowLeft01Icon" :size="16" />
          Back to home
        </NuxtLink>
      </div>
    </div>
  </div>
</template>