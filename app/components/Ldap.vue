<script setup lang='ts'>
import type { DropdownMenuItem, TreeItem } from '@nuxt/ui'

const { start, finish } = useLoadingIndicator()
const { selectedLdap } = await useLdapConnection()

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

  items.value = buildLdapTree(dn, childs, onSelect)

  expanded.value = [dn]
  finish()
}

async function refreshAttributes(){
  start()
  const res = await $fetch<{ attributes: Record<string, string[]> }>('/server/ldap-attributes', { method: 'post', body: { id: selectedLdap.value, dn: selected.value?.fullDn } })
  attributes.value = res.attributes
  finish()
}

watch(selectedLdap, refreshTree, { immediate: true })
watch(selected, refreshAttributes)

async function refreshAll(){
  await refreshTree()
  await refreshAttributes()
}

const options = ref<DropdownMenuItem[][]>([[
  { icon: 'i-lucide-rotate-ccw', label: 'Refresh', onSelect: refreshAll },
  { icon: 'i-lucide-circle-plus', label: 'Add attribute' },
]])
</script>

<template>
  <div class="mb-4 flex w-full items-center justify-end">
    <UDropdownMenu :items="options" size="sm">
      <UIcon name="i-lucide-ellipsis-vertical" class="cursor-pointer transition-all duration-300 hover:text-gray-400 hover:dark:text-gray-600" />
    </UDropdownMenu>
  </div>
  <div class="space-between h-ldap max-h-ldap flex w-full justify-between space-x-4 overflow-x-auto">
    <div class="max-h-ldap min-w-64 overflow-auto">
      <UTree v-model="selected" v-model:expanded="expanded" :items />
    </div>
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
                <UContextMenu :items="$buildContextMenu(attrK, val)">
                  <td class="px-6 py-4">
                    {{ attrK }}
                  </td>
                </UContextMenu>
                <UContextMenu :items="$buildContextMenu(attrK, val)">
                  <td class="max-w-20 truncate px-6 py-4">
                    {{ val }}
                  </td>
                </UContextMenu>
                <UContextMenu :items="$buildContextMenu(attrK, val)">
                  <td class="px-6 py-4 text-end">
                    {{ val.length }}
                  </td>
                </UContextMenu>
              </tr>

              <template #content>
                <img v-if="attrK === 'jpegPhoto'" :src="`data:image/jpeg;base64,${val}`" :alt="attrK" class="max-w-32 rounded-md">
                <div v-else class=" max-w-64 min-w-36 rounded-md px-4 py-6 text-center break-words">
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
</template>
