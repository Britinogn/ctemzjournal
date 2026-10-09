// composables/usePageSeo.ts
export function usePageSeo(title: string, description: string) {
  const full = `${title} · Ctemz Journal`
  useHead({ title })
  useSeoMeta({
    description,
    ogTitle: full,
    ogDescription: description,
    twitterTitle: full,
    twitterDescription: description,
  })
}