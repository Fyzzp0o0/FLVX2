<script setup lang="ts">
import { h, onMounted, ref } from 'vue'
import { NCard, NButton, NDataTable, NModal, NForm, NFormItem, NInput, NSelect, NInputNumber, NTag, NRadioGroup, NRadioButton, NSpace, NDescriptions, NDescriptionsItem, useMessage } from 'naive-ui'
import { tunnelList, tunnelCreate, tunnelUpdate, tunnelDelete, tunnelDiagnose, nodeList } from '@/api'

const message = useMessage()
const rows = ref<any[]>([])
const nodes = ref<any[]>([])
const showModal = ref(false)
const showDiagModal = ref(false)
const diag = ref<any>(null)
const diagLoading = ref(false)
const editing = ref<any>(null)

interface TunnelForm {
  name: string; type: number; flow: number; trafficRatio: number; inIp: string
  inNodeId: number[]; chainNodes: number[][]; outNodeId: number[]
}
const form = ref<TunnelForm>({ name: '', type: 1, flow: 1, trafficRatio: 1, inIp: '', inNodeId: [], chainNodes: [[]], outNodeId: [] })
// 每组转发链的协议/策略(组级),与节点选择分开维护
const chainMeta = ref<{ protocol: string; strategy: string }[]>([])

const nodeOptions = () => nodes.value.map((n) => ({ label: `${n.name} (${n.serverIp})${n.status === 1 ? '' : ' [离线]'}`, value: n.id }))
const chainProtocols = [
  { label: 'tls', value: 'tls' }, { label: 'wss', value: 'wss' }, { label: 'tcp', value: 'tcp' },
  { label: 'mtls', value: 'mtls' }, { label: 'mwss', value: 'mwss' }, { label: 'mtcp', value: 'mtcp' }
]
const chainStrategies = [
  { label: 'fifo(顺序)', value: 'fifo' }, { label: 'round(轮询)', value: 'round' }, { label: 'rand(随机)', value: 'rand' }
]

const columns = [
  { title: '名称', key: 'name' },
  { title: '类型', key: 'typeView', render: (r: any) => h(NTag, { type: r.type === 1 ? 'default' : 'info', size: 'small' }, { default: () => (r.type === 1 ? '端口转发' : '隧道转发') }) },
  { title: '入站 IP', key: 'inIp' },
  {
    title: '链路', key: 'links', render: (r: any) => {
      if (r.type === 1) return '—'
      const hop = (r.chainNodes || []).length
      const out = (r.outNodeId || []).map((x: any) => x.nodeId).join(',')
      return `入口${(r.inNodeId || []).length} · ${hop}跳 · 出口${out || '无'}`
    }
  },
  { title: '状态', key: 'status', render: (r: any) => h(NTag, { type: r.status === 1 ? 'success' : 'error', size: 'small' }, { default: () => (r.status === 1 ? '启用' : '停用') }) },
  {
    title: '操作', key: 'actions', render: (r: any) =>
      h('div', null, [
        h(NButton, { size: 'tiny', type: 'info', quaternary: true, onClick: () => edit(r) }, { default: () => '编辑' }),
        h(NButton, { size: 'tiny', type: 'warning', quaternary: true, onClick: () => diagnose(r) }, { default: () => '诊断' }),
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

function resetForm() {
  form.value = { name: '', type: 1, flow: 1, trafficRatio: 1, inIp: '', inNodeId: [], chainNodes: [[]], outNodeId: [] }
  chainMeta.value = [{ protocol: 'tls', strategy: 'fifo' }]
}

function openCreate() {
  editing.value = null
  resetForm()
  showModal.value = true
}

function edit(r: any) {
  editing.value = r
  form.value = {
    name: r.name, type: r.type, flow: r.flow, trafficRatio: r.trafficRatio, inIp: r.inIp || '',
    inNodeId: (r.inNodeId || []).map((x: any) => x.nodeId),
    chainNodes: (r.chainNodes && r.chainNodes.length ? r.chainNodes : [[]]).map((g: any) => g.map((x: any) => x.nodeId)),
    outNodeId: (r.outNodeId || []).map((x: any) => x.nodeId)
  }
  chainMeta.value = (r.chainNodes && r.chainNodes.length ? r.chainNodes : [[]]).map((g: any) => ({
    protocol: (g[0] && g[0].protocol) || 'tls',
    strategy: (g[0] && g[0].strategy) || 'fifo'
  }))
  showModal.value = true
}

async function save() {
  // 校验
  if (!form.value.name.trim()) { message.error('请输入隧道名称'); return }
  if (form.value.inNodeId.length === 0) { message.error('请至少选择一个入口节点'); return }
  let chainPayload: any[] = []
  if (form.value.type === 2) {
    if (form.value.outNodeId.length === 0) { message.error('隧道转发必须选择出口节点'); return }
    const validChains = form.value.chainNodes.filter((g) => g.length > 0)
    if (validChains.length === 0) { message.error('请至少配置一跳转发链节点'); return }
    form.value.chainNodes = validChains
    chainPayload = validChains.map((g, gi) => {
      const meta = chainMeta.value[gi] || { protocol: 'tls', strategy: 'fifo' }
      return g.map((nodeId, xi) => (xi === 0 ? { nodeId, protocol: meta.protocol, strategy: meta.strategy } : { nodeId }))
    })
  }
  try {
    if (editing.value) {
      await tunnelUpdate({ id: editing.value.id, name: form.value.name, flow: form.value.flow, trafficRatio: form.value.trafficRatio, inIp: form.value.inIp })
      message.success('隧道已更新(链路如需调整请删除重建)')
      showModal.value = false
      load()
      return
    }
    await tunnelCreate({
      name: form.value.name, type: form.value.type, flow: form.value.flow, trafficRatio: form.value.trafficRatio, inIp: form.value.inIp,
      inNodeId: form.value.inNodeId.map((nodeId) => ({ nodeId })),
      outNodeId: form.value.outNodeId.map((nodeId) => ({ nodeId })),
      chainNodes: chainPayload
    })
    message.success('隧道已创建,转发链服务已下发')
    showModal.value = false
    load()
  } catch (e: any) { message.error(e.message) }
}


async function remove(r: any) {
  try { await tunnelDelete(r.id); message.success('已删除(关联转发/授权/服务已清理)'); load() } catch (e: any) { message.error(e.message) }
}

async function diagnose(r: any) {
  diag.value = null
  showDiagModal.value = true
  diagLoading.value = true
  try {
    diag.value = await tunnelDiagnose(r.id)
    diag.value.tunnelName = r.name
  } catch (e: any) { message.error(e.message) } finally { diagLoading.value = false }
}

onMounted(load)
</script>

<template>
  <n-card title="隧道管理">
    <template #header-extra>
      <n-button type="primary" size="small" @click="openCreate">新建隧道</n-button>
    </template>
    <n-data-table :columns="columns" :data="rows" :bordered="false" size="small" :pagination="{ pageSize: 20 }" />

    <!-- 创建/编辑 -->
    <n-modal v-model:show="showModal" preset="card" :title="editing ? '编辑隧道' : '新建隧道'" style="width: 620px">
      <n-form label-placement="left" label-width="90px">
        <n-form-item label="名称"><n-input v-model:value="form.name" placeholder="2-50 个字符" /></n-form-item>
        <n-form-item v-if="!editing" label="类型">
          <n-radio-group v-model:value="form.type">
            <n-radio-button :value="1">端口转发</n-radio-button>
            <n-radio-button :value="2">隧道转发(多级链路)</n-radio-button>
          </n-radio-group>
        </n-form-item>
        <n-form-item label="入口节点">
          <n-select v-model:value="form.inNodeId" multiple :options="nodeOptions()" placeholder="选择在线节点" />
        </n-form-item>

        <!-- 转发链(仅隧道转发) -->
        <template v-if="form.type === 2 && !editing">
          <n-form-item label="转发链">
            <div style="width: 100%">
              <div v-for="(group, gi) in form.chainNodes" :key="gi" style="margin-bottom: 8px">
                <n-space align="center">
                  <n-tag size="small" :type="'info'">第{{ gi + 1 }}跳</n-tag>
                  <n-select v-model:value="form.chainNodes[gi]" multiple :options="nodeOptions()" placeholder="该跳的节点(可多节点)" style="flex: 1" />
                  <n-select v-model:value="chainMeta[gi].protocol" :options="chainProtocols" style="width: 110px" placeholder="协议" />
                  <n-select v-model:value="chainMeta[gi].strategy" :options="chainStrategies" style="width: 120px" placeholder="策略" />
                  <n-button size="tiny" quaternary type="error" :disabled="form.chainNodes.length <= 1" @click="form.chainNodes.splice(gi, 1); chainMeta.splice(gi, 1)">删</n-button>
                </n-space>
              </div>
              <n-button size="tiny" type="primary" ghost @click="form.chainNodes.push([]); chainMeta.push({ protocol: 'tls', strategy: 'fifo' })">+ 添加一跳</n-button>
            </div>
          </n-form-item>
          <n-form-item label="出口节点">
            <n-select v-model:value="form.outNodeId" multiple :options="nodeOptions()" placeholder="必填,选择在线节点" />
          </n-form-item>
        </template>

        <n-form-item label="流量类型">
          <n-radio-group v-model:value="form.flow">
            <n-radio-button :value="1">单向</n-radio-button>
            <n-radio-button :value="2">双向</n-radio-button>
          </n-radio-group>
        </n-form-item>
        <n-form-item label="流量倍率"><n-input-number v-model:value="form.trafficRatio" :min="0" :max="100" :step="0.1" style="width: 160px" /></n-form-item>
        <n-form-item label="入站 IP"><n-input v-model:value="form.inIp" placeholder="留空自动使用入口节点 IP" /></n-form-item>
      </n-form>
      <template #footer><n-button type="primary" block :loading="false" @click="save">{{ editing ? '保存' : '创建' }}</n-button></template>
    </n-modal>

    <!-- 诊断结果 -->
    <n-modal v-model:show="showDiagModal" preset="card" title="隧道诊断" style="width: 620px">
      <n-descriptions v-if="diag" :column="1" bordered size="small" style="margin-bottom: 12px">
        <n-descriptions-item label="隧道">{{ diag.tunnelName }}</n-descriptions-item>
        <n-descriptions-item label="类型">{{ diag.tunnelType }}</n-descriptions-item>
        <n-descriptions-item label="耗时">{{ diag.timestamp }}ms</n-descriptions-item>
      </n-descriptions>
      <n-data-table v-if="diag" :data="diag.results || []" :bordered="false" size="small" :columns="[
        { title: '节点', key: 'nodeName' },
        { title: '目标', key: 'target', render: (r: any) => `${r.targetIp}${r.targetPort ? ':' + r.targetPort : ''}` },
        { title: '结果', key: 'result', render: (r: any) => h(NTag, { type: r.success ? 'success' : 'error', size: 'small' }, { default: () => (r.success ? `✅ ${r.averageTime}ms · 丢包${r.packetLoss}%` : '❌ 失败') }) },
        { title: '说明', key: 'description' }
      ]" />
      <n-spin v-if="diagLoading" />
    </n-modal>
  </n-card>
</template>
