<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  columns: {
    type: Array,
    default: () => [],
  },
  rows: {
    type: Array,
    default: () => [],
  },
  rowKey: {
    type: String,
    default: 'id',
  },
  clickableRows: {
    type: Boolean,
    default: false,
  },
  rowClass: {
    type: Function,
    default: null,
  },
})

const emit = defineEmits(['row-click'])

const sortKey = ref('')
const sortDirection = ref('asc')

const sortedRows = computed(() => {
  if (!sortKey.value) return props.rows

  const column = props.columns.find((item) => item.key === sortKey.value)
  if (!column) return props.rows

  const sorted = [...props.rows].sort((left, right) => {
    const leftValue = typeof column.sortValue === 'function' ? column.sortValue(left) : left[sortKey.value]
    const rightValue = typeof column.sortValue === 'function' ? column.sortValue(right) : right[sortKey.value]

    if (leftValue === null || leftValue === undefined) return 1
    if (rightValue === null || rightValue === undefined) return -1

    if (typeof leftValue === 'number' && typeof rightValue === 'number') {
      return leftValue - rightValue
    }

    return String(leftValue).localeCompare(String(rightValue), 'ru')
  })

  return sortDirection.value === 'asc' ? sorted : sorted.reverse()
})

function toggleSort(column) {
  if (column.sortable === false) return

  if (sortKey.value === column.key) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
    return
  }

  sortKey.value = column.key
  sortDirection.value = 'asc'
}
</script>

<template>
  <div class="table-wrap card-surface">
    <table class="table">
      <thead>
        <tr>
          <th v-for="column in columns" :key="column.key">
            <button
              class="table__sort-btn"
              :class="{ 'table__sort-btn--active': sortKey === column.key }"
              :disabled="column.sortable === false"
              type="button"
              @click="toggleSort(column)"
            >
              {{ column.label }}
            </button>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="row in sortedRows"
          :key="row[rowKey] || row.id"
          :class="[clickableRows && 'table__row--clickable', rowClass ? rowClass(row) : '']"
          @click="clickableRows && emit('row-click', row)"
        >
          <td v-for="column in columns" :key="column.key" :class="column.className">
            <slot :name="`cell-${column.key}`" :row="row">
              {{ row[column.key] ?? '—' }}
            </slot>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.table-wrap {
  overflow: auto;
  padding: 8px;
}

.table {
  width: 100%;
  border-collapse: collapse;
  min-width: 700px;
}

th,
td {
  text-align: left;
  padding: 12px;
  border-bottom: 1px solid var(--line);
  vertical-align: middle;
}

th {
  color: var(--ink-soft);
  font-size: 0.82rem;
  font-weight: 700;
}

.table__sort-btn {
  border: 0;
  background: transparent;
  color: inherit;
  font-size: inherit;
  font-weight: inherit;
  padding: 0;
  cursor: pointer;
}

.table__sort-btn:disabled {
  cursor: default;
}

.table__sort-btn--active {
  color: var(--accent);
}

tbody tr:last-child td {
  border-bottom: 0;
}

.table__row--clickable {
  cursor: pointer;
}

.table__row--clickable:hover {
  background: #ffffff90;
}

:deep(.table__row--high-risk) td {
  background: color-mix(in srgb, var(--risk-high-bg), white 40%);
}

:deep(.table__row--high-risk.table__row--clickable:hover) td {
  background: color-mix(in srgb, var(--risk-high-bg), white 28%);
}
</style>
