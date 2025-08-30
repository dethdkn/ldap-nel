<script setup lang="ts">
const error = useError()
</script>

<template>
  <UApp>
    <NuxtLayout>
      <UError :error :clear="{ color: 'neutral', size: 'xl', icon: 'i-lucide-arrow-left', class: 'rounded-full'}">
        <template #message>
          <p>
            If you believe this is a mistake, please open an issue on
            <NuxtLink to="https://github.com/dethdkn/ldap-nel/issues" external target="_blank" class="cursor-pointer hover:underline">
              GitHub
            </NuxtLink>
          </p>
        </template>
      </UError>
    </NuxtLayout>
  </UApp>
</template>
