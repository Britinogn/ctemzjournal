<script setup lang="ts">
type Step = {
  n: number
  title: string
  paragraphs: string[]
  list?: string[]
}

// Twelve sections straight from the How it works doc. Single open at a time.
const steps: Step[] = [
  {
    n: 1,
    title: 'Learn',
    paragraphs: [
      'Your learning starts outside Ctemz Journal.',
      'Study market structure. Learn technical analysis. Understand risk management. Study trading psychology. Learn how different strategies work. Read books, take courses, learn from mentors, study charts, watch the market, and pay attention to your own experiences.',
      'There is no single source that makes you a trader. Ctemz Journal gives you a place to document what you take from those experiences.',
      'Record the concepts you are learning, the ideas you want to test, the rules you are developing, and the lessons that matter to you.',
      'The goal is not to collect information. The goal is to turn what you learn into something you can practice and eventually evaluate.',
    ],
  },
  {
    n: 2,
    title: 'Build your trading plan',
    paragraphs: [
      'Learning something does not mean you should immediately trade it with real money.',
      'Before taking a setup to the market, understand what you are actually looking for. Define the parts of your process so your decisions are not based entirely on how you feel when the market is moving.',
    ],
    list: [
      'Your strategy',
      'What conditions need to exist before you enter',
      'Where your entry comes from',
      'Where your stop loss belongs',
      'Your target',
      'How much you are willing to risk',
      'What would invalidate your idea',
    ],
  },
  {
    n: 3,
    title: 'Practice and test',
    paragraphs: [
      'An idea can sound perfect until you put it against real market data. Test it.',
      'Backtest historical charts. Practice your entries. Study different market conditions. Look for situations where your strategy works and situations where it does not.',
      'Do not assume a strategy works because someone online says it does. Build evidence for yourself.',
      'Ctemz Journal can become the record of what you tested, what you observed, and what you learned from the process. That information is valuable.',
    ],
  },
  {
    n: 4,
    title: 'Journal your trades',
    paragraphs: [
      'When you take a trade, record it. Do not only record whether you won or lost.',
      'A trade without context only tells you the outcome. A properly documented trade can tell you why the outcome happened.',
    ],
    list: [
      'The pair',
      'Your setup',
      'Your entry',
      'Your stop loss and take profit',
      'Your position size and risk',
      'Your chart screenshots',
      'Why you entered',
      'What you were thinking',
      'Whether the trade followed your rules',
      'What happened',
      'What you learned',
    ],
  },
  {
    n: 5,
    title: 'Record the decision, not just the result',
    paragraphs: [
      'This is one of the most important parts of journaling.',
      'You can take a trade that follows your strategy perfectly and still lose. You can break every rule and still make money. The result alone does not tell you whether the decision was good.',
      'That is why Ctemz Journal focuses on the process behind your trades. Your journal gives you somewhere to be honest about those decisions.',
    ],
    list: [
      'Did you wait for your setup?',
      'Did you respect your risk?',
      'Did you move your stop because you were afraid?',
      'Did you enter because of FOMO?',
      'Did you close early because you became nervous?',
      'Did you increase your position because you wanted to recover a previous loss?',
    ],
  },
  {
    n: 6,
    title: 'Track your behaviour',
    paragraphs: [
      'Trading is not only about market analysis. It is also about how you behave when your money is at risk.',
      'These behaviours can become patterns. When you consistently record them, you can begin to see those patterns instead of treating every mistake as an isolated event.',
      'Your journal gives you evidence about yourself.',
    ],
    list: [
      'Fear',
      'Greed',
      'FOMO',
      'Revenge trading',
      'Overconfidence',
      'Impatience',
      'Hesitation',
    ],
  },
  {
    n: 7,
    title: 'Review your trades',
    paragraphs: [
      'Journaling is not useful if you never look back at what you recorded. Review is where the record starts becoming useful.',
      'Ask yourself what keeps working and what keeps going wrong.',
      'The purpose is not to judge yourself. The purpose is to understand yourself.',
    ],
    list: [
      'Look through your previous trades',
      'Compare your setups',
      'Look at your wins and losses',
      'Look at the trades where you followed your rules',
      'Look at the trades where you did not',
      'Look for repeated mistakes',
      'Look for patterns in your behaviour',
      'Look at your risk',
      'Look at your performance',
    ],
  },
  {
    n: 8,
    title: 'Measure your performance',
    paragraphs: [
      'Your trading results should eventually tell you more than how much money you made.',
      'Numbers can reveal things that emotions hide. You may believe that one setup is your best because you remember a few large wins. Your records might show something completely different.',
      'That is why evidence matters.',
    ],
    list: [
      'Your win rate',
      'Your expectancy',
      'Your drawdown',
      'Your results in R',
      'Your performance by setup',
      'Your performance by pair',
      'How consistently you follow your rules',
    ],
  },
  {
    n: 9,
    title: 'Learn from the record',
    paragraphs: [
      'After enough trades and enough honest reviews, your journal starts becoming a source of personal knowledge.',
      'You begin to understand which setups suit you. You understand where your strategy struggles. You recognize your common mistakes. You understand how you behave after losses. You see how your risk affects your results.',
      'You begin to recognize the conditions where you trade well and the conditions where you should probably stay out.',
      'This is where journaling becomes more than record keeping. Your past decisions become something you can learn from.',
    ],
  },
  {
    n: 10,
    title: 'Improve your process',
    paragraphs: [
      'The goal is not to make every trade a winner. The goal is to make your process better.',
      'Your journal helps you make these decisions using your own trading history instead of assumptions.',
    ],
    list: [
      'Maybe you need to reduce your risk',
      'Maybe you need stricter entry conditions',
      'Maybe you need to stop trading a particular setup',
      'Maybe you need more backtesting',
      'Maybe your strategy is fine but your execution needs work',
      'Maybe your biggest problem is not your strategy at all. Maybe it is your discipline.',
    ],
  },
  {
    n: 11,
    title: 'Build your trading history',
    paragraphs: [
      'Over time, your journal becomes a record of your development. Your first trades will not look like your hundredth.',
      'Instead of relying on memory, you have a record of what actually happened. You can look back and see where you started, what you struggled with, what you learned, and how your process developed.',
    ],
    list: [
      'Your understanding will change',
      'Your strategy will change',
      'Your risk management will improve',
      'Your mistakes will change',
      'Your confidence will change',
      'There will be good periods and bad periods',
    ],
  },
  {
    n: 12,
    title: 'Decide what you want to share',
    paragraphs: [
      'Your trading journal is yours.',
      'You may want to keep everything private while you learn. You may eventually want to share a setup, a lesson, a chart, or part of your journey with other traders.',
      'Ctemz Journal lets you decide what becomes public. Sharing your journey should not require exposing everything about your finances. Your private trading information can remain private while the parts of your journey you choose to share can become visible to others.',
    ],
  },
]

const openStep = ref(1)
function toggle(n: number) {
  openStep.value = openStep.value === n ? 0 : n
}

const loop = [
  'Learn',
  'Practice',
  'Test',
  'Trade',
  'Journal',
  'Review',
  'Improve',
  'Repeat',
]
</script>

<template>
  <section
    class="mx-auto w-full max-w-5xl px-4 pt-16 md:px-6 md:pt-24"
    aria-labelledby="how-title"
  >
    <!-- Header -->
    <header class="mx-auto max-w-2xl text-center">
      <p class="text-xs font-semibold uppercase tracking-wide text-primary">
        How it works
      </p>
      <h2
        id="how-title"
        class="mt-1 text-balance text-2xl font-bold tracking-tight md:text-3xl"
      >
        Your trading journey, documented.
      </h2>
      <p class="mt-3 text-pretty text-sm leading-relaxed text-muted md:text-base">
        You learn from the market and from the people, resources, and experiences around you.
        Then you journal what you learn. The journal becomes your record of that process.
      </p>
    </header>

    <!-- Steps -->
    <ol class="mt-10 space-y-2.5 md:mt-14">
      <li
        v-for="s in steps"
        :key="s.n"
        class="overflow-hidden rounded-2xl border border-border bg-surface"
      >
        <button
          type="button"
          class="flex w-full items-center gap-4 px-4 py-4 text-left transition-colors hover:bg-bg md:px-5 md:py-5"
          :aria-expanded="openStep === s.n"
          :aria-controls="`step-${s.n}`"
          @click="toggle(s.n)"
        >
          <span
            class="tnum grid h-8 w-8 shrink-0 place-items-center rounded-full text-sm font-bold"
            :class="
              openStep === s.n
                ? 'bg-primary text-on-primary'
                : 'bg-bg text-muted'
            "
          >
            {{ s.n }}
          </span>

          <span
            class="min-w-0 flex-1 text-base font-bold tracking-tight md:text-lg"
          >
            {{ s.title }}
          </span>

          <svg
            class="h-4 w-4 shrink-0 text-muted transition-transform duration-200 motion-reduce:transition-none"
            :class="openStep === s.n ? 'rotate-180' : ''"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            aria-hidden="true"
          >
            <path d="m6 9 6 6 6-6" />
          </svg>
        </button>

        <div
          :id="`step-${s.n}`"
          class="grid transition-[grid-template-rows] duration-300 ease-out motion-reduce:transition-none"
          :class="openStep === s.n ? 'grid-rows-[1fr]' : 'grid-rows-[0fr]'"
        >
          <div class="overflow-hidden">
            <div class="space-y-3 px-4 pb-5 pl-16 text-sm leading-relaxed text-muted md:px-5 md:pb-6 md:pl-17 md:text-[15px]">
              <p v-for="(p, i) in s.paragraphs" :key="i" class="text-pretty">
                {{ p }}
              </p>

              <ul
                v-if="s.list"
                class="space-y-1.5 pt-1"
              >
                <li
                  v-for="item in s.list"
                  :key="item"
                  class="flex items-start gap-2.5"
                >
                  <span
                    class="mt-2 h-1 w-1 shrink-0 rounded-full bg-primary"
                    aria-hidden="true"
                  />
                  <span class="text-pretty">{{ item }}</span>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </li>
    </ol>

    <!-- Loop -->
    <div
      class="mt-12 rounded-3xl border border-border bg-surface p-6 text-center md:mt-16 md:p-10"
    >
      <p class="text-xs font-semibold uppercase tracking-wide text-primary">
        The loop
      </p>

      <ul class="mt-4 flex flex-wrap items-center justify-center gap-x-2 gap-y-2">
        <template v-for="(word, i) in loop" :key="word">
          <li
            class="rounded-full border border-border bg-bg px-3.5 py-1.5 text-sm font-semibold"
            :class="word === 'Repeat' ? 'text-primary' : ''"
          >
            {{ word }}
          </li>
          <li
            v-if="i < loop.length - 1"
            class="text-muted"
            aria-hidden="true"
          >
            →
          </li>
        </template>
      </ul>

      <div
        class="mx-auto mt-6 max-w-2xl space-y-3 text-pretty text-sm leading-relaxed text-muted md:text-base"
      >
        <p>
          You learn something. You practice it. You test whether it actually works. You trade it.
          You record what happened. You review the evidence. You improve your process. Then you
          repeat the cycle.
        </p>
        <p>That process never really ends.</p>
      </div>

      <p
        class="mx-auto mt-6 max-w-2xl text-balance text-base font-semibold md:text-lg"
      >
        Ctemz Journal is not the source of your trading knowledge. It is the place where your
        knowledge, decisions, experiences, mistakes, and progress come together.
      </p>

      <ul
        class="mx-auto mt-6 flex max-w-2xl flex-wrap justify-center gap-2 text-sm"
      >
        <li
          v-for="line in ['You learn from wherever you choose', 'You journal what you learn', 'You review what you do', 'You improve from what you discover']"
          :key="line"
          class="rounded-full border border-primary/30 bg-primary/5 px-3.5 py-1.5 font-medium text-primary"
        >
          {{ line }}
        </li>
      </ul>

      <p class="mt-6 text-sm font-semibold uppercase tracking-wide text-primary">
        That is how Ctemz Journal works.
      </p>
    </div>
  </section>
</template>