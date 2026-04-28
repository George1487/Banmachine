<script setup>
import { computed } from 'vue'

const props = defineProps({
  value: {
    type: String,
    default: '',
  },
  tone: {
    type: String,
    default: '',
  },
})

const toneByValue = {
  teacher: 'warning',
  student: 'success',
  active: 'success',
  closed: 'default',
  uploaded: 'default',
  parsing: 'warning',
  parsed: 'success',
  failed: 'danger',
  pending: 'default',
  processing: 'warning',
  done: 'success',
  low: 'risk-low',
  medium: 'risk-medium',
  high: 'risk-high',
}

const normalizedTone = computed(() => props.tone || toneByValue[props.value] || 'default')

const label = computed(() => {
  const map = {
    teacher: 'Преподаватель',
    student: 'Студент',
    active: 'Активна',
    closed: 'Закрыта',
    uploaded: 'Загружено',
    parsing: 'Обрабатывается',
    parsed: 'Готово',
    failed: 'Ошибка',
    pending: 'В очереди',
    processing: 'Выполняется',
    done: 'Завершено',
    low: 'Низкий',
    medium: 'Средний',
    high: 'Высокий',
  }

  return map[props.value] || props.value || '—'
})
</script>

<template>
  <span class="badge" :class="`badge--${normalizedTone}`">{{ label }}</span>
</template>

<style scoped>
.badge {
  display: inline-flex;
  align-items: center;
  min-height: 24px;
  border-radius: var(--radius-pill);
  padding: 0 10px;
  font-size: 0.76rem;
  font-weight: 700;
  border: 1px solid transparent;
}

.badge--default {
  color: var(--ink-soft);
  border-color: var(--line);
  background: var(--surface-muted);
}

.badge--success {
  color: #1e6b2b;
  border-color: #bde3c2;
  background: var(--risk-low-bg);
}

.badge--warning {
  color: #9a5a00;
  border-color: #f0d3a5;
  background: var(--risk-medium-bg);
}

.badge--danger {
  color: #a61b1b;
  border-color: #f0c2c2;
  background: var(--risk-high-bg);
}

.badge--risk-low {
  color: var(--risk-low-text);
  border-color: #bde3c2;
  background: var(--risk-low-bg);
}

.badge--risk-medium {
  color: var(--risk-medium-text);
  border-color: #f0d3a5;
  background: var(--risk-medium-bg);
}

.badge--risk-high {
  color: var(--risk-high-text);
  border-color: #f0c2c2;
  background: var(--risk-high-bg);
}
</style>
