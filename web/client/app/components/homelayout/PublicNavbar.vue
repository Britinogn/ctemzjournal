<script setup lang="ts">
import { onClickOutside, useEventListener, useWindowScroll } from '@vueuse/core'
import { Cancel01Icon, Menu01Icon } from '~/utils/icons'

withDefaults(defineProps<{
  siteName?: string
  logoUrl?: string | null
  allowSignups?: boolean
  /** Logged-in target (/admin or /dashboard); null for guests. */
  dashboardTo?: string | null
}>(), {
  siteName: 'Ctemz Journal',
  logoUrl: null,
  allowSignups: true,
  dashboardTo: null,
})

const route = useRoute()
const open = ref(false)
const menuRef = ref<HTMLElement | null>(null)
const toggleRef = ref<HTMLButtonElement | null>(null)

// Links should point at routes that exist. If /prices and /about are not built
// yet, use home-page anchors instead (for example '/#prices').
const links = [
  { label: 'Journals', to: '/journals' },
  { label: 'Live prices', to: '/prices' },
  { label: 'About', to: '/about' },
]

// Nested routes keep their parent highlighted (/journals/123 -> Journals).
function isActive(to: string): boolean {
  return route.path === to || route.path.startsWith(`${to}/`)
}

// A soft shadow appears only once the page has scrolled under the header.
const { y } = useWindowScroll()
const scrolled = computed(() => y.value > 4)

// Close the menu on navigation, outside tap and Escape.
watch(() => route.path, () => { open.value = false })
onClickOutside(menuRef, () => { open.value = false }, { ignore: [toggleRef] })
useEventListener('keydown', (e: KeyboardEvent) => {
  if (e.key === 'Escape' && open.value) {
    open.value = false
    toggleRef.value?.focus()
  }
})

const primaryBtn
  = 'inline-flex items-center justify-center rounded-xl bg-primary px-4 text-sm font-semibold text-on-primary transition-colors hover:opacity-90'
</script>

<template>
  <header
    class="sticky top-0 z-40 border-b border-border bg-surface/95 pt-[env(safe-area-inset-top)] backdrop-blur transition-shadow"
    :class="scrolled ? 'shadow-sm' : ''"
  >
    <!-- Three columns so the nav sits in the true centre of the page -->
    <div class="mx-auto grid h-16 w-full max-w-6xl grid-cols-[1fr_auto] items-center gap-3 px-4 md:grid-cols-[1fr_auto_1fr] md:px-6">
      
      <NuxtLink to="/" class="flex w-fit min-w-0 items-center gap-2.5">
        <BrandLogo :logo-url="logoUrl" :site-name="siteName" :size="32" />
        <span class="truncate text-base font-extrabold tracking-tight text-text">
          {{ siteName }}
        </span>
      </NuxtLink>

      <nav class="hidden items-center gap-1 md:flex" aria-label="Public">
        <NuxtLink
          v-for="l in links"
          :key="l.to"
          :to="l.to"
          :aria-current="isActive(l.to) ? 'page' : undefined"
          class="relative rounded-lg px-3 py-2 text-sm font-medium transition-colors"
          :class="isActive(l.to) ? 'text-primary' : 'text-muted hover:text-text'"
        >
          {{ l.label }}
          <!-- Active marker, so the current page is not shown by colour alone -->
          <span
            v-if="isActive(l.to)"
            class="absolute inset-x-3 -bottom-[13px] h-0.5 rounded-full bg-primary"
            aria-hidden="true"
          />
        </NuxtLink>
      </nav>

      <div class="flex items-center justify-end gap-2">
        <ThemeToggle />

        <!-- Desktop actions -->
        <div class="hidden items-center gap-1 md:flex">
          <NuxtLink v-if="dashboardTo" :to="dashboardTo" :class="[primaryBtn, 'h-10']">
            Dashboard
          </NuxtLink>
          <template v-else>
            <NuxtLink
              to="/auth/login"
              class="inline-flex h-10 items-center rounded-xl px-4 text-sm font-semibold text-text transition-colors hover:bg-bg"
            >
              Log in
            </NuxtLink>
            <NuxtLink v-if="allowSignups" to="/auth/signup" :class="[primaryBtn, 'h-10']">
              Get started
            </NuxtLink>
          </template>
        </div>

        <!-- Mobile menu button -->
        <button
          ref="toggleRef"
          type="button"
          class="inline-flex h-11 w-11 items-center justify-center rounded-xl border border-border text-muted transition-colors hover:text-text md:hidden"
          :aria-expanded="open"
          aria-controls="public-mobile-menu"
          :aria-label="open ? 'Close menu' : 'Open menu'"
          @click="open = !open"
        >
          <UiAppIcon :icon="open ? Cancel01Icon : Menu01Icon" :size="20" />
        </button>
      </div>
    </div>

    <!-- Mobile menu: overlays the page so content below does not jump -->
    <Transition
      enter-active-class="transition duration-150 ease-out motion-reduce:transition-none"
      enter-from-class="-translate-y-1 opacity-0"
      leave-active-class="transition duration-100 ease-in motion-reduce:transition-none"
      leave-to-class="-translate-y-1 opacity-0"
    >
      <div
        v-if="open"
        id="public-mobile-menu"
        ref="menuRef"
        class="absolute inset-x-0 top-full border-b border-border bg-surface px-4 pb-4 pt-2 shadow-lg md:hidden"
      >
        <nav class="grid gap-1" aria-label="Public mobile">
          <NuxtLink
            v-for="l in links"
            :key="l.to"
            :to="l.to"
            :aria-current="isActive(l.to) ? 'page' : undefined"
            class="flex h-12 items-center justify-between rounded-xl px-3 text-base font-medium transition-colors"
            :class="isActive(l.to) ? 'bg-primary/10 text-primary' : 'text-text hover:bg-bg'"
          >
            {{ l.label }}
            <span v-if="isActive(l.to)" class="h-2 w-2 rounded-full bg-primary" aria-hidden="true" />
          </NuxtLink>
        </nav>

        <div class="mt-3 grid gap-2 border-t border-border pt-3">
          <NuxtLink v-if="dashboardTo" :to="dashboardTo" :class="[primaryBtn, 'h-12 text-base']">
            Dashboard
          </NuxtLink>
          <template v-else>
            <NuxtLink v-if="allowSignups" to="/auth/signup" :class="[primaryBtn, 'h-12 text-base']">
              Get started
            </NuxtLink>
            <NuxtLink
              to="/auth/login"
              class="inline-flex h-12 items-center justify-center rounded-xl border border-border text-base font-semibold text-text transition-colors hover:bg-bg"
            >
              Log in
            </NuxtLink>
          </template>
        </div>
      </div>
    </Transition>
  </header>
</template>