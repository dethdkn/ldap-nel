import { defineConfig } from 'vitepress'
import { description, version } from '../../package.json'

export default defineConfig({
  title: 'Ldap Nel',
  description,
  lang: 'en-US',
  lastUpdated: true,
  cleanUrls: true,
  sitemap: {
    hostname: 'https://ldap-nel.vercel.app',
  },
  head: [
    ['meta', { name: 'theme-color', content: '#00c850' }],
    ['meta', { name: 'og:type', content: 'website' }],
    ['meta', { name: 'og:locale', content: 'en-US' }],
    ['meta', { name: 'og:site_name', content: 'Ldap Nel' }],
  ],
  themeConfig: {
    siteTitle: 'Ldap Nel',
    logo: '/nel.png',
    search: {
      provider: 'local',
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/dethdkn/ldap-nel' },
      { icon: { svg: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M16.36 14c.08-.66.14-1.32.14-2c0-.68-.06-1.34-.14-2h3.38c.16.64.26 1.31.26 2s-.1 1.36-.26 2m-5.15 5.56c.6-1.11 1.06-2.31 1.38-3.56h2.95a8.03 8.03 0 0 1-4.33 3.56M14.34 14H9.66c-.1-.66-.16-1.32-.16-2c0-.68.06-1.35.16-2h4.68c.09.65.16 1.32.16 2c0 .68-.07 1.34-.16 2M12 19.96c-.83-1.2-1.5-2.53-1.91-3.96h3.82c-.41 1.43-1.08 2.76-1.91 3.96M8 8H5.08A7.923 7.923 0 0 1 9.4 4.44C8.8 5.55 8.35 6.75 8 8m-2.92 8H8c.35 1.25.8 2.45 1.4 3.56A8.008 8.008 0 0 1 5.08 16m-.82-2C4.1 13.36 4 12.69 4 12s.1-1.36.26-2h3.38c-.08.66-.14 1.32-.14 2c0 .68.06 1.34.14 2M12 4.03c.83 1.2 1.5 2.54 1.91 3.97h-3.82c.41-1.43 1.08-2.77 1.91-3.97M18.92 8h-2.95a15.65 15.65 0 0 0-1.38-3.56c1.84.63 3.37 1.9 4.33 3.56M12 2C6.47 2 2 6.5 2 12a10 10 0 0 0 10 10a10 10 0 0 0 10-10A10 10 0 0 0 12 2"/></svg>' }, link: 'https://rosa.dev.br', ariaLabel: 'Me ;)' },
    ],
    nav: [
      { text: 'Home', link: '/' },
      {
        text: `V${version}`,
        items: [
          { text: 'Changelog', link: 'https://github.com/dethdkn/ldap-nel/releases', target: '_blank' },
          { text: 'Bugs', link: 'https://github.com/dethdkn/ldap-nel/issues', target: '_blank' },
          { text: 'Sponsor', link: 'https://github.com/sponsors/dethdkn', target: '_blank' },
        ],
      },
    ],
    sidebar: [
      {
        text: 'Get started',
        items: [
          { text: 'Reason', link: '/reason' },
          { text: 'Setup', link: '/setup' },
          { text: 'Upgrade', link: '/upgrade' },
          { text: 'Showcase', link: '/showcase' },
          { text: 'License', link: 'https://github.com/dethdkn/ldap-nel/blob/main/LICENSE', target: '_blank' },
        ],
      },
    ],
    editLink: {
      pattern: 'https://github.com/dethdkn/ldap-nel/edit/main/docs/:path',
      text: 'Edit this page on GitHub',
    },
    docFooter: {
      prev: 'Previous page',
      next: 'Next page',
    },
    footer: {
      message: 'Released under the MIT License.',
      copyright: '© 2025 Gabriel \'DethDKN\' Rosa',
    },
  },
})
