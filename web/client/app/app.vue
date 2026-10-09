<script setup lang="ts">
const siteUrl = 'https://ctemzjournal.onrender.com'
const siteName = 'Ctemz Journal'
const defaultTitle = 'Ctemz Journal — Forex Trading Journal & Analytics'
// Keep under ~160 characters so Google doesn't cut it off.
const defaultDescription
  = 'Log your forex trades, track risk, analyze win rates and cut emotional trading with Ctemz Journal, your free digital trading journal.'
const ogImage = `${siteUrl}/img/og-image.png`
const ogImageAlt = 'Ctemz Journal — Know why you win, and why you lose.'

const route = useRoute()

// Reactive so it updates on every navigation.
const canonical = computed(
  () => `${siteUrl}${route.path === '/' ? '' : route.path}`,
)

// Private areas should never be indexed.
const isPrivate = computed(() =>
  ['/admin', '/dashboard', '/auth'].some(p => route.path === p || route.path.startsWith(`${p}/`)),
)
const robots = computed(() =>
  isPrivate.value
    ? 'noindex, nofollow'
    : 'index, follow, max-snippet:-1, max-image-preview:large, max-video-preview:-1',
)

useHead({
  htmlAttrs: { lang: 'en' },

  // Pages set a plain title (e.g. 'Privacy Policy'); the site name is added here.
  titleTemplate: title =>
    !title || title === defaultTitle || title === siteName
      ? defaultTitle
      : `${title} · ${siteName}`,

  link: [
    { rel: 'icon', type: 'image/png', sizes: '96x96', href: '/img/favicon-96x96.png' },
    { rel: 'icon', type: 'image/svg+xml', href: '/img/favicon.svg' },
    { rel: 'shortcut icon', href: '/img/favicon.ico' },
    { rel: 'apple-touch-icon', sizes: '180x180', href: '/img/apple-touch-icon.png' },
    { rel: 'manifest', href: '/img/site.webmanifest' },
    { rel: 'canonical', href: canonical },
    { rel: 'mask-icon', href: '/img/safari-pinned-tab.svg', color: '#0E7C86' },
    { rel: 'preconnect', href: 'https://www.sabilytics.com', crossorigin: '' },
    { rel: 'dns-prefetch', href: 'https://www.sabilytics.com' },
  ],

  meta: [
    // 'color-scheme: light' removed so browsers can follow the user's theme.
    // If your app is light-only, add it back: { name: 'color-scheme', content: 'light' }
    { name: 'theme-color', content: '#0E7C86' },
    { name: 'apple-mobile-web-app-title', content: siteName },
    { name: 'apple-mobile-web-app-capable', content: 'yes' },
    { name: 'mobile-web-app-capable', content: 'yes' },
    { name: 'apple-mobile-web-app-status-bar-style', content: 'default' },
    { name: 'referrer', content: 'strict-origin-when-cross-origin' },
    { name: 'format-detection', content: 'telephone=no' },
  ],

  script: [
    {
      key: 'ld-website',
      type: 'application/ld+json',
      innerHTML: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'WebSite',
        name: siteName,
        url: siteUrl,
        description: defaultDescription,
        publisher: {
          '@type': 'Organization',
          name: siteName,
          logo: { '@type': 'ImageObject', url: `${siteUrl}/img/logo.png` },
        },
      }),
    },
    {
      key: 'ld-app',
      type: 'application/ld+json',
      innerHTML: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'SoftwareApplication',
        name: siteName,
        url: siteUrl,
        applicationCategory: 'FinanceApplication',
        operatingSystem: 'Web',
        description: defaultDescription,
        offers: { '@type': 'Offer', price: '0', priceCurrency: 'USD' },
      }),
    },
    // Analytics
    {
      key: 'analytics',
      src: 'https://www.sabilytics.com/script.js',
      async: true,
      'data-site': '0cffnjqpm646',
      'data-domain': 'ctemzjournal.onrender.com',
    },
  ],
})

// Defaults. Pages can override these (see composables/usePageSeo.ts).
useSeoMeta({
  title: defaultTitle,
  description: defaultDescription,
  robots,
  author: siteName,

  ogType: 'website',
  ogSiteName: siteName,
  ogTitle: defaultTitle,
  ogDescription: defaultDescription,
  ogUrl: canonical,
  ogImage,
  ogImageType: 'image/png',
  ogImageWidth: 1200,
  ogImageHeight: 630,
  ogImageAlt,
  ogLocale: 'en_US',

  twitterCard: 'summary_large_image',
  twitterTitle: defaultTitle,
  twitterDescription: defaultDescription,
  twitterImage: ogImage,
  twitterImageAlt: ogImageAlt,
  // twitterSite: '@yourhandle',
})
</script>

<template>
  <div>
    <NuxtLoadingIndicator />
    <NuxtRouteAnnouncer />
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <Toaster position="top-right" :gap="8" />
  </div>
</template>