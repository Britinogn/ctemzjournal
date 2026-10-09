<script setup lang="ts">
import {
  DiscordIcon,
  InstagramIcon,
  Link02Icon,
  Linkedin01Icon,
  TelegramIcon,
  TiktokIcon,
  TwitterIcon,
  YoutubeIcon,
} from '~/utils/icons'
import BrandLogo from '../BrandLogo.vue'

const props = withDefaults(defineProps<{
  siteName?: string
  tagline?: string | null
  contactEmail?: string | null
  footerText?: string | null
  riskDisclaimer?: string | null
  socialLinks?: Record<string, unknown>
}>(), {
  siteName: 'Ctemz Journal',
  tagline: null,
  contactEmail: null,
  footerText: null,
  riskDisclaimer: null,
  socialLinks: () => ({}),
})

const SOCIAL_ICONS: Record<string, unknown> = {
  x: TwitterIcon,
  twitter: TwitterIcon,
  instagram: InstagramIcon,
  youtube: YoutubeIcon,
  telegram: TelegramIcon,
  discord: DiscordIcon,
  linkedin: Linkedin01Icon,
  tiktok: TiktokIcon,
}

const socials = computed(() =>
  Object.entries(props.socialLinks)
    .filter((entry): entry is [string, string] => typeof entry[1] === 'string' && entry[1] !== '')
    .map(([platform, url]) => ({
      platform,
      url,
      icon: SOCIAL_ICONS[platform.toLowerCase()] ?? Link02Icon,
    })),
)

const exploreLinks = [
  { to: '/journals', label: 'Public journals' },
  { to: '/prices', label: 'Live prices' },
  { to: '/about', label: 'About' },
]

const accountLinks = [
  { to: '/auth/login', label: 'Log in' },
  { to: '/auth/signup', label: 'Sign up' },
  { to: '/dashboard', label: 'Dashboard' },
]
</script>

<template>
  <footer class="border-t border-border bg-surface">
    <div class="mx-auto w-full max-w-6xl px-5 pb-[max(1.5rem,env(safe-area-inset-bottom))] pt-10 sm:px-6 md:pt-12">
      <!-- Top: brand + link columns -->
      <div class="grid grid-cols-2 gap-x-6 gap-y-10 md:grid-cols-[1.5fr_1fr_1fr] md:gap-16">
        <!-- Brand (full width on mobile) -->
        <div class="col-span-2 md:col-span-1">
          <BrandLogo :site-name="siteName" :size="150" />

          <p v-if="tagline" class="mt-4 max-w-sm text-sm leading-6 text-muted">
            {{ tagline }}
          </p>

          <a
            v-if="contactEmail"
            :href="`mailto:${contactEmail}`"
            class="mt-4 inline-flex min-h-[44px] items-center break-all text-sm font-medium text-text transition-colors hover:text-primary md:min-h-0"
          >
            {{ contactEmail }}
          </a>

          <div v-if="socials.length > 0" class="mt-3 flex flex-wrap items-center gap-2.5 md:mt-5">
            <a
              v-for="s in socials"
              :key="s.platform"
              :href="s.url"
              target="_blank"
              rel="noopener noreferrer"
              :aria-label="s.platform"
              class="inline-flex h-11 w-11 items-center justify-center rounded-xl border border-border bg-bg text-muted transition-all hover:border-primary/40 hover:bg-primary/5 hover:text-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary md:h-9 md:w-9 md:rounded-lg"
            >
              <UiAppIcon :icon="s.icon" :size="18" aria-hidden="true" />
            </a>
          </div>
        </div>

        <!-- Explore -->
        <nav aria-label="Explore">
          <p class="text-xs font-semibold uppercase tracking-wider text-muted">
            Explore
          </p>
          <ul class="mt-3 text-sm">
            <li v-for="l in exploreLinks" :key="l.to">
              <NuxtLink
                :to="l.to"
                class="block py-2.5 text-muted transition-colors hover:text-text md:py-1.5"
              >
                {{ l.label }}
              </NuxtLink>
            </li>
          </ul>
        </nav>

        <!-- Account -->
        <nav aria-label="Account">
          <p class="text-xs font-semibold uppercase tracking-wider text-muted">
            Account
          </p>
          <ul class="mt-3 text-sm">
            <li v-for="l in accountLinks" :key="l.to">
              <NuxtLink
                :to="l.to"
                class="block py-2.5 text-muted transition-colors hover:text-text md:py-1.5"
              >
                {{ l.label }}
              </NuxtLink>
            </li>
          </ul>
        </nav>
      </div>

      <!-- Footer text + risk disclaimer -->
      <div
        v-if="footerText || riskDisclaimer"
        class="mt-10 space-y-3 border-t border-border pt-6 md:mt-12"
      >
        <p v-if="footerText" class="max-w-3xl text-xs leading-5 text-muted">
          {{ footerText }}
        </p>

        <p
          v-if="riskDisclaimer"
          class="max-w-4xl rounded-xl bg-bg p-3.5 text-xs leading-5 text-muted/90 md:rounded-none md:bg-transparent md:p-0"
        >
          <span class="font-medium text-muted">Risk disclaimer:</span>
          {{ riskDisclaimer }}
        </p>
      </div>

      <!-- Bottom bar -->
      <div class="mt-8 flex flex-col items-center gap-3 border-t border-border pt-6 text-center sm:flex-row sm:justify-between sm:text-left">
        <div class="flex items-center gap-6 text-xs text-muted sm:order-last">
          <NuxtLink
            to="/privacy"
            target="_blank"
            rel="noopener noreferrer"
            class="py-2 transition-colors hover:text-text"
          >
            Privacy
          </NuxtLink>

          <NuxtLink
            to="/terms"
            target="_blank"
            rel="noopener noreferrer"
            class="py-2 transition-colors hover:text-text"
          >
            Terms
          </NuxtLink>
        </div>

        <p class="text-xs text-muted">
          © {{ new Date().getFullYear() }} {{ siteName }}. All rights reserved.
        </p>
      </div>
    </div>
  </footer>
</template>