import { defineStore } from 'pinia'

import { labService } from '@/services/labService'

export const useLabsStore = defineStore('labs', {
  state: () => ({
    labs: [],
    teacherLabs: [],
    currentLab: null,
    isLoading: false,
  }),

  actions: {
    async fetchLabs() {
      this.isLoading = true

      try {
        const data = await labService.getLabs()
        this.labs = Array.isArray(data) ? data : []
        return this.labs
      } finally {
        this.isLoading = false
      }
    },

    async fetchTeacherLabs() {
      this.isLoading = true

      try {
        const data = await labService.getTeacherLabs()
        this.teacherLabs = Array.isArray(data) ? data : []
        return this.teacherLabs
      } finally {
        this.isLoading = false
      }
    },

    async fetchLab(labId) {
      this.isLoading = true

      try {
        const data = await labService.getLab(labId)
        this.currentLab = data
        return data
      } finally {
        this.isLoading = false
      }
    },

    async createLab(payload) {
      const data = await labService.createLab(payload)
      this.teacherLabs.unshift(data)
      return data
    },

    async updateLab(labId, payload) {
      const data = await labService.updateLab(labId, payload)
      this.currentLab = data

      this.teacherLabs = this.teacherLabs.map((item) =>
        item.id === labId ? { ...item, ...data } : item,
      )

      return data
    },
  },
})
