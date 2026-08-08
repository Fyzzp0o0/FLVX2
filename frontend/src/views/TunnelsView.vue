<script setup lang="ts">
import { h, onMounted, ref } from 'vue'
import { NCard, NButton, NDataTable, NModal, NForm, NFormItem, NInput, NSelect, NInputNumber, NTag, useMessage } from 'naive-ui'
import { tunnelList, tunnelCreate, tunnelUpdate, tunnelDelete, tunnelDiagnose, nodeList } from '@/api'

const message = useMessage()
const rows = ref<any[]>([])
const nodes = ref<any[]>([])
const showModal = ref(false)
const form = ref({ name: '', type: 1, flow: 1, trafficRatio: 1, inNodeIds: [] as number[], inIp: '' })
const editing = ref<any>(null)

const columns = [
  { title: '名称', key: 'name' },
  { title: '类型', key: 'typeView', render: (r: any) => (r.type === 1 ? '端口转发' : '隧道转发') },
  { title: '入口', key: 'inIp' },
  { title: '入口节点', key: 'inNodes', render: (r: any) => (r.inNodeId || []).map((x: any) => x.nodeId).join(',') },
  { title: '状态', key: 'status', render: (r: any) => h(NTag, { type: r.status === 1 ? 'success' : 'error', size: 'small' }, { default: () => (r.status === 1 ? '启用' : '停用') }) },
  {
    title: '操作', key: 'actions', render: (r: any) =>
      h('div', null, [
        h(NButton, { size: 'tiny', type: 'info', quaternary: true, onClick: () => edit(r) }, { default: () => '编辑' }),
        h(NButton, { size: 'tiny', type: 'info', quaternary: true, onClick: () => diagnose(r) }, { default: () => '诊断' }),
        h(NButton, { size: 'tiny', type: 'error', quaternary: true, onClick: () => remove(r) }, { default: () => '删除' })
      ])
  }
]

async function load() {
  try {
    rows.value = await tunnelList()
    nodes.value = await nodeList()
  } catch (e: any) { message.error(e.message) }
}

async function save() {
  try {
    if (editing.value) {
      await tunnelUpdate({ id: editing.value.id, name: form.value.name, flow: form.value.flow, trafficRatio: form.value.trafficRatio, inIp: form.value.inIp })
      message.success('隧道已更新')
    } else {
      await tunnelCreate({ ...form.value, inNodeId: form.value.inNodeIds.map((id) => ({ nodeId: id })) })
      message.success('隧道已创建')
    }
    showModal.value = false
    load()
  } catch (e: any) { message.error(e.message) }
}

function edit(r: any) {
  editing.value = r
  form.value = {
    name: r.name, type: r.type, flow: r.flow, trafficRatio: r.trafficRatio,
    inNodeIds: (r.inNodeId || []).map((x: any) => x.nodeId), inIp: r.inIp || ''
  }
  showModal.value = true
}

async function remove(r: any) {
  try { await tunnelDelete(r.id); message.success('已删除'); load() } catch (e: any) { message.error(e.message) }
}

async function diagnose(r: any) {
  try {
    const d = await tunnelDiagnose(r.id)
    const lines = (d.results || []).map((x: any) => `${x.description}: ${x.success ? `✅ ${x.averageTime}ms` : '❌'}`)
    message.info(lines.join('\n'), { duration: 10000, closable: true })
  } catch (e: any) { message.error(e.message) }
}

onMounted(load)
</script>

<template>
  <n-card title="隧道管理">
    <template #header-extra>
      <n-button type="primary" size="small" @click="showModal = true">新建隧道</n-button>
    </template>
    <n-data-table :columns="columns" :data="rows" :bordered="false" size="small" />

    <n-modal v-model:show="showModal" preset="card" title="新建隧道(端口转发)" style="width: 480px">
      <n-form label-placement="left" label-width="90px">
        <n-form-item label="名称"><n-input v-model:value="form.name" /></n-form-item>
        <n-form-item v-if="!editing" label="入口节点">
          <n-select v-model:value="form.inNodeIds" multiple :options="nodes.filter((n) => n.status === 1).map((n) => ({ label: `${n.name} (${n.serverIp})`, value: n.id }))" />
        </n-form-item>
        <n-form-item label="流量倍率"><n-input-number v-model:value="form.trafficRatio" :min="0" :max="100" :step="0.1" /></n-form-item>
        <n-form-item label="入站 IP"><n-input v-model:value="form.inIp" placeholder="留空自动使用节点 IP" /></n-form-item>
      </n-form>
      <template #footer><n-button type="primary" block @click="save">创建</n-button></template>
    </n-modal>
  </n-card>
</template>
