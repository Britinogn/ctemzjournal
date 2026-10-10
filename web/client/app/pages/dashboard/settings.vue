<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { toast } from 'vue-sonner'
import imageCompression from 'browser-image-compression'
import type { SupabaseClient } from '@supabase/supabase-js'
import { Logout02Icon, Moon02Icon, Sun03Icon, Tag01Icon, Target01Icon, Wallet01Icon } from '~/utils/icons'
import { dashboardKey, meKey, type Profile, type ProfileUpdate } from '~/types'
import { DEFAULT_TIMEZONE, timezoneOptions } from '~/utils/timezones'
import { AVATAR_BUCKET, avatarStoragePath, resolveAvatarUrl } from '~/utils/avatar'

definePageMeta({ middleware: 'auth', layout: 'dashboard' })

const api = useApi()
const queryClient = useQueryClient()
const config = useRuntimeConfig()
const { $supabase } = useNuxtApp() as unknown as { $supabase: SupabaseClient }
const { logout } = useAuth()
const colorMode = useColorMode()

const { data: me, isPending } = useQuery({
  queryKey: meKey(),
  queryFn: () => api.get<Profile>('/me'),
  staleTime: 5 * 60_000,
})

const displayName = ref('')
const timezone = ref(DEFAULT_TIMEZONE)
const avatarPath = ref('')
const avatarBroken = ref(false)
const uploadingAvatar = ref(false)
const saving = ref(false)
const email = ref('')
const signingOut = ref(false)

const THEMES = [
  { value: 'light', label: 'Light', icon: Sun03Icon },
  { value: 'dark', label: 'Dark', icon: Moon02Icon },
  { value: 'system', label: 'System', icon: null },
] as const

// On phones this page is the "More" tab, and the sidebar is not there, so these three lists need a way in.
const LISTS = [
  { label: 'Accounts', to: '/dashboard/accounts', icon: Wallet01Icon },
  { label: 'Setups', to: '/dashboard/setups', icon: Target01Icon },
  { label: 'Tags', to: '/dashboard/tags', icon: Tag01Icon },
]

watchEffect(() => {
  if (me.value) {
    displayName.value = me.value.DisplayName ?? ''
    timezone.value = me.value.Timezone || DEFAULT_TIMEZONE
    avatarPath.value = me.value.AvatarPath ?? ''
    avatarBroken.value = false
  }
})

const tzOptions = computed(() => timezoneOptions(timezone.value))
const avatarUrl = computed(() => resolveAvatarUrl(String(config.public.supabaseUrl || ''), avatarPath.value, me.value?.UpdatedAt))
const showAvatar = computed(() => avatarUrl.value !== '' && !avatarBroken.value)

/** Upload a profile picture: compress → avatars/{userId}/avatar.webp → PATCH /me. */
async function onAvatarFile(event: Event): Promise<void> {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file || uploadingAvatar.value || !me.value)
    return
  if (!file.type.startsWith('image/')) {
    toast.error('Choose an image file')
    return
  }
  if (file.size > 5 * 1024 * 1024) {
    toast.error('Image too large — max 5 MB')
    return
  }
  uploadingAvatar.value = true
  try {
    const compressed = await imageCompression(file, {
      maxSizeMB: 0.5,
      maxWidthOrHeight: 512,
      fileType: 'image/webp',
    })
    const webp = new File([compressed], 'avatar.webp', { type: 'image/webp' })
    const path = avatarStoragePath(me.value.ID)
    const { error } = await $supabase.storage
      .from(AVATAR_BUCKET)
      .upload(path, webp, { upsert: true, contentType: 'image/webp' })
    if (error)
      throw new Error(error.message)
    const updated = await api.patch<Profile>('/me', { avatar_path: path })
    queryClient.setQueryData(meKey(), updated)
    avatarPath.value = updated.AvatarPath ?? path
    avatarBroken.value = false
    toast.success('Profile picture updated')
  }
  catch (err) {
    toast.error(err instanceof Error ? err.message : 'Upload failed — try again')
  }
  finally {
    uploadingAvatar.value = false
  }
}

async function onRemoveAvatar(): Promise<void> {
  if (uploadingAvatar.value || !me.value || !avatarPath.value)
    return
  uploadingAvatar.value = true
  try {
    // Best-effort storage cleanup — legacy http(s) URLs have no stored file.
    if (!/^https?:\/\//i.test(avatarPath.value)) {
      await $supabase.storage.from(AVATAR_BUCKET).remove([avatarPath.value])
    }
    const updated = await api.patch<Profile>('/me', { avatar_path: '' })
    queryClient.setQueryData(meKey(), updated)
    avatarPath.value = ''
    toast.success('Profile picture removed')
  }
  catch {
    toast.error('Could not remove picture')
  }
  finally {
    uploadingAvatar.value = false
  }
}

onMounted(async () => {
  const { data: { session } } = await $supabase.auth.getSession()
  email.value = session?.user?.email ?? ''
})

const memberSince = computed(() => {
  if (!me.value?.CreatedAt)
    return '—'
  const tz = me.value?.Timezone || timezone.value || DEFAULT_TIMEZONE
  try {
    return new Date(me.value.CreatedAt).toLocaleDateString('en-GB', { month: 'short', year: 'numeric', timeZone: tz })
  }
  catch {
    return new Date(me.value.CreatedAt).toLocaleDateString('en-GB', { month: 'short', year: 'numeric', timeZone: 'UTC' })
  }
})

// The save button wakes up only when name/timezone changed (picture saves instantly).
const dirty = computed(() =>
  displayName.value.trim() !== (me.value?.DisplayName ?? '')
  || timezone.value !== (me.value?.Timezone || DEFAULT_TIMEZONE),
)

async function onSave(): Promise<void> {
  if (saving.value)
    return
  saving.value = true
  try {
    const body: ProfileUpdate = {
      display_name: displayName.value.trim(),
      timezone: timezone.value,
    }
    const updated = await api.patch<Profile>('/me', { ...body })
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

/* ---------- Shared look, defined once ---------- */
const cardCls = 'rounded-2xl border border-border bg-surface p-4 md:p-6'
const inputCls
  = 'h-11 w-full rounded-xl border border-border bg-bg px-4 text-sm outline-none transition-colors placeholder:text-muted focus:border-primary focus-visible:ring-2 focus-visible:ring-primary/30'
const skel = 'animate-pulse rounded-2xl border border-border bg-surface'
</script>

<template>
  <div>
    <!-- The top bar shows "Settings" from tablet up; this heading is for phones and screen readers -->
    <h1 class="mb-4 text-xl font-bold tracking-tight md:sr-only">
      Settings
    </h1>

    <!-- Phones only. Tablet and desktop have these in the sidebar. -->
    <nav class="mb-4 md:hidden" aria-label="Your lists">
      <ul :class="[cardCls, 'divide-y divide-border p-0!']">
        <li v-for="l in LISTS" :key="l.to">
          <NuxtLink
            :to="l.to"
            class="flex min-h-14 items-center gap-3 px-4 text-sm font-semibold transition-colors hover:bg-bg focus-visible:outline-2-2 focus-visible:-outline-offset-2 focus-visible:outline-2-primary"
          >
            <span class="inline-flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
              <UiAppIcon :icon="l.icon" :size="18" />
            </span>
            <span class="flex-1">{{ l.label }}</span>
            <svg class="h-4 w-4 text-muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
              <path d="M9 5l7 7-7 7" />
            </svg>
          </NuxtLink>
        </li>
      </ul>
    </nav>

    <!-- Loading: the same shape as the real page -->
    <div v-if="isPending" class="grid gap-4 md:grid-cols-2" aria-hidden="true">
      <div :class="[skel, 'h-72']" />
      <div class="space-y-4">
        <div :class="[skel, 'h-44']" />
        <div :class="[skel, 'h-48']" />
      </div>
    </div>

    <div v-else class="grid items-start gap-4 md:grid-cols-2">
      <!-- Profile (display name + timezone + profile picture) -->
      <section :class="cardCls" aria-label="Profile">
        <h2 class="text-sm font-semibold">
          Profile
        </h2>
        <p class="mt-1 text-sm text-muted">
          Display name shows on your public journals. Calendar days group in this timezone.
        </p>
        <form class="mt-5 space-y-4" @submit.prevent="onSave">
          <div class="flex items-center gap-3">
            <img
              v-if="showAvatar"
              :src="avatarUrl"
              alt="Your profile picture"
              class="h-12 w-12 rounded-full object-cover"
              @error="avatarBroken = true"
            >
            <span
              v-else
              class="inline-flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-primary/15 font-semibold text-primary"
              aria-hidden="true"
            >
              {{ (displayName.trim()?.[0] ?? 'T').toUpperCase() }}
            </span>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-semibold">
                {{ displayName.trim() || 'Trader' }}
              </p>
              <div class="mt-1 flex gap-2">
                <label
                  class="inline-flex h-9 cursor-pointer items-center rounded-xl border border-border px-3 text-xs font-semibold transition hover:border-primary hover:text-primary"
                >
                  <input type="file" accept="image/png,image/jpeg,image/webp" class="sr-only" :disabled="uploadingAvatar" @change="onAvatarFile">
                  {{ uploadingAvatar ? 'Uploading…' : showAvatar ? 'Change picture' : 'Upload picture' }}
                </label>
                <button
                  v-if="showAvatar"
                  type="button"
                  :disabled="uploadingAvatar"
                  class="inline-flex h-9 items-center rounded-xl border border-border px-3 text-xs font-semibold text-muted transition hover:text-text disabled:opacity-50"
                  @click="onRemoveAvatar"
                >
                  Remove
                </button>
              </div>
            </div>
          </div>
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
            <div class="relative">
              <select id="settings-tz" v-model="timezone" :class="[inputCls, 'appearance-none pr-10']">
                <option v-for="tz in tzOptions" :key="tz" :value="tz">
                  {{ tz }}
                </option>
              </select>
              <svg class="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                <path d="M6 9l6 6 6-6" />
              </svg>
            </div>
          </div>
          <!-- Full width on phones, natural width from tablet up. Off until something changes. -->
          <button
            type="submit"
            :disabled="saving || !dirty"
            class="inline-flex h-12 w-full items-center justify-center rounded-xl bg-primary px-6 text-sm font-semibold text-on-primary transition-colors hover:opacity-90 focus-visible:outline-2-2 focus-visible:outline-2-offset-2 focus-visible:outline-2-primary disabled:cursor-not-allowed disabled:opacity-50 sm:w-auto"
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
          <p class="mt-1 text-sm text-muted">
            Follows your system until you pick one.
          </p>
          <div class="mt-4 grid grid-cols-3 gap-2" role="group" aria-label="Theme">
            <button
              v-for="t in THEMES"
              :key="t.value"
              type="button"
              :aria-pressed="colorMode.preference === t.value"
              :class="[
                'flex h-11 items-center justify-center gap-1.5 rounded-xl border text-sm font-medium transition-colors focus-visible:outline-2-2 focus-visible:outline-2-offset-2 focus-visible:outline-2-primary',
                colorMode.preference === t.value
                  ? 'border-primary bg-primary/10 font-semibold text-primary'
                  : 'border-border text-muted hover:text-text',
              ]"
              @click="colorMode.preference = t.value"
            >
              <UiAppIcon v-if="t.icon" :icon="t.icon" :size="16" />
              <!-- System has no icon in the set, so it gets a small screen -->
              <svg v-else class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                <rect x="3" y="4" width="18" height="12" rx="2" /><path d="M8 20h8M12 16v4" />
              </svg>
              {{ t.label }}
            </button>
          </div>
        </section>

        <!-- Account -->
        <section :class="cardCls" aria-label="Account">
          <h2 class="text-sm font-semibold">
            Account
          </h2>
          <dl class="mt-4 divide-y divide-border text-sm">
            <div class="flex min-h-11 items-center justify-between gap-3">
              <dt class="shrink-0 text-muted">
                Email
              </dt>
              <dd class="min-w-0 truncate font-medium">
                {{ email || '—' }}
              </dd>
            </div>
            <div class="flex min-h-11 items-center justify-between gap-3">
              <dt class="text-muted">
                Role
              </dt>
              <dd
                class="rounded-full px-2.5 py-0.5 text-xs font-semibold capitalize"
                :class="me?.Role === 'admin' ? 'bg-warning/10 text-warning-text' : 'bg-primary/10 text-primary'"
              >
                {{ me?.Role ?? 'user' }}
              </dd>
            </div>
            <div class="flex min-h-11 items-center justify-between gap-3">
              <dt class="text-muted">
                Member since
              </dt>
              <dd class="tnum font-medium">
                {{ memberSince }}
              </dd>
            </div>
          </dl>
          <!-- Neutral, not red: your plan keeps red for trade results -->
          <button
            type="button"
            :disabled="signingOut"
            :aria-busy="signingOut"
            class="mt-4 inline-flex h-12 w-full items-center justify-center gap-2 rounded-xl border border-border text-sm font-semibold text-muted transition-colors hover:border-muted hover:text-text focus-visible:outline-2-2 focus-visible:outline-2-offset-2 focus-visible:outline-2-primary disabled:cursor-wait disabled:opacity-60"
            @click="onLogout"
          >
            <UiAppIcon :icon="Logout02Icon" :size="18" />
            {{ signingOut ? 'Signing out…' : 'Sign out' }}
          </button>
        </section>
      </div>
    </div>
  </div>
</template>