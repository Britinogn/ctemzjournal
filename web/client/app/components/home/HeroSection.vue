<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    siteName?: string

    /** From /public/site-settings. Used as the headline. */
    tagline?: string | null

    dashboardTo?: string | null
    allowSignups?: boolean
  }>(),
  {
    siteName: 'Ctemz Journal',
    tagline: null,
    dashboardTo: null,
    allowSignups: true,
  },
)

const DEFAULT_HEADLINE = 'Know why you win, and why you lose.'

const headline = computed(() => props.tagline?.trim() || DEFAULT_HEADLINE)

// The primary button always says where it goes.
const primary = computed(() => {
  if (props.dashboardTo) {
    return {
      to: props.dashboardTo,
      label: 'Open your journal',
    }
  }

  if (props.allowSignups) {
    return {
      to: '/auth/signup',
      label: 'Start your journal',
    }
  }

  return {
    to: '/auth/login',
    label: 'Log in to your journal',
  }
})
</script>

<template>
  <section
    class="mx-auto w-full max-w-6xl px-4 pt-12 text-center md:px-6 md:pt-24"
    :aria-label="`${siteName} hero`"
  >
    <!-- Trust / positioning -->
    <p
      class="inline-flex items-center gap-1.5 rounded-full border border-border bg-surface px-3 py-1 text-xs font-semibold text-primary"
    >
      <UiAppIcon :icon="LockKeyIcon" :size="14" />
      Private by default
    </p>

    <!-- Headline -->
    <h1
      class="mx-auto mt-5 max-w-3xl text-balance text-4xl font-extrabold leading-[1.08] tracking-tight sm:text-5xl md:text-6xl"
    >
      {{ headline }}
    </h1>

    <!-- Main value proposition -->
    <p
      class="mx-auto mt-4 max-w-xl text-pretty text-base text-muted md:text-lg"
    >
      Log every forex trade with its setup, risk and screenshots. See your real
      numbers, and whether you followed your own rules.
    </p>

    <!-- Product positioning -->
    <p class="mx-auto mt-3 max-w-xl text-pretty text-sm font-medium text-primary">
      Your pipbook. Every trade, written down.
    </p>

    <!-- Actions -->
    <div
      class="mx-auto mt-8 flex max-w-sm flex-col items-stretch gap-2.5 sm:max-w-none sm:flex-row sm:items-center sm:justify-center"
    >
      <NuxtLink
        :to="primary.to"
        class="inline-flex h-12 items-center justify-center rounded-xl bg-primary px-7 text-sm font-semibold text-on-primary transition-colors hover:opacity-90 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
      >
        {{ primary.label }}
      </NuxtLink>

      <NuxtLink
        to="/journals"
        class="inline-flex h-12 items-center justify-center rounded-xl border border-border bg-surface px-7 text-sm font-semibold transition-colors hover:border-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
      >
        Browse public journals
      </NuxtLink>
    </div>

    <p class="mt-4 text-sm text-muted">
      Journals stay private until you choose to share one.
    </p>

    <!-- Product preview -->
    <div
      class="mx-auto mt-12 w-full max-w-3xl text-left mask-[linear-gradient(to_bottom,black_72%,transparent)]"
      role="img"
      :aria-label="`Preview of the ${siteName} dashboard: an equity curve rising to +$7,928, a 60% win rate, +0.83R average and 3 open trades.`"
    >
      <div
        class="overflow-hidden rounded-t-2xl border border-b-0 border-border bg-surface shadow-xl"
      >
        <!-- Browser header -->
        <div
          class="flex items-center gap-1.5 border-b border-border px-4 py-2.5"
          aria-hidden="true"
        >
          <span class="h-2.5 w-2.5 rounded-full bg-border" />
          <span class="h-2.5 w-2.5 rounded-full bg-border" />
          <span class="h-2.5 w-2.5 rounded-full bg-border" />

          <span
            class="ml-2 rounded-md bg-bg px-2 py-0.5 text-[11px] text-muted"
          >
            Dashboard
          </span>
        </div>

        <!-- Dashboard preview -->
        <div
          class="grid gap-3 p-4 pb-10 sm:grid-cols-5"
          aria-hidden="true"
        >
          <!-- Equity -->
          <div class="sm:col-span-3">
            <p class="text-xs font-medium text-muted">
              Equity curve
            </p>

            <p class="tnum text-2xl font-bold">
              +$7,928
            </p>

            <div class="relative mt-2">
              <!-- The curve is revealed left to right once. -->
              <svg
                viewBox="0 0 200 60"
                class="hero-reveal h-28 w-full"
                preserveAspectRatio="none"
              >
                <defs>
                  <linearGradient
                    id="hero-fill"
                    x1="0"
                    y1="0"
                    x2="0"
                    y2="1"
                  >
                    <stop
                      offset="0"
                      stop-color="var(--primary)"
                      stop-opacity="0.22"
                    />
                    <stop
                      offset="1"
                      stop-color="var(--primary)"
                      stop-opacity="0"
                    />
                  </linearGradient>
                </defs>

                <path
                  d="M0,48 L25,44 L50,46 L75,34 L100,38 L125,26 L150,30 L175,16 L200,8 L200,60 L0,60 Z"
                  fill="url(#hero-fill)"
                />

                <path
                  d="M0,48 L25,44 L50,46 L75,34 L100,38 L125,26 L150,30 L175,16 L200,8"
                  fill="none"
                  stroke="var(--primary)"
                  stroke-width="2.5"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  vector-effect="non-scaling-stroke"
                />
              </svg>

              <!-- Drawn in HTML so the stretched SVG does not squash it into an oval -->
              <span
                class="hero-dot absolute right-0 top-[13%] h-3 w-3 -translate-y-1/2 translate-x-1/2 rounded-full border-2 border-primary bg-surface"
              />
            </div>
          </div>

          <!-- Stats -->
          <div
            class="grid grid-cols-3 gap-2 sm:col-span-2 sm:grid-cols-1"
          >
            <div class="rounded-xl bg-bg p-3">
              <p class="text-xs text-muted">
                Win rate
              </p>

              <p class="tnum text-base font-bold">
                60%
              </p>
            </div>

            <div class="rounded-xl bg-bg p-3">
              <p class="text-xs text-muted">
                Average R
              </p>

              <p class="tnum text-base font-bold text-profit-text">
                +0.83R
              </p>
            </div>

            <div class="rounded-xl bg-bg p-3">
              <p class="text-xs text-muted">
                Open trades
              </p>

              <p class="tnum text-base font-bold">
                3
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
    
  </section>
</template>

<style scoped>
.hero-reveal {
  clip-path: inset(0 100% 0 0);
  animation: hero-reveal 1.2s 0.25s cubic-bezier(0.22, 0.8, 0.3, 1) forwards;
}

.hero-dot {
  opacity: 0;
  animation: hero-dot 0.3s 1.35s ease-out forwards;
}

@keyframes hero-reveal {
  to {
    clip-path: inset(0 0 0 0);
  }
}

@keyframes hero-dot {
  to {
    opacity: 1;
  }
}

@media (prefers-reduced-motion: reduce) {
  .hero-reveal {
    animation: none;
    clip-path: none;
  }

  .hero-dot {
    animation: none;
    opacity: 1;
  }
}
</style>