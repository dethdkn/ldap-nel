<script setup lang='ts'>
const props = defineProps({
  baseDn: { type: String, required: true },
})
const emit = defineEmits<(e: 'refresh', payload: string)=> void>()
const model = defineModel<boolean>()

const { selectedLdap } = await useLdapConnection()
const toast = useToast()
const { start, finish, isLoading } = useLoadingIndicator()

const state = ref<NewDn>({
  id: 0,
  dn: '',
  attributes: [],
})

const fullDn = computed(() => `${state.value.dn},${props.baseDn}`)

const currentState = ref({ attribute: '', val: '' })

const AttributeModal = ref(false)
const possibleAttributes = ref<string[]>([])
const uploadedFile = ref<File>()
const encryptionMethod = ref(passwordEncryptionMethods[0] || 'Plain Text')

whenever(uploadedFile, async f => currentState.value.val = await imageToB64(f))

function createAttr(item: string){
  possibleAttributes.value.push(item)
  possibleAttributes.value.sort()
  currentState.value.attribute = item
}

function addAttribute(){
  if(!currentState.value.attribute || !currentState.value.val){
    toast.add({ title: 'Attribute and Value are required', icon: 'i-lucide-shield-alert', color: 'error' })
    return
  }

  let { val } = currentState.value
  const { attribute } = currentState.value
  if(attribute === 'userPassword') val = `${encryptionMethod.value}:${val}`

  state.value.attributes.push({ attribute, value: val })
  state.value.attributes.sort((a, b) => a.attribute.localeCompare(b.attribute))

  currentState.value = { attribute: '', val: '' }
  AttributeModal.value = false
}

async function addDn(){
  start()

  const body = newDnSchema.safeParse({ ...state.value, id: selectedLdap.value || 0, dn: fullDn.value })
  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/dn', { method: 'post', body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  emit('refresh', fullDn.value)
  model.value = false
}

whenever(() => !model.value, () => {
  state.value = { id: 0, dn: '', attributes: [] }
  currentState.value = { attribute: '', val: '' }
  possibleAttributes.value = []
  uploadedFile.value = undefined
  encryptionMethod.value = passwordEncryptionMethods[0] || 'Plain Text'
})
</script>

<template>
  <UModal v-model:open="model" title="Add DN" :ui="{ footer: 'justify-end' }">
    <template #body>
      <div class="flex flex-col items-center justify-center space-y-3">
        <UFieldGroup>
          <UInput v-model="state.dn" :icon="getLdapIcon(state.dn, true)" size="lg" placeholder="uid=example,ou=people" class="w-full" />
          <UBadge color="neutral" variant="outline" size="lg" :label="`,${baseDn}`" />
        </UFieldGroup>
        <div class="mt-6 flex w-full items-center justify-end">
          <UButton label="Add Attribute" color="neutral" :disabled="isLoading" @click="AttributeModal = true" />
        </div>
        <table class="w-full table-fixed text-left text-sm text-gray-500 dark:text-gray-400">
          <thead class="bg-gray-200 text-xs text-gray-700 uppercase dark:bg-gray-700 dark:text-gray-400">
            <tr>
              <th scope="col" class="rounded-tl-md px-6 py-3">
                Attribute
              </th>
              <th scope="col" class="px-6 py-3">
                Value
              </th>
              <th scope="col" class="px-6 py-3 text-end">
                Size
              </th>
              <th scope="col" class="rounded-tr-md px-6 py-3 text-end">
                Actions
              </th>
            </tr>
          </thead>
          <tbody>
            <UPopover v-for="({attribute, value}, key) in state.attributes" :key="`${attribute}-${key}`" :mode="attribute === 'jpegPhoto' ? 'hover' : 'click'" :content="{side: 'top'}">
              <tr class="bg-gray-50 hover:bg-gray-300 dark:border-gray-700 dark:bg-gray-800 dark:hover:bg-gray-600">
                <td class="px-6 py-4">
                  {{ attribute }}
                </td>
                <td class="max-w-20 truncate px-6 py-4">
                  {{ attribute === 'userPassword' ? '***' : value }}
                </td>
                <td class="px-6 py-4 text-end">
                  {{ attribute === 'userPassword' ? '*' : value.length }}
                </td>
                <td class="px-6 py-4 text-center">
                  <UButton variant="ghost" size="xs" class="rounded-full" icon="i-lucide-trash" color="error" @click="state.attributes.splice(key, 1)" />
                </td>
              </tr>

              <template #content>
                <img v-if="attribute === 'jpegPhoto'" :src="`data:image/jpeg;base64,${value}`" :alt="attribute" class="max-w-32 rounded-md">
                <div v-else class=" max-w-64 min-w-20 rounded-md px-4 py-6 text-center text-sm break-words">
                  <p>
                    {{ value }}
                  </p>
                </div>
              </template>
            </UPopover>
          </tbody>
        </table>
      </div>
    </template>
    <template #footer="{ close }">
      <UButton label="Cancel" color="neutral" variant="outline" :disabled="isLoading" @click="close" />
      <UButton label="Submit" color="neutral" :loading="isLoading" @click="addDn" />
    </template>
  </UModal>

  <UModal v-model:open="AttributeModal" title="Add Attribute" :ui="{ footer: 'justify-end' }">
    <template #body>
      <div>
        <UFormField label="Attribute" name="attribute" class="mb-2">
          <UInputMenu v-model="currentState.attribute" icon="i-lucide-tags" size="lg" placeholder="Attribute" :items="possibleAttributes" class="w-full" create-item @create="createAttr" />
        </UFormField>
        <UTabs v-if="currentState.attribute === 'jpegPhoto'" :items="jpegPhotoTabs" color="neutral">
          <template #text>
            <UFormField label="Value" name="value">
              <UInput v-model="currentState.val" icon="i-lucide-scan-text" size="lg" placeholder="Value" class="w-full" />
            </UFormField>
            <div class="flex items-center justify-center px-2 pt-4">
              <img v-if="currentState.val" :src="`data:image/jpeg;base64,${currentState.val}`" :alt="currentState.attribute" class="size-32 rounded-md">
              <div v-else class="size-32 animate-pulse rounded-md bg-gray-300 dark:bg-gray-700" />
            </div>
          </template>
          <template #file>
            <div class="flex items-center justify-center px-2 py-4">
              <UFileUpload v-model="uploadedFile" accept="image/*" label="Drop your image here" class="min-h-48 w-96" />
            </div>
          </template>
        </UTabs>
        <template v-else-if="currentState.attribute === 'userPassword'">
          <UFormField label="Encryption Method" name="encryptionMethod">
            <USelect v-model="encryptionMethod" icon="i-lucide-lock" :items="passwordEncryptionMethods" size="lg" class="w-full" />
          </UFormField>
          <UFormField label="Password" name="value">
            <UInput v-model="currentState.val" icon="i-lucide-key-round" size="lg" placeholder="Value" :type="encryptionMethod === 'Plain Text' ? 'text' : 'password'" class="w-full" />
          </UFormField>
        </template>
        <UFormField v-else label="Value" name="value">
          <UInput v-model="currentState.val" icon="i-lucide-scan-text" size="lg" placeholder="Value" class="w-full" />
        </UFormField>
      </div>
    </template>
    <template #footer="{ close }">
      <UButton label="Cancel" color="neutral" variant="outline" :disabled="isLoading" @click="close" />
      <UButton label="Add" color="neutral" :loading="isLoading" @click="addAttribute" />
    </template>
  </UModal>
</template>
