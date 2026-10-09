<!-- pages/privacy.vue -->
<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { siteSettingsKey, type PublicSettings } from '~/types'

const LAST_UPDATED = 'October 9, 2026'

useHead({ title: 'Privacy Policy | Ctemz Journal' })
useSeoMeta({
  description:
    'How Ctemz Journal collects, uses and protects your information and trading journal data.',
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
  bullets?: string[]
}

const sections: Section[] = [
  {
    id: 'information-we-collect',
    title: 'Information we collect',
    paragraphs: [
      'Ctemz Journal collects the information needed to provide and maintain your account and trading journal.',
    ],
    bullets: [
      'Account details: your name, email address and authentication information.',
      'Journal content you choose to record: currency pairs, setups, entry and exit prices, stop loss, take profit, risk, results, screenshots, emotions, mistakes and notes.',
      'Basic technical data needed to keep your session signed in and the service secure.',
    ],
  },
  {
    id: 'how-we-use-it',
    title: 'How we use your information',
    paragraphs: [
      'We use your information to operate your journal, calculate your trading statistics, help you review your decisions and performance, keep your account secure, and improve the reliability of the service.',
    ],
  },
  {
    id: 'cookies-and-sessions',
    title: 'Cookies and sessions',
    paragraphs: [
      'We use browser storage and cookies that are necessary to keep you signed in and to remember basic preferences such as your theme. We do not use them to sell your data.',
    ],
  },
  {
    id: 'storage-and-security',
    title: 'Storage and security',
    paragraphs: [
      'Your information is stored with the infrastructure providers used to run Ctemz Journal. We apply access controls and reasonable safeguards to help protect it. No online service can guarantee absolute security, so please use a strong, unique password.',
    ],
  },
  {
    id: 'sharing',
    title: 'Sharing of information',
    paragraphs: [
      'We do not sell your personal information.',
      'Service providers may process information on our behalf where necessary to run the app, such as authentication, database hosting, file storage and web hosting. We may also disclose information where the law requires it or where necessary to protect the service and its users.',
    ],
  },
  {
    id: 'journal-privacy',
    title: 'Trading journal privacy',
    paragraphs: [
      'Your trade records, notes, screenshots and statistics are intended for your personal use. If you use a feature that shares journal information publicly, what you include in that shared view may be seen by other people. Do not put sensitive personal information in public notes or shared content.',
    ],
  },
  {
    id: 'your-choices',
    title: 'Your choices and data requests',
    paragraphs: [
      'You are responsible for the information you enter into your journal. You can ask us to:',
    ],
    bullets: [
      'Access the personal information we hold about you.',
      'Correct information that is inaccurate.',
      'Export a copy of your data.',
      'Delete your account and associated data.',
    ],
  },
  {
    id: 'retention',
    title: 'Data retention',
    paragraphs: [
      'We keep your information for as long as needed to provide the service, maintain your account, meet legal obligations and resolve disputes. Some information may remain in backups for a limited period after deletion.',
    ],
  },
  {
    id: 'age',
    title: 'Eligibility',
    paragraphs: [
      'Ctemz Journal is intended for adults. You must be at least 18 years old to create an account.',
    ],
  },
  {
    id: 'changes',
    title: 'Changes to this policy',
    paragraphs: [
      'We may update this Privacy Policy as Ctemz Journal evolves. When we do, we will update the date at the top of this page.',
    ],
  },
]

const jumpItems = computed(() => [
  ...sections.map(s => ({ id: s.id, title: s.title })),
  { id: 'contact', title: 'Contact' },
])

const glance = [
  { title: 'We never sell your data', body: 'Your journal is yours.' },
  { title: 'You stay in control', body: 'Ask us to access, export or delete it.' },
  { title: 'Only what is needed', body: 'We collect what the journal needs to work.' },
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
        Privacy Policy
      </h1>
      <p class="mx-auto mt-4 max-w-xl text-sm leading-relaxed text-muted sm:text-base">
        How Ctemz Journal collects, uses and protects your information.
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
              <ul v-if="s.bullets" class="space-y-2">
                <li v-for="(b, bi) in s.bullets" :key="bi" class="flex items-start gap-2.5">
                  <span class="mt-2 h-1.5 w-1.5 shrink-0 rounded-full bg-primary" />
                  <span>{{ b }}</span>
                </li>
              </ul>
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
          Questions about your data?
        </h2>
        <p class="mx-auto mt-2 max-w-md text-sm leading-relaxed text-muted">
          To access, correct, export or delete your information, or to ask anything about this policy, get in touch.
        </p>
        <a
          v-if="settings?.contact_email"
          :href="`mailto:${settings.contact_email}?subject=Privacy%20request`"
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
      This policy describes how the service works today and may change as the product evolves.
    </p>

    <p class="mt-8 text-center text-xs text-muted">
      See also our
      <NuxtLink to="/terms" class="font-medium text-primary hover:underline">Terms of Use</NuxtLink>.
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