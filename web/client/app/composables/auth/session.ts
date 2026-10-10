import type { QueryClient } from '@tanstack/vue-query'
import type { SupabaseClient } from '@supabase/supabase-js'
import type { Profile, SyncRequest } from '~/types'
import { detectBrowserTimezone } from '~/utils/timezones'

interface SupabaseNuxtApp {
  $supabase: SupabaseClient;
  $queryClient: QueryClient;
}

/**
 * Email auth (Supabase) + backend handshake.
 * login/signup → POST /auth/sync → GET /me → role redirect.
 */
export function useAuth() {
  const api = useApi()

  function client(): SupabaseClient {
    return (useNuxtApp() as unknown as SupabaseNuxtApp).$supabase
  }

  async function syncAndRedirect(displayName?: string, timezone?: string): Promise<Profile> {
    const body: SyncRequest = {}
    if (displayName)
      body.display_name = displayName
    // Browser zone so new profiles skip the Lagos default. The backend only
    // applies this on insert — existing custom timezones are never overwritten.
    body.timezone = timezone || detectBrowserTimezone()
    await api.post<Profile>('/auth/sync', {...body})
    const me = await api.get<Profile>('/me')
    await navigateTo(me.Role === 'admin' ? '/admin' : '/dashboard')
    return me
  }

  async function login(email: string, password: string): Promise<Profile> {
    const { data, error } = await client().auth.signInWithPassword({ email, password })
    if (error)
      throw new Error(error.message)
    return syncAndRedirect(data.user?.user_metadata?.display_name as string | undefined)
  }

  async function signup(email: string, password: string, displayName?: string, timezone?: string): Promise<void> {
    const tz = timezone || detectBrowserTimezone()
    const { data, error } = await client().auth.signUp({
      email,
      password,
      options: { data: { display_name: displayName ?? '', timezone: tz } },
    })
    if (error)
      throw new Error(error.message)
    // Supabase answers a duplicate email with fake success: an obfuscated
    // user with no identities and NO confirmation email (anti-enumeration).
    // Catch it here instead of stranding them on /auth/verify-email.
    const identities = (data.user as unknown as { identities?: unknown[] } | null)?.identities
    if (data.user && Array.isArray(identities) && identities.length === 0)
      // throw new Error('This email is already registered — log in instead.')
    throw new Error('An account with this email already exists. Try logging in instead.')
    if (!data.session) {
      await navigateTo('/auth/verify-email')
      return
    }
    await syncAndRedirect(displayName, tz)
  }

  async function logout(): Promise<void> {
    const nuxtApp = useNuxtApp() as unknown as SupabaseNuxtApp
    await nuxtApp.$supabase.auth.signOut()
    nuxtApp.$queryClient.clear()
    await navigateTo('/auth/login')
  }

  return {
    login,
    signup,
    logout,
    syncAndRedirect,
    me: () => api.get<Profile>('/me'),
  }
}
