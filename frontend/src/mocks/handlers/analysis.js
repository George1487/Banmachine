import { http, HttpResponse, delay } from 'msw'
import { state } from '../state.js'
import { resolveUser, unauthorized, forbidden, notFound } from './auth.js'

const BASE = '/api/v1'

// Через сколько poll-запросов job переходит в processing / done
const POLLS_TO_PROCESSING = 2
const POLLS_TO_DONE = 5

export const analysisHandlers = [
  // POST /labs/:labId/analysis/run
  http.post(`${BASE}/labs/:labId/analysis/run`, async ({ params, request }) => {
    await delay(400)
    const user = resolveUser(request)
    if (!user) return unauthorized()
    if (user.role !== 'teacher') return forbidden()

    const lab = state.labs.find((l) => l.id === params.labId)
    if (!lab) return notFound('lab_not_found', 'Лабораторная не найдена')

    const activeJob = state.analysisJobs.find(
      (j) => j.labId === params.labId && (j.status === 'pending' || j.status === 'processing'),
    )
    if (activeJob) {
      return HttpResponse.json(
        { error: { code: 'analysis_job_already_running', message: 'Анализ уже выполняется' } },
        { status: 409 },
      )
    }

    const job = {
      id: `job_${state.nextJobId++}`,
      labId: params.labId,
      status: 'pending',
      createdBy: user.id,
      createdAt: new Date().toISOString(),
      startedAt: null,
      finishedAt: null,
      errorMessage: null,
    }
    state.analysisJobs.push(job)
    state.jobPollCount[job.id] = 0

    return HttpResponse.json({
      data: { analysisJobId: job.id, labId: job.labId, status: job.status },
    })
  }),

  // GET /analysis/jobs/:jobId
  http.get(`${BASE}/analysis/jobs/:jobId`, async ({ params, request }) => {
    await delay(300)
    const user = resolveUser(request)
    if (!user) return unauthorized()

    const job = state.analysisJobs.find((j) => j.id === params.jobId)
    if (!job) return notFound('not_found', 'Job не найден')

    // Имитация прогресса по числу poll-запросов
    if (job.status !== 'done' && job.status !== 'failed') {
      state.jobPollCount[job.id] = (state.jobPollCount[job.id] || 0) + 1
      const count = state.jobPollCount[job.id]

      if (count >= POLLS_TO_DONE) {
        job.status = 'done'
        job.finishedAt = new Date().toISOString()
      } else if (count >= POLLS_TO_PROCESSING) {
        job.status = 'processing'
        if (!job.startedAt) job.startedAt = new Date().toISOString()
      }
    }

    const creator = state.users.find((u) => u.id === job.createdBy)

    return HttpResponse.json({
      data: {
        id: job.id,
        labId: job.labId,
        status: job.status,
        createdBy: { id: creator.id, fullName: creator.fullName },
        createdAt: job.createdAt,
        startedAt: job.startedAt,
        finishedAt: job.finishedAt,
        errorMessage: job.errorMessage,
      },
    })
  }),

  // GET /labs/:labId/analysis
  http.get(`${BASE}/labs/:labId/analysis`, async ({ params, request }) => {
    await delay(300)
    const user = resolveUser(request)
    if (!user) return unauthorized()
    if (user.role !== 'teacher') return forbidden()

    const result = state.analysisResults[params.labId]
    if (!result) {
      return HttpResponse.json({
        data: { labId: params.labId, lastAnalysisJob: null, stats: null, items: [] },
      })
    }

    return HttpResponse.json({ data: { labId: params.labId, ...result } })
  }),

  // GET /submissions/:submissionId/matches
  http.get(`${BASE}/submissions/:submissionId/matches`, async ({ params, request }) => {
    await delay(300)
    const user = resolveUser(request)
    if (!user) return unauthorized()
    if (user.role !== 'teacher') return forbidden()

    const result = state.matches[params.submissionId]
    if (!result) {
      return HttpResponse.json({
        data: { submissionId: params.submissionId, analysisJobId: null, matches: [] },
      })
    }

    return HttpResponse.json({ data: result })
  }),
]
