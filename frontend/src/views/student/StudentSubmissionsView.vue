<script setup>
import { computed, onMounted, ref } from 'vue'

import AppButton from '@/components/button/AppButton.vue'
import SubmissionStatusBadge from '@/components/submissions/SubmissionStatusBadge.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import DataTable from '@/components/ui/DataTable.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import TextInput from '@/components/form/TextInput.vue'
import SelectInput from '@/components/form/Select.vue'
import { formatDateTime } from '@/utils/format'
import { useSubmissionsStore } from '@/stores/submissions'

const submissionsStore = useSubmissionsStore()

const isLoading = ref(true)
const errorMessage = ref('')
const search = ref('')
const statusFilter = ref('all')

const columns = [
  { key: 'labTitle', label: 'Лабораторная' },
  { key: 'submittedAt', label: 'Дата' },
  { key: 'status', label: 'Статус' },
  { key: 'actions', label: 'Действия', sortable: false },
]

const filteredRows = computed(() => {
  return submissionsStore.mySubmissions.filter((item) => {
    const byStatus = statusFilter.value === 'all' || item.status === statusFilter.value
    const byTitle = item.labTitle?.toLowerCase().includes(search.value.toLowerCase())
    return byStatus && byTitle
  })
})

async function loadSubmissions() {
  isLoading.value = true
  errorMessage.value = ''

  try {
    await submissionsStore.fetchMySubmissions()
  } catch (error) {
    errorMessage.value = error?.message || 'Не удалось загрузить работы'
  } finally {
    isLoading.value = false
  }
}

onMounted(loadSubmissions)
</script>

<template>
  <section class="view">
    <PageHeader title="Мои работы" subtitle="История отправленных submission" />

    <div class="view__filters card-surface">
      <TextInput v-model="search" label="Поиск по лабораторной" placeholder="Введите название" />
      <SelectInput
        v-model="statusFilter"
        label="Статус"
        :options="[
          { value: 'all', label: 'Все' },
          { value: 'uploaded', label: 'Загружено' },
          { value: 'parsing', label: 'Обрабатывается' },
          { value: 'parsed', label: 'Готово' },
          { value: 'failed', label: 'Ошибка' },
        ]"
      />
    </div>

    <Skeleton v-if="isLoading" height="260px" />

    <ErrorState v-else-if="errorMessage" :message="errorMessage" @retry="loadSubmissions" />

    <EmptyState
      v-else-if="filteredRows.length === 0"
      title="Сдач пока нет"
      description="Загрузите первую лабораторную работу"
    />

    <DataTable v-else :columns="columns" :rows="filteredRows">
      <template #cell-submittedAt="{ row }">
        {{ formatDateTime(row.submittedAt) }}
      </template>
      <template #cell-status="{ row }">
        <SubmissionStatusBadge :status="row.status" />
      </template>
      <template #cell-actions="{ row }">
        <AppButton size="sm" variant="ghost" :as="'router-link'" :to="`/student/labs/${row.labId}/submit`">
          Открыть
        </AppButton>
      </template>
    </DataTable>
  </section>
</template>

<style scoped>
.view {
  display: grid;
  gap: 14px;
}

.view__filters {
  padding: 14px;
  display: grid;
  grid-template-columns: 1fr 220px;
  gap: 10px;
}

@media (max-width: 900px) {
  .view__filters {
    grid-template-columns: 1fr;
  }
}
</style>
