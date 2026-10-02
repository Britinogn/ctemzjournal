import type { SupabaseClient } from '@supabase/supabase-js'

/**
 * User-area guard. Client-side: Supabase session lives in localStorage,
 * so the server pass lets the request through and the client enforces.
 * /me failures already redirect via the API client (401 → /login,
 * 403 → /suspended).
 */
export default defineNuxtRouteMiddleware(async () => {
  if (import.meta.server)
    return
  const { $supabase } = useNuxtApp() as unknown as { $supabase: SupabaseClient }
  const { data: { session } } = await $supabase.auth.getSession()
  if (!session)
    return navigateTo('/auth/login')
})
