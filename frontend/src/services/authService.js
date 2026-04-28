import { api, unwrapData } from './api'

function normalizeRole(role) {
  if (role === 0 || role === '0' || role === 'Teacher' || role === 'teacher') return 'teacher'
  if (role === 1 || role === '1' || role === 'Student' || role === 'student') return 'student'
  return typeof role === 'string' ? role.toLowerCase() : null
}

function toBackendRole(role) {
  return normalizeRole(role) === 'teacher' ? 0 : 1
}

function normalizeUser(raw) {
  if (!raw) return null
  return {
    id: raw.userId ?? raw.id,
    email: raw.email,
    fullName: raw.fullName,
    role: normalizeRole(raw.role),
    groupName: raw.groupName ?? null,
  }
}

export const authService = {
  async register(payload) {
    const response = await api.post('/auth/register', {
      ...payload,
      role: toBackendRole(payload?.role),
    })
    return normalizeUser(unwrapData(response))
  },

  async login(payload) {
    const response = await api.post('/auth/login', payload)
    const data = unwrapData(response)
    return {
      accessToken: data.token ?? data.accessToken,
      user: normalizeUser(data.user),
    }
  },

  async getMe() {
    const response = await api.get('/auth/me')
    return normalizeUser(unwrapData(response))
  },
}
