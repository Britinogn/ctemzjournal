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
    const updated = await api.patch<SiteSettingRow>('/admin/site-settings', body)
    queryClient.setQueryData(['admin', 'site-settings'], updated)
    queryClient.invalidateQueries({ queryKey: siteSettingsKey() })
    toast.success('Site settings saved — home page and header update')
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
    <h1 class="mb-4 text-xl font-bold tracking-tight md:text-2xl">
      Site settings
    </h1>

    <div v-if="isPending" class="h-96 animate-pulse rounded-2xl bg-surface" />

    <form v-else class="grid items-start gap-4 xl:grid-cols-2" @submit.prevent="onSave">
      <!-- Identity -->
      <section class="rounded-2xl border border-border bg-surface p-4 md:p-6" aria-label="Identity">
        <h2 class="text-sm font-semibold">Identity</h2>
        <p class="mt-0.5 text-xs text-muted">Shown on the home page and header everywhere.</p>
        <div class="mt-4 space-y-4">
          <div>
            <label for="ss-name" :class="labelCls">Site name</label>
            <input id="ss-name" v-model="form.site_name" type="text" required :class="inputCls">
          </div>
          <div>
            <label for="ss-tag" :class="labelCls">Tagline</label>
            <input id="ss-tag" v-model="form.tagline" type="text" :class="inputCls">
          </div>
          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <span :class="labelCls">Logo</span>
              <div class="flex items-center gap-3">
                <NuxtImg
                  v-if="form.logo_path"
                  :src="assetUrl(form.logo_path)"
                  alt="Logo preview"
                  width="48" height="48"
                  class="h-12 w-12 rounded-xl border border-border object-contain"
                />
                <label class="cursor-pointer rounded-xl border border-border px-4 py-2 text-sm font-medium transition hover:border-primary">
                  <input type="file" accept="image/png,image/jpeg,image/webp,image/svg+xml" class="sr-only" @change="uploadAsset('logo', $event)">
                  {{ uploading === 'logo' ? 'Uploading…' : 'Upload logo' }}
                </label>
              </div>
            </div>
            <div>
              <span :class="labelCls">Favicon</span>
              <div class="flex items-center gap-3">
                <NuxtImg
                  v-if="form.favicon_path"
                  :src="assetUrl(form.favicon_path)"
                  alt="Favicon preview"
                  width="32" height="32"
                  class="h-8 w-8 rounded-lg border border-border object-contain"
                />
                <label class="cursor-pointer rounded-xl border border-border px-4 py-2 text-sm font-medium transition hover:border-primary">
                  <input type="file" accept="image/png,image/jpeg,image/webp,image/svg+xml,image/x-icon" class="sr-only" @change="uploadAsset('favicon', $event)">
                  {{ uploading === 'favicon' ? 'Uploading…' : 'Upload' }}
                </label>
              </div>
            </div>
          </div>
          <p class="text-[11px] text-muted">Uploads straight to the site-assets bucket, then save to apply.</p>
        </div>
      </section>

      <!-- Content + access -->
      <div class="space-y-4">
        <section class="rounded-2xl border border-border bg-surface p-4 md:p-6" aria-label="Content">
          <h2 class="text-sm font-semibold">Content</h2>
          <div class="mt-4 space-y-4">
            <div>
              <label for="ss-email" :class="labelCls">Contact email</label>
              <input id="ss-email" v-model="form.contact_email" type="email" :class="inputCls">
            </div>
            <div>
              <label for="ss-footer" :class="labelCls">Footer text</label>
              <textarea id="ss-footer" v-model="form.footer_text" rows="2" :class="inputCls" />
            </div>
            <div>
              <label for="ss-risk" :class="labelCls">Risk disclaimer</label>
              <textarea id="ss-risk" v-model="form.risk_disclaimer" rows="3" :class="inputCls" />
            </div>
            <div>
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
              <p v-if="socials.length === 0" class="text-xs text-muted">No social links yet.</p>
            </div>
          </div>
        </section>

        <section class="rounded-2xl border border-border bg-surface p-4 md:p-6" aria-label="Access">
          <h2 class="text-sm font-semibold">Access</h2>
          <div class="mt-3 space-y-3">
            <label class="flex cursor-pointer items-center justify-between gap-3">
              <span class="text-sm font-medium">Allow signups</span>
              <input v-model="form.allow_signups" type="checkbox" class="peer sr-only">
              <span class="relative h-6 w-11 shrink-0 rounded-full bg-border transition peer-checked:bg-primary peer-focus-visible:ring-2 peer-focus-visible:ring-primary after:absolute after:left-0.5 after:top-0.5 after:h-5 after:w-5 after:rounded-full after:bg-white after:transition peer-checked:after:translate-x-5" />
            </label>
            <label class="flex cursor-pointer items-center justify-between gap-3">
              <span>
                <span class="block text-sm font-medium">Maintenance mode</span>
                <span class="block text-xs text-muted">Visitors see the maintenance page.</span>
              </span>
              <input v-model="form.maintenance_mode" type="checkbox" class="peer sr-only">
              <span class="relative h-6 w-11 shrink-0 rounded-full bg-border transition peer-checked:bg-warning peer-focus-visible:ring-2 peer-focus-visible:ring-warning after:absolute after:left-0.5 after:top-0.5 after:h-5 after:w-5 after:rounded-full after:bg-white after:transition peer-checked:after:translate-x-5" />
            </label>
          </div>
          <button
            type="submit"
            :disabled="saving"
            class="mt-4 w-full rounded-xl bg-primary py-2.5 text-sm font-semibold text-on-primary transition hover:opacity-90 disabled:opacity-50"
          >
            {{ saving ? 'Saving…' : 'Save settings' }}
          </button>
        </section>
      </div>
    </form>
  </div>
</template>
