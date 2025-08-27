<script setup lang='ts'>
const emit = defineEmits<(e: 'refresh')=> void>()

const { updateModal, updateState } = useCrudModal()
const toast = useToast()
const { start, finish, isLoading } = useLoadingIndicator()

async function updateAttributeValue(){
  start()

  const body = attributeNewValueSchema.safeParse(updateState.value)
  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/attribute-value', { method: 'put', body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  emit('refresh')
  updateModal.value = false
}
</script>

<template>
  <UModal v-model:open="updateModal" title="Update Attribute Value" :ui="{ footer: 'justify-end' }">
    <template #body>
      <UForm :schema="attributeNewValueSchema" :state="updateState" class="w-full space-y-4 md:space-y-6" @submit="updateAttributeValue">
        <UFormField label="New Value" name="newValue">
          <UInput v-model="updateState.newValue" icon="i-lucide-scan-text" size="lg" placeholder="Value" class="w-full" />
        </UFormField>
      </UForm>
    </template>
    <template #footer="{ close }">
      <UButton label="Cancel" color="neutral" variant="outline" :disabled="isLoading" @click="close" />
      <UButton label="Submit" color="neutral" :loading="isLoading" @click="updateAttributeValue" />
    </template>
  </UModal>
</template>
