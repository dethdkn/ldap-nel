<script setup lang='ts'>
import type { DropdownMenuItem, TreeItem } from '@nuxt/ui'

const { user } = useUserSession()
const { start, finish } = useLoadingIndicator()
const { selectedLdap } = await useLdapConnection()
const { openAddModal, openDeleteDnModal, openCopyModal, openMoveModal } = useCrudModal()

const searchModal = ref(false)
const addDnModal = ref(false)
const importModal = ref(false)

const treeWrapper = ref<HTMLElement>()
const items = ref<TreeItem[]>([])
const attributes = ref<Record<string, string[]>>({})
const selected = ref()
const expanded = ref<string[]>([])

// eslint-disable-next-line typescript/no-explicit-any
function onSelect(e: any){
  if(selected.value?.label === e?.detail?.value?.label) e?.preventDefault()
}

async function refreshTree(){
  start()
  const { dn, childs } = await $fetch<{ childs: string[], dn: string }>('/server/ldap-childs', { method: 'post', body: { id: selectedLdap.value } })
    .catch(error => { throw createError({ statusMessage: error?.data?.message || 'Server Error', fatal: true }) })

  items.value = buildLdapTree(dn, childs, onSelect)

  expanded.value = [dn]
  finish()
}

async function refreshAttributes(){
  start()
  const res = await $fetch<{ attributes: Record<string, string[]> }>('/server/ldap-attributes', { method: 'post', body: { id: selectedLdap.value, dn: selected.value?.fullDn } })
    .catch(error => { throw createError({ statusMessage: error?.data?.message || 'Server Error', fatal: true }) })
  attributes.value = res.attributes
  finish()
}

watch(selectedLdap, refreshTree, { immediate: true })
watch(selected, refreshAttributes)

async function refreshAll(unselect?: boolean){
  if(unselect) selected.value = undefined
  await refreshTree()
  await refreshAttributes()
}

const options = ref<DropdownMenuItem[][]>([[
  { icon: 'i-lucide-search', label: 'Search', kbds: ['meta', 'k'], onSelect: () => searchModal.value = true },
  { icon: 'i-lucide-rotate-ccw', label: 'Refresh', onSelect: () => refreshAll() },
  { icon: 'i-lucide-folder-plus', label: 'Add DN', disabled: !user.value.admin, onSelect: () => addDnModal.value = true },
  { icon: 'i-lucide-folder-minus', label: 'Delete DN', disabled: !user.value.admin, onSelect: () => openDeleteDnModal(selectedLdap.value || 0, selected.value?.fullDn || '') },
  { icon: 'i-lucide-folders', label: 'Copy DN', disabled: !user.value.admin, onSelect: () => openCopyModal(selectedLdap.value || 0, selected.value?.fullDn || '') },
  { icon: 'i-lucide-folder-symlink', label: 'Move DN', disabled: !user.value.admin, onSelect: () => openMoveModal(selectedLdap.value || 0, selected.value?.fullDn || '') },
  { icon: 'i-lucide-circle-plus', label: 'Add attribute', disabled: !user.value.admin, onSelect: () => openAddModal(selectedLdap.value || 0, selected.value?.fullDn || '') },
  { icon: 'i-lucide-upload', label: 'Export', onSelect: () => globalThis.location.href = `/server/ldap-export/${selectedLdap.value}/${selected.value?.fullDn?.replaceAll(' ', '')}` },
  { icon: 'i-lucide-download', label: 'Import', disabled: !user.value.admin, onSelect: () => importModal.value = true },
]])

defineShortcuts({ meta_k: () => searchModal.value = !searchModal.value })

function search(fullDn: string){
  handleSearch(treeWrapper, items.value[0]?.label || '', fullDn)
}

function rightClick(e: Event){
  const target = e.target as HTMLElement
  const button = target.tagName === 'BUTTON' ? target : target.closest('button')
  if(button instanceof HTMLButtonElement) button.click()
}

async function updateSearch(fullDn: string){
  await refreshTree()
  search(fullDn)
}
</script>

<template>
  <div class="mb-4 flex w-full items-center justify-between">
    <div />
    <UBreadcrumb v-if="selected?.fullDn" :items="selected?.fullDn?.split(',')?.toReversed()?.map((dn:string) => ({label: dn, icon: getLdapIcon(dn, true)}))" />
    <UDropdownMenu :items="options" size="sm">
      <UIcon name="i-lucide-ellipsis-vertical" class="cursor-pointer transition-all duration-300 hover:text-gray-400 hover:dark:text-gray-600" />
    </UDropdownMenu>
  </div>
  <div class="space-between h-ldap max-h-ldap flex w-full justify-between space-x-4 overflow-x-auto">
    <UContextMenu :items="options" size="sm">
      <div ref="treeWrapper" class="max-h-ldap min-w-64 overflow-auto">
        <UTree v-model="selected" v-model:expanded="expanded" :items @click.right="rightClick" />
      </div>
    </UContextMenu>
    <div class="max-h-ldap w-full min-w-92 overflow-auto">
      <table class="w-full table-fixed text-left text-sm text-gray-500 dark:text-gray-400">
        <thead class="bg-gray-200 text-xs text-gray-700 uppercase dark:bg-gray-700 dark:text-gray-400">
          <tr>
            <th scope="col" class="rounded-tl-md px-6 py-3">
              Attribute
            </th>
            <th scope="col" class="px-6 py-3">
              Value
            </th>
            <th scope="col" class="rounded-tr-md px-6 py-3 text-end">
              Size
            </th>
          </tr>
        </thead>
        <tbody>
          <template v-for="(attrK, key) in Object.keys(attributes)" :key="`${attrK}-${key}`">
            <UPopover v-for="(val, k) in attributes[attrK]" :key="`${val}-${k}`" :mode="attrK === 'jpegPhoto' ? 'hover' : 'click'" :content="{side: 'top'}">
              <tr class="bg-gray-50 hover:bg-gray-300 dark:border-gray-700 dark:bg-gray-800 dark:hover:bg-gray-600">
                <UContextMenu :items="$buildContextMenu(selectedLdap || 0, selected?.fullDn || '', attrK, val)" size="sm">
                  <td class="px-6 py-4">
                    {{ attrK }}
                  </td>
                </UContextMenu>
                <UContextMenu :items="$buildContextMenu(selectedLdap || 0, selected?.fullDn || '', attrK, val)" size="sm">
                  <td class="max-w-20 truncate px-6 py-4">
                    {{ val }}
                  </td>
                </UContextMenu>
                <UContextMenu :items="$buildContextMenu(selectedLdap || 0, selected?.fullDn || '', attrK, val)" size="sm">
                  <td class="px-6 py-4 text-end">
                    {{ val.length }}
                  </td>
                </UContextMenu>
              </tr>

              <template #content>
                <img v-if="attrK === 'jpegPhoto'" :src="`data:image/jpeg;base64,${val}`" :alt="attrK" class="max-w-32 rounded-md">
                <div v-else class=" max-w-64 min-w-20 rounded-md px-4 py-6 text-center text-sm break-words">
                  <p>
                    {{ val }}
                  </p>
                </div>
              </template>
            </UPopover>
          </template>
        </tbody>
      </table>
    </div>
  </div>
  <SearchModal v-model="searchModal" :items @searched="search" />
  <DnAddModal v-if="user.admin" v-model="addDnModal" :base-dn="items?.[0]?.fullDn || ''" @refresh="updateSearch" />
  <DnCopyModal v-if="user.admin" :items @refresh="() => refreshAll(true)" />
  <DnMoveModal v-if="user.admin" :items @refresh="() => refreshAll(true)" />
  <DnDeleteModal v-if="user.admin" @refresh="() => refreshAll(true)" />
  <ImportModal v-if="user.admin" v-model="importModal" @refresh="refreshAll" />
  <AttributeAddModal v-if="user.admin" @refresh="refreshAttributes" />
  <AttributeUpdateModal v-if="user.admin" @refresh="refreshAttributes" />
  <AttributeJpegPhotoModal v-if="user.admin" @refresh="refreshAttributes" />
  <AttributePasswordModal v-if="user.admin" @refresh="refreshAttributes" />
  <AttributeDeleteModal v-if="user.admin" @refresh="refreshAttributes" />
</template>
