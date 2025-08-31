import type { DropdownMenuItem } from '@nuxt/ui'

export default () => {
  const { user, clearUserSession } = useUserSession()
  const passwordModal = usePasswordModal()
  const settingsModal = useSettingsModal()
  const aboutModal = useAboutModal()

  const logout = () => {
    clearUserSession()
    return navigateTo('/login')
  }

  return computed(() => {
    const menu: DropdownMenuItem[][] = [
      [{ label: user.value.username || '???', icon: 'i-lucide-user', disabled: true }],
      [{ icon: 'i-lucide-key-round', label: 'Update Password', onSelect: () => passwordModal.value = true }],
    ]

    const baseMenu = [
      { icon: 'i-lucide-info', label: 'About  ', onSelect: () => aboutModal.value = true },
      { icon: 'i-lucide-log-out', label: 'Logout', onSelect: logout },
    ]

    const adminMenu = [
      { icon: 'i-lucide-settings', label: 'Settings  ', onSelect: () => settingsModal.value = true },
      ...baseMenu,
    ]

    if(user.value.admin) menu[1]?.push(...adminMenu)
    else menu[1]?.push(...baseMenu)

    return menu
  })
}
