<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import type { AdminOverview } from '~/types'
import { adminAuditKey, adminOverviewKey } from '~/types'
import type { AuditEntry } from '~/types'

definePageMeta({ middleware: 'admin', layout: 'admin' })

const api = useApi()

const { data: overview, isPending, isError, refetch } = useQuery({
  queryKey: adminOverviewKey(),
  queryFn: () => api.get<AdminOverview>('/admin/overview'),
})

const { data: recentAudit } = useQuery({
  queryKey: adminAuditKey(0),
  queryFn: () => api.get<AuditEntry[]>('/admin/audit?limit=5'),
  staleTime: 60_000,
})

const cards = computed(() => [
  { label: 'Total users', value: String(overview.value?.user_count ?? 0) },
  { label: 'Total trades', value: String(overview.value?.trade_count ?? 0) },
  { label: 'New users this week', value: String(overview.value?.new_users_this_week ?? 0) },
])

function actionLabel(action: string): string {
  return action
    .split('.')
    .map(part => part.replace(/_/g, ' '))
    .map(part => part.charAt(0).toUpperCase() + part.slice(1))
    .join(' · ')
}
</script>

<template>
  <div>
    <h1 class="mb-4 text-xl font-bold tracking-tight md:text-2xl">
      Admin overview
    </h1>

    <div v-if="isPending" class="grid grid-cols-1 gap-3 sm:grid-cols-3">
      <div v-for="i in 3" :key="i" class="h-28 animate-pulse rounded-2xl bg-surface" />
    </div>
    <div v-else-if="isError" class="rounded-2xl border border-border bg-surface p-8 text-center">
      <p class="text-sm text-muted">
        Couldn't load overview.
      </p>
      <button
        type="button"
        class="mt-3 rounded-xl bg-primary px-5 py-2 text-sm font-semibold text-on-primary transition hover:opacity-90"
        @click="() => refetch()"
      >
        Retry
      </button>
    </div>

    <template v-else>
      <div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
        <div v-for="card in cards" :key="card.label" class="rounded-2xl border border-border bg-surface p-4">
          <p class="text-xs font-medium text-muted">
            {{ card.label }}
          </p>
          <p class="tnum mt-1 text-2xl font-bold tracking-tight">
            {{ card.value }}
          </p>
        </div>
      </div>

      <div class="mt-4 grid gap-4 xl:grid-cols-2">
        <section class="rounded-2xl border border-border bg-surface p-4" aria-label="Recent admin activity">
          <div class="mb-2 flex items-center justify-between">
            <h2 class="text-sm font-semibold">
              Recent activity
            </h2>
            <NuxtLink to="/admin/audit" class="text-sm font-medium text-primary hover:underline">
              View all
            </NuxtLink>
          </div>
          <ul v-if="(recentAudit ?? []).length > 0" class="divide-y divide-border">
            <li v-for="row in recentAudit" :key="row.ID" class="flex items-center justify-between gap-2 py-2 text-sm">
              <span class="font-medium">{{ actionLabel(row.Action) }}</span>
              <span class="tnum shrink-0 text-xs text-muted">{{ row.TargetID.slice(0, 8) }}</span>
            </li>
          </ul>
          <p v-else class="py-4 text-center text-sm text-muted">
            No admin actions yet.
          </p>
        </section>

        <section class="rounded-2xl border border-border bg-surface p-4" aria-label="Quick links">
          <h2 class="mb-2 text-sm font-semibold">
            Manage
          </h2>
          <div class="grid grid-cols-2 gap-2">
            <NuxtLink
              v-for="link in [
                { to: '/admin/users', label: 'Users' },
                { to: '/admin/journals', label: 'Journals' },
                { to: '/admin/settings', label: 'Settings' },
                { to: '/admin/audit', label: 'Audit log' },
              ]"
              :key="link.to"
              :to="link.to"
              class="rounded-xl border border-border px-4 py-3 text-center text-sm font-semibold transition hover:border-primary hover:text-primary"
            >
              {{ link.label }}
            </NuxtLink>
          </div>
        </section>
      </div>
    </template>
  </div>
</template>
