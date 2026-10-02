<script setup lang="ts">
/**
 * Brand lockup. The logo ALWAYS renders through NuxtImg (local fallback
 * until backend settings provide the live logo_url).
 */
withDefaults(defineProps<{
  logoUrl?: string | null;
  size?: number;
  showName?: boolean;
  siteName?: string;
}>(), {
  logoUrl: null,
  size: 32,
  showName: true,
  siteName: 'Ctemz Journal',
})
</script>

<template>
  <span class="inline-flex items-center gap-2">
    <NuxtImg
      :src="logoUrl || '/logo.svg'"
      :width="size"
      :height="size"
      :alt="siteName"
      class="rounded-lg"
    />
    <span v-if="showName" class="text-lg font-bold tracking-tight">{{ siteName }}</span>
  </span>
</template>
