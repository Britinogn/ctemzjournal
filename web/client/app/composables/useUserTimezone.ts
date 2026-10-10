import { useQuery } from '@tanstack/vue-query'
import type { Profile } from '~/types'
import { meKey } from '~/types'
import { DEFAULT_TIMEZONE, detectBrowserTimezone } from '~/utils/timezones'

/**
 * The viewer's display timezone: profile setting first, then the browser
 * zone, then Lagos. Same `/me` key + stale time as every other reader, so
 * this adds no extra requests.
 */
export function useUserTimezone() {
  const api = useApi()
  const { data: me } = useQuery({
    queryKey: meKey(),
    queryFn: () => api.get<Profile>('/me'),
    staleTime: 5 * 60_000,
  })
  const timezone = computed(() => me.value?.Timezone || detectBrowserTimezone() || DEFAULT_TIMEZONE)
  return { timezone, me }
}
