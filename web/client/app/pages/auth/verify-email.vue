<template>
  <div class="rounded-2xl border border-border bg-surface p-6 text-center shadow-sm sm:p-8">
    <h1 class="text-2xl font-bold tracking-tight">
      Check your inbox
    </h1>
    <p class="mt-2 text-sm text-muted">
      We sent a verification link to your email. Open it to activate your journal,
      then log in.
    </p>
    <NuxtLink
      to="/auth/login"
      class="mt-6 inline-block w-full rounded-xl bg-primary py-2.5 text-sm font-semibold text-on-primary transition hover:opacity-90"
    >
      Go to login
    </NuxtLink>
  </div>
</template>
