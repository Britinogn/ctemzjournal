<script setup lang="ts">
import imageCompression from 'browser-image-compression'
import type { Account, PendingImage, Setup, Tag, TradeCreate, TradeDraft } from '~/types'
import { fromDatetimeLocalToISO } from '~/utils/format'
import { detectBrowserTimezone } from '~/utils/timezones'

const props = withDefaults(defineProps<{
  accounts: Account[];
  setups: Setup[];
  tags: Tag[];
  initial?: Partial<TradeDraft>;
  saving?: boolean;
  submitLabel?: string;
  /** Profile timezone — typed wall times are interpreted in this zone. */
  timezone?: string;
}>(), {
  saving: false,
  submitLabel: 'Save trade',
  timezone: '',
})

const tz = computed(() => props.timezone || detectBrowserTimezone())

const emit = defineEmits<{
  submit: [input: TradeCreate, images: PendingImage[], makePublic: boolean];
  cancel: [];
}>()

// const PAIRS = ['EUR/USD', 'GBP/USD', 'USD/JPY', 'USD/CHF', 'AUD/USD', 'USD/CAD', 'XAU/USD', 'USD/NGN']
const PAIRS = [
  'EUR/USD', 'GBP/USD', 'USD/JPY', 'USD/CHF',
  'AUD/USD', 'USD/CAD', 'NZD/USD',

  'EUR/GBP', 'EUR/JPY', 'EUR/CHF', 'EUR/AUD', 'EUR/CAD', 'EUR/NZD',

  'GBP/JPY', 'GBP/CHF', 'GBP/AUD', 'GBP/CAD', 'GBP/NZD',

  'AUD/JPY', 'AUD/CHF', 'AUD/CAD', 'AUD/NZD',

  'NZD/JPY', 'NZD/CHF', 'NZD/CAD',

  'CAD/JPY', 'CAD/CHF', 'CHF/JPY',

  'XAU/USD', 'USD/NGN',

  'BTC/USD', 'ETH/USD', 'SOL/USD',
]
const TIMEFRAMES = ['M1', 'M5', 'M15', 'M30', 'H1', 'H4', 'D1', 'W1']
const EMOTIONS = ['Calm', 'Confident', 'Neutral', 'Anxious', 'FOMO', 'Frustrated']
const MAX_IMAGES = 5

const draft = reactive<TradeDraft>({
  account_id: props.initial?.account_id ?? '',
  setup_id: props.initial?.setup_id ?? '',
  pair: props.initial?.pair ?? '',
  direction: props.initial?.direction ?? 'long',
  timeframe: props.initial?.timeframe ?? 'H1',
  opened_at: props.initial?.opened_at ?? '',
  closed_at: props.initial?.closed_at ?? '',
  entry: props.initial?.entry,
  stop_loss: props.initial?.stop_loss,
  take_profit: props.initial?.take_profit,
  exit_price: props.initial?.exit_price,
  lot_size: props.initial?.lot_size,
  commission: props.initial?.commission,
  swap: props.initial?.swap,
  followed_rules: props.initial?.followed_rules,
  emotion: props.initial?.emotion ?? '',
  notes: props.initial?.notes ?? '',
  status: props.initial?.status ?? 'open',
  tag_ids: props.initial?.tag_ids ?? [],
  make_public: props.initial?.make_public ?? false,
})

const pending = ref<PendingImage[]>([])
const fileError = ref<string | null>(null)

function pipSize(pair: string): number {
  const p = pair.toUpperCase()
  if (p.includes('JPY') || p.includes('XAU'))
    return 0.01
  return 0.0001
}

function contractSize(pair: string): number {
  return pair.toUpperCase().includes('XAU') ? 100 : 100_000
}

/** Display-only preview (the server computes the real values on save). */
const preview = computed(() => {
  const none = { riskPips: null as number | null, rewardPips: null as number | null, estRisk: null as number | null, rr: null as number | null }
  const { pair } = draft
  const entry = numOrUndef(draft.entry)
  const stopLoss = numOrUndef(draft.stop_loss)
  const lotSize = numOrUndef(draft.lot_size)
  if (entry === undefined || stopLoss === undefined || lotSize === undefined || !pair)
    return none
  const pip = pipSize(pair)
  const contract = contractSize(pair)
  const takeProfit = numOrUndef(draft.take_profit)
  const riskPips = Math.abs(entry - stopLoss) / pip
  const rewardPips = takeProfit === undefined ? null : Math.abs(takeProfit - entry) / pip
  const estRisk = riskPips * pip * lotSize * contract
  const rr = rewardPips !== null && riskPips > 0 ? rewardPips / riskPips : null
  return { riskPips, rewardPips, estRisk, rr }
})

function toggleTag(id: string): void {
  const i = draft.tag_ids.indexOf(id)
  if (i >= 0)
    draft.tag_ids.splice(i, 1)
  else
    draft.tag_ids.push(id)
}

async function onFiles(event: Event): Promise<void> {
  fileError.value = null
  const input = event.target as HTMLInputElement
  const files = [...(input.files ?? [])].slice(0, MAX_IMAGES - pending.value.length)
  input.value = ''
  for (const file of files) {
    try {
      const compressed = await imageCompression(file, {
        maxSizeMB: 2,
        maxWidthOrHeight: 1600,
        fileType: 'image/webp',
      })
      const webp = new File([compressed], file.name.replace(/\.\w+$/, '') + '.webp', { type: 'image/webp' })
      pending.value.push({
        preview: URL.createObjectURL(webp),
        file: webp,
        kind: 'entry',
        progress: 0,
      })
    }
    catch {
      fileError.value = 'Could not compress an image — skipped.'
    }
  }
  if (pending.value.length >= MAX_IMAGES)
    fileError.value = 'Maximum 5 images per trade.'
}

function removePending(index: number): void {
  const [removed] = pending.value.splice(index, 1)
  if (removed)
    URL.revokeObjectURL(removed.preview)
}

const valid = computed(() =>
  draft.account_id !== '' && draft.pair.trim() !== ''
  && isNum(draft.entry) && isNum(draft.stop_loss) && isNum(draft.lot_size),
)

function isNum(v: unknown): v is number {
  return typeof v === 'number' && !Number.isNaN(v)
}

function numOrUndef(v: unknown): number | undefined {
  return isNum(v) ? v : undefined
}

function onSubmit(): void {
  if (!valid.value || props.saving)
    return
  const input: TradeCreate = {
    account_id: draft.account_id,
    pair: draft.pair.trim().toUpperCase(),
    direction: draft.direction,
    status: numOrUndef(draft.exit_price) !== undefined ? 'closed' : 'open',
  }
  if (draft.setup_id)
    input.setup_id = draft.setup_id
  if (draft.timeframe)
    input.timeframe = draft.timeframe
  if (draft.opened_at)
    input.opened_at = fromDatetimeLocalToISO(draft.opened_at, tz.value)
  const entry = numOrUndef(draft.entry)
  const stopLoss = numOrUndef(draft.stop_loss)
  const takeProfit = numOrUndef(draft.take_profit)
  const exitPrice = numOrUndef(draft.exit_price)
  const lotSize = numOrUndef(draft.lot_size)
  const commission = numOrUndef(draft.commission)
  const swap = numOrUndef(draft.swap)
  if (entry !== undefined)
    input.entry = entry
  if (stopLoss !== undefined)
    input.stop_loss = stopLoss
  if (takeProfit !== undefined)
    input.take_profit = takeProfit
  if (exitPrice !== undefined)
    input.exit_price = exitPrice
  if (lotSize !== undefined)
    input.lot_size = lotSize
  if (commission !== undefined)
    input.commission = commission
  if (swap !== undefined)
    input.swap = swap
  if (draft.followed_rules !== undefined)
    input.followed_rules = draft.followed_rules
  if (draft.emotion)
    input.emotion = draft.emotion
  if (draft.notes.trim())
    input.notes = draft.notes.trim()
  if (draft.tag_ids.length > 0)
    input.tag_ids = [...draft.tag_ids]
  emit('submit', input, [...pending.value], draft.make_public)
}

const inputCls = 'w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary'
const labelCls = 'mb-1.5 block text-sm font-medium'
</script>

<template>
  <form @submit.prevent="onSubmit">
    <div class="grid gap-4 xl:grid-cols-3">
      <div class="space-y-4 xl:col-span-2">
        <!-- The trade -->
        <section class="rounded-2xl border border-border bg-surface p-4 md:p-5" aria-label="The trade">
          <h2 class="text-sm font-semibold">The trade</h2>
          <div class="mt-3 grid gap-3 sm:grid-cols-2">
            <div>
              <label for="tf-account" :class="labelCls">Account</label>
              <select id="tf-account" v-model="draft.account_id" required :class="inputCls">
                <option value="" disabled>Select account</option>
                <option v-for="a in accounts" :key="a.ID" :value="a.ID">
                  {{ a.Name }} ({{ a.Currency }})
                </option>
              </select>
            </div>
            <div>
              <label for="tf-pair" :class="labelCls">Pair</label>
              <input
                id="tf-pair" v-model="draft.pair" list="tf-pairs" required
                placeholder="EUR/USD" :class="inputCls" class="tnum uppercase"
              >
              <datalist id="tf-pairs">
                <option v-for="p in PAIRS" :key="p" :value="p" />
              </datalist>
            </div>
            <div>
              <span id="tf-dir-label" :class="labelCls">Direction</span>
              <div class="grid grid-cols-2 gap-2" role="group" aria-labelledby="tf-dir-label">
                <button
                  v-for="d in (['long', 'short'] as const)" :key="d" type="button"
                  :aria-pressed="draft.direction === d"
                  :class="[
                    'rounded-xl border py-2.5 text-sm font-semibold capitalize transition',
                    draft.direction === d
                      ? 'border-primary bg-primary/10 text-primary'
                      : 'border-border text-muted hover:text-text',
                  ]"
                  @click="draft.direction = d"
                >
                  {{ d === 'long' ? '↑ Long' : '↓ Short' }}
                </button>
              </div>
            </div>
            <div>
              <label for="tf-tf" :class="labelCls">Timeframe</label>
              <select id="tf-tf" v-model="draft.timeframe" :class="inputCls">
                <option v-for="t in TIMEFRAMES" :key="t" :value="t">{{ t }}</option>
              </select>
            </div>
            <div>
              <label for="tf-setup" :class="labelCls">Setup</label>
              <select id="tf-setup" v-model="draft.setup_id" :class="inputCls">
                <option value="">No setup</option>
                <option v-for="s in setups" :key="s.ID" :value="s.ID">{{ s.Name }}</option>
              </select>
            </div>
            <div>
              <label for="tf-opened" :class="labelCls">Opened at</label>
              <input
                id="tf-opened" v-model="draft.opened_at" type="datetime-local"
                :class="`${inputCls} tnum`"
              >
              <p class="mt-1 text-[11px] text-muted">Times are in {{ tz }}.</p>
            </div>
          </div>
        </section>

        <!-- Prices and size -->
        <section class="rounded-2xl border border-border bg-surface p-4 md:p-5" aria-label="Prices and size">
          <div class="flex items-baseline justify-between">
            <h2 class="text-sm font-semibold">Prices and size</h2>
            <p class="text-[11px] text-muted">Risk and R are calculated by the server on save</p>
          </div>
          <div class="mt-3 grid gap-3 sm:grid-cols-3">
            <div>
              <label for="tf-entry" :class="labelCls">Entry</label>
              <input id="tf-entry" v-model.number="draft.entry" type="number" step="any" required placeholder="1.08520" :class="`${inputCls} tnum`">
            </div>
            <div>
              <label for="tf-sl" :class="labelCls">Stop loss</label>
              <input id="tf-sl" v-model.number="draft.stop_loss" type="number" step="any" required placeholder="1.08320" :class="`${inputCls} tnum`">
            </div>
            <div>
              <label for="tf-tp" :class="labelCls">Take profit</label>
              <input id="tf-tp" v-model.number="draft.take_profit" type="number" step="any" placeholder="1.08940" :class="`${inputCls} tnum`">
            </div>
            <div>
              <label for="tf-lot" :class="labelCls">Lot size</label>
              <input id="tf-lot" v-model.number="draft.lot_size" type="number" step="any" min="0" required placeholder="0.50" :class="`${inputCls} tnum`">
            </div>
            <div>
              <label for="tf-exit" :class="labelCls">Exit <span class="font-normal text-muted">(when closed)</span></label>
              <input id="tf-exit" v-model.number="draft.exit_price" type="number" step="any" placeholder="Leave empty if open" :class="`${inputCls} tnum`">
            </div>
            <div class="grid grid-cols-2 gap-2">
              <div>
                <label for="tf-comm" :class="labelCls">Comm.</label>
                <input id="tf-comm" v-model.number="draft.commission" type="number" step="any" min="0" placeholder="0.00" :class="`${inputCls} tnum`">
              </div>
              <div>
                <label for="tf-swap" :class="labelCls">Swap</label>
                <input id="tf-swap" v-model.number="draft.swap" type="number" step="any" placeholder="0.00" :class="`${inputCls} tnum`">
              </div>
            </div>
          </div>
        </section>

        <!-- Review -->
        <section class="rounded-2xl border border-border bg-surface p-4 md:p-5" aria-label="Review">
          <h2 class="text-sm font-semibold">Review</h2>
          <label class="mt-3 flex cursor-pointer items-center justify-between gap-3">
            <span>
              <span class="block text-sm font-medium">I followed my rules</span>
              <span class="block text-xs text-muted">Be honest. This feeds your rule-following rate.</span>
            </span>
            <input v-model="draft.followed_rules" type="checkbox" class="peer sr-only">
            <span class="relative h-6 w-11 shrink-0 rounded-full bg-border transition peer-checked:bg-primary peer-focus-visible:ring-2 peer-focus-visible:ring-primary after:absolute after:left-0.5 after:top-0.5 after:h-5 after:w-5 after:rounded-full after:bg-white after:transition peer-checked:after:translate-x-5" />
          </label>
          <div class="mt-3 grid gap-3 sm:grid-cols-2">
            <div>
              <label for="tf-emotion" :class="labelCls">Emotion</label>
              <select id="tf-emotion" v-model="draft.emotion" :class="inputCls">
                <option value="">—</option>
                <option v-for="e in EMOTIONS" :key="e" :value="e">{{ e }}</option>
              </select>
            </div>
            <div>
              <span :class="labelCls">Tags</span>
              <div class="flex flex-wrap gap-1.5">
                <button
                  v-for="t in tags" :key="t.ID" type="button"
                  :aria-pressed="draft.tag_ids.includes(t.ID)"
                  :class="[
                    'rounded-full border px-2.5 py-1 text-xs font-medium transition',
                    draft.tag_ids.includes(t.ID)
                      ? 'border-primary bg-primary/10 text-primary'
                      : 'border-border text-muted hover:text-text',
                  ]"
                  @click="toggleTag(t.ID)"
                >
                  {{ t.Kind === 'mistake' ? '⚠ ' : '' }}{{ t.Name }}
                </button>
                <span v-if="tags.length === 0" class="text-xs text-muted">No tags yet — add them under Tags.</span>
              </div>
            </div>
          </div>
          <div class="mt-3">
            <label for="tf-notes" :class="labelCls">Notes</label>
            <textarea
              id="tf-notes" v-model="draft.notes" rows="3"
              placeholder="What did you see, and what would you repeat?"
              class="w-full resize-y rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary"
            />
          </div>
        </section>

        <!-- Screenshots -->
        <section class="rounded-2xl border border-border bg-surface p-4 md:p-5" aria-label="Screenshots">
          <div class="flex items-baseline justify-between">
            <h2 class="text-sm font-semibold">Screenshots</h2>
            <p class="text-xs text-muted">{{ pending.length }} of 5</p>
          </div>
          <ul v-if="pending.length > 0" class="mt-3 space-y-2">
            <li v-for="(img, i) in pending" :key="img.preview" class="flex items-center gap-3 rounded-xl border border-border p-2.5">
              <img :src="img.preview" alt="Screenshot preview" class="h-11 w-11 rounded-lg object-cover">
              <div class="min-w-0 flex-1">
                <p class="truncate text-xs font-medium">{{ img.file.name }}</p>
                <p class="text-[11px] text-muted">{{ (img.file.size / 1024).toFixed(0) }} KB · WebP</p>
              </div>
              <select
                v-model="img.kind" aria-label="Image kind"
                class="rounded-lg border border-border bg-bg px-2 py-1.5 text-xs outline-none focus:border-primary"
              >
                <option value="entry">Entry</option>
                <option value="exit">Exit</option>
              </select>
              <button
                type="button" aria-label="Remove image"
                class="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-muted transition hover:bg-loss/10 hover:text-loss"
                @click="removePending(i)"
              >
                ✕
              </button>
            </li>
          </ul>
          <label
            class="mt-3 flex cursor-pointer items-center justify-center gap-2 rounded-xl border border-dashed border-border px-4 py-3 text-sm font-medium text-muted transition hover:border-primary hover:text-text"
          >
            <input type="file" accept="image/*" multiple class="sr-only" @change="onFiles">
            + Add screenshot
          </label>
          <p v-if="fileError" class="mt-2 text-xs text-warning-text">{{ fileError }}</p>
          <p class="mt-2 text-[11px] text-muted">Tag each image as entry or exit. Images are shrunk to WebP under 2 MB before upload.</p>
        </section>
      </div>

      <!-- Side rail: preview + visibility + save -->
      <div class="space-y-4 xl:sticky xl:top-4 xl:self-start">
        <section class="rounded-2xl border border-border bg-surface p-4" aria-label="Preview">
          <div class="flex items-center justify-between">
            <h2 class="text-sm font-semibold">Preview</h2>
            <span class="rounded-full bg-bg px-2 py-0.5 text-[11px] text-muted">Display only</span>
          </div>
          <dl class="tnum mt-3 space-y-1.5 text-sm">
            <div class="flex justify-between"><dt class="text-muted">Risk</dt><dd class="font-semibold">{{ preview.riskPips === null ? '—' : `${preview.riskPips.toFixed(1)} pips` }}</dd></div>
            <div class="flex justify-between"><dt class="text-muted">Reward</dt><dd class="font-semibold">{{ preview.rewardPips === null ? '—' : `${preview.rewardPips.toFixed(1)} pips` }}</dd></div>
            <div class="flex justify-between"><dt class="text-muted">Estimated risk</dt><dd class="font-semibold">{{ preview.estRisk === null ? '—' : `$${preview.estRisk.toFixed(0)}` }}</dd></div>
            <div class="flex justify-between border-t border-border pt-1.5">
              <dt class="text-muted">Planned R:R</dt>
              <dd class="text-base font-bold">{{ preview.rr === null ? '—' : `1 : ${preview.rr.toFixed(2)}` }}</dd>
            </div>
          </dl>
        </section>

        <section class="rounded-2xl border border-border bg-surface p-4" aria-label="Visibility">
          <label class="flex cursor-pointer items-center justify-between gap-3">
            <span>
              <span class="block text-sm font-medium">Make public</span>
              <span class="block text-xs text-muted">Shows pair, setup and R only. Never money or lots.</span>
            </span>
            <input v-model="draft.make_public" type="checkbox" class="peer sr-only">
            <span class="relative h-6 w-11 shrink-0 rounded-full bg-border transition peer-checked:bg-primary peer-focus-visible:ring-2 peer-focus-visible:ring-primary after:absolute after:left-0.5 after:top-0.5 after:h-5 after:w-5 after:rounded-full after:bg-white after:transition peer-checked:after:translate-x-5" />
          </label>
          <button
            type="submit"
            :disabled="!valid || saving"
            class="mt-3 w-full rounded-xl bg-primary py-2.5 text-sm font-semibold text-on-primary transition hover:opacity-90 disabled:opacity-50"
          >
            {{ saving ? 'Saving…' : submitLabel }}
          </button>
          <button
            type="button"
            class="mt-2 w-full rounded-xl py-2 text-sm font-medium text-muted transition hover:text-text"
            @click="$emit('cancel')"
          >
            Cancel
          </button>
        </section>
      </div>
    </div>
  </form>
</template>
