<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { NForm, NFormItem, NInput, NButton, NCard, NTabPane, NTabs, NAlert, useMessage } from 'naive-ui'
import { login, register } from '@/api'
import { useAuthStore } from '@/store/auth'
import { useAppStore } from '@/store/app'

const router = useRouter()
const message = useMessage()
const auth = useAuthStore()
const app = useAppStore()
app.load()

const tab = ref('login')
const loading = ref(false)
const form = ref({ username: 'admin_user', password: 'admin_user', user: '', pwd: '', confirm: '' })

async function doLogin() {
  loading.value = true
  try {
    const d = await login(form.value.username, form.value.password)
    auth.setAuth(d.token, d.name, d.role_id, d.requirePasswordChange)
    message.success('登录成功')
    router.push(d.requirePasswordChange ? '/profile' : '/dashboard')
  } catch (e: any) {
    message.error(e.message)
  } finally {
    loading.value = false
  }
}

async function doRegister() {
  if (form.value.pwd !== form.value.confirm) {
    message.error('两次密码不一致')
    return
  }
  loading.value = true
  try {
    await register(form.value.user, form.value.pwd)
    message.success('注册成功,请登录')
    tab.value = 'login'
  } catch (e: any) {
    message.error(e.message)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-bg">
    <n-card class="login-card" :bordered="false">
      <h1 class="brand">{{ app.appName }}</h1>
      <n-tabs v-model:value="tab" type="line" animated>
        <n-tab-pane name="login" tab="登录">
          <n-form @keyup.enter="doLogin">
            <n-form-item label="用户名">
              <n-input v-model:value="form.username" placeholder="用户名" />
            </n-form-item>
            <n-form-item label="密码">
              <n-input v-model:value="form.password" type="password" show-password-on="click" placeholder="密码" />
            </n-form-item>
            <n-button type="primary" block :loading="loading" @click="doLogin">登 录</n-button>
          </n-form>
        </n-tab-pane>
        <n-tab-pane name="register" tab="注册">
          <n-form @keyup.enter="doRegister">
            <n-form-item label="用户名">
              <n-input v-model:value="form.user" placeholder="3-20 位" />
            </n-form-item>
            <n-form-item label="密码">
              <n-input v-model:value="form.pwd" type="password" show-password-on="click" placeholder="6-32 位" />
            </n-form-item>
            <n-form-item label="确认密码">
              <n-input v-model:value="form.confirm" type="password" show-password-on="click" />
            </n-form-item>
            <n-button type="primary" block :loading="loading" @click="doRegister">注 册</n-button>
          </n-form>
        </n-tab-pane>
      </n-tabs>
    </n-card>
  </div>
</template>

<style scoped>
.login-bg {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
}
.login-card {
  width: 380px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.97);
}
.brand {
  text-align: center;
  letter-spacing: 0.2em;
  margin: 8px 0 20px;
  color: #0f172a;
}
</style>
