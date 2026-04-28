import { defineStore } from 'pinia'

import { authService } from '@/services/authService'
import { clearTokens, getAccessToken, setAccessToken } from '@/services/tokenStorage'
import router from '@/router'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    accessToken: getAccessToken(),
    isRestoring: false,
  }),

  getters: {
    isAuthenticated: (state) => Boolean(state.accessToken),
    isTeacher: (state) => state.user?.role === 'teacher',
    isStudent: (state) => state.user?.role === 'student',
  },

  actions: {
    setUser(user) {
      this.user = user
    },

    setToken(accessToken) {
      this.accessToken = accessToken || null
      if (this.accessToken) {
        setAccessToken(this.accessToken)
      } else {
        clearTokens()
      }
    },

    async register(payload) {
      await authService.register(payload)
      return this.login({ email: payload.email, password: payload.password })
    },

    async login({ email, password }) {
      const data = await authService.login({ email, password })
      this.setToken(data.accessToken)
      this.setUser(data.user)

      try {
        const me = await authService.getMe()
        this.setUser(me)
        return me
      } catch {
        return data.user
      }
    },

    async fetchMe() {
      const user = await authService.getMe()
      this.setUser(user)
      return user
    },

    async restoreSession() {
      if (this.isRestoring) return this.isAuthenticated
      this.isRestoring = true

      try {
        if (!this.accessToken) return false
        await this.fetchMe()
        return true
      } catch {
        this.clearSession()
        return false
      } finally {
        this.isRestoring = false
      }
    },

    clearSession() {
      this.user = null
      this.accessToken = null
      clearTokens()
    },

    async logout() {
      this.clearSession()
      await router.push('/login')
    },
  },
})
