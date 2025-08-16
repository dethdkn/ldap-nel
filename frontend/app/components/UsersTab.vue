<script setup lang='ts'>
import type { TableColumn } from '@nuxt/ui'

const createUserModal = ref(false)

const UBadge = resolveComponent('UBadge')
const UButton = resolveComponent('UButton')

const { data, refresh } = await useFetch<User[]>('/server/users')

const columns: TableColumn<User>[] = [
  { accessorKey: 'id', header: 'ID', cell: ({ row }) => h(UBadge, { variant: 'subtle', color: 'neutral' }, () => row.getValue('id')) },
  { accessorKey: 'username', header: 'Username', cell: ({ row }) => h(UBadge, { variant: 'subtle', color: 'primary' }, () => row.getValue('username')) },
  { accessorKey: 'admin', header: 'Level', cell: ({ row }) => {
    const level = row.getValue('admin') ? 'Admin' : 'User'
    const color = row.getValue('admin') ? 'error' : 'secondary'
    return h(UBadge, { variant: 'subtle', color }, () => level)
  } },
  { accessorKey: 'actions', header: 'Actions', cell: () => {
    const buttonProps = { variant: 'subtle', size: 'xs', class: 'rounded-full' }
    const updateButton = h(UButton, { icon: 'i-lucide-user-pen', color: 'secondary', ...buttonProps }, () => '')
    const deleteButton = h(UButton, { icon: 'i-lucide-user-x', color: 'error', ...buttonProps }, () => '')
    return h('div', { class: 'flex items-center space-x-2' }, [updateButton, deleteButton])
  } },
]
</script>

<template>
  <div class="my-4 flex items-center justify-end">
    <UButton label="Create User" icon="i-lucide-user-plus" color="neutral" variant="outline" size="sm" @click="createUserModal = true" />
  </div>
  <UTable :data :columns />
  <UserCreateModal v-model="createUserModal" @refresh="refresh" />
</template>
