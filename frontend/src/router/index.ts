import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'login', component: () => import('@/views/LoginView.vue') },
    { path: '/dashboard', name: 'dashboard', component: () => import('@/views/DashboardView.vue'), meta: { auth: true } },
    { path: '/forward', name: 'forward', component: () => import('@/views/ForwardsView.vue'), meta: { auth: true } },
    { path: '/node', name: 'node', component: () => import('@/views/NodesView.vue'), meta: { auth: true, admin: true } },
    { path: '/tunnel', name: 'tunnel', component: () => import('@/views/TunnelsView.vue'), meta: { auth: true, admin: true } },
    { path: '/user', name: 'user', component: () => import('@/views/UsersView.vue'), meta: { auth: true, admin: true } },
    { path: '/limit', name: 'limit', component: () => import('@/views/LimitsView.vue'), meta: { auth: true, admin: true } },
    { path: '/config', name: 'config', component: () => import('@/views/ConfigView.vue'), meta: { auth: true, admin: true } },
    { path: '/profile', name: 'profile', component: () => import('@/views/ProfileView.vue'), meta: { auth: true } },
    { path: '/:pathMatch(.*)*', redirect: '/' }
  ]
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.auth && !auth.isLoggedIn) return { name: 'login' }
  if (to.meta.admin && auth.roleId !== 0) return { name: 'dashboard' }
  if (to.name === 'login' && auth.isLoggedIn) return { name: 'dashboard' }
})

export default router
