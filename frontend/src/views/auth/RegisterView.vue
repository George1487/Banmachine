<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import RegisterForm from '@/components/auth/RegisterForm.vue'
import { toUserMessage } from '@/services/api'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()

const loading = ref(false)
const formError = ref('')

function mapRegisterError(error) {
  if (error?.code === 'email_already_exists') {
    return 'Email уже зарегистрирован'
  }

  if (error?.code === 'group_name_required') {
    return 'Укажите группу для студента'
  }

  return toUserMessage(error, 'Не удалось зарегистрироваться')
}

async function handleRegister(payload) {
  formError.value = ''
  loading.value = true

  try {
    const user = await authStore.register(payload)
    await router.replace(user.role === 'teacher' ? '/teacher/labs' : '/student/labs')
  } catch (error) {
    formError.value = mapRegisterError(error)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <RegisterForm :loading="loading" :error="formError" @submit="handleRegister" />
</template>
