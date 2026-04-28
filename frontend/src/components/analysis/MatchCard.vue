<script setup>
import Badge from '@/components/ui/Badge.vue'
import Card from '@/components/ui/Card.vue'
import ScoreBar from '@/components/ui/ScoreBar.vue'
import { formatScore } from '@/utils/format'

const props = defineProps({
  match: {
    type: Object,
    required: true,
  },
})
</script>

<template>
  <Card class="match-card">
    <div class="match-card__top">
      <div>
        <h3>{{ props.match.student?.fullName || props.match.otherSubmissionId }}</h3>
        <p class="match-card__meta mono">{{ props.match.otherSubmissionId }}</p>
      </div>
      <Badge :value="props.match.riskLevel" />
    </div>

    <div class="match-card__scores">
      <ScoreBar label="Текст" :value="props.match.textScore" />
      <ScoreBar label="Расчеты" :value="props.match.calculationScore" />
      <ScoreBar label="Изображения" :value="props.match.imagesScore" />
    </div>

    <p class="match-card__final mono">Итоговый балл: {{ formatScore(props.match.finalScore) }}</p>
  </Card>
</template>

<style scoped>
.match-card {
  display: grid;
  gap: 12px;
}

.match-card__top {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.match-card__meta {
  color: var(--ink-faint);
  font-size: 0.8rem;
}

.match-card__scores {
  display: grid;
  gap: 10px;
}

.match-card__final {
  font-weight: 700;
}
</style>
