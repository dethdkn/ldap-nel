<script setup lang="ts">
const { user } = useUserSession()
</script>

<template>
  <Nav />
  <UMain class="flex flex-col items-center justify-center space-y-1">
    <slot />
    <PasswordModal />
    <SettingsModal v-if="user.admin" />
    <AboutModal />
  </UMain>
</template>
