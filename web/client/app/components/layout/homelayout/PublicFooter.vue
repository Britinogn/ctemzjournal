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
</script>

<template>
  <footer class="border-t border-border bg-surface">
    <div class="mx-auto w-full max-w-6xl px-4 py-12 md:px-6">
      <div class="grid gap-10 md:grid-cols-[1.5fr_1fr_1fr] md:gap-16">
        <div>
          <BrandLogo :site-name="siteName" :size="150" />

          <p v-if="tagline" class="mt-4 max-w-sm text-sm leading-6 text-muted">
            {{ tagline }}
          </p>

          <a
            v-if="contactEmail"
            :href="`mailto:${contactEmail}`"
            class="mt-4 inline-flex text-sm font-medium text-text transition-colors hover:text-primary"
          >
            {{ contactEmail }}
          </a>

          <div v-if="socials.length > 0" class="mt-5 flex items-center gap-2">
            <a
              v-for="s in socials"
              :key="s.platform"
              :href="s.url"
              target="_blank"
              rel="noopener noreferrer"
              :aria-label="s.platform"
              class="inline-flex h-9 w-9 items-center justify-center rounded-lg border border-border bg-background text-muted transition-all hover:border-primary/40 hover:bg-primary/5 hover:text-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
            >
              <UiAppIcon :icon="s.icon" :size="17" aria-hidden="true" />
            </a>
          </div>
        </div>

        <nav aria-label="Explore">
          <p class="text-xs font-semibold uppercase tracking-wider text-muted">
            Explore
          </p>

          <ul class="mt-4 space-y-3 text-sm">
            <li>
              <NuxtLink to="/journals" class="text-muted transition-colors hover:text-text">
                Public journals
              </NuxtLink>
            </li>
            <li>
              <NuxtLink to="/prices" class="text-muted transition-colors hover:text-text">
                Live prices
              </NuxtLink>
            </li>
            <li>
              <NuxtLink to="/about" class="text-muted transition-colors hover:text-text">
                About
              </NuxtLink>
            </li>
          </ul>
        </nav>

        <nav aria-label="Account">
          <p class="text-xs font-semibold uppercase tracking-wider text-muted">
            Account
          </p>

          <ul class="mt-4 space-y-3 text-sm">
            <li>
              <NuxtLink to="/auth/login" class="text-muted transition-colors hover:text-text">
                Log in
              </NuxtLink>
            </li>
            <li>
              <NuxtLink to="/auth/signup" class="text-muted transition-colors hover:text-text">
                Sign up
              </NuxtLink>
            </li>
            <li>
              <NuxtLink to="/dashboard" class="text-muted transition-colors hover:text-text">
                Dashboard
              </NuxtLink>
            </li>
          </ul>
        </nav>
      </div>

      <div
        v-if="footerText || riskDisclaimer"
        class="mt-12 border-t border-border pt-6"
      >
        <p v-if="footerText" class="max-w-3xl text-xs leading-5 text-muted">
          {{ footerText }}
        </p>

        <p
          v-if="riskDisclaimer"
          class="mt-3 max-w-4xl text-xs leading-5 text-muted/80"
        >
          <span class="font-medium text-muted">Risk disclaimer:</span>
          {{ riskDisclaimer }}
        </p>
      </div>

      <div class="mt-8 flex flex-col gap-3 border-t border-border pt-6 sm:flex-row sm:items-center sm:justify-between">
        <p class="text-xs text-muted">
          © {{ new Date().getFullYear() }} {{ siteName }}. All rights reserved.
        </p>

        <div class="flex items-center gap-4 text-xs text-muted">
          <NuxtLink to="/privacy" class="transition-colors hover:text-text">
            Privacy
          </NuxtLink>

          <NuxtLink to="/terms" class="transition-colors hover:text-text">
            Terms
          </NuxtLink>
        </div>
      </div>
    </div>
  </footer>
</template>