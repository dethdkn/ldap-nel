import { description, version } from '../package.json'

export default defineNuxtConfig({
  modules: ['@nuxt/image', '@nuxt/ui', '@nuxtjs/seo', '@vueuse/nuxt', 'nuxt-security'],
  $development: {
    security: { headers: { crossOriginEmbedderPolicy: 'unsafe-none' } },
  },
  devtools: { enabled: true },
  app: { head: { templateParams: { separator: '•' } } },
  css: ['~/assets/main.css'],
  site: {
    name: 'Ldap Nel',
    description,
    indexable: false,
  },
  runtimeConfig: {
    public: { version },
  },
  routeRules: {
    '/server/**': { proxy: { to: 'http://localhost:8080/**' } },
  },
  compatibilityDate: '2025-07-15',
  linkChecker: { enabled: false },
  sitemap: { enabled: false },
})
