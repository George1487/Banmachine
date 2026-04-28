<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

import MatchCard from '@/components/analysis/MatchCard.vue'
import Badge from '@/components/ui/Badge.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useAnalysisStore } from '@/stores/analysis'

const route = useRoute()
const analysisStore = useAnalysisStore()

const isLoading = ref(true)
const pageError = ref('')

const matches = computed(() => analysisStore.matches?.matches || [])

const topRisk = computed(() => matches.value[0]?.riskLevel || null)

async function loadMatches() {
  isLoading.value = true
  pageError.value = ''

  try {
    await analysisStore.fetchMatches(route.params.submissionId)
  } catch (error) {
    pageError.value = error?.message || 'Не удалось загрузить совпадения'
  } finally {
    isLoading.value = false
  }
}

onMounted(loadMatches)
</script>

<template>
  <section class="view">
    <PageHeader title="Совпадения работы" subtitle="Детальные совпадения по выбранному submission" />

    <Skeleton v-if="isLoading" height="260px" />
    <ErrorState v-else-if="pageError" :message="pageError" @retry="loadMatches" />

    <template v-else>
      <section class="summary card-surface">
        <p>Submission ID: <span class="mono">{{ analysisStore.matches?.submissionId || route.params.submissionId }}</span></p>
        <p>Найдено совпадений: <strong>{{ matches.length }}</strong></p>
        <p>
          Топ-риск:
          <Badge :value="topRisk" />
        </p>
      </section>

      <EmptyState
        v-if="matches.length === 0"
        title="Совпадений не найдено"
        description="Для этой работы нет пересечений в последнем анализе"
      />

      <div v-else class="matches">
        <MatchCard v-for="match in matches" :key="match.otherSubmissionId" :match="match" />
      </div>
    </template>
  </section>
</template>

<style scoped>
.view {
  display: grid;
  gap: 14px;
}

.summary {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
}

.summary p {
  color: var(--ink-soft);
}

.matches {
  display: grid;
  gap: 12px;
}
</style>
