import { api, unwrapData } from './api'

const JOB_STATUS_MAP = {
  0: 'pending',
  1: 'processing',
  2: 'done',
  3: 'failed',
  Created: 'pending',
  Pending: 'pending',
  Processing: 'processing',
  Done: 'done',
  Finished: 'done',
  Failed: 'failed',
}

function normalizeJobStatus(status) {
  return JOB_STATUS_MAP[status] ?? status?.toLowerCase() ?? status
}

function normalizeRisk(risk) {
  return risk?.toLowerCase() ?? null
}

function normalizeJob(raw) {
  if (!raw) return null
  return {
    id: raw.jobId ?? raw.id,
    labId: raw.labId,
    status: normalizeJobStatus(raw.status),
    createdBy: raw.createdBy
      ? { id: raw.createdBy.teacherId ?? raw.createdBy.id, fullName: raw.createdBy.fullName }
      : null,
    createdAt: raw.createdAt,
    startedAt: raw.startedAt,
    finishedAt: raw.finishedAt,
    errorMessage: raw.errorMessage,
  }
}

export const analysisService = {
  async runAnalysis(labId) {
    const response = await api.post(`/analysis/${labId}/run`)
    const data = unwrapData(response)
    return {
      analysisJobId: data.analysisJobId,
      labId: data.labId,
      status: normalizeJobStatus(data.status),
    }
  },

  async getJob(jobId) {
    const response = await api.get(`/analysis/${jobId}/jobs`)
    return normalizeJob(unwrapData(response))
  },

  async getAnalysis(labId) {
    const response = await api.get(`/analysis/${labId}/labs`)
    const data = unwrapData(response)
    if (!data) return null

    return {
      labId: data.labId,
      lastAnalysisJob: data.lastAnalysisJob
        ? {
            id: data.lastAnalysisJob.id,
            status: normalizeJobStatus(data.lastAnalysisJob.status),
            createdAt: data.lastAnalysisJob.createdAt,
            finishedAt: data.lastAnalysisJob.finishedAt,
          }
        : null,
      stats: data.stats ?? null,
      items: Array.isArray(data.items)
        ? data.items.map((item) => ({
            submissionId: item.submissionId,
            student: item.student
              ? { id: item.student.studentId ?? item.student.id, fullName: item.student.fullName }
              : null,
            topMatchSubmissionId: item.topMatchSubmissionId,
            topMatchScore: item.topMatchScore,
            finalScoreRiskLevel: normalizeRisk(item.finalScoreRiskLevel),
          }))
        : [],
    }
  },

  async getMatches(submissionId) {
    const response = await api.get(`/submissions/${submissionId}/matches`)
    const data = unwrapData(response)
    if (!data) return null

    return {
      submissionId: data.submissionId,
      analysisJobId: data.analysisJobId,
      matches: Array.isArray(data.matches)
        ? data.matches.map((m) => ({
            otherSubmissionId: m.otherSubmissionId,
            student: m.student
              ? { id: m.student.studentId ?? m.student.id, fullName: m.student.fullName }
              : null,
            textScore: m.textScore,
            calculationScore: m.calculationScore,
            imagesScore: m.imagesScore,
            finalScore: m.finalScore,
            riskLevel: normalizeRisk(m.riskLevel),
          }))
        : [],
    }
  },
}
