<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: {
    type: String,
    default: '',
  },
})

const steps = computed(() => {
  const state = props.status
  const map = {
    uploaded: ['done', 'current', 'idle'],
    parsing: ['done', 'done', 'current'],
    parsed: ['done', 'done', 'done'],
    failed: ['done', 'done', 'failed'],
  }

  const visual = map[state] || ['idle', 'idle', 'idle']

  return [
    { label: 'Загружено', state: visual[0] },
    { label: 'Обработка', state: visual[1] },
    { label: state === 'failed' ? 'Ошибка' : 'Готово', state: visual[2] },
  ]
})
</script>

<template>
  <ol class="timeline">
    <li v-for="step in steps" :key="step.label" class="timeline__item" :class="`timeline__item--${step.state}`">
      <span class="timeline__dot" />
      <span class="timeline__label">{{ step.label }}</span>
    </li>
  </ol>
</template>

<style scoped>
.timeline {
  display: grid;
  gap: 10px;
  padding: 0;
  margin: 0;
  list-style: none;
}

.timeline__item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--ink-faint);
}

.timeline__dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: var(--surface-muted);
  border: 1px solid var(--line);
}

.timeline__item--done {
  color: var(--success);
}

.timeline__item--done .timeline__dot {
  background: var(--success);
  border-color: var(--success);
}

.timeline__item--current {
  color: var(--warning);
}

.timeline__item--current .timeline__dot {
  background: var(--warning);
  border-color: var(--warning);
}

.timeline__item--failed {
  color: var(--danger);
}

.timeline__item--failed .timeline__dot {
  background: var(--danger);
  border-color: var(--danger);
}
</style>
