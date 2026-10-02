<template>
  <div class="flex min-h-screen flex-col bg-bg text-text">
    <header class="border-b border-border bg-surface">
      <div class="mx-auto flex h-16 w-full max-w-md items-center justify-between px-4 sm:max-w-lg">
        <BrandLogo />
        <ThemeToggle />
      </div>
    </header>
    <main class="flex flex-1 items-start justify-center px-4 py-10">
      <div class="w-full max-w-md sm:max-w-lg">
        <slot />
      </div>
    </main>
    <footer class="border-t border-border py-4 text-center text-xs text-muted">
      Trading involves risk. Past performance is not indicative of future results.
    </footer>
  </div>
</template>
