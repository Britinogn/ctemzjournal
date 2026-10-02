<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { adminAuditKey, type AuditEntry } from '~/types'

definePageMeta({ middleware: 'admin', layout: 'admin' })

const api = useApi()
const page = ref(0)
const LIMIT = 20

const { data: rows, isPending, isError, refetch } = useQuery({
  queryKey: computed(() => adminAuditKey(page.value)),
  queryFn: () => api.get<AuditEntry[]>(`/admin/audit?limit=${LIMIT}&offset=${page.value * LIMIT}`),
})

function actionLabel(action: string): string {
  return action
    .split('.')
    .map(part => part.replace(/_/g, ' '))
    .map(part => part.charAt(0).toUpperCase() + part.slice(1))
    .join(' · ')
}

function fmtTime(iso: string): string {
  return new Date(iso).toLocaleString('en-GB', {
    day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit', timeZone: 'UTC',
  })
}

/** Meta arrives base64-encoded (jsonb bytes) — decode for display. */
function metaText(meta: string): string {
  if (!meta)
    return '—'
  try {
    const obj = JSON.parse(atob(meta)) as Record<string, unknown>
    const parts = Object.entries(obj).map(([k, v]) => `${k}: ${String(v)}`)
    return parts.length > 0 ? parts.join(', ') : '—'
  }
  catch {
    return '—'
  }
}
</script>

<template>
  <div>
    <h1 class="mb-1 text-xl font-bold tracking-tight md:text-2xl">
      Audit log
    </h1>
    <p class="mb-4 text-sm text-muted">
      Every admin write, newest first.
    </p>

    <div v-if="isPending" class="h-96 animate-pulse rounded-2xl bg-surface" />
    <div v-else-if="isError" class="rounded-2xl border border-border bg-surface p-8 text-center">
      <p class="text-sm text-muted">
        Couldn't load the audit log.
      </p>
      <button
        type="button"
        class="mt-3 rounded-xl bg-primary px-5 py-2 text-sm font-semibold text-on-primary transition hover:opacity-90"
        @click="() => refetch()"
      >
        Retry
      </button>
    </div>

    <section v-else class="overflow-hidden rounded-2xl border border-border bg-surface" aria-label="Audit log">
      <ul v-if="(rows ?? []).length > 0" class="divide-y divide-border">
        <li v-for="row in rows" :key="row.ID" class="px-4 py-3">
          <div class="flex flex-wrap items-center justify-between gap-2">
            <p class="text-sm font-semibold">
              {{ actionLabel(row.Action) }}
            </p>
            <p class="tnum text-xs text-muted">
              {{ fmtTime(row.CreatedAt) }}
            </p>
          </div>
          <p class="tnum mt-0.5 truncate text-xs text-muted">
            {{ row.TargetType }} · {{ row.TargetID.slice(0, 13) }} · {{ metaText(row.Meta) }}
          </p>
        </li>
      </ul>
      <p v-else class="p-8 text-center text-sm text-muted">
        No admin actions recorded yet.
      </p>
      <div class="flex items-center justify-end gap-2 border-t border-border px-4 py-3">
        <button
          type="button" :disabled="page === 0"
          class="rounded-xl border border-border px-4 py-1.5 text-sm font-medium transition hover:border-primary disabled:opacity-40"
          @click="page--"
        >
          Previous
        </button>
        <button
          type="button" :disabled="(rows ?? []).length < LIMIT"
          class="rounded-xl border border-border px-4 py-1.5 text-sm font-medium transition hover:border-primary disabled:opacity-40"
          @click="page++"
        >
          Next
        </button>
      </div>
    </section>
  </div>
</template>
