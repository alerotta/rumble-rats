import { getAccessToken, clearSession } from "./session"

export async function apiFetch<T>(
  input: RequestInfo,
  init: RequestInit = {}
): Promise<T> {
  const headers = new Headers(init.headers)

  if (!headers.has("Content-Type") && init.body) {
    headers.set("Content-Type", "application/json")
  }

  const token = getAccessToken()
  if (token) {
    headers.set("Authorization", `Bearer ${token}`)
  }

  const res = await fetch(input, { ...init, headers })

  if (res.status === 401) {
    clearSession()
    throw new Error("Unauthorized")
  }

  if (!res.ok) {
    let message = `Request failed (${res.status})`
    try {
      const data = await res.json()
      message = data?.message || data?.error || message
    } catch {
        //empty
    }
    throw new Error(message)
  }

  if (res.status === 204) return undefined as T

  return res.json() as Promise<T>
}
