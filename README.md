# FLVX2

Go 重写版流量端口转发 / 隧道转发面板。保留 go-gost v3.2.6 定制内核(FLVX2-Agent)与 PostgreSQL 数据,后端由 Spring Boot 替换为 **Go 单二进制**(含 Vue 3 前端),内存占用从 ~500M 降至 ~200-250M。

## 特性

- 🚀 **单二进制**:前端(Vue 3 + Naive UI)经 `go:embed` 嵌入,免 nginx,双端口 6635(前端)/ 6636(后端)
- 🧩 **协议兼容**:与 go-gost 定制内核(agent)的 WebSocket + AES-256-GCM 协议逐字段对齐,**agent 零改动**;旧 JWT 存量 token 无缝迁移
- 📊 多用户 / 多节点 / 端口转发 / 隧道转发 / 限速 / 流量统计(单向/双向)/ 订阅(subscription-userinfo)
- ⚙️ 定时任务:每日流量重置 + 到期停服、每小时统计快照
- 🔒 兼容原架构:PostgreSQL 表结构一字不动,数据零迁移

## 快速开始

```bash
# 面板
curl -fsSL https://raw.githubusercontent.com/Fyzzp0o0/FLVX2/main/scripts/install-flvx2.sh | bash -s -- install
# 节点内核(需先创建节点获得 NODE_SECRET)
NODE_SECRET=xxxx curl -fsSL https://raw.githubusercontent.com/Fyzzp0o0/FLVX2/main/scripts/install-flvx2.sh | bash -s -- install-agent
```

详细文档见 [docs/DEPLOY.md](docs/DEPLOY.md)。

## 开发

```bash
# 后端(Go 1.22+)
DB_PASSWORD=xxx JWT_SECRET=xxx go run ./cmd/FLVX2

# 前端(Vue 3)
cd frontend && npm install && npm run dev   # 开发(API 指向 6636)
npm run build                                # 产物输出至 ../internal/web/dist(嵌入二进制)
```

## 目录结构

```
cmd/FLVX2/          入口(双端口监听 + cron)
internal/
  api/              REST 层(登录/节点/隧道/转发/限速/配置/open_api)
  ws/               WebSocket(AES-256-GCM)+ Gost 命令下发
  service/          业务逻辑(含定时任务/流量闭环/孤儿清理/诊断)
  db/               连接池 + schema/data 幂等初始化
  web/              go:embed 前端产物 + SPA fallback
frontend/           Vue 3 + Naive UI 前端
scripts/            一键脚本(install / install-agent / update / uninstall)
docs/               部署手册与迁移方案
```

## License

AGPL-3.0
