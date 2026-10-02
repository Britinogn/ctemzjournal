import type { SupabaseClient } from '@supabase/supabase-js'
import type { Profile } from '~/types'

/**
 * Admin-area guard: session + DB role must be admin, otherwise the user
 * is sent to /dashboard (experience only — the API enforces roles).
 */
export default defineNuxtRouteMiddleware(async () => {
  if (import.meta.server)
    return
  const { $supabase } = useNuxtApp() as unknown as { $supabase: SupabaseClient }
  const { data: { session } } = await $supabase.auth.getSession()
  if (!session)
    return navigateTo('/auth/login')
  try {
    const me = await useApi().get<Profile>('/me')
    if (me.Role !== 'admin')
      return navigateTo('/dashboard')
  }
  catch {
    return navigateTo('/auth/login')
  }
})
