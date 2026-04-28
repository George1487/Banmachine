<script setup>
import { computed } from 'vue'

const props = defineProps({
  label: {
    type: String,
    required: true,
  },
  value: {
    type: Number,
    default: 0,
  },
})

const clamped = computed(() => Math.max(0, Math.min(1, Number(props.value) || 0)))

const color = computed(() => {
  if (clamped.value >= 0.8) return 'var(--score-5)'
  if (clamped.value >= 0.6) return 'var(--score-4)'
  if (clamped.value >= 0.4) return 'var(--score-3)'
  if (clamped.value >= 0.2) return 'var(--score-2)'
  return 'var(--score-1)'
})
</script>

<template>
  <div class="score">
    <div class="score__header">
      <span>{{ label }}</span>
      <strong class="mono">{{ (clamped * 100).toFixed(0) }}%</strong>
    </div>
    <div class="score__track">
      <div class="score__fill" :style="{ width: `${clamped * 100}%`, background: color }" />
    </div>
  </div>
</template>

<style scoped>
.score {
  display: grid;
  gap: 6px;
}

.score__header {
  display: flex;
  justify-content: space-between;
  color: var(--ink-soft);
  font-size: 0.85rem;
}

.score__track {
  height: 8px;
  border-radius: var(--radius-pill);
  background: var(--surface-muted);
  border: 1px solid var(--line);
  overflow: hidden;
}

.score__fill {
  height: 100%;
  transition: width 300ms ease-out;
}
</style>
