<script setup lang="ts">
import { toast } from 'vue-sonner'
import type { Setup } from '~/types'

definePageMeta({ middleware: 'auth', layout: 'dashboard' })

const api = useApi()
const { data: setups, isPending, refetch } = useSetups()

const name = ref('')
const rules = ref('')
const invalidation = ref('')
const adding = ref(false)

const toDelete = ref<Setup | null>(null)
const deleting = ref(false)

async function onAdd(): Promise<void> {
  if (!name.value.trim() || adding.value)
    return
  adding.value = true
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
  }
  catch {
    toast.error('Could not add setup — names must be unique')
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
</script>

<template>
  <div>
    <h1 class="mb-4 text-xl font-bold tracking-tight md:text-2xl">
      Setups
    </h1>
    <div class="grid items-start gap-4 xl:grid-cols-2">
      <section class="rounded-2xl border border-border bg-surface p-4" aria-label="Setups list">
        <div class="mb-2 flex items-center justify-between">
          <h2 class="text-sm font-semibold">Setups</h2>
          <span class="text-xs text-muted">{{ setups?.length ?? 0 }} total</span>
        </div>
        <div v-if="isPending" class="h-32 animate-pulse rounded-xl bg-bg" />
        <ul v-else-if="(setups ?? []).length > 0" class="divide-y divide-border">
          <li v-for="s in setups" :key="s.ID" class="flex items-center justify-between gap-2 py-3">
            <div class="min-w-0">
              <p class="truncate text-sm font-medium">{{ s.Name }}</p>
              <p v-if="s.Rules" class="truncate text-xs text-muted">{{ s.Rules }}</p>
            </div>
            <button
              type="button"
              class="shrink-0 rounded-lg border border-border px-2.5 py-1 text-xs font-medium text-loss transition hover:border-loss"
              @click="toDelete = s"
            >
              Delete
            </button>
          </li>
        </ul>
        <p v-else class="py-6 text-center text-sm text-muted">
          No setups yet — define your first edge.
        </p>
      </section>

      <section class="rounded-2xl border border-border bg-surface p-4" aria-label="Add setup">
        <h2 class="text-sm font-semibold">Add setup</h2>
        <form class="mt-3 space-y-3" @submit.prevent="onAdd">
          <div>
            <label for="setup-name" class="mb-1.5 block text-sm font-medium">Name</label>
            <input id="setup-name" v-model="name" type="text" required placeholder="Setup name" class="w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary">
          </div>
          <div>
            <label for="setup-rules" class="mb-1.5 block text-sm font-medium">Rules <span class="font-normal text-muted">(optional)</span></label>
            <textarea id="setup-rules" v-model="rules" rows="2" placeholder="Entry rules…" class="w-full resize-y rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary" />
          </div>
          <div>
            <label for="setup-inv" class="mb-1.5 block text-sm font-medium">Invalidation <span class="font-normal text-muted">(optional)</span></label>
            <textarea id="setup-inv" v-model="invalidation" rows="2" placeholder="When the setup is void…" class="w-full resize-y rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary" />
          </div>
          <button
            type="submit" :disabled="adding"
            class="w-full rounded-xl bg-primary py-2.5 text-sm font-semibold text-on-primary transition hover:opacity-90 disabled:opacity-50"
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
