import {
  DashboardSquare02Icon,
  GlobeIcon,
  Settings01Icon,
  ShieldCheckIcon,
  UsersIcon,
} from '~/utils/icons'

export interface AdminNavItem {
  label: string
  /** Shorter label for the phone bottom tabs. */
  short: string
  to: string
  icon: typeof DashboardSquare02Icon
  exact: boolean
  /** Only controls where the sidebar dividers go. */
  group: number
}

/** One list for the sidebar, the top bar and the bottom tabs, so they never disagree. */
export const ADMIN_NAV: AdminNavItem[] = [
  { label: 'Overview', short: 'Overview', to: '/admin', icon: DashboardSquare02Icon, exact: true, group: 0 },
  { label: 'Users', short: 'Users', to: '/admin/users', icon: UsersIcon, exact: false, group: 1 },
  { label: 'Journals', short: 'Journals', to: '/admin/journals', icon: GlobeIcon, exact: false, group: 1 },
  { label: 'Site settings', short: 'Site', to: '/admin/settings', icon: Settings01Icon, exact: false, group: 2 },
  { label: 'Audit log', short: 'Audit', to: '/admin/audit', icon: ShieldCheckIcon, exact: false, group: 2 },
]

/** The order the phone tabs have always used. */
const TAB_ORDER = ['/admin', '/admin/users', '/admin/journals', '/admin/audit', '/admin/settings']
export const ADMIN_TABS: AdminNavItem[] = TAB_ORDER
  .map(to => ADMIN_NAV.find(n => n.to === to))
  .filter((n): n is AdminNavItem => !!n)

export function useAdminNav() {
  const route = useRoute()

  function isActive(item: Pick<AdminNavItem, 'to' | 'exact'>): boolean {
    if (item.exact)
      return route.path === item.to
    return route.path === item.to || route.path.startsWith(`${item.to}/`)
  }

  /** The section the admin is in, for example Users while on /admin/users/123. */
  const current = computed(() => ADMIN_NAV.find(isActive) ?? null)

  /** True on pages below a section, where a back link helps. */
  const isChild = computed(() => {
    const c = current.value
    return !!c && route.path.replace(/\/$/, '') !== c.to
  })

  return { nav: ADMIN_NAV, tabs: ADMIN_TABS, isActive, current, isChild }
}