/**
 * Stats + dashboard aggregates. All accept `?account=`; all math is
 * server-side. Calendar days are grouped in Lagos time by the backend —
 * the browser never regroups. Money here is quote-currency units.
 */

export interface StatsSummary {
  total: number;
  wins: number;
  losses: number;
  win_rate: number;
  avg_r: number;
  expectancy: number;
  total_pnl: number;
  max_drawdown: number;
  rule_rate: number;
  open_trades: number;
  total_trades: number;
}

export interface EquityPoint {
  date: string;
  equity: number;
}

export interface CalendarDay {
  date: string;
  trades: number;
  wins: number;
  losses: number;
  pnl: number;
}

export interface SetupStat {
  setup_id: string | null;
  /** Null-setup group arrives as "No setup". */
  setup_name: string;
  trades: number;
  pnl: number;
  avg_r: number;
}

export interface PairStat {
  pair: string;
  trades: number;
  pnl: number;
  avg_r: number;
}
