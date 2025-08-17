<script setup lang='ts'>
const emit = defineEmits<(e: 'refresh')=> void>()
const model = defineModel<boolean>()

const toast = useToast()
const { start, finish, isLoading } = useLoadingIndicator()

const state = ref<Ldap>({ id: 0, name: '', url: '', port: 389, ssl: false, base_dn: '', bind_dn: '', bind_pass: '' })
async function createLdap(){
  start()

  const body = ldapSchema.safeParse(state.value)
  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string }>('/server/ldap', { method: 'post', body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  emit('refresh')
  model.value = false
  state.value = { id: 0, name: '', url: '', port: 389, ssl: false, base_dn: '', bind_dn: '', bind_pass: '' }
}

whenever(() => !model.value, () => state.value = { id: 0, name: '', url: '', port: 389, ssl: false, base_dn: '', bind_dn: '', bind_pass: '' })
</script>

<template>
  <UModal v-model:open="model" title="Create Ldap Connection" :ui="{ footer: 'justify-end' }">
    <template #body>
      <UForm :schema="ldapSchema" :state class="w-full space-y-4 md:space-y-6" @submit="createLdap">
        <UFormField label="Name" name="name">
          <UInput v-model="state.name" icon="i-lucide-star" size="lg" placeholder="My Connection" class="w-full" />
        </UFormField>
        <UFormField label="URL" name="url">
          <UInput v-model="state.url" icon="i-lucide-link" size="lg" placeholder="example.com" class="w-full" />
        </UFormField>
        <UFormField label="Port" name="port">
          <UInputNumber v-model="state.port" icon="i-lucide-ethernet-port" orientation="vertical" size="lg" placeholder="389" class="w-full" />
        </UFormField>
        <UFormField name="ssl">
          <UCheckbox v-model="state.ssl" label="SSL" name="ssl" variant="card" class="w-full" />
        </UFormField>
        <UFormField label="Base DN" name="base_dn">
          <UInput v-model="state.base_dn" icon="i-lucide-house" size="lg" placeholder="dc=example,dc=com" class="w-full" />
        </UFormField>
        <UFormField label="Bind DN" name="bind_dn">
          <UInput v-model="state.bind_dn" icon="i-lucide-user" size="lg" class="w-full" />
        </UFormField>
        <UFormField label="Bind Password" name="bind_pass">
          <UInput v-model="state.bind_pass" icon="i-lucide-key-round" size="lg" type="password" class="w-full" />
        </UFormField>
      </UForm>
    </template>
    <template #footer="{ close }">
      <UButton label="Cancel" color="neutral" variant="outline" :disabled="isLoading" @click="close" />
      <UButton label="Test & Save" color="neutral" :loading="isLoading" @click="createLdap" />
    </template>
  </UModal>
</template>
