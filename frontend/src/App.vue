<script setup lang="ts">
import { onMounted, computed, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NLayout, NLayoutSider, NLayoutHeader, NLayoutContent, NMenu, NButton, NTag, NSpace, NIcon,
  NMessageProvider, NDialogProvider, NNotificationProvider
} from 'naive-ui'
import type { MenuOption } from 'naive-ui'
import { LayoutDashboard, ArrowLeftRight, Server, GitBranch, Users, Gauge, Settings, UserCircle2 } from 'lucide-vue-next'
import { useAuthStore } from '@/store/auth'
import { useAppStore } from '@/store/app'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const app = useAppStore()
onMounted(() => app.load())

// 登录页全屏(无侧边栏),其余页面走主布局
const isLogin = computed(() => route.name === 'login')
const activeKey = computed(() => route.path)

const renderIcon = (icon: any) => () => h(NIcon, null, { default: () => h(icon) })

// 必须 computed:登录前后 roleId 变化时菜单动态重建(否则管理员菜单不出现)
const menuOptions = computed<MenuOption[]>(() => [
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
])

function logout() {
  auth.logout()
  router.push({ name: 'login' })
}
</script>

<template>
  <n-message-provider>
    <n-dialog-provider>
      <n-notification-provider>
        <router-view v-if="isLogin" />
        <n-layout v-else has-sider style="min-height: 100vh">
          <n-layout-sider bordered collapse-mode="width" :width="200" :collapsed-width="64">
            <div class="logo">{{ app.appName }}</div>
            <n-menu v-model:value="activeKey" :options="menuOptions" :collapsed-width="64" @update:value="(k) => router.push(k as string)" />
          </n-layout-sider>
          <n-layout>
            <n-layout-header bordered class="header">
              <span>{{ app.appName }} Panel</span>
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
      </n-notification-provider>
    </n-dialog-provider>
  </n-message-provider>
</template>

<style scoped>
.logo { padding: 16px 20px; font-size: 20px; font-weight: 700; letter-spacing: 0.15em; }
.header { display: flex; justify-content: space-between; align-items: center; padding: 0 16px; height: 48px; }
.content { padding: 16px; }
</style>
