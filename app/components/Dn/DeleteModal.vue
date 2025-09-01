<script setup lang='ts'>
const emit = defineEmits<(e: 'refresh')=> void>()

const { deleteDnModal, deleteDnState } = useCrudModal()
const toast = useToast()
const { start, finish, isLoading } = useLoadingIndicator()

async function deleteDn(){
  start()

  const res = await $fetch<{ message: string }>('/server/dn', { method: 'delete', body: deleteDnState.value })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  emit('refresh')
  deleteDnModal.value = false
}
</script>

<template>
  <UModal v-model:open="deleteDnModal" title="Delete DN" :ui="{ footer: 'justify-between' }">
    <template #body>
      <div class="text-center break-words">
        <p>Are you sure you want to delete <strong> {{ deleteDnState.dn }}</strong>?</p>
      </div>
    </template>
    <template #footer="{ close }">
      <UCheckbox v-model="deleteDnState.smart" label="Smart Delete" />
      <div class="flex items-center justify-center space-x-1">
        <UButton label="Cancel" color="neutral" variant="outline" :disabled="isLoading" @click="close" />
        <UButton label="Delete" color="error" :loading="isLoading" @click="deleteDn" />
      </div>
    </template>
  </UModal>
</template>
