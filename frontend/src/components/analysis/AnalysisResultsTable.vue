<script setup>
import DataTable from '@/components/ui/DataTable.vue'
import Badge from '@/components/ui/Badge.vue'
import { formatScore } from '@/utils/format'

const props = defineProps({
  rows: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['open'])

const columns = [
  { key: 'student', label: 'Студент' },
  { key: 'topMatchScore', label: 'Макс. совпадение' },
  { key: 'finalScoreRiskLevel', label: 'Уровень риска' },
]
</script>

<template>
  <DataTable :columns="columns" :rows="rows" row-key="submissionId" clickable-rows @row-click="emit('open', $event)">
    <template #cell-student="{ row }">
      {{ row.student?.fullName || '—' }}
    </template>
    <template #cell-topMatchScore="{ row }">
      <span class="mono">{{ formatScore(row.topMatchScore) }}</span>
    </template>
    <template #cell-finalScoreRiskLevel="{ row }">
      <Badge :value="row.finalScoreRiskLevel" />
    </template>
  </DataTable>
</template>
