import { chromium } from 'playwright'
const BASE = process.argv[2] || 'http://127.0.0.1:3000'
const browser = await chromium.launch()
const ctx = await browser.newContext()
const page = await ctx.newPage()
const errors = []
page.on('pageerror', (e) => errors.push(`[pageerror] ${e.stack || e}`))
page.on('console', (m) => { if (m.type() === 'error') errors.push(`[console] ${m.text()}`) })
await page.goto(BASE, { waitUntil: 'networkidle' })
await page.click('button:has-text("登 录")')
await page.waitForTimeout(1500)
for (const p of ['/forward','/node','/tunnel','/user','/limit','/config','/profile','/dashboard']) {
  errors.length = 0
  await page.goto(`${BASE}${p}`, { waitUntil: 'networkidle' })
  await page.waitForTimeout(1500)
  console.log(`=== 直达 ${p}(${errors.length} 错误) ===`)
  console.log(errors[0] || '(no error)')
}
await browser.close()
