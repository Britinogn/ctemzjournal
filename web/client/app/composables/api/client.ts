import { toast } from 'vue-sonner'
import type { SupabaseClient } from '@supabase/supabase-js'
import type { ApiData } from '~/types'

type Method = 'GET' | 'POST' | 'PATCH' | 'PUT' | 'DELETE'
type ApiBody = Record<string, unknown> | Array<unknown> | string | number | boolean | null

interface ApiClient {
  get: <T>(path: string) => Promise<T>;
  post: <T>(path: string, body?: ApiBody) => Promise<T>;
  patch: <T>(path: string, body?: ApiBody) => Promise<T>;
  put: <T>(path: string, body?: ApiBody) => Promise<T>;
  del: <T>(path: string) => Promise<T>;
}

/**
 * Authenticated Go API client. Attaches the Supabase access token,
 * unwraps the `{ data }` envelope, and handles auth failures in one place:
 * 401 → sign out + /login, 403 → /suspended, 429 → toast.
 */
export function useApi(): ApiClient {
  const baseURL = useRuntimeConfig().public.apiUrl as string

  async function request<T>(method: Method, path: string, body?: ApiBody): Promise<T> {
    const { $supabase } = useNuxtApp() as unknown as { $supabase: SupabaseClient }
    const { data: { session } } = await $supabase.auth.getSession()
    try {
      const res = await $fetch<ApiData<T>>(path, {
        baseURL,
        method,
        body: body as Record<string, unknown> | undefined,
        headers: session?.access_token ? { Authorization: `Bearer ${session.access_token}` } : {},
      })
      return res.data
    }
    catch (err: unknown) {
      const status = (err as { status?: number })?.status
        ?? (err as { response?: { status?: number } })?.response?.status
      if (status === 401) {
        await $supabase.auth.signOut()
        await navigateTo('/auth/login')
      }
      else if (status === 403) {
        await navigateTo('/auth/suspended')
      }
      else if (status === 429) {
        toast.error('Too many requests, try again shortly')
      }
      throw err
    }
  }

  return {
    get: <T>(path: string) => request<T>('GET', path),
    post: <T>(path: string, body?: ApiBody) => request<T>('POST', path, body),
    patch: <T>(path: string, body?: ApiBody) => request<T>('PATCH', path, body),
    put: <T>(path: string, body?: ApiBody) => request<T>('PUT', path, body),
    del: <T>(path: string) => request<T>('DELETE', path),
  }
}
