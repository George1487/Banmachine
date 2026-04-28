<script setup>
import { onBeforeUnmount, onMounted } from 'vue'

import AppShell from '@/components/app/AppShell.vue'
import ToastContainer from '@/components/app/ToastContainer.vue'
import { useNotificationsStore } from '@/stores/notifications'

const notificationsStore = useNotificationsStore()

function onApiError(event) {
  const { status } = event.detail || {}

  if (status === 403) {
    notificationsStore.warning('Недостаточно прав')
    return
  }

  if (status === 404) {
    notificationsStore.warning('Сущность не найдена')
    return
  }

  if (status >= 500) {
    notificationsStore.error('Ошибка сервера, попробуйте позже')
  }
}

onMounted(() => {
  window.addEventListener('banmachine:api-error', onApiError)
})

onBeforeUnmount(() => {
  window.removeEventListener('banmachine:api-error', onApiError)
})
</script>

<template>
  <AppShell>
    <RouterView />
  </AppShell>
  <ToastContainer />
</template>
