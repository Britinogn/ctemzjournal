/** Avatar storage helpers. Profile pictures live in the public Supabase
 *  `avatars` bucket at `{user_id}/avatar.webp`. `AvatarPath` in the DB may
 *  hold that relative path (new) or a full http(s) URL (legacy). */

export const AVATAR_BUCKET = 'avatars'

/** Storage path for a user's picture — one file per user, overwritten on change. */
export function avatarStoragePath(userId: string): string {
  return `${userId}/avatar.webp`
}

/**
 * Resolve `AvatarPath` to a renderable URL.
 * '' → '' (show initials), http(s) → as-is, otherwise public storage URL.
 * `rev` (e.g. profile UpdatedAt) is appended as `?v=` so an overwritten
 * picture at the same path is not served stale from cache.
 */
export function resolveAvatarUrl(supabaseUrl: string, path: string | null | undefined, rev?: string | null): string {
  const p = (path ?? '').trim()
  if (p === '')
    return ''
  const query = (rev ?? '').trim() ? `v=${encodeURIComponent((rev ?? '').trim())}` : ''
  const withRev = (url: string): string => {
    if (!query)
      return url
    return `${url}${url.includes('?') ? '&' : '?'}${query}`
  }
  if (/^https?:\/\//i.test(p))
    return withRev(p)
  const base = supabaseUrl.replace(/\/$/, '')
  if (!base)
    return ''
  return withRev(`${base}/storage/v1/object/public/${AVATAR_BUCKET}/${p.replace(/^\//, '')}`)
}
