<!-- pages/terms.vue -->
<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { siteSettingsKey, type PublicSettings } from '~/types'

const LAST_UPDATED = 'October 9, 2026'

useHead({ title: 'Terms of Service | Ctemz Journal' })
useSeoMeta({
  description:
    'The terms that apply when you use Ctemz Journal, a trading journal for recording and reviewing your trades.',
})

const api = useApi()
const { data: settings } = useQuery({
  queryKey: siteSettingsKey(),
  queryFn: () => api.get<PublicSettings>('/public/site-settings'),
  staleTime: 30 * 60_000,
})

interface Section {
  id: string
  title: string
  paragraphs: string[]
}

const sections: Section[] = [
  {
    id: 'acceptance',
    title: 'Acceptance of these terms',
    paragraphs: [
      'By accessing or using Ctemz Journal, you agree to these Terms of Service. If you do not agree to these terms, please do not use the service.',
    ],
  },
  {
    id: 'about',
    title: 'About Ctemz Journal',
    paragraphs: [
      'Ctemz Journal is a trading journal designed to help you record trades, document your setups, track your risk, review your decisions and analyse your trading performance. It is a record keeping and learning tool, not a trading platform or brokerage service.',
    ],
  },
  {
    id: 'eligibility',
    title: 'Eligibility',
    paragraphs: [
      'You must be at least 18 years old to create an account and use the service.',
    ],
  },
  {
    id: 'your-account',
    title: 'Your account',
    paragraphs: [
      'You are responsible for maintaining the security of your account and for all activity carried out through it. You agree to provide accurate information and to tell us if you suspect unauthorised access to your account.',
      'You must not misuse the service or attempt to gain unauthorised access to another user\'s information.',
    ],
  },
  {
    id: 'your-data',
    title: 'Your trading data',
    paragraphs: [
      'You remain responsible for the trade records, notes, screenshots and other information you enter into Ctemz Journal. Make sure what you record is accurate and that you have the right to upload any content you submit.',
      'We recommend keeping your own backups of important records.',
    ],
  },
  {
    id: 'no-advice',
    title: 'No financial advice',
    paragraphs: [
      'Ctemz Journal provides journaling, record keeping and analytical tools for educational and informational purposes only. Nothing provided by the service is financial, investment or trading advice, or a recommendation to buy or sell any financial instrument.',
      'Trading forex, gold, cryptocurrencies and other financial instruments involves risk and may result in the loss of some or all of your capital. Past performance and journal statistics do not guarantee future results. You are solely responsible for your trading decisions and their outcomes.',
    ],
  },
  {
    id: 'acceptable-use',
    title: 'Acceptable use',
    paragraphs: [
      'You agree not to use Ctemz Journal for unlawful activities, interfere with its operation, attempt to bypass security measures, distribute malicious software, or access data that does not belong to you. We may restrict access where necessary to protect the service and its users.',
    ],
  },
  {
    id: 'public-sharing',
    title: 'Public sharing',
    paragraphs: [
      'If the service provides features for sharing journal entries or trading statistics, you are responsible for the content you choose to make public. Before sharing, review what will be displayed and avoid exposing sensitive personal or account information.',
    ],
  },
  {
    id: 'availability',
    title: 'Service availability',
    paragraphs: [
      'We aim to keep Ctemz Journal available and reliable, but we do not guarantee uninterrupted or error free operation. The service may be updated, modified, temporarily suspended or discontinued for maintenance, security, technical or business reasons.',
    ],
  },
  {
    id: 'termination',
    title: 'Suspension and termination',
    paragraphs: [
      'We may suspend or terminate access to the service if you violate these terms, misuse the application, threaten its security, or engage in activities that may harm the service or other users. Where reasonably possible, we will notify affected users.',
    ],
  },
  {
    id: 'liability',
    title: 'Disclaimer and limitation of liability',
    paragraphs: [
      'To the extent permitted by applicable law, Ctemz Journal is provided on an "as is" and "as available" basis. We do not guarantee that calculations, statistics, records or other information will always be accurate, complete or available. Please verify important information independently before relying on it.',
      'To the extent permitted by applicable law, Ctemz Journal and its operators are not liable for trading losses, decisions made using information recorded or displayed by the service, loss of data, or indirect damages arising from your use of the service. Nothing in these terms excludes liability that cannot legally be excluded.',
    ],
  },
  {
    id: 'changes',
    title: 'Changes to these terms',
    paragraphs: [
      'We may update these Terms of Service as Ctemz Journal evolves. Updated terms will be published on this page with a revised date. Your continued use of the service after changes take effect means you accept the updated terms, where permitted by applicable law.',
    ],
  },
  {
    id: 'governing-law',
    title: 'Governing law',
    paragraphs: [
      'These terms are subject to applicable laws and regulations. Any disputes will be handled in accordance with the laws and jurisdiction applicable to the service and the parties involved.',
    ],
  },
]

const jumpItems = computed(() => [
  ...sections.map(s => ({ id: s.id, title: s.title })),
  { id: 'contact', title: 'Contact' },
])

const glance = [
  { title: 'A journal, not a broker', body: 'We record and analyse. We do not trade for you.' },
  { title: 'Not financial advice', body: 'Trading carries risk. Decisions are yours.' },
  { title: 'Keep your account safe', body: 'You are responsible for activity on it.' },
]

const showTop = ref(false)

function prefersReducedMotion() {
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

function goTo(id: string, e?: Event) {
  e?.preventDefault()
  const el = document.getElementById(id)
  if (!el)
    return
  el.scrollIntoView({ behavior: prefersReducedMotion() ? 'auto' : 'smooth', block: 'start' })
  history.replaceState(null, '', `#${id}`)
}

function toTop() {
  window.scrollTo({ top: 0, behavior: prefersReducedMotion() ? 'auto' : 'smooth' })
}

function printPage() {
  window.print()
}

function onScroll() {
  showTop.value = window.scrollY > 700
}

onMounted(() => {
  window.addEventListener('scroll', onScroll, { passive: true })
  onScroll()
})
onBeforeUnmount(() => window.removeEventListener('scroll', onScroll))
</script>

<template>
  <main class="mx-auto w-full max-w-3xl flex-1 px-4 py-12 sm:px-6 sm:py-16 md:py-20">
    <!-- Header -->
    <header class="text-center">
      <span class="inline-flex items-center rounded-full bg-primary/10 px-3 py-1 text-xs font-medium text-primary">
        Legal
      </span>
      <h1 class="mt-4 text-3xl font-bold tracking-tight sm:text-4xl md:text-5xl">
        Terms of Service
      </h1>
      <p class="mx-auto mt-4 max-w-xl text-sm leading-relaxed text-muted sm:text-base">
        The rules that apply when you use Ctemz Journal.
      </p>
      <div class="mt-4 flex flex-wrap items-center justify-center gap-x-4 gap-y-2 text-xs text-muted">
        <span>Last updated: {{ LAST_UPDATED }}</span>
        <span aria-hidden="true">·</span>
        <button
          type="button"
          class="font-medium text-primary hover:underline print:hidden"
          @click="printPage"
        >
          Print or save as PDF
        </button>
      </div>
    </header>

    <!-- At a glance -->
    <div class="mt-10 grid gap-3 sm:grid-cols-3">
      <div
        v-for="g in glance"
        :key="g.title"
        class="rounded-xl border border-border bg-surface p-4 text-center sm:text-left"
      >
        <p class="text-sm font-semibold">
          {{ g.title }}
        </p>
        <p class="mt-1 text-xs leading-relaxed text-muted">
          {{ g.body }}
        </p>
      </div>
    </div>

    <!-- Jump to -->
    <nav class="mt-10 rounded-2xl border border-border bg-surface p-5 print:hidden sm:p-6" aria-label="Jump to section">
      <p class="mb-3 text-xs font-semibold uppercase tracking-wider text-muted">
        Jump to
      </p>
      <ol class="grid gap-x-6 gap-y-1 sm:grid-cols-2">
        <li v-for="(item, i) in jumpItems" :key="item.id">
          <a
            :href="`#${item.id}`"
            class="flex items-baseline gap-2 rounded-lg px-2 py-1.5 text-sm text-muted transition hover:bg-bg hover:text-primary"
            @click="goTo(item.id, $event)"
          >
            <span class="w-5 shrink-0 text-xs font-semibold text-primary">{{ i + 1 }}</span>
            <span>{{ item.title }}</span>
          </a>
        </li>
      </ol>
    </nav>

    <!-- Sections -->
    <div class="mt-10 space-y-4">
      <section
        v-for="(s, i) in sections"
        :id="s.id"
        :key="s.id"
        class="scroll-mt-24 rounded-2xl border border-border bg-surface p-5 sm:p-7"
      >
        <div class="flex items-start gap-4">
          <span class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-sm font-bold text-primary">
            {{ i + 1 }}
          </span>
          <div class="min-w-0 flex-1">
            <h2 class="text-lg font-semibold tracking-tight sm:text-xl">
              {{ s.title }}
            </h2>
            <div class="mt-3 space-y-3 text-sm leading-relaxed text-muted sm:text-[15px]">
              <p v-for="(p, pi) in s.paragraphs" :key="pi">
                {{ p }}
              </p>
            </div>
          </div>
        </div>
      </section>

      <!-- Contact -->
      <section
        id="contact"
        class="scroll-mt-24 rounded-2xl border border-primary/30 bg-primary/5 p-5 text-center sm:p-8"
      >
        <h2 class="text-lg font-semibold tracking-tight sm:text-xl">
          Questions about these terms?
        </h2>
        <p class="mx-auto mt-2 max-w-md text-sm leading-relaxed text-muted">
          If anything here is unclear, get in touch and we will help.
        </p>
        <a
          v-if="settings?.contact_email"
          :href="`mailto:${settings.contact_email}?subject=Terms%20question`"
          class="mt-5 inline-flex items-center justify-center rounded-xl bg-primary px-6 py-3 text-sm font-semibold text-on-primary transition hover:opacity-90"
        >
          {{ settings.contact_email }}
        </a>
        <p v-else class="mt-5 text-sm text-muted">
          Use the contact details listed in the site footer.
        </p>
      </section>
    </div>

    <p class="mt-8 text-center text-xs text-muted">
      See also our
      <NuxtLink to="/privacy" class="font-medium text-primary hover:underline">Privacy Policy</NuxtLink>.
    </p>

    <!-- Back to top -->
    <Transition
      enter-active-class="transition duration-200"
      enter-from-class="translate-y-2 opacity-0"
      leave-active-class="transition duration-150"
      leave-to-class="translate-y-2 opacity-0"
    >
      <button
        v-if="showTop"
        type="button"
        aria-label="Back to top"
        class="fixed bottom-6 right-6 z-40 flex h-11 w-11 items-center justify-center rounded-full border border-border bg-surface text-lg shadow-md transition hover:bg-bg print:hidden"
        @click="toTop"
      >
        ↑
      </button>
    </Transition>
  </main>
</template>