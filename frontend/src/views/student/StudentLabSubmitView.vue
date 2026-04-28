<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import SubmissionUpload from '@/components/submissions/SubmissionUpload.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useFileUpload } from '@/composables/useFileUpload'
import { usePolling } from '@/composables/usePolling'
import { useLabsStore } from '@/stores/labs'
import { useSubmissionsStore } from '@/stores/submissions'

const route = useRoute()
const labsStore = useLabsStore()
const submissionsStore = useSubmissionsStore()

const isLoading = ref(true)
const pageError = ref('')
const uploadError = ref('')
const status = ref('')
const submissionId = ref('')

const labId = computed(() => route.params.labId)
const lab = computed(() => labsStore.currentLab)
const isClosed = computed(() => lab.value?.status === 'closed')

const { upload, isUploading, progress } = useFileUpload((targetLabId, file, onProgress) =>
  submissionsStore.uploadSubmission(targetLabId, file, onProgress),
)

const {
  data: pollingSubmission,
  start: startPolling,
  isPolling,
} = usePolling(() => submissionsStore.fetchSubmission(submissionId.value), {
  interval: 3500,
  stopWhen: (result) => ['parsed', 'failed'].includes(result?.status),
})

watch(pollingSubmission, (value) => {
  if (!value) return

  status.value = value.status

  if (value.status === 'failed') {
    uploadError.value = value.errorMessage || 'Не удалось обработать файл. Попробуйте загрузить заново.'
  }
})

async function loadLab() {
  isLoading.value = true
  pageError.value = ''

  try {
    await labsStore.fetchLab(labId.value)
  } catch (error) {
    pageError.value = error?.message || 'Не удалось загрузить лабораторную'
  } finally {
    isLoading.value = false
  }
}

async function handleFileSelected(file) {
  uploadError.value = ''

  try {
    const result = await upload(labId.value, file)
    submissionId.value = result.id
    status.value = result.status || 'uploaded'

    await startPolling()
  } catch (error) {
    if (error?.code === 'lab_not_active') {
      uploadError.value = 'Лабораторная закрыта для сдачи'
      return
    }

    if (error?.code === 'unsupported_file_type') {
      uploadError.value = 'Формат файла не поддерживается. Попробуйте docx'
      return
    }

    uploadError.value = error?.message || 'Не удалось загрузить файл'
  }
}

function handleRetry() {
  uploadError.value = ''
  status.value = ''
}

onMounted(loadLab)
</script>

<template>
  <section class="view">
    <PageHeader
      title="Сдача лабораторной"
      :subtitle="lab ? `${lab.title} • дедлайн ${new Date(lab.deadlineAt).toLocaleString('ru-RU')}` : ''"
    />

    <Skeleton v-if="isLoading" height="300px" />

    <ErrorState v-else-if="pageError" :message="pageError" @retry="loadLab" />

    <template v-else-if="lab">
      <EmptyState
        v-if="isClosed"
        title="Прием работ закрыт"
        description="Эта лабораторная закрыта преподавателем, загрузка недоступна"
      />

      <SubmissionUpload
        v-else
        :disabled="isClosed"
        :is-uploading="isUploading"
        :progress="progress"
        :status="status"
        :error-message="uploadError"
        @file-selected="handleFileSelected"
        @retry="handleRetry"
      />

      <p v-if="isPolling" class="view__polling">Проверяем статус обработки...</p>
    </template>
  </section>
</template>

<style scoped>
.view {
  display: grid;
  gap: 14px;
}

.view__polling {
  color: var(--ink-soft);
  font-size: 0.9rem;
}
</style>
