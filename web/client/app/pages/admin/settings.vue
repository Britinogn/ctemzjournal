<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { toast } from 'vue-sonner'
import type { SupabaseClient } from '@supabase/supabase-js'
import type { LogoTarget, SettingsUpdate } from '~/types'
import { siteSettingsKey } from '~/types'

definePageMeta({ middleware: 'admin', layout: 'admin' })

interface SiteSettingRow {
  ID: number;
  SiteName: string;
  Tagline: string | null;
  LogoPath: string | null;
  FaviconPath: string | null;
  ContactEmail: string | null;
  FooterText: string | null;
  RiskDisclaimer: string | null;
  SocialLinks: string | null;
  AllowSignups: boolean;
  MaintenanceMode: boolean;
}

const api = useApi()
const queryClient = useQueryClient()
const config = useRuntimeConfig()
const { $supabase } = useNuxtApp() as unknown as { $supabase: SupabaseClient }

const { data: settings, isPending } = useQuery({
  queryKey: ['admin', 'site-settings'],
  queryFn: () => api.get<SiteSettingRow>('/admin/site-settings'),
  staleTime: 60_000,
})

const form = reactive({
  site_name: '',
  tagline: '',
  logo_path: '',
  favicon_path: '',
  contact_email: '',
  footer_text: '',
  risk_disclaimer: '',
  allow_signups: true,
  maintenance_mode: false,
})
const socials = ref<Array<{ platform: string; url: string }>>([])
const saving = ref(false)
const uploading = ref<'logo' | 'favicon' | null>(null)

watchEffect(() => {
  if (!settings.value)
    return
  form.site_name = settings.value.SiteName ?? ''
  form.tagline = settings.value.Tagline ?? ''
  form.logo_path = settings.value.LogoPath ?? ''
  form.favicon_path = settings.value.FaviconPath ?? ''
  form.contact_email = settings.value.ContactEmail ?? ''
  form.footer_text = settings.value.FooterText ?? ''
  form.risk_disclaimer = settings.value.RiskDisclaimer ?? ''
  form.allow_signups = settings.value.AllowSignups
  form.maintenance_mode = settings.value.MaintenanceMode
  try {
    const parsed = settings.value.SocialLinks ? JSON.parse(atob(settings.value.SocialLinks)) : {}
    socials.value = Object.entries(parsed).map(([platform, url]) => ({ platform, url: String(url) }))
  }
  catch {
    socials.value = []
  }
})

function assetUrl(path: string): string {
  const base = String(config.public.supabaseUrl || '').replace(/\/$/, '')
  return `${base}/storage/v1/object/public/site-assets/${path.replace(/^\//, '')}`
}

function addSocial(): void {
  socials.value.push({ platform: '', url: '' })
}

function removeSocial(index: number): void {
  socials.value.splice(index, 1)
}

async function uploadAsset(kind: 'logo' | 'favicon', event: Event): Promise<void> {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file || uploading.value)
    return
  const ext = file.name.split('.').pop()?.toLowerCase() ?? ''
  uploading.value = kind
  try {
    const target = await api.post<LogoTarget>('/admin/site-settings/logo', { kind, extension: ext })
    const { error } = await $supabase.storage
      .from(target.bucket)
      .upload(target.path, file, { upsert: true, contentType: file.type || undefined })
    if (error)
      throw new Error(error.message)
    if (kind === 'logo')
      form.logo_path = target.path
    else
      form.favicon_path = target.path
    toast.success(`${kind === 'logo' ? 'Logo' : 'Favicon'} uploaded — save settings to apply`)
  }
  catch (err) {
    toast.error(err instanceof Error ? err.message : 'Upload failed — check the storage bucket policy')
  }
  finally {
    uploading.value = null
  }
}

async function onSave(): Promise<void> {
  if (saving.value)
    return
  saving.value = true
  try {
    const links: Record<string, string> = {}
    for (const { platform, url } of socials.value) {
      if (platform.trim() && url.trim())
        links[platform.trim().toLowerCase()] = url.trim()
    }
    const body: SettingsUpdate = {
      site_name: form.site_name.trim(),
      tagline: form.tagline.trim(),
      logo_path: form.logo_path.trim(),
      favicon_path: form.favicon_path.trim(),
      contact_email: form.contact_email.trim(),
      footer_text: form.footer_text.trim(),
      risk_disclaimer: form.risk_disclaimer.trim(),
      social_links: links,
      allow_signups: form.allow_signups,
      maintenance_mode: form.maintenance_mode,
    }
    const updated = await api.patch<SiteSettingRow>('/admin/site-settings', {...body})
    queryClient.setQueryData(['admin', 'site-settings'], updated)
    queryClient.invalidateQueries({ queryKey: siteSettingsKey() })
    toast.success('Settings saved — home page and header update')
  }
  catch {
    toast.error('Could not save settings')
  }
  finally {
    saving.value = false
  }
}

const inputCls = 'w-full rounded-xl border border-border bg-bg px-4 py-2.5 text-sm outline-none transition placeholder:text-muted focus:border-primary'
const labelCls = 'mb-1.5 block text-sm font-medium'
</script>

<template>
  <div>
    <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
      <h1 class="text-xl font-bold tracking-tight md:text-2xl">
        Settings
      </h1>
      <ThemeToggle />
    </div>

    <div v-if="isPending" class="h-96 animate-pulse rounded-2xl bg-surface" />

    <form v-else class="mx-auto w-full max-w-2xl space-y-4" @submit.prevent="onSave">
      <!-- Brand -->
      <section class="rounded-2xl border border-border bg-surface p-4 md:p-6" aria-label="Brand">
        <h2 class="text-sm font-semibold">Brand</h2>
        <div class="mt-3 flex items-center gap-4">
          <span class="inline-flex h-16 w-16 items-center justify-center overflow-hidden rounded-2xl border border-border">
            <BrandLogo
              :logo-url="form.logo_path ? assetUrl(form.logo_path) : null"
              :show-name="false"
              :size="56"
              site-name="Site logo"
            />
          </span>
          <div>
            <label class="inline-block cursor-pointer rounded-xl border border-border px-4 py-2 text-sm font-semibold transition hover:border-primary">
              <input type="file" accept="image/png,image/jpeg,image/webp,image/svg+xml" class="sr-only" @change="uploadAsset('logo', $event)">
              {{ uploading === 'logo' ? 'Uploading…' : 'Change logo' }}
            </label>
            <p class="mt-1.5 text-[11px] text-muted">Square, PNG or WebP. Saved with the settings.</p>
          </div>
        </div>
        <div class="mt-3 flex items-center gap-3">
          <span class="inline-flex h-10 w-10 items-center justify-center overflow-hidden rounded-xl border border-border">
            <BrandLogo
              :logo-url="form.favicon_path ? assetUrl(form.favicon_path) : null"
              :show-name="false"
              :size="32"
              site-name="Site favicon"
            />
          </span>
          <div>
            <label class="inline-block cursor-pointer rounded-xl border border-border px-4 py-2 text-sm font-semibold transition hover:border-primary">
              <input type="file" accept="image/png,image/jpeg,image/webp,image/svg+xml,image/x-icon" class="sr-only" @change="uploadAsset('favicon', $event)">
              {{ uploading === 'favicon' ? 'Uploading…' : 'Change favicon' }}
            </label>
          </div>
        </div>
        <div class="mt-4 grid gap-3 sm:grid-cols-2">
          <div>
            <label for="ss-name" :class="labelCls">Site name</label>
            <input id="ss-name" v-model="form.site_name" type="text" required :class="inputCls">
          </div>
          <div>
            <label for="ss-email" :class="labelCls">Contact email</label>
            <input id="ss-email" v-model="form.contact_email" type="email" :class="inputCls">
          </div>
        </div>
        <div class="mt-3">
          <label for="ss-tag" :class="labelCls">Tagline</label>
          <input id="ss-tag" v-model="form.tagline" type="text" :class="inputCls">
        </div>
        <div class="mt-3">
          <label for="ss-risk" :class="labelCls">Risk disclaimer</label>
          <textarea id="ss-risk" v-model="form.risk_disclaimer" rows="4" :class="inputCls" />
          <p class="mt-1 text-[11px] text-muted">Shown in the footer of every public page.</p>
        </div>
        <div class="mt-3">
          <label for="ss-footer" :class="labelCls">Footer text</label>
          <textarea id="ss-footer" v-model="form.footer_text" rows="2" :class="inputCls" />
        </div>
        <div class="mt-3">
          <div class="mb-1.5 flex items-center justify-between">
            <span class="text-sm font-medium">Social links</span>
            <button type="button" class="text-sm font-medium text-primary hover:underline" @click="addSocial">
              + Add
            </button>
          </div>
          <ul class="space-y-2">
            <li v-for="(s, i) in socials" :key="i" class="flex gap-2">
              <input v-model="s.platform" type="text" placeholder="x" aria-label="Platform" class="w-28 shrink-0 rounded-xl border border-border bg-bg px-3 py-2 text-sm outline-none transition placeholder:text-muted focus:border-primary">
              <input v-model="s.url" type="url" placeholder="https://…" aria-label="URL" class="min-w-0 flex-1 rounded-xl border border-border bg-bg px-3 py-2 text-sm outline-none transition placeholder:text-muted focus:border-primary">
              <button type="button" aria-label="Remove link" class="shrink-0 rounded-xl px-2 text-muted transition hover:text-loss" @click="removeSocial(i)">
                ✕
              </button>
            </li>
          </ul>
          <p v-if="socials.length === 0" class="text-xs text-muted">No social links yet — they appear as icons in the footer.</p>
        </div>
      </section>

      <!-- Access -->
      <section class="rounded-2xl border border-border bg-surface p-4 md:p-6" aria-label="Access">
        <h2 class="text-sm font-semibold">Access</h2>
        <div class="mt-2 divide-y divide-border">
          <label class="flex cursor-pointer items-center justify-between gap-3 py-3">
            <span>
              <span class="block text-sm font-medium">Allow signups</span>
              <span class="block text-xs text-muted">When off, the signup page and its links are hidden.</span>
            </span>
            <input v-model="form.allow_signups" type="checkbox" class="peer sr-only">
            <span class="relative h-6 w-11 shrink-0 rounded-full bg-border transition peer-checked:bg-primary peer-focus-visible:ring-2 peer-focus-visible:ring-primary after:absolute after:left-0.5 after:top-0.5 after:h-5 after:w-5 after:rounded-full after:bg-white after:transition peer-checked:after:translate-x-5" />
          </label>
          <label class="flex cursor-pointer items-center justify-between gap-3 py-3">
            <span>
              <span class="block text-sm font-medium">Maintenance mode</span>
              <span class="block text-xs text-muted">Traders see a maintenance page. Admins can still sign in.</span>
            </span>
            <input v-model="form.maintenance_mode" type="checkbox" class="peer sr-only">
            <span class="relative h-6 w-11 shrink-0 rounded-full bg-border transition peer-checked:bg-warning peer-focus-visible:ring-2 peer-focus-visible:ring-warning after:absolute after:left-0.5 after:top-0.5 after:h-5 after:w-5 after:rounded-full after:bg-white after:transition peer-checked:after:translate-x-5" />
          </label>
        </div>
      </section>

      <button
        type="submit"
        :disabled="saving"
        class="rounded-xl bg-primary px-6 py-2.5 text-sm font-semibold text-on-primary transition hover:opacity-90 disabled:opacity-50"
      >
        {{ saving ? 'Saving…' : 'Save settings' }}
      </button>
    </form>
  </div>
</template>
