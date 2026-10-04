<script setup lang="ts">
const { tabs, isActive } = useAdminNav()

const focusRing
  = 'focus-visible:outline-2-2 focus-visible:-outline-offset-2 focus-visible:outline-2-primary'
</script>

<template>
  <!-- Mobile bottom tabs (tablet and desktop use the admin sidebar). -->
  <nav class="fixed inset-x-0 bottom-0 z-30 border-t border-border bg-surface/95 pb-[env(safe-area-inset-bottom)] backdrop-blur md:hidden" aria-label="Admin">
    <ul class="grid grid-cols-5 px-2">
      <li v-for="tab in tabs" :key="tab.to">
        <NuxtLink
          :to="tab.to"
          :aria-current="isActive(tab) ? 'page' : undefined"
          :class="[
            'relative flex min-h-14 flex-col items-center justify-center gap-0.5 rounded-xl text-[11px] font-medium transition-colors',
            focusRing,
            isActive(tab) ? 'font-semibold text-primary' : 'text-muted',
          ]"
        >
          <!-- Marker on the top edge, so the current tab is not shown by colour alone -->
          <span v-if="isActive(tab)" class="absolute -top-px h-0.5 w-8 rounded-b-full bg-primary" aria-hidden="true" />
          <UiAppIcon :icon="tab.icon" :size="22" />
          {{ tab.short }}
        </NuxtLink>
      </li>
    </ul>
  </nav>
</template>