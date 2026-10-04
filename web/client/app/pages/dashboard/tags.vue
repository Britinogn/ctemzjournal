<script setup lang="ts">
import { toast } from 'vue-sonner'
import { Tag01Icon } from '~/utils/icons'
import type { Tag, TagKind } from '~/types'

definePageMeta({ middleware: 'auth', layout: 'dashboard' })

const api = useApi()
const { data: tags, isPending, refetch } = useTags()

const name = ref('')
const kind = ref<TagKind>('general')
const adding = ref(false)
const nameError = ref('')
const nameInput = ref<HTMLInputElement | null>(null)

const toDelete = ref<Tag | null>(null)
const deleting = ref(false)

async function onAdd(): Promise<void> {
  if (!name.value.trim() || adding.value)
    return
  adding.value = true
  nameError.value = ''
  try {
    await api.post<Tag>('/tags', { name: name.value.trim(), kind: kind.value })
    name.value = ''
    await refetch()
    toast.success('Tag added')
    nameInput.value?.focus() // ready for the next one
  }
  catch {
    nameError.value = 'Names must be unique'
    toast.error('Could not add tag — names must be unique')
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
    await api.del(`/tags/${toDelete.value.ID}`)
    toDelete.value = null
    await refetch()
    toast.success('Tag deleted')
  }
  catch {
    toast.error('Could not delete tag')
  }
  finally {
    deleting.value = false
  }
}

/* ---------- Shared look, defined once ---------- */
const panel = 'rounded-2xl border border-border bg-surface p-4 md:p-5'
const input
  = 'h-11 w-full rounded-xl border border-border bg-bg px-4 text-sm outline-none transition-colors placeholder:text-muted focus:border-primary focus-visible:ring-2 focus-visible:ring-primary/30'
const choice
  = 'flex h-10 cursor-pointer items-center justify-center gap-1.5 rounded-lg text-sm font-semibold capitalize text-muted transition-colors hover:text-text peer-checked:bg-surface peer-checked:shadow-sm peer-checked:ring-1 peer-checked:ring-border peer-focus-visible:outline-2-2 peer-focus-visible:outline-2-offset-2 peer-focus-visible:outline-2-primary'
</script>

<template>
  <div>
    <!-- The top bar shows "Tags" from tablet up; this heading is for phones and screen readers -->
    <h1 class="mb-4 text-xl font-bold tracking-tight md:sr-only">
      Tags
    </h1>

    <div class="grid items-start gap-4 xl:grid-cols-[minmax(0,1.3fr)_minmax(0,1fr)]">
      <!-- List -->
      <section :class="panel" aria-label="Tags list">
        <div class="mb-3 flex items-center justify-between">
          <h2 class="text-sm font-semibold">
            Tags
          </h2>
          <span class="tnum text-xs text-muted">{{ tags?.length ?? 0 }} total</span>
        </div>

        <div v-if="isPending" class="space-y-2" aria-hidden="true">
          <div v-for="i in 3" :key="i" class="h-16 animate-pulse rounded-xl border border-border bg-bg" />
        </div>

        <ul v-else-if="(tags ?? []).length > 0" class="space-y-2">
          <li
            v-for="t in tags"
            :key="t.ID"
            class="flex items-center gap-3 rounded-xl border border-border bg-bg px-3 py-2.5"
          >
            <!-- Mistake tags get the warning colour on the icon too, so they stand out down the list -->
            <span
              class="inline-flex h-10 w-10 shrink-0 items-center justify-center rounded-xl"
              :class="t.Kind === 'mistake' ? 'bg-warning/10 text-warning-text' : 'bg-primary/10 text-primary'"
            >
              <UiAppIcon :icon="Tag01Icon" :size="20" />
            </span>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-semibold">
                {{ t.Name }}
              </p>
              <!-- The word and the icon carry the meaning, not just the colour -->
              <p
                class="mt-0.5 inline-flex items-center gap-1 text-xs font-medium"
                :class="t.Kind === 'mistake' ? 'text-warning-text' : 'text-muted'"
              >
                <svg v-if="t.Kind === 'mistake'" class="h-3.5 w-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                  <path d="M12 4l9 16H3L12 4zM12 10v4M12 17h.01" />
                </svg>
                {{ t.Kind === 'mistake' ? 'mistake' : 'general' }}
              </p>
            </div>
            <button
              type="button"
              class="inline-flex h-10 shrink-0 items-center gap-1.5 rounded-xl border border-border bg-surface px-3 text-sm font-medium text-muted transition-colors hover:border-muted hover:text-text focus-visible:outline-2-2 focus-visible:outline-2-offset-2 focus-visible:outline-2-primary"
              :aria-label="`Delete ${t.Name}`"
              @click="toDelete = t"
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
            <UiAppIcon :icon="Tag01Icon" :size="24" />
          </span>
          <p class="mt-3 text-sm text-muted">
            No tags yet — flag your first mistake.
          </p>
          <a
            href="#tag-name"
            class="mt-3 inline-flex h-11 items-center rounded-xl border border-border bg-surface px-5 text-sm font-semibold transition-colors hover:border-primary xl:hidden"
            @click.prevent="nameInput?.focus()"
          >
            Go to the form
          </a>
        </div>
      </section>

      <!-- Form -->
      <section :class="panel" aria-label="Add tag">
        <h2 class="text-sm font-semibold">
          Add tag
        </h2>
        <form class="mt-4 space-y-4" @submit.prevent="onAdd">
          <div>
            <label for="tag-name" class="mb-1.5 block text-sm font-medium">Name</label>
            <input
              id="tag-name"
              ref="nameInput"
              v-model="name"
              type="text"
              required
              autocomplete="off"
              placeholder="Tag name"
              :class="[input, nameError ? 'border-loss' : '']"
              :aria-invalid="nameError ? 'true' : undefined"
              :aria-describedby="nameError ? 'tag-name-error' : undefined"
              @input="nameError = ''"
            >
            <!-- The backend's rule, shown on the field itself as well as in the toast -->
            <p v-if="nameError" id="tag-name-error" class="mt-1.5 flex items-center gap-1.5 text-xs font-medium text-loss" role="alert">
              <svg class="h-3.5 w-3.5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                <path d="M12 4l9 16H3L12 4zM12 10v4M12 17h.01" />
              </svg>
              {{ nameError }}
            </p>
          </div>

          <!-- Two options, so big tap targets beat a dropdown. Mistake picks up the warning colour. -->
          <fieldset>
            <legend class="mb-1.5 text-sm font-medium">
              Type
            </legend>
            <div class="grid grid-cols-2 gap-1 rounded-xl border border-border bg-bg p-1">
              <label>
                <input v-model="kind" type="radio" name="tag-kind" value="general" class="peer sr-only">
                <span :class="[choice, 'peer-checked:text-primary']">general</span>
              </label>
              <label>
                <input v-model="kind" type="radio" name="tag-kind" value="mistake" class="peer sr-only">
                <span :class="[choice, 'peer-checked:text-warning-text']">
                  <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                    <path d="M12 4l9 16H3L12 4zM12 10v4M12 17h.01" />
                  </svg>
                  mistake
                </span>
              </label>
            </div>
          </fieldset>

          <button
            type="submit"
            :disabled="adding"
            class="inline-flex h-12 w-full items-center justify-center rounded-xl bg-primary text-sm font-semibold text-on-primary transition-colors hover:opacity-90 focus-visible:outline-2-2 focus-visible:outline-2-offset-2 focus-visible:outline-2-primary disabled:cursor-wait disabled:opacity-60"
          >
            {{ adding ? 'Adding…' : '+ Add tag' }}
          </button>
        </form>
      </section>
    </div>

    <UiConfirmDialog
      :open="toDelete !== null"
      title="Delete tag?"
      :message="`“${toDelete?.Name}” will be removed from its trades.`"
      :busy="deleting"
      @confirm="onDelete"
      @close="toDelete = null"
    />
  </div>
</template>