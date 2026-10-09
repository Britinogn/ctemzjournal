/**
 * Chunk-load failure helpers.
 *
 * Browsers throw raw messages like:
 * "Failed to fetch dynamically imported module: http://.../_nuxt/...".
 * That happens after a new deploy (old hashed chunk gone), a dev-server
 * restart, offline, or a blocked CDN request. Never show the raw URL —
 * map it to a human message + a reload action instead.
 */

const CHUNK_ERROR_PATTERN
  = /failed to fetch dynamically imported module|importing a module script failed|loading chunk [\w-]+ failed|chunkloaderror|failed to load module/i

export function isChunkLoadError(err: unknown): boolean {
  const message
    = err instanceof Error
      ? `${err.message} ${err.stack ?? ''}`
      : typeof err === 'string'
        ? err
        : (() => {
            try {
              return JSON.stringify(err) ?? ''
            }
            catch {
              return ''
            }
          })()
  return CHUNK_ERROR_PATTERN.test(message)
}

export function friendlyChunkMessage(): string {
  return 'A new version is available or your connection dropped. Refresh to continue.'
}
