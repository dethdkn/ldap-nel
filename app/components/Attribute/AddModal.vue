<script setup lang='ts'>
const emit = defineEmits<(e: 'refresh')=> void>()

const { addModal, addState } = useCrudModal()
const toast = useToast()
const { start, finish, isLoading } = useLoadingIndicator()
const { selectedLdap } = await useLdapConnection()

const possibleAttributes = ref<string[]>([])
const uploadedFile = ref<File>()
const encryptionMethod = ref(passwordEncryptionMethods[0] || 'Plain Text')

whenever(uploadedFile, async f => addState.value.value = await imageToB64(f))

whenever(() => addState.value.dn, async () => {
  const { possibleAttributes: pa } = await $fetch<{ possibleAttributes: string[] }>('/server/possible-attributes', { method: 'post', body: { id: selectedLdap.value, dn: addState.value.dn } })
    .catch(() => ({ possibleAttributes: [] }))
  possibleAttributes.value = pa
})

function createAttr(item: string){
  possibleAttributes.value.push(item)
  possibleAttributes.value.sort()
  addState.value.attribute = item
}

async function addAttributeValue(){
  start()

  let state = addState.value
  if(state.attribute === 'userPassword') state = { ...state, value: `${encryptionMethod.value}:${state.value}` }

  const body = attributeValueSchema.safeParse(state)
  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/attribute-value', { method: 'post', body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  emit('refresh')
  addModal.value = false
}

whenever(() => !addModal.value, () => {
  possibleAttributes.value = []
  uploadedFile.value = undefined
  encryptionMethod.value = passwordEncryptionMethods[0] || 'Plain Text'
})
</script>

<template>
  <UModal v-model:open="addModal" title="Add Attribute" :ui="{ footer: 'justify-end' }">
    <template #body>
      <p class="text-center">
        {{ addState.dn?.split(',')?.join(', ') || '-' }}
      </p>
      <UForm :schema="attributeValueSchema" :state="addState" class="w-full space-y-4 md:space-y-6" @submit="addAttributeValue">
        <UFormField label="Attribute" name="attribute">
          <UInputMenu v-model="addState.attribute" icon="i-lucide-tags" size="lg" placeholder="Attribute" :items="possibleAttributes" class="w-full" create-item @create="createAttr" />
        </UFormField>
        <UTabs v-if="addState.attribute === 'jpegPhoto'" :items="jpegPhotoTabs" color="neutral">
          <template #text>
            <UFormField label="Value" name="value">
              <UInput v-model="addState.value" icon="i-lucide-scan-text" size="lg" placeholder="Value" class="w-full" />
            </UFormField>
            <div class="flex items-center justify-center px-2 pt-4">
              <img v-if="addState.value" :src="`data:image/jpeg;base64,${addState.value}`" :alt="addState.attribute" class="size-32 rounded-md">
              <div v-else class="size-32 animate-pulse rounded-md bg-gray-300 dark:bg-gray-700" />
            </div>
          </template>
          <template #file>
            <div class="flex items-center justify-center px-2 py-4">
              <UFileUpload v-model="uploadedFile" accept="image/*" label="Drop your image here" class="min-h-48 w-96" />
            </div>
          </template>
        </UTabs>
        <template v-else-if="addState.attribute === 'userPassword'">
          <UFormField label="Encryption Method" name="encryptionMethod">
            <USelect v-model="encryptionMethod" icon="i-lucide-lock" :items="passwordEncryptionMethods" size="lg" class="w-full" />
          </UFormField>
          <UFormField label="Password" name="value">
            <UInput v-model="addState.value" icon="i-lucide-key-round" size="lg" placeholder="Value" :type="encryptionMethod === 'Plain Text' ? 'text' : 'password'" class="w-full" />
          </UFormField>
        </template>
        <UFormField v-else label="Value" name="value">
          <UInput v-model="addState.value" icon="i-lucide-scan-text" size="lg" placeholder="Value" class="w-full" />
        </UFormField>
      </UForm>
    </template>
    <template #footer="{ close }">
      <UButton label="Cancel" color="neutral" variant="outline" :disabled="isLoading" @click="close" />
      <UButton label="Submit" color="neutral" :loading="isLoading" @click="addAttributeValue" />
    </template>
  </UModal>
</template>
