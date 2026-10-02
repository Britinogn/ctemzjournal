import type { SupabaseClient } from '@supabase/supabase-js'
import type { Profile } from '~/types'

/** Keeps signed-in users off /login and /signup (role-aware landing). */
export default defineNuxtRouteMiddleware(async () => {
  if (import.meta.server)
    return
  const { $supabase } = useNuxtApp() as unknown as { $supabase: SupabaseClient }
  const { data: { session } } = await $supabase.auth.getSession()
  if (!session)
    return
  try {
    const me = await useApi().get<Profile>('/me')
    return navigateTo(me.Role === 'admin' ? '/admin' : '/dashboard')
  }
  catch {
    return
  }
})
