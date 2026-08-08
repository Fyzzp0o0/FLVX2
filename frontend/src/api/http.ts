import axios from 'axios'
import { useAuthStore } from '@/store/auth'
import router from '@/router'

// 生产同源 /api/v1/(Go 在 6635 上处理);开发指向 6636
export const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE ? `${import.meta.env.VITE_API_BASE}/api/v1/` : '/api/v1/',
  timeout: 30000
})

http.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) config.headers.Authorization = auth.token // 裸 token(无 Bearer)
  return config
})

http.interceptors.response.use(
  (resp) => resp,
  (err) => Promise.reject(err)
)

// 统一 R 包装解析:code===0 成功;401 清 token 跳登录
export async function post<T = any>(url: string, data?: any): Promise<T> {
  const resp = await http.post(url, data ?? {})
  const body = resp.data
  if (body && typeof body === 'object' && 'code' in body) {
    if (body.code === 0) return body.data as T
    if (body.code === 401) {
      useAuthStore().logout()
      router.push({ name: 'login' })
    }
    throw new Error(body.msg || '请求失败')
  }
  return body as T
}

export async function getRaw(url: string): Promise<any> {
  const resp = await http.get(url)
  return resp.data
}
