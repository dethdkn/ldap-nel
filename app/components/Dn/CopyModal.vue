<script setup lang='ts'>
import type { TreeItem } from '@nuxt/ui'

const props = defineProps({
  items: { type: Array as PropType<TreeItem[]>, required: true },
})
const emit = defineEmits<(e: 'refresh')=> void>()

const { copyModal, copyState } = useCrudModal()
const toast = useToast()
const { start, finish, isLoading } = useLoadingIndicator()

const search = ref('')

function flatten(items: TreeItem[]): TreeItem[]{
  const result: TreeItem[] = []

  const traverse = (node: TreeItem) => {
    result.push({ label: node.value, value: node.value })

    if(node.children) for(const child of node.children) traverse(child)
  }

  for(const item of items) traverse(item)

  return result
}

const flatItems = computed(() => flatten(props.items[0]?.children || []).filter(i => i?.label?.toLowerCase()?.includes(search.value.toLowerCase())).slice(0, 100))

async function copyDn(){
  start()

  const body = CopyMoveSchema.safeParse({ ...copyState.value, targetDn: `${copyState.value.dn.split(',')[0]},${copyState.value.targetDn}` })
  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/copy-dn', { method: 'post', body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  emit('refresh')
  copyModal.value = false
}
</script>

<template>
  <UModal v-model:open="copyModal" title="Copy DN" :ui="{ footer: 'justify-end' }">
    <template #body>
      <div class="space-y-4 text-center">
        <p>Copy to</p>
        <UForm :schema="CopyMoveSchema" :state="copyState" class="w-full space-y-4 md:space-y-6" @submit="copyDn">
          <UFieldGroup>
            <UBadge color="neutral" variant="outline" size="lg" :label="`${copyState.dn.split(',')[0]},`" />
            <USelectMenu v-model="copyState.targetDn" v-model:search-term="search" :items="flatItems" :multiple="false" ignore-filter label-key="label" value-key="value" search-input size="lg" class="w-full min-w-36" />
          </UFieldGroup>
        </UForm>
      </div>
    </template>
    <template #footer="{ close }">
      <UButton label="Cancel" color="neutral" variant="outline" :disabled="isLoading" @click="close" />
      <UButton label="Submit" color="neutral" :loading="isLoading" @click="copyDn" />
    </template>
  </UModal>
</template>
