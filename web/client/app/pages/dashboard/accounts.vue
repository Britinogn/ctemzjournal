<script setup lang="ts">
import { toast } from 'vue-sonner'
import type { Account, AccountCurrency, AccountType } from '~/types'

definePageMeta({ middleware: 'auth', layout: 'dashboard' })

const api = useApi()
const { data: accounts, isPending, refetch } = useAccounts()

const name = ref('')
const type = ref<AccountType>('demo')
const currency = ref<AccountCurrency>('USD')
const balance = ref<number | undefined>(undefined)
const adding = ref(false)

const toDelete = ref<Account | null>(null)
const deleting = ref(false)

async function onAdd(): Promise<void> {
  if (!name.value.trim() || adding.value)
    return
  adding.value = true
  try {
    await api.post<Account>('/accounts', {
      name: name.value.trim(),
      type: type.value,
      currency: currency.value,
      starting_balance: balance.value ?? 0,
    })
    name.value = ''
    balance.value = undefined
    await refetch()
    toast.success('Account added')
  }
  catch {
    toast.error('Could not add account — names must be unique')
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
    await api.del(`/accounts/${toDelete.value.ID}`)
    toDelete.value = null
    await refetch()
    toast.success('Account deleted')
  }
  catch {
    toast.error('Could not delete account')
  }
  finally {
    deleting.value = false
  }
}
</script>

<template>
  <div>
    <h1 class="mb-4 text-xl font-bold tracking-tight md:text-2xl">
      Accounts
    </h1>
    <div class="grid items-start gap-4 xl:grid-cols-2">
      <section class="rounded-2xl border border-border bg-surface p-4" aria-label="Accounts list">
        <div class="mb-2 flex items-center justify-between">
          <h2 class="text-sm font-semibold">Accounts</h2>
          <span class="text-xs text-muted">{{ accounts?.length ?? 0 }} total</span>
        </div>
        <div v-if="isPending" class="h-32 animate-pulse rounded-xl bg-bg" />
        <ul v-else-if="(accounts ?? []).length > 0" class="divide-y divide-border">
          <li v-for="a in accounts" :key="a.ID" class="flex items-center justify-between gap-2 py-3">
            <p class="truncate text-sm font-medium">{{ a.Name }}</p>
            <span class="flex shrink-0 items-center gap-2">
              <span class="rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium lowercase text-primary">
                {{ a.Type }}
              </span>
              <span class="tnum hidden text-xs text-muted sm:inline">{{ a.Currency }}</span>
              <button
                type="button"
                class="rounded-lg border border-border px-2.5 py-1 text-xs font-medium text-loss transition hover:border-loss"
                @click="toDelete = a"
              >
                Delete
              </button>
            </span>
          </li>
        </ul>
        <p v-else class="py-6 text-center text-sm text-muted">
          No accounts yet — add your first one.
        </p>
      </section>

      <section class="rounded-2xl border border-border bg-surface p-4" aria-label="Add account">
        <h2 class="text-sm font-semibold">Add account</h2>
        <form class="mt-3 space-y-3" @submit.prevent="onAdd">
          <div>
            <label for="acc-name" class="mb-1.5 block text-sm font-medium">Name</label>
            <input id="acc-name" v-model="name" type="text" required placeholder="Account name" class="w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary">
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label for="acc-type" class="mb-1.5 block text-sm font-medium">Type</label>
              <select id="acc-type" v-model="type" class="w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition focus:border-primary">
                <option value="demo">demo</option>
                <option value="live">live</option>
              </select>
            </div>
            <div>
              <label for="acc-cur" class="mb-1.5 block text-sm font-medium">Currency</label>
              <select id="acc-cur" v-model="currency" class="w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition focus:border-primary">
                <option value="USD">USD</option>
                <option value="NGN">NGN</option>
              </select>
            </div>
          </div>
          <div>
            <label for="acc-bal" class="mb-1.5 block text-sm font-medium">Starting balance</label>
            <input id="acc-bal" v-model.number="balance" type="number" min="0" step="any" placeholder="0" class="tnum w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary">
          </div>
          <button
            type="submit" :disabled="adding"
            class="w-full rounded-xl bg-primary py-2.5 text-sm font-semibold text-on-primary transition hover:opacity-90 disabled:opacity-50"
          >
            {{ adding ? 'Adding…' : '+ Add account' }}
          </button>
        </form>
      </section>
    </div>

    <UiConfirmDialog
      :open="toDelete !== null"
      title="Delete account?"
      :message="`“${toDelete?.Name}” and its trades will be gone for good.`"
      :busy="deleting"
      @confirm="onDelete"
      @close="toDelete = null"
    />
  </div>
</template>
