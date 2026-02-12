export type RegisterRequest = {
  username: string
  email: string
  password: string
}

export type LoginRequest = {
  username: string
  password: string
}

export type RefreshRequest = {
  refreshToken: string
}

export type AuthResponse = {
  token: string
  username: string
}

