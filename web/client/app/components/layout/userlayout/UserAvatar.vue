<script setup lang="ts">
import { initials } from '~/utils/format';
import { resolveAvatarUrl } from '~/utils/avatar';

const props = withDefaults(defineProps<{
  name?: string | null;
  email?: string;
  role?: string;
  size?: number;
  /** Raw AvatarPath: storage path (`{userId}/avatar.webp`) or legacy http(s) URL. */
  src?: string | null;
  /** Cache-buster (profile UpdatedAt) — same path re-rendered fresh after re-upload. */
  rev?: string | null;
}>(), {
  name: null,
  email: '',
  role: 'user',
  size: 40,
  src: null,
  rev: null,
});

const config = useRuntimeConfig();
const label = computed(() => initials(props.name, props.email));
const hasImage = ref(true);

watch([() => props.src, () => props.rev], () => {
  hasImage.value = true;
});

const imageUrl = computed(() => resolveAvatarUrl(String(config.public.supabaseUrl || ''), props.src, props.rev));
const showImage = computed(() => imageUrl.value !== '' && hasImage.value);
</script>

<template>
  <img
    v-if="showImage"
    :src="imageUrl"
    :alt="name || email || 'Account'"
    :width="size"
    :height="size"
    loading="lazy"
    referrerpolicy="no-referrer"
    class="shrink-0 rounded-full object-cover"
    :style="{ width: `${size}px`, height: `${size}px` }"
    @error="hasImage = false"
  >
  <span
    v-else
    class="inline-flex shrink-0 items-center justify-center rounded-full bg-primary/15 font-semibold text-primary"
    :style="{ width: `${size}px`, height: `${size}px`, fontSize: `${Math.round(size * 0.36)}px` }"
    :title="name || email || 'Account'"
  >
    {{ label }}
  </span>
</template>
