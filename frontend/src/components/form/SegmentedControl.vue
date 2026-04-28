<script setup>
const props = defineProps({
  modelValue: {
    type: String,
    required: true,
  },
  options: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <div class="segmented">
    <button
      v-for="option in options"
      :key="option.value"
      class="segmented__item focus-ring"
      :class="{ 'segmented__item--active': props.modelValue === option.value }"
      type="button"
      @click="emit('update:modelValue', option.value)"
    >
      {{ option.label }}
    </button>
  </div>
</template>

<style scoped>
.segmented {
  border: 1px solid var(--line);
  border-radius: var(--radius-control);
  background: var(--surface-muted);
  padding: 4px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 4px;
}

.segmented__item {
  min-height: 40px;
  border: 0;
  border-radius: 9px;
  background: transparent;
  color: var(--ink-soft);
  font-weight: 700;
  cursor: pointer;
}

.segmented__item--active {
  background: var(--surface);
  color: var(--ink);
  box-shadow: 0 2px 10px rgba(31, 41, 55, 0.12);
}
</style>
