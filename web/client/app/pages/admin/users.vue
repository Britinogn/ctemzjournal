<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { toast } from 'vue-sonner'
import { adminUsersKey, meKey, type AdminUser } from '~/types'

definePageMeta({ middleware: 'admin', layout: 'admin' })

const api = useApi()
const queryClient = useQueryClient()
const search = ref('')
const debounced = ref('')
const page = ref(0)
const LIMIT = 20

let timer: ReturnType<typeof setTimeout> | null = null
watch(search, (v) => {
  if (timer)
    clearTimeout(timer)
  timer = setTimeout(() => {
    debounced.value = v.trim()
    page.value = 0
  }, 400)
})

const { data: me } = useQuery({
  queryKey: meKey(),
  queryFn: () => api.get<{ ID: string }>('/me'),
  staleTime: 5 * 60_000,
})

const { data: users, isPending, isError, refetch } = useQuery({
  queryKey: computed(() => adminUsersKey(debounced.value, page.value)),
  queryFn: () => {
    const q = new URLSearchParams({ limit: String(LIMIT), offset: String(page.value * LIMIT) })
    if (debounced.value)
      q.set('search', debounced.value)
    return api.get<AdminUser[]>(`/admin/users?${q.toString()}`)
  },
})

const acting = ref<string | null>(null)

async function setStatus(user: AdminUser, status: 'active' | 'suspended'): Promise<void> {
  if (acting.value)
    return
  acting.value = user.ID
  try {
    await api.patch(`/admin/users/${user.ID}/status`, { status })
    await refetch()
    toast.success(status === 'suspended' ? 'User suspended' : 'User reactivated')
  }
  catch {
    toast.error('Could not update user')
  }
  finally {
    acting.value = null
  }
}

function fmtDate(iso: string): string {
  return new Date(iso).toLocaleDateString('en-GB', { day: 'numeric', month: 'short', year: 'numeric', timeZone: 'UTC' })
}
</script>

<template>
  <div>
    <h1 class="mb-4 text-xl font-bold tracking-tight md:text-2xl">
      Users
    </h1>

    <div class="mb-4">
      <label for="admin-user-search" class="sr-only">Search users</label>
      <input
        id="admin-user-search"
        v-model="search"
        type="search"
        placeholder="Search by display name"
        class="w-full rounded-xl border border-border bg-surface px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary md:max-w-sm"
      >
    </div>

    <div v-if="isPending" class="h-96 animate-pulse rounded-2xl bg-surface" />
    <div v-else-if="isError" class="rounded-2xl border border-border bg-surface p-8 text-center">
      <p class="text-sm text-muted">
        Couldn't load users.
      </p>
      <button
        type="button"
        class="mt-3 rounded-xl bg-primary px-5 py-2 text-sm font-semibold text-on-primary transition hover:opacity-90"
        @click="() => refetch()"
      >
        Retry
      </button>
    </div>

    <section v-else class="overflow-hidden rounded-2xl border border-border bg-surface" aria-label="Users">
      <div class="overflow-x-auto">
        <table class="w-full min-w-[640px] text-left text-sm">
          <thead>
            <tr class="border-b border-border text-xs text-muted">
              <th class="px-4 py-2.5 font-medium">User</th>
              <th class="px-4 py-2.5 font-medium">Role</th>
              <th class="px-4 py-2.5 font-medium">Status</th>
              <th class="px-4 py-2.5 font-medium">Joined</th>
              <th class="px-4 py-2.5 text-right font-medium">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border">
            <tr v-for="u in users ?? []" :key="u.ID">
              <td class="px-4 py-2.5">
                <p class="font-medium">{{ u.DisplayName || '—' }}</p>
                <p class="tnum text-xs text-muted">{{ u.ID.slice(0, 8) }}</p>
              </td>
              <td class="px-4 py-2.5 capitalize text-muted">{{ u.Role }}</td>
              <td class="px-4 py-2.5">
                <span
                  :class="[
                    'rounded-full px-2.5 py-0.5 text-xs font-semibold',
                    u.Status === 'suspended' ? 'bg-warning/10 text-warning-text' : 'bg-profit/10 text-profit-text',
                  ]"
                >
                  {{ u.Status }}
                </span>
              </td>
              <td class="whitespace-nowrap px-4 py-2.5 text-muted">{{ fmtDate(u.CreatedAt) }}</td>
              <td class="px-4 py-2.5 text-right">
                <button
                  v-if="u.ID !== me?.ID"
                  type="button"
                  :disabled="acting === u.ID"
                  :class="[
                    'rounded-lg border px-3 py-1.5 text-xs font-semibold transition disabled:opacity-50',
                    u.Status === 'suspended'
                      ? 'border-border hover:border-profit hover:text-profit-text'
                      : 'border-border hover:border-warning hover:text-warning-text',
                  ]"
                  @click="setStatus(u, u.Status === 'suspended' ? 'active' : 'suspended')"
                >
                  {{ acting === u.ID ? 'Working…' : u.Status === 'suspended' ? 'Reactivate' : 'Suspend' }}
                </button>
                <span v-else class="text-xs text-muted">You</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="flex items-center justify-between border-t border-border px-4 py-3 text-sm">
        <p class="text-muted">
          Page {{ page + 1 }}
        </p>
        <div class="flex gap-2">
          <button
            type="button" :disabled="page === 0"
            class="rounded-xl border border-border px-4 py-1.5 text-sm font-medium transition hover:border-primary disabled:opacity-40"
            @click="page--"
          >
            Previous
          </button>
          <button
            type="button" :disabled="(users ?? []).length < LIMIT"
            class="rounded-xl border border-border px-4 py-1.5 text-sm font-medium transition hover:border-primary disabled:opacity-40"
            @click="page++"
          >
            Next
          </button>
        </div>
      </div>
    </section>
  </div>
</template>
