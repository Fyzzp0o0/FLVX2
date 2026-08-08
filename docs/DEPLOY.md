# FLVX2 部署与维护手册

> 架构:Go 单二进制面板(前端 Vue3 已嵌入)+ PostgreSQL + go-gost 定制内核(FLVX2-Agent)。
> 端口:**6635** 前端页面 / **6636** 后端 API 与 agent 对接(无需 nginx)。

## 1. 快速安装(一键脚本)

```bash
# 安装面板(自动装 PostgreSQL、下载二进制、注册 systemd 服务)
curl -fsSL https://raw.githubusercontent.com/Fyzzp0o0/FLVX2/main/scripts/install-flvx2.sh | bash -s -- install

# 安装节点内核(agent):NODE_SECRET 在面板「节点管理 → 创建节点」后获得
curl -fsSL https://raw.githubusercontent.com/Fyzzp0o0/FLVX2/main/scripts/install-flvx2.sh | bash -s -- install-agent
```

安装完成后:

- 前端:访问 `http://<服务器IP>:6635`,默认账号 **admin_user / admin_user**(登录后请立即修改)
- 节点:面板「节点管理」创建节点 → 获得 secret → 在节点机执行 `install-agent`(设 `NODE_SECRET=xxx`)

### 常用环境变量

| 变量 | 默认 | 说明 |
|---|---|---|
| `FRONTEND_PORT` | 6635 | 前端端口 |
| `BACKEND_PORT` | 6636 | 后端端口(agent 对接) |
| `DB_NAME` / `DB_USER` / `DB_PASSWORD` | flvx2 / flvx2 / 随机 | PostgreSQL(首次生成后保存在 `/opt/FLVX2/.env`,重跑脚本会复用) |
| `JWT_SECRET` | 随机 | JWT 密钥(首次生成后固定;**更新/重装必须保持,否则存量 token 全部失效**) |
| `FLVX2_VERSION` | latest | 指定安装版本(如 `2.0.0`) |
| `NODE_SECRET` | — | agent 安装必填(节点密钥) |
| `NODE_PANEL_ADDR` | `本机IP:6636` | agent 连接的面板地址 |

## 2. 更新

```bash
curl -fsSL https://raw.githubusercontent.com/Fyzzp0o0/FLVX2/main/scripts/install-flvx2.sh | bash -s -- update
```

更新仅替换二进制并重启,配置(`/opt/FLVX2/.env`)与数据库保留。

## 3. 完全卸载

```bash
curl -fsSL https://raw.githubusercontent.com/Fyzzp0o0/FLVX2/main/scripts/install-flvx2.sh | bash -s -- uninstall
```

删除服务、`/opt/FLVX2`、`/opt/FLVX2-Agent`;**数据库保留**。如需连数据库一并删除:

```bash
sudo -u postgres dropdb flvx2
sudo -u postgres psql -c "DROP USER IF EXISTS \"flvx2\""
```

## 4. 日常维护命令

```bash
# 服务状态 / 日志
systemctl status flvx2-backend
journalctl -u flvx2-backend -f          # 实时日志
journalctl -u flvx2-backend --since "10 min ago"

# 重启 / 停止 / 启动
systemctl restart flvx2-backend
systemctl stop flvx2-backend
systemctl start flvx2-backend

# Agent 侧
systemctl status flvx2-agent
journalctl -u flvx2-agent -f
```

### 数据库备份 / 恢复

```bash
# 备份
pg_dump -h 127.0.0.1 -U flvx2 flvx2 > flvx2-backup-$(date +%F).sql

# 恢复
psql -h 127.0.0.1 -U flvx2 -d flvx2 < flvx2-backup-2026-01-01.sql
```

### 手动运行(调试)

```bash
cd /opt/FLVX2 && set -a && . ./.env && set +a && ./FLVX2
```

## 5. 配置说明(`/opt/FLVX2/.env`)

| 键 | 说明 |
|---|---|
| `DB_HOST/DB_PORT/DB_NAME/DB_USER/DB_PASSWORD` | PostgreSQL 连接 |
| `JWT_SECRET` | JWT 签名密钥(**勿改**,改了全部用户要重新登录) |
| `LOG_DIR` | 日志目录 |
| `FRONTEND_PORT` / `BACKEND_PORT` | 端口 |

## 6. 定时任务(内置于面板,无需配置)

| 任务 | 时间 | 行为 |
|---|---|---|
| 流量重置 + 到期停服 | 每天 00:00:05 | 按 `flow_reset_time` 重置流量;到期用户/隧道暂停转发并停用 |
| 流量统计快照 | 每小时整点 | 写 `statistics_flow`(保留 48h) |

## 7. 目录结构

```
/opt/FLVX2/
├── FLVX2           # 面板二进制(含前端)
├── .env            # 配置(权限 600)
└── logs/           # 日志
/opt/FLVX2-Agent/
├── FLVX2-Agent     # 节点内核二进制
└── config.json     # addr/secret/http/tls/socks
```

## 8. 常见问题

- **节点一直离线**:检查 `config.json` 的 `addr` 是否为 `面板IP:6636`、`secret` 是否与面板节点一致;`journalctl -u flvx2-agent -f` 查看连接日志。
- **登录后立即被登出**:`JWT_SECRET` 与安装时不一致(更新/迁移时务必保留 `.env`)。
- **流量不统计**:确认转发服务名格式(`{转发ID}_{用户ID}_{隧道权限ID}_tcp`),agent 每 5s 上报一次。
- **端口被占用**:修改 `.env` 的 `FRONTEND_PORT/BACKEND_PORT` 后 `systemctl restart flvx2-backend`。
