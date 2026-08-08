import { defineStore } from 'pinia'

// localStorage keys 与旧前端一致(token/role_id/name/admin)
export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    name: localStorage.getItem('name') || '',
    roleId: Number(localStorage.getItem('role_id') ?? -1),
    requirePasswordChange: false
  }),
  getters: {
    isLoggedIn: (s) => !!s.token && !isExpired(s.token),
    isAdmin: (s) => s.roleId === 0
  },
  actions: {
    setAuth(token: string, name: string, roleId: number, requirePasswordChange: boolean) {
      this.token = token
      this.name = name
      this.roleId = roleId
      this.requirePasswordChange = requirePasswordChange
      localStorage.setItem('token', token)
      localStorage.setItem('name', name)
      localStorage.setItem('role_id', String(roleId))
      if (roleId === 0) localStorage.setItem('admin', '1')
    },
    logout() {
      this.token = ''
      this.name = ''
      this.roleId = -1
      ;['token', 'role_id', 'name', 'admin'].forEach((k) => localStorage.removeItem(k))
    }
  }
})

// 本地 UI 级 JWT 过期判断(不验签,后端全量校验)
function isExpired(token: string): boolean {
  try {
    const payload = JSON.parse(atob(token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')))
    return payload.exp * 1000 < Date.now()
  } catch {
    return true
  }
}
