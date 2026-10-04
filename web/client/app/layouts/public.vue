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
    <div v-if="underMaintenance" class="mx-auto flex w-full max-w-xl flex-1 flex-col items-center justify-center px-4 text-center">
      <h1 class="text-2xl font-bold tracking-tight">
        Down for maintenance
      </h1>
      <p class="mt-2 text-sm text-muted">
        We're tuning things up. Check back shortly — your journal is safe.
      </p>
      <NuxtLink
        to="/auth/login"
        class="mt-6 rounded-xl bg-primary px-6 py-2.5 text-sm font-semibold text-on-primary transition hover:opacity-90"
      >
        Sign in
      </NuxtLink>
    </div>
    <slot v-else />
    <LayoutHomelayoutPublicFooter
      :site-name="settings?.site_name"
      :tagline="settings?.tagline"
      :footer-text="settings?.footer_text"
      :risk-disclaimer="settings?.risk_disclaimer"
      :social-links="settings?.social_links"
    />
  </div>
</template>
