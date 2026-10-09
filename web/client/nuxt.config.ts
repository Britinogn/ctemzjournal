// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  css: ['~/assets/css/main.css'],

  imports: {
    dirs: ['composables/**'],
  },

  vite: {
    plugins: [
      tailwindcss(),
    ],
  },

  modules: [
    '@nuxtjs/color-mode',
    '@nuxt/fonts',
    '@nuxt/eslint',
    'vue-sonner/nuxt',
    '@nuxt/image',
    '@nuxtjs/sitemap',
  ],

  site: {
    url: 'https://ctemzjournal.onrender.com',
    name: 'Ctemz Journal',
  },

  sitemap: {
    exclude: ['/admin/**', '/dashboard/**', '/auth/**'],
  },

  // Remote images (Supabase site-assets logos) must be allowlisted or
  // NuxtImg's optimizer rejects them. Derived from env so local and
  // production each permit their own Supabase host.
  image: {
    domains: [
      (process.env.NUXT_PUBLIC_SUPABASE_URL || '')
        .replace(/^https?:\/\//, '')
        .split('/')[0],
    ].filter((d): d is string => d !== '' && d !== undefined),

    // screens: {
    //   sm: 640,
    //   md: 768,
    //   lg: 1024,
    //   xl: 1280,
    //   '2xl': 1536,
    //   '3xl': 1920,
    //   '4xl': 2560,
    // },
  },

  // Plain `.dark` class (no suffix) — the design tokens switch on it.
  // Follows the OS, remembers the choice, no flash on load.
  colorMode: {
    classSuffix: '',
    preference: 'system',
    fallback: 'light',
    storageKey: 'ctemzjournal-theme',
  },

  fonts: {
    defaults: {
      weights: [400, 500, 600, 700, 800],
    },
    families: [
      { name: 'Plus Jakarta Sans', provider: 'google' },
      { name: 'JetBrains Mono', provider: 'google' },
    ],
  },

  runtimeConfig: {
    public: {
      // Go API base URL (no trailing slash).
      apiUrl: process.env.NUXT_PUBLIC_API_URL || 'http://localhost:8080',
      supabaseUrl: process.env.NUXT_PUBLIC_SUPABASE_URL || '',
      supabaseAnonKey: process.env.NUXT_PUBLIC_SUPABASE_ANON_KEY || '',
      cloudinaryCloudName: process.env.NUXT_PUBLIC_CLOUDINARY_CLOUD_NAME || '',
    },
  },
})