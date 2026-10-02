import { useQuery } from '@tanstack/vue-query'
import { accountsKey, setupsKey, tagsKey, type Account, type Setup, type Tag } from '~/types'

const LONG_STALE = 5 * 60_000;

/** Cached reference lists for selects and filters. */
export function useAccounts() {
  const api = useApi()
  return useQuery({
    queryKey: accountsKey(),
    queryFn: () => api.get<Account[]>('/accounts'),
    staleTime: LONG_STALE,
  })
}

export function useSetups() {
  const api = useApi()
  return useQuery({
    queryKey: setupsKey(),
    queryFn: () => api.get<Setup[]>('/setups'),
    staleTime: LONG_STALE,
  })
}

export function useTags() {
  const api = useApi()
  return useQuery({
    queryKey: tagsKey(),
    queryFn: () => api.get<Tag[]>('/tags'),
    staleTime: LONG_STALE,
  })
}
