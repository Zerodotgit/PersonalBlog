import { defineStore } from 'pinia'
import { ref } from 'vue'

const TOKEN_KEY = 'blog_admin_token'
const USERNAME_KEY = 'blog_admin_username'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem(TOKEN_KEY) || '')
  const username = ref<string>(localStorage.getItem(USERNAME_KEY) || '')

  function setAuth(t: string, u: string) {
    token.value = t
    username.value = u
    localStorage.setItem(TOKEN_KEY, t)
    localStorage.setItem(USERNAME_KEY, u)
  }

  function logout() {
    token.value = ''
    username.value = ''
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USERNAME_KEY)
  }

  function isAuthed() {
    return !!token.value
  }

  return { token, username, setAuth, logout, isAuthed }
})
