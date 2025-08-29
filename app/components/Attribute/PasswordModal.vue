<script setup lang='ts'>
const emit = defineEmits<(e: 'refresh')=> void>()

const { passwordModal, passwordState } = useCrudModal()
const toast = useToast()
const { start, finish, isLoading } = useLoadingIndicator()

const encryptionMethod = ref(passwordEncryptionMethods[0] || 'Plain Text')

watch(encryptionMethod, em => {
  if(em === 'Plain Text') passwordState.value.newValue = passwordState.value.value
  else passwordState.value.newValue = ''
})

async function updatePasswordValue(){
  start()

  const body = attributeNewValueSchema.safeParse({ ...passwordState.value, newValue: `${encryptionMethod.value}:${passwordState.value.newValue}` })
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
  passwordModal.value = false
}
</script>

<template>
  <UModal v-model:open="passwordModal" title="Update Password" :ui="{ footer: 'justify-end' }">
    <template #body>
      <UForm :schema="attributeNewValueSchema" :state="passwordState" class="w-full space-y-4 md:space-y-6" @submit="updatePasswordValue">
        <UFormField label="Encryption Method" name="encryptionMethod">
          <USelect v-model="encryptionMethod" icon="i-lucide-lock" :items="passwordEncryptionMethods" size="lg" class="w-full" />
        </UFormField>
        <UFormField label="New Password" name="newValue">
          <UInput v-model="passwordState.newValue" icon="i-lucide-key-round" size="lg" placeholder="Value" :type="encryptionMethod === 'Plain Text' ? 'text' : 'password'" class="w-full" />
        </UFormField>
      </UForm>
    </template>
    <template #footer="{ close }">
      <UButton label="Cancel" color="neutral" variant="outline" :disabled="isLoading" @click="close" />
      <UButton label="Submit" color="neutral" :loading="isLoading" @click="updatePasswordValue" />
    </template>
  </UModal>
</template>
