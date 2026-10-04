<script setup lang="ts">
import { useEventListener } from '@vueuse/core'

const props = defineProps<{
  open: boolean
  title: string
  message: string
  confirmLabel?: string
  busy?: boolean
}>()

const emit = defineEmits<{
  confirm: []
  close: []
}>()

// Ids that tie the dialog to its title and message, so screen readers read both when it opens.
const uid = useId()
const titleId = `${uid}-title`
const descId = `${uid}-desc`

const panel = ref<HTMLElement | null>(null)
const cancelBtn = ref<HTMLButtonElement | null>(null)
let opener: HTMLElement | null = null

// Nothing can close the dialog while the action is running.
function requestClose(): void {
  if (!props.busy)
    emit('close')
}

// Escape closes. Tab stays inside the dialog instead of reaching the page behind it.
useEventListener('keydown', (e: KeyboardEvent) => {
  if (!props.open)
    return
  if (e.key === 'Escape') {
    e.stopPropagation()
    requestClose()
    return
  }
  if (e.key !== 'Tab' || !panel.value)
    return
  const items = [...panel.value.querySelectorAll<HTMLElement>('button:not([disabled]), [href], input, select, textarea, [tabindex]:not([tabindex="-1"])')]
  const first = items[0]
  const last = items[items.length - 1]
  if (!first || !last)
    return
  const active = document.activeElement
  if (!panel.value.contains(active) || (e.shiftKey && active === first)) {
    e.preventDefault()
    ;(e.shiftKey ? last : first).focus()
  }
  else if (!e.shiftKey && active === last) {
    e.preventDefault()
    first.focus()
  }
})

// On open: remember what was focused, lock the page behind, and focus Cancel (the safe choice).
// On close: give focus back to the button that opened it.
watch(() => props.open, async (isOpen) => {
  if (isOpen) {
    opener = document.activeElement as HTMLElement | null
    document.documentElement.style.overflow = 'hidden'
    await nextTick()
    cancelBtn.value?.focus()
  }
  else {
    document.documentElement.style.overflow = ''
    opener?.focus?.()
    opener = null
  }
})
onBeforeUnmount(() => {
  document.documentElement.style.overflow = ''
})
</script>

<template>
  <!-- Teleported to <body>, so a sticky header or a transformed parent can never trap it -->
  <ClientOnly>
    <Teleport to="body">
      <Transition name="dialog">
        <div
          v-if="open"
          class="fixed inset-0 z-50 flex items-end justify-center bg-text/40 p-4 pb-[max(1rem,env(safe-area-inset-bottom))] backdrop-blur-xs sm:items-center"
          role="presentation"
          @click.self="requestClose"
        >
          <div
            ref="panel"
            role="alertdialog"
            aria-modal="true"
            :aria-labelledby="titleId"
            :aria-describedby="descId"
            class="dialog-panel w-full max-w-sm rounded-2xl border border-border bg-surface p-5 shadow-xl"
          >
            <div class="flex items-start gap-3">
              <span class="inline-flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-loss/10 text-loss" aria-hidden="true">
                <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M12 4l9 16H3L12 4zM12 10v4M12 17h.01" />
                </svg>
              </span>
              <div class="min-w-0">
                <h2 :id="titleId" class="text-base font-bold">
                  {{ title }}
                </h2>
                <p :id="descId" class="mt-1 break-words text-sm text-muted">
                  {{ message }}
                </p>
              </div>
            </div>

            <!-- Side by side and full width on phones, right-aligned on larger screens -->
            <div class="mt-5 grid grid-cols-2 gap-2 sm:flex sm:justify-end">
              <button
                ref="cancelBtn"
                type="button"
                :disabled="busy"
                class="inline-flex h-11 items-center justify-center rounded-xl border border-border px-5 text-sm font-semibold transition-colors hover:border-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary disabled:cursor-not-allowed disabled:opacity-50"
                @click="requestClose"
              >
                Cancel
              </button>
              <!-- on-primary is white in light mode and dark in dark mode, so the text stays readable on both reds -->
              <button
                type="button"
                :disabled="busy"
                :aria-busy="busy"
                class="inline-flex h-11 items-center justify-center rounded-xl bg-loss px-5 text-sm font-semibold text-on-primary transition-colors hover:opacity-90 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-loss disabled:cursor-wait disabled:opacity-60"
                @click="emit('confirm')"
              >
                {{ busy ? 'Working…' : (confirmLabel ?? 'Delete') }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </ClientOnly>
</template>

<style scoped>
.dialog-enter-active,
.dialog-leave-active {
  transition: opacity 0.16s ease;
}
.dialog-enter-active .dialog-panel,
.dialog-leave-active .dialog-panel {
  transition: transform 0.2s cubic-bezier(0.22, 0.8, 0.3, 1);
}
.dialog-enter-from,
.dialog-leave-to {
  opacity: 0;
}
.dialog-enter-from .dialog-panel,
.dialog-leave-to .dialog-panel {
  transform: translateY(16px) scale(0.98);
}
@media (prefers-reduced-motion: reduce) {
  .dialog-enter-active,
  .dialog-leave-active,
  .dialog-enter-active .dialog-panel,
  .dialog-leave-active .dialog-panel {
    transition: none;
  }
}
</style>