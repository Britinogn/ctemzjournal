import { createClient, type SupabaseClient } from '@supabase/supabase-js'

/**
 * Supabase client — login, signup, session and the admin logo upload ONLY.
 * All API data goes through the Go API (TanStack Query), never Supabase.
 */
export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  const supabase: SupabaseClient = createClient(
    config.public.supabaseUrl as string,
    config.public.supabaseAnonKey as string,
  )
  return { provide: { supabase } }
})
