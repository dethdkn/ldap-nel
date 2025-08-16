import type { DropdownMenuItem } from '@nuxt/ui'

export default () => {
  const { user, clearUserSession } = useUserSession()
  const passwordModal = usePasswordModal()
  const { modal: settingsModal } = useSettingsModal()
  const aboutModal = useAboutModal()

  const logout = () => {
    clearUserSession()
    navigateTo('/login')
  }

  return computed(() => {
    const LogoutItem = { icon: 'i-lucide-log-out', label: 'Logout', onSelect: logout }

    const adminMenu = [
      { icon: 'i-lucide-settings', label: 'Settings  ', onSelect: () => settingsModal.value = true },
      { icon: 'i-lucide-info', label: 'About  ', onSelect: () => aboutModal.value = true },
      LogoutItem,
    ]

    const menu: DropdownMenuItem[][] = [
      [{ label: user.value.username || '???', icon: 'i-lucide-user', disabled: true }],
      [{ icon: 'i-lucide-key', label: 'Update Password', onSelect: () => passwordModal.value = true }],
    ]

    if(user.value.admin) menu[1]?.push(...adminMenu)
    else menu[1]?.push(LogoutItem)

    return menu
  })
}
