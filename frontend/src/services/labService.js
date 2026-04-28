import { api, unwrapData } from './api'

function normalizeLabStatus(status) {
  if (status === 0 || status === '0' || status === 'Active' || status === 'active') return 'active'
  if (status === 1 || status === '1' || status === 'Closed' || status === 'closed') return 'closed'
  return typeof status === 'string' ? status.toLowerCase() : status
}

function toBackendLabStatus(status) {
  return normalizeLabStatus(status) === 'closed' ? 1 : 0
}

function normalizeLab(raw) {
  if (!raw) return null
  return {
    id: raw.labId ?? raw.id,
    title: raw.title,
    description: raw.description,
    status: normalizeLabStatus(raw.labStatus ?? raw.status),
    deadlineAt: raw.deadlineAt,
  }
}

function normalizeTeacherLab(raw) {
  if (!raw) return null
  return {
    id: raw.labId ?? raw.id,
    title: raw.title,
    status: normalizeLabStatus(raw.labStatus ?? raw.status),
    submissionsCount: raw.submissionCount ?? raw.submissionsCount ?? 0,
    parsedSubmissionsCount: raw.parsedSubmissionCount ?? raw.parsedSubmissionsCount ?? 0,
  }
}

export const labService = {
  async createLab(payload) {
    const response = await api.post('/labs', payload)
    return normalizeLab(unwrapData(response))
  },

  async getLabs() {
    const response = await api.get('/labs')
    const data = unwrapData(response)
    return Array.isArray(data) ? data.map(normalizeLab) : []
  },

  async getLab(labId) {
    const response = await api.get(`/labs/${labId}`)
    return normalizeLab(unwrapData(response))
  },

  async getTeacherLabs() {
    const response = await api.get('/labs/me')
    const data = unwrapData(response)
    return Array.isArray(data) ? data.map(normalizeTeacherLab) : []
  },

  async updateLab(labId, payload) {
    const requestBody = {
      labId,
      ...payload,
    }

    if (Object.prototype.hasOwnProperty.call(requestBody, 'status')) {
      requestBody.status = toBackendLabStatus(requestBody.status)
    }

    const response = await api.patch(`/labs/${labId}`, {
      ...requestBody,
    })
    return normalizeLab(unwrapData(response))
  },
}
