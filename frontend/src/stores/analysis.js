import { defineStore } from 'pinia'

import { analysisService } from '@/services/analysisService'

export const useAnalysisStore = defineStore('analysis', {
  state: () => ({
    currentJob: null,
    analysisResult: null,
    matches: null,
    isLoading: false,
  }),

  actions: {
    async runAnalysis(labId) {
      const data = await analysisService.runAnalysis(labId)
      this.currentJob = data
      return data
    },

    async pollJob(jobId) {
      const data = await analysisService.getJob(jobId)
      this.currentJob = data
      return data
    },

    async fetchAnalysis(labId) {
      this.isLoading = true

      try {
        const data = await analysisService.getAnalysis(labId)
        this.analysisResult = data
        return data
      } finally {
        this.isLoading = false
      }
    },

    async fetchMatches(submissionId) {
      this.isLoading = true

      try {
        const data = await analysisService.getMatches(submissionId)
        this.matches = data
        return data
      } finally {
        this.isLoading = false
      }
    },

    clearJob() {
      this.currentJob = null
    },
  },
})
