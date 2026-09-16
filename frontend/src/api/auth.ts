import http from './http'

export interface LoginPayload {
  username: string
  password: string
}

export interface LoginResult {
  token: string
  username: string
}

export function login(payload: LoginPayload) {
  return http.post<LoginResult>('/auth/login', payload).then((r) => r.data)
}
