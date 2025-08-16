<script setup lang='ts'>
import type { TableColumn } from '@nuxt/ui'

const { user } = useUserSession()

const createUserModal = ref(false)

const updateUserModal = ref(false)
const updateUserState = ref<User>({ id: 0, username: '', password: '', repeatPassword: '', admin: false })
function openUpdateModal(user: User){
  updateUserState.value = { ...user }
  updateUserModal.value = true
}

const deleteUserModal = ref(false)
const deleteUserState = ref<User>({ id: 0, username: '', password: '', repeatPassword: '', admin: false })
function openDeleteModal(user: User){
  deleteUserState.value = { ...user }
  deleteUserModal.value = true
}

const UBadge = resolveComponent('UBadge')
const UButton = resolveComponent('UButton')
const UTooltip = resolveComponent('UTooltip')

const { data, refresh } = await useFetch<User[]>('/server/users')

const columns: TableColumn<User>[] = [
  { accessorKey: 'id', header: 'ID', cell: ({ row }) => h(UBadge, { variant: 'subtle', color: 'neutral' }, () => row.getValue('id')) },
  { accessorKey: 'username', header: 'Username', cell: ({ row }) => h(UBadge, { variant: 'subtle', color: 'primary' }, () => row.getValue('username')) },
  { accessorKey: 'admin', header: 'Level', cell: ({ row }) => {
    const level = row.getValue('admin') ? 'Admin' : 'User'
    const color = row.getValue('admin') ? 'error' : 'secondary'
    return h(UBadge, { variant: 'subtle', color }, () => level)
  } },
  { accessorKey: 'actions', header: 'Actions', cell: ({ row }) => {
    const buttonProps = { variant: 'subtle', size: 'xs', class: 'rounded-full' }
    const updateButton = h(UButton, { icon: 'i-lucide-user-pen', color: 'secondary', ...buttonProps, onClick: () => openUpdateModal(row.original) }, () => '')
    const deleteButton = h(UButton, { icon: 'i-lucide-user-x', color: 'error', ...buttonProps, onClick: () => openDeleteModal(row.original), disabled: user.value.username === row.getValue('username') }, () => '')
    const tooltipProps = { delay: 0, content: { side: 'top' } }
    const updateTooltip = h(UTooltip, { text: `Update ${row.getValue('username')}`, ...tooltipProps }, () => updateButton)
    const deleteTooltip = h(UTooltip, { text: `Delete ${row.getValue('username')}`, ...tooltipProps }, () => deleteButton)
    return h('div', { class: 'flex items-center space-x-2' }, [updateTooltip, deleteTooltip])
  } },
]
</script>

<template>
  <div class="my-4 flex items-center justify-end">
    <UButton label="Create User" icon="i-lucide-user-plus" color="neutral" variant="outline" size="sm" @click="createUserModal = true" />
  </div>
  <UTable :data :columns />
  <UserCreateModal v-model="createUserModal" @refresh="refresh" />
  <UserUpdateModal v-model="updateUserModal" v-model:state="updateUserState" @refresh="refresh" />
  <UserDeleteModal v-model="deleteUserModal" v-model:state="deleteUserState" @refresh="refresh" />
</template>
