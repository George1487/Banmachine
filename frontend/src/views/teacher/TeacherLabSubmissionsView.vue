<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import AppButton from '@/components/button/AppButton.vue'
import SubmissionStatusBadge from '@/components/submissions/SubmissionStatusBadge.vue'
import Badge from '@/components/ui/Badge.vue'
import DataTable from '@/components/ui/DataTable.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { formatDateTime, formatScore } from '@/utils/format'
import { useLabsStore } from '@/stores/labs'
import { useSubmissionsStore } from '@/stores/submissions'

const route = useRoute()
const router = useRouter()
const labsStore = useLabsStore()
const submissionsStore = useSubmissionsStore()

const isLoading = ref(true)
const pageError = ref('')

const columns = [
  { key: 'student', label: 'Студент' },
  { key: 'submittedAt', label: 'Дата' },
  { key: 'status', label: 'Статус' },
  { key: 'topMatchScore', label: 'Макс. совпадение' },
  { key: 'finalScoreRiskLevel', label: 'Уровень риска' },
]

async function loadPage() {
  isLoading.value = true
  pageError.value = ''

  try {
    await Promise.all([
      labsStore.fetchLab(route.params.labId),
      submissionsStore.fetchLabSubmissions(route.params.labId),
    ])
  } catch (error) {
    pageError.value = error?.message || 'Не удалось загрузить работы'
  } finally {
    isLoading.value = false
  }
}

function resolveRowClass(row) {
  return row?.finalScoreRiskLevel === 'high' ? 'table__row--high-risk' : ''
}

function openMatches(row) {
  if (!row?.submissionId || !row.finalScoreRiskLevel) return
  router.push(`/teacher/submissions/${row.submissionId}/matches`)
}

onMounted(loadPage)
</script>

<template>
  <section class="view">
    <PageHeader
      :title="labsStore.currentLab?.title || 'Работы лабораторной'"
      subtitle="Список работ по лабораторной"
    >
      <template #actions>
        <AppButton variant="ghost" :as="'router-link'" :to="`/teacher/labs/${route.params.labId}/edit`">
          Редактировать
        </AppButton>
        <AppButton :as="'router-link'" :to="`/teacher/labs/${route.params.labId}/analysis`">
          Перейти к анализу
        </AppButton>
      </template>
    </PageHeader>

    <Skeleton v-if="isLoading" height="280px" />
    <ErrorState v-else-if="pageError" :message="pageError" @retry="loadPage" />

    <EmptyState
      v-else-if="submissionsStore.labSubmissions.length === 0"
      title="Пока нет работ"
      description="Студенты ещё не загрузили работы для этой лабораторной"
    />

    <DataTable
      v-else
      :columns="columns"
      :rows="submissionsStore.labSubmissions"
      :row-class="resolveRowClass"
      clickable-rows
      @row-click="openMatches"
    >
      <template #cell-student="{ row }">
        {{ row.student?.fullName || '—' }}
      </template>
      <template #cell-submittedAt="{ row }">
        {{ formatDateTime(row.submittedAt) }}
      </template>
      <template #cell-status="{ row }">
        <SubmissionStatusBadge :status="row.status" />
      </template>
      <template #cell-topMatchScore="{ row }">
        <span v-if="row.topMatchScore === null || row.topMatchScore === undefined">Нет данных</span>
        <span v-else class="mono">{{ formatScore(row.topMatchScore) }}</span>
      </template>
      <template #cell-finalScoreRiskLevel="{ row }">
        <span v-if="!row.finalScoreRiskLevel">—</span>
        <span v-else class="risk-cell">
          <span v-if="row.finalScoreRiskLevel === 'high'" class="risk-cell__icon" aria-hidden="true">!</span>
          <Badge :value="row.finalScoreRiskLevel" />
        </span>
      </template>
    </DataTable>
  </section>
</template>

<style scoped>
.view {
  display: grid;
  gap: 14px;
}

.risk-cell {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.risk-cell__icon {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  font-size: 0.72rem;
  font-weight: 800;
  color: var(--risk-high-text);
  background: var(--risk-high-bg);
  border: 1px solid #f0c2c2;
}
</style>
