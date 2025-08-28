<script setup lang='ts'>
import type { TableColumn } from '@nuxt/ui'

const { refreshLdaps } = await useLdapConnection()

const createLdapModal = ref(false)

const updateLdapModal = ref(false)
const updateLdapState = ref<Ldap>({ id: 0, name: '', url: '', port: 389, ssl: false, base_dn: '', bind_dn: '', bind_pass: '' })
function openUpdateModal(ldap: Ldap){
  updateLdapState.value = { ...ldap }
  updateLdapModal.value = true
}

const deleteLdapModal = ref(false)
const deleteLdapState = ref<Ldap>({ id: 0, name: '', url: '', port: 389, ssl: false, base_dn: '', bind_dn: '', bind_pass: '' })
function openDeleteModal(ldap: Ldap){
  deleteLdapState.value = { ...ldap }
  deleteLdapModal.value = true
}

const UBadge = resolveComponent('UBadge')
const UButton = resolveComponent('UButton')
const UTooltip = resolveComponent('UTooltip')

const { data, refresh } = await useFetch<Ldap[]>('/server/ldaps')

async function globalRefresh(){
  await refreshLdaps()
  await refresh()
}

const columns: TableColumn<Ldap>[] = [
  { accessorKey: 'id', header: 'ID', cell: ({ row }) => h(UBadge, { variant: 'subtle', color: 'neutral' }, () => row.getValue('id')) },
  { accessorKey: 'name', header: 'Name', cell: ({ row }) => h(UBadge, { variant: 'subtle', color: 'primary' }, () => row.getValue('name')) },
  { accessorKey: 'url', header: 'Url', cell: ({ row }) => h(UBadge, { variant: 'subtle', color: 'secondary' }, () => row.getValue('url')) },
  { accessorKey: 'port', header: 'Port', cell: ({ row }) => h(UBadge, { variant: 'subtle', color: 'warning' }, () => row.getValue('port')) },
  { accessorKey: 'ssl', header: 'SSL', cell: ({ row }) => {
    const level = row.getValue('ssl') ? 'Enabled' : 'Disabled'
    const color = row.getValue('ssl') ? 'primary' : 'error'
    return h(UBadge, { variant: 'subtle', color }, () => level)
  } },
  { accessorKey: 'base_dn', header: 'Bind DN', cell: ({ row }) => h(UBadge, { variant: 'subtle', color: 'secondary' }, () => row.getValue('base_dn')) },
  { accessorKey: 'bind_dn', header: 'Bind DN', cell: ({ row }) => h(UBadge, { variant: 'subtle', color: 'warning' }, () => row.getValue('bind_dn') || '-') },
  { accessorKey: 'actions', header: 'Actions', cell: ({ row }) => {
    const buttonProps = { variant: 'subtle', size: 'xs', class: 'rounded-full' }
    const updateButton = h(UButton, { icon: 'i-lucide-folder-pen', color: 'secondary', ...buttonProps, onClick: () => openUpdateModal(row.original) }, () => '')
    const deleteButton = h(UButton, { icon: 'i-lucide-folder-x', color: 'error', ...buttonProps, onClick: () => openDeleteModal(row.original) }, () => '')
    const tooltipProps = { delay: 0, content: { side: 'top' } }
    const updateTooltip = h(UTooltip, { text: `Update ${row.getValue('name')}`, ...tooltipProps }, () => updateButton)
    const deleteTooltip = h(UTooltip, { text: `Delete ${row.getValue('name')}`, ...tooltipProps }, () => deleteButton)
    return h('div', { class: 'flex items-center space-x-2' }, [updateTooltip, deleteTooltip])
  } },
]
</script>

<template>
  <div class="my-4 flex items-center justify-end">
    <UButton label="Create Ldap Connection" icon="i-lucide-folder-plus" color="neutral" variant="outline" size="sm" @click="createLdapModal = true" />
  </div>
  <UTable :data :columns />
  <LdapCreateModal v-model="createLdapModal" @refresh="globalRefresh" />
  <LdapUpdateModal v-model="updateLdapModal" v-model:state="updateLdapState" @refresh="globalRefresh" />
  <LdapDeleteModal v-model="deleteLdapModal" v-model:state="deleteLdapState" @refresh="globalRefresh" />
</template>
