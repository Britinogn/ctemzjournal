/**
 * Admin area. Every write is audit-logged server-side.
 * NOTE: `meta` arrives base64-encoded (jsonb bytes) — decode before display.
 * Flagged as a backend fix candidate (decode to object server-side).
 */
import type { Profile } from './auth';

export interface SignupDay {
  date: string;
  count: number;
}

export interface AdminOverview {
  user_count: number;
  trade_count: number;
  new_users_this_week: number;
  active_users: number;
  suspended_users: number;
  public_trades: number;
  hidden_journals: number;
  new_users_prev_week: number;
  signups_last_7d: SignupDay[];
}

/** GET /admin/users row: profile + email + trade count. */
export interface AdminUser {
  ID: string;
  DisplayName: string | null;
  Email: string | null;
  Role: 'user' | 'admin';
  Status: 'active' | 'suspended';
  Timezone: string;
  AvatarPath: string | null;
  CreatedAt: string;
  UpdatedAt: string;
  TradeCount: number;
}

/** PATCH /admin/users/{id}/status */
export interface UserStatusUpdate {
  status: 'active' | 'suspended';
}

/** PATCH /admin/journals/{id}/hide — false restores. */
export interface JournalHide {
  hidden: boolean;
}

/** PATCH /admin/site-settings — all optional; social_links presence
 *  (even `{}`) replaces, absence leaves unchanged. */
export interface SettingsUpdate {
  site_name?: string;
  tagline?: string;
  logo_path?: string;
  favicon_path?: string;
  contact_email?: string;
  footer_text?: string;
  risk_disclaimer?: string;
  social_links?: Record<string, unknown>;
  allow_signups?: boolean;
  maintenance_mode?: boolean;
}

/** POST /admin/site-settings/logo — reserves the site-assets path the
 *  browser uploads to itself (Supabase JS), then PATCH saves the path. */
export interface LogoRequest {
  kind: 'logo' | 'favicon';
  extension: string;
}

export interface LogoTarget {
  bucket: 'site-assets';
  path: string;
}

export interface AuditEntry {
  ID: string;
  AdminID: string | null;
  Action: string;
  TargetType: string;
  TargetID: string;
  /** Base64-encoded JSON — decode before display (see note above). */
  Meta: string;
  CreatedAt: string;
}
