# FLVX2 → Go 重写方案(总览)

> 目标:保留 go-gost v3.2.6 定制内核(FLVX2-Agent)与 PostgreSQL 数据不动,把 Spring Boot 面板后端替换为 Go 单二进制,内存从 ~500M 降到 ~250M(面板进程 ~80M)。
> 依据:对 `springboot-backend`(112 个 Java 文件)与 `go-gost/panel`(agent 侧)逐行核对的双向协议规格。**agent 零改动。**

## 0. 命名与目录规范(项目级约定)

| 项 | 约定 | 说明 |
|---|---|---|
| 项目名 | **FLVX2** | 面板 Go 重写项目的正式名称 |
| Go 模块路径 | `github.com/Fyzzp0o0/FLVX2` | 建议单独仓库,或并入 FLVX2_Panel 的 `go/` 子目录 |
| 面板安装目录 | **`/opt/FLVX2`** | 安装脚本默认值(含二进制、前端 embed、logs) |
| Agent 安装目录 | **`/opt/FLVX2-Agent`** | 安装脚本默认值(现 go-gost 内核不变,仅目录名规范) |
| 面板二进制 | `FLVX2`(可执行文件名) | `ExecStart=/opt/FLVX2/FLVX2` |
| Agent 二进制 | `FLVX2-Agent`(现名 flux_agent,可保留原名) | 内核代码零改动,改名仅影响安装脚本 |
| systemd 服务名 | `flvx2-backend` / `flvx2-agent` | **保留小写**,与现有脚本一致,避免存量机器升级冲突 |
| 前端端口 | **6635** | 前端页面入口(原 nginx :6366 → **6635**);Go 进程内监听,含静态文件 + SPA fallback + 同源 /api 转发 |
| 后端端口 | **6636** | API/WebSocket/flow 上报(原 :6365 → **6636**);**agent 连 `面板IP:6636`** |
| 数据库 | flvx2(库/用户) | 不变,数据零迁移 |

> 安装脚本(新 `scripts/install-flvx2.sh` 或新仓库 install.sh)默认路径统一为 `APP_DIR=/opt/FLVX2`、`AGENT_DIR=/opt/FLVX2-Agent`,可通过环境变量覆盖。

---

## 1. 兼容性边界(重写不动的东西)

| 部分 | 处理方式 |
|---|---|
| go-gost 定制内核(flux_agent) | **零改动**。agent 通过 WebSocket + HTTP 与面板通信,协议全部开放,见 `02-protocol.md` |
| PostgreSQL 数据 | **零迁移**。表结构一字不动(直接执行原 `schema.sql`),存量数据直接可用 |
| 前端(React dist) | **由 Vue 3 重写**(已选定方案 B,见 `05-frontend-ui.md`),API/WS 契约不变,产物 `go:embed` 嵌入二进制(可选保留 nginx 方案) |
| 现有节点/隧道/转发运行状态 | 无感切换:停 Java 后端 → 起 Go 后端 → agent 自动重连 |

## 2. 技术选型

| 项 | 选择 | 说明 |
|---|---|---|
| 语言 | Go 1.22+ | 编译型原生二进制,无 VM 固定开销(对比 JVM `-Xms256m`) |
| HTTP 框架 | gin-gonic | 与 PFGo 同款,成熟轻量 |
| 数据库 | pgx v5(原生驱动) | 直接写 SQL,不用 ORM;`UPDATE ... SET in_flow = in_flow + $1` 原子累加 |
| WebSocket | gorilla/websocket | 与 agent 同库,行为对齐最稳 |
| JWT | 手写 HMAC-SHA256 | 兼容现有 token(用户无需重新登录,90 天有效期内无缝) |
| AES | crypto/aes + cipher.NewGCM | 兼容 agent 的 AES-256-GCM,密钥 = SHA-256(node secret) |
| 定时任务 | robfig/cron | 每天 00:00:05 流量重置+到期停服;每小时统计快照 |
| 静态托管 | go:embed `dist/` + **双端口监听** | 前端 **6635**(静态+SPA fallback,同源 `/api/*`、`/system-info`、`/flow/*` 由 Go 进程内转发到后端 handler);后端 **6636**(API/WS/flow,agent 对接)。免 nginx;需自维护 MIME 表 |
| 配置 | 环境变量 + 可选 .env | DB_HOST/DB_PORT/DB_NAME/DB_USER/DB_PASSWORD/JWT_SECRET/LOG_DIR |

## 3. 项目目录结构(建议)

```
FLVX2/
├── cmd/FLVX2/
│   └── main.go              # 入口:加载配置 → 初始化 DB → 启动 HTTP/WS → cron
├── internal/
│   ├── config/              # 环境变量解析(对照 application.yml)
│   ├── db/                  # pgx 连接池 + schema 初始化(embed 原 schema.sql/data.sql)
│   ├── model/               # 10 张表的结构体(int64 对齐 BIGINT)
│   ├── api/                 # REST 层(handlers,对应各 Controller)
│   │   ├── user.go          # /api/v1/user/*
│   │   ├── node.go          # /api/v1/node/*
│   │   ├── tunnel.go        # /api/v1/tunnel/*
│   │   ├── forward.go       # /api/v1/forward/*
│   │   ├── speed_limit.go   # /api/v1/speed-limit/*
│   │   ├── config.go        # /api/v1/config/*
│   │   ├── captcha.go       # /api/v1/captcha/*(或替换为自研滑块验证码)
│   │   ├── open_api.go      # /api/v1/open_api/sub_store
│   │   ├── flow.go          # /flow/upload /flow/config /flow/test
│   │   └── middleware.go    # JWT 拦截 + 角色校验 + CORS + R 包装
│   ├── ws/                  # /system-info WebSocket(节点 + 管理员双通道)
│   │   ├── hub.go           # 会话管理(节点覆盖旧连接、广播)
│   │   ├── crypto.go        # AES-256-GCM(SHA-256(secret) 派生)
│   │   └── commands.go      # 12 种命令构造与响应匹配(requestId → future)
│   ├── service/             # 业务逻辑(对照 ServiceImpl)
│   │   ├── forward.go       # 创建/暂停/恢复/删除 + 配额检查 + get_port 分配
│   │   ├── flow.go          # 流量累加 + 倍率 + 超限停服
│   │   ├── tunnel.go        # 链/服务下发(AddChains/AddChainService)
│   │   ├── node.go          # 节点 CRUD + SetProtocol
│   │   └── job.go           # cron:重置流量/到期停服/统计快照(孤儿清理由 /flow/config 触发)
│   └── web/                 # go:embed dist(前端静态资源,Vue 构建产物)
├── frontend/                # 前端源码(Vue 3 + Vite + TS,重写自 vite-frontend)
│   ├── src/                 # 页面/组件/状态(见 05-frontend-ui.md)
│   └── dist/                # 构建产物(拷贝至 internal/web/embed 或直接 embed frontend/dist)
├── web/dist/                # 前端构建产物(嵌入)
├── schema.sql / data.sql    # 原仓库资源,embed 使用
└── go.mod
```

## 4. 与原代码对应关系(迁移地图)

| Java 侧 | Go 侧 |
|---|---|
| `controller/*Controller` | `internal/api/*.go` |
| `service/impl/*ServiceImpl` | `internal/service/*.go` |
| `WebSocketServer` + `WebSocketInterceptor` | `internal/ws/hub.go` |
| `AESCrypto` / `JwtUtil` / `Md5Util` | `internal/ws/crypto.go` + `internal/api/middleware.go` |
| `GostUtil`(命令构造) | `internal/ws/commands.go` |
| `FlowController` | `internal/api/flow.go` + `internal/service/flow.go` |
| `ResetFlowAsync` / `StatisticsFlowAsync` / `CheckGostConfigAsync` | `internal/service/job.go`(cron) |
| `mapper/*.xml` + MyBatis-Plus 自动 CRUD | `internal/db/` 原生 SQL |

## 5. 实施阶段(里程碑)

| 阶段 | 内容 | 验收 |
|---|---|---|
| M1 骨架 | config/db/model/中间件/R 包装、`/flow/test`、登录+JWT | 启动 <1s,登录成功,旧 token 可用 |
| M2 节点通道 | `/system-info` WS(双类型握手)、AES 加解密、广播 | 现有 agent 上线,节点状态实时 |
| M3 流量闭环 | `/flow/upload` 累加+倍率+超限停服、`/flow/config` 孤儿清理 | 流量计数与 Java 版一致 |
| M4 业务 CRUD | 用户/节点/隧道/转发/限速/配置全部接口 | 与 Java 版行为逐项对比通过 |
| M5 命令下发 | AddService/Chains/Limiters/Pause/Resume/TcpPing/SetProtocol | 创建转发真实生效 |
| M6 前端与部署 | Vue 3 前端重写上线、go:embed dist、双端口 6635(前端)/6636(后端)、systemd 替换、数据/序列校验 | 浏览器全流程可用 |
| M7 定时任务+收尾 | cron 三件套、bug 修正清单(见 00 §7)、压测 | 跑 48h 无异常 |

## 6. 切换步骤(生产)

```bash
# 1) 备份(保险)
pg_dump -h 127.0.0.1 -U flvx2 flvx2 > flvx2-backup.sql

# 2) 停 Java 后端(agent 会进入 5s 重连等待)
systemctl stop flvx2-backend

# 3) 部署 Go 二进制(可复用同一数据库,零迁移)
install -m755 FLVX2 /opt/FLVX2/FLVX2
cat > /etc/systemd/system/flvx2-backend.service <<'EOF'
[Unit]
Description=FLVX2 Panel (Go)
After=network.target postgresql.service
[Service]
WorkingDirectory=/opt/FLVX2
Environment=DB_HOST=127.0.0.1 DB_PORT=5432 DB_NAME=flvx2 DB_USER=flvx2 DB_PASSWORD=<原密码>
Environment=JWT_SECRET=<原值,保证旧 token 可用> LOG_DIR=/opt/FLVX2/logs
Environment=FRONTEND_PORT=6635 BACKEND_PORT=6636
ExecStart=/opt/FLVX2/FLVX2
Restart=on-failure
[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload && systemctl enable --now flvx2-backend

# 4) 验证
curl -s http://127.0.0.1:6636/flow/test   # → "test"(后端/agent 对接端口)
curl -s http://127.0.0.1:6635/            # → 前端页面(免 nginx)
# agent config.json 的 addr 改为 面板IP:6636,5s 内自动重连,节点状态恢复在线
# 若重新执行 agent 安装脚本,设 NODE_PANEL_ADDR=<面板IP>:6636(install-flvx2.sh 默认值同步更新)
```

> ⚠️ `JWT_SECRET` 必须沿用原值,否则所有已登录用户的 token 失效;`DB_PASSWORD` 沿用原值。
> 若机器内存紧张,可同步执行 PostgreSQL 调优:`shared_buffers=32MB`、`max_connections=30`(整机可再降 ~100M)。

## 7. 重写时修正的已知 bug(来自源码核对)

1. `UserTunnelServiceImpl.updateUserTunnel` 成功路径返回 `R.err("用户隧道权限更新失败")` → 改为成功返回
2. `FlowController` 锁 Map(`locks`/`userLocks`/`tunnelLocks`)只增不减 → Go 用 `sync.Map` 或分片锁 + 定期清理
3. `WebSocketServer.sessionLocks` 同样只增不减 → 同上
4. JWT 注释"7 天"与实际 90 天不符 → 以 90 天为准,文档写清楚
5. 登录/注册日志记录明文密码(`LogAspect`) → 移除密码字段
6. `/api/v1/open_api/sub_store` 明文凭据走 GET → 至少提示改用 POST/加密(兼容性优先则保留并告警)
7. 密码无盐 MD5(`Md5Util.md5`;`admin_salt_2024` 盐是死代码)→ 新密码可平滑升级为 bcrypt(存量校验保留 MD5 兼容路径)
8. 端口正则允许 `00000` 等异常值 → 严格校验
9. `CheckGostConfigAsync` 孤儿清理按 `_tls` 结尾推断 tunnel_id 有误删风险 → 改为同时校验 forward/tunnel 存在性

## 8. 收益与成本

| 指标 | Java 现状 | Go 重写 |
|---|---|---|
| 面板进程常驻 | ~300-350M | **~80M** |
| 整机(含 PG + agent) | ~500M | **~250M**(PG 调优后 ~200M) |
| 启动时间 | 15s+ | **<1s** |
| 依赖服务 | JVM + PG + nginx | PG(可选 SQLite 化,需改 schema) |
| 二进制 | admin.jar(~100M+) | 单二进制 `FLVX2`(~40-60M,含前端) |
| 维护语言 | Java 21 + Spring | Go |
| 工作量 | — | 约 2-4 周(单人或 1 周全力) |
