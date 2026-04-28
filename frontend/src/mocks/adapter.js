import MockAdapter from 'axios-mock-adapter'
import { api } from '@/services/api'
import { state } from './state.js'

const mock = new MockAdapter(api, { delayResponse: 300, onNoMatch: 'passthrough' })

// ─── Helpers ──────────────────────────────────────────────────────────────────

function generateToken(prefix) {
  return `${prefix}_${Date.now()}_${Math.random().toString(36).slice(2)}`
}

function normalizeRole(role) {
  if (role === 'Teacher' || role === 'teacher') return 'Teacher'
  if (role === 'Student' || role === 'student') return 'Student'
  return null
}

function normalizeLabStatus(status) {
  if (status === 'Active' || status === 'active') return 'Active'
  if (status === 'Closed' || status === 'closed') return 'Closed'
  return null
}

function resolveUser(config) {
  const auth = config.headers?.Authorization || config.headers?.authorization || ''
  const token = auth.replace('Bearer ', '')
  const userId = state.sessions[token]
  if (!userId) return null
  return state.users.find((u) => u.userId === userId) ?? null
}

function ok(data) {
  return [200, { data }]
}

function err(status, code, message) {
  return [status, { error: { code, message } }]
}

// ─── Auth ─────────────────────────────────────────────────────────────────────

mock.onPost('/auth/register').reply((config) => {
  const { email, password, fullName, role, groupName } = JSON.parse(config.data)
  const normalizedRole = normalizeRole(role)

  if (state.users.find((u) => u.email === email)) {
    return err(409, 'email_already_exists', 'Email уже зарегистрирован')
  }
  if (!normalizedRole) {
    return err(422, 'invalid_role', 'Роль должна быть Teacher или Student')
  }
  if (normalizedRole === 'Student' && !groupName) {
    return err(422, 'group_name_required', 'Укажите группу для студента')
  }
  if (normalizedRole === 'Teacher' && groupName) {
    return err(422, 'group_name_must_be_null_for_teacher', 'У преподавателя не может быть группы')
  }

  const user = {
    userId: `user_${state.nextUserId++}`,
    email, password, fullName, role: normalizedRole,
    groupName: normalizedRole === 'Student' ? groupName : null,
  }
  state.users.push(user)

  return ok({
    userId: user.userId,
    email: user.email,
    fullName: user.fullName,
    role: user.role,
    groupName: user.groupName,
  })
})

mock.onPost('/auth/login').reply((config) => {
  const { email, password } = JSON.parse(config.data)
  const user = state.users.find((u) => u.email === email && u.password === password)

  if (!user) return err(401, 'invalid_credentials', 'Неверный email или пароль')

  const token = generateToken('token')
  state.sessions[token] = user.userId

  return ok({
    token,
    user: { userId: user.userId, fullName: user.fullName, role: user.role },
  })
})

mock.onGet('/auth/me').reply((config) => {
  const user = resolveUser(config)
  if (!user) return err(401, 'unauthorized', 'Не авторизован')

  return ok({
    userId: user.userId,
    email: user.email,
    fullName: user.fullName,
    role: user.role,
    groupName: user.groupName,
  })
})

// ─── Labs ─────────────────────────────────────────────────────────────────────

mock.onPost('/labs').reply((config) => {
  const user = resolveUser(config)
  if (!user) return err(401, 'unauthorized', 'Не авторизован')
  if (user.role !== 'Teacher') return err(403, 'forbidden', 'Недостаточно прав')

  const { title, description, deadlineAt } = JSON.parse(config.data)
  const lab = {
    labId: `lab_${state.nextLabId++}`,
    title, description,
    labStatus: 'Active',
    deadlineAt,
    createdBy: user.userId,
  }
  state.labs.push(lab)

  return ok({
    labId: lab.labId,
    title: lab.title,
    description: lab.description,
    labStatus: lab.labStatus,
    deadlineAt: lab.deadlineAt,
  })
})

mock.onGet('/labs/me').reply((config) => {
  const user = resolveUser(config)
  if (!user) return err(401, 'unauthorized', 'Не авторизован')
  if (user.role !== 'Teacher') return err(403, 'forbidden', 'Недостаточно прав')

  const result = state.labs
    .filter((l) => l.createdBy === user.userId)
    .map((lab) => {
      const subs = state.submissions.filter((s) => s.labId === lab.labId)
      return {
        labId: lab.labId,
        title: lab.title,
        labStatus: lab.labStatus,
        submissionCount: subs.length,
        parsedSubmissionCount: subs.filter((s) => s.status === 'Parsed').length,
      }
    })

  return ok(result)
})

mock.onGet('/labs').reply((config) => {
  const user = resolveUser(config)
  if (!user) return err(401, 'unauthorized', 'Не авторизован')

  return ok(
    state.labs.map((l) => ({
      labId: l.labId,
      title: l.title,
      labStatus: l.labStatus,
      deadlineAt: l.deadlineAt,
    }))
  )
})

mock.onGet(/\/labs\/([^/]+)$/).reply((config) => {
  const user = resolveUser(config)
  if (!user) return err(401, 'unauthorized', 'Не авторизован')

  const labId = config.url.split('/').pop()
  const lab = state.labs.find((l) => l.labId === labId)
  if (!lab) return err(404, 'lab_not_found', 'Лабораторная не найдена')

  return ok({
    labId: lab.labId,
    title: lab.title,
    description: lab.description,
    labStatus: lab.labStatus,
    deadlineAt: lab.deadlineAt,
  })
})

mock.onPatch(/\/labs\/([^/]+)$/).reply((config) => {
  const user = resolveUser(config)
  if (!user) return err(401, 'unauthorized', 'Не авторизован')
  if (user.role !== 'Teacher') return err(403, 'forbidden', 'Недостаточно прав')

  const labId = config.url.split('/').pop()
  const lab = state.labs.find((l) => l.labId === labId)
  if (!lab) return err(404, 'lab_not_found', 'Лабораторная не найдена')
  if (lab.createdBy !== user.userId) return err(403, 'forbidden', 'Недостаточно прав')

  const body = JSON.parse(config.data)
  if (body.title !== undefined) lab.title = body.title
  if (body.description !== undefined) lab.description = body.description
  if (body.status !== undefined) {
    const normalizedStatus = normalizeLabStatus(body.status)
    if (!normalizedStatus) {
      return err(422, 'invalid_lab_status', 'Статус должен быть Active или Closed')
    }
    lab.labStatus = normalizedStatus
  }
  if (body.deadlineAt !== undefined) lab.deadlineAt = body.deadlineAt

  return ok({
    labId: lab.labId,
    title: lab.title,
    description: lab.description,
    labStatus: lab.labStatus,
    deadlineAt: lab.deadlineAt,
  })
})

// ─── Submissions ──────────────────────────────────────────────────────────────

mock.onPost(/\/submissions\/([^/]+)$/).reply((config) => {
  const user = resolveUser(config)
  if (!user) return err(401, 'unauthorized', 'Не авторизован')
  if (user.role !== 'Student') return err(403, 'forbidden', 'Недостаточно прав')

  const labId = config.url.split('/').pop()
  const lab = state.labs.find((l) => l.labId === labId)
  if (!lab) return err(404, 'lab_not_found', 'Лабораторная не найдена')
  if (lab.labStatus !== 'Active') return err(409, 'lab_not_active', 'Лабораторная закрыта для сдачи')

  const formData = config.data
  if (formData instanceof FormData) {
    const file = formData.get('file')
    if (file?.name) {
      const allowed = ['.docx', '.pdf', '.doc']
      if (!allowed.some((ext) => file.name.toLowerCase().endsWith(ext))) {
        return err(422, 'unsupported_file_type', 'Формат файла не поддерживается. Попробуйте docx или pdf')
      }
    }
  }

  const submission = {
    submissionId: `sub_${state.nextSubId++}`,
    labId,
    studentId: user.userId,
    status: 'Pending',
    submittedAt: new Date().toISOString(),
  }
  state.submissions.push(submission)

  // Имитация lifecycle: Pending → Parsed
  setTimeout(() => { submission.status = 'Parsed' }, 6000)

  return ok({
    submissionId: submission.submissionId,
    labId: submission.labId,
    studentId: submission.studentId,
    status: submission.status,
    submittedAt: submission.submittedAt,
  })
})

mock.onGet('/submissions/me').reply((config) => {
  const user = resolveUser(config)
  if (!user) return err(401, 'unauthorized', 'Не авторизован')
  if (user.role !== 'Student') return err(403, 'forbidden', 'Недостаточно прав')

  const result = state.submissions
    .filter((s) => s.studentId === user.userId)
    .map((s) => {
      const lab = state.labs.find((l) => l.labId === s.labId)
      return {
        submissionId: s.submissionId,
        labId: s.labId,
        labTitle: lab?.title || '—',
        status: s.status,
        submittedAt: s.submittedAt,
      }
    })

  return ok(result)
})

mock.onGet(/\/submissions\/([^/]+)\/matches$/).reply((config) => {
  const user = resolveUser(config)
  if (!user) return err(401, 'unauthorized', 'Не авторизован')
  if (user.role !== 'Teacher') return err(403, 'forbidden', 'Недостаточно прав')

  const submissionId = config.url.split('/')[2]
  const result = state.matches[submissionId]
  if (!result) return ok({ submissionId, analysisJobId: null, matches: [] })

  return ok(result)
})

mock.onGet(/\/submissions\/labs\/([^/]+)$/).reply((config) => {
  const user = resolveUser(config)
  if (!user) return err(401, 'unauthorized', 'Не авторизован')
  if (user.role !== 'Teacher') return err(403, 'forbidden', 'Недостаточно прав')

  const labId = config.url.split('/').pop()
  const lab = state.labs.find((l) => l.labId === labId)
  if (!lab) return err(404, 'lab_not_found', 'Лабораторная не найдена')

  const analysisResult = state.analysisResults[labId]
  const result = state.submissions
    .filter((s) => s.labId === labId)
    .map((s) => {
      const student = state.users.find((u) => u.userId === s.studentId)
      let topMatchScore = null
      let finalScoreRiskLevel = null
      if (analysisResult) {
        const item = analysisResult.items.find((i) => i.submissionId === s.submissionId)
        if (item) {
          topMatchScore = item.topMatchScore
          finalScoreRiskLevel = item.finalScoreRiskLevel
        }
      }
      return {
        submissionId: s.submissionId,
        students: { userId: student.userId, fullName: student.fullName },
        status: s.status,
        submittedAt: s.submittedAt,
        topMatchScore,
        finalScoreRiskLevel,
      }
    })

  return ok(result)
})

mock.onGet(/\/submissions\/([^/]+)$/).reply((config) => {
  const user = resolveUser(config)
  if (!user) return err(401, 'unauthorized', 'Не авторизован')

  const submissionId = config.url.split('/').pop()
  const sub = state.submissions.find((s) => s.submissionId === submissionId)
  if (!sub) return err(404, 'submission_not_found', 'Submission не найден')

  const student = state.users.find((u) => u.userId === sub.studentId)
  return ok({
    submissionId: sub.submissionId,
    labId: sub.labId,
    students: { userId: student.userId, fullName: student.fullName },
    status: sub.status,
    submittedAt: sub.submittedAt,
  })
})

// ─── Analysis ─────────────────────────────────────────────────────────────────

mock.onPost(/\/analysis\/([^/]+)\/run$/).reply((config) => {
  const user = resolveUser(config)
  if (!user) return err(401, 'unauthorized', 'Не авторизован')
  if (user.role !== 'Teacher') return err(403, 'forbidden', 'Недостаточно прав')

  const labId = config.url.split('/')[2]
  const lab = state.labs.find((l) => l.labId === labId)
  if (!lab) return err(404, 'lab_not_found', 'Лабораторная не найдена')

  const activeJob = state.analysisJobs.find(
    (j) => j.labId === labId && (j.status === 'Created' || j.status === 'Processing'),
  )
  if (activeJob) return err(409, 'analysis_job_already_running', 'Анализ уже выполняется')

  const job = {
    jobId: `job_${state.nextJobId++}`,
    labId,
    status: 'Created',
    createdBy: user.userId,
    createdAt: new Date().toISOString(),
    startedAt: null,
    finishedAt: null,
    errorMessage: null,
  }
  state.analysisJobs.push(job)
  state.jobPollCount[job.jobId] = 0

  return ok({ analysisJobId: job.jobId, labId: job.labId, status: job.status })
})

mock.onGet(/\/analysis\/([^/]+)\/jobs$/).reply((config) => {
  const user = resolveUser(config)
  if (!user) return err(401, 'unauthorized', 'Не авторизован')

  const jobId = config.url.split('/')[2]
  const job = state.analysisJobs.find((j) => j.jobId === jobId)
  if (!job) return err(404, 'not_found', 'Job не найден')

  if (job.status !== 'Finished' && job.status !== 'Failed') {
    state.jobPollCount[job.jobId] = (state.jobPollCount[job.jobId] || 0) + 1
    const count = state.jobPollCount[job.jobId]

    if (count >= 5) {
      job.status = 'Finished'
      job.finishedAt = new Date().toISOString()
    } else if (count >= 2) {
      job.status = 'Processing'
      if (!job.startedAt) job.startedAt = new Date().toISOString()
    }
  }

  const creator = state.users.find((u) => u.userId === job.createdBy)
  return ok({
    jobId: job.jobId,
    labId: job.labId,
    status: job.status,
    createdBy: { teacherId: creator.userId, fullName: creator.fullName },
    createdAt: job.createdAt,
    startedAt: job.startedAt,
    finishedAt: job.finishedAt,
    errorMessage: job.errorMessage,
  })
})

mock.onGet(/\/analysis\/([^/]+)\/labs$/).reply((config) => {
  const user = resolveUser(config)
  if (!user) return err(401, 'unauthorized', 'Не авторизован')
  if (user.role !== 'Teacher') return err(403, 'forbidden', 'Недостаточно прав')

  const labId = config.url.split('/')[2]
  const result = state.analysisResults[labId]
  if (!result) return ok({ labId, lastAnalysisJob: null, stats: null, items: [] })

  return ok({ labId, ...result })
})
