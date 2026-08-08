<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NCard, NButton, NForm, NFormItem, NInput, NGrid, NGridItem, NSpace, NAlert, useMessage } from 'naive-ui'
import { configList, configUpdateSingle } from '@/api'

const message = useMessage()
const cfg = ref<Record<string, string>>({})
const key = ref('')
const value = ref('')

// 常用配置(固定展示,新增配置会追加到下方"全部配置")
const commonKeys = ['app_name', 'ip', 'url', 'announcement']

const commonLabels: Record<string, string> = {
  app_name: '站点名称(app_name)',
  ip: '服务器IP(ip)',
  url: '站点地址(url)',
  announcement: '公告(announcement)'
}

const commonPlaceholders: Record<string, string> = {
  app_name: '显示在页面标题/登录页的品牌名',
  ip: '必填:节点安装命令中的面板地址(如 1.2.3.4 或 1.2.3.4:6636)',
  url: '面板对外访问地址',
  announcement: '登录页公告内容'
}

async function load() {
  try {
    cfg.value = await configList()
  } catch (e: any) { message.error(e.message) }
}

async function saveItem(name: string) {
  try {
    await configUpdateSingle(name, cfg.value[name] || '')
    message.success(`已保存: ${name}`)
    load()
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
    <n-alert type="info" style="margin-bottom: 12px" :show-icon="false">
      「服务器IP」用于节点安装命令:在节点管理点击「安装命令」前,请先在这里填写面板对外 IP(如 1.2.3.4 或 1.2.3.4:6636)并保存。
    </n-alert>

    <n-form label-placement="left" label-width="130px" style="max-width: 680px">
      <n-form-item v-for="k in commonKeys" :key="k" :label="commonLabels[k] || k">
        <n-input v-model:value="cfg[k]" :placeholder="commonPlaceholders[k] || ''" />
        <n-button size="small" type="primary" style="margin-left: 8px" @click="saveItem(k)">保存</n-button>
      </n-form-item>
    </n-form>

    <n-card title="全部配置(含自定义)" size="small" style="max-width: 680px; margin-top: 12px">
      <n-form label-placement="left" label-width="130px">
        <n-form-item v-for="(v, k) in cfg" :key="k" :label="k">
          <n-input v-model:value="cfg[k]" />
          <n-button size="small" type="primary" style="margin-left: 8px" @click="saveItem(k)">保存</n-button>
        </n-form-item>
      </n-form>
    </n-card>

    <n-card title="新增配置" size="small" style="max-width: 680px; margin-top: 12px">
      <n-grid :cols="3" :x-gap="8">
        <n-grid-item><n-input v-model:value="key" placeholder="配置名(如 ip)" /></n-grid-item>
        <n-grid-item><n-input v-model:value="value" placeholder="值" /></n-grid-item>
        <n-grid-item><n-button type="primary" block @click="addItem">添加</n-button></n-grid-item>
      </n-grid>
    </n-card>
  </n-card>
</template>
