<script setup lang='ts'>
import type { TreeItem } from '@nuxt/ui'

const { selectedLdap } = await useLdapConnection()

const items = ref<TreeItem[]>([])
const selected = ref()
const expanded = ref<string[]>([])

// eslint-disable-next-line typescript/no-explicit-any
function onSelect(e: any){
  if(selected.value?.label === e?.detail?.value?.label) e?.preventDefault()
}

watch(selectedLdap, async () => {
  const { dn, childs } = await $fetch<{ childs: string[], dn: string }>('/server/ldap-childs', { method: 'post', body: { id: selectedLdap.value } })

  items.value = buildLdapTree(dn, childs, onSelect)

  expanded.value = [dn]
}, { immediate: true })
</script>

<template>
  <div class="space-between h-ldap max-h-ldap flex w-full justify-between space-x-4">
    <div class="max-h-ldap overflow-scroll md:min-w-64">
      <UTree v-model="selected" v-model:expanded="expanded" :items />
    </div>
    <div class="max-h-ldap w-full overflow-scroll bg-red-500">
      sadasd {{ selected }}
    </div>
  </div>
</template>
