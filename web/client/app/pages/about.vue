<script setup lang="ts">
definePageMeta({ layout: 'public' })

const title = 'About | Ctemz Journal'
const description =
  'A journal for the trader you are becoming. Document your trades, track what you learn, and see your progress over time.'
useSeoMeta({ title, description, ogTitle: title, ogDescription: description })

/* ---------- 1. Where learning comes from ---------- */
const sources = [
  { text: 'Market structure from a course.' },
  { text: 'Risk management from a book.' },
  { text: 'Technical analysis from a mentor.' },
  { text: 'Trading psychology from experience.' },
  { text: 'Something you discover from studying charts or reviewing one of your own trades.', emphasis: true },
]

/* ---------- 2. The process ---------- */
const loop = [
  { step: 'Learn', body: 'Understand how the market works, how to read price, and what risk actually does to an account.' },
  { step: 'Practice', body: 'Apply what you have learned before any real money is involved.' },
  { step: 'Test', body: 'Backtest your ideas against historical data instead of assuming they work.' },
  { step: 'Trade', body: 'Take trades with defined risk.' },
  { step: 'Journal', body: 'Record the setup, entry, stop loss, target, risk, result, screenshots and how you felt.' },
  { step: 'Review', body: 'See which setups work, where you repeat mistakes, and whether you kept your rules.' },
]
const activeStep = ref(0)

const activeStepData = computed(
  () => loop[activeStep.value] ?? { step: '', body: '' },
)

/* ---------- 3. Your trades tell a story ---------- */
const questions = [
  'Why did you take it?',
  'What did you see?',
  'What was your setup?',
  'Where was your stop loss?',
  'How much were you willing to risk?',
  'Did the trade follow your rules?',
  'What were you thinking before entering?',
  'What happened after the trade?',
  'What did you learn?',
]

const insights = [
  'The setups you trade best',
  'The mistakes you repeat',
  'Whether you are following your rules',
  'How your risk affects your results',
  'How your expectations compare with what actually happens',
]

/* ---------- 4. Profit is not the whole story ---------- */
const outcomes = [
  {
    tag: 'Followed the rules and lost',
    goodProcess: true,
    trade: 'EUR/USD, long',
    rules: 'Yes',
    result: 'Loss -1.0R',
    winning: false,
    verdict: 'A good decision. You managed risk and executed your setup. The loss was the cost of doing business.',
  },
  {
    tag: 'Broke the rules and won',
    goodProcess: false,
    trade: 'GBP/JPY, long',
    rules: 'No. Stop moved, size doubled',
    result: 'Win 2.4R',
    winning: true,
    verdict: 'A bad decision. The market rewarded a mistake. That is not the same as a good trade.',
  },
]

/* ---------- 5. Risk comes first ---------- */
const risk = 20
const reward = 42
const rMultiple = (reward / risk).toFixed(1)
const riskShare = (risk / (risk + reward)) * 100

/* ---------- 6. Psychology ---------- */
const emotions = [
  { name: 'Fear', body: 'can make you close a good trade too early' },
  { name: 'Greed', body: 'can make you take more risk than you planned' },
  { name: 'FOMO', body: 'can make you enter a trade that was never part of your setup' },
  { name: 'Revenge trading', body: 'can turn one loss into several' },
  { name: 'Overconfidence', body: 'can make a winning streak dangerous' },
]

/* ---------- 7. Hype ---------- */
const hype = [
  'Huge profit screenshots',
  'Luxury lifestyles',
  'Guaranteed strategies',
  'Signal sellers',
  'Unrealistic win rates',
  'Promises of effortless money',
]

/* ---------- 8. Privacy demo ---------- */
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

/* ---------- 9. See how far you have come ---------- */
const changes = [
  'Your understanding will change.',
  'Your strategy will evolve.',
  'Your risk management will improve.',
  'Your mistakes will change.',
  'You will have good periods and bad periods.',
  'There will be times when you question yourself and times when things finally start to make sense.',
]

/* ---------- 10. Built for the long term ---------- */
const audience = [
  'The beginner who is still learning the basics.',
  'The trader testing their first strategy.',
  'The trader trying to become more disciplined.',
  'The trader who wants to stop relying on memory and start relying on evidence.',
  'Anyone who wants to take their trading development seriously.',
]

const signature = [
  'Learn from wherever you choose',
  'Practice what you learn',
  'Journal what you do',
  'Review what happened',
  'Improve from the evidence',
]
</script>

<template>
  <div class="flex-1 pb-4">
    <!-- ============ 1. Hero ============ -->
    <header
      class="mx-auto grid w-full max-w-6xl items-center gap-10 px-4 pt-12 md:px-6 md:pt-20 lg:grid-cols-[minmax(0,1.05fr)_minmax(0,1fr)] lg:gap-14"
    >
      <div>
        <p class="text-xs font-semibold uppercase tracking-wide text-primary">About Ctemz Journal</p>
        <h1 class="mt-2 text-balance text-4xl font-extrabold leading-[1.06] tracking-tight md:text-5xl">
          A journal for the trader you are becoming.
        </h1>
        <div class="mt-5 max-w-xl space-y-3 text-pretty text-base leading-relaxed text-muted md:text-lg">
          <p>
            Forex trading is not something you master by finding one perfect strategy, copying someone
            else's trades, or watching enough videos to feel ready.
          </p>
          <p>
            You learn by studying the market, testing what you learn, making decisions, taking trades,
            reviewing those decisions, and doing it again.
          </p>
          <p>That is what Ctemz Journal is built around.</p>
        </div>
      </div>

      <aside
        class="rounded-2xl border border-border bg-surface p-5 shadow-xl md:p-6"
        aria-label="Where learning comes from"
      >
        <p class="text-sm leading-relaxed text-muted">
          Ctemz Journal is a personal space for traders to document their journey, keep track of what
          they are learning, record their trades, and look back at their progress over time.
        </p>

        <ol class="mt-5">
          <li
            v-for="(s, i) in sources"
            :key="s.text"
            class="relative flex gap-4 pb-5 last:pb-0"
          >
            <span
              v-if="i < sources.length - 1"
              class="absolute bottom-0 left-2.75 top-7 w-px bg-border"
              aria-hidden="true"
            />
            <span
              class="tnum relative z-10 mt-0.5 grid h-6 w-6 shrink-0 place-items-center rounded-full border border-border bg-bg text-[11px] font-bold text-muted"
            >
              {{ i + 1 }}
            </span>
            <p
              class="text-sm leading-relaxed"
              :class="s.emphasis ? 'font-semibold text-text' : 'text-muted'"
            >
              {{ s.text }}
            </p>
          </li>
        </ol>

        <p class="mt-5 border-t border-border pt-5 text-sm font-semibold leading-relaxed text-text">
          You don't learn from Ctemz Journal. You keep a journal of what you learn.
        </p>
      </aside>
    </header>

    <!-- ============ 2. The process ============ -->
    <section
      class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6"
      aria-labelledby="process-title"
    >
      <div class="max-w-2xl">
        <p class="text-xs font-semibold uppercase tracking-wide text-primary">
          Learn. Practice. Test. Trade. Journal. Review.
        </p>
        <h2 id="process-title" class="mt-1 text-balance text-2xl font-bold tracking-tight md:text-3xl">
          Becoming a better trader is a process.
        </h2>
      </div>

      <div class="mt-8 rounded-3xl border border-border bg-surface p-4 md:p-6">
        <div class="flex flex-wrap items-center gap-2" role="group" aria-label="The trading process">
          <button
            v-for="(s, i) in loop"
            :key="s.step"
            type="button"
            class="flex items-center gap-2 rounded-full px-3.5 py-2 text-sm font-semibold transition-colors"
            :class="
              activeStep === i
                ? 'bg-bg text-text ring-1 ring-border'
                : 'text-muted hover:bg-bg hover:text-text'
            "
            :aria-pressed="activeStep === i"
            @click="activeStep = i"
          >
            <span
              class="tnum grid h-5 w-5 place-items-center rounded-full text-[11px]"
              :class="activeStep === i ? 'bg-primary/15 text-primary' : 'bg-bg'"
            >
              {{ i + 1 }}
            </span>
            {{ s.step }}
          </button>
        </div>

        <div class="mt-4 rounded-2xl bg-bg px-4 py-3">
          <p class="text-sm leading-relaxed text-muted" aria-live="polite">
            <span class="font-semibold text-text">{{ activeStepData.step }}.</span>
            {{ activeStepData.body }}
          </p>
        </div>
      </div>

      <p class="mt-6 max-w-2xl text-pretty text-base font-medium leading-relaxed">
        Over time, your journal becomes more than a collection of trades. It becomes evidence of how you
        think, how you trade, where you improve, and where you keep making the same mistakes.
      </p>
    </section>

    <!-- ============ 3. Your trades tell a story ============ -->
    <section
      class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6"
      aria-labelledby="story-title"
    >
      <p class="text-xs font-semibold uppercase tracking-wide text-primary">Your trades tell a story</p>
      <h2 id="story-title" class="mt-1 max-w-2xl text-balance text-2xl font-bold tracking-tight md:text-3xl">
        A trade is more than an entry and an exit.
      </h2>

      <div class="mt-8 grid gap-4 lg:grid-cols-2 lg:gap-6">
        <div class="rounded-3xl border border-border bg-surface p-5 md:p-6">
          <h3 class="text-xs font-semibold uppercase tracking-wide text-muted">
            Ask yourself after every trade
          </h3>
          <ul class="mt-4 space-y-2.5">
            <li
              v-for="q in questions"
              :key="q"
              class="flex items-start gap-3 text-sm"
            >
              <span
                class="mt-2 h-1 w-1 shrink-0 rounded-full bg-primary"
                aria-hidden="true"
              />
              <span class="leading-snug">{{ q }}</span>
            </li>
          </ul>
        </div>

        <div class="rounded-3xl border border-border bg-surface p-5 md:p-6">
          <h3 class="text-xs font-semibold uppercase tracking-wide text-muted">
            What the record gives you
          </h3>
          <p class="mt-4 text-sm leading-relaxed text-muted">
            Recording these things gives you something that memory cannot. A record.
          </p>
          <p class="mt-4 text-sm leading-relaxed text-muted">
            Instead of looking back and saying, “I think I have improved,” you can actually see it.
          </p>
          <ul class="mt-4 space-y-2.5">
            <li v-for="i in insights" :key="i" class="flex items-start gap-3 text-sm">
              <svg
                class="mt-0.5 h-4 w-4 shrink-0 text-primary"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2.2"
                stroke-linecap="round"
                stroke-linejoin="round"
                aria-hidden="true"
              >
                <path d="M20 6 9 17l-5-5" />
              </svg>
              <span class="leading-snug">{{ i }}</span>
            </li>
          </ul>
        </div>
      </div>
    </section>

    <!-- ============ 4. Profit is not the whole story ============ -->
    <section
      class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6"
      aria-labelledby="outcome-title"
    >
      <p class="text-xs font-semibold uppercase tracking-wide text-primary">
        Profit is not the whole story
      </p>
      <h2 id="outcome-title" class="mt-1 max-w-2xl text-balance text-2xl font-bold tracking-tight md:text-3xl">
        A winning trade is not automatically a good trade. A losing trade is not automatically a bad trade.
      </h2>

      <div class="mt-8 grid gap-4 md:grid-cols-2 md:gap-6">
        <article
          v-for="o in outcomes"
          :key="o.tag"
          class="rounded-3xl border border-border bg-surface p-5 md:p-6"
        >
          <p
            class="text-xs font-semibold uppercase tracking-wide"
            :class="o.goodProcess ? 'text-primary' : 'text-muted'"
          >
            {{ o.tag }}
          </p>

          <dl class="mt-4 space-y-2 text-sm">
            <div class="flex items-baseline justify-between gap-4">
              <dt class="text-muted">Trade</dt>
              <dd class="font-semibold">{{ o.trade }}</dd>
            </div>
            <div class="flex items-baseline justify-between gap-4">
              <dt class="text-muted">Rules followed</dt>
              <dd class="text-right font-semibold">{{ o.rules }}</dd>
            </div>
          </dl>

          <div class="mt-4 flex items-center justify-between rounded-xl bg-bg px-4 py-3">
            <span class="text-sm text-muted">Result</span>
            <span
              class="tnum rounded-full px-2.5 py-0.5 text-sm font-semibold"
              :class="o.winning ? 'bg-profit/10 text-profit-text' : 'bg-border text-text'"
            >
              {{ o.result }}
            </span>
          </div>

          <p
            class="mt-4 text-sm leading-relaxed"
            :class="o.goodProcess ? 'font-medium text-text' : 'text-muted'"
          >
            {{ o.verdict }}
          </p>
        </article>
      </div>

      <p class="mt-6 max-w-2xl text-pretty leading-relaxed text-muted">
        The goal is to help you understand the decisions behind your results. If you followed your plan,
        managed your risk correctly, and executed your setup properly, a losing trade can still be a
        good trade. If you ignored your rules and somehow made money, that does not make the decision
        good. The journal helps you separate the outcome from the process.
      </p>
    </section>

    <!-- ============ 5. Risk comes first ============ -->
    <section
      class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6"
      aria-labelledby="risk-title"
    >
      <div
        class="grid gap-8 rounded-3xl border border-border bg-surface p-6 md:p-10 lg:grid-cols-[minmax(0,1fr)_minmax(0,1.1fr)] lg:gap-14"
      >
        <div>
          <p class="text-xs font-semibold uppercase tracking-wide text-primary">Risk comes first</p>
          <h2 id="risk-title" class="mt-1 text-balance text-2xl font-bold tracking-tight md:text-3xl">
            Before asking how much you can make, understand how much you can lose.
          </h2>
          <div class="mt-4 space-y-3 text-pretty leading-relaxed text-muted">
            <p>Trading is about managing risk before chasing returns.</p>
            <p>
              Ctemz Journal gives you a place to record the risk behind your trades, understand your
              reward relative to that risk, and see how your decisions affect your overall performance.
            </p>
            <p>
              Staying in the market long enough to learn is more important than trying to get rich from
              one trade.
            </p>
          </div>
        </div>

        <figure class="self-center" aria-label="A trade that risked 20 pips and made 42 pips">
          <div class="relative pb-14 pt-9">
            <div class="tnum absolute inset-x-0 top-0 text-xs text-muted">
              <span class="absolute left-0">Stop 1.08320</span>
              <span class="absolute -translate-x-1/2" :style="{ left: riskShare + '%' }">Entry 1.08520</span>
              <span class="absolute right-0">Exit 1.08940</span>
            </div>

            <div class="flex h-3.5 w-full overflow-hidden rounded-full">
              <div class="bg-border" :style="{ width: riskShare + '%' }" />
              <div class="bg-primary" :style="{ width: 100 - riskShare + '%' }" />
            </div>
            <span
              class="absolute top-8.5 h-6 w-0.5 -translate-x-1/2 bg-text"
              :style="{ left: riskShare + '%' }"
              aria-hidden="true"
            />

            <div class="tnum absolute inset-x-0 bottom-0 flex text-sm font-semibold">
              <span class="text-muted" :style="{ width: riskShare + '%' }">Risk {{ risk }} pips</span>
              <span class="text-primary" :style="{ width: 100 - riskShare + '%' }">
                Reward {{ reward }} pips
              </span>
            </div>
          </div>

          <figcaption class="mt-4 flex items-center gap-3 rounded-xl bg-bg px-4 py-3 text-sm">
            <span class="tnum text-muted">{{ reward }} ÷ {{ risk }} =</span>
            <span
              class="tnum inline-flex items-center gap-1 rounded-full bg-profit/10 px-2.5 py-0.5 font-semibold text-profit-text"
            >
              Win +{{ rMultiple }}R
            </span>
          </figcaption>
        </figure>
      </div>
    </section>

    <!-- ============ 6. Test your ideas ============ -->
    <section
      class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6"
      aria-labelledby="test-title"
    >
      <p class="text-xs font-semibold uppercase tracking-wide text-primary">Test your ideas</p>
      <h2 id="test-title" class="mt-1 max-w-2xl text-balance text-2xl font-bold tracking-tight md:text-3xl">
        A strategy does not become yours because someone on the internet says it works.
      </h2>

      <div class="mt-8 grid items-stretch gap-4 md:grid-cols-[minmax(0,1fr)_auto_minmax(0,1fr)] md:gap-6">
        <div class="rounded-3xl border border-border bg-bg p-5 md:p-6">
          <p class="text-xs font-semibold uppercase tracking-wide text-muted">Before</p>
          <p class="mt-3 text-lg font-semibold italic text-muted">
            “I think this strategy works.”
          </p>
        </div>

        <div class="grid place-items-center text-muted" aria-hidden="true">
          <svg
            class="h-5 w-5 rotate-90 md:rotate-0"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M5 12h14M13 6l6 6-6 6" />
          </svg>
        </div>

        <div class="rounded-3xl border border-primary/30 bg-primary/5 p-5 md:p-6">
          <p class="text-xs font-semibold uppercase tracking-wide text-primary">After</p>
          <p class="mt-3 text-lg font-semibold">
            “What happened when I tested these rules myself?”
          </p>
        </div>
      </div>

      <p class="mt-6 max-w-2xl text-pretty leading-relaxed text-muted">
        Backtest it. Study it. Understand its conditions. See where it works and where it does not.
        Then decide whether it deserves a place in your trading plan. Ctemz Journal gives you somewhere
        to keep that process documented. Your journal can become a record of the ideas you tested, the
        strategies you practiced, and the evidence you collected along the way.
      </p>
    </section>

    <!-- ============ 7. The psychological side matters too ============ -->
    <section
      class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6"
      aria-labelledby="psychology-title"
    >
      <p class="text-xs font-semibold uppercase tracking-wide text-primary">
        The psychological side matters too
      </p>
      <h2 id="psychology-title" class="mt-1 max-w-2xl text-balance text-2xl font-bold tracking-tight md:text-3xl">
        Trading is not only about charts.
      </h2>

      <ul class="mt-8 grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        <li
          v-for="e in emotions"
          :key="e.name"
          class="rounded-2xl border border-border bg-surface p-5"
        >
          <p class="text-base font-bold tracking-tight">{{ e.name }}</p>
          <p class="mt-1 text-sm leading-relaxed text-muted">{{ e.body }}</p>
        </li>
      </ul>

      <p class="mt-6 max-w-2xl text-pretty leading-relaxed text-muted">
        These behaviors are difficult to see when you only look at your account balance. Writing them
        down makes them easier to recognize. The more honestly you journal, the easier it becomes to
        identify the behavior behind your decisions.
      </p>
    </section>

    <!-- ============ 8. No hype ============ -->
    <section class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6" aria-labelledby="hype-title">
      <div class="rounded-3xl border border-border bg-surface p-6 md:p-10">
        <p class="text-xs font-semibold uppercase tracking-wide text-primary">
          No hype. Just the work.
        </p>
        <h2 id="hype-title" class="mt-1 max-w-2xl text-balance text-2xl font-bold tracking-tight md:text-3xl">
          There is a lot of noise around forex.
        </h2>

        <ul class="mt-6 flex flex-wrap gap-2">
          <li
            v-for="h in hype"
            :key="h"
            class="rounded-full border border-border bg-bg px-3.5 py-1.5 text-sm text-muted line-through decoration-border"
          >
            {{ h }}
          </li>
        </ul>

        <div class="mt-8 grid gap-6 border-t border-border pt-8 md:grid-cols-2 md:gap-10">
          <div class="space-y-3 text-pretty leading-relaxed text-muted">
            <p>Ctemz Journal is not built around any of that.</p>
            <p>It is built around the work.</p>
            <p>
              Study the market. Practice. Take calculated risks. Make mistakes. Record them. Review
              them. Improve. Repeat.
            </p>
            <p>There is no shortcut that can replace experience and disciplined practice.</p>
          </div>
          <p class="self-center text-balance text-3xl font-extrabold leading-tight tracking-tight md:text-4xl">
            No hype. Just the work.
          </p>
        </div>
      </div>
    </section>

    <!-- ============ 9. Your journal belongs to you ============ -->
    <section
      class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6"
      aria-labelledby="privacy-title"
    >
      <div class="grid gap-10 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)] lg:gap-14">
        <div>
          <p class="text-xs font-semibold uppercase tracking-wide text-primary">
            Your journal belongs to you
          </p>
          <h2 id="privacy-title" class="mt-1 text-balance text-2xl font-bold tracking-tight md:text-3xl">
            Your trading journey is personal.
          </h2>
          <div class="mt-4 space-y-3 text-pretty leading-relaxed text-muted">
            <p>
              Your trades, notes, mistakes, strategies, results, and observations belong to you. That
              is why Ctemz Journal is private by default.
            </p>
            <p>
              You decide what you want to keep to yourself and what you want to share. Your journal
              should be a place where you can be completely honest about your trading, including the
              decisions you would rather forget.
            </p>
            <p class="font-semibold text-text">
              Because those decisions may be the ones that teach you the most.
            </p>
          </div>
          <p class="mt-5 max-w-xl text-sm leading-relaxed text-muted">
            Try the switch to see exactly what another person gets.
          </p>
        </div>

        <section
          class="rounded-2xl border border-border bg-surface p-4 shadow-xl md:p-5"
          aria-label="Example trade, private and public views"
        >
          <div
            class="mb-4 grid grid-cols-2 gap-1 rounded-xl border border-border bg-bg p-1"
            role="group"
            aria-label="Choose who is looking"
          >
            <button
              type="button"
              class="flex h-10 items-center justify-center gap-2 rounded-lg text-sm font-semibold transition-colors"
              :class="
                view === 'you'
                  ? 'bg-surface text-text shadow-sm ring-1 ring-border'
                  : 'text-muted hover:text-text'
              "
              :aria-pressed="view === 'you'"
              @click="view = 'you'"
            >
              <svg
                class="h-4 w-4"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.8"
                stroke-linecap="round"
                stroke-linejoin="round"
                aria-hidden="true"
              >
                <rect x="5" y="11" width="14" height="9" rx="2" />
                <path d="M8 11V8a4 4 0 018 0v3" />
              </svg>
              You see
            </button>
            <button
              type="button"
              class="flex h-10 items-center justify-center gap-2 rounded-lg text-sm font-semibold transition-colors"
              :class="
                view === 'public'
                  ? 'bg-surface text-text shadow-sm ring-1 ring-border'
                  : 'text-muted hover:text-text'
              "
              :aria-pressed="view === 'public'"
              @click="view = 'public'"
            >
              <svg
                class="h-4 w-4"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.8"
                stroke-linecap="round"
                stroke-linejoin="round"
                aria-hidden="true"
              >
                <path d="M2 12s4-7 10-7 10 7 10 7-4 7-10 7S2 12 2 12z" />
                <circle cx="12" cy="12" r="3" />
              </svg>
              Others see
            </button>
          </div>

          <div class="overflow-hidden rounded-xl border border-border bg-bg">
            <svg viewBox="0 0 120 40" class="h-28 w-full" preserveAspectRatio="none" aria-hidden="true">
              <path
                d="M0,34 L30,30 L60,32 L90,20 L120,10 L120,40 L0,40 Z"
                fill="var(--profit)"
                opacity="0.14"
              />
              <path
                d="M0,34 L30,30 L60,32 L90,20 L120,10"
                fill="none"
                stroke="var(--profit)"
                stroke-width="2.5"
                stroke-linecap="round"
                stroke-linejoin="round"
                vector-effect="non-scaling-stroke"
              />
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
                <svg
                  v-if="!r.shared"
                  class="h-3.5 w-3.5"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="1.8"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  aria-hidden="true"
                >
                  <rect x="5" y="11" width="14" height="9" rx="2" />
                  <path d="M8 11V8a4 4 0 018 0v3" />
                </svg>
              </dt>
              <dd class="min-w-0 text-right text-sm font-semibold">
                <Transition name="swap" mode="out-in">
                  <span
                    v-if="r.shared || view === 'you'"
                    :key="`v-${r.label}`"
                    class="tnum inline-block max-w-full"
                    :class="
                      r.win ? 'rounded-full bg-profit/10 px-2.5 py-0.5 text-profit-text' : ''
                    "
                  >
                    {{ r.value }}
                  </span>
                  <span v-else :key="`m-${r.label}`" class="inline-flex items-center justify-end">
                    <span class="block h-3 rounded bg-border" :class="r.mask" aria-hidden="true" />
                    <span class="sr-only">Hidden from others</span>
                  </span>
                </Transition>
              </dd>
            </div>
          </dl>

          <p
            class="mt-3 rounded-lg bg-bg px-3 py-2 text-xs leading-relaxed text-muted"
            aria-live="polite"
          >
            {{ announce }}
            <template v-if="view === 'public'">
              Crop balances out of your chart before you share it.
            </template>
          </p>
        </section>
      </div>
    </section>

    <!-- ============ 10. See how far you have come ============ -->
    <section
      class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6"
      aria-labelledby="journey-title"
    >
      <div class="grid gap-10 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)] lg:gap-14">
        <div>
          <p class="text-xs font-semibold uppercase tracking-wide text-primary">
            See how far you have come
          </p>
          <h2 id="journey-title" class="mt-1 text-balance text-2xl font-bold tracking-tight md:text-3xl">
            Your first trade will not look like your hundredth.
          </h2>
          <div class="mt-4 space-y-3 text-pretty leading-relaxed text-muted">
            <p>
              Without a record, much of that disappears into memory. With a journal, you can look back
              and see the journey.
            </p>
            <p class="font-semibold text-text">
              Not just how much money you made or lost, but how you became a different trader.
            </p>
          </div>
        </div>

        <ul class="self-center">
          <li
            v-for="c in changes"
            :key="c"
            class="flex items-start gap-3 border-b border-border py-3.5 first:border-t"
          >
            <span class="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary" aria-hidden="true" />
            <span class="text-sm leading-relaxed">{{ c }}</span>
          </li>
        </ul>
      </div>
    </section>

    <!-- ============ 11. Built for the long term ============ -->
    <section class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6">
      <div class="rounded-3xl border border-border bg-surface p-6 text-center md:p-12">
        <p class="text-xs font-semibold uppercase tracking-wide text-primary">Built for the long term</p>
        <h2
          class="mx-auto mt-2 max-w-3xl text-balance text-2xl font-extrabold leading-tight tracking-tight md:text-3xl"
        >
          Ctemz Journal is for traders who understand that becoming good at trading takes time.
        </h2>

        <ul class="mx-auto mt-8 flex max-w-3xl flex-wrap justify-center gap-2">
          <li
            v-for="a in audience"
            :key="a"
            class="rounded-full border border-border bg-bg px-3.5 py-1.5 text-sm font-medium"
          >
            {{ a }}
          </li>
        </ul>

        <div class="mx-auto mt-8 max-w-xl space-y-3 text-pretty text-sm leading-relaxed text-muted">
          <p>
            Ctemz Journal does not promise profitable trades or a winning strategy.
            It does not claim that journaling alone will make you a successful trader.
          </p>
          <p class="text-base font-semibold text-text">
            It gives you something more fundamental. A place to document the work.
          </p>
          <p>
            Because you cannot improve what you never review. And you cannot review what you never
            record.
          </p>
        </div>

        <ul class="mx-auto mt-8 flex max-w-3xl flex-wrap justify-center gap-2">
          <li
            v-for="s in signature"
            :key="s"
            class="rounded-full border border-primary/30 bg-primary/5 px-3.5 py-1.5 text-sm font-medium text-primary"
          >
            {{ s }}
          </li>
        </ul>

        <p class="mt-6 text-sm font-semibold uppercase tracking-wide text-primary">
          That is Ctemz Journal.
        </p>
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