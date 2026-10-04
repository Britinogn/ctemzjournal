<script setup lang="ts">
import { useEventListener, useMediaQuery } from '@vueuse/core'
import {
  Add01Icon,
  Home01Icon,
  Menu01Icon,
  Settings01Icon,
  Tag01Icon,
  Target01Icon,
  TradeUpIcon,
  Wallet01Icon,
} from '~/utils/icons'

const route = useRoute()
const { isActive } = useUserNav()

const tabs = [
  { label: 'Home', to: '/dashboard', icon: Home01Icon, exact: true, fab: false, more: false },
  { label: 'Trades', to: '/dashboard/trades', icon: TradeUpIcon, exact: false, fab: false, more: false },
  { label: 'Log', to: '/dashboard/trades/new', icon: Add01Icon, exact: true, fab: true, more: false },
  { label: 'Accounts', to: '/dashboard/accounts', icon: Wallet01Icon, exact: false, fab: false, more: false },
  { label: 'More', to: '/dashboard/settings', icon: Menu01Icon, exact: false, fab: false, more: true },
]

// The two pages the tab bar had no room for, plus Settings. "More" opens a sheet with these.
const moreItems = [
  { label: 'Setups', to: '/dashboard/setups', icon: Target01Icon, exact: false },
  { label: 'Tags', to: '/dashboard/tags', icon: Tag01Icon, exact: false },
  { label: 'Settings', to: '/dashboard/settings', icon: Settings01Icon, exact: false },
]

/* ---------- The "More" sheet ---------- */
const open = ref(false)
const panel = ref<HTMLElement | null>(null)
const moreBtn = ref<HTMLButtonElement | null>(null)

// "More" counts as the current tab on any page that lives inside it.
const moreActive = computed(() => moreItems.some(i => isActive(i)))
const tabActive = (tab: (typeof tabs)[number]) => (tab.more ? moreActive.value : isActive(tab))

function close(): void {
  open.value = false
}

// Close after navigating, and if the screen grows to tablet size (where the sidebar takes over).
watch(() => route.path, close)
const isWide = useMediaQuery('(min-width: 768px)')
watch(isWide, (wide) => {
  if (wide)
    close()
})

// Escape closes. Tab stays inside the sheet.
useEventListener('keydown', (e: KeyboardEvent) => {
  if (!open.value)
    return
  if (e.key === 'Escape') {
    close()
    return
  }
  if (e.key !== 'Tab' || !panel.value)
    return
  const items = [...panel.value.querySelectorAll<HTMLElement>('a[href], button:not([disabled])')]
  const first = items[0]
  const last = items[items.length - 1]
  if (!first || !last)
    return
  const active = document.activeElement
  if (!panel.value.contains(active) || (e.shiftKey && active === first)) {
    e.preventDefault()
    ;(e.shiftKey ? last : first).focus()
  }
  else if (!e.shiftKey && active === last) {
    e.preventDefault()
    first.focus()
  }
})

// On open: lock the page behind and focus the first link. On close: focus goes back to "More".
watch(open, async (isOpen) => {
  if (isOpen) {
    document.documentElement.style.overflow = 'hidden'
    await nextTick()
    panel.value?.querySelector<HTMLElement>('a[href]')?.focus()
  }
  else {
    document.documentElement.style.overflow = ''
    moreBtn.value?.focus()
  }
})
onBeforeUnmount(() => {
  document.documentElement.style.overflow = ''
})

const focusRing
  = 'focus-visible:outline-2-2 focus-visible:-outline-offset-2 focus-visible:outline-2-primary'
const tabCls
  = 'relative flex min-h-14 w-full flex-col items-center justify-center gap-0.5 rounded-xl text-[11px] font-medium transition-colors'
</script>

<template>
  <!-- Mobile bottom nav (tablet and desktop use the sidebar rail). -->
  <nav class="fixed inset-x-0 bottom-0 z-30 border-t border-border bg-surface/95 pb-[env(safe-area-inset-bottom)] backdrop-blur md:hidden" aria-label="Primary">
    <ul class="grid grid-cols-5 px-2">
      <li v-for="tab in tabs" :key="tab.to">
        <!-- "More" opens a sheet, so it is a button, not a link -->
        <button
          v-if="tab.more"
          ref="moreBtn"
          type="button"
          aria-haspopup="dialog"
          :aria-expanded="open"
          :aria-current="moreActive ? 'page' : undefined"
          :class="[tabCls, focusRing, tabActive(tab) || open ? 'font-semibold text-primary' : 'text-muted']"
          @click="open = !open"
        >
          <span v-if="tabActive(tab)" class="absolute -top-px h-0.5 w-8 rounded-b-full bg-primary" aria-hidden="true" />
          <UiAppIcon :icon="tab.icon" :size="22" />
          {{ tab.label }}
        </button>

        <NuxtLink
          v-else
          :to="tab.to"
          :aria-current="isActive(tab) ? 'page' : undefined"
          :class="[tabCls, focusRing, isActive(tab) ? 'font-semibold text-primary' : 'text-muted']"
        >
          <!-- Marker on the top edge, so the current tab is not shown by colour alone -->
          <span v-if="isActive(tab) && !tab.fab" class="absolute -top-px h-0.5 w-8 rounded-b-full bg-primary" aria-hidden="true" />

          <!-- The thick border in the bar's colour cuts the button out of the bar's top edge -->
          <span
            v-if="tab.fab"
            class="-mt-6 inline-flex h-12 w-12 items-center justify-center rounded-2xl border-4 border-surface bg-primary text-on-primary shadow-lg"
          >
            <UiAppIcon :icon="tab.icon" :size="22" />
          </span>
          <UiAppIcon v-else :icon="tab.icon" :size="22" />
          {{ tab.label }}
        </NuxtLink>
      </li>
    </ul>
  </nav>

  <!-- The sheet. Teleported to <body> so the bar's blur and z-index cannot trap it. -->
  <ClientOnly>
    <Teleport to="body">
      <Transition name="sheet">
        <div
          v-if="open"
          class="fixed inset-0 z-50 flex items-end bg-text/40 backdrop-blur-xs md:hidden"
          role="presentation"
          @click.self="close"
        >
          <div
            ref="panel"
            role="dialog"
            aria-modal="true"
            aria-labelledby="more-title"
            class="sheet-panel w-full rounded-t-3xl border-t border-border bg-surface px-4 pb-[calc(1rem+env(safe-area-inset-bottom))] pt-3 shadow-xl"
          >
            <div class="mx-auto mb-2 h-1 w-10 rounded-full bg-border" aria-hidden="true" />
            <div class="mb-2 flex items-center justify-between">
              <h2 id="more-title" class="text-base font-bold">
                More
              </h2>
              <button
                type="button"
                aria-label="Close"
                class="-mr-2 inline-flex h-11 w-11 items-center justify-center rounded-xl text-muted transition-colors hover:bg-bg hover:text-text focus-visible:outline-2-2 focus-visible:outline-2-offset-2 focus-visible:outline-2-primary"
                @click="close"
              >
                <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                  <path d="M6 6l12 12M18 6L6 18" />
                </svg>
              </button>
            </div>

            <ul class="divide-y divide-border overflow-hidden rounded-2xl border border-border">
              <li v-for="item in moreItems" :key="item.to">
                <NuxtLink
                  :to="item.to"
                  :aria-current="isActive(item) ? 'page' : undefined"
                  :class="[
                    'flex min-h-14 items-center gap-3 px-4 text-sm font-semibold transition-colors hover:bg-bg focus-visible:outline-2-2 focus-visible:-outline-offset-2 focus-visible:outline-2-primary',
                    isActive(item) ? 'bg-primary/10 text-primary' : '',
                  ]"
                >
                  <span class="inline-flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
                    <UiAppIcon :icon="item.icon" :size="18" />
                  </span>
                  <span class="flex-1">{{ item.label }}</span>
                  <svg class="h-4 w-4 text-muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                    <path d="M9 5l7 7-7 7" />
                  </svg>
                </NuxtLink>
              </li>
            </ul>
          </div>
        </div>
      </Transition>
    </Teleport>
  </ClientOnly>
</template>

<style scoped>
.sheet-enter-active,
.sheet-leave-active {
  transition: opacity 0.16s ease;
}
.sheet-enter-active .sheet-panel,
.sheet-leave-active .sheet-panel {
  transition: transform 0.22s cubic-bezier(0.22, 0.8, 0.3, 1);
}
.sheet-enter-from,
.sheet-leave-to {
  opacity: 0;
}
.sheet-enter-from .sheet-panel,
.sheet-leave-to .sheet-panel {
  transform: translateY(100%);
}
@media (prefers-reduced-motion: reduce) {
  .sheet-enter-active,
  .sheet-leave-active,
  .sheet-enter-active .sheet-panel,
  .sheet-leave-active .sheet-panel {
    transition: none;
  }
}
</style>