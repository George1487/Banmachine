<script setup>
import AppButton from '@/components/button/AppButton.vue'

const props = defineProps({
  disabled: {
    type: Boolean,
    default: false,
  },
  loading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['run'])
</script>

<template>
  <AppButton :disabled="disabled || loading" @click="emit('run')">
    {{ loading ? 'Запуск...' : 'Запустить анализ' }}
  </AppButton>
</template>
