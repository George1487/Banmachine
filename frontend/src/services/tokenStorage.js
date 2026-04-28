const ACCESS_TOKEN_KEY = 'banmachine.accessToken'
const REFRESH_TOKEN_KEY = 'banmachine.refreshToken'

function safeLocalStorage() {
  try {
    return window.localStorage
  } catch {
    return null
  }
}

export function getAccessToken() {
  return safeLocalStorage()?.getItem(ACCESS_TOKEN_KEY) ?? null
}

export function getRefreshToken() {
  return safeLocalStorage()?.getItem(REFRESH_TOKEN_KEY) ?? null
}

export function setAccessToken(token) {
  const storage = safeLocalStorage()
  if (!storage) return

  if (token) {
    storage.setItem(ACCESS_TOKEN_KEY, token)
  } else {
    storage.removeItem(ACCESS_TOKEN_KEY)
  }
}

export function setRefreshToken(token) {
  const storage = safeLocalStorage()
  if (!storage) return

  if (token) {
    storage.setItem(REFRESH_TOKEN_KEY, token)
  } else {
    storage.removeItem(REFRESH_TOKEN_KEY)
  }
}

export function saveTokens({ accessToken, refreshToken }) {
  setAccessToken(accessToken)
  setRefreshToken(refreshToken)
}

export function clearTokens() {
  const storage = safeLocalStorage()
  if (!storage) return

  storage.removeItem(ACCESS_TOKEN_KEY)
  storage.removeItem(REFRESH_TOKEN_KEY)
}
