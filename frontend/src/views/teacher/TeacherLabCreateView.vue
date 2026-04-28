<script setup>
import { ref } from 'vue'

import LabForm from '@/components/labs/LabForm.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import { useLabsStore } from '@/stores/labs'
import { useNotificationsStore } from '@/stores/notifications'

const labsStore = useLabsStore()
const notificationsStore = useNotificationsStore()

const loading = ref(false)

async function goToLabsList() {
  window.location.assign('/teacher/labs')
}

async function handleSubmit(payload) {
  loading.value = true

  try {
    await labsStore.createLab(payload)
    await goToLabsList()
  } catch (error) {
    notificationsStore.error(error?.message || 'Не удалось создать лабораторную')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <section class="view">
    <PageHeader title="Создать лабораторную" subtitle="Новый слот для сдачи работ" />
    <LabForm :loading="loading" @submit="handleSubmit" />
  </section>
</template>

<style scoped>
.view {
  display: grid;
  gap: 14px;
}
</style>
