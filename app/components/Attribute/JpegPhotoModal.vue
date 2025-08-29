<script setup lang='ts'>
const emit = defineEmits<(e: 'refresh')=> void>()

const { jpegPhotoModal, jpegPhotoState } = useCrudModal()
const toast = useToast()
const { start, finish, isLoading } = useLoadingIndicator()

const uploadedFile = ref<File>()

whenever(uploadedFile, async f => jpegPhotoState.value.newValue = await imageToB64(f))

async function updateJpegPhotoValue(){
  start()

  const body = attributeNewValueSchema.safeParse(jpegPhotoState.value)
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
  jpegPhotoModal.value = false
}

whenever(() => !jpegPhotoModal.value, () => uploadedFile.value = undefined)
</script>

<template>
  <UModal v-model:open="jpegPhotoModal" title="Update Photo" :ui="{ footer: 'justify-end' }">
    <template #body>
      <UTabs :items="jpegPhotoTabs" color="neutral">
        <template #text>
          <UForm :schema="attributeNewValueSchema" :state="jpegPhotoState" class="w-full space-y-4 md:space-y-6" @submit="updateJpegPhotoValue">
            <UFormField label="New Value" name="newValue">
              <UInput v-model="jpegPhotoState.newValue" icon="i-lucide-scan-text" size="lg" placeholder="Value" class="w-full" />
            </UFormField>
          </UForm>
          <div class="flex items-center justify-center px-2 pt-4">
            <img :src="`data:image/jpeg;base64,${jpegPhotoState.newValue}`" :alt="jpegPhotoState.attribute" class="size-32 rounded-md">
          </div>
        </template>
        <template #file>
          <div class="flex items-center justify-center px-2 py-4">
            <UFileUpload v-model="uploadedFile" accept="image/*" label="Drop your image here" class="min-h-48 w-96" />
          </div>
        </template>
      </UTabs>
    </template>
    <template #footer="{ close }">
      <UButton label="Cancel" color="neutral" variant="outline" :disabled="isLoading" @click="close" />
      <UButton label="Submit" color="neutral" :loading="isLoading" @click="updateJpegPhotoValue" />
    </template>
  </UModal>
</template>
