<script setup lang='ts'>
const { avaiableLdaps } = await useLdapConnection()
const settingsModal = useSettingsModal()
</script>

<template>
  <div>
    <div v-if="(avaiableLdaps?.length || 0) > 0">
      {{ avaiableLdaps }}
    </div>
    <div v-else class="min-h-content flex flex-col items-center justify-center space-y-1">
      <NuxtImg src="/nel.png" alt="Ldap Nel" class="w-32" preload />
      <h1 class="my-4 text-xl font-bold">
        Welcome to Ldap Nel!
      </h1>
      <p>It looks like you don’t have any connections yet.</p>
      <p>Go to settings to create your first connection.</p>
      <UButton label="Open Settings" icon="i-lucide-settings" color="neutral" variant="outline" class="mt-4" @click="settingsModal = true" />
    </div>
  </div>
</template>
