<script setup lang="ts">
import { h, onMounted, ref } from 'vue'
import { NCard, NGrid, NGridItem, NStatistic, NDataTable, NProgress, NTag } from 'naive-ui'
import * as echarts from 'echarts'
import { userPackage } from '@/api'

const pkg = ref<any>(null)
const stats = ref<any[]>([])
const chartEl = ref<HTMLDivElement | null>(null)
let chart: echarts.ECharts | null = null

const statColumns = [
  { title: '时间', key: 'time' },
  { title: '增量', key: 'flow', render: (r: any) => `${(r.flow / 1024 / 1024).toFixed(2)} MB` },
  { title: '累计', key: 'totalFlow', render: (r: any) => `${(r.totalFlow / 1024 / 1024 / 1024).toFixed(2)} GB` }
]

const forwardCols = [
  { title: '名称', key: 'name' },
  { title: '入口', key: 'entry', render: (r: any) => `${r.inIp}:${r.inPort}` },
  { title: '状态', key: 'status', render: (r: any) => h(NTag, { type: r.status === 1 ? 'success' : 'warning', size: 'small' }, { default: () => (r.status === 1 ? '运行中' : '已暂停') }) }
]

function renderChart() {
  if (!chartEl.value || !stats.value?.length) return
  if (!chart) chart = echarts.init(chartEl.value)
  const list = stats.value.slice(-24)
  chart.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: 40, right: 16, top: 20, bottom: 24 },
    xAxis: { type: 'category', data: list.map((x: any) => x.time) },
    yAxis: { type: 'value', name: 'MB' },
    series: [{
      name: '流量',
      type: 'bar',
      data: list.map((x: any) => +(x.flow / 1024 / 1024).toFixed(1)),
      itemStyle: { color: '#2080f0', borderRadius: [3, 3, 0, 0] },
      barMaxWidth: 24
    }]
  })
}

onMounted(async () => {
  try {
    pkg.value = await userPackage()
    stats.value = pkg.value.statisticsFlows || []
    renderChart()
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
      <div ref="chartEl" style="height: 240px"></div>
      <n-data-table :columns="statColumns" :data="stats.slice(-24)" :bordered="false" size="small" style="margin-top: 8px" />
    </n-card>
  </div>
</template>
