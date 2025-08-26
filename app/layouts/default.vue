<script setup lang="ts">
const { user } = useUserSession()
</script>

<template>
  <Nav />
  <slot />
  <PasswordModal />
  <SettingsModal v-if="user.admin" />
  <AboutModal />
</template>
