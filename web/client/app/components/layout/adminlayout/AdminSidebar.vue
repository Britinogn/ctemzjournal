<script setup lang="ts">
import type { Profile } from '~/types'

defineProps<{ profile: Profile | null }>()

const { nav, isActive } = useAdminNav()

const focusRing
  = 'focus-visible:outline-2-2 focus-visible:outline-2-offset-2 focus-visible:outline-2-primary'
</script>

<template>
  <aside class="fixed inset-y-0 left-0 z-30 hidden w-20 flex-col border-r border-border bg-surface md:flex xl:w-64" aria-label="Admin sidebar">
    <!-- Same height as the top bar, with a bottom border, so the two lines join up -->
    <div class="flex h-16 items-center justify-center border-b border-border px-3 xl:justify-start xl:px-5">
      <NuxtLink to="/admin" class="flex items-center rounded-lg" :class="focusRing">
        <BrandLogo :show-name="false" :size="34" />
        <!-- max-xl:sr-only keeps the name for screen readers when the sidebar is the icon rail -->
        <span class="ml-2 text-lg font-bold tracking-tight max-xl:sr-only">Ctemz Journal</span>
      </NuxtLink>
    </div>

    <div class="mt-4 px-3 xl:px-4">
      <span class="flex items-center justify-center gap-1.5 rounded-xl bg-warning/10 px-3 py-2 text-sm font-semibold text-warning-text">
        <svg class="h-4 w-4 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          <path d="M12 3l8 3v6c0 5-3.5 8-8 9-4.5-1-8-4-8-9V6l8-3z" />
        </svg>
        <span class="max-xl:sr-only">Admin</span>
      </span>
    </div>

    <p class="mt-4 px-5 text-xs font-medium uppercase tracking-wide text-muted max-xl:hidden">
      Admin area
    </p>
    <nav class="mt-1 flex-1 space-y-1 overflow-y-auto px-3 xl:px-4" aria-label="Admin">
      <template v-for="(item, i) in nav" :key="item.to">
        <div v-if="i > 0 && nav[i - 1]!.group !== item.group" class="my-3! border-t border-border" aria-hidden="true" />
        <NuxtLink
          :to="item.to"
          :title="item.label"
          :aria-current="isActive(item) ? 'page' : undefined"
          :class="[
            'relative flex h-11 items-center justify-center gap-3 rounded-xl text-sm font-medium transition-colors xl:justify-start xl:px-3',
            focusRing,
            isActive(item)
              ? 'bg-primary/10 font-semibold text-primary'
              : 'text-muted hover:bg-bg hover:text-text',
          ]"
        >
          <!-- Edge marker, so the current page is not shown by colour alone -->
          <span
            v-if="isActive(item)"
            class="absolute inset-y-2 -left-3 w-1 rounded-r-full bg-primary xl:-left-4"
            aria-hidden="true"
          />
          <UiAppIcon :icon="item.icon" :size="20" class="shrink-0" />
          <span class="max-xl:sr-only">{{ item.label }}</span>
        </NuxtLink>
      </template>
    </nav>

    <!-- Settings is not passed to the card any more: the Site settings tab above already covers it. -->
    <div class="border-t border-border p-3 xl:border-t-0 xl:p-4">
      <LayoutUserlayoutUserCard :profile="profile" />
    </div>
  </aside>
</template>