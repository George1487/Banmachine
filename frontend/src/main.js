import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import './assets/styles/main.css'

import { bindAuthTokenBridge } from '@/services/api'
import { useAuthStore } from '@/stores/auth'

async function bootstrap() {
  // Снять регистрацию старых Service Workers (MSW и др.)
  if ('serviceWorker' in navigator) {
    const registrations = await navigator.serviceWorker.getRegistrations()
    await Promise.all(registrations.map((r) => r.unregister()))
  }

  const app = createApp(App)
  const pinia = createPinia()
  const authStore = useAuthStore(pinia)

  bindAuthTokenBridge({
    getAccessTokenFromStore: () => authStore.accessToken,
    onAuthFailure: async ({ intendedRoute }) => {
      authStore.clearSession()

      if (router.currentRoute.value.path.startsWith('/login')) {
        return
      }

      await router.replace({
        path: '/login',
        query: intendedRoute ? { redirect: intendedRoute } : {},
      })
    },
  })

  app.use(pinia)
  app.use(router)

  app.mount('#app')
}

bootstrap()
