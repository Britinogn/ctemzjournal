<script setup lang="ts">
type SizeMap = number | {
  base?: number
  sm?: number
  md?: number
  lg?: number
  xl?: number
  '2xl'?: number
}

const props = withDefaults(defineProps<{
  logoUrl?: string | null
  size?: SizeMap
  showName?: boolean
  siteName?: string
}>(), {
  logoUrl: null,
  size: 32,
  showName: true,
  siteName: 'Ctemz Journal',
})

const src = computed(() => props.logoUrl?.trim() || '/logo.svg')

// Narrow the union once. `null` means "size is a plain number".
const sizeMap = computed<Record<string, number> | null>(() => {
  if (typeof props.size === 'number') return null
  const out: Record<string, number> = {}
  for (const [bp, px] of Object.entries(props.size)) {
    if (typeof px === 'number') out[bp] = px
  }
  return out
})

const maxSize = computed(() => {
  const map = sizeMap.value
  if (!map) return typeof props.size === 'number' ? props.size : 32
  const values = Object.values(map)
  return values.length ? Math.max(...values) : 32
})

const bpPrefix: Record<string, string> = {
  base: '', sm: 'sm:', md: 'md:', lg: 'lg:', xl: 'xl:', '2xl': '2xl:',
}

const iconClasses = computed(() => {
  const map = sizeMap.value
  if (!map) return ''
  return Object.entries(map)
    .map(([bp, px]) => `${bpPrefix[bp] ?? ''}size-[${px}px]`)
    .join(' ')
})

// Text scales with the icon so the lockup stays visually balanced.
const nameClasses = computed(() => {
  const map = sizeMap.value
  if (!map) return 'text-lg'
  const textMap: Record<string, string> = {
    sm: 'sm:text-xl',
    md: 'md:text-2xl',
    lg: 'lg:text-3xl',
    xl: 'xl:text-4xl',
    '2xl': '2xl:text-5xl',
  }
  return ['text-lg', ...Object.keys(map).map((k) => textMap[k] ?? '')]
    .filter(Boolean)
    .join(' ')
})

// Gap scales too, otherwise the icon and text visually touch at large sizes.
const gapClasses = computed(() => {
  const map = sizeMap.value
  if (!map) return 'gap-2'
  const gapMap: Record<string, string> = {
    sm: 'sm:gap-3',
    md: 'md:gap-4',
    lg: 'lg:gap-5',
    xl: 'xl:gap-6',
    '2xl': '2xl:gap-7',
  }
  return ['gap-2', ...Object.keys(map).map((k) => gapMap[k] ?? '')]
    .filter(Boolean)
    .join(' ')
})
</script>

<template>
  <span
    class="inline-flex items-center whitespace-nowrap"
    :class="gapClasses"
  >
    <NuxtImg
      :src="src"
      :alt="siteName"
      :width="maxSize"
      :height="maxSize"
      class="shrink-0 rounded-lg"
      :class="iconClasses"
    />
    <span
      v-if="showName"
      class="font-bold tracking-tight"
      :class="nameClasses"
    >
      <!-- {{ siteName }} -->
    </span>
  </span>
</template>