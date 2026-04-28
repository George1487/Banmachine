import { defineStore } from 'pinia'

import { submissionService } from '@/services/submissionService'

export const useSubmissionsStore = defineStore('submissions', {
  state: () => ({
    mySubmissions: [],
    labSubmissions: [],
    currentSubmission: null,
    isLoading: false,
  }),

  actions: {
    async fetchMySubmissions() {
      this.isLoading = true

      try {
        const data = await submissionService.getMySubmissions()
        this.mySubmissions = Array.isArray(data) ? data : []
        return this.mySubmissions
      } finally {
        this.isLoading = false
      }
    },

    async fetchLabSubmissions(labId) {
      this.isLoading = true

      try {
        const data = await submissionService.getLabSubmissions(labId)
        this.labSubmissions = Array.isArray(data) ? data : []
        return this.labSubmissions
      } finally {
        this.isLoading = false
      }
    },

    async fetchSubmission(submissionId) {
      const data = await submissionService.getSubmission(submissionId)
      this.currentSubmission = data
      return data
    },

    async uploadSubmission(labId, file, onProgress) {
      const data = await submissionService.uploadSubmission(labId, file, onProgress)
      this.currentSubmission = data
      return data
    },
  },
})
