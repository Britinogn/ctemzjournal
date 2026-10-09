<script setup lang="ts">
const siteUrl = 'https://ctemzjournal.onrender.com'
const siteName = 'Ctemz Journal'
const defaultTitle = 'Ctemz Journal — Know why you win, and why you lose.'
const defaultDescription =
  'Journal every forex trade, track your performance, review your setups, and build a record of what you learn from the market.'

const route = useRoute()

const canonical = computed(
  () => `${siteUrl}${route.path === '/' ? '' : route.path}`,
)

useHead({
  htmlAttrs: { lang: 'en' },

  titleTemplate: (title) => (title ? `${title} · ${siteName}` : defaultTitle),

  link: [
    { rel: 'icon', type: 'image/png', sizes: '96x96', href: '/img/favicon-96x96.png' },
    { rel: 'icon', type: 'image/svg+xml', href: '/img/favicon.svg' },
    { rel: 'shortcut icon', href: '/img/favicon.ico' },
    { rel: 'apple-touch-icon', sizes: '180x180', href: '/img/apple-touch-icon.png' },
    { rel: 'manifest', href: '/img/site.webmanifest' },
    { rel: 'canonical', href: canonical },
    { rel: 'mask-icon', href: '/img/safari-pinned-tab.svg', color: '#0E7C86' },
    { rel: 'sitemap', type: 'application/xml', href: '/sitemap.xml' },
    { rel: 'preconnect', href: 'https://www.sabilytics.com', crossorigin: '' },
    { rel: 'dns-prefetch', href: 'https://www.sabilytics.com' },
  ],

  meta: [
    { name: 'color-scheme', content: 'light' },
    { name: 'theme-color', content: '#0E7C86' },
    { name: 'apple-mobile-web-app-title', content: 'Ctemz Journal' },
    { name: 'apple-mobile-web-app-capable', content: 'yes' },
    { name: 'mobile-web-app-capable', content: 'yes' },
    { name: 'apple-mobile-web-app-status-bar-style', content: 'default' },
    { name: 'referrer', content: 'strict-origin-when-cross-origin' },
    { name: 'format-detection', content: 'telephone=no' },
    {
      name: 'keywords',
      content:
        'forex journal, trading journal, forex trading journal, trade journal, trading log, forex trades, trading performance',
    },
  ],

  script: [
    {
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
     // ── Analytics ───────────────────────────────────────────────────
    {
      src: 'https://www.sabilytics.com/script.js',
      async: true,
      'data-site': '0cffnjqpm646',
      'data-domain': 'ctemzjournal.onrender.com',
      // 'strategy'="afterInteractive",
    },
  ],
})

useSeoMeta({
  title: defaultTitle,
  description: defaultDescription,
  robots: 'index, follow, max-snippet:-1, max-image-preview:large, max-video-preview:-1',
  author: 'Ctemz Journal',

  ogType: 'website',
  ogSiteName: siteName,
  ogTitle: defaultTitle,
  ogDescription: defaultDescription,
  ogUrl: canonical,
  ogImage: `${siteUrl}/img/og-image.png`,
  ogImageType: 'image/png',
  ogImageWidth: 1200,
  ogImageHeight: 630,
  ogImageAlt: 'Ctemz Journal — Know why you win, and why you lose.',
  ogLocale: 'en_US',

  twitterCard: 'summary_large_image',
  twitterTitle: defaultTitle,
  twitterDescription: defaultDescription,
  twitterImage: `${siteUrl}/img/og-image.png`,
  twitterImageAlt: 'Ctemz Journal — Know why you win, and why you lose.',
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