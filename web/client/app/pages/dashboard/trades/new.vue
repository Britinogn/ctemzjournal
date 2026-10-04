<script setup lang="ts">
import { toast } from 'vue-sonner'
import type { PendingImage, Trade, TradeCreate, UploadSignature } from '~/types'
import { NoteIcon } from '~/utils/icons' // or whatever icon you prefer for trades

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
      const uploaded = await $fetch<{ public_id: string }>(sig.upload_url, {
        method: 'POST',
        body: form,
      })
      await api.post(`/trades/${tradeId}/images`, {
        public_id: uploaded.public_id,
        kind: img.kind,
      })
      done += 1
      uploadStatus.value = `Uploading screenshots ${done}/${images.length}…`
    }
    catch {
      toast.warning(`Screenshot ${img.file.name} failed — trade saved without it`)
    }
  }
  uploadStatus.value = ''
}

async function onSubmit(
  input: TradeCreate,
  images: PendingImage[],
  makePublic: boolean,
): Promise<void> {
  if (saving.value)
    return
  saving.value = true
  try {
    const trade = await api.post<Trade>('/trades', { ...input })
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

/* ---------- Shared look ---------- */
const panel = 'rounded-2xl border border-border bg-surface p-4 md:p-5'
</script>

<template>
  <div>
    <!-- Mobile heading (desktop uses top bar) -->
    <h1 class="mb-4 text-xl font-bold tracking-tight md:sr-only">
      Log a trade
    </h1>

    <!-- A form reads best at a limited width. On a big screen it stays a comfortable column instead of stretching across. -->
    <section :class="[panel, 'mx-auto max-w-4xl']" aria-label="Log trade form">
      <div class="mb-5 flex items-center gap-3">
        <span class="inline-flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
          <UiAppIcon :icon="NoteIcon" :size="20" />
        </span>
        <div>
          <h2 class="text-sm font-semibold">
            New trade
          </h2>
          <p class="text-xs text-muted">
            Record the setup, risk and outcome
          </p>
        </div>
      </div>

      <TradesTradeForm
        :accounts="accounts ?? []"
        :setups="setups ?? []"
        :tags="tags ?? []"
        :saving="saving"
        submit-label="Save trade"
        @submit="onSubmit"
        @cancel="navigateTo('/dashboard/trades')"
      />

      <!-- Upload progress, with a spinner so it is clear the page is still working -->
      <p
        v-if="uploadStatus"
        role="status"
        class="mt-4 flex items-center justify-center gap-2 text-sm text-muted"
      >
        <svg class="h-4 w-4 animate-spin motion-reduce:animate-none" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" aria-hidden="true">
          <path d="M12 3a9 9 0 109 9" />
        </svg>
        {{ uploadStatus }}
      </p>
    </section>
  </div>
</template>