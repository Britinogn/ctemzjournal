<script setup lang="ts">
import { useWindowScroll } from '@vueuse/core'
import type { Profile } from '~/types'

defineProps<{ profile: Profile | null }>()

const { current, isChild } = useUserNav()

// A soft shadow appears only once content has scrolled under the bar.
const { y } = useWindowScroll()
const scrolled = computed(() => y.value > 4)

const focusRing
  = 'focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary'
</script>

<template>
  <!--
    One bar for every screen size, and it works with both navigation surfaces:
    Tablet and desktop: the page title. The sidebar carries the logo and the links.
    Phone: the logo. The bottom tabs carry the links.
    Below a section (new trade, trade detail) a back button shows on every size.
    The theme switch lives in the account menu, so there is no second toggle here.
    Its inner box has the same max width and padding as <main>, so the left edge lines up with the page.
    On tablet and desktop it is as tall as the sidebar's logo row (h-16), so the two bottom borders form one line.
  -->
  <header
    class="sticky top-0 z-30 border-b border-border bg-surface/95 pt-[env(safe-area-inset-top)] backdrop-blur transition-shadow"
    :class="scrolled ? 'shadow-sm' : ''"
  >
    <div class="mx-auto flex h-14 w-full max-w-6xl items-center gap-3 px-4 md:h-16 md:px-6">
      <NuxtLink
        v-if="isChild && current"
        :to="current.to"
        :aria-label="`Back to ${current.label}`"
        :class="['-ml-2 inline-flex h-11 w-11 shrink-0 items-center justify-center rounded-xl text-muted transition-colors hover:bg-bg hover:text-text', focusRing]"
      >
        <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          <path d="M15 5l-7 7 7 7" />
        </svg>
      </NuxtLink>
      <!-- Phone only: the sidebar has the logo from tablet up -->
      <NuxtLink
        v-else
        to="/dashboard"
        :class="['-ml-1 flex min-h-11 min-w-0 items-center rounded-lg px-1 md:hidden', focusRing]"
      >
        <BrandLogo :size="28" />
      </NuxtLink>

      <!-- Not an h1: each page keeps its own heading -->
      <p
        v-if="current"
        class="min-w-0 truncate text-lg font-bold tracking-tight"
        :class="isChild ? '' : 'hidden md:block'"
      >
        {{ current.label }}
      </p>

      <div class="ml-auto flex shrink-0 items-center">
        <LayoutUserlayoutUserMenu :profile="profile" placement="down" />
      </div>
    </div>
  </header>
</template>