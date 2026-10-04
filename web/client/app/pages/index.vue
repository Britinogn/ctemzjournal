<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import type { SupabaseClient } from '@supabase/supabase-js'
import { meKey, siteSettingsKey, type Profile, type PublicSettings } from '~/types'

definePageMeta({ layout: 'public' })

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

useSeoMeta({
  title: 'Ctemz Journal — Know why you win, and why you lose',
  description: 'Forex trading journal: log trades, track setups and rules, share results without showing money.',
})
</script>

<template>
  <div class="flex-1">
    <HomeHeroSection
      :site-name="settings?.site_name"
      :tagline="settings?.tagline"
      :dashboard-to="dashboardTo"
      :allow-signups="settings?.allow_signups ?? true"
    />
    <HomeLivePricesStrip />
    <HomeRecentJournals />
    <HomeHowItWorks />
    <HomeCtaBanner :dashboard-to="dashboardTo" :allow-signups="settings?.allow_signups ?? true" />
  </div>
</template>
