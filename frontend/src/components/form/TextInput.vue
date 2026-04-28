<script setup>
const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: '',
  },
  label: {
    type: String,
    required: true,
  },
  type: {
    type: String,
    default: 'text',
  },
  placeholder: {
    type: String,
    default: '',
  },
  error: {
    type: String,
    default: '',
  },
  required: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <label class="field">
    <span class="field__label">{{ label }}</span>
    <input
      class="field__control focus-ring"
      :value="modelValue"
      :type="type"
      :placeholder="placeholder"
      :required="required"
      @input="emit('update:modelValue', $event.target.value)"
    />
    <small v-if="error" class="field__error">{{ error }}</small>
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
  color: var(--ink);
}

.field__error {
  color: var(--danger);
  font-size: 0.78rem;
}
</style>
