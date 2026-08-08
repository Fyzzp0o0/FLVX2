import { ref } from 'vue'
import { defineStore } from 'pinia'
import { configGet } from '@/api'

// 站点全局配置(app_name 等),免鉴权读取
export const useAppStore = defineStore('app', () => {
  const appName = ref('FLVX2')
  const announcement = ref('')
  let loaded = false

  async function load() {
    if (loaded) return
    loaded = true
    try {
      const [name, ann] = await Promise.all([
        configGet('app_name').catch(() => null),
        configGet('announcement').catch(() => null)
      ])
      if (name?.value) {
        appName.value = name.value
        document.title = `${name.value} Panel`
      }
      if (ann?.value) announcement.value = ann.value
    } catch { /* 免鉴权失败时保持默认 */ }
  }

  return { appName, announcement, load }
})
