# 认证、配置与前端托管规格(Go 重写参考)

## 1. JWT 规格(自实现,无第三方库;Go 手写即可)

- **算法**:`HmacSHA256`(Go:`crypto/hmac` + `sha256`);Header 写 `{"alg":"HmacSHA256","typ":"JWT"}`(非标准 HS256 字符串,但验证只重算签名,不影响)
- **编码**:三段 `Base64URL(无padding)` 拼接 `header.payload.signature`;签名 = `HMAC-SHA256(secret, "header.payload")`
- **Claims**:
  - `sub`:userId 字符串
  - `iat`:秒;`exp`:秒,= now + **90 天**(代码 `90L*24*60*60*1000`;注释写 7 天是错的,以代码为准)
  - `user`、`name`:用户名;`role_id`:0=管理员,1=普通用户
- **密钥**:环境变量 `JWT_SECRET`,无默认值必须设置;**重写部署时必须沿用原值**,否则存量 token 全部失效
- **校验**:非空 → 3 段 → 重算签名字符串相等 → payload exp > now
- **无 Bearer 前缀**:`Authorization` 头直接裸 token
- 前端(`utils/jwt.ts`)只 `atob` 看 exp 不验签名(仅 UI 判断用,后端必须验)

## 2. 登录流程时序(Go 复刻)

```
POST /api/v1/captcha/check                 [免鉴权]
  vite_config.captcha_enabled=="true" → data=1,否则 data=0
        │ data=1 时前端弹 TAC 滑块验证码:
        │   POST /api/v1/captcha/generate → {id, 背景图数据...}
        │   POST /api/v1/captcha/verify {id, data:轨迹} → {validToken:<captchaId>}
        │   前端把 validToken 存 form.captchaId
        ▼
POST /api/v1/user/login {username, password, captchaId}   [免鉴权]
  ① captcha_enabled=true 时:captchaId 非空 + 二次校验(有效期 120s)
  ② SELECT * FROM "user" WHERE "user"=$1;不存在 → R.err("账号或密码错误")
  ③ md5(password) != pwd(**无盐** `Md5Util.md5`;`admin_salt_2024` 盐为死代码) → 同上
  ④ status==0 → R.err("账号被停用")
  ⑤ 签发 JWT(90 天)
  ⑥ requirePasswordChange = (username=="admin_user" || password=="admin_user")
  ⑦ → {code:0, data:{token, name, role_id, requirePasswordChange}}
```

- **改密**(`/user/updatePassword`,JWT):confirm 一致 → 当前密码比对 → 改用户名查重 → 同时更新用户名+新密码
- **注册**(`/user/register`,免鉴权):captcha 二次校验 → 用户名查重 → 默认配额(register_default_flow=100GB / num=5 / exp_days=30 / flow_reset_days=30,可被 vite_config 覆盖)→ 建 roleId=1 用户

## 3. 密码

- **无盐 MD5**:`md5(明文)` 32 位小写 hex(Md5Util.md5);源码另有 `md5WithSalt`/`DEFAULT_SALT="admin_salt_2024"` 但**无任何业务调用点**(死代码)
- 已实测:`printf 'admin_user' | md5` = `3c85cdebade1c51cf64ca9f3c09d182d` = data.sql 内置管理员哈希,证明登录链路无盐
- **建议**:Go 版新密码可选升级 bcrypt(`$2a$` 前缀识别),存量校验保留无盐 MD5 兼容路径

## 4. 验证码(tianai-captcha 1.5.2 → 自研滑块,见 05-frontend-ui.md §4)

- 原配置:vite_config `captcha_enabled`(开关)、`captcha_type`(SLIDER/WORD_IMAGE_CLICK/ROTATE/CONCAT/RANDOM)
- 原资源:26 张背景图(bgimages/)、10 组滑块模板(slide/)、字体 SIMSUN.TTC——均在 `resources/` 下,jar 内自带,**不依赖系统字体**
- 原流程:generate(返回含 id + 图片数据)→ verify(matching 轨迹,成功返回 `{validToken:id}`)→ login/register 二次校验(validToken 有效期 120s)
- **Go + Vue 方案**:三端点(`/api/v1/captcha/check|generate|verify`)保留,内部替换为自研滑块(Go 生成背景图+缺口坐标,Vue 组件拖拽校验,偏差阈值判定);`captcha_enabled/captcha_type` 配置语义保留;前后端同步重写无兼容负担;首版可先关闭 captcha 延后实现

## 5. CORS 与响应约定

- CORS:`allowedOrigins("*")`、方法 GET/POST/DELETE/PUT、暴露 `Authorization` 头(gin 中间件)
- 统一响应 `{code, msg, ts, data}`;401/403 放 body(HTTP 200)
- 例外:captcha generate/verify、open_api 成功、flow 三端点、/error

## 6. 配置与环境变量(对照 application.yml)

| 环境变量 | 默认 | 说明 |
|---|---|---|
| DB_HOST / DB_PORT | 127.0.0.1 / 5432 | PostgreSQL |
| DB_NAME / DB_USER / DB_PASSWORD | flvx2 | 连接串 `jdbc:postgresql://host:port/db` 的等价 |
| JWT_SECRET | 无(必填) | 沿用原值 |
| LOG_DIR | logs/ | 日志目录 |
| FRONTEND_PORT | 6635 | 前端静态页面入口(原 nginx 6366 → 6635),含 SPA fallback 与同源 /api 转发 |
| BACKEND_PORT | 6636 | API + WebSocket + flow 上报(原 6365 → 6636);**agent 对接端口** |
| 固定参数(可省) | — | HikariCP:max 20/min 5 → Go 连接池 2-5 即可;Tomcat 800 线程 → Go 默认 |

> 数据库初始化:启动时执行 schema.sql + data.sql(`sql.init.mode=always`,幂等 `ON CONFLICT DO NOTHING`),Go 用 embed + 事务执行,同样幂等。

## 7. 前端调用面(兼容性依据)

- axios baseURL:生产 `VITE_API_BASE` 为空 → 同源 `/api/v1/`(浏览器访问 6635,Go 在 6635 上同源处理 API);开发指向 `http://127.0.0.1:6636`
- token 放 `Authorization` 头;`code===401` 清 token 跳登录
- WebSocket:`ws://<host>/system-info?type=0&secret=${token}`(node 页);每 2s 收 info 推送,断线指数退避重连最多 5 次
- 前端页面路由:`/`(登录)、`/dashboard`、`/node`、`/tunnel`、`/forward`、`/limit`、`/user`、`/profile`、`/change-password`、`/config`、`/settings`

## 8. 静态托管(go:embed 方案,双端口)

- 构建产物:`vite build`(base='/'、outDir='dist'、默认 minify)→ 产物 ~1.5-2.5MB
- 结构:`dist/index.html` + `dist/assets/*` + `dist/favicon.*`(public/ 复制)
- **Go 进程双端口监听**:
  - **6635(前端)**:embed.FS 托管静态文件 + SPA fallback(排除 /api、/flow、/system-info 前缀);`/api/v1/*`、`/system-info`、`/flow/*` 由 Go 进程内转发到后端 handler(同引擎双端口,等价原 nginx 反代,但免 nginx 进程)
  - **6636(后端)**:原生挂载 API + WebSocket + flow(agent 对接端口);`GET /flow/test` → "test"
  - 两端口共享同一 gin 引擎与中间件(JWT 拦截、CORS),仅监听地址不同
- 备选:不嵌入,现有 nginx 方案原样保留(nginx :6635 反代 :6636),Go 仅替换后端端口 6636(前端零改动)

## 9. 前端配套(已定 Vue 3 整体重写,见 05-frontend-ui.md)

- 前端将整体重写为 Vue 3,`tac.min.js`、`@heroui`、`recharts` 均被替换;API/WS 契约不变
- 验证码随重写一并替换为自研滑块(§4),无独立兼容风险
- 若验证码首版未完成,`captcha_enabled` 保持默认 false 即可上线
