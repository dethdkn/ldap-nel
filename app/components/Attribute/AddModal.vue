<script setup lang='ts'>
const emit = defineEmits<(e: 'refresh')=> void>()

const { addModal, addState } = useCrudModal()
const toast = useToast()
const { start, finish, isLoading } = useLoadingIndicator()
const { selectedLdap } = await useLdapConnection()

const possibleAttributes = ref<string[]>([])

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

  const body = attributeValueSchema.safeParse(addState.value)
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
        <UFormField label="Value" name="value">
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
