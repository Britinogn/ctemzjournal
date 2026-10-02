/**
 * API envelope + shared transport types.
 *
 * IMPORTANT — key casing differs by origin:
 * - Database rows (Profile, Account, Trade, …) serialize with Go's default
 *   CAPITALIZED keys: `{ "ID": "...", "UserID": "...", "Pair": "EUR/USD" }`.
 *   Nullable columns arrive as `null` (pgtype renders invalid values null).
 * - Hand-built DTOs (stats, public, admin, signature) use snake_case keys
 *   from explicit `json` tags: `{ "win_rate": 0.6 }`.
 * The types below follow each shape exactly — do not "normalize" casing
 * client-side; the backend contract is the source of truth.
 */
import type { TradeFilters } from './trades';

/** Success envelope: every 200/201 response. */
export interface ApiData<T> {
  data: T;
}

/** Error envelope: `{ "error": { "message": "..." } }`. */
export interface ApiError {
  error: {
    message: string;
  };
}

/** List query: `?limit=` clamped 1–100 by the server (default 20), `?offset=` (default 0). */
export interface PageParams {
  limit?: number;
  offset?: number;
}

/** TanStack Query key factory inputs. */
export type QueryKey = readonly unknown[];

export function dashboardKey(accountId?: string): QueryKey {
  return ['dashboard', accountId ?? 'all'] as const;
}

export function meKey(): QueryKey {
  return ['me'] as const;
}

export function ratesKey(): QueryKey {
  return ['rates'] as const;
}

export function siteSettingsKey(): QueryKey {
  return ['site-settings'] as const;
}

export function journalsKey(limit = 20): QueryKey {
  return ['journals', limit] as const;
}

export function tradesKey(filters?: TradeFilters): QueryKey {
  return ['trades', filters ?? {}] as const;
}

export function accountsKey(): QueryKey {
  return ['accounts'] as const;
}

export function setupsKey(): QueryKey {
  return ['setups'] as const;
}

export function tagsKey(): QueryKey {
  return ['tags'] as const;
}

export function statsKey(accountId?: string): QueryKey {
  return ['stats', accountId ?? 'all'] as const;
}

export function adminOverviewKey(): QueryKey {
  return ['admin', 'overview'] as const;
}

export function adminUsersKey(search = '', page = 0): QueryKey {
  return ['admin', 'users', search, page] as const;
}

export function adminJournalsKey(page = 0): QueryKey {
  return ['admin', 'journals', page] as const;
}

export function adminAuditKey(page = 0): QueryKey {
  return ['admin', 'audit', page] as const;
}
