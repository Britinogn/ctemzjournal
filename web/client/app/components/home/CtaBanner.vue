<script setup lang="ts">
const props = withDefaults(defineProps<{
  dashboardTo?: string | null
  allowSignups?: boolean
}>(), {
  dashboardTo: null,
  allowSignups: true,
})

// Same rule as the hero: the button always says where it goes.
const primary = computed(() => {
  if (props.dashboardTo) return { to: props.dashboardTo, label: 'Open your journal' }
  if (props.allowSignups) return { to: '/auth/signup', label: 'Start your journal' }
  return { to: '/auth/login', label: 'Log in to your journal' }
})
</script>

<template>
  <section class="mx-auto w-full max-w-6xl px-4 py-16 md:px-6 md:py-20" aria-labelledby="cta-title">
    <div class="rounded-3xl border border-primary/20 bg-primary/5 px-6 py-12 text-center md:py-14">
      <h2 id="cta-title" class="mx-auto max-w-lg text-balance text-2xl font-extrabold tracking-tight md:text-3xl">
        Your edge is in your history. Start reading it.
      </h2>
      <p class="mx-auto mt-3 max-w-md text-pretty text-base text-muted">
        Your trades stay private until you choose to share them.
      </p>
      <div class="mx-auto mt-7 flex max-w-sm flex-col gap-2.5 sm:max-w-none sm:flex-row sm:justify-center">
        <NuxtLink
          :to="primary.to"
          class="inline-flex h-12 items-center justify-center rounded-xl bg-primary px-7 text-sm font-semibold text-on-primary transition-colors hover:opacity-90 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
        >
          {{ primary.label }}
        </NuxtLink>
        <NuxtLink
          to="/journals"
          class="inline-flex h-12 items-center justify-center rounded-xl border border-border bg-surface px-7 text-sm font-semibold transition-colors hover:border-primary focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
        >
          Browse public journals
        </NuxtLink>
      </div>
    </div>
  </section>
</template>