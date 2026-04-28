<script setup>
import { computed } from 'vue'

import Badge from '@/components/ui/Badge.vue'

const props = defineProps({
  status: {
    type: String,
    default: '',
  },
  errorMessage: {
    type: String,
    default: '',
  },
})

const message = computed(() => {
  if (props.status === 'pending') return 'Анализ в очереди...'
  if (props.status === 'processing') return 'Анализ выполняется...'
  if (props.status === 'failed') return props.errorMessage || 'Анализ завершился ошибкой'
  if (props.status === 'done') return 'Анализ завершен'
  return 'Ожидание запуска анализа'
})
</script>

<template>
  <section class="banner card-surface">
    <Badge :value="status" />
    <p>{{ message }}</p>
  </section>
</template>

<style scoped>
.banner {
  display: flex;
  align-items: center;
  gap: 10px;
}

p {
  color: var(--ink-soft);
}
</style>
