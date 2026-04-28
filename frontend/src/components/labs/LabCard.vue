<script setup>
import { computed } from 'vue'

import AppButton from '@/components/button/AppButton.vue'
import Card from '@/components/ui/Card.vue'
import { formatDateTime } from '@/utils/format'

import LabStatusBadge from './LabStatusBadge.vue'

const props = defineProps({
  lab: {
    type: Object,
    required: true,
  },
})

const canSubmit = computed(() => props.lab.status === 'active')
</script>

<template>
  <Card class="lab-card">
    <div class="lab-card__top">
      <h3>{{ lab.title }}</h3>
      <LabStatusBadge :status="lab.status" />
    </div>

    <p class="lab-card__description">{{ lab.description || 'Описание отсутствует' }}</p>
    <p class="lab-card__deadline">Дедлайн: {{ formatDateTime(lab.deadlineAt) }}</p>

    <AppButton
      v-if="canSubmit"
      :as="'router-link'"
      :to="`/student/labs/${lab.id}/submit`"
      variant="secondary"
    >
      Сдать отчет
    </AppButton>
  </Card>
</template>

<style scoped>
.lab-card {
  display: grid;
  gap: 12px;
}

.lab-card__top {
  display: flex;
  align-items: start;
  justify-content: space-between;
  gap: 10px;
}

.lab-card__description {
  color: var(--ink-soft);
}

.lab-card__deadline {
  color: var(--ink-faint);
  font-size: 0.9rem;
}
</style>
