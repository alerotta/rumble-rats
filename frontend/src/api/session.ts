let accessToken: string | null = null
let username: string | null = null

export function setSession(data: { token: string; username?: string }) {
  accessToken = data.token
  if (data.username !== undefined) username = data.username
}

export function clearSession() {
  accessToken = null
  username = null
}

export function getAccessToken() {
  return accessToken
}

export function getUsername() {
  return username
}

export function isLoggedIn() {
  return !!accessToken
}
