/**
 * Trades. Money math (RiskAmount/Pnl/RMultiple) is SERVER-calculated —
 * the form preview is display only. Reopening a closed trade clears
 * its results server-side.
 */

export type Direction = 'long' | 'short';
export type TradeStatus = 'open' | 'closed';
export type TradeResult = 'win' | 'loss';
export type ImageKind = 'entry' | 'exit';

export interface Trade {
  ID: string;
  UserID: string;
  AccountID: string;
  SetupID: string | null;
  Pair: string;
  Direction: Direction;
  Timeframe: string | null;
  OpenedAt: string | null;
  ClosedAt: string | null;
  Entry: number | null;
  StopLoss: number | null;
  TakeProfit: number | null;
  ExitPrice: number | null;
  LotSize: number | null;
  Commission: number;
  Swap: number;
  RiskAmount: number | null;
  Pnl: number | null;
  RMultiple: number | null;
  FollowedRules: boolean | null;
  Emotion: string | null;
  Notes: string | null;
  Status: TradeStatus;
  IsPublic: boolean;
  HiddenByAdmin: boolean;
  CreatedAt: string;
  UpdatedAt: string;
}

/** POST /trades body. `status: closed` requires `exit_price`. */
export interface TradeCreate {
  account_id: string;
  setup_id?: string;
  pair: string;
  direction: Direction;
  timeframe?: string;
  opened_at?: string;
  closed_at?: string;
  entry?: number;
  stop_loss?: number;
  take_profit?: number;
  exit_price?: number;
  lot_size?: number;
  commission?: number;
  swap?: number;
  followed_rules?: boolean;
  emotion?: string;
  notes?: string;
  status?: TradeStatus;
  tag_ids?: string[];
}

/** PATCH /trades/{id} — all optional; `setup_clear` unlinks the setup,
 *  `tag_ids` null = unchanged, array = replace. */
export interface TradeUpdate {
  account_id?: string;
  setup_id?: string;
  setup_clear?: boolean;
  pair?: string;
  direction?: Direction;
  timeframe?: string;
  opened_at?: string;
  closed_at?: string;
  entry?: number;
  stop_loss?: number;
  take_profit?: number;
  exit_price?: number;
  lot_size?: number;
  commission?: number;
  swap?: number;
  followed_rules?: boolean;
  emotion?: string;
  notes?: string;
  status?: TradeStatus;
  tag_ids?: string[] | null;
}

/** PATCH /trades/{id}/visibility — the "Make public" switch (off by default). */
export interface VisibilityUpdate {
  is_public: boolean;
}

/** GET /trades list filters (same set drives the CSV export). */
export interface TradeFilters {
  pair?: string;
  setup?: string;
  account?: string;
  status?: TradeStatus | '';
  direction?: Direction | '';
  result?: TradeResult | '';
  tag?: string;
  from?: string;
  to?: string;
  limit?: number;
  offset?: number;
}
