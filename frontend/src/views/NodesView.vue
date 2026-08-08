<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NCard, NButton, NDataTable, NModal, NForm, NFormItem, NInput, NTag, NSpace, useMessage } from 'naive-ui'
import { h } from 'vue'
import { nodeList, nodeCreate, nodeDelete, nodeInstall } from '@/api'

const message = useMessage()
const rows = ref<any[]>([])
const showModal = ref(false)
const form = ref({ name: '', serverIp: '', port: '', http: 1, tls: 0, socks: 1 })

function renderProto(r: any) {
  return h(NSpace, { size: 2 }, { default: () => protoTags(r) })
}
function protoTags(r: any) {
  const out: any[] = []
  const vals = [r.http, r.tls, r.socks]
  ;['http', 'tls', 'socks'].forEach((p, i) => {
    out.push(h(NTag, { size: 'tiny', type: vals[i] ? 'success' : 'default' }, { default: () => p }))
  })
  return out
}

function toggle(i: number) {
  const keys = ['http', 'tls', 'socks'] as const
  const f = form.value as any
  f[keys[i]] = f[keys[i]] === 1 ? 0 : 1
}

const columns = [
  { title: '名称', key: 'name' },
  { title: 'IP', key: 'serverIp' },
  { title: '端口段', key: 'port' },
  {
    title: '协议', key: 'proto', render: renderProto
  },
  { title: '状态', key: 'status', render: (r: any) => h(NTag, { type: r.status === 1 ? 'success' : 'default', size: 'small' }, { default: () => (r.status === 1 ? '在线' : '离线') }) },
  {
    title: '操作', key: 'actions', render: (r: any) =>
      h(NSpace, { size: 4 }, { default: () => [
        h(NButton, { size: 'tiny', type: 'info', quaternary: true, onClick: () => install(r) }, { default: () => '安装命令' }),
        h(NButton, { size: 'tiny', type: 'error', quaternary: true, onClick: () => remove(r) }, { default: () => '删除' })
      ] })
  }
]

async function load() {
  try { rows.value = await nodeList() } catch (e: any) { message.error(e.message) }
}

async function save() {
  try {
    await nodeCreate(form.value)
    message.success('节点已创建(离线,需安装 agent 上线)')
    showModal.value = false
    load()
  } catch (e: any) { message.error(e.message) }
}

async function remove(r: any) {
  try { await nodeDelete(r.id); message.success('已删除'); load() } catch (e: any) { message.error(e.message) }
}

async function install(r: any) {
  try {
    const cmd = await nodeInstall(r.id)
    message.info(cmd, { duration: 15000, closable: true })
  } catch (e: any) { message.error(e.message) }
}

onMounted(load)
</script>

<template>
  <n-card title="节点管理">
    <template #header-extra>
      <n-button type="primary" size="small" @click="showModal = true">新建节点</n-button>
    </template>
    <n-data-table :columns="columns" :data="rows" :bordered="false" size="small" />

    <n-modal v-model:show="showModal" preset="card" title="新建节点" style="width: 480px">
      <n-form label-placement="left" label-width="90px">
        <n-form-item label="名称"><n-input v-model:value="form.name" /></n-form-item>
        <n-form-item label="服务器 IP"><n-input v-model:value="form.serverIp" /></n-form-item>
        <n-form-item label="端口段"><n-input v-model:value="form.port" placeholder="如 1000-2000,3000" /></n-form-item>
        <n-form-item label="支持协议">
          <n-space>
            <n-tag v-for="(p, i) in ['http', 'tls', 'socks']" :key="p" :type="[form.http, form.tls, form.socks][i] ? 'success' : 'default'" style="cursor: pointer" @click="toggle(i)">{{ p }}</n-tag>
          </n-space>
        </n-form-item>
      </n-form>
      <template #footer><n-button type="primary" block @click="save">创建</n-button></template>
    </n-modal>
  </n-card>
</template>
