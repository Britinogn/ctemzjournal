<script setup lang="ts">
definePageMeta({ layout: 'public' })

const title = 'About | Ctemz Journal'
const description = 'What Ctemz Journal is, who it is for, and exactly what others can see when you share a trade.'
useSeoMeta({ title, description, ogTitle: title, ogDescription: description })

/* ---------- The memorable part: one trade, two views ---------- */
type View = 'you' | 'public'
const view = ref<View>('you')

const rows = [
  { label: 'Pair', value: 'EUR/USD', shared: true },
  { label: 'Direction', value: 'Long, H1', shared: true },
  { label: 'Setup', value: 'Breakout retest', shared: true },
  { label: 'Result', value: '+2.1R', shared: true, win: true },
  { label: 'Lot size', value: '0.50', shared: false, mask: 'w-10' },
  { label: 'Profit', value: '+$420', shared: false, mask: 'w-16' },
  { label: 'Notes', value: 'Waited for the retest. Kept to my plan.', shared: false, mask: 'w-36' },
]

const announce = computed(() =>
  view.value === 'you'
    ? 'Showing everything only you can see.'
    : 'Showing what other people see. Lot size, profit and notes are hidden.',
)

/* ---------- Four habits, each with a small specimen ---------- */
const points = [
  {
    title: 'Built for forex traders',
    body: 'Pairs, pips, lots and gold, logged the way FX actually trades. The server checks the risk, P&L and R-multiple on every position.',
    kind: 'pairs',
  },
  {
    title: 'Your rules, measured',
    body: 'Tag setups, mistakes and emotions once. Your rule-following rate, expectancy, drawdown and breakdowns by setup and pair update on their own.',
    kind: 'ring',
  },
  {
    title: 'Private by default',
    body: 'Everything is yours alone until you switch on Make public. Even then, only the pair, direction, setup, result in R and your first chart appear. Money and lot size never do.',
    kind: 'locks',
  },
  {
    title: 'Review on any device',
    body: 'Log a trade from your phone in under a minute, then review the calendar and equity curve on a bigger screen. Days are grouped in Lagos time, so daily results match your trading day.',
    kind: 'day',
  },
]
const pairs = ['EUR/USD', 'GBP/JPY', 'XAU/USD', 'USD/NGN']
const locked = ['Lot size', 'Profit', 'Notes']

/* ---------- Worked example: 20 pips risked, 42 made ---------- */
const risk = 20
const reward = 42
const rMultiple = (reward / risk).toFixed(1)
const riskShare = (risk / (risk + reward)) * 100
</script>

<template>
  <div class="flex-1 pb-4">
    <!-- 1. Hero: the promise, then proof of it -->
    <header class="mx-auto grid w-full max-w-6xl items-center gap-10 px-4 pt-12 md:px-6 md:pt-20 lg:grid-cols-[minmax(0,1.05fr)_minmax(0,1fr)] lg:gap-14">
      <div>
        <p class="text-xs font-semibold uppercase tracking-wide text-primary">
          About
        </p>
        <h1 class="mt-2 text-balance text-4xl font-extrabold leading-[1.06] tracking-tight md:text-5xl">
          A journal that shows you why you win, and why you lose.
        </h1>
        <div class="mt-5 max-w-xl space-y-3 text-pretty text-base leading-relaxed text-muted md:text-lg">
          <p>
            Ctemz Journal is a forex trading journal for traders who want evidence, not gut feel.
            Log every position with its setup, risk and screenshots. The server calculates the real
            numbers, so your win rate, expectancy and drawdown are facts, not guesses.
          </p>
          <p>
            Sharing is a choice you make one trade at a time. Try the switch on the right to see
            exactly what other people get.
          </p>
        </div>
      </div>

      <!-- Interactive demo: the same trade, seen by you and seen by everyone else -->
      <section class="rounded-2xl border border-border bg-surface p-4 shadow-xl md:p-5" aria-label="Example trade, private and public views">
        <div class="mb-4 grid grid-cols-2 gap-1 rounded-xl border border-border bg-bg p-1" role="group" aria-label="Choose who is looking">
          <button
            type="button"
            class="flex h-10 items-center justify-center gap-2 rounded-lg text-sm font-semibold transition-colors"
            :class="view === 'you' ? 'bg-surface text-text shadow-sm ring-1 ring-border' : 'text-muted hover:text-text'"
            :aria-pressed="view === 'you'"
            @click="view = 'you'"
          >
            <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="5" y="11" width="14" height="9" rx="2" /><path d="M8 11V8a4 4 0 018 0v3" /></svg>
            You see
          </button>
          <button
            type="button"
            class="flex h-10 items-center justify-center gap-2 rounded-lg text-sm font-semibold transition-colors"
            :class="view === 'public' ? 'bg-surface text-text shadow-sm ring-1 ring-border' : 'text-muted hover:text-text'"
            :aria-pressed="view === 'public'"
            @click="view = 'public'"
          >
            <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M2 12s4-7 10-7 10 7 10 7-4 7-10 7S2 12 2 12z" /><circle cx="12" cy="12" r="3" /></svg>
            Others see
          </button>
        </div>

        <!-- The chart is shared, so it is always visible -->
        <div class="overflow-hidden rounded-xl border border-border bg-bg">
          <svg viewBox="0 0 120 40" class="h-28 w-full" preserveAspectRatio="none" aria-hidden="true">
            <path d="M0,34 L30,30 L60,32 L90,20 L120,10 L120,40 L0,40 Z" fill="var(--profit)" opacity="0.14" />
            <path d="M0,34 L30,30 L60,32 L90,20 L120,10" fill="none" stroke="var(--profit)" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" vector-effect="non-scaling-stroke" />
          </svg>
        </div>

        <dl class="mt-2">
          <div
            v-for="r in rows"
            :key="r.label"
            class="flex min-h-12 items-center justify-between gap-4 border-b border-border py-2 last:border-b-0"
          >
            <dt class="flex items-center gap-1.5 text-sm text-muted">
              {{ r.label }}
              <svg v-if="!r.shared" class="h-3.5 w-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="5" y="11" width="14" height="9" rx="2" /><path d="M8 11V8a4 4 0 018 0v3" /></svg>
            </dt>
            <dd class="min-w-0 text-right text-sm font-semibold">
              <Transition name="swap" mode="out-in">
                <span
                  v-if="r.shared || view === 'you'"
                  :key="`v-${r.label}`"
                  class="tnum inline-block max-w-full"
                  :class="r.win ? 'rounded-full bg-profit/10 px-2.5 py-0.5 text-profit-text' : ''"
                >
                  <template v-if="r.win"><span aria-hidden="true">↑ </span>Win </template>{{ r.value }}
                </span>
                <span v-else :key="`m-${r.label}`" class="inline-flex items-center justify-end">
                  <span class="block h-3 rounded bg-border" :class="r.mask" aria-hidden="true" />
                  <span class="sr-only">Hidden from others</span>
                </span>
              </Transition>
            </dd>
          </div>
        </dl>

        <p class="mt-3 rounded-lg bg-bg px-3 py-2 text-xs leading-relaxed text-muted" aria-live="polite">
          {{ announce }}
          <template v-if="view === 'public'"> Crop balances out of your chart before you share it.</template>
        </p>
      </section>
    </header>

    <!-- 2. Four habits as a ledger, not four identical cards -->
    <section class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6" aria-labelledby="habits-title">
      <p class="text-xs font-semibold uppercase tracking-wide text-primary">
        Why it works
      </p>
      <h2 id="habits-title" class="mt-1 max-w-2xl text-balance text-2xl font-bold tracking-tight md:text-3xl">
        Made for the way forex traders work.
      </h2>

      <ul class="mt-8 border-b border-border">
        <li
          v-for="p in points"
          :key="p.title"
          class="grid gap-4 border-t border-border py-7 md:grid-cols-[minmax(0,1fr)_minmax(0,1.3fr)] md:gap-10"
        >
          <div>
            <h3 class="text-lg font-bold tracking-tight">
              {{ p.title }}
            </h3>

            <!-- A small specimen of the thing the habit is about -->
            <div class="mt-3 flex flex-wrap items-center gap-2" aria-hidden="true">
              <template v-if="p.kind === 'pairs'">
                <span v-for="x in pairs" :key="x" class="tnum rounded-full border border-border bg-surface px-3 py-1 text-xs font-semibold">{{ x }}</span>
              </template>
              <template v-else-if="p.kind === 'ring'">
                <span class="grid h-14 w-14 place-items-center rounded-full" style="background: conic-gradient(var(--primary) 82%, var(--border) 0)">
                  <span class="tnum grid h-10 w-10 place-items-center rounded-full bg-bg text-xs font-semibold">82%</span>
                </span>
                <span class="text-xs text-muted">Rules followed (example)</span>
              </template>
              <template v-else-if="p.kind === 'locks'">
                <span v-for="x in locked" :key="x" class="inline-flex items-center gap-1.5 rounded-full border border-border bg-surface px-3 py-1 text-xs font-medium text-muted">
                  <svg class="h-3 w-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="5" y="11" width="14" height="9" rx="2" /><path d="M8 11V8a4 4 0 018 0v3" /></svg>
                  {{ x }}
                </span>
              </template>
              <template v-else>
                <span class="tnum rounded-full border border-border bg-surface px-3 py-1 text-xs font-semibold">Day closes 00:00 Lagos</span>
              </template>
            </div>
          </div>
          <p class="max-w-prose text-pretty text-base leading-relaxed text-muted">
            {{ p.body }}
          </p>
        </li>
      </ul>
    </section>

    <!-- 3. The numbers, worked out -->
    <section class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6" aria-labelledby="maths-title">
      <div class="grid gap-8 rounded-3xl border border-border bg-surface p-6 md:p-10 lg:grid-cols-[minmax(0,1fr)_minmax(0,1.2fr)] lg:gap-14">
        <div>
          <p class="text-xs font-semibold uppercase tracking-wide text-primary">
            No black box
          </p>
          <h2 id="maths-title" class="mt-1 text-balance text-2xl font-bold tracking-tight md:text-3xl">
            Every number is worked out the same way.
          </h2>
          <p class="mt-3 text-pretty text-base leading-relaxed text-muted">
            R is your result measured in units of what you risked. Risk 20 pips and make 42,
            and that trade is <span class="tnum font-semibold text-text">+{{ rMultiple }}R</span>.
            Expectancy is the average R you earn per trade:
            win rate times average win, minus loss rate times average loss.
          </p>
        </div>

        <figure class="self-center" aria-label="A trade that risked 20 pips and made 42 pips">
          <div class="relative pb-14 pt-9">
            <!-- Labels above the markers -->
            <div class="tnum absolute inset-x-0 top-0 text-xs text-muted">
              <span class="absolute left-0">Stop 1.08320</span>
              <span class="absolute -translate-x-1/2" :style="{ left: riskShare + '%' }">Entry 1.08520</span>
              <span class="absolute right-0">Exit 1.08940</span>
            </div>
            <!-- Risk and reward drawn to scale -->
            <div class="flex h-3.5 w-full overflow-hidden rounded-full">
              <div class="bg-border" :style="{ width: riskShare + '%' }" />
              <div class="bg-primary" :style="{ width: 100 - riskShare + '%' }" />
            </div>
            <span class="absolute top-8.5 h-6 w-0.5 -translate-x-1/2 bg-text" :style="{ left: riskShare + '%' }" aria-hidden="true" />
            <!-- Distances below -->
            <div class="tnum absolute inset-x-0 bottom-0 flex text-sm font-semibold">
              <span class="text-muted" :style="{ width: riskShare + '%' }">Risk {{ risk }} pips</span>
              <span class="text-primary" :style="{ width: 100 - riskShare + '%' }">Reward {{ reward }} pips</span>
            </div>
          </div>
          <figcaption class="mt-4 flex items-center gap-3 rounded-xl bg-bg px-4 py-3 text-sm">
            <span class="tnum text-muted">{{ reward }} ÷ {{ risk }} =</span>
            <span class="tnum inline-flex items-center gap-1 rounded-full bg-profit/10 px-2.5 py-0.5 font-semibold text-profit-text">
              <span aria-hidden="true">↑</span> Win +{{ rMultiple }}R
            </span>
          </figcaption>
        </figure>
      </div>
    </section>

    <HomeCta />
  </div>
</template>

<style scoped>
.swap-enter-active,
.swap-leave-active {
  transition: opacity 0.16s ease;
}
.swap-enter-from,
.swap-leave-to {
  opacity: 0;
}
@media (prefers-reduced-motion: reduce) {
  .swap-enter-active,
  .swap-leave-active {
    transition: none;
  }
}
</style>