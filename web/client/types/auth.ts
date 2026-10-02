/** Auth + profile types. Login/session itself is Supabase Auth (email);
 * the backend only verifies the JWT and resolves the role from the DB. */

export type Role = 'user' | 'admin';
export type ProfileStatus = 'active' | 'suspended';

/** GET /me — also the redirect source (admin → /admin, user → /dashboard). */
export interface Profile {
  ID: string;
  DisplayName: string | null;
  Role: Role;
  Status: ProfileStatus;
  Timezone: string;
  AvatarPath: string | null;
  CreatedAt: string;
  UpdatedAt: string;
}

/** PATCH /me body — all optional. */
export interface ProfileUpdate {
  display_name?: string;
  timezone?: string;
  avatar_path?: string;
}

/** POST /auth/sync body (fallback profile bootstrap on first login). */
export interface SyncRequest {
  display_name?: string;
}

export interface LoginInput {
  email: string;
  password: string;
}

export interface SignupInput {
  email: string;
  password: string;
  display_name?: string;
}

export interface ResetRequestInput {
  email: string;
}

export interface NewPasswordInput {
  password: string;
}
