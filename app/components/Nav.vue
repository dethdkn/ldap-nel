<script setup lang="ts">
const themeMenu = useThemeMenu()
const userMenu = useUserMenu()
const { user } = useUserSession()
const { avaiableLdaps, selectedLdap } = await useLdapConnection()
</script>

<template>
  <UHeader :toggle="false">
    <template #title>
      <div class="group flex cursor-default items-center">
        <NuxtImg src="/nel.png" alt="Ldap Nel" class="mr-3 w-8 group-hover:animate-pulse" preload />
        <span class="self-center text-2xl font-semibold whitespace-nowrap text-slate-700 transition-all duration-300 group-hover:text-green-500 dark:text-slate-200 dark:group-hover:text-green-400">Ldap Nel</span>
      </div>
    </template>

    <template #right>
      <USelect v-model="selectedLdap" :items="avaiableLdaps" size="sm" value-key="id" label-key="name" placeholder="Ldap..." class="w-24" />
      <UDropdownMenu :items="themeMenu" size="sm">
        <UButton icon="i-lucide-swatch-book" size="sm" color="neutral" square variant="ghost" />
      </UDropdownMenu>
      <UDropdownMenu :items="userMenu" size="sm">
        <div class="cursor-pointer">
          <UAvatar :alt="user.username || '?'" size="xs" />
        </div>
      </UDropdownMenu>
    </template>
  </UHeader>
</template>
