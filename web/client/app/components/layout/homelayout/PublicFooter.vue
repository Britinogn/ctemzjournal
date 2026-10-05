<script setup lang="ts">
import {
  DiscordIcon,
  InstagramIcon,
  Link02Icon,
  TelegramIcon,
  TwitterIcon,
  YoutubeIcon,
} from '~/utils/icons'
import BrandLogo from '../BrandLogo.vue'

const props = withDefaults(defineProps<{
  siteName?: string;
  tagline?: string | null;
  footerText?: string | null;
  riskDisclaimer?: string | null;
  socialLinks?: Record<string, unknown>;
}>(), {
  siteName: 'Ctemz Journal',
})

const SOCIAL_ICONS: Record<string, unknown> = {
  x: TwitterIcon,
  twitter: TwitterIcon,
  instagram: InstagramIcon,
  youtube: YoutubeIcon,
  telegram: TelegramIcon,
  discord: DiscordIcon,
}

const socials = computed(() =>
  Object.entries(props.socialLinks ?? {})
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
    <div class="mx-auto w-full max-w-6xl px-4 py-10 md:px-6">
      <div class="grid gap-8 md:grid-cols-3">
        <div>
          <BrandLogo :site-name="siteName" :size="150" />
          <p v-if="tagline" class="mt-2 text-sm text-muted">
            {{ tagline }}
          </p>
          <div v-if="socials.length > 0" class="mt-3 flex gap-2">
            <a
              v-for="s in socials"
              :key="s.platform"
              :href="s.url"
              target="_blank"
              rel="noopener"
              :aria-label="s.platform"
              class="inline-flex h-9 w-9 items-center justify-center rounded-xl border border-border text-muted transition hover:border-primary hover:text-primary"
            >
              <UiAppIcon :icon="s.icon" :size="18" />
            </a>
          </div>
        </div>
        <nav aria-label="Explore">
          <p class="text-sm font-semibold">Explore</p>
          <ul class="mt-2 space-y-1.5 text-sm text-muted">
            <li><NuxtLink to="/journals" class="transition hover:text-text">Public journals</NuxtLink></li>
            <li><NuxtLink to="/prices" class="transition hover:text-text">Live prices</NuxtLink></li>
            <li><NuxtLink to="/about" class="transition hover:text-text">About</NuxtLink></li>
          </ul>
        </nav>
        <nav aria-label="Account">
          <p class="text-sm font-semibold">Account</p>
          <ul class="mt-2 space-y-1.5 text-sm text-muted">
            <li><NuxtLink to="/auth/login" class="transition hover:text-text">Log in</NuxtLink></li>
            <li><NuxtLink to="/auth/signup" class="transition hover:text-text">Sign up</NuxtLink></li>
            <li><NuxtLink to="/dashboard" class="transition hover:text-text">Dashboard</NuxtLink></li>
          </ul>
        </nav>
      </div>
      <p v-if="footerText" class="mt-8 text-xs text-muted">
        {{ footerText }}
      </p>
      <p v-if="riskDisclaimer" class="mt-2 border-t border-border pt-4 text-xs leading-relaxed text-muted">
        Risk disclaimer: {{ riskDisclaimer }}
      </p>
      <p class="mt-4 text-xs text-muted">
        © {{ new Date().getFullYear() }} {{ siteName }}. All rights reserved.
      </p>
    </div>
  </footer>
</template>
