<script setup lang="ts">
import { toast } from 'vue-sonner'
import type { SupabaseClient } from '@supabase/supabase-js'

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
  <div class="rounded-2xl border border-border bg-surface p-6 shadow-sm sm:p-8">
    <h1 class="text-2xl font-bold tracking-tight">
      Forgot password
    </h1>
    <p class="mt-1 text-sm text-muted">
      We'll email you a reset link
    </p>

    <form v-if="!sent" class="mt-6 space-y-4" @submit.prevent="onSubmit">
      <div>
        <label for="forgot-email" class="mb-1.5 block text-sm font-medium">Email</label>
        <input
          id="forgot-email"
          v-model="email"
          type="email"
          required
          autocomplete="email"
          placeholder="you@example.com"
          class="w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary"
        >
      </div>
      <button
        type="submit"
        :disabled="loading"
        class="w-full rounded-xl bg-primary py-2.5 text-sm font-semibold text-on-primary transition hover:opacity-90 disabled:opacity-50"
      >
        {{ loading ? 'Sending…' : 'Send reset link' }}
      </button>
    </form>
    <p v-else class="mt-6 rounded-xl bg-profit/10 px-4 py-2.5 text-sm text-profit-text">
      Reset link sent. Check your inbox, then set a new password.
    </p>

    <p class="mt-6 text-center text-sm text-muted">
      <NuxtLink to="/auth/login" class="font-medium text-primary hover:underline">
        Back to login
      </NuxtLink>
    </p>
  </div>
</template>
