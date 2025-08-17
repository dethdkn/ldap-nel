<script setup lang='ts'>
const emit = defineEmits<(e: 'refresh')=> void>()
const model = defineModel<boolean>()
const state = defineModel<User>('state', { required: true })

const validation1 = ref('')
const validation2 = ref('')

const toast = useToast()
const { start, finish, isLoading } = useLoadingIndicator()

async function deleteUser(){
  start()

  if(validation1.value !== state.value.username || validation2.value !== 'delete user'){
    toast.add({ title: 'Validation failed', icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const body = updateUserSchema.safeParse(state.value)
  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/user', { method: 'delete', body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  emit('refresh')
  model.value = false
  state.value = { id: 0, username: '', password: '', repeatPassword: '', admin: false }
}

whenever(() => !model.value, () => {
  state.value = { id: 0, username: '', password: '', repeatPassword: '', admin: false }
  validation1.value = ''
  validation2.value = ''
})
</script>

<template>
  <UModal v-model:open="model" title="Delete User" :ui="{ footer: 'justify-end' }">
    <template #body>
      <div class="w-full space-y-4 md:space-y-6">
        <UFormField :label="`Enter the username ${state.username} to continue:`">
          <UInput v-model="validation1" icon="i-lucide-user-search" size="lg" class="w-full" />
        </UFormField>
        <UFormField label="To verify, type delete user below:">
          <UInput v-model="validation2" icon="i-lucide-shield-alert" size="lg" class="w-full" />
        </UFormField>
      </div>
    </template>
    <template #footer="{ close }">
      <UButton label="Cancel" color="neutral" variant="outline" :loading="isLoading" @click="close" />
      <UButton label="Delete" color="error" :loading="isLoading" :disabled="!(validation1 === state.username && validation2 === 'delete user')" @click="deleteUser" />
    </template>
  </UModal>
</template>
