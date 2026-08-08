<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  NLayout, NLayoutSider, NLayoutHeader, NLayoutContent, NMenu, NButton, NTag,
  NGrid, NGridItem, NCard, NStatistic, NDataTable, NSpace, useMessage
} from 'naive-ui'
import type { MenuOption } from 'naive-ui'
import { DashboardOutlined, SwapOutlined, ClusterOutlined, ApartmentOutlined, TeamOutlined, SpeedOutlined, SettingOutlined, UserOutlined } from '@vicons/antd'
import { NIcon } from 'naive-ui'
import { h } from 'vue'
import { userPackage, forwardList, forwardPause, forwardResume, forwardDelete, forwardDiagnose } from '@/api'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const message = useMessage()
const auth = useAuthStore()

const renderIcon = (icon: any) => () => h(NIcon, null, { default: () => h(icon) })

const menuOptions: MenuOption[] = [
  { label: '仪表盘', key: '/dashboard', icon: renderIcon(LayoutDashboard) },
  { label: '我的转发', key: '/forward', icon: renderIcon(ArrowLeftRight) },
  ...(auth.isAdmin
    ? [
        { label: '节点管理', key: '/node', icon: renderIcon(Server) },
        { label: '隧道管理', key: '/tunnel', icon: renderIcon(GitBranch) },
        { label: '用户管理', key: '/user', icon: renderIcon(Users) },
        { label: '限速规则', key: '/limit', icon: renderIcon(Gauge) },
        { label: '网站配置', key: '/config', icon: renderIcon(Settings) }
      ]
    : []),
  { label: '个人中心', key: '/profile', icon: renderIcon(UserCircle2) }
]

const activeKey = ref(router.currentRoute.value.path)

const pkg = ref<any>(null)
const forwards = ref<any[]>([])
const columns = [
  { title: '名称', key: 'name' },
  { title: '入口', key: 'entry', render: (r: any) => `${r.inIp}:${r.inPort}` },
  { title: '目标', key: 'remoteAddr' },
  { title: '隧道', key: 'tunnelName' },
  { title: '流量 ↓/↑', key: 'flow', render: (r: any) => `${fmtBytes(r.inFlow)} / ${fmtBytes(r.outFlow)}` },
  { title: '状态', key: 'status', render: (r: any) => h(NTag, { type: r.status === 1 ? 'success' : 'warning', size: 'small' }, { default: () => (r.status === 1 ? '运行中' : '已暂停') }) },
  {
    title: '操作', key: 'actions', render: (r: any) =>
      h(NSpace, { size: 4 }, {
        default: () => [
          h(NButton, { size: 'tiny', type: 'primary', quaternary: true, onClick: () => doPause(r) }, { default: () => (r.status === 1 ? '暂停' : '恢复') }),
          h(NButton, { size: 'tiny', type: 'warning', quaternary: true, onClick: () => doDiagnose(r) }, { default: () => '诊断' }),
          h(NButton, { size: 'tiny', type: 'error', quaternary: true, onClick: () => doDelete(r) }, { default: () => '删除' })
        ]
      })
  }
]

function fmtBytes(v: number) {
  if (!v) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0
  while (v >= 1024 && i < units.length - 1) { v /= 1024; i++ }
  return `${v.toFixed(1)} ${units[i]}`
}

async function load() {
  try {
    pkg.value = await userPackage()
    forwards.value = await forwardList()
  } catch (e: any) {
    message.error(e.message)
  }
}

async function doPause(r: any) {
  try {
    if (r.status === 1) { await forwardPause(r.id); message.success('已暂停') }
    else { await forwardResume(r.id); message.success('已恢复') }
    load()
  } catch (e: any) { message.error(e.message) }
}
async function doDelete(r: any) {
  try { await forwardDelete(r.id); message.success('已删除'); load() } catch (e: any) { message.error(e.message) }
}
async function doDiagnose(r: any) {
  try {
    const d = await forwardDiagnose(r.id)
    const lines = (d.results || []).map((x: any) => `${x.description}: ${x.success ? `✅ ${x.averageTime}ms` : '❌ ' + (x.message || '失败')}`)
    message.info(lines.join('\n'), { duration: 8000, closable: true })
  } catch (e: any) { message.error(e.message) }
}

function logout() {
  auth.logout()
  router.push({ name: 'login' })
}

onMounted(load)
</script>

<template>
  <n-layout has-sider style="min-height: 100vh">
    <n-layout-sider bordered collapse-mode="width" :width="200" :collapsed-width="64">
      <div class="logo">FLVX2</div>
      <n-menu v-model:value="activeKey" :options="menuOptions" :collapsed-width="64" @update:value="(k) => router.push(k as string)" />
    </n-layout-sider>
    <n-layout>
      <n-layout-header bordered class="header">
        <span>FLVX2 Panel</span>
        <n-space align="center">
          <n-tag size="small" :type="auth.isAdmin ? 'error' : 'default'">{{ auth.isAdmin ? '管理员' : '用户' }}</n-tag>
          <n-button size="small" quaternary @click="logout">退出</n-button>
        </n-space>
      </n-layout-header>
      <n-layout-content class="content">
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<style scoped>
.logo { padding: 16px 20px; font-size: 20px; font-weight: 700; letter-spacing: 0.15em; }
.header { display: flex; justify-content: space-between; align-items: center; padding: 0 16px; height: 48px; }
.content { padding: 16px; }
</style>
