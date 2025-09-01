<script setup lang='ts'>
const emit = defineEmits<(e: 'refresh')=> void>()

const model = defineModel<boolean>()

const { selectedLdap } = await useLdapConnection()
const toast = useToast()
const { start, finish, isLoading } = useLoadingIndicator()

const uploadedFile = ref<File>()

async function importLdap(){
  start()

  const formData = new FormData()
  formData.append('id', String(selectedLdap.value))
  if(uploadedFile.value) formData.append('file', uploadedFile.value)

  const res = await $fetch<{ message: string }>('/server/ldap-import', { method: 'put', body: formData })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  emit('refresh')
  model.value = false
}

whenever(() => !model.value, () => uploadedFile.value = undefined)
</script>

<template>
  <UModal v-model:open="model" title="Import DN" :ui="{ footer: 'justify-end' }">
    <template #body>
      <div class="flex items-center justify-center">
        <UFileUpload v-model="uploadedFile" accept=".ldif" label="Drop your ldif here" class="min-h-48 w-96" />
      </div>
    </template>
    <template #footer="{ close }">
      <UButton label="Cancel" color="neutral" variant="outline" :disabled="isLoading" @click="close" />
      <UButton label="Submit" color="neutral" :loading="isLoading" @click="importLdap" />
    </template>
  </UModal>
</template>
