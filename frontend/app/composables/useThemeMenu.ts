import type { DropdownMenuItem } from '@nuxt/ui'

export default () => {
  const colorMode = useColorMode()

  return computed<DropdownMenuItem[][]>(() => [[
    { label: 'System', icon: 'i-lucide-monitor', onSelect: () => colorMode.preference = 'system', disabled: colorMode.preference === 'system' },
    { label: 'Light', icon: 'i-lucide-sun', onSelect: () => colorMode.preference = 'light', disabled: colorMode.preference === 'light' },
    { label: 'Dark', icon: 'i-lucide-moon', onSelect: () => colorMode.preference = 'dark', disabled: colorMode.preference === 'dark' },
  ]])
}
