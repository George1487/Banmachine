<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import LoginForm from '@/components/auth/LoginForm.vue'
import { toUserMessage } from '@/services/api'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()

const loading = ref(false)
const formError = ref('')

function roleRedirect(role) {
  return role === 'teacher' ? '/teacher/labs' : '/student/labs'
}

function mapLoginError(error) {
  if (error?.code === 'invalid_credentials') {
    return 'Неверный email или пароль'
  }

  return toUserMessage(error, 'Не удалось выполнить вход')
}

async function handleLogin(payload) {
  formError.value = ''
  loading.value = true

  try {
    const user = await authStore.login(payload)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : null

    await router.replace(redirect || roleRedirect(user.role))
  } catch (error) {
    formError.value = mapLoginError(error)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <LoginForm :loading="loading" :error="formError" @submit="handleLogin" />
</template>
