<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import AnalysisProgressBanner from '@/components/analysis/AnalysisProgressBanner.vue'
import AnalysisResultsTable from '@/components/analysis/AnalysisResultsTable.vue'
import AnalysisRunButton from '@/components/analysis/AnalysisRunButton.vue'
import AnalysisStatsCard from '@/components/analysis/AnalysisStatsCard.vue'
import AppButton from '@/components/button/AppButton.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { usePolling } from '@/composables/usePolling'
import { formatDateTime } from '@/utils/format'
import { useAnalysisStore } from '@/stores/analysis'
import { useLabsStore } from '@/stores/labs'
import { useNotificationsStore } from '@/stores/notifications'

const route = useRoute()
const router = useRouter()
const labsStore = useLabsStore()
const analysisStore = useAnalysisStore()
const notificationsStore = useNotificationsStore()

const isLoading = ref(true)
const runLoading = ref(false)
const pageError = ref('')
const activeJobId = ref('')
const conflictJobId = ref('')

const NON_FATAL_ANALYSIS_CODES = new Set(['analysis_job_not_found', 'submission_analysis_summary_not_found'])

const lastAnalysisJob = computed(() => analysisStore.analysisResult?.lastAnalysisJob || null)

const {
  data: polledJob,
  start: startPolling,
  isPolling,
} = usePolling(() => analysisStore.pollJob(activeJobId.value), {
  interval: 3500,
  stopWhen: (job) => ['done', 'failed'].includes(job?.status),
})

watch(polledJob, async (job) => {
  if (!job) return

  if (job.status === 'done') {
    conflictJobId.value = ''

    try {
      await analysisStore.fetchAnalysis(route.params.labId)
      notificationsStore.success('Анализ завершен')
    } catch (error) {
      if (isNonFatalAnalysisError(error)) {
        notificationsStore.warning('Результаты анализа еще формируются, обновите страницу через несколько секунд')
        return
      }

      notificationsStore.error(error?.message || 'Не удалось загрузить результаты анализа')
    }
  }

  if (job.status === 'failed') {
    conflictJobId.value = ''
    notificationsStore.error(job.errorMessage || 'Анализ завершился ошибкой')
  }
})

function extractJobIdFromError(error) {
  const payload = error?.raw?.response?.data?.data
  return payload?.analysisJobId || payload?.jobId || payload?.id || ''
}

function isNonFatalAnalysisError(error) {
  return NON_FATAL_ANALYSIS_CODES.has(error?.code)
}

async function loadPage() {
  isLoading.value = true
  pageError.value = ''

  try {
    await labsStore.fetchLab(route.params.labId)

    try {
      await analysisStore.fetchAnalysis(route.params.labId)
    } catch (error) {
      if (!isNonFatalAnalysisError(error)) {
        throw error
      }

      analysisStore.analysisResult = null
    }
  } catch (error) {
    pageError.value = error?.message || 'Не удалось загрузить данные анализа'
  } finally {
    isLoading.value = false
  }
}

async function goToJobStatus() {
  if (!conflictJobId.value) return

  activeJobId.value = conflictJobId.value
  await startPolling()
}

async function runAnalysis() {
  runLoading.value = true

  try {
    const job = await analysisStore.runAnalysis(route.params.labId)
    conflictJobId.value = ''
    activeJobId.value = job.analysisJobId || job.id
    await startPolling()
  } catch (error) {
    if (error?.code === 'analysis_job_already_running') {
      const jobId = extractJobIdFromError(error)
      conflictJobId.value = jobId
      notificationsStore.warning('Анализ уже выполняется')
      return
    }

    notificationsStore.error(error?.message || 'Не удалось запустить анализ')
  } finally {
    runLoading.value = false
  }
}

function openMatches(row) {
  if (!row?.submissionId) return
  router.push(`/teacher/submissions/${row.submissionId}/matches`)
}

onMounted(loadPage)
</script>

<template>
  <section class="view">
    <PageHeader :title="labsStore.currentLab?.title || 'Анализ лабораторной'" subtitle="Запуск и результаты анализа">
      <template #actions>
        <AppButton variant="ghost" :as="'router-link'" :to="`/teacher/labs/${route.params.labId}/submissions`">
          К работам
        </AppButton>
        <AppButton v-if="conflictJobId && !isPolling" variant="secondary" @click="goToJobStatus">
          Перейти к статусу
        </AppButton>
        <AnalysisRunButton :loading="runLoading" :disabled="isPolling" @run="runAnalysis" />
      </template>
    </PageHeader>

    <Skeleton v-if="isLoading" height="280px" />
    <ErrorState v-else-if="pageError" :message="pageError" @retry="loadPage" />

    <template v-else>
      <section v-if="lastAnalysisJob" class="last-job card-surface">
        <p>Последний запуск:</p>
        <p class="mono">{{ formatDateTime(lastAnalysisJob.createdAt) }}</p>
        <p>Статус: {{ lastAnalysisJob.status }}</p>
        <p>Завершён: {{ formatDateTime(lastAnalysisJob.finishedAt) }}</p>
      </section>

      <AnalysisProgressBanner
        v-if="analysisStore.currentJob || polledJob"
        :status="(polledJob || analysisStore.currentJob)?.status"
        :error-message="(polledJob || analysisStore.currentJob)?.errorMessage"
      />

      <EmptyState
        v-if="!analysisStore.analysisResult"
        title="Анализ пока не запускался"
        description="Нажмите «Запустить анализ», чтобы получить результаты"
      />

      <template v-else>
        <AnalysisStatsCard :stats="analysisStore.analysisResult.stats" />

        <AnalysisResultsTable
          v-if="analysisStore.analysisResult.items?.length"
          :rows="analysisStore.analysisResult.items"
          @open="openMatches"
        />

        <EmptyState
          v-else
          title="Совпадений пока нет"
          description="После запуска анализа здесь появится сводка по работам"
        />
      </template>
    </template>
  </section>
</template>

<style scoped>
.view {
  display: grid;
  gap: 14px;
}

.last-job {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
}

.last-job p {
  color: var(--ink-soft);
}
</style>
