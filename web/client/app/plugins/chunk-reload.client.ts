import { toast } from 'vue-sonner'
import { friendlyChunkMessage, isChunkLoadError } from '~/utils/chunk-error'

/**
 * Turns raw chunk-load failures ("Failed to fetch dynamically imported
 * module: .../_nuxt/...") into a readable toast + one-click refresh.
 *
 * Triggers: new deploy with stale hashed chunks, dev-server restart,
 * offline / flaky network. The raw module URL is never shown to users.
 */
export default defineNuxtPlugin((nuxtApp) => {
  if (import.meta.server)
    return

  const RELOAD_KEY = 'ctemz-chunk-reloaded-at'
  let toastShown = false

  function showReloadPrompt(): void {
    if (toastShown)
      return
    toastShown = true
    toast.error(friendlyChunkMessage(), {
      description: 'Your page is out of date.',
      duration: 15000,
      action: {
        label: 'Refresh',
        onClick: () => window.location.reload(),
      },
    })
    // Reset so a later, unrelated navigation can prompt again.
    window.setTimeout(() => {
      toastShown = false
    }, 15000)
  }

  function handleChunkError(): void {
    try {
      const last = Number(sessionStorage.getItem(RELOAD_KEY) ?? 0)
      // Auto-reload once if the failure looks like a fresh deploy and we
      // haven't just reloaded (prevents loops on persistent offline).
      if (!Number.isFinite(last) || Date.now() - last > 10000) {
        sessionStorage.setItem(RELOAD_KEY, String(Date.now()))
        window.location.reload()
        return
      }
    }
    catch {
      // sessionStorage may be blocked — fall through to manual prompt.
    }
    showReloadPrompt()
  }

  // Nuxt chunk hook (covers middleware / page lazy-import failures).
  try {
    nuxtApp.hook('app:chunkError', handleChunkError)
  }
  catch {
    // Older Nuxt — fallback listeners below still cover it.
  }

  // Router-level async chunk failures.
  nuxtApp.hook('app:created', () => {
    try {
      const router = useRouter()
      router.onError((err) => {
        if (isChunkLoadError(err))
          handleChunkError()
      })
    }
    catch {
      // Router not ready — window listeners still cover it.
    }
  })

  // Final safety net: preload errors + unhandled import() rejections.
  window.addEventListener('vite:preloadError', (event) => {
    event.preventDefault()
    handleChunkError()
  })
  window.addEventListener('unhandledrejection', (event) => {
    if (isChunkLoadError(event.reason)) {
      event.preventDefault()
      handleChunkError()
    }
  })
  window.addEventListener('error', (event) => {
    const target = event.target as HTMLElement | null
    if (target && (target.tagName === 'SCRIPT' || target.tagName === 'LINK'))
      showReloadPrompt()
  }, true)
})
