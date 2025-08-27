import type { ContextMenuItem } from '@nuxt/ui'

export default defineNuxtPlugin(() => {
  const toast = useToast()
  const { copy } = useClipboard()
  const { user } = useUserSession()
  const { openDeleteModal, updateModal } = useCrudModal()

  const clipboard = async (v: string, m: string) => {
    await copy(v)
    toast.add({ title: m, icon: 'i-lucide-badge-check', color: 'success' })
  }

  const buildContextMenu = (ldap: number, dn: string, attr: string, val: string): ContextMenuItem[][] => [
    [
      { icon: 'i-lucide-pen-line', label: 'Edit value', disabled: !user.value.admin, onSelect: () => updateModal.value = true },
      { icon: 'i-lucide-trash', label: 'Delete value', disabled: !user.value.admin, onSelect: () => openDeleteModal(ldap, dn, attr, val) },
    ],
    [
      { icon: 'i-lucide-clipboard', label: 'Copy', onSelect: async () => await clipboard(`${attr}: ${val}`, 'Attribute and Value copied to clipboard') },
      { icon: 'i-lucide-tags', label: 'Copy Attribute Name', onSelect: async () => await clipboard(attr, 'Attribute name copied to clipboard') },
      { icon: 'i-lucide-scan-text', label: 'Copy Value', onSelect: async () => await clipboard(val, 'Value copied to clipboard') },
    ],
  ]

  return { provide: { buildContextMenu } }
})
