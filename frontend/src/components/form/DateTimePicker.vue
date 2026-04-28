<script setup>
defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  label: {
    type: String,
    required: true,
  },
})

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <label class="field">
    <span class="field__label">{{ label }}</span>
    <input
      class="field__control focus-ring"
      type="datetime-local"
      :value="modelValue"
      @input="emit('update:modelValue', $event.target.value)"
    />
  </label>
</template>

<style scoped>
.field {
  display: grid;
  gap: 6px;
}

.field__label {
  color: var(--ink-soft);
  font-size: 0.85rem;
  font-weight: 700;
}

.field__control {
  min-height: 44px;
  border-radius: var(--radius-control);
  border: 1px solid var(--line);
  background: var(--surface);
  padding: 0 12px;
}
</style>
