<script setup lang="ts">
import { toast } from 'vue-sonner'
import { Target01Icon } from '~/utils/icons'
import type { Setup } from '~/types'

definePageMeta({ middleware: 'auth', layout: 'dashboard' })

const api = useApi()
const { data: setups, isPending, refetch } = useSetups()

const name = ref('')
const rules = ref('')
const invalidation = ref('')
const adding = ref(false)
const nameError = ref('')
const nameInput = ref<HTMLInputElement | null>(null)

const toDelete = ref<Setup | null>(null)
const deleting = ref(false)

async function onAdd(): Promise<void> {
  if (!name.value.trim() || adding.value)
    return
  adding.value = true
  nameError.value = ''
  try {
    await api.post<Setup>('/setups', {
      name: name.value.trim(),
      ...(rules.value.trim() ? { rules: rules.value.trim() } : {}),
      ...(invalidation.value.trim() ? { invalidation: invalidation.value.trim() } : {}),
    })
    name.value = ''
    rules.value = ''
    invalidation.value = ''
    await refetch()
    toast.success('Setup added')
    nameInput.value?.focus() // ready for the next one
  }
  catch {
    nameError.value = 'Names must be unique'
    toast.error('Could not add setup — names must be unique')
    nameInput.value?.focus()
  }
  finally {
    adding.value = false
  }
}

async function onDelete(): Promise<void> {
  if (!toDelete.value || deleting.value)
    return
  deleting.value = true
  try {
    await api.del(`/setups/${toDelete.value.ID}`)
    toDelete.value = null
    await refetch()
    toast.success('Setup deleted')
  }
  catch {
    toast.error('Could not delete setup')
  }
  finally {
    deleting.value = false
  }
}

/* ---------- Shared look, defined once ---------- */
const panel = 'rounded-2xl border border-border bg-surface p-4 md:p-5'
const field
  = 'w-full rounded-xl border border-border bg-bg px-4 text-sm outline-none transition-colors placeholder:text-muted focus:border-primary focus-visible:ring-2 focus-visible:ring-primary/30'
</script>

<template>
  <div>
    <!-- The top bar shows "Setups" from tablet up; this heading is for phones and screen readers -->
    <h1 class="mb-4 text-xl font-bold tracking-tight md:sr-only">
      Setups
    </h1>

    <div class="grid items-start gap-4 xl:grid-cols-[minmax(0,1.3fr)_minmax(0,1fr)]">
      <!-- List -->
      <section :class="panel" aria-label="Setups list">
        <div class="mb-3 flex items-center justify-between">
          <h2 class="text-sm font-semibold">
            Setups
          </h2>
          <span class="tnum text-xs text-muted">{{ setups?.length ?? 0 }} total</span>
        </div>

        <div v-if="isPending" class="space-y-2" aria-hidden="true">
          <div v-for="i in 3" :key="i" class="h-16 animate-pulse rounded-xl border border-border bg-bg" />
        </div>

        <ul v-else-if="(setups ?? []).length > 0" class="space-y-2">
          <li
            v-for="s in setups"
            :key="s.ID"
            class="flex items-center gap-3 rounded-xl border border-border bg-bg px-3 py-2.5"
          >
            <span class="inline-flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
              <UiAppIcon :icon="Target01Icon" :size="20" />
            </span>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-semibold">
                {{ s.Name }}
              </p>
              <!-- Two lines of the rules, not one cut-off line -->
              <p v-if="s.Rules" class="mt-0.5 line-clamp-2 text-xs text-muted">
                {{ s.Rules }}
              </p>
            </div>
            <button
              type="button"
              class="inline-flex h-10 shrink-0 items-center gap-1.5 rounded-xl border border-loss/40 bg-surface px-3 text-sm font-medium text-loss transition-colors hover:border-loss hover:bg-loss/10 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-loss"
              :aria-label="`Delete ${s.Name}`"
              @click="toDelete = s"
            >
              <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                <path d="M4 7h16M9 7V4h6v3M6 7l1 13h10l1-13" />
              </svg>
              <span class="max-sm:sr-only">Delete</span>
            </button>
          </li>
        </ul>

        <div v-else class="flex flex-col items-center rounded-xl border border-dashed border-border px-4 py-10 text-center">
          <span class="inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-primary/10 text-primary">
            <UiAppIcon :icon="Target01Icon" :size="24" />
          </span>
          <p class="mt-3 text-sm text-muted">
            No setups yet — define your first edge.
          </p>
          <a
            href="#setup-name"
            class="mt-3 inline-flex h-11 items-center rounded-xl border border-border bg-surface px-5 text-sm font-semibold transition-colors hover:border-primary xl:hidden"
            @click.prevent="nameInput?.focus()"
          >
            Go to the form
          </a>
        </div>
      </section>

      <!-- Form -->
      <section :class="panel" aria-label="Add setup">
        <h2 class="text-sm font-semibold">
          Add setup
        </h2>
        <form class="mt-4 space-y-4" @submit.prevent="onAdd">
          <div>
            <label for="setup-name" class="mb-1.5 block text-sm font-medium">Name</label>
            <input
              id="setup-name"
              ref="nameInput"
              v-model="name"
              type="text"
              required
              autocomplete="off"
              placeholder="Setup name"
              :class="[field, 'h-11', nameError ? 'border-loss' : '']"
              :aria-invalid="nameError ? 'true' : undefined"
              :aria-describedby="nameError ? 'setup-name-error' : undefined"
              @input="nameError = ''"
            >
            <!-- The backend's rule, shown on the field itself as well as in the toast -->
            <p v-if="nameError" id="setup-name-error" class="mt-1.5 flex items-center gap-1.5 text-xs font-medium text-loss" role="alert">
              <svg class="h-3.5 w-3.5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                <path d="M12 4l9 16H3L12 4zM12 10v4M12 17h.01" />
              </svg>
              {{ nameError }}
            </p>
          </div>

          <div>
            <label for="setup-rules" class="mb-1.5 block text-sm font-medium">Rules <span class="font-normal text-muted">(optional)</span></label>
            <textarea id="setup-rules" v-model="rules" rows="3" placeholder="Entry rules…" :class="[field, 'min-h-24 resize-y py-2.5']" />
          </div>

          <div>
            <label for="setup-inv" class="mb-1.5 block text-sm font-medium">Invalidation <span class="font-normal text-muted">(optional)</span></label>
            <textarea id="setup-inv" v-model="invalidation" rows="3" placeholder="When the setup is void…" :class="[field, 'min-h-24 resize-y py-2.5']" />
          </div>

          <button
            type="submit"
            :disabled="adding"
            class="inline-flex h-12 w-full items-center justify-center rounded-xl bg-primary text-sm font-semibold text-on-primary transition-colors hover:opacity-90 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary disabled:cursor-wait disabled:opacity-60"
          >
            {{ adding ? 'Adding…' : '+ Add setup' }}
          </button>
        </form>
      </section>
    </div>

    <UiConfirmDialog
      :open="toDelete !== null"
      title="Delete setup?"
      :message="`“${toDelete?.Name}” will be unlinked from its trades.`"
      :busy="deleting"
      @confirm="onDelete"
      @close="toDelete = null"
    />
  </div>
</template>