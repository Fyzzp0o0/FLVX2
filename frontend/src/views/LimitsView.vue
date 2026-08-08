<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NCard, NButton, NDataTable, NModal, NForm, NFormItem, NInput, NInputNumber, NSelect, useMessage } from 'naive-ui'
import { speedLimitList, speedLimitCreate, speedLimitDelete, tunnelList } from '@/api'

const message = useMessage()
const rows = ref<any[]>([])
const tunnels = ref<any[]>([])
const showModal = ref(false)
const form = ref({ name: '', speed: 1048576, tunnelId: null as number | null, tunnelName: '' })

const columns = [
  { title: '名称', key: 'name' },
  { title: '限速', key: 'speedView', render: (r: any) => `${(r.speed / 1048576).toFixed(1)} Mbps` },
  { title: '隧道', key: 'tunnelName' },
  {
    title: '操作', key: 'actions', render: (r: any) =>
      h(NButton, { size: 'tiny', type: 'error', quaternary: true, onClick: () => remove(r) }, { default: () => '删除' })
  }
]

async function load() {
  try {
    rows.value = await speedLimitList()
    tunnels.value = await tunnelList()
  } catch (e: any) { message.error(e.message) }
}

async function save() {
  const t = tunnels.value.find((x) => x.id === form.value.tunnelId)
  try {
    await speedLimitCreate({ ...form.value, tunnelName: t?.name || '' })
    message.success('限速规则已创建')
    showModal.value = false
    load()
  } catch (e: any) { message.error(e.message) }
}

async function remove(r: any) {
  try { await speedLimitDelete(r.id); message.success('已删除'); load() } catch (e: any) { message.error(e.message) }
}

onMounted(load)
</script>

<template>
  <n-card title="限速规则">
    <template #header-extra>
      <n-button type="primary" size="small" @click="showModal = true">新建规则</n-button>
    </template>
    <n-data-table :columns="columns" :data="rows" :bordered="false" size="small" />

    <n-modal v-model:show="showModal" preset="card" title="新建限速规则" style="width: 460px">
      <n-form label-placement="left" label-width="90px">
        <n-form-item label="名称"><n-input v-model:value="form.name" /></n-form-item>
        <n-form-item label="限速(bps)"><n-input-number v-model:value="form.speed" :min="1" /></n-form-item>
        <n-form-item label="绑定隧道">
          <n-select v-model:value="form.tunnelId" :options="tunnels.map((t) => ({ label: t.name, value: t.id }))" />
        </n-form-item>
      </n-form>
      <template #footer><n-button type="primary" block @click="save">创建</n-button></template>
    </n-modal>
  </n-card>
</template>
