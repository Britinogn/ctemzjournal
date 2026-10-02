<script setup lang="ts">
import { toast } from 'vue-sonner'
import type { SupabaseClient } from '@supabase/supabase-js'

const password = ref('')
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
  <div class="rounded-2xl border border-border bg-surface p-6 shadow-sm sm:p-8">
    <h1 class="text-2xl font-bold tracking-tight">
      Set new password
    </h1>
    <p class="mt-1 text-sm text-muted">
      Choose something strong
    </p>

    <form class="mt-6 space-y-4" @submit.prevent="onSubmit">
      <div>
        <label for="reset-password" class="mb-1.5 block text-sm font-medium">New password</label>
        <input
          id="reset-password"
          v-model="password"
          type="password"
          required
          minlength="6"
          autocomplete="new-password"
          placeholder="Minimum 6 characters"
          class="w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary"
        >
      </div>
      <button
        type="submit"
        :disabled="loading"
        class="w-full rounded-xl bg-primary py-2.5 text-sm font-semibold text-on-primary transition hover:opacity-90 disabled:opacity-50"
      >
        {{ loading ? 'Saving…' : 'Save new password' }}
      </button>
    </form>
  </div>
</template>
