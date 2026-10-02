<script setup lang="ts">
import { toast } from 'vue-sonner'
import type { Tag, TagKind } from '~/types'

definePageMeta({ middleware: 'auth', layout: 'dashboard' })

const api = useApi()
const { data: tags, isPending, refetch } = useTags()

const name = ref('')
const kind = ref<TagKind>('general')
const adding = ref(false)

const toDelete = ref<Tag | null>(null)
const deleting = ref(false)

async function onAdd(): Promise<void> {
  if (!name.value.trim() || adding.value)
    return
  adding.value = true
  try {
    await api.post<Tag>('/tags', { name: name.value.trim(), kind: kind.value })
    name.value = ''
    await refetch()
    toast.success('Tag added')
  }
  catch {
    toast.error('Could not add tag — names must be unique')
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
</script>

<template>
  <div>
    <h1 class="mb-4 text-xl font-bold tracking-tight md:text-2xl">
      Tags
    </h1>
    <div class="grid items-start gap-4 xl:grid-cols-2">
      <section class="rounded-2xl border border-border bg-surface p-4" aria-label="Tags list">
        <div class="mb-2 flex items-center justify-between">
          <h2 class="text-sm font-semibold">Tags</h2>
          <span class="text-xs text-muted">{{ tags?.length ?? 0 }} total</span>
        </div>
        <div v-if="isPending" class="h-32 animate-pulse rounded-xl bg-bg" />
        <ul v-else-if="(tags ?? []).length > 0" class="divide-y divide-border">
          <li v-for="t in tags" :key="t.ID" class="flex items-center justify-between gap-2 py-3">
            <p class="truncate text-sm font-medium">{{ t.Name }}</p>
            <span class="flex shrink-0 items-center gap-2">
              <span
                :class="[
                  'rounded-full px-2.5 py-0.5 text-xs font-medium',
                  t.Kind === 'mistake' ? 'bg-warning/10 text-warning-text' : 'bg-bg text-muted',
                ]"
              >
                {{ t.Kind === 'mistake' ? '⚠ mistake' : 'general' }}
              </span>
              <button
                type="button"
                class="rounded-lg border border-border px-2.5 py-1 text-xs font-medium text-loss transition hover:border-loss"
                @click="toDelete = t"
              >
                Delete
              </button>
            </span>
          </li>
        </ul>
        <p v-else class="py-6 text-center text-sm text-muted">
          No tags yet — flag your first mistake.
        </p>
      </section>

      <section class="rounded-2xl border border-border bg-surface p-4" aria-label="Add tag">
        <h2 class="text-sm font-semibold">Add tag</h2>
        <form class="mt-3 space-y-3" @submit.prevent="onAdd">
          <div>
            <label for="tag-name" class="mb-1.5 block text-sm font-medium">Name</label>
            <input id="tag-name" v-model="name" type="text" required placeholder="Tag name" class="w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary">
          </div>
          <div>
            <label for="tag-kind" class="mb-1.5 block text-sm font-medium">Type</label>
            <select id="tag-kind" v-model="kind" class="w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition focus:border-primary">
              <option value="general">general</option>
              <option value="mistake">mistake</option>
            </select>
          </div>
          <button
            type="submit" :disabled="adding"
            class="w-full rounded-xl bg-primary py-2.5 text-sm font-semibold text-on-primary transition hover:opacity-90 disabled:opacity-50"
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
