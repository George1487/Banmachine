import axios from 'axios'

import { clearTokens } from './tokenStorage'

function resolveApiBaseUrl() {
  const raw = import.meta.env.VITE_API_BASE_URL
  const explicit = typeof raw === 'string' ? raw.trim() : ''

  if (import.meta.env.DEV) {
    return explicit || '/api'
  }

  return explicit || '/api'
}

const API_BASE_URL = resolveApiBaseUrl()

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 15000,
})

let authTokenResolver = () => null
let authFailureHandler = null

const BACKEND_REASON_CODE_MAP = {
  wrong_password: 'invalid_credentials',
  user_not_found: 'invalid_credentials',
  conflict: 'analysis_job_already_running',
  not_allowed: 'forbidden',
  you_are_not_the_owner_of_this_lab: 'forbidden',
  you_do_not_have_sufficient_rights: 'forbidden',
}

export function bindAuthTokenBridge({ getAccessTokenFromStore, onAuthFailure }) {
  if (typeof getAccessTokenFromStore === 'function') {
    authTokenResolver = getAccessTokenFromStore
  }

  if (typeof onAuthFailure === 'function') {
    authFailureHandler = onAuthFailure
  }
}

function toSnakeCase(value) {
  if (!value) return ''
  return String(value)
    .trim()
    .replace(/([a-z0-9])([A-Z])/g, '$1_$2')
    .replace(/[^a-zA-Z0-9]+/g, '_')
    .replace(/^_+|_+$/g, '')
    .toLowerCase()
}

function normalizeBackendReasonCode(value) {
  const base = toSnakeCase(value)
  return BACKEND_REASON_CODE_MAP[base] || base || 'unknown_error'
}

function extractValidationMessage(payload) {
  const errors = payload?.errors
  if (!errors || typeof errors !== 'object') return null

  const firstField = Object.keys(errors)[0]
  const firstItem = firstField ? errors[firstField] : null
  if (Array.isArray(firstItem) && firstItem.length) {
    return firstItem[0]
  }

  return null
}

function normalizeApiError(error) {
  const responseData = error?.response?.data
  const nestedPayload = responseData?.error

  if (nestedPayload && typeof nestedPayload === 'object') {
    return {
      status: error?.response?.status ?? 0,
      code: nestedPayload?.code ?? 'unknown_error',
      message: nestedPayload?.message ?? error?.message ?? 'Неизвестная ошибка',
      raw: error,
    }
  }

  if (typeof responseData === 'string') {
    const code = normalizeBackendReasonCode(responseData)
    return {
      status: error?.response?.status ?? 0,
      code,
      message: responseData || error?.message || 'Неизвестная ошибка',
      raw: error,
    }
  }

  if (responseData && typeof responseData === 'object') {
    const validationMessage = extractValidationMessage(responseData)
    const objectCode =
      responseData?.code ||
      (validationMessage ? 'validation_error' : normalizeBackendReasonCode(responseData?.detail))

    return {
      status: error?.response?.status ?? 0,
      code: objectCode || 'unknown_error',
      message:
        validationMessage ||
        responseData?.message ||
        responseData?.detail ||
        responseData?.title ||
        error?.message ||
        'Неизвестная ошибка',
      raw: error,
    }
  }

  return {
    status: error?.response?.status ?? 0,
    code: 'unknown_error',
    message: error?.message ?? 'Неизвестная ошибка',
    raw: error,
  }
}

function emitGlobalApiError(normalizedError) {
  if (typeof window === 'undefined') return

  const status = normalizedError?.status
  if (status === 403 || status === 404 || status >= 500) {
    window.dispatchEvent(
      new CustomEvent('banmachine:api-error', {
        detail: normalizedError,
      }),
    )
  }
}

api.interceptors.request.use((config) => {
  const rawToken = authTokenResolver?.()
  const token = rawToken?.startsWith?.('Bearer ') ? rawToken : rawToken ? `Bearer ${rawToken}` : null

  if (token) {
    config.headers = config.headers || {}
    config.headers.Authorization = token
  }

  return config
})

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const status = error?.response?.status
    const normalizedError = normalizeApiError(error)

    if (status === 401) {
      clearTokens()
      const intendedRoute = `${window.location.pathname}${window.location.search}`
      await authFailureHandler?.({ intendedRoute })
    }

    emitGlobalApiError(normalizedError)
    return Promise.reject(normalizedError)
  },
)

export function unwrapData(response) {
  return response?.data?.data ?? response?.data
}

export function toUserMessage(error, fallback = 'Произошла ошибка') {
  if (!error) return fallback
  return error.message || fallback
}

export { api }
