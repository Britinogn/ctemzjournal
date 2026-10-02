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
  ],

  // No public home build yet: / lands on the (guarded) dashboard.
  routeRules: {
    '/': { redirect: '/dashboard' },
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