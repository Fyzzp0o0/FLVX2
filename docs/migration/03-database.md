# 数据库规格(Spring Boot → Go 重写参考)

> 结论先行:**表结构一字不动**(直接执行原 `schema.sql`),存量数据零迁移;Go 用 pgx 写原生 SQL 替代 MyBatis-Plus。

## 0. 全局要点

- 10 张表,**全部无显式外键、无 CHECK、无触发器、无二级索引**(唯一索引仅 `vite_config.name`);id 均为 `BIGSERIAL`(隐式序列 `{table}_id_seq`);时间戳均为**毫秒级 BIGINT**
- 表名 `"user"`、列名 `"user"` 是 PG 保留字,**SQL 必须带双引号**
- **status 语义**:`1`=启用/正常,`0`=停用/暂停/删除(代码为准,注释不可信)
- 流量单位:`flow` 字段是 **GB**,字节换算 `1024³`;`in_flow/out_flow` 以字节计
- Go 类型映射:所有 BIGINT → `int64`;`traffic_ratio`(DOUBLE PRECISION)→ `float64`;`chain_type`(VARCHAR)→ 按字符串处理

## 1. 表结构清单

### 1.1 `forward`(转发条目,seq `forward_id_seq`)

| 列 | 类型 | 默认 | 说明 |
|---|---|---|---|
| id | BIGINT | BIGSERIAL | 主键 |
| user_id / user_name | BIGINT / VARCHAR(100) | — | 所属用户(冗余用户名) |
| name | VARCHAR(100) | — | 转发名(含 `_tcp`/`_udp` 后缀) |
| tunnel_id | BIGINT | — | 关联 tunnel |
| remote_addr | TEXT | — | 转发目标(可逗号多目标) |
| strategy | VARCHAR(100) | `'fifo'` | 负载策略 |
| in_flow / out_flow | BIGINT | 0 | 下行/上行字节(原子累加) |
| created_time / updated_time | BIGINT | — | 毫秒 |
| status | BIGINT | — | 1=启用,0=暂停 |
| inx | BIGINT | 0 | 排序 |

### 1.2 `forward_port`(转发占用端口)

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT | 主键 |
| forward_id / node_id | BIGINT | 关联 |
| port | BIGINT | 该转发在节点上占用的端口 |

### 1.3 `node`(节点)

| 列 | 类型 | 默认 | 说明 |
|---|---|---|---|
| id | BIGINT | — | 主键 |
| name / secret | VARCHAR(100) | — | secret=鉴权密钥(上报校验) |
| server_ip | VARCHAR(100) | — | 服务器 IP |
| port | TEXT | — | **端口范围串**,如 `"1000-2000,3000,4000-5000"`(逗号分隔,`a-b` 区间) |
| interface_name | VARCHAR(200) | NULL | 网卡名 |
| version | VARCHAR(100) | NULL | 客户端版本 |
| http / tls / socks | BIGINT | 0 | 协议支持标志 |
| created_time / updated_time | BIGINT | — | 毫秒 |
| status | BIGINT | — | 1=启用,0=停用 |
| tcp_listen_addr / udp_listen_addr | VARCHAR(100) | `'[::]'` | 监听地址 |

### 1.4 `speed_limit`(限速规则)

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT | 主键 |
| name / tunnel_name | VARCHAR(100) | 规则名/冗余隧道名 |
| speed | BIGINT | 限速值(bps,代码仅展示换算) |
| tunnel_id | BIGINT | 关联 tunnel |
| created_time / updated_time / status | BIGINT | — |

### 1.5 `statistics_flow`(统计快照)

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT | 主键 |
| user_id | BIGINT | 关联 |
| flow | BIGINT | 本小时增量(与上一条快照 total_flow 之差) |
| total_flow | BIGINT | 快照时刻累计 = in_flow+out_flow |
| time | VARCHAR(100) | `"HH:mm"` |
| created_time | BIGINT | 毫秒(清理 48h 前数据依据) |

### 1.6 `tunnel`(隧道模板)

| 列 | 类型 | 默认 | 说明 |
|---|---|---|---|
| id | BIGINT | — | 主键 |
| name | VARCHAR(100) | — | 隧道名 |
| traffic_ratio | DOUBLE PRECISION | 1.0 | 流量比率(上报乘数) |
| type | BIGINT | — | 1=端口转发,其他=隧道转发 |
| protocol | VARCHAR(10) | `'tls'` | 协议 |
| flow | BIGINT | — | 流量倍率系数(上报再乘;admin 视图硬编码 99999=不限制) |
| created_time / updated_time / status | BIGINT | — | — |
| in_ip | TEXT | NULL | 入站 IP |

### 1.7 `chain_tunnel`(隧道链节点)

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT | 主键 |
| tunnel_id / node_id | BIGINT | 关联 |
| chain_type | VARCHAR(10) | **字符串** `"1"`=入口、`"2"`=转发链、`"3"`=出口(按字符串比较!) |
| port | BIGINT | NULL=入口无端口配置 |
| strategy | VARCHAR(10) | 入口为 NULL |
| inx | BIGINT | 转发链序号,仅转发链有 |
| protocol | VARCHAR(10) | 入口为 NULL |

### 1.8 `"user"`(用户,表名/列名带引号)

| 列 | 类型 | 默认 | 说明 |
|---|---|---|---|
| id | BIGINT | — | 主键 |
| "user" | VARCHAR(100) | — | **用户名(保留字,必须 `"user"`)**;业务唯一 |
| pwd | VARCHAR(100) | — | `MD5(明文)` **无盐** 32 位小写(源码 `Md5Util.md5`;`md5WithSalt`/`admin_salt_2024` 为死代码,无业务调用) |
| role_id | BIGINT | — | 0=管理员(列表查询 `role_id != 0` 排除) |
| exp_time | BIGINT | — | 到期毫秒;admin 用 2727251700000(≈2056 年) |
| flow | BIGINT | — | 总流量上限 GB;超限 `flow×1024³ ≤ in+out` |
| in_flow / out_flow | BIGINT | 0 | 字节 |
| flow_reset_time | BIGINT | — | **每月第几号重置:0=不重置,1-31=每月 N 号**。⚠ 注册接口实际写入 `now + 30天` 毫秒时间戳(值恒 >31,月末分支恒命中 → 注册用户每月最后一天都会被重置流量,现状行为,Go 侧需决策保留或修正) |
| num | BIGINT | — | 转发数量上限 |
| created_time / updated_time / status | BIGINT | — | — |

### 1.9 `user_tunnel`(用户-隧道授权)

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT | 主键 |
| user_id / tunnel_id | BIGINT | 组合业务唯一 |
| speed_id | BIGINT | NULL 可空;**更新时必须显式写 NULL**(MyBatis-Plus IGNORED 语义) |
| num / flow | BIGINT | 该隧道内转发上限 / 流量上限 GB |
| in_flow / out_flow | BIGINT | 默认 0 |
| flow_reset_time / exp_time / status | BIGINT | 同 user |

### 1.10 `vite_config`(键值配置)

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT | 主键 |
| name | VARCHAR(200) | **UNIQUE**(upsert 依据) |
| value | VARCHAR(200) | 配置值(≤200 字符) |
| time | BIGINT | 毫秒 |

## 2. 初始数据(data.sql 等价)

- `INSERT INTO "user" VALUES (1,'admin_user','3c85cdebade1c51cf64ca9f3c09d182d',0,2727251700000,99999,0,0,1,99999,...) ON CONFLICT (id) DO NOTHING`(密码 = **无盐** `md5("admin_user")`,已实测验证;每月 1 号重置)
- `SELECT setval('user_id_seq', (SELECT MAX(id) FROM "user"))` — **必须保留**,否则后续插入主键冲突
- `INSERT INTO vite_config VALUES (1,'app_name','flux',...) ON CONFLICT (id) DO NOTHING` + setval 同理

## 3. 关键业务 SQL(Go 必须原样保留)

### 3.1 ⭐ 流量累加(每次上报,原子)

```sql
UPDATE forward     SET in_flow = in_flow + $d, out_flow = out_flow + $u WHERE id = $forwardId;
UPDATE "user"      SET in_flow = in_flow + $d, out_flow = out_flow + $u WHERE id = $userId;
UPDATE user_tunnel SET in_flow = in_flow + $d, out_flow = out_flow + $u WHERE id = $userTunnelId;
```
- 入库前先算倍率:`d' = trunc(d × traffic_ratio) × tunnel.flow`、`u'` 同理(源码先 BigDecimal 乘 ratio → `longValue()` 截断 → 再 ×flow;勿整体相乘后一次截断)
- 建议单事务 + 依赖 UPDATE 原子性(Java 还叠加了 synchronized 分锁,Go 可省)

### 3.2 ⭐ 超限/到期停服(userTunnelId≠0 时)

| 检查 | 条件 | 动作 |
|---|---|---|
| 用户流量 | `user.flow×1024³ < in+out`(**严格小于**,FlowController:265) | 暂停该用户全部 forward(status=0)+ `PauseService` |
| 用户到期 | `exp_time ≤ now` | 同上 |
| 用户停用 | `status != 1` | 同上 |
| 隧道流量 | `in+out ≥ user_tunnel.flow×1024³`(**含等号**,FlowController:292) | 暂停该用户该隧道的 forward |
| 隧道到期/停用 | user_tunnel 同条件 | 同上 |

恢复入口复查用 `flow×1024³ <= in+out`(含等号)判"已用完"。

### 3.3 ⭐ 流量重置(每天 00:00:05)

```sql
-- 当天非月末:
WHERE flow_reset_time != 0 AND flow_reset_time = $today
-- 当天是月末(30 天月设 31 号 → 30 号执行):
WHERE flow_reset_time != 0 AND (flow_reset_time = $today OR flow_reset_time > $lastDayOfMonth)
```
清零 `in_flow/out_flow`;**不碰 status**。

### 3.4 ⭐ 到期停服(同 cron 追加)

- `SELECT * FROM "user" WHERE role_id != 0 AND status = 1 AND exp_time < $now` → 暂停其 forward(status=0)→ `UPDATE "user" SET status=0`
- user_tunnel 同理;停服对每个 forward 查 `chain_tunnel WHERE tunnel_id=? AND chain_type='1'` 入口节点发 `PauseService`

### 3.5 ⭐ 端口分配 get_port

```sql
-- 已占用:
SELECT port FROM chain_tunnel WHERE node_id = ?;        -- port 非 NULL
SELECT port FROM forward_port WHERE node_id = ?;        -- 创建转发时排除自身:AND forward_id != ?
```
可用端口 = `parsePorts(node.port)` − 已占用,取升序第一个;parsePorts 按逗号拆、`a-b` 展开闭区间。

### 3.6 ⭐ 统计快照(每小时整点)

```sql
DELETE FROM statistics_flow WHERE created_time < $now-48h;
-- 每用户:上一条快照 ORDER BY id DESC LIMIT 1
INSERT INTO statistics_flow (user_id, flow, total_flow, time, created_time)
VALUES ($uid, $flow, $total, 'HH:mm', $now);
```
`total_flow = in+out`;`flow = 本次 − 上次`,**负则回退为本次 total_flow**;面板查询 `LIMIT 24`,不足补零。

### 3.7 其他

- 用户删除级联:`DELETE FROM forward/user_tunnel/statistics_flow/"user" WHERE user_id=?`
- 节点删除前:`SELECT ... FROM chain_tunnel WHERE node_id=? GROUP BY tunnel_id`(查引用)
- admin 隧道列表:`t.status=1`,流量/数量字段**查询时虚拟值** 99999/`'无限制'`(不落库)
- 三表联查(user_tunnel LEFT JOIN tunnel LEFT JOIN speed_limit)用于 package/user/list 详情

## 4. Go 实现注意清单

1. 直接执行原 schema.sql(PG 方言);不新增索引/外键
2. 迁移后 `setval('{table}_id_seq', MAX(id))` 对每表执行
3. 流量累加三连 UPDATE + 超限 + 停服放一个事务
4. `flow(GB) × 1024³` 比较字节,勿用 1e9;99999 GB × 1024³ ≈ 1.07e14 不溢出
5. `"user"` 表/列双引号(pgx:`SELECT ... FROM "user" WHERE "user" = $1`)
6. chain_type 按字符串 `'1'` 比较
7. 端口分配建议加 `SELECT ... FOR UPDATE` 或应用层锁防并发重复
8. cron 三件套:00:00:05 重置+停服;每小时整点快照+48h 清理;孤儿清理
