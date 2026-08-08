<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NCard, NButton, NDataTable, NModal, NForm, NFormItem, NInput, NInputNumber, NTag, useMessage } from 'naive-ui'
import { h } from 'vue'
import { userList, userCreate, userUpdate, userDelete, userReset, tunnelUserList, tunnelUserAssign, tunnelUserRemove, tunnelList } from '@/api'

const message = useMessage()
const rows = ref<any[]>([])
const tunnels = ref<any[]>([])
const showModal = ref(false)
const showAuthModal = ref(false)
const form = ref({ user: '', pwd: '', flow: 10, num: 5, expTime: 0, flowResetTime: 0 })
const authForm = ref({ userId: 0, tunnelId: null as number | null, flow: 10, num: 5, expTime: 0, flowResetTime: 0, speedId: null as number | null })

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '用户名', key: 'user' },
  { title: '流量', key: 'flowView', render: (r: any) => `${fmt(r.inFlow)} / ${fmt(r.outFlow)} (上限 ${r.flow}GB)` },
  { title: '转发数', key: 'num' },
  { title: '到期', key: 'expView', render: (r: any) => new Date(r.expTime).toLocaleDateString() },
  { title: '状态', key: 'status', render: (r: any) => h(NTag, { type: r.status === 1 ? 'success' : 'error', size: 'small' }, { default: () => (r.status === 1 ? '正常' : '停用') }) },
  {
    title: '操作', key: 'actions', render: (r: any) =>
      h('div', null, [
        h(NButton, { size: 'tiny', type: 'primary', quaternary: true, onClick: () => openAuth(r) }, { default: () => '隧道授权' }),
        h(NButton, { size: 'tiny', type: 'warning', quaternary: true, onClick: () => resetFlow(r) }, { default: () => '重置流量' }),
        h(NButton, { size: 'tiny', type: 'error', quaternary: true, onClick: () => remove(r) }, { default: () => '删除' })
      ])
  }
]

function fmt(v: number) { return v > 1048576 ? `${(v / 1048576).toFixed(1)}MB` : `${(v / 1024).toFixed(0)}KB` }

async function load() {
  try {
    rows.value = await userList()
    tunnels.value = await tunnelList()
  } catch (e: any) { message.error(e.message) }
}

async function save() {
  try {
    await userCreate({ ...form.value, expTime: Date.now() + 30 * 24 * 3600 * 1000, flowResetTime: 1 })
    message.success('用户已创建')
    showModal.value = false
    load()
  } catch (e: any) { message.error(e.message) }
}

async function remove(r: any) {
  try { await userDelete(r.id); message.success('已删除'); load() } catch (e: any) { message.error(e.message) }
}

async function resetFlow(r: any) {
  try { await userReset(r.id, 1); message.success('流量已重置'); load() } catch (e: any) { message.error(e.message) }
}

async function openAuth(r: any) {
  authForm.value = { userId: r.id, tunnelId: null, flow: 10, num: 5, expTime: 0, flowResetTime: 0, speedId: null }
  showAuthModal.value = true
}

async function saveAuth() {
  try {
    await tunnelUserAssign({ ...authForm.value, expTime: Date.now() + 30 * 24 * 3600 * 1000, flowResetTime: 1 })
    message.success('授权成功')
    showAuthModal.value = false
  } catch (e: any) { message.error(e.message) }
}

onMounted(load)
</script>

<template>
  <n-card title="用户管理">
    <template #header-extra>
      <n-button type="primary" size="small" @click="showModal = true">新建用户</n-button>
    </template>
    <n-data-table :columns="columns" :data="rows" :bordered="false" size="small" :pagination="{ pageSize: 20 }" />

    <n-modal v-model:show="showModal" preset="card" title="新建用户" style="width: 460px">
      <n-form label-placement="left" label-width="90px">
        <n-form-item label="用户名"><n-input v-model:value="form.user" /></n-form-item>
        <n-form-item label="密码"><n-input v-model:value="form.pwd" type="password" /></n-form-item>
        <n-form-item label="流量上限(GB)"><n-input-number v-model:value="form.flow" :min="0" /></n-form-item>
        <n-form-item label="转发上限"><n-input-number v-model:value="form.num" :min="0" /></n-form-item>
      </n-form>
      <template #footer><n-button type="primary" block @click="save">创建</n-button></template>
    </n-modal>

    <n-modal v-model:show="showAuthModal" preset="card" title="隧道授权" style="width: 460px">
      <n-form label-placement="left" label-width="90px">
        <n-form-item label="隧道">
          <n-select v-model:value="authForm.tunnelId" :options="tunnels.map((t) => ({ label: t.name, value: t.id }))" />
        </n-form-item>
        <n-form-item label="流量上限(GB)"><n-input-number v-model:value="authForm.flow" :min="0" /></n-form-item>
        <n-form-item label="转发上限"><n-input-number v-model:value="authForm.num" :min="0" /></n-form-item>
      </n-form>
      <template #footer><n-button type="primary" block @click="saveAuth">授权</n-button></template>
    </n-modal>
  </n-card>
</template>
