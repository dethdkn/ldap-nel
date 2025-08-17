<script setup lang='ts'>
import type { TabsItem } from '@nuxt/ui'

const settingsModal = useSettingsModal()

const items: TabsItem[] = [
  { label: 'Ldaps', icon: 'i-lucide-folder', slot: 'ldaps' },
  { label: 'Users', icon: 'i-lucide-user', slot: 'users' },
]
</script>

<template>
  <UModal v-model:open="settingsModal" title="Settings">
    <template #body>
      <UTabs :items class="w-full">
        <template #ldaps>
          <LdapsTab />
        </template>
        <template #users>
          <UsersTab />
        </template>
      </UTabs>
    </template>
  </UModal>
</template>
