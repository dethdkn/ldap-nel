<script setup lang='ts'>
const { user } = useUserSession()
const { avaiableLdaps, selectedLdap } = await useLdapConnection()
const settingsModal = useSettingsModal()
</script>

<template>
  <div v-if="selectedLdap" class="mx-auto mt-4 flex max-w-7xl flex-col items-center justify-center px-2 md:px-4">
    <Ldap />
  </div>
  <template v-else-if="(avaiableLdaps?.length || 0) > 0">
    <NuxtImg src="/nel.png" alt="Ldap Nel" class="w-32" />
    <h1 class="my-4 text-xl font-bold">
      Welcome to Ldap Nel!
    </h1>
    <p>Please choose a connection from the navigation bar above.</p>
  </template>
  <template v-else>
    <NuxtImg src="/nel.png" alt="Ldap Nel" class="w-32" />
    <h1 class="my-4 text-xl font-bold">
      Welcome to Ldap Nel!
    </h1>
    <p>It looks like you don't have any connections yet.</p>
    <template v-if="user.admin">
      <p>Go to settings to create your first connection.</p>
      <UButton label="Open Settings" icon="i-lucide-settings" color="neutral" variant="outline" class="mt-4" @click="settingsModal = true" />
    </template>
    <p v-else>
      Please contact an administrator to create one.
    </p>
  </template>
</template>
