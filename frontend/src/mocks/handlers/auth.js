import { http, HttpResponse, delay } from 'msw'
import { state } from '../state.js'

const BASE = '/api/v1'

function generateToken(prefix) {
  return `${prefix}_${Date.now()}_${Math.random().toString(36).slice(2)}`
}

export function resolveUser(request) {
  const auth = request.headers.get('Authorization')
  if (!auth) return null
  const token = auth.replace('Bearer ', '')
  const userId = state.sessions[token]
  if (!userId) return null
  return state.users.find((u) => u.id === userId) ?? null
}

export function unauthorized() {
  return HttpResponse.json(
    { error: { code: 'unauthorized', message: 'Не авторизован' } },
    { status: 401 },
  )
}

export function forbidden() {
  return HttpResponse.json(
    { error: { code: 'forbidden', message: 'Недостаточно прав' } },
    { status: 403 },
  )
}

export function notFound(code, message) {
  return HttpResponse.json(
    { error: { code, message } },
    { status: 404 },
  )
}

export const authHandlers = [
  // POST /auth/register
  http.post(`${BASE}/auth/register`, async ({ request }) => {
    await delay(300)
    const { email, password, fullName, role, groupName } = await request.json()

    if (state.users.find((u) => u.email === email)) {
      return HttpResponse.json(
        { error: { code: 'email_already_exists', message: 'Email уже зарегистрирован' } },
        { status: 409 },
      )
    }
    if (role === 'student' && !groupName) {
      return HttpResponse.json(
        { error: { code: 'group_name_required', message: 'Укажите группу для студента' } },
        { status: 422 },
      )
    }
    if (role === 'teacher' && groupName) {
      return HttpResponse.json(
        { error: { code: 'group_name_must_be_null_for_teacher', message: 'У преподавателя не может быть группы' } },
        { status: 422 },
      )
    }

    const user = {
      id: `user_${state.nextUserId++}`,
      email,
      password,
      fullName,
      role,
      groupName: role === 'student' ? groupName : null,
    }
    state.users.push(user)

    return HttpResponse.json({
      data: {
        id: user.id,
        email: user.email,
        fullName: user.fullName,
        role: user.role,
        groupName: user.groupName,
      },
    })
  }),

  // POST /auth/login
  http.post(`${BASE}/auth/login`, async ({ request }) => {
    await delay(400)
    const { email, password } = await request.json()
    const user = state.users.find((u) => u.email === email && u.password === password)

    if (!user) {
      return HttpResponse.json(
        { error: { code: 'invalid_credentials', message: 'Неверный email или пароль' } },
        { status: 401 },
      )
    }

    const accessToken = generateToken('access')
    const refreshToken = generateToken('refresh')
    state.sessions[accessToken] = user.id
    state.sessions[refreshToken] = user.id

    return HttpResponse.json({
      data: {
        accessToken,
        refreshToken,
        user: { id: user.id, fullName: user.fullName, role: user.role },
      },
    })
  }),

  // POST /auth/refresh
  http.post(`${BASE}/auth/refresh`, async ({ request }) => {
    await delay(200)
    const { refreshToken } = await request.json()
    const userId = state.sessions[refreshToken]

    if (!userId) {
      return HttpResponse.json(
        { error: { code: 'invalid_token', message: 'Невалидный refresh token' } },
        { status: 401 },
      )
    }

    const newAccessToken = generateToken('access')
    state.sessions[newAccessToken] = userId

    return HttpResponse.json({ data: { accessToken: newAccessToken } })
  }),

  // GET /me
  http.get(`${BASE}/me`, async ({ request }) => {
    await delay(200)
    const user = resolveUser(request)
    if (!user) return unauthorized()

    return HttpResponse.json({
      data: {
        id: user.id,
        email: user.email,
        fullName: user.fullName,
        role: user.role,
        groupName: user.groupName,
      },
    })
  }),
]
