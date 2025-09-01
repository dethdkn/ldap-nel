<script setup lang='ts'>
const emit = defineEmits<(e: 'refresh')=> void>()

const { deleteModal, deleteState } = useCrudModal()
const toast = useToast()
const { start, finish, isLoading } = useLoadingIndicator()

async function deleteAttributeValue(){
  start()

  const body = attributeValueSchema.safeParse(deleteState.value)
  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/attribute-value', { method: 'delete', body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  emit('refresh')
  deleteModal.value = false
}
</script>

<template>
  <UModal v-model:open="deleteModal" title="Delete Attribute Value" :ui="{ footer: 'justify-end' }">
    <template #body>
      <div class="text-center break-words">
        <p>Are you sure you want to delete</p>
        <p><strong>{{ deleteState.attribute }}: {{ deleteState.value }}</strong></p>
        <p>from <strong>{{ deleteState.dn }}</strong>?</p>
      </div>
    </template>
    <template #footer="{ close }">
      <UButton label="Cancel" color="neutral" variant="outline" :disabled="isLoading" @click="close" />
      <UButton label="Delete" color="error" :loading="isLoading" @click="deleteAttributeValue" />
    </template>
  </UModal>
</template>
