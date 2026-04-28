<script setup>
import { computed, reactive } from 'vue'

import AppButton from '@/components/button/AppButton.vue'
import PasswordInput from '@/components/form/PasswordInput.vue'
import SegmentedControl from '@/components/form/SegmentedControl.vue'
import TextInput from '@/components/form/TextInput.vue'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  error: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['submit'])

const form = reactive({
  email: '',
  password: '',
  fullName: '',
  role: 'Student',
  groupName: '',
})

const roleOptions = [
  { value: 'Student', label: 'Студент' },
  { value: 'Teacher', label: 'Преподаватель' },
]

const showGroup = computed(() => form.role === 'Student')

function handleSubmit() {
  emit('submit', {
    ...form,
    groupName: showGroup.value ? form.groupName : null,
  })
}
</script>

<template>
  <form class="auth-form" @submit.prevent="handleSubmit">
    <h1>Регистрация</h1>
    <p class="auth-form__subtitle">Создайте аккаунт для работы с BanMachine</p>

    <SegmentedControl v-model="form.role" :options="roleOptions" />

    <TextInput v-model="form.fullName" label="ФИО" placeholder="Иван Иванов" required />
    <TextInput v-model="form.email" label="Email" type="email" placeholder="you@example.com" required />
    <PasswordInput v-model="form.password" label="Пароль" placeholder="Минимум 6 символов" />
    <TextInput
      v-if="showGroup"
      v-model="form.groupName"
      label="Группа"
      placeholder="M3207"
      required
    />

    <p v-if="props.error" class="auth-form__error">{{ props.error }}</p>

    <AppButton :disabled="loading" :block="true" type="submit">
      {{ loading ? 'Создаем...' : 'Создать аккаунт' }}
    </AppButton>

    <p class="auth-form__link-row">
      Уже есть аккаунт?
      <RouterLink to="/login">Войти</RouterLink>
    </p>
  </form>
</template>

<style scoped>
.auth-form {
  max-width: 520px;
  display: grid;
  gap: 12px;
}

.auth-form__subtitle {
  color: var(--ink-soft);
  margin-bottom: 6px;
}

.auth-form__error {
  color: var(--danger);
  font-size: 0.87rem;
  font-weight: 600;
}

.auth-form__link-row {
  color: var(--ink-soft);
  font-size: 0.9rem;
}
</style>
