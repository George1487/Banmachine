import { http, HttpResponse, delay } from 'msw'
import { state } from '../state.js'
import { resolveUser, unauthorized, forbidden, notFound } from './auth.js'

const BASE = '/api/v1'

export const labsHandlers = [
  // POST /labs
  http.post(`${BASE}/labs`, async ({ request }) => {
    await delay(300)
    const user = resolveUser(request)
    if (!user) return unauthorized()
    if (user.role !== 'teacher') return forbidden()

    const { title, description, deadlineAt } = await request.json()
    const lab = {
      id: `lab_${state.nextLabId++}`,
      title,
      description,
      status: 'active',
      deadlineAt,
      createdBy: user.id,
    }
    state.labs.push(lab)

    return HttpResponse.json({
      data: {
        id: lab.id,
        title: lab.title,
        description: lab.description,
        status: lab.status,
        deadlineAt: lab.deadlineAt,
      },
    })
  }),

  // GET /labs
  http.get(`${BASE}/labs`, async ({ request }) => {
    await delay(300)
    const user = resolveUser(request)
    if (!user) return unauthorized()

    const list = state.labs.map((lab) => ({
      id: lab.id,
      title: lab.title,
      status: lab.status,
      deadlineAt: lab.deadlineAt,
    }))

    return HttpResponse.json({ data: list })
  }),

  // GET /teachers/me/labs  — должен быть ДО /labs/:labId чтобы не перехватился
  http.get(`${BASE}/teachers/me/labs`, async ({ request }) => {
    await delay(300)
    const user = resolveUser(request)
    if (!user) return unauthorized()
    if (user.role !== 'teacher') return forbidden()

    const teacherLabs = state.labs
      .filter((l) => l.createdBy === user.id)
      .map((lab) => {
        const subs = state.submissions.filter((s) => s.labId === lab.id)
        return {
          id: lab.id,
          title: lab.title,
          status: lab.status,
          submissionsCount: subs.length,
          parsedSubmissionsCount: subs.filter((s) => s.status === 'parsed').length,
        }
      })

    return HttpResponse.json({ data: teacherLabs })
  }),

  // GET /labs/:labId
  http.get(`${BASE}/labs/:labId`, async ({ params, request }) => {
    await delay(200)
    const user = resolveUser(request)
    if (!user) return unauthorized()

    const lab = state.labs.find((l) => l.id === params.labId)
    if (!lab) return notFound('lab_not_found', 'Лабораторная не найдена')

    return HttpResponse.json({
      data: {
        id: lab.id,
        title: lab.title,
        description: lab.description,
        status: lab.status,
        deadlineAt: lab.deadlineAt,
      },
    })
  }),

  // PATCH /labs/:labId
  http.patch(`${BASE}/labs/:labId`, async ({ params, request }) => {
    await delay(300)
    const user = resolveUser(request)
    if (!user) return unauthorized()
    if (user.role !== 'teacher') return forbidden()

    const lab = state.labs.find((l) => l.id === params.labId)
    if (!lab) return notFound('lab_not_found', 'Лабораторная не найдена')
    if (lab.createdBy !== user.id) return forbidden()

    const body = await request.json()
    if (body.title !== undefined) lab.title = body.title
    if (body.description !== undefined) lab.description = body.description
    if (body.status !== undefined) lab.status = body.status
    if (body.deadlineAt !== undefined) lab.deadlineAt = body.deadlineAt

    return HttpResponse.json({
      data: {
        id: lab.id,
        title: lab.title,
        description: lab.description,
        status: lab.status,
        deadlineAt: lab.deadlineAt,
      },
    })
  }),
]
