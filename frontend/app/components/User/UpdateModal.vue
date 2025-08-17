<script setup lang='ts'>
const emit = defineEmits<(e: 'refresh')=> void>()
const model = defineModel<boolean>()
const state = defineModel<User>('state', { required: true })

const toast = useToast()
const { start, finish, isLoading } = useLoadingIndicator()

async function updateUser(){
  start()

  const body = updateUserSchema.safeParse(state.value)
  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/user', { method: 'put', body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  emit('refresh')
  model.value = false
  state.value = { id: 0, username: '', password: '', repeatPassword: '', admin: false }
}

whenever(() => !model.value, () => state.value = { id: 0, username: '', password: '', repeatPassword: '', admin: false })
</script>

<template>
  <UModal v-model:open="model" title="Update User" :ui="{ footer: 'justify-end' }">
    <template #body>
      <UForm :schema="userSchema" :state class="w-full space-y-4 md:space-y-6" @submit="updateUser">
        <UFormField label="Username" name="username">
          <UInput v-model="state.username" icon="i-lucide-user" size="lg" class="w-full" />
        </UFormField>
        <UFormField label="Password" name="password">
          <UInput v-model="state.password" icon="i-lucide-key-round" size="lg" type="password" class="w-full" />
        </UFormField>
        <UFormField label="Repeat Password" name="repeatPassword">
          <UInput v-model="state.repeatPassword" icon="i-lucide-shield-check" size="lg" type="password" class="w-full" />
        </UFormField>
        <UFormField name="admin">
          <UCheckbox v-model="state.admin" label="Admin" name="admin" variant="card" class="w-full" />
        </UFormField>
      </UForm>
    </template>
    <template #footer="{ close }">
      <UButton label="Cancel" color="neutral" variant="outline" :disabled="isLoading" @click="close" />
      <UButton label="Submit" color="neutral" :loading="isLoading" @click="updateUser" />
    </template>
  </UModal>
</template>
