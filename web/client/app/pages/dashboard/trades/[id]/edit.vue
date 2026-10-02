<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { toast } from 'vue-sonner'
import type { PendingImage, Trade, TradeCreate, TradeImageWithUrl, TradeUpdate } from '~/types'

definePageMeta({ middleware: 'auth', layout: 'dashboard' })

const route = useRoute()
const tradeId = computed(() => String(route.params.id ?? ''))
const api = useApi()
const queryClient = useQueryClient()
const { data: accounts } = useAccounts()
const { data: setups } = useSetups()
const { data: tags } = useTags()

const { data: trade, isPending, isError } = useQuery({
  queryKey: ['trade', tradeId.value],
  queryFn: () => api.get<Trade>(`/trades/${tradeId.value}`),
})
const { data: images, refetch: refetchImages } = useQuery({
  queryKey: ['trade-images', tradeId.value],
  queryFn: () => api.get<TradeImageWithUrl[]>(`/trades/${tradeId.value}/images`),
})

/** Backend ISO → datetime-local value. */
function toLocalInput(iso: string | null): string {
  if (!iso)
    return ''
  const d = new Date(iso)
  const pad = (n: number): string => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const initial = computed(() => {
  const t = trade.value
  if (!t)
    return undefined
  return {
    account_id: t.AccountID,
    setup_id: t.SetupID ?? '',
    pair: t.Pair,
    direction: t.Direction,
    timeframe: t.Timeframe ?? 'H1',
    opened_at: toLocalInput(t.OpenedAt),
    entry: t.Entry ?? undefined,
    stop_loss: t.StopLoss ?? undefined,
    take_profit: t.TakeProfit ?? undefined,
    exit_price: t.ExitPrice ?? undefined,
    lot_size: t.LotSize ?? undefined,
    commission: t.Commission,
    swap: t.Swap,
    followed_rules: t.FollowedRules ?? undefined,
    emotion: t.Emotion ?? '',
    notes: t.Notes ?? '',
    status: t.Status,
    tag_ids: [] as string[],
    make_public: false,
  }
})

const saving = ref(false)
const uploadStatus = ref('')
const deletingImage = ref<string | null>(null)

async function uploadImages(id: string, pending: PendingImage[]): Promise<void> {
  let done = 0
  for (const img of pending) {
    try {
      const sig = await api.post<{ upload_url: string; api_key: string; timestamp: string; signature: string; folder: string; type: string }>(`/trades/${id}/images/upload-signature`)
      const form = new FormData()
      form.append('file', img.file)
      form.append('api_key', sig.api_key)
      form.append('timestamp', sig.timestamp)
      form.append('signature', sig.signature)
      form.append('folder', sig.folder)
      form.append('type', sig.type)
      const uploaded = await $fetch<{ public_id: string }>(sig.upload_url, { method: 'POST', body: form })
      await api.post(`/trades/${id}/images`, { public_id: uploaded.public_id, kind: img.kind })
      done += 1
      uploadStatus.value = `Uploading screenshots ${done}/${pending.length}…`
    }
    catch {
      toast.warning(`Screenshot ${img.file.name} failed`)
    }
  }
  uploadStatus.value = ''
}

async function onSubmit(input: TradeCreate, pending: PendingImage[]): Promise<void> {
  if (saving.value)
    return
  saving.value = true
  try {
    const body: TradeUpdate = { ...input }
    await api.patch<Trade>(`/trades/${tradeId.value}`, body)
    if (pending.length > 0)
      await uploadImages(tradeId.value, pending)
    queryClient.invalidateQueries({ queryKey: ['trade', tradeId.value] })
    queryClient.invalidateQueries({ queryKey: ['trade-images', tradeId.value] })
    queryClient.invalidateQueries({ queryKey: ['trades'] })
    queryClient.invalidateQueries({ queryKey: ['dashboard'] })
    toast.success('Trade updated')
    await navigateTo(`/dashboard/trades/${tradeId.value}`)
  }
  catch {
    toast.error('Could not update the trade')
  }
  finally {
    saving.value = false
  }
}

async function onDeleteImage(imageId: string): Promise<void> {
  if (deletingImage.value)
    return
  deletingImage.value = imageId
  try {
    await api.del(`/trades/${tradeId.value}/images/${imageId}`)
    await refetchImages()
    toast.success('Screenshot deleted')
  }
  catch {
    toast.error('Could not delete screenshot')
  }
  finally {
    deletingImage.value = null
  }
}
</script>

<template>
  <div>
    <div class="mb-4 flex items-center justify-between">
      <NuxtLink :to="`/dashboard/trades/${tradeId}`" class="text-sm font-medium text-primary hover:underline">
        ← Back to trade
      </NuxtLink>
    </div>
    <h1 class="mb-4 text-xl font-bold tracking-tight md:text-2xl">
      Edit trade
    </h1>

    <div v-if="isPending" class="h-96 animate-pulse rounded-2xl bg-surface" />
    <div v-else-if="isError || !trade" class="rounded-2xl border border-border bg-surface p-8 text-center">
      <p class="text-sm text-muted">
        Trade not found.
      </p>
      <NuxtLink to="/dashboard/trades" class="mt-3 inline-block text-sm font-medium text-primary hover:underline">
        Back to trades
      </NuxtLink>
    </div>

    <template v-else>
      <!-- Existing screenshots -->
      <section class="mb-4 rounded-2xl border border-border bg-surface p-4 md:p-5" aria-label="Existing screenshots">
        <h2 class="text-sm font-semibold">Screenshots on this trade</h2>
        <div v-if="(images ?? []).length > 0" class="mt-3 grid grid-cols-2 gap-2 sm:grid-cols-4">
          <div v-for="img in images" :key="img.ID" class="overflow-hidden rounded-xl border border-border">
            <img :src="img.URL" :alt="`${img.Kind ?? 'Trade'} screenshot`" class="aspect-video w-full object-cover" loading="lazy">
            <div class="flex items-center justify-between px-2 py-1">
              <span class="text-[11px] capitalize text-muted">{{ img.Kind ?? 'chart' }}</span>
              <button
                type="button"
                :disabled="deletingImage === img.ID"
                class="text-[11px] font-medium text-loss transition hover:underline disabled:opacity-50"
                @click="onDeleteImage(img.ID)"
              >
                {{ deletingImage === img.ID ? 'Removing…' : 'Remove' }}
              </button>
            </div>
          </div>
        </div>
        <p v-else class="mt-2 text-sm text-muted">
          No screenshots yet — add some below.
        </p>
        <p class="mt-2 text-[11px] text-muted">Clearing the exit price reopens the trade (results are cleared server-side).</p>
      </section>

      <TradesTradeForm
        :accounts="accounts ?? []"
        :setups="setups ?? []"
        :tags="tags ?? []"
        :initial="initial"
        :saving="saving"
        submit-label="Save changes"
        @submit="onSubmit"
        @cancel="navigateTo(`/dashboard/trades/${tradeId}`)"
      />
      <p v-if="uploadStatus" role="status" class="mt-3 text-center text-sm text-muted">
        {{ uploadStatus }}
      </p>
    </template>
  </div>
</template>
