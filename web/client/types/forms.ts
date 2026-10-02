/**
 * Form + draft states. Mirrors TradeCreate but form-friendly:
 * selects as plain strings, numbers as `number | undefined`,
 * tag ids as string array. The draft autosaves to localStorage.
 */
import type { Direction, ImageKind, TradeStatus } from './trades';
import type { AccountCurrency, AccountType } from './accounts';
import type { TagKind } from './tags';

export interface TradeDraft {
  account_id: string;
  setup_id: string;
  pair: string;
  direction: Direction;
  timeframe: string;
  opened_at: string;
  closed_at: string;
  entry?: number;
  stop_loss?: number;
  take_profit?: number;
  exit_price?: number;
  lot_size?: number;
  commission?: number;
  swap?: number;
  followed_rules?: boolean;
  emotion: string;
  notes: string;
  status: TradeStatus;
  tag_ids: string[];
  make_public: boolean;
}

/** Live risk/R preview (display only — server computes the real values). */
export interface RiskPreview {
  risk: number | null;
  r_multiple: number | null;
}

export interface PendingImage {
  /** local preview URL before upload */
  preview: string;
  file: File;
  kind: ImageKind;
  progress: number;
  error?: string;
}

export interface AccountForm {
  name: string;
  type: AccountType;
  currency: AccountCurrency;
  starting_balance?: number;
}

export interface SetupForm {
  name: string;
  rules: string;
  invalidation: string;
}

export interface TagForm {
  name: string;
  kind: TagKind;
}
