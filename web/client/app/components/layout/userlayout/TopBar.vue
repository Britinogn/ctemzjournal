<script setup lang="ts">
import { useWindowScroll } from '@vueuse/core'
import type { Profile } from '~/types'
import { InternetIcon } from '~/utils/icons'
import BrandLogo from '../BrandLogo.vue'

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
    Phone: a sticky bar across the top, as before.
    Tablet and desktop: the bar is fixed to the top of the screen, in the space to the right of the sidebar.
    It starts where the sidebar ends (md:left-20 for the icon rail, xl:left-64 for the full sidebar) and runs to the right edge,
    so it sits level with the sidebar's logo row (both are h-16) and the two bottom borders form one line.
    Those two numbers must match the sidebar's widths (w-20 and xl:w-64).
  -->
  <header
    class="sticky top-0 z-30 border-b border-border bg-surface/95 pt-[env(safe-area-inset-top)] backdrop-blur transition-shadow md:fixed md:left-20 md:right-0 md:top-0 xl:left-64"
    :class="scrolled ? 'shadow-sm' : ''"
  >
    <!-- Same max width and padding as <main>, so the left edge lines up with the page -->
    <div class="mx-auto flex h-14 w-full max-w-6xl items-center gap-3 px-4 md:h-16 md:px-6">
      <!-- Phone only: the sidebar has the logo from tablet up. Remove md:hidden to show it here as well. -->
      <NuxtLink
        to="/dashboard"
        :class="['-ml-1 flex min-h-11 min-w-0 items-center rounded-lg px-1 md:hidden', focusRing]"
      >
        <BrandLogo :size="150" />
      </NuxtLink>

      <!-- Not an h1: each page keeps its own heading -->
      <p
        v-if="current"
        class="min-w-0 truncate text-lg font-bold tracking-tight"
        :class="isChild ? '' : 'hidden md:block'"
      >
        {{ current.label }}
      </p>

      <div class="ml-auto flex shrink-0 items-center gap-1">
        <NuxtLink
          to="/"
          aria-label="Home"
          title="Home"
          :class="['inline-flex h-11 w-11 items-center justify-center rounded-xl text-muted transition-colors hover:bg-bg hover:text-text', focusRing]"
        >
          <UiAppIcon :icon="InternetIcon" :size="22" aria-hidden="true" />
        </NuxtLink>
        <LayoutUserlayoutUserMenu :profile="profile" placement="down" />
      </div>
    </div>
  </header>

  <!--
    A fixed bar takes no space in the page, so this empty block holds its place and the content starts below it.
    It is not needed on a phone, where the bar is sticky and does take space.
  -->
  <div class="hidden md:block md:h-[calc(4rem+env(safe-area-inset-top))]" aria-hidden="true" />
</template>