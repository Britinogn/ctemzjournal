<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { toast } from 'vue-sonner'
import type { SupabaseClient } from '@supabase/supabase-js'
import { Moon02Icon, Sun03Icon } from '~/utils/icons'
import { dashboardKey, meKey, type Profile, type ProfileUpdate } from '~/types'

definePageMeta({ middleware: 'auth', layout: 'dashboard' })

const api = useApi()
const queryClient = useQueryClient()
const { $supabase } = useNuxtApp() as unknown as { $supabase: SupabaseClient }
const { logout } = useAuth()
const colorMode = useColorMode()

const { data: me, isPending } = useQuery({
  queryKey: meKey(),
  queryFn: () => api.get<Profile>('/me'),
  staleTime: 5 * 60_000,
})

const displayName = ref('')
const timezone = ref('Africa/Lagos')
const saving = ref(false)
const email = ref('')
const signingOut = ref(false)

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

const THEMES = [
  { value: 'light', label: 'Light', icon: Sun03Icon },
  { value: 'dark', label: 'Dark', icon: Moon02Icon },
  { value: 'system', label: 'System', icon: null },
] as const

watchEffect(() => {
  if (me.value) {
    displayName.value = me.value.DisplayName ?? ''
    timezone.value = me.value.Timezone || 'Africa/Lagos'
  }
})

onMounted(async () => {
  const { data: { session } } = await $supabase.auth.getSession()
  email.value = session?.user?.email ?? ''
})

const memberSince = computed(() => {
  if (!me.value?.CreatedAt)
    return '—'
  return new Date(me.value.CreatedAt).toLocaleDateString('en-GB', { month: 'short', year: 'numeric', timeZone: 'UTC' })
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

async function onLogout(): Promise<void> {
  if (signingOut.value)
    return
  signingOut.value = true
  try {
    await logout()
  }
  finally {
    signingOut.value = false
  }
}

const cardCls = 'rounded-2xl border border-border bg-surface p-4 md:p-6'
const inputCls = 'w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary'
</script>

<template>
  <div>
    <h1 class="mb-4 text-xl font-bold tracking-tight md:text-2xl">
      Settings
    </h1>

    <div v-if="isPending" class="grid gap-4 md:grid-cols-2">
      <div class="h-72 animate-pulse rounded-2xl bg-surface" />
      <div class="space-y-4">
        <div class="h-44 animate-pulse rounded-2xl bg-surface" />
        <div class="h-24 animate-pulse rounded-2xl bg-surface" />
      </div>
    </div>

    <div v-else class="grid items-start gap-4 md:grid-cols-2">
      <!-- Profile (everything PATCH /me supports, minus avatar upload) -->
      <section :class="cardCls" aria-label="Profile">
        <h2 class="text-sm font-semibold">
          Profile
        </h2>
        <p class="mt-0.5 text-xs text-muted">
          Display name shows on your public journals. Calendar days group in this timezone.
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
              :class="inputCls"
            >
          </div>
          <div>
            <label for="settings-tz" class="mb-1.5 block text-sm font-medium">Timezone</label>
            <select id="settings-tz" v-model="timezone" :class="inputCls">
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

      <div class="space-y-4">
        <!-- Appearance -->
        <section :class="cardCls" aria-label="Appearance">
          <h2 class="text-sm font-semibold">
            Appearance
          </h2>
          <p class="mt-0.5 text-xs text-muted">
            Follows your system until you pick one.
          </p>
          <div class="mt-3 grid grid-cols-3 gap-2" role="group" aria-label="Theme">
            <button
              v-for="t in THEMES"
              :key="t.value"
              type="button"
              :aria-pressed="colorMode.preference === t.value"
              :class="[
                'flex items-center justify-center gap-1.5 rounded-xl border py-2.5 text-sm font-medium transition',
                colorMode.preference === t.value
                  ? 'border-primary bg-primary/10 text-primary'
                  : 'border-border text-muted hover:text-text',
              ]"
              @click="colorMode.preference = t.value"
            >
              <UiAppIcon v-if="t.icon" :icon="t.icon" :size="16" />
              {{ t.label }}
            </button>
          </div>
        </section>

        <!-- Account -->
        <section :class="cardCls" aria-label="Account">
          <h2 class="text-sm font-semibold">
            Account
          </h2>
          <dl class="mt-3 space-y-2 text-sm">
            <div class="flex items-center justify-between gap-2">
              <dt class="text-muted">Email</dt>
              <dd class="truncate font-medium">{{ email || '—' }}</dd>
            </div>
            <div class="flex items-center justify-between gap-2">
              <dt class="text-muted">Role</dt>
              <dd class="rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-semibold capitalize text-primary">
                {{ me?.Role ?? 'user' }}
              </dd>
            </div>
            <div class="flex items-center justify-between gap-2">
              <dt class="text-muted">Member since</dt>
              <dd class="font-medium">{{ memberSince }}</dd>
            </div>
          </dl>
          <button
            type="button"
            :disabled="signingOut"
            class="mt-4 w-full rounded-xl border border-loss/40 py-2.5 text-sm font-semibold text-loss transition hover:bg-loss/10 disabled:opacity-50"
            @click="onLogout"
          >
            {{ signingOut ? 'Signing out…' : 'Sign out' }}
          </button>
        </section>
      </div>
    </div>
  </div>
</template>
