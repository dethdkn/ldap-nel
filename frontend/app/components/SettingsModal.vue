<script setup lang='ts'>
import type { TabsItem } from '@nuxt/ui'

const settingsModal = useSettingsModal()

const items: TabsItem[] = [
  { label: 'Ldaps', icon: 'i-lucide-lock', slot: 'ldaps' },
  { label: 'Users', icon: 'i-lucide-user', slot: 'users' },
]
</script>

<template>
  <UModal v-model:open="settingsModal" title="Settings">
    <template #body>
      <UTabs :items class="w-full">
        <template #ldaps>
          <div class="p-4">
            <h2 class="mb-4 text-lg font-semibold">
              Ldaps Settings
            </h2>
          </div>
        </template>
        <template #users>
          <UsersTab />
        </template>
      </UTabs>
    </template>
  </UModal>
</template>
