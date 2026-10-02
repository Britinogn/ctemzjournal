/**
 * GET /dashboard — one call feeding the whole first screen.
 * Keys are snake_case (explicit `json` tags on the backend aggregate),
 * matching the live response exactly.
 */
import type { Account } from './accounts';
import type { Profile } from './auth';
import type { CalendarDay, EquityPoint, PairStat, SetupStat, StatsSummary } from './stats';
import type { Trade } from './trades';

export interface DashboardOverview {
  me: Profile;
  accounts: Account[];
  summary: StatsSummary;
  equity_curve: EquityPoint[];
  calendar: CalendarDay[];
  by_setup: SetupStat[];
  by_pair: PairStat[];
  /** Last 10 trades (account-scoped when `?account=` is set). */
  recent_trades: Trade[];
}
