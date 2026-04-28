import { api, unwrapData } from './api'

const STATUS_MAP = {
  0: 'uploaded',
  1: 'parsing',
  2: 'parsed',
  3: 'failed',
  Uploaded: 'uploaded',
  Parsing: 'parsing',
  Parsed: 'parsed',
  Failed: 'failed',
  Pending: 'uploaded',
}

function normalizeStatus(status) {
  return STATUS_MAP[status] ?? status?.toLowerCase() ?? status
}

function normalizeRisk(risk) {
  return risk?.toLowerCase() ?? null
}

function normalizeStudent(raw) {
  if (!raw) return null
  if (typeof raw !== 'object') {
    return { id: raw, fullName: null }
  }

  // backend может вернуть "students", "student" или только "studentId"
  return {
    id: raw.userId ?? raw.studentId ?? raw.id,
    fullName: raw.fullName ?? null,
  }
}

export const submissionService = {
  async uploadSubmission(labId, file, onProgress) {
    const formData = new FormData()
    formData.append('file', file)

    const response = await api.post(`/submissions/${labId}`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (event) => {
        if (!onProgress || !event.total) return
        const percent = Math.round((event.loaded * 100) / event.total)
        onProgress(percent)
      },
    })

    const data = unwrapData(response)
    return {
      id: data.submissionId ?? data.id,
      labId: data.labId,
      studentId: data.studentId,
      status: normalizeStatus(data.status),
      submittedAt: data.submittedAt,
    }
  },

  async getMySubmissions() {
    const response = await api.get('/submissions/me')
    const data = unwrapData(response)
    return Array.isArray(data)
      ? data.map((s) => ({
          id: s.submissionId ?? s.id,
          labId: s.labId,
          labTitle: s.labTitle,
          status: normalizeStatus(s.status),
          submittedAt: s.submittedAt,
        }))
      : []
  },

  async getSubmission(submissionId) {
    const response = await api.get(`/submissions/${submissionId}`)
    const data = unwrapData(response)
    return {
      id: data.submissionId ?? data.id,
      labId: data.labId,
      student: normalizeStudent(data.students ?? data.student ?? { studentId: data.studentId }),
      status: normalizeStatus(data.status),
      submittedAt: data.submittedAt,
    }
  },

  async getLabSubmissions(labId) {
    const response = await api.get(`/submissions/labs/${labId}`)
    const data = unwrapData(response)
    return Array.isArray(data)
      ? data.map((s) => ({
          submissionId: s.submissionId ?? s.id,
          student: normalizeStudent(s.students ?? s.student),
          status: normalizeStatus(s.status),
          submittedAt: s.submittedAt,
          topMatchScore: s.topMatchScore ?? null,
          finalScoreRiskLevel: normalizeRisk(s.finalScoreRiskLevel),
        }))
      : []
  },
}
