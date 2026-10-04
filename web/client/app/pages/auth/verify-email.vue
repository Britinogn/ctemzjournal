<script setup lang="ts">
import {
  Mail01Icon,
  ArrowLeft01Icon,
} from '~/utils/icons'
</script>

<template>
  <div class="w-full max-w-md">
    <div class="rounded-2xl border border-border bg-surface p-7 text-center shadow-sm sm:p-8">

      <!-- Header -->
      <div class="mb-6">
        <div class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 text-primary">
          <UiAppIcon :icon="Mail01Icon" :size="22" />
        </div>
        <h1 class="text-2xl font-bold tracking-tight">
          Check your inbox
        </h1>
        <p class="mt-2 text-sm text-muted">
          We sent a verification link to your email.
          Open it to activate your journal, then log in.
        </p>
      </div>

      <!-- Action -->
      <NuxtLink
        to="/auth/login"
        class="inline-flex w-full items-center justify-center rounded-xl bg-primary py-3 text-sm font-semibold text-on-primary transition hover:opacity-90"
      >
        Go to login
      </NuxtLink>

      <!-- Footer -->
      <div class="mt-6">
        <NuxtLink
          to="/"
          class="inline-flex items-center gap-1.5 text-sm font-medium text-muted transition hover:text-primary"
        >
          <UiAppIcon :icon="ArrowLeft01Icon" :size="16" />
          Back to home
        </NuxtLink>
      </div>
    </div>
  </div>
</template>