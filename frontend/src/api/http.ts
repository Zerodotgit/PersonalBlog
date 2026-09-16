import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

const http = axios.create({
  // 开发环境通过 vite proxy 转发 /api 到后端；生产由 nginx 反代
  baseURL: '/api',
  timeout: 15000,
})

// 请求拦截：附带 JWT
http.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers.Authorization = `Bearer ${auth.token}`
  }
  return config
})

// 响应拦截：401 时清理登录态并跳转登录页
http.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      const auth = useAuthStore()
      auth.logout()
      if (router.currentRoute.value.path.startsWith('/admin')) {
        router.replace('/admin/login')
      }
    }
    return Promise.reject(err)
  },
)

export default http
