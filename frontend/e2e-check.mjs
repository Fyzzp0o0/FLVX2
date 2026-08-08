// 前端端到端冒烟:登录 → 逐菜单点击 → 验证页面变化与 console 错误
// 用法: node scripts/e2e-check.mjs <面板地址>
import { chromium } from 'playwright'

const BASE = process.argv[2] || 'http://127.0.0.1:6635'
const results = []
const errors = []

const browser = await chromium.launch()
const page = await browser.newPage()
page.on('console', (msg) => {
  if (msg.type() === 'error') errors.push(`[console.error] ${msg.text()}`)
})
page.on('pageerror', (err) => errors.push(`[pageerror] ${err.stack || err}`))

// 1. 登录页
await page.goto(BASE, { waitUntil: 'networkidle' })
await page.waitForTimeout(1500)
results.push(['登录页标题', await page.title()])
const loginVisible = await page.isVisible('button:has-text("登 录")')
results.push(['登录按钮可见', loginVisible])

// 2. 登录
await page.click('button:has-text("登 录")')
await page.waitForTimeout(2000)
results.push(['登录后 URL', page.url()])

// 3. 逐菜单点击,验证页面切换
const menus = [
  ['仪表盘', '/dashboard'],
  ['我的转发', '/forward'],
  ['节点管理', '/node'],
  ['隧道管理', '/tunnel'],
  ['用户管理', '/user'],
  ['限速规则', '/limit'],
  ['网站配置', '/config'],
  ['个人中心', '/profile']
]
for (const [label, path] of menus) {
  const before = errors.length
  const menuVisible = await page.isVisible(`.n-menu-item-content:has-text("${label}")`)
  if (menuVisible) {
    await page.click(`.n-menu-item-content:has-text("${label}")`)
    await page.waitForTimeout(1500)
    const url = page.url()
    const body = await page.evaluate(() => document.body.innerText.length)
    const newErrs = errors.slice(before)
    results.push([`点击「${label}」`, `${url} | 正文${body}字${url.includes(path) ? '' : ' ❌路径不符'}${newErrs.length ? ` | 新增错误${newErrs.length}条: ${newErrs[0]}` : ''}`])
  } else {
    results.push([`点击「${label}」`, '❌ 菜单项不存在'])
  }
}

// 4. 截图
await page.screenshot({ path: '/tmp/e2e-final.png', fullPage: true })

console.log('===== 页面流程结果 =====')
for (const [k, v] of results) console.log(`  ${k}: ${v}`)
console.log('===== 控制台/页面错误 =====')
if (errors.length === 0) console.log('  无错误 ✅')
else for (const e of errors) console.log(`  ${e}`)

await browser.close()
