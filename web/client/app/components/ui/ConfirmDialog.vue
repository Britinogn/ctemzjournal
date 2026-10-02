<script setup lang="ts">
defineProps<{
  open: boolean;
  title: string;
  message: string;
  confirmLabel?: string;
  busy?: boolean;
}>()

const emit = defineEmits<{
  confirm: [];
  close: [];
}>()
</script>

<template>
  <div v-if="open" class="fixed inset-0 z-50 flex items-end justify-center bg-text/40 p-4 sm:items-center" role="presentation" @click.self="emit('close')">
    <div role="alertdialog" aria-modal="true" :aria-label="title" class="w-full max-w-sm rounded-2xl border border-border bg-surface p-5 shadow-xl" @keydown.escape="emit('close')">
      <h2 class="text-base font-bold">
        {{ title }}
      </h2>
      <p class="mt-1 text-sm text-muted">
        {{ message }}
      </p>
      <div class="mt-4 flex justify-end gap-2">
        <button
          type="button"
          class="rounded-xl border border-border px-4 py-2 text-sm font-medium transition hover:border-primary"
          @click="emit('close')"
        >
          Cancel
        </button>
        <button
          type="button"
          :disabled="busy"
          class="rounded-xl bg-loss px-4 py-2 text-sm font-semibold text-white transition hover:opacity-90 disabled:opacity-50"
          @click="emit('confirm')"
        >
          {{ busy ? 'Working…' : (confirmLabel ?? 'Delete') }}
        </button>
      </div>
    </div>
  </div>
</template>
