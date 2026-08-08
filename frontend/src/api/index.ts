import { post } from './http'

// ---- 用户 ----
export const login = (username: string, password: string, captchaId?: string) =>
  post<{ token: string; name: string; role_id: number; requirePasswordChange: boolean }>('/user/login', { username, password, captchaId })
export const register = (user: string, pwd: string) => post('/user/register', { user, pwd })
export const userPackage = () => post<any>('/user/package')
export const updatePassword = (d: any) => post('/user/updatePassword', d)
export const userList = () => post<any[]>('/user/list')
export const userCreate = (d: any) => post('/user/create', d)
export const userUpdate = (d: any) => post('/user/update', d)
export const userDelete = (id: number) => post('/user/delete', { id })
export const userReset = (id: number, type: number) => post('/user/reset', { id, type })

// ---- 节点 ----
export const nodeList = () => post<any[]>('/node/list')
export const nodeCreate = (d: any) => post('/node/create', d)
export const nodeUpdate = (d: any) => post('/node/update', d)
export const nodeDelete = (id: number) => post('/node/delete', { id })
export const nodeInstall = (id: number) => post<string>('/node/install', { id })

// ---- 隧道 ----
export const tunnelList = () => post<any[]>('/tunnel/list')
export const tunnelCreate = (d: any) => post('/tunnel/create', d)
export const tunnelUpdate = (d: any) => post('/tunnel/update', d)
export const tunnelDelete = (id: number) => post('/tunnel/delete', { id })
export const tunnelDiagnose = (tunnelId: number) => post<any>('/tunnel/diagnose', { tunnelId })
export const tunnelUserAssign = (d: any) => post('/tunnel/user/assign', d)
export const tunnelUserList = (userId: number) => post<any[]>('/tunnel/user/list', { userId })
export const tunnelUserRemove = (id: number) => post('/tunnel/user/remove', { id })
export const tunnelUserUpdate = (d: any) => post('/tunnel/user/update', d)
export const myTunnels = () => post<any[]>('/tunnel/user/tunnel')

// ---- 转发 ----
export const forwardList = () => post<any[]>('/forward/list')
export const forwardCreate = (d: any) => post('/forward/create', d)
export const forwardUpdate = (d: any) => post('/forward/update', d)
export const forwardDelete = (id: number) => post('/forward/delete', { id })
export const forwardForceDelete = (id: number) => post('/forward/force-delete', { id })
export const forwardPause = (id: number) => post('/forward/pause', { id })
export const forwardResume = (id: number) => post('/forward/resume', { id })
export const forwardDiagnose = (forwardId: number) => post<any>('/forward/diagnose', { forwardId })
export const forwardUpdateOrder = (forwards: { id: number; inx: number }[]) => post('/forward/update-order', { forwards })

// ---- 限速 ----
export const speedLimitList = () => post<any[]>('/speed-limit/list')
export const speedLimitCreate = (d: any) => post('/speed-limit/create', d)
export const speedLimitUpdate = (d: any) => post('/speed-limit/update', d)
export const speedLimitDelete = (id: number) => post('/speed-limit/delete', { id })

// ---- 配置 ----
export const configList = () => post<Record<string, string>>('/config/list')
export const configUpdateSingle = (name: string, value: string) => post('/config/update-single', { name, value })
