<script setup lang="ts">
import { CheckmarkCircle01Icon, EyeIcon, ZapIcon } from '~/utils/icons'

// These really are a sequence (log, then learn, then share), so the list is numbered.
const steps = [
  {
    icon: ZapIcon,
    title: 'Log a trade in under a minute',
    body: 'Enter the pair, entry, stop and size. Risk and R preview as you type, and the server checks every number.',
  },
  {
    icon: CheckmarkCircle01Icon,
    title: 'See if you kept your rules',
    body: 'Tag setups, mistakes and emotions. Your rule-following rate, expectancy and drawdown update with every trade.',
  },
  {
    icon: EyeIcon,
    title: 'Share without showing money',
    body: 'Switch on Make public and only the pair, setup, result in R and your first chart appear on the home page. Lot size and P&L stay private.',
  },
]
</script>

<template>
  <section class="mx-auto w-full max-w-6xl px-4 pt-16 md:px-6" aria-labelledby="how-title">
    <div class="mb-5">
      <p class="text-xs font-semibold uppercase tracking-wide text-primary">
        How it works
      </p>
      <h2 id="how-title" class="text-2xl font-bold tracking-tight">
        Log it, learn from it, share it.
      </h2>
    </div>

    <ol class="grid gap-3 md:grid-cols-3">
      <li
        v-for="(s, i) in steps"
        :key="s.title"
        class="rounded-2xl border border-border bg-surface p-5"
      >
        <div class="flex items-start justify-between">
          <span class="inline-flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 text-primary">
            <UiAppIcon :icon="s.icon" :size="20" />
          </span>
          <!-- Hidden from screen readers: the list already announces "1 of 3" -->
          <span class="tnum text-sm font-semibold text-muted" aria-hidden="true">{{ i + 1 }}</span>
        </div>
        <h3 class="mt-4 text-base font-bold tracking-tight">
          {{ s.title }}
        </h3>
        <p class="mt-1.5 text-[15px] leading-relaxed text-muted">
          {{ s.body }}
        </p>
      </li>
    </ol>
  </section>
</template>