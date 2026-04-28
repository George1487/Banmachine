<script setup>
import { computed, onMounted, ref } from 'vue'

import LabCard from '@/components/labs/LabCard.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useLabsStore } from '@/stores/labs'

const labsStore = useLabsStore()

const isLoading = ref(true)
const errorMessage = ref('')

const labs = computed(() => {
  return [...labsStore.labs].sort((a, b) => new Date(a.deadlineAt) - new Date(b.deadlineAt))
})

async function loadLabs() {
  isLoading.value = true
  errorMessage.value = ''

  try {
    await labsStore.fetchLabs()
  } catch (error) {
    errorMessage.value = error?.message || 'Не удалось загрузить лабораторные'
  } finally {
    isLoading.value = false
  }
}

onMounted(loadLabs)
</script>

<template>
  <section class="view">
    <PageHeader title="Лабораторные" subtitle="Выберите лабораторную и загрузите отчет" />

    <div v-if="isLoading" class="view__grid">
      <Skeleton v-for="index in 3" :key="index" height="180px" />
    </div>

    <ErrorState v-else-if="errorMessage" :message="errorMessage" @retry="loadLabs" />

    <EmptyState
      v-else-if="labs.length === 0"
      title="Нет доступных лабораторных"
      description="Как только преподаватель создаст лабораторную, она появится здесь"
    />

    <div v-else class="view__grid">
      <LabCard v-for="lab in labs" :key="lab.id" :lab="lab" />
    </div>
  </section>
</template>

<style scoped>
.view {
  display: grid;
  gap: 14px;
}

.view__grid {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  align-items: stretch;
}
</style>
