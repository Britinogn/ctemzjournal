<script setup lang="ts">
import { toast } from 'vue-sonner'
import { Wallet01Icon } from '~/utils/icons'
import type { Account, AccountCurrency, AccountType } from '~/types'

definePageMeta({ middleware: 'auth', layout: 'dashboard' })

const api = useApi()
const { data: accounts, isPending, refetch } = useAccounts()

const name = ref('')
const type = ref<AccountType>('demo')
const currency = ref<AccountCurrency>('USD')
const balance = ref<number | undefined>(undefined)
const adding = ref(false)
const nameError = ref('')
const nameInput = ref<HTMLInputElement | null>(null)

const toDelete = ref<Account | null>(null)
const deleting = ref(false)

async function onAdd(): Promise<void> {
  if (!name.value.trim() || adding.value)
    return
  adding.value = true
  nameError.value = ''
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
    nameInput.value?.focus() // ready for the next one
  }
  catch {
    nameError.value = 'Names must be unique'
    toast.error('Could not add account — names must be unique')
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

/* ---------- Shared look, defined once ---------- */
const panel = 'rounded-2xl border border-border bg-surface p-4 md:p-5'
const input
  = 'h-11 w-full rounded-xl border border-border bg-bg px-4 text-sm outline-none transition-colors placeholder:text-muted focus:border-primary focus-visible:ring-2 focus-visible:ring-primary/30'
const choice
  = 'flex h-10 cursor-pointer items-center justify-center rounded-lg text-sm font-semibold text-muted transition-colors hover:text-text peer-checked:bg-surface peer-checked:text-primary peer-checked:shadow-sm peer-checked:ring-1 peer-checked:ring-border peer-focus-visible:outline-2 peer-focus-visible:outline-offset-2 peer-focus-visible:outline-primary'
const symbol = computed(() => (currency.value === 'NGN' ? '₦' : '$'))
</script>

<template>
  <div>
    <!-- The top bar shows "Accounts" from tablet up; this heading is for phones and screen readers -->
    <h1 class="mb-4 text-xl font-bold tracking-tight md:sr-only">
      Accounts
    </h1>

    <div class="grid items-start gap-4 xl:grid-cols-[minmax(0,1.3fr)_minmax(0,1fr)]">
      <!-- List -->
      <section :class="panel" aria-label="Accounts list">
        <div class="mb-3 flex items-center justify-between">
          <h2 class="text-sm font-semibold">
            Accounts
          </h2>
          <span class="tnum text-xs text-muted">{{ accounts?.length ?? 0 }} total</span>
        </div>

        <div v-if="isPending" class="space-y-2" aria-hidden="true">
          <div v-for="i in 3" :key="i" class="h-16 animate-pulse rounded-xl border border-border bg-bg" />
        </div>

        <ul v-else-if="(accounts ?? []).length > 0" class="space-y-2">
          <li
            v-for="a in accounts"
            :key="a.ID"
            class="flex items-center gap-3 rounded-xl border border-border bg-bg px-3 py-2.5"
          >
            <span class="inline-flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
              <UiAppIcon :icon="Wallet01Icon" :size="20" />
            </span>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-semibold">
                {{ a.Name }}
              </p>
              <!-- Currency used to disappear on phones, so it moved under the name where it always fits -->
              <p class="mt-0.5 flex items-center gap-2 text-xs text-muted">
                <span
                  class="rounded-full px-2 py-0.5 font-semibold capitalize"
                  :class="a.Type === 'live' ? 'bg-primary/10 text-primary' : 'border border-border bg-surface'"
                >
                  {{ a.Type }}
                </span>
                <span class="tnum">{{ a.Currency }}</span>
              </p>
            </div>
            <button
              type="button"
              class="inline-flex h-10 shrink-0 items-center gap-1.5 rounded-xl border border-border bg-surface px-3 text-sm font-medium text-muted transition-colors hover:border-muted hover:text-text focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
              :aria-label="`Delete ${a.Name}`"
              @click="toDelete = a"
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
            <UiAppIcon :icon="Wallet01Icon" :size="24" />
          </span>
          <p class="mt-3 text-sm text-muted">
            No accounts yet — add your first one.
          </p>
          <!-- Sends focus to the form, which on a phone sits below this list -->
          <a
            href="#acc-name"
            class="mt-3 inline-flex h-11 items-center rounded-xl border border-border bg-surface px-5 text-sm font-semibold transition-colors hover:border-primary xl:hidden"
            @click.prevent="nameInput?.focus()"
          >
            Go to the form
          </a>
        </div>
      </section>

      <!-- Form -->
      <section :class="panel" aria-label="Add account">
        <h2 class="text-sm font-semibold">
          Add account
        </h2>
        <form class="mt-4 space-y-4" @submit.prevent="onAdd">
          <div>
            <label for="acc-name" class="mb-1.5 block text-sm font-medium">Name</label>
            <input
              id="acc-name"
              ref="nameInput"
              v-model="name"
              type="text"
              required
              autocomplete="off"
              placeholder="Account name"
              :class="[input, nameError ? 'border-loss' : '']"
              :aria-invalid="nameError ? 'true' : undefined"
              :aria-describedby="nameError ? 'acc-name-error' : undefined"
              @input="nameError = ''"
            >
            <!-- The backend's rule, shown on the field itself as well as in the toast -->
            <p v-if="nameError" id="acc-name-error" class="mt-1.5 flex items-center gap-1.5 text-xs font-medium text-loss" role="alert">
              <svg class="h-3.5 w-3.5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                <path d="M12 4l9 16H3L12 4zM12 10v4M12 17h.01" />
              </svg>
              {{ nameError }}
            </p>
          </div>

          <!-- Two options each, so big tap targets beat a dropdown -->
          <div class="grid grid-cols-2 gap-3">
            <fieldset>
              <legend class="mb-1.5 text-sm font-medium">
                Type
              </legend>
              <div class="grid grid-cols-2 gap-1 rounded-xl border border-border bg-bg p-1">
                <label v-for="t in (['demo', 'live'] as const)" :key="t">
                  <input v-model="type" type="radio" name="acc-type" :value="t" class="peer sr-only">
                  <span :class="[choice, 'capitalize']">{{ t }}</span>
                </label>
              </div>
            </fieldset>
            <fieldset>
              <legend class="mb-1.5 text-sm font-medium">
                Currency
              </legend>
              <div class="grid grid-cols-2 gap-1 rounded-xl border border-border bg-bg p-1">
                <label v-for="c in (['USD', 'NGN'] as const)" :key="c">
                  <input v-model="currency" type="radio" name="acc-cur" :value="c" class="peer sr-only">
                  <span :class="choice">{{ c }}</span>
                </label>
              </div>
            </fieldset>
          </div>

          <div>
            <label for="acc-bal" class="mb-1.5 block text-sm font-medium">Starting balance</label>
            <div class="relative">
              <span class="tnum pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-sm text-muted" aria-hidden="true">{{ symbol }}</span>
              <input
                id="acc-bal"
                v-model.number="balance"
                type="number"
                inputmode="decimal"
                min="0"
                step="any"
                placeholder="0"
                :class="[input, 'tnum pl-9']"
              >
            </div>
          </div>

          <button
            type="submit"
            :disabled="adding"
            class="inline-flex h-12 w-full items-center justify-center rounded-xl bg-primary text-sm font-semibold text-on-primary transition-colors hover:opacity-90 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary disabled:cursor-wait disabled:opacity-60"
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