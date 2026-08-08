# 前端 Vue 3 重写方案(已选定方案 B)

> 决策:前端由 React 18(vite-frontend)整体重写为 **Vue 3 + Vite + TypeScript**。
> 硬约束:**后端 API 契约、页面功能、交互行为完全不变**,只换实现与视觉。后端 Go 重写不受影响。

## 1. 技术选型

| 项 | 选择 | 说明 |
|---|---|---|
| 框架 | Vue 3.4+ + `<script setup>` + TypeScript | 组合式 API |
| 构建 | Vite 6(默认 minify/treeshake) | 产物 ~1.5-2.5MB,远优于现 React 的 4-6MB |
| 路由 | vue-router 4 | 11 条路由与现有路径一一对应 |
| 状态 | Pinia + localStorage 持久化 | **key 与现值保持一致**:`token`、`role_id`、`name`、`admin` |
| HTTP | axios(同现有 network.ts 行为) | baseURL = `VITE_API_BASE || ''` → 生产同源 `/api/v1/`(浏览器访问 6635,Go 同源处理);开发指向 `http://127.0.0.1:6636`;`Authorization` 裸 token;`code===401` 清 token 跳登录;30s 超时 |
| UI 组件 | **Naive UI**(推荐:TS 友好、暗色主题一流、组件现代)备选 Element Plus | 替换 @heroui |
| 样式 | Tailwind CSS 4 + 暗色模式 | 系统跟随 + 手动切换(localStorage) |
| 图表 | **ECharts**(vue-echarts) | 替换 recharts;流量曲线/24h 统计/延迟趋势 |
| 动效 | Vue Transition + 轻量库(可选 @vueuse/motion) | 替代 framer-motion |
| 拓扑图 | **xyflow 的 Vue 版(@vue-flow/core)** | 转发链「入口→中转→出口」链路可视化 |
| 验证码 | **自研滑块验证码**(Vue 组件 + Go 后端) | 替代 tianai-captcha + tac.min.js,前后端一起重写无兼容负担 |

## 2. 必须复刻的现有行为(功能不变清单)

| 行为 | 现状(React) | Vue 复刻方式 |
|---|---|---|
| 三布局切换 | `useH5Mode`:屏宽≤768 / UA 移动端 / `?h5=true` → H5 / H5Simple 布局 | 同一逻辑抽成 composable `useH5Mode` |
| 登录后跳转 | 已登录访问 `/` → `/dashboard`;`requirePasswordChange=true` → `/change-password` | router 守卫等价实现 |
| JWT 本地校验 | `utils/jwt.ts`:`atob` 解 payload 看 exp(不验签名)→ `isLoggedIn()` / `isAdmin()` | 同一实现(仅 UI 判断,后端仍全量验签) |
| 401 处理 | axios 拦截器:清 localStorage + 跳登录 | 同一拦截器 |
| 节点实时状态 | `/system-info?type=0&secret=<token>` WebSocket,2s 推送 status/info,断线指数退避重连 ≤5 次 | vue 页面内 WebSocket 封装,行为一致 |
| WebView 桥接 | `utils/panel.ts`:JsInterface / webkit.messageHandlers 获取面板地址(settings 页) | 原样保留(Android/WebView 场景) |
| 强制改密 | `/change-password` 路由 + updatePassword 流程 | 原样 |
| 页面清单 | `/`、`/dashboard`、`/node`、`/tunnel`、`/forward`、`/limit`、`/user`、`/profile`、`/change-password`、`/config`、`/settings` | **路径不变,一一对应** |
| API 调用集 | api/index.ts 约 40 个端点,全部 POST | 按 01-api.md 逐端点对齐 |

## 3. 页面 → Vue 组件映射

| 路由 | React 原文件 | Vue 组件 | 功能要点 |
|---|---|---|---|
| `/` | pages/index.tsx | `views/LoginView.vue` | 登录/注册 + 滑块验证码(新) |
| `/dashboard` | pages/dashboard.tsx | `views/DashboardView.vue` | 套餐/流量(单向/双向)、转发列表、24h 曲线(ECharts)、到期提醒 |
| `/node` | pages/node.tsx | `views/NodeView.vue` | 节点 CRUD、安装命令、WS 实时状态、系统信息 |
| `/tunnel` | pages/tunnel.tsx | `views/TunnelView.vue` | 隧道 CRUD(端口/隧道转发)、入口/出口/转发链编辑、诊断;新增拓扑图(可选) |
| `/forward` | pages/forward.tsx | `views/ForwardView.vue` | 转发 CRUD、暂停/恢复/诊断/强删、多出口、dnd 排序(换 @vueuse/dnd 或简单上下移) |
| `/limit` | pages/limit.tsx | `views/LimitView.vue` | 限速规则 CRUD |
| `/user` | pages/user.tsx | `views/UserView.vue` | 用户 CRUD、流量重置、隧道授权 |
| `/profile` `/change-password` | profile.tsx / change-password.tsx | `views/ProfileView.vue` / `views/ChangePasswordView.vue` | 资料/强制改密 |
| `/config` | config.tsx | `views/ConfigView.vue` | 站点配置(仅管理员) |
| `/settings` | settings.tsx | `views/SettingsView.vue` | WebView 面板地址管理 |

## 4. 验证码自研方案(替换 tianai)

- **后端**(Go):`/api/v1/captcha/check`、`/generate`、`/verify` 三端点保留,内部换成自研滑块:
  - generate:生成一张背景图(内置资源或程序绘制)+ 随机缺口 x 坐标,返回 `{id, image(base64), x(可选混淆)}`
  - verify:前端提交滑动距离/轨迹,后端校验偏差 ≤ 阈值;成功返回 `{validToken:id}`(沿用现有二次校验流程,有效期 120s)
  - captcha_enabled / captcha_type 配置语义保留(captcha_type 简化为 slider)
- **前端**:Vue 滑块组件(拖拽、轨迹收集),替换 tac.min.js;契约自定(前后端同步重写,无兼容负担)
- 工作量约 2-3 天;若首版赶工可先 `captcha_enabled=false`(默认即关闭)延后

## 5. 里程碑与工作量

| 步骤 | 内容 | 时长 |
|---|---|---|
| F1 | 工程搭建:Vite + Vue3 + TS + Pinia + Router + Naive UI + Tailwind + axios 封装(401/白名单) | 1 天 |
| F2 | 登录/注册 + 滑块验证码 + 三布局骨架 + router 守卫 | 2 天 |
| F3 | Dashboard(套餐/流量/ECharts 曲线)+ Profile/改密 | 2 天 |
| F4 | 节点页(CRUD + WS 实时)+ 用户页 | 2 天 |
| F5 | 隧道页(链路编辑/诊断 + 拓扑图)+ 转发页(CRUD/暂停/恢复/dnd) | 3 天 |
| F6 | 限速页 + 配置页 + settings(WebView)+ 暗色模式收尾 | 1-2 天 |
| F7 | 回归:11 路由全流程、H5/H5Simple 布局、401/WS 断线重连、与 Go 后端联调 | 1-2 天 |

总计约 **12-14 人日**,与 Go 后端(M1-M7)可并行推进,联调集中在 M6。

## 6. 与 Go 后端的接口(不变)

- 全部按 `01-api.md` 端点对接;响应 `{code, msg, ts, data}` 约定不变
- WebSocket `/system-info?type=0&secret=<token>` 行为不变
- **端口**:前端页面入口 **6635**,后端 API/WS **6636**(agent 对接);浏览器访问 `http://<IP>:6635`,前端同源请求 `/api/v1/` 由 Go 在 6635 上处理(见 04-auth.md §8)
- `dist/` 产物由 Go 面板 `go:embed` 托管,免 nginx
