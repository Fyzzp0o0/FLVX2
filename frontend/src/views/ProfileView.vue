<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NForm, NFormItem, NInput, NButton, NAlert, useMessage } from 'naive-ui'
import { updatePassword } from '@/api'
import { useAuthStore } from '@/store/auth'

const message = useMessage()
const router = useRouter()
const auth = useAuthStore()

const form = ref({ newUsername: auth.name, currentPassword: '', newPassword: '', confirmPassword: '' })

async function save() {
  if (form.value.newPassword !== form.value.confirmPassword) {
    message.error('两次密码不一致')
    return
  }
  try {
    await updatePassword(form.value)
    message.success('修改成功,请重新登录')
    auth.logout()
    router.push({ name: 'login' })
  } catch (e: any) { message.error(e.message) }
}
</script>

<template>
  <n-card title="个人中心" style="max-width: 520px">
    <n-alert v-if="auth.requirePasswordChange" type="warning" style="margin-bottom: 16px">
      首次登录,请修改默认密码(admin_user / admin_user)
    </n-alert>
    <n-form label-placement="left" label-width="90px">
      <n-form-item label="用户名"><n-input v-model:value="form.newUsername" /></n-form-item>
      <n-form-item label="当前密码"><n-input v-model:value="form.currentPassword" type="password" show-password-on="click" /></n-form-item>
      <n-form-item label="新密码"><n-input v-model:value="form.newPassword" type="password" show-password-on="click" /></n-form-item>
      <n-form-item label="确认密码"><n-input v-model:value="form.confirmPassword" type="password" show-password-on="click" /></n-form-item>
      <n-button type="primary" block @click="save">保存修改</n-button>
    </n-form>
  </n-card>
</template>
