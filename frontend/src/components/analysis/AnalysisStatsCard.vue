<script setup>
import StatTile from '@/components/ui/StatTile.vue'
import { formatScore } from '@/utils/format'

const props = defineProps({
  stats: {
    type: Object,
    required: true,
  },
})
</script>

<template>
  <section class="stats card-surface">
    <div class="stats__grid">
      <StatTile label="Всего работ" :value="props.stats.totalSubmissions || 0" />
      <StatTile label="Актуальные" :value="props.stats.actualSubmissions || 0" />
      <StatTile label="Обработано" :value="props.stats.parsedSubmissions || 0" />
      <StatTile label="Пиковый риск" :value="formatScore(props.stats.maxFinalScore)" />
    </div>

    <div class="stats__risk">
      <p class="stats__risk-item stats__risk-item--high">Высокий: {{ props.stats.highRiskCount || 0 }}</p>
      <p class="stats__risk-item stats__risk-item--medium">Средний: {{ props.stats.mediumRiskCount || 0 }}</p>
      <p class="stats__risk-item stats__risk-item--low">Низкий: {{ props.stats.lowRiskCount || 0 }}</p>
    </div>
  </section>
</template>

<style scoped>
.stats {
  display: grid;
  gap: 12px;
}

.stats__grid {
  display: grid;
  gap: 10px;
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.stats__risk {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.stats__risk-item {
  border-radius: var(--radius-pill);
  padding: 6px 10px;
  font-size: 0.85rem;
  font-weight: 700;
}

.stats__risk-item--high {
  background: var(--risk-high-bg);
  color: var(--risk-high-text);
}

.stats__risk-item--medium {
  background: var(--risk-medium-bg);
  color: var(--risk-medium-text);
}

.stats__risk-item--low {
  background: var(--risk-low-bg);
  color: var(--risk-low-text);
}

@media (max-width: 900px) {
  .stats__grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
