/**
 * GET /dashboard — one call feeding the whole first screen.
 * Keys are CAPITALIZED (Go default, no json tags on the aggregate).
 */
import type { Account } from './accounts';
import type { Profile } from './auth';
import type { CalendarDay, EquityPoint, PairStat, SetupStat, StatsSummary } from './stats';
import type { Trade } from './trades';

export interface DashboardOverview {
  Me: Profile;
  Accounts: Account[];
  Summary: StatsSummary;
  Equity: EquityPoint[];
  Calendar: CalendarDay[];
  BySetup: SetupStat[];
  ByPair: PairStat[];
  /** Last 10 trades (account-scoped when `?account=` is set). */
  Recent: Trade[];
}
