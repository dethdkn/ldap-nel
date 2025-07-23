import process from 'node:process'

export default defineNuxtConfig({
  modules: ['@nuxt/image', '@nuxt/ui', '@nuxtjs/seo', '@vueuse/nuxt', 'nuxt-security'],
  $development: {
    security: { headers: { crossOriginEmbedderPolicy: 'unsafe-none' } },
  },
  devtools: { enabled: true },
  app: { head: { templateParams: { separator: '•' } } },
  css: ['~/assets/main.css'],
  site: {
    url: process.env.SITE_URL || 'http://localhost:3000',
    name: 'Ldap Nel',
    description: '🔐 Modern web interface for LDAP administration',
    indexable: false,
  },
  routeRules: {
    '/api/**': { proxy: { to: 'http://localhost:8080/**' } },
  },
  compatibilityDate: '2025-07-15',
  linkChecker: { enabled: false },
})
