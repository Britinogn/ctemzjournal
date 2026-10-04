import {
  DashboardSquare02Icon,
  Settings01Icon,
  Tag01Icon,
  Target01Icon,
  TradeUpIcon,
  Wallet01Icon,
} from '~/utils/icons'

export interface UserNavItem {
  label: string
  to: string
  icon: typeof DashboardSquare02Icon
  exact: boolean
  /** Only controls where the sidebar dividers go. */
  group: number
}

/** One list for the sidebar and the top bar, so they never disagree. */
export const USER_NAV: UserNavItem[] = [
  { label: 'Dashboard', to: '/dashboard', icon: DashboardSquare02Icon, exact: true, group: 0 },
  { label: 'Trades', to: '/dashboard/trades', icon: TradeUpIcon, exact: false, group: 0 },
  { label: 'Accounts', to: '/dashboard/accounts', icon: Wallet01Icon, exact: false, group: 1 },
  { label: 'Setups', to: '/dashboard/setups', icon: Target01Icon, exact: false, group: 1 },
  { label: 'Tags', to: '/dashboard/tags', icon: Tag01Icon, exact: false, group: 1 },
  { label: 'Settings', to: '/dashboard/settings', icon: Settings01Icon, exact: false, group: 2 },
]

export function useUserNav() {
  const route = useRoute()

  function isActive(item: Pick<UserNavItem, 'to' | 'exact'>): boolean {
    if (item.exact)
      return route.path === item.to
    return route.path === item.to || route.path.startsWith(`${item.to}/`)
  }

  /** The section the visitor is in, for example Trades while on /dashboard/trades/new. */
  const current = computed(() => USER_NAV.find(isActive) ?? null)

  /** True on pages below a section (new trade, trade detail), where a back link helps. */
  const isChild = computed(() => {
    const c = current.value
    return !!c && route.path.replace(/\/$/, '') !== c.to
  })

  return { nav: USER_NAV, isActive, current, isChild }
}