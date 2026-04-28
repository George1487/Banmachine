<script setup>
import { reactive } from 'vue'

import AppButton from '@/components/button/AppButton.vue'
import PasswordInput from '@/components/form/PasswordInput.vue'
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
})

function handleSubmit() {
  emit('submit', { ...form })
}
</script>

<template>
  <form class="auth-form" @submit.prevent="handleSubmit">
    <h1>Вход</h1>
    <p class="auth-form__subtitle">Войдите в BanMachine</p>

    <TextInput v-model="form.email" label="Email" type="email" placeholder="you@example.com" required />
    <PasswordInput v-model="form.password" label="Пароль" placeholder="Введите пароль" />

    <p v-if="props.error" class="auth-form__error">{{ props.error }}</p>

    <AppButton :disabled="loading" :block="true" type="submit">
      {{ loading ? 'Входим...' : 'Войти' }}
    </AppButton>

    <p class="auth-form__link-row">
      Нет аккаунта?
      <RouterLink to="/register">Зарегистрироваться</RouterLink>
    </p>
  </form>
</template>

<style scoped>
.auth-form {
  max-width: 460px;
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
