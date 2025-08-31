<script setup lang='ts'>
const aboutModal = useAboutModal()
const { version } = useRuntimeConfig().public
</script>

<template>
  <UModal v-model:open="aboutModal" title="About">
    <template #body>
      <div class="flex flex-col items-center justify-center space-y-4">
        <NuxtImg src="/nel.png" alt="Ldap Nel" class="w-32" />
        <h1 class="text-xl font-bold">
          Ldap Nel V{{ version }}
        </h1>
        <div class="flex items-center justify-center space-x-2">
          <UButton icon="i-lucide-globe" size="md" color="neutral" variant="outline" to="https://ldap-nel.vercel.app" external target="_blank" />
          <UButton icon="i-lucide-github" size="md" color="neutral" variant="outline" to="https://github.com/dethdkn/ldap-nel" external target="_blank" />
          <UButton icon="i-lucide-history" size="md" color="neutral" variant="outline" to="https://github.com/dethdkn/ldap-nel/releases" external target="_blank" />
        </div>
        <div class="text-center">
          <p>
            Copyright © 2025 <NuxtLink class="cursor-pointer text-blue-500 hover:underline" to="https://rosa.dev.br" external target="_blank">
              Gabriel 'DethDKN' Rosa
            </NuxtLink>
          </p>
          <p>
            <!-- eslint-disable-next-line harlanzw/link-lowercase -->
            This project is under <NuxtLink class="cursor-pointer text-blue-500 hover:underline" to="https://github.com/dethdkn/ldap-nel/blob/main/LICENSE" external target="_blank">
              MIT license
            </NuxtLink>
          </p>
        </div>
      </div>
    </template>
  </UModal>
</template>
