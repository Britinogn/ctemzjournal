<script setup lang="ts">
import { useWindowScroll } from '@vueuse/core'
import type { Profile } from '~/types'

defineProps<{ profile: Profile | null }>()

const { current, isChild } = useAdminNav()

// A soft shadow appears only once content has scrolled under the bar.
const { y } = useWindowScroll()
const scrolled = computed(() => y.value > 4)

const focusRing
  = 'focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary'
</script>

<template>
  <!--
    One bar for every screen size, and it works with both navigation surfaces:
    Tablet and desktop: a breadcrumb (Admin > page). The sidebar carries the logo and the links.
    Phone: logo and an Admin badge. The bottom tabs carry the links.
    Below a section (for example a user's detail page) a back button shows on every size.
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
      <template v-else>
        <NuxtLink
          to="/admin"
          :class="['-ml-1 flex min-h-11 min-w-0 items-center rounded-lg px-1 md:hidden', focusRing]"
        >
          <BrandLogo :size="28" />
        </NuxtLink>
        <span class="inline-flex shrink-0 items-center gap-1 rounded-full bg-warning/10 px-2.5 py-1 text-xs font-semibold text-warning-text md:hidden">
          <svg class="h-3 w-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M12 3l8 3v6c0 5-3.5 8-8 9-4.5-1-8-4-8-9V6l8-3z" />
          </svg>
          Admin
        </span>
      </template>

      <!-- Phone, below a section: the section name next to the back button -->
      <p v-if="isChild && current" class="min-w-0 truncate text-base font-bold tracking-tight md:hidden">
        {{ current.label }}
      </p>

      <!-- Tablet and desktop: a breadcrumb, not an h1, so each page keeps its own heading -->
      <nav v-if="current" class="hidden min-w-0 md:block" aria-label="Breadcrumb">
        <ol class="flex min-w-0 items-center gap-2">
          <li class="shrink-0">
            <NuxtLink to="/admin" :class="['rounded text-sm font-medium text-muted transition-colors hover:text-text', focusRing]">
              Admin
            </NuxtLink>
          </li>
          <li class="shrink-0 text-muted" aria-hidden="true">
            <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M9 5l7 7-7 7" />
            </svg>
          </li>
          <li class="min-w-0 truncate text-lg font-bold tracking-tight">
            <NuxtLink v-if="isChild" :to="current.to" :class="['rounded transition-colors hover:text-primary', focusRing]">
              {{ current.label }}
            </NuxtLink>
            <span v-else aria-current="page">{{ current.label }}</span>
          </li>
        </ol>
      </nav>

      <div class="ml-auto flex shrink-0 items-center gap-2">
        <ThemeToggle />
        <LayoutUserlayoutUserMenu :profile="profile" placement="down" />
      </div>
    </div>
  </header>
</template>