<script setup lang="ts">
import type { CommandPaletteGroup, CommandPaletteItem, TreeItem } from '@nuxt/ui'

const props = defineProps({
  items: { type: Array as PropType<TreeItem[]>, required: true },
})
const emit = defineEmits<(e: 'searched', payload: string)=> void>()
const model = defineModel<boolean>()

const searchTerm = ref('')

function flatten(items: TreeItem[]): TreeItem[]{
  const result: TreeItem[] = []

  const traverse = (node: TreeItem) => {
    result.push({ label: node.label, icon: node.icon || 'i-lucide-folder', onSelect: () => {
      emit('searched', node.fullDn || '')
      model.value = false
    } })

    if(node.children) for(const child of node.children) traverse(child)
  }

  for(const item of items) traverse(item)

  return result
}

const groups = computed(() => {
  const g: CommandPaletteGroup<CommandPaletteItem>[] = [{
    id: 'entrys',
    label: searchTerm.value ? `Entrys matching “${searchTerm.value}”...` : 'Entrys',
    items: flatten(props.items[0]?.children || []).filter(i => i.label?.toLowerCase()?.includes(searchTerm.value.toLowerCase())),
    ignoreFilter: true,
  }]
  return g
})

whenever(() => !model.value, () => searchTerm.value = '')
</script>

<template>
  <UModal v-model:open="model">
    <template #content>
      <UCommandPalette
        v-model:search-term="searchTerm"
        :groups
        placeholder="Search entrys..."
        class="h-80"
      />
    </template>
  </UModal>
</template>
