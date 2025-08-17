<script setup lang='ts'>
import type { TreeItem } from '@nuxt/ui'

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
  const { dn, childs } = await $fetch<{ childs: string[], dn: string }>('/server/ldap-childs', { method: 'post', body: { id: selectedLdap.value } })

  items.value = buildLdapTree(dn, childs, onSelect)

  expanded.value = [dn]
}

async function refreshAttributes(){
  const res = await $fetch<{ attributes: Record<string, string[]> }>('/server/ldap-attributes', { method: 'post', body: { id: selectedLdap.value, dn: selected.value?.fullDn } })
  attributes.value = res.attributes
}

watch(selectedLdap, refreshTree, { immediate: true })
watch(selected, refreshAttributes)
</script>

<template>
  <div class="space-between h-ldap max-h-ldap flex w-full justify-between space-x-4 overflow-x-scroll">
    <div class="max-h-ldap min-w-64 overflow-scroll">
      <UTree v-model="selected" v-model:expanded="expanded" :items />
    </div>
    <div class="max-h-ldap w-full min-w-92 overflow-scroll">
      <table class="w-full text-left text-sm text-gray-500 dark:text-gray-400">
        <thead class="bg-gray-200 text-xs text-gray-700 uppercase dark:bg-gray-700 dark:text-gray-400">
          <tr>
            <th scope="col" class="rounded-tl-md px-6 py-3">
              Attribute
            </th>
            <th scope="col" class="px-6 py-3">
              Value
            </th>
            <th scope="col" class="rounded-tr-md px-6 py-3">
              Size
            </th>
          </tr>
        </thead>
        <tbody>
          <template v-for="(attrK, key) in Object.keys(attributes)" :key="`${attrK}-${key}`">
            <template v-for="(val, k) in attributes[attrK]" :key="`${val}-${k}`">
              <tr class="bg-gray-50 hover:bg-gray-300 dark:border-gray-700 dark:bg-gray-800 dark:hover:bg-gray-600">
                <td class="px-6 py-4">
                  {{ attrK }}
                </td>
                <td class="px-6 py-4">
                  {{ val }}
                </td>
                <td class="px-6 py-4">
                  {{ val.length }}
                </td>
              </tr>
            </template>
          </template>
        </tbody>
      </table>
    </div>
  </div>
</template>
