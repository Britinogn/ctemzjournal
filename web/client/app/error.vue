<script setup lang="ts">
import { isChunkLoadError } from '~/utils/chunk-error'

const props = defineProps<{
  error?: { statusCode?: number; message?: string; stack?: string }
}>()

const nuxtError = useError()
const error = computed(() => props.error ?? nuxtError.value)

const isChunkError = computed(() => isChunkLoadError(
  error.value?.message ?? error.value,
))

const title = computed(() => {
  if (isChunkError.value)
    return 'Please refresh the page'
  if (error.value?.statusCode === 404)
    return 'Page not found'
  return 'Something went wrong'
})

const description = computed(() => {
  if (isChunkError.value)
    return 'We just released an update or your connection dropped. Refresh to get the latest version — your journal is safe.'
  if (error.value?.statusCode === 404)
    return 'This page moved or never existed. Head back home to keep journaling.'
  return 'An unexpected error happened. Try again, or go back home.'
})

function refresh(): void {
  window.location.reload()
}

async function goHome(): Promise<void> {
  await clearError({ redirect: '/' })
}
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center bg-bg px-4 py-10 text-text">
    <div class="w-full max-w-md rounded-2xl border border-border bg-surface p-7 text-center shadow-sm sm:p-8">
      <div class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 text-primary">
        <svg class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          <path d="M21 12a9 9 0 1 1-2.64-6.36M21 3v6h-6" />
        </svg>
      </div>
      <h1 class="text-2xl font-bold tracking-tight">
        {{ title }}
      </h1>
      <p class="mt-2 text-sm text-muted">
        {{ description }}
      </p>
      <div class="mt-6 flex flex-col gap-3">
        <button
          v-if="isChunkError"
          type="button"
          class="flex h-12 w-full items-center justify-center rounded-xl bg-primary text-sm font-semibold text-on-primary transition-colors hover:opacity-90"
          @click="refresh"
        >
          Refresh now
        </button>
        <button
          type="button"
          class="flex h-12 w-full items-center justify-center rounded-xl border border-border text-sm font-semibold transition hover:border-primary hover:text-primary"
          @click="goHome"
        >
          Back to home
        </button>
      </div>
    </div>
  </div>
</template>
