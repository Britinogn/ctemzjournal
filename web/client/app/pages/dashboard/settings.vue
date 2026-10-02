<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { toast } from 'vue-sonner'
import { dashboardKey, meKey, type Profile, type ProfileUpdate } from '~/types'

definePageMeta({ middleware: 'auth', layout: 'dashboard' })

const api = useApi()
const queryClient = useQueryClient()

const { data: me, isPending } = useQuery({
  queryKey: meKey(),
  queryFn: () => api.get<Profile>('/me'),
  staleTime: 5 * 60_000,
})

const displayName = ref('')
const timezone = ref('Africa/Lagos')
const saving = ref(false)

const TIMEZONES = [
  'Africa/Lagos',
  'UTC',
  'Europe/London',
  'Europe/Berlin',
  'America/New_York',
  'America/Chicago',
  'Asia/Dubai',
  'Asia/Tokyo',
  'Australia/Sydney',
]

watchEffect(() => {
  if (me.value) {
    displayName.value = me.value.DisplayName ?? ''
    timezone.value = me.value.Timezone || 'Africa/Lagos'
  }
})

async function onSave(): Promise<void> {
  if (saving.value)
    return
  saving.value = true
  try {
    const body: ProfileUpdate = {
      display_name: displayName.value.trim(),
      timezone: timezone.value,
    }
    const updated = await api.patch<Profile>('/me', body)
    queryClient.setQueryData(meKey(), updated)
    queryClient.invalidateQueries({ queryKey: dashboardKey() })
    toast.success('Profile saved')
  }
  catch {
    toast.error('Could not save profile')
  }
  finally {
    saving.value = false
  }
}
</script>

<template>
  <div>
    <h1 class="mb-4 text-xl font-bold tracking-tight md:text-2xl">
      Settings
    </h1>

    <div v-if="isPending" class="h-64 animate-pulse rounded-2xl bg-surface" />

    <section v-else class="max-w-xl rounded-2xl border border-border bg-surface p-4 md:p-6" aria-label="Profile">
      <h2 class="text-sm font-semibold">
        Profile
      </h2>
      <p class="mt-0.5 text-xs text-muted">
        Calendar days are grouped in this timezone.
      </p>

      <form class="mt-4 space-y-4" @submit.prevent="onSave">
        <div>
          <label for="settings-name" class="mb-1.5 block text-sm font-medium">Display name</label>
          <input
            id="settings-name"
            v-model="displayName"
            type="text"
            autocomplete="nickname"
            placeholder="Ada Obi"
            class="w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary"
          >
        </div>
        <div>
          <label for="settings-tz" class="mb-1.5 block text-sm font-medium">Timezone</label>
          <select
            id="settings-tz"
            v-model="timezone"
            class="w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition focus:border-primary"
          >
            <option v-for="tz in TIMEZONES" :key="tz" :value="tz">
              {{ tz }}
            </option>
          </select>
        </div>
        <button
          type="submit"
          :disabled="saving"
          class="rounded-xl bg-primary px-6 py-2.5 text-sm font-semibold text-on-primary transition hover:opacity-90 disabled:opacity-50"
        >
          {{ saving ? 'Saving…' : 'Save changes' }}
        </button>
      </form>
    </section>
  </div>
</template>
