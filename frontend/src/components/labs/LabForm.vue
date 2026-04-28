<script setup>
import { reactive, watch } from 'vue'

import AppButton from '@/components/button/AppButton.vue'
import DateTimePicker from '@/components/form/DateTimePicker.vue'
import SelectInput from '@/components/form/Select.vue'
import TextInput from '@/components/form/TextInput.vue'

const props = defineProps({
  model: {
    type: Object,
    default: () => ({
      title: '',
      description: '',
      deadlineAt: '',
      status: 'active',
    }),
  },
  mode: {
    type: String,
    default: 'create',
  },
  loading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['submit'])

const form = reactive({
  title: '',
  description: '',
  deadlineAt: '',
  status: 'active',
})

function toApiLabStatus(status) {
  return status === 'closed' ? 'Closed' : 'Active'
}

function toDateTimeLocalValue(value) {
  if (!value) return ''

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''

  const pad = (num) => String(num).padStart(2, '0')
  const year = date.getFullYear()
  const month = pad(date.getMonth() + 1)
  const day = pad(date.getDate())
  const hours = pad(date.getHours())
  const minutes = pad(date.getMinutes())

  return `${year}-${month}-${day}T${hours}:${minutes}`
}

watch(
  () => props.model,
  (value) => {
    form.title = value?.title || ''
    form.description = value?.description || ''
    form.deadlineAt = toDateTimeLocalValue(value?.deadlineAt)
    form.status = value?.status || 'active'
  },
  { immediate: true, deep: true },
)

function handleSubmit() {
  const payload = {
    title: form.title,
    description: form.description,
    deadlineAt: form.deadlineAt ? new Date(form.deadlineAt).toISOString() : null,
  }

  if (props.mode === 'edit') {
    payload.status = toApiLabStatus(form.status)
  }

  emit('submit', payload)
}
</script>

<template>
  <form class="lab-form card-surface" @submit.prevent="handleSubmit">
    <TextInput v-model="form.title" label="Название" placeholder="Лабораторная работа №1" required />

    <label class="lab-form__textarea-field">
      <span>Описание</span>
      <textarea v-model="form.description" class="lab-form__textarea focus-ring" rows="5" />
    </label>

    <DateTimePicker v-model="form.deadlineAt" label="Дедлайн" />

    <SelectInput
      v-if="mode === 'edit'"
      v-model="form.status"
      label="Статус"
      :options="[
        { value: 'active', label: 'Активна' },
        { value: 'closed', label: 'Закрыта' },
      ]"
    />

    <AppButton :disabled="loading" type="submit" variant="primary">
      {{ loading ? 'Сохраняем...' : mode === 'edit' ? 'Сохранить' : 'Создать' }}
    </AppButton>
  </form>
</template>

<style scoped>
.lab-form {
  display: grid;
  gap: 12px;
  padding: 18px;
}

.lab-form__textarea-field {
  display: grid;
  gap: 6px;
}

.lab-form__textarea-field span {
  color: var(--ink-soft);
  font-size: 0.85rem;
  font-weight: 700;
}

.lab-form__textarea {
  border: 1px solid var(--line);
  border-radius: var(--radius-control);
  padding: 10px 12px;
  resize: vertical;
  min-height: 120px;
}
</style>
