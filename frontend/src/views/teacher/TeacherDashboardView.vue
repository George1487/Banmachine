<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import AppButton from '@/components/button/AppButton.vue'
import LabStatusBadge from '@/components/labs/LabStatusBadge.vue'
import DataTable from '@/components/ui/DataTable.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import StatTile from '@/components/ui/StatTile.vue'
import { useLabsStore } from '@/stores/labs'

const labsStore = useLabsStore()
const router = useRouter()

const isLoading = ref(true)
const errorMessage = ref('')

const columns = [
  { key: 'title', label: 'Лабораторная' },
  { key: 'status', label: 'Статус' },
  { key: 'submissionsCount', label: 'Работ сдано' },
  { key: 'parsedSubmissionsCount', label: 'Обработано' },
  { key: 'actions', label: 'Действия', sortable: false },
]

const totals = computed(() => {
  const labs = labsStore.teacherLabs
  const submissionsCount = labs.reduce((sum, item) => sum + (item.submissionsCount || 0), 0)
  const parsedCount = labs.reduce((sum, item) => sum + (item.parsedSubmissionsCount || 0), 0)
  const mostActive = [...labs].sort((a, b) => (b.submissionsCount || 0) - (a.submissionsCount || 0))[0]

  return {
    labs: labs.length,
    submissionsCount,
    parsedCount,
    mostActive: mostActive?.title || '—',
  }
})

async function loadLabs() {
  isLoading.value = true
  errorMessage.value = ''

  try {
    await labsStore.fetchTeacherLabs()
  } catch (error) {
    errorMessage.value = error?.message || 'Не удалось загрузить лабораторные'
  } finally {
    isLoading.value = false
  }
}

function openSubmissions(labId) {
  router.push(`/teacher/labs/${labId}/submissions`)
}

onMounted(loadLabs)
</script>

<template>
  <section class="view">
    <PageHeader title="Лабораторные преподавателя" subtitle="Управление лабораторными и анализом">
      <template #actions>
        <AppButton :as="'router-link'" to="/teacher/labs/create">Создать лабораторную</AppButton>
      </template>
    </PageHeader>

    <div class="view__stats">
      <StatTile label="Всего лабораторных" :value="totals.labs" />
      <StatTile label="Всего работ" :value="totals.submissionsCount" />
      <StatTile label="Обработано" :value="totals.parsedCount" />
      <StatTile label="Самая активная" :value="totals.mostActive" />
    </div>

    <Skeleton v-if="isLoading" height="280px" />

    <ErrorState v-else-if="errorMessage" :message="errorMessage" @retry="loadLabs" />

    <EmptyState
      v-else-if="labsStore.teacherLabs.length === 0"
      title="Нет лабораторных"
      description="Создайте первую лабораторную для приема работ"
    >
      <template #action>
        <AppButton :as="'router-link'" to="/teacher/labs/create">Создать</AppButton>
      </template>
    </EmptyState>

    <DataTable
      v-else
      :columns="columns"
      :rows="labsStore.teacherLabs"
      row-key="id"
      clickable-rows
      @row-click="openSubmissions($event.id)"
    >
      <template #cell-status="{ row }">
        <LabStatusBadge :status="row.status" />
      </template>
      <template #cell-actions="{ row }">
        <div class="actions">
          <AppButton size="sm" variant="secondary" @click.stop="openSubmissions(row.id)">Работы</AppButton>
          <AppButton size="sm" variant="ghost" :as="'router-link'" :to="`/teacher/labs/${row.id}/analysis`" @click.stop>
            Анализ
          </AppButton>
          <AppButton size="sm" variant="ghost" :as="'router-link'" :to="`/teacher/labs/${row.id}/edit`" @click.stop>
            Редактировать
          </AppButton>
        </div>
      </template>
    </DataTable>
  </section>
</template>

<style scoped>
.view {
  display: grid;
  gap: 14px;
}

.view__stats {
  display: grid;
  gap: 10px;
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.actions {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

@media (max-width: 900px) {
  .view__stats {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
