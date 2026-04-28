import { defineStore } from 'pinia'

export const useNotificationsStore = defineStore('notifications', {
  state: () => ({
    items: [],
  }),

  actions: {
    show(type, message, timeout = 5000) {
      const canUseRandomUuid =
        typeof globalThis !== 'undefined' &&
        globalThis.crypto &&
        typeof globalThis.crypto.randomUUID === 'function'

      const id = canUseRandomUuid
        ? globalThis.crypto.randomUUID()
        : `toast_${Date.now()}_${Math.random().toString(16).slice(2)}`

      this.items.push({ id, type, message, timeout })

      if (timeout > 0) {
        setTimeout(() => this.dismiss(id), timeout)
      }

      return id
    },

    success(message, timeout) {
      return this.show('success', message, timeout)
    },

    warning(message, timeout) {
      return this.show('warning', message, timeout)
    },

    error(message, timeout) {
      return this.show('error', message, timeout)
    },

    dismiss(id) {
      this.items = this.items.filter((item) => item.id !== id)
    },

    clear() {
      this.items = []
    },
  },
})
