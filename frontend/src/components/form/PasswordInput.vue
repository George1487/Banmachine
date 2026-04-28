<script setup>
import { ref } from 'vue'

import TextInput from './TextInput.vue'

defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  label: {
    type: String,
    default: 'Пароль',
  },
  placeholder: {
    type: String,
    default: 'Введите пароль',
  },
  error: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['update:modelValue'])

const show = ref(false)
</script>

<template>
  <div class="password-field">
    <TextInput
      :model-value="modelValue"
      :label="label"
      :placeholder="placeholder"
      :error="error"
      :type="show ? 'text' : 'password'"
      @update:model-value="emit('update:modelValue', $event)"
    />
    <button class="password-field__toggle focus-ring" type="button" @click="show = !show">
      {{ show ? 'Скрыть' : 'Показать' }}
    </button>
  </div>
</template>

<style scoped>
.password-field {
  position: relative;
}

.password-field__toggle {
  position: absolute;
  right: 12px;
  top: 30px;
  min-height: 28px;
  border: 0;
  background: transparent;
  color: var(--ink-soft);
  font-weight: 700;
  cursor: pointer;
}
</style>
