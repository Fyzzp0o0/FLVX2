<script setup lang="ts">
import { h, onMounted, ref } from 'vue'
import { NCard, NButton, NDataTable, NModal, NForm, NFormItem, NInput, NSelect, NSpace, NTag, useMessage } from 'naive-ui'
import { forwardList, forwardCreate, forwardUpdate, forwardDelete, forwardForceDelete, forwardPause, forwardResume, forwardDiagnose, forwardUpdateOrder, myTunnels } from '@/api'
import { useAuthStore } from '@/store/auth'

const message = useMessage()
const auth = useAuthStore()
const rows = ref<any[]>([])
const tunnels = ref<any[]>([])
const showModal = ref(false)
const editing = ref<any>(null)
const form = ref({ name: '', tunnelId: null as number | null, remoteAddr: '', strategy: 'fifo', inPort: 0 })

const columns = [
  { title: '名称', key: 'name' },
  { title: '入口', key: 'entry', render: (r: any) => `${r.inIp}:${r.inPort}` },
  { title: '目标', key: 'remoteAddr' },
  { title: '隧道', key: 'tunnelName' },
  { title: '流量 ↓ / ↑', key: 'flow', render: (r: any) => `${fmt(r.inFlow)} / ${fmt(r.outFlow)}` },
  { title: '状态', key: 'status', render: (r: any) => hTag(r.status === 1 ? '运行中' : '已暂停', r.status === 1 ? 'success' : 'warning') },
  { title: '操作', key: 'actions', render: (r: any) => hActions(r) }
]

const orderUp = (r: any) => { move(r, -1) }
const orderDown = (r: any) => { move(r, 1) }

function hTag(text: string, type: any) { return h(NTag, { type, size: 'small' }, { default: () => text }) }
function fmt(v: number) { return v > 1048576 ? `${(v / 1048576).toFixed(1)}MB` : `${(v / 1024).toFixed(0)}KB` }

function hActions(r: any) {
  return h(NSpace, { size: 4 }, { default: () => [
    h(NButton, { size: 'tiny', type: 'info', quaternary: true, onClick: () => edit(r) }, { default: () => '编辑' }),
    h(NButton, { size: 'tiny', type: 'primary', quaternary: true, onClick: () => orderUp(r) }, { default: () => '↑' }),
    h(NButton, { size: 'tiny', type: 'primary', quaternary: true, onClick: () => orderDown(r) }, { default: () => '↓' }),
    h(NButton, { size: 'tiny', type: r.status === 1 ? 'warning' : 'primary', quaternary: true, onClick: () => toggle(r) }, { default: () => (r.status === 1 ? '暂停' : '恢复') }),
    h(NButton, { size: 'tiny', type: 'info', quaternary: true, onClick: () => diagnose(r) }, { default: () => '诊断' }),
    h(NButton, { size: 'tiny', type: 'error', quaternary: true, onClick: () => remove(r, false) }, { default: () => '删除' }),
    ...(auth.isAdmin ? [h(NButton, { size: 'tiny', type: 'error', text: true, onClick: () => remove(r, true) }, { default: () => '强删' })] : [])
  ] })
}

function edit(r: any) {
  editing.value = r
  form.value = { name: r.name, tunnelId: r.tunnelId, remoteAddr: r.remoteAddr, strategy: r.strategy || 'fifo', inPort: 0 }
  showModal.value = true
}

async function move(r: any, dir: number) {
  const idx = rows.value.findIndex((x) => x.id === r.id)
  const target = idx + dir
  if (target < 0 || target >= rows.value.length) return
  const list = rows.value.map((x) => ({ id: x.id, inx: x.inx }))
  ;[list[idx].inx, list[target].inx] = [list[target].inx, list[idx].inx]
  try {
    await forwardUpdateOrder(list)
    rows.value = rows.value.map((x) => ({ ...x, inx: x.inx }))
    message.success('排序已更新')
    load()
  } catch (e: any) { message.error(e.message) }
}

async function load() {
  try {
    rows.value = await forwardList()
    tunnels.value = await myTunnels()
  } catch (e: any) { message.error(e.message) }
}

function openCreate() {
  editing.value = null
  form.value = { name: '', tunnelId: null, remoteAddr: '', strategy: 'fifo', inPort: 0 }
  showModal.value = true
}

async function save() {
  try {
    if (editing.value) {
      await forwardUpdate({ id: editing.value.id, userId: editing.value.userId, ...form.value })
    } else {
      await forwardCreate(form.value)
    }
    message.success('保存成功')
    showModal.value = false
    load()
  } catch (e: any) { message.error(e.message) }
}

async function toggle(r: any) {
  try {
    if (r.status === 1) await forwardPause(r.id)
    else await forwardResume(r.id)
    load()
  } catch (e: any) { message.error(e.message) }
}
async function remove(r: any, force: boolean) {
  try {
    if (force) await forwardForceDelete(r.id)
    else await forwardDelete(r.id)
    message.success('已删除')
    load()
  } catch (e: any) { message.error(e.message) }
}
async function diagnose(r: any) {
  try {
    const d = await forwardDiagnose(r.id)
    const lines = (d.results || []).map((x: any) => `${x.description}: ${x.success ? `✅ ${x.averageTime}ms` : '❌'}`)
    message.info(lines.join('\n'), { duration: 10000, closable: true })
  } catch (e: any) { message.error(e.message) }
}

onMounted(load)
</script>

<template>
  <n-card title="转发管理">
    <template #header-extra>
      <n-button type="primary" size="small" @click="openCreate">新建转发</n-button>
    </template>
    <n-data-table :columns="columns" :data="rows" :bordered="false" size="small" :pagination="{ pageSize: 20 }" />

    <n-modal v-model:show="showModal" preset="card" title="新建转发" style="width: 480px">
      <n-form label-placement="left" label-width="90px">
        <n-form-item label="名称"><n-input v-model:value="form.name" placeholder="转发名称" /></n-form-item>
        <n-form-item label="隧道">
          <n-select v-model:value="form.tunnelId" :options="tunnels.map((t) => ({ label: t.name, value: t.id }))" placeholder="选择隧道" />
        </n-form-item>
        <n-form-item label="目标地址"><n-input v-model:value="form.remoteAddr" placeholder="ip:port,多个用逗号分隔" /></n-form-item>
        <n-form-item label="策略">
          <n-select v-model:value="form.strategy" :options="[{ label: 'fifo', value: 'fifo' }, { label: 'round', value: 'round' }, { label: 'random', value: 'random' }]" />
        </n-form-item>
        <n-form-item label="指定端口"><n-input-number v-model:value="form.inPort" :min="0" :max="65535" placeholder="0=自动分配" /></n-form-item>
      </n-form>
      <template #footer><n-button type="primary" block @click="save">保存</n-button></template>
    </n-modal>
  </n-card>
</template>
