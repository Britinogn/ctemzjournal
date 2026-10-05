<script setup lang="ts">
definePageMeta({ layout: 'public' })

const title = 'About | Ctemz Journal'
const description =
  'Ctemz Journal is a forex learning and trading journal built for people who want to become better traders through understanding, practice, discipline and honest review.'
useSeoMeta({ title, description, ogTitle: title, ogDescription: description })

const activeStepData = computed(
  () => loop[activeStep.value] ?? { step: '', body: '' },
)

/* ---------- 1. Where it starts ---------- */
const ladder = [
  { text: 'It starts with understanding how the market works.' },
  { text: 'Then learning how to read price.' },
  { text: 'Then managing risk.' },
  { text: 'Then testing your ideas.' },
  { text: 'Then keeping a record of what you actually did.' },
  { text: 'And most importantly, learning from it.', emphasis: true },
]

/* ---------- 2. The loop ---------- */
const loop = [
  { step: 'Learn', body: 'Understand how the market works, how to read price, and what risk actually does to an account.' },
  { step: 'Practise', body: 'Apply what you have learned before any real money is involved.' },
  { step: 'Test', body: 'Backtest your rules against historical data instead of assuming they work.' },
  { step: 'Journal', body: 'Record the setup, entry, stop, target, risk, result, screenshots and how you felt.' },
  { step: 'Review', body: 'See which setups work, where you repeat mistakes, and whether you kept your rules.' },
  { step: 'Improve', body: 'Turn that evidence into a repeatable process you can keep sharpening.' },
]
const activeStep = ref(0)

/* ---------- 3. The beginner path ---------- */
const path = [
  'Learn the fundamentals of forex',
  'Understand market structure',
  'Learn technical analysis',
  'Understand risk management',
  'Learn the difference between trading and gambling',
  'Study trading psychology',
  'Choose a strategy',
  'Backtest it',
  'Practise it',
  'Journal your results',
  'Review your mistakes',
]

/* ---------- 4. What you record / what you get back ---------- */
const entryFields = [
  'Setup',
  'Entry',
  'Stop loss',
  'Take profit',
  'Risk',
  'Result',
  'Screenshots',
  'Emotions',
  'Mistakes',
  'Rules followed',
]

const insights = [
  'Which setups perform best for you',
  'Where you repeatedly make mistakes',
  'Whether you are following your trading rules',
  'How your risk management is performing',
  'Your win rate',
  'Your expectancy',
  'Your drawdown',
  'Your performance across different pairs and setups',
]

/* ---------- 5. Process over outcome ---------- */
const outcomes = [
  {
    tag: 'Followed the rules — and lost',
    goodProcess: true,
    trade: 'EUR/USD, long',
    rules: 'Yes',
    result: '−1.0R',
    winning: false,
    verdict: 'A good decision. The loss was the cost of doing business.',
  },
  {
    tag: 'Broke the rules — and won',
    goodProcess: false,
    trade: 'GBP/JPY, long',
    rules: 'No — stop moved, size doubled',
    result: '+2.4R',
    winning: true,
    verdict: 'A bad decision. The market rewarded a mistake.',
  },
]

/* ---------- 6. Risk comes before profit ---------- */
const risk = 20
const reward = 42
const rMultiple = (reward / risk).toFixed(1)
const riskShare = (risk / (risk + reward)) * 100

/* ---------- 7. Psychology ---------- */
const emotions = [
  { name: 'Fear', body: 'can make you close a good trade too early' },
  { name: 'Greed', body: 'can make you risk more than you planned' },
  { name: 'FOMO', body: 'can make you enter late' },
  { name: 'Revenge trading', body: 'can turn one loss into several' },
  { name: 'Overconfidence', body: 'can make a winning streak dangerous' },
]

/* ---------- 8. Hype ---------- */
const hype = [
  'Luxury lifestyles',
  'Screenshots of huge profits',
  '“Guaranteed” strategies',
  'Signal sellers',
  'Unrealistic win rates',
  'Accounts showing effortless money',
]

/* ---------- 9. Privacy demo ---------- */
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

/* ---------- 10. A record of your journey ---------- */
const changes = [
  'Your understanding will change.',
  'Your strategy may change.',
  'Your risk management will improve.',
  'You’ll make mistakes.',
  'You’ll discover patterns in your behaviour.',
  'You’ll have good periods and bad periods.',
]

/* ---------- 11. Closing ---------- */
const closingSteps = [
  'Learn the market',
  'Understand your risk',
  'Test your ideas',
  'Record your decisions',
  'Review your results',
  'Learn from your mistakes',
  'Repeat',
]
</script>

<template>
  <div class="flex-1 pb-4">
    <!-- ============ 1. Hero: the promise, and where it starts ============ -->
    <header
      class="mx-auto grid w-full max-w-6xl items-center gap-10 px-4 pt-12 md:px-6 md:pt-20 lg:grid-cols-[minmax(0,1.05fr)_minmax(0,1fr)] lg:gap-14"
    >
      <div>
        <p class="text-xs font-semibold uppercase tracking-wide text-primary">About Ctemz Journal</p>
        <h1 class="mt-2 text-balance text-4xl font-extrabold leading-[1.06] tracking-tight md:text-5xl">
          Trading should be learned, not guessed.
        </h1>
        <div class="mt-5 max-w-xl space-y-3 text-pretty text-base leading-relaxed text-muted md:text-lg">
          <p>
            Ctemz Journal is a forex learning and trading journal built for people who want to become
            better traders through understanding, practice, discipline and honest review.
          </p>
          <p>
            Becoming a trader is more than finding an entry, placing a trade and hoping the market moves
            your way. It is a skill — and skills are built in steps, not in a single lucky week.
          </p>
        </div>
      </div>

      <aside
        class="rounded-2xl border border-border bg-surface p-5 shadow-xl md:p-6"
        aria-label="Where the journey starts"
      >
        <ol>
          <li
            v-for="(s, i) in ladder"
            :key="s.text"
            class="relative flex gap-4 pb-5 last:pb-0"
          >
            <span
              v-if="i < ladder.length - 1"
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
      </aside>
    </header>

    <!-- ============ 2. Start with understanding ============ -->
    <section
      class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6"
      aria-labelledby="understanding-title"
    >
      <div class="max-w-2xl">
        <p class="text-xs font-semibold uppercase tracking-wide text-primary">Start with understanding</p>
        <h2 id="understanding-title" class="mt-1 text-balance text-2xl font-bold tracking-tight md:text-3xl">
          Forex can look complicated when you're just getting started.
        </h2>
        <div class="mt-4 space-y-3 text-pretty text-base leading-relaxed text-muted">
          <p>
            Charts are moving. Candles are forming. People are talking about market structure,
            liquidity, support and resistance, risk to reward, strategies, indicators, psychology and
            countless other concepts.
          </p>
          <p>It is easy to consume information without actually understanding it.</p>
          <p>Ctemz Journal is built around a simpler approach.</p>
        </div>
      </div>

      <!-- The loop, made tangible -->
      <div class="mt-8 rounded-3xl border border-border bg-surface p-4 md:p-6">
        <div class="flex flex-wrap items-center gap-2" role="group" aria-label="The learning loop">
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
        You don't need to know everything on day one. You need a clear path to keep learning.
      </p>
    </section>

    <!-- ============ 3. Built for the beginner journey ============ -->
    <section
      class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6"
      aria-labelledby="beginner-title"
    >
      <p class="text-xs font-semibold uppercase tracking-wide text-primary">Built for the beginner journey</p>
      <h2 id="beginner-title" class="mt-1 max-w-2xl text-balance text-2xl font-bold tracking-tight md:text-3xl">
        Your first goal should be understanding.
      </h2>
      <p class="mt-3 max-w-2xl text-pretty text-base leading-relaxed text-muted">
        Ctemz Journal isn't designed around the idea that you should immediately start risking real
        money. These are the steps, in order.
      </p>

      <ol class="mt-8 grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
        <li
          v-for="(p, i) in path"
          :key="p"
          class="flex items-start gap-3 rounded-xl border border-border bg-surface px-4 py-3"
        >
          <span
            class="tnum mt-0.5 grid h-6 w-6 shrink-0 place-items-center rounded-full bg-bg text-xs font-bold text-muted"
          >
            {{ i + 1 }}
          </span>
          <span class="text-sm font-medium leading-snug">{{ p }}</span>
        </li>
      </ol>

      <div class="mt-8 max-w-2xl space-y-3 border-l-2 border-primary pl-5 text-pretty leading-relaxed">
        <p class="text-muted">
          Then, when you have enough understanding and experience, decide whether you're ready to trade
          real money.
        </p>
        <p class="font-semibold">
          The goal isn't to rush you into the market. The goal is to help you become prepared for it.
        </p>
      </div>
    </section>

    <!-- ============ 4. Your journal is part of the learning process ============ -->
    <section
      class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6"
      aria-labelledby="journal-title"
    >
      <p class="text-xs font-semibold uppercase tracking-wide text-primary">
        Your journal is part of the learning process
      </p>
      <h2 id="journal-title" class="mt-1 max-w-2xl text-balance text-2xl font-bold tracking-tight md:text-3xl">
        A journal shouldn't just tell you whether you made or lost money. It should help you understand
        why.
      </h2>

      <div class="mt-8 grid gap-4 lg:grid-cols-2 lg:gap-6">
        <div class="rounded-3xl border border-border bg-surface p-5 md:p-6">
          <h3 class="text-xs font-semibold uppercase tracking-wide text-muted">What you record</h3>
          <ul class="mt-4 flex flex-wrap gap-2">
            <li
              v-for="f in entryFields"
              :key="f"
              class="rounded-full border border-border bg-bg px-3 py-1.5 text-sm font-medium"
            >
              {{ f }}
            </li>
          </ul>
          <p class="mt-5 text-sm leading-relaxed text-muted">
            Over time, your journal becomes a record of your development — not a memory test.
          </p>
        </div>

        <div class="rounded-3xl border border-border bg-surface p-5 md:p-6">
          <h3 class="text-xs font-semibold uppercase tracking-wide text-muted">What you get back</h3>
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

      <p class="mt-6 text-base font-semibold">Instead of relying on memory or feelings, you have evidence.</p>
    </section>

    <!-- ============ 5. Don't confuse winning with good trading ============ -->
    <section
      class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6"
      aria-labelledby="process-title"
    >
      <p class="text-xs font-semibold uppercase tracking-wide text-primary">
        Don't confuse winning with good trading
      </p>
      <h2 id="process-title" class="mt-1 max-w-2xl text-balance text-2xl font-bold tracking-tight md:text-3xl">
        The result of a trade is not the same thing as the quality of the decision.
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
              <span v-if="o.winning" aria-hidden="true">↑ </span>{{ o.result }}
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
        You can follow your strategy perfectly and lose. You can break your rules and still win. That's
        why Ctemz Journal focuses on the process behind every trade, not just the final result — the
        objective is a repeatable process you can evaluate and improve over time.
      </p>
    </section>

    <!-- ============ 6. Risk comes before profit ============ -->
    <section
      class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6"
      aria-labelledby="risk-title"
    >
      <div
        class="grid gap-8 rounded-3xl border border-border bg-surface p-6 md:p-10 lg:grid-cols-[minmax(0,1fr)_minmax(0,1.1fr)] lg:gap-14"
      >
        <div>
          <p class="text-xs font-semibold uppercase tracking-wide text-primary">Risk comes before profit</p>
          <h2 id="risk-title" class="mt-1 text-balance text-2xl font-bold tracking-tight md:text-3xl">
            The second question matters more than the first.
          </h2>
          <div class="mt-4 space-y-3 text-pretty leading-relaxed text-muted">
            <p>The question shouldn't always be: “How much can I make?”</p>
            <p>It should also be: “How much can I lose if I'm wrong?”</p>
            <p>
              Ctemz Journal helps you think about every trade in terms of risk, R multiple, reward and
              overall performance — because protecting your capital gives you something even more
              valuable: the ability to keep learning.
            </p>
          </div>
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
            <span
              class="absolute top-8.5 h-6 w-0.5 -translate-x-1/2 bg-text"
              :style="{ left: riskShare + '%' }"
              aria-hidden="true"
            />

            <!-- Distances below -->
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
              <span aria-hidden="true">↑</span> Win +{{ rMultiple }}R
            </span>
          </figcaption>
        </figure>
      </div>
    </section>

    <!-- ============ 7. Backtest before you believe ============ -->
    <section
      class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6"
      aria-labelledby="backtest-title"
    >
      <p class="text-xs font-semibold uppercase tracking-wide text-primary">Backtest before you believe</p>
      <h2 id="backtest-title" class="mt-1 max-w-2xl text-balance text-2xl font-bold tracking-tight md:text-3xl">
        A strategy can sound perfect when someone explains it online. That doesn't mean it works for you.
      </h2>

      <div class="mt-8 grid items-stretch gap-4 md:grid-cols-[minmax(0,1fr)_auto_minmax(0,1fr)] md:gap-6">
        <div class="rounded-3xl border border-border bg-bg p-5 md:p-6">
          <p class="text-xs font-semibold uppercase tracking-wide text-muted">Before</p>
          <p class="mt-3 text-lg font-semibold italic text-muted">“I think this strategy works.”</p>
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
            “What happened when these rules were actually followed?”
          </p>
        </div>
      </div>

      <p class="mt-6 max-w-2xl text-pretty leading-relaxed text-muted">
        Backtesting gives you a way to take a strategy and test its rules against historical market data.
        It doesn't guarantee future results — but it helps you understand your strategy, identify its
        weaknesses, and build a stronger foundation before putting real money at risk.
      </p>
    </section>

    <!-- ============ 8. Your psychology matters ============ -->
    <section
      class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6"
      aria-labelledby="psychology-title"
    >
      <p class="text-xs font-semibold uppercase tracking-wide text-primary">Your psychology matters</p>
      <h2 id="psychology-title" class="mt-1 max-w-2xl text-balance text-2xl font-bold tracking-tight md:text-3xl">
        Trading isn't only about charts.
      </h2>
      <p class="mt-3 max-w-2xl leading-relaxed text-muted">
        Your decisions are affected by your emotions.
      </p>

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
        Ctemz Journal gives you a place to record those behaviours instead of ignoring them — because
        becoming a better trader also means becoming better at managing yourself.
      </p>
    </section>

    <!-- ============ 9. Ignore the hype ============ -->
    <section class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6" aria-labelledby="hype-title">
      <div class="rounded-3xl border border-border bg-surface p-6 md:p-10">
        <p class="text-xs font-semibold uppercase tracking-wide text-primary">Ignore the hype</p>
        <h2 id="hype-title" class="mt-1 max-w-2xl text-balance text-2xl font-bold tracking-tight md:text-3xl">
          Forex is surrounded by a lot of noise.
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
            <p>Ctemz Journal isn't built around that.</p>
            <p>We don't believe trading should be sold as a shortcut to wealth.</p>
            <p>
              We believe it should be approached as a skill that requires education, practice, discipline,
              risk management and continuous improvement.
            </p>
          </div>
          <p class="self-center text-balance text-3xl font-extrabold leading-tight tracking-tight md:text-4xl">
            No hype. Just the work.
          </p>
        </div>
      </div>
    </section>

    <!-- ============ 10. Private by default (interactive) ============ -->
    <section
      class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6"
      aria-labelledby="privacy-title"
    >
      <div class="grid gap-10 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)] lg:gap-14">
        <div>
          <p class="text-xs font-semibold uppercase tracking-wide text-primary">Private by default</p>
          <h2 id="privacy-title" class="mt-1 text-balance text-2xl font-bold tracking-tight md:text-3xl">
            Your journal is personal.
          </h2>
          <div class="mt-4 space-y-3 text-pretty leading-relaxed text-muted">
            <p>Your trades, notes, mistakes, performance and account information belong to you.</p>
            <p>
              Ctemz Journal is designed with privacy in mind, while giving you the option to share
              selected journal entries when you choose.
            </p>
            <p class="font-semibold text-text">Sharing should be intentional. Not automatic.</p>
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

          <!-- The chart is shared, so it is always visible -->
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

    <!-- ============ 11. A record of your journey ============ -->
    <section
      class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6"
      aria-labelledby="journey-title"
    >
      <div class="grid gap-10 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)] lg:gap-14">
        <div>
          <p class="text-xs font-semibold uppercase tracking-wide text-primary">
            A record of your journey
          </p>
          <h2 id="journey-title" class="mt-1 text-balance text-2xl font-bold tracking-tight md:text-3xl">
            Your first trade won't look like your hundredth.
          </h2>
          <div class="mt-4 space-y-3 text-pretty leading-relaxed text-muted">
            <p>
              That's part of learning — and it's the part most traders only remember through
              rose-tinted hindsight.
            </p>
            <p class="font-semibold text-text">
              Ctemz Journal exists so you can actually see that progression. Not just remember it.
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

    <!-- ============ 12. Closing ============ -->
    <section class="mx-auto mt-20 w-full max-w-6xl px-4 md:mt-28 md:px-6">
      <div class="rounded-3xl border border-border bg-surface p-6 text-center md:p-12">
        <h2
          class="mx-auto max-w-2xl text-balance text-3xl font-extrabold leading-tight tracking-tight md:text-4xl"
        >
          Trade less emotionally. Learn more deliberately.
        </h2>
        <p class="mx-auto mt-4 max-w-xl text-pretty leading-relaxed text-muted">
          You don't need to become a profitable trader overnight. You need to become a more informed
          trader than you were yesterday.
        </p>

        <ul class="mx-auto mt-8 flex max-w-3xl flex-wrap justify-center gap-2">
          <li
            v-for="s in closingSteps"
            :key="s"
            class="rounded-full border border-border bg-bg px-3.5 py-1.5 text-sm font-medium"
          >
            {{ s }}
          </li>
        </ul>

        <p class="mt-8 text-sm text-muted">Ctemz Journal is built for that journey.</p>
        <p class="mt-2 text-sm font-semibold uppercase tracking-wide text-primary">
          Learn. Practice. Journal. Improve.
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