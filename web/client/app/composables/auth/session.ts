import type { QueryClient } from '@tanstack/vue-query'
import type { SupabaseClient } from '@supabase/supabase-js'
import type { Profile, SyncRequest } from '~/types'

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

  async function syncAndRedirect(displayName?: string): Promise<Profile> {
    const body: SyncRequest = displayName ? { display_name: displayName } : {}
    await api.post<Profile>('/auth/sync', body)
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

  async function signup(email: string, password: string, displayName?: string): Promise<void> {
    const { data, error } = await client().auth.signUp({
      email,
      password,
      options: { data: { display_name: displayName ?? '' } },
    })
    if (error)
      throw new Error(error.message)
    if (!data.session) {
      await navigateTo('/auth/verify-email')
      return
    }
    await syncAndRedirect(displayName)
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
