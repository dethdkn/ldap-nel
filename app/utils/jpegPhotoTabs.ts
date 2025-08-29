import type { TabsItem } from '@nuxt/ui'

export default <TabsItem[]> [
  { label: 'Plain Text', icon: 'i-lucide-case-lower', slot: 'text' },
  { label: 'File Upload', icon: 'i-lucide-image', slot: 'file' },
]
