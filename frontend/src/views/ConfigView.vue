<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NCard, NButton, NForm, NFormItem, NInput, NGrid, NGridItem, useMessage } from 'naive-ui'
import { configList, configUpdateSingle } from '@/api'

const message = useMessage()
const cfg = ref<Record<string, string>>({})
const key = ref('')
const value = ref('')

async function load() {
  try { cfg.value = await configList() } catch (e: any) { message.error(e.message) }
}

async function saveItem(name: string) {
  try {
    await configUpdateSingle(name, cfg.value[name] || '')
    message.success('已保存: ' + name)
  } catch (e: any) { message.error(e.message) }
}

async function addItem() {
  if (!key.value) { message.error('请输入配置名'); return }
  try {
    await configUpdateSingle(key.value, value.value)
    message.success('已添加')
    key.value = ''
    value.value = ''
    load()
  } catch (e: any) { message.error(e.message) }
}

onMounted(load)
</script>

<template>
  <n-card title="网站配置">
    <n-form label-placement="left" label-width="120px" style="max-width: 640px">
      <n-form-item v-for="(v, k) in cfg" :key="k" :label="k">
        <n-input v-model:value="cfg[k]" />
        <n-button size="small" type="primary" style="margin-left: 8px" @click="saveItem(k)">保存</n-button>
      </n-form-item>
    </n-form>
    <n-card title="新增配置" size="small" style="max-width: 640px; margin-top: 12px">
      <n-grid :cols="3" :x-gap="8">
        <n-grid-item><n-input v-model:value="key" placeholder="配置名(如 ip)" /></n-grid-item>
        <n-grid-item><n-input v-model:value="value" placeholder="值" /></n-grid-item>
        <n-grid-item><n-button type="primary" block @click="addItem">添加</n-button></n-grid-item>
      </n-grid>
    </n-card>
  </n-card>
</template>
