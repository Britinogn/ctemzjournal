<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import type { SupabaseClient } from '@supabase/supabase-js'
import { meKey, siteSettingsKey, type Profile, type PublicSettings } from '~/types'

const api = useApi()
const { $supabase } = useNuxtApp() as unknown as { $supabase: SupabaseClient }

const { data: settings } = useQuery({
  queryKey: siteSettingsKey(),
  queryFn: () => api.get<PublicSettings>('/public/site-settings'),
  staleTime: 30 * 60_000,
})

const hasSession = ref(false)
onMounted(async () => {
  const { data: { session } } = await $supabase.auth.getSession()
  hasSession.value = !!session
})

const { data: me } = useQuery({
  queryKey: meKey(),
  queryFn: () => api.get<Profile>('/me'),
  enabled: hasSession,
  retry: false,
  staleTime: 5 * 60_000,
})

const dashboardTo = computed<string | null>(() => {
  if (!me.value)
    return null
  return me.value.Role === 'admin' ? '/admin' : '/dashboard'
})

/** Maintenance locks out visitors; admins still browse and always sign in. */
const underMaintenance = computed(
  () => settings.value?.maintenance_mode === true && me.value?.Role !== 'admin',
)
</script>

<template>
  <div class="flex min-h-screen flex-col bg-bg text-text">
    <LayoutHomelayoutPublicNavbar
      :site-name="settings?.site_name"
      :logo-url="settings?.logo_url"
      :allow-signups="settings?.allow_signups ?? true"
      :dashboard-to="dashboardTo"
    />

    
    <div
  v-if="underMaintenance"
  class="mx-auto flex w-full max-w-md flex-1 flex-col items-center justify-center px-4 py-16 sm:py-24"
>
  <div class="w-full rounded-2xl border border-border bg-surface p-7 text-center shadow-sm sm:p-8">

    <!-- Header -->
    <div class="mb-6">
      <div class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 text-primary">
        <UiAppIcon :icon="Wrench01Icon" :size="22" />
      </div>

      <span class="mb-3 inline-flex items-center gap-1.5 rounded-full bg-primary/10 px-3 py-1 text-xs font-medium text-primary">
        <span class="relative flex h-2 w-2">
          <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-primary opacity-60" />
          <span class="relative inline-flex h-2 w-2 rounded-full bg-primary" />
        </span>
        Scheduled maintenance
      </span>

      <h1 class="text-2xl font-bold tracking-tight">
        We'll be right back
      </h1>
      <p class="mt-2 text-sm text-muted">
        {{ settings?.site_name ?? 'The site' }} is temporarily down while we make some improvements
        to keep things fast, secure, and reliable.
      </p>
    </div>

    <!-- Reassurance -->
    <div class="mb-6 rounded-xl border border-border bg-bg/50 p-4 text-left">
      <ul class="space-y-2.5 text-sm text-muted">
        <li class="flex items-start gap-2.5">
          <span class="mt-0.5 text-primary">✓</span>
          <span>Your journal entries and account data are safe and untouched.</span>
        </li>
        <li class="flex items-start gap-2.5">
          <span class="mt-0.5 text-primary">✓</span>
          <span>No action is needed on your part.</span>
        </li>
        <li class="flex items-start gap-2.5">
          <span class="mt-0.5 text-primary">✓</span>
          <span>This usually takes only a short while. Please check back soon.</span>
        </li>
      </ul>
    </div>

    <!-- Action -->
    <NuxtLink
      to="/auth/login"
      class="inline-flex w-full items-center justify-center rounded-xl bg-primary py-3 text-sm font-semibold text-on-primary transition hover:opacity-90"
    >
      Admin sign in
    </NuxtLink>

    <!-- Contact -->
    <p v-if="settings?.contact_email" class="mt-5 text-xs text-muted">
      Need help in the meantime?
      <a
        :href="`mailto:${settings.contact_email}`"
        class="font-medium text-primary hover:underline"
      >
        {{ settings.contact_email }}
      </a>
    </p>
  </div>
</div>
    
    <slot v-else />
    <LayoutHomelayoutPublicFooter
      :site-name="settings?.site_name"
      :tagline="settings?.tagline"
      :contact-email="settings?.contact_email"
      :footer-text="settings?.footer_text"
      :risk-disclaimer="settings?.risk_disclaimer"
      :social-links="settings?.social_links"
    />
  </div>
</template>
