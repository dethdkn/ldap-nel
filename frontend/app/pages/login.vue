<script setup lang="ts">
definePageMeta({
  layout: false,
  middleware: 'setup',
})
useHead({ title: 'Login' })

const toast = useToast()
const themeMenu = useThemeMenu()
const { start, finish, isLoading } = useLoadingIndicator()
const { setUserSession } = useUserSession()

const state = ref<Auth>({ username: '', password: '' })
async function login(){
  start()

  const body = authSchema.safeParse(state.value)
  if(!body.success){
    for(const e of body.error.issues) toast.add({ title: e.message, icon: 'i-lucide-shield-alert', color: 'error' })
    return finish({ error: true })
  }

  const res = await $fetch<{ message: string, token: string, isAdmin: boolean }>('/server/login', { method: 'post', body: body.data })
    .catch(error => { toast.add({ title: error.data.message, icon: 'i-lucide-shield-alert', color: 'error' }) })

  if(!res) return finish({ error: true })

  setUserSession({ username: state.value.username, admin: res.isAdmin, token: res.token })
  finish({ force: true })
  toast.add({ title: res.message, icon: 'i-lucide-badge-check', color: 'success' })
  await navigateTo('/')
}
</script>

<template>
  <main class="min-h-screen bg-white dark:bg-slate-900">
    <section class="mx-auto flex min-h-screen max-w-7xl flex-col items-center justify-center space-y-5 px-6">
      <NuxtImg src="/nel.png" alt="Ldap Nel" class="w-32" preload />
      <div class="w-full rounded-lg bg-slate-50 shadow sm:max-w-md md:mt-0 xl:p-0 dark:border dark:border-gray-700 dark:bg-slate-800">
        <div class="space-y-4 p-6 sm:p-8 md:space-y-6">
          <div class="flex w-full items-center justify-between">
            <h1 mclass="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white">
              Ldap Nel
            </h1>
            <div class="flex items-center justify-center space-x-2">
              <UDropdownMenu :items="themeMenu">
                <UButton icon="i-lucide-swatch-book" size="sm" color="neutral" square variant="ghost" />
              </UDropdownMenu>
            </div>
          </div>
          <UForm :schema="authSchema" :state class="w-full space-y-4 md:space-y-6" @submit="login">
            <UFormField label="Username" name="username">
              <UInput v-model="state.username" icon="i-lucide-user" size="lg" class="w-full" />
            </UFormField>
            <UFormField label="Password" name="password">
              <UInput v-model="state.password" icon="i-lucide-key-round" size="lg" type="password" class="w-full" />
            </UFormField>
            <UButton label="Login" icon="i-lucide-log-in" color="success" variant="solid" block type="submit" :loading="isLoading" />
          </UForm>
        </div>
      </div>
    </section>
  </main>
</template>
