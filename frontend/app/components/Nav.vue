<script setup lang="ts">
const themeMenu = useThemeMenu()
const userMenu = useUserMenu()
const { user } = useUserSession()
const { avaiableLdaps, selectedLdap } = await useLdapConnection()
</script>

<template>
  <nav class="border-b border-gray-200 bg-white dark:border-gray-800 dark:bg-slate-900">
    <div class="mx-auto flex max-w-screen-xl flex-wrap items-center justify-between p-4">
      <div class="group flex cursor-default items-center">
        <NuxtImg src="/nel.png" alt="Ldap Nel" class="mr-3 w-8 group-hover:animate-pulse" preload />
        <span class="self-center text-2xl font-semibold whitespace-nowrap text-slate-700 transition-all duration-300 group-hover:text-green-500 dark:text-slate-200 dark:group-hover:text-green-400">Ldap Nel</span>
      </div>
      <div class="flex items-center space-x-2 sm:mt-0">
        <USelect v-model="selectedLdap" :items="avaiableLdaps" size="sm" class="w-24" />
        <UDropdownMenu :items="themeMenu">
          <UButton icon="i-lucide-swatch-book" size="sm" color="neutral" square variant="ghost" />
        </UDropdownMenu>
        <UDropdownMenu :items="userMenu">
          <div class="cursor-pointer">
            <UAvatar :alt="user.username || '?'" size="xs" />
          </div>
        </UDropdownMenu>
      </div>
    </div>
  </nav>
</template>
