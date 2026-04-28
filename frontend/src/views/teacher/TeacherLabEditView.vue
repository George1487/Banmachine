<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import LabForm from '@/components/labs/LabForm.vue'
import ConfirmDialog from '@/components/ui/ConfirmDialog.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useLabsStore } from '@/stores/labs'
import { useNotificationsStore } from '@/stores/notifications'

const route = useRoute()
const router = useRouter()
const labsStore = useLabsStore()
const notificationsStore = useNotificationsStore()

const isLoading = ref(true)
const isSaving = ref(false)
const pageError = ref('')
const showCloseConfirm = ref(false)
const pendingPayload = ref(null)

async function goToLabsList() {
  window.location.assign('/teacher/labs')
}

async function loadLab() {
  isLoading.value = true
  pageError.value = ''

  try {
    await labsStore.fetchLab(route.params.labId)
  } catch (error) {
    pageError.value = error?.message || 'Не удалось загрузить лабораторную'
  } finally {
    isLoading.value = false
  }
}

async function save(payload) {
  isSaving.value = true

  try {
    await labsStore.updateLab(route.params.labId, payload)
    await goToLabsList()
  } catch (error) {
    notificationsStore.error(error?.message || 'Не удалось обновить лабораторную')
  } finally {
    isSaving.value = false
  }
}

function handleSubmit(payload) {
  if (payload.status === 'Closed' && labsStore.currentLab?.status !== 'closed') {
    pendingPayload.value = payload
    showCloseConfirm.value = true
    return
  }

  save(payload)
}

function confirmClose() {
  if (!pendingPayload.value) return
  save(pendingPayload.value)
  pendingPayload.value = null
}

onMounted(loadLab)
</script>

<template>
  <section class="view">
    <PageHeader title="Редактирование лабораторной" subtitle="Обновите параметры лабораторной" />

    <Skeleton v-if="isLoading" height="280px" />
    <ErrorState v-else-if="pageError" :message="pageError" @retry="loadLab" />
    <LabForm v-else mode="edit" :model="labsStore.currentLab" :loading="isSaving" @submit="handleSubmit" />

    <ConfirmDialog
      v-model="showCloseConfirm"
      title="Закрыть лабораторную?"
      message="Студенты больше не смогут сдавать работы. Продолжить?"
      @confirm="confirmClose"
    />
  </section>
</template>

<style scoped>
.view {
  display: grid;
  gap: 14px;
}
</style>
