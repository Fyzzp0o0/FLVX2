<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NCard, NGrid, NGridItem, NStatistic, NDataTable, NProgress, NTag } from 'naive-ui'
import { userPackage } from '@/api'

const pkg = ref<any>(null)
const stats = ref<any[]>([])

const statColumns = [
  { title: '时间', key: 'time' },
  { title: '增量', key: 'flow', render: (r: any) => `${(r.flow / 1024 / 1024).toFixed(2)} MB` },
  { title: '累计', key: 'totalFlow', render: (r: any) => `${(r.totalFlow / 1024 / 1024 / 1024).toFixed(2)} GB` }
]

const forwardCols = [
  { title: '名称', key: 'name' },
  { title: '入口', key: 'entry', render: (r: any) => `${r.inIp}:${r.inPort}` },
  { title: '状态', key: 'status', render: (r: any) => (r.status === 1 ? '运行中' : '已暂停') }
]

onMounted(async () => {
  try {
    pkg.value = await userPackage()
    stats.value = pkg.value.statisticsFlows || []
  } catch (e: any) { /* 路由守卫兜底 */ }
})

function pct(flow: number) {
  const limit = pkg.value?.userInfo?.flow || 1
  return Math.min(100, Math.round(((pkg.value?.userInfo?.inFlow + pkg.value?.userInfo?.outFlow) / (limit * 1024 ** 3)) * 100))
}
</script>

<template>
  <div v-if="pkg">
    <n-grid :cols="4" :x-gap="12" responsive="screen">
      <n-grid-item>
        <n-card><n-statistic label="账户" :value="pkg.userInfo.user" /></n-card>
      </n-grid-item>
      <n-grid-item>
        <n-card><n-statistic label="流量使用" :value="pct(0)" suffix="%"><template #prefix /></n-statistic>
          <n-progress type="line" :percentage="pct(0)" :show-indicator="false" /></n-card>
      </n-grid-item>
      <n-grid-item>
        <n-card><n-statistic label="到期时间" :value="new Date(pkg.userInfo.expTime).toLocaleDateString()" /></n-card>
      </n-grid-item>
      <n-grid-item>
        <n-card><n-statistic label="转发数量" :value="pkg.forwards.length" :suffix="`/ ${pkg.userInfo.num}`" /></n-card>
      </n-grid-item>
    </n-grid>
    <n-card title="我的转发" style="margin-top: 12px">
      <n-data-table :columns="forwardCols" :data="pkg.forwards" :bordered="false" size="small" />
    </n-card>
    <n-card title="近 24 小时流量" style="margin-top: 12px">
      <n-data-table :columns="statColumns" :data="stats.slice(-24)" :bordered="false" size="small" />
    </n-card>
  </div>
</template>
