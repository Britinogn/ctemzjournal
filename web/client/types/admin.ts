/**
 * Admin area. Every write is audit-logged server-side.
 * NOTE: `meta` arrives base64-encoded (jsonb bytes) — decode before display.
 * Flagged as a backend fix candidate (decode to object server-side).
 */
import type { Profile } from './auth';

export interface AdminOverview {
  user_count: number;
  trade_count: number;
  new_users_this_week: number;
}

export type AdminUser = Profile;

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
