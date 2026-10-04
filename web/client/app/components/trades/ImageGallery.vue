<script setup lang="ts">
import { useEventListener } from '@vueuse/core'

// Only the fields this component needs, so it accepts your TradeImageWithUrl as it is.
interface GalleryImage {
  ID: string
  Kind?: string | null
}

const props = withDefaults(defineProps<{
  images: readonly GalleryImage[]
  /** Shows a red Remove button under each image (the edit page). */
  removable?: boolean
  removingId?: string | null
  gridClass?: string
}>(), {
  removable: false,
  removingId: null,
  gridClass: 'grid-cols-2',
})

const emit = defineEmits<{ remove: [id: string] }>()

/**
 * The address of the picture. It reads URL first, then url and Url, because the key the API sends
 * can differ from the one in the TypeScript type, and an empty address is what makes the alt text show.
 */
function srcOf(img: object): string {
  const r = img as Record<string, unknown>
  const v = r.URL ?? r.url ?? r.Url
  return typeof v === 'string' ? v : ''
}
const kindOf = (img: GalleryImage) => img.Kind ?? 'chart'

// Pictures that fail to load get a tidy tile instead of the browser's broken-image alt text.
const failed = ref<Record<string, boolean>>({})
const broken = (img: GalleryImage) => !srcOf(img) || failed.value[img.ID]

/* ---------- Viewer ---------- */
const current = ref<number | null>(null)
const loaded = ref(false)
const panel = ref<HTMLElement | null>(null)
const closeBtn = ref<HTMLButtonElement | null>(null)
let opener: HTMLElement | null = null

const active = computed(() => (current.value === null ? null : props.images[current.value] ?? null))
const many = computed(() => props.images.length > 1)

function openAt(i: number, e: MouseEvent): void {
  opener = e.currentTarget as HTMLElement
  loaded.value = false
  current.value = i
}
function close(): void {
  current.value = null
}
function go(delta: number): void {
  if (current.value === null || !many.value)
    return
  loaded.value = false
  current.value = (current.value + delta + props.images.length) % props.images.length
}

useEventListener('keydown', (e: KeyboardEvent) => {
  if (current.value === null)
    return
  if (e.key === 'Escape') {
    close()
    return
  }
  if (e.key === 'ArrowLeft')
    go(-1)
  else if (e.key === 'ArrowRight')
    go(1)
  else if (e.key === 'Tab' && panel.value) {
    const items = [...panel.value.querySelectorAll<HTMLElement>('a[href], button:not([disabled])')]
    const first = items[0]
    const last = items[items.length - 1]
    if (!first || !last)
      return
    const el = document.activeElement
    if (!panel.value.contains(el) || (e.shiftKey && el === first)) {
      e.preventDefault()
      ;(e.shiftKey ? last : first).focus()
    }
    else if (!e.shiftKey && el === last) {
      e.preventDefault()
      first.focus()
    }
  }
})

// Swipe left or right on a phone.
let touchX = 0
function onTouchStart(e: TouchEvent): void {
  touchX = e.changedTouches[0]?.clientX ?? 0
}
function onTouchEnd(e: TouchEvent): void {
  const dx = (e.changedTouches[0]?.clientX ?? 0) - touchX
  if (Math.abs(dx) > 50)
    go(dx < 0 ? 1 : -1)
}

// Lock the page behind while open, and give focus back to the thumbnail when it closes.
watch(current, async (now, before) => {
  if (now !== null && before === null) {
    document.documentElement.style.overflow = 'hidden'
    await nextTick()
    closeBtn.value?.focus()
  }
  else if (now === null) {
    document.documentElement.style.overflow = ''
    opener?.focus?.()
    opener = null
  }
})
onBeforeUnmount(() => {
  document.documentElement.style.overflow = ''
})

const ringCls = 'focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary'
const roundBtn
  = 'inline-flex h-11 w-11 items-center justify-center rounded-full border border-border bg-surface text-text shadow-sm transition-colors hover:border-primary'
</script>

<template>
  <ul class="grid gap-2" :class="gridClass">
    <li
      v-for="(img, i) in images"
      :key="img.ID"
      class="overflow-hidden rounded-xl border border-border bg-bg"
    >
      <button
        type="button"
        class="group relative block w-full focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-primary"
        :aria-label="`View ${kindOf(img)} screenshot larger`"
        @click="openAt(i, $event)"
      >
        <img
          v-if="!broken(img)"
          :src="srcOf(img)"
          alt=""
          class="aspect-video w-full object-cover transition-opacity group-hover:opacity-90"
          loading="lazy"
          decoding="async"
          @error="failed[img.ID] = true"
        >
        <span v-else class="flex aspect-video flex-col items-center justify-center gap-1 px-2 text-center text-xs text-muted">
          <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <rect x="3" y="4" width="18" height="16" rx="2" /><circle cx="9" cy="10" r="1.5" /><path d="M21 16l-5-5-8 8" />
          </svg>
          Image not available
        </span>

        <!-- A small magnifier shows on hover, so it is clear the picture opens -->
        <span v-if="!broken(img)" class="pointer-events-none absolute right-2 top-2 inline-flex h-7 w-7 items-center justify-center rounded-full bg-surface/90 text-text opacity-0 shadow-sm transition-opacity group-hover:opacity-100 group-focus-visible:opacity-100" aria-hidden="true">
          <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="6" /><path d="M20 20l-4-4M11 8v6M8 11h6" />
          </svg>
        </span>
      </button>

      <div class="flex items-center justify-between gap-2 px-2.5 py-1.5" :class="removable ? '' : 'justify-center'">
        <span class="text-xs font-medium capitalize text-muted">{{ kindOf(img) }}</span>
        <button
          v-if="removable"
          type="button"
          :disabled="removingId === img.ID"
          :aria-label="`Remove ${kindOf(img)} screenshot`"
          :class="['inline-flex h-9 items-center rounded-lg border border-loss/40 px-3 text-xs font-semibold text-loss transition-colors hover:bg-loss/10 disabled:cursor-wait disabled:opacity-50', 'focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-loss']"
          @click="emit('remove', img.ID)"
        >
          {{ removingId === img.ID ? 'Removing…' : 'Remove' }}
        </button>
      </div>
    </li>
  </ul>

  <!-- The viewer. Teleported to <body> so nothing can clip it. -->
  <ClientOnly>
    <Teleport to="body">
      <Transition name="lb">
        <div
          v-if="active"
          class="fixed inset-0 z-50 bg-bg/95 backdrop-blur-sm"
          role="dialog"
          aria-modal="true"
          aria-label="Screenshot viewer"
        >
          <div ref="panel" class="flex h-full flex-col">
            <!-- Top bar -->
            <div class="flex items-center justify-between gap-3 px-4 pb-2 pt-[max(0.75rem,env(safe-area-inset-top))] md:px-6">
              <p class="min-w-0 truncate text-sm font-semibold capitalize">
                {{ kindOf(active) }}
                <span v-if="many" class="tnum ml-1 font-medium normal-case text-muted">{{ (current ?? 0) + 1 }} / {{ images.length }}</span>
              </p>
              <div class="flex items-center gap-2">
                <a
                  v-if="srcOf(active)"
                  :href="srcOf(active)"
                  target="_blank"
                  rel="noopener"
                  :class="[roundBtn, ringCls]"
                  aria-label="Open full size in a new tab"
                >
                  <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                    <path d="M14 4h6v6M20 4l-9 9M18 14v5a1 1 0 01-1 1H5a1 1 0 01-1-1V7a1 1 0 011-1h5" />
                  </svg>
                </a>
                <button ref="closeBtn" type="button" :class="[roundBtn, ringCls]" aria-label="Close" @click="close">
                  <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                    <path d="M6 6l12 12M18 6L6 18" />
                  </svg>
                </button>
              </div>
            </div>

            <!-- Picture. A tap on the empty space around it closes the viewer. -->
            <div
              class="relative flex min-h-0 flex-1 items-center justify-center px-4 pb-[max(1rem,env(safe-area-inset-bottom))] md:px-16"
              @click.self="close"
              @touchstart.passive="onTouchStart"
              @touchend.passive="onTouchEnd"
            >
              <button v-if="many" type="button" :class="[roundBtn, ringCls, 'absolute left-2 z-10 md:left-4']" aria-label="Previous screenshot" @click="go(-1)">
                <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M15 5l-7 7 7 7" /></svg>
              </button>

              <img
                v-if="srcOf(active) && !failed[active.ID]"
                :key="active.ID"
                :src="srcOf(active)"
                :alt="`${kindOf(active)} screenshot`"
                class="max-h-full max-w-full rounded-xl border border-border object-contain shadow-xl transition-opacity duration-150"
                :class="loaded ? 'opacity-100' : 'opacity-0'"
                @load="loaded = true"
                @error="failed[active.ID] = true"
              >
              <p v-else class="text-sm text-muted">
                Image not available
              </p>

              <!-- Spinner while the large picture arrives -->
              <svg
                v-if="srcOf(active) && !failed[active.ID] && !loaded"
                class="absolute h-6 w-6 animate-spin text-muted motion-reduce:animate-none"
                viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" aria-hidden="true"
              >
                <path d="M12 3a9 9 0 109 9" />
              </svg>

              <button v-if="many" type="button" :class="[roundBtn, ringCls, 'absolute right-2 z-10 md:right-4']" aria-label="Next screenshot" @click="go(1)">
                <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M9 5l7 7-7 7" /></svg>
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </ClientOnly>
</template>

<style scoped>
.lb-enter-active,
.lb-leave-active {
  transition: opacity 0.15s ease;
}
.lb-enter-from,
.lb-leave-to {
  opacity: 0;
}
@media (prefers-reduced-motion: reduce) {
  .lb-enter-active,
  .lb-leave-active {
    transition: none;
  }
}
</style>