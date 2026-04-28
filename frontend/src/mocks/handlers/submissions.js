import { http, HttpResponse, delay } from 'msw'
import { state } from '../state.js'
import { resolveUser, unauthorized, forbidden, notFound } from './auth.js'

const BASE = '/api/v1'

export const submissionsHandlers = [
  // POST /labs/:labId/submissions
  http.post(`${BASE}/labs/:labId/submissions`, async ({ params, request }) => {
    await delay(800)
    const user = resolveUser(request)
    if (!user) return unauthorized()
    if (user.role !== 'student') return forbidden()

    const lab = state.labs.find((l) => l.id === params.labId)
    if (!lab) return notFound('lab_not_found', 'Лабораторная не найдена')

    if (lab.status !== 'active') {
      return HttpResponse.json(
        { error: { code: 'lab_not_active', message: 'Лабораторная закрыта для сдачи' } },
        { status: 409 },
      )
    }

    // Имитация: принимаем только .docx / .pdf — проверяем Content-Type или имя файла
    const formData = await request.formData()
    const file = formData.get('file')
    if (file && file.name) {
      const allowed = ['.docx', '.pdf', '.doc']
      const hasAllowed = allowed.some((ext) => file.name.toLowerCase().endsWith(ext))
      if (!hasAllowed) {
        return HttpResponse.json(
          { error: { code: 'unsupported_file_type', message: 'Формат файла не поддерживается. Попробуйте docx или pdf' } },
          { status: 422 },
        )
      }
    }

    const submission = {
      id: `sub_${state.nextSubId++}`,
      labId: params.labId,
      studentId: user.id,
      status: 'uploaded',
      submittedAt: new Date().toISOString(),
    }
    state.submissions.push(submission)

    // Имитация lifecycle: uploaded → parsing → parsed
    setTimeout(() => {
      submission.status = 'parsing'
      setTimeout(() => {
        submission.status = 'parsed'
      }, 3000)
    }, 3000)

    return HttpResponse.json({
      data: {
        id: submission.id,
        labId: submission.labId,
        studentId: submission.studentId,
        status: submission.status,
        submittedAt: submission.submittedAt,
      },
    })
  }),

  // GET /students/me/submissions
  http.get(`${BASE}/students/me/submissions`, async ({ request }) => {
    await delay(300)
    const user = resolveUser(request)
    if (!user) return unauthorized()
    if (user.role !== 'student') return forbidden()

    const mySubs = state.submissions
      .filter((s) => s.studentId === user.id)
      .map((s) => {
        const lab = state.labs.find((l) => l.id === s.labId)
        return {
          id: s.id,
          labId: s.labId,
          labTitle: lab ? lab.title : 'Неизвестная лабораторная',
          status: s.status,
          submittedAt: s.submittedAt,
        }
      })

    return HttpResponse.json({ data: mySubs })
  }),

  // GET /submissions/:submissionId
  http.get(`${BASE}/submissions/:submissionId`, async ({ params, request }) => {
    await delay(200)
    const user = resolveUser(request)
    if (!user) return unauthorized()

    const sub = state.submissions.find((s) => s.id === params.submissionId)
    if (!sub) return notFound('submission_not_found', 'Submission не найден')

    const student = state.users.find((u) => u.id === sub.studentId)

    return HttpResponse.json({
      data: {
        id: sub.id,
        labId: sub.labId,
        student: { id: student.id, fullName: student.fullName },
        status: sub.status,
        submittedAt: sub.submittedAt,
      },
    })
  }),

  // GET /labs/:labId/submissions
  http.get(`${BASE}/labs/:labId/submissions`, async ({ params, request }) => {
    await delay(300)
    const user = resolveUser(request)
    if (!user) return unauthorized()
    if (user.role !== 'teacher') return forbidden()

    const lab = state.labs.find((l) => l.id === params.labId)
    if (!lab) return notFound('lab_not_found', 'Лабораторная не найдена')

    const analysisResult = state.analysisResults[params.labId]

    const subs = state.submissions
      .filter((s) => s.labId === params.labId)
      .map((s) => {
        const student = state.users.find((u) => u.id === s.studentId)
        let topMatchScore = null
        let finalScoreRiskLevel = null

        if (analysisResult) {
          const item = analysisResult.items.find((i) => i.submissionId === s.id)
          if (item) {
            topMatchScore = item.topMatchScore
            finalScoreRiskLevel = item.finalScoreRiskLevel
          }
        }

        return {
          submissionId: s.id,
          student: { id: student.id, fullName: student.fullName },
          status: s.status,
          submittedAt: s.submittedAt,
          topMatchScore,
          finalScoreRiskLevel,
        }
      })

    return HttpResponse.json({ data: subs })
  }),
]
