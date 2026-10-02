<script setup lang="ts">
import { toast } from 'vue-sonner'
import type { PendingImage, Trade, TradeCreate, UploadSignature } from '~/types'

definePageMeta({ middleware: 'auth', layout: 'dashboard' })

const api = useApi()
const { data: accounts } = useAccounts()
const { data: setups } = useSetups()
const { data: tags } = useTags()

const saving = ref(false)
const uploadStatus = ref('')

async function uploadImages(tradeId: string, images: PendingImage[]): Promise<void> {
  let done = 0
  for (const img of images) {
    try {
      const sig = await api.post<UploadSignature>(`/trades/${tradeId}/images/upload-signature`)
      const form = new FormData()
      form.append('file', img.file)
      form.append('api_key', sig.api_key)
      form.append('timestamp', sig.timestamp)
      form.append('signature', sig.signature)
      form.append('folder', sig.folder)
      form.append('type', sig.type)
      const uploaded = await $fetch<{ public_id: string }>(sig.upload_url, { method: 'POST', body: form })
      await api.post(`/trades/${tradeId}/images`, { public_id: uploaded.public_id, kind: img.kind })
      done += 1
      uploadStatus.value = `Uploading screenshots ${done}/${images.length}…`
    }
    catch {
      toast.warning(`Screenshot ${img.file.name} failed — trade saved without it`)
    }
  }
  uploadStatus.value = ''
}

async function onSubmit(input: TradeCreate, images: PendingImage[], makePublic: boolean): Promise<void> {
  if (saving.value)
    return
  saving.value = true
  try {
    const trade = await api.post<Trade>('/trades', input)
    if (makePublic) {
      await api.patch(`/trades/${trade.ID}/visibility`, { is_public: true })
    }
    if (images.length > 0)
      await uploadImages(trade.ID, images)
    toast.success('Trade logged')
    await navigateTo('/dashboard/trades')
  }
  catch {
    toast.error('Could not save the trade')
  }
  finally {
    saving.value = false
  }
}
</script>

<template>
  <div>
    <h1 class="mb-4 text-xl font-bold tracking-tight md:text-2xl">
      Log a trade
    </h1>
    <TradesTradeForm
      :accounts="accounts ?? []"
      :setups="setups ?? []"
      :tags="tags ?? []"
      :saving="saving"
      submit-label="Save trade"
      @submit="onSubmit"
      @cancel="navigateTo('/dashboard/trades')"
    />
    <p v-if="uploadStatus" role="status" class="mt-3 text-center text-sm text-muted">
      {{ uploadStatus }}
    </p>
  </div>
</template>
