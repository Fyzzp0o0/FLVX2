# WebSocket 与流量上报协议规格(agent 零改动的依据)

> 依据 `go-gost/panel/`(agent 侧)与 `springboot-backend`(面板侧)逐行双向核对。
> **Go 重写面板必须按本规格逐字段对齐;agent 保持不动。**

## 1. 连接建立

### 1.1 URL 与 Query(agent → 面板)

```
ws://<host>:<port>/system-info?type=1&secret=<节点密钥>&version=<version>&http=<0|1>&tls=<0|1>&socks=<0|1>
```

| 参数 | 值 | 说明 |
|---|---|---|
| `type` | `1`=节点(secret 查 node 表);`0`=管理员(secret 即 JWT) | 决定握手校验分支 |
| `secret` | 节点密钥 / JWT | type=1 查不到 → 握手拒绝(HTTP 400) |
| `version` | agent 传空串 | 非空时握手后写 node.version |
| `http/tls/socks` | 0/1 | 节点启动/重连时**每次重读** config.json;握手后写 node 表 |

agent 侧 `config.json`:`{"addr":"1.2.3.4:6636","secret":"xxxx","http":1,"tls":0,"socks":1}`。`addr` 同时用于 WS 与 HTTP 上报;**必须指向面板后端端口 6636**(原 6365 → 6636;安装脚本 `NODE_PANEL_ADDR` 默认同步更新)。

### 1.2 会话管理(面板侧)

- type=1:nodeId → session;**同节点新连接覆盖旧连接并主动 close 旧连接**;更新 DB status=1;向所有管理员会话广播 `{"id":"<nodeId>","type":"status","data":1}`
- 连接关闭:仅当前活跃 session 被关时更新 DB status=0 并广播 `data:0`(旧连接被覆盖后关闭不触发)
- 心跳:agent 每 2s 发系统信息;面板收到**解密后 payload 含子串 `memory_usage`** 即回 `{"type":"call"}`(加密发送);agent 忽略 type=="call"
- agent 重连:失败 5s 重试;读超时 30s;写超时 5s(>1MB 消息 30s);握手超时 10s
- 面板向离线节点 send_msg → 立即返回 `{msg:"节点不在线"}`

## 2. 消息框架

### 2.1 上行(agent → 面板)

**系统信息(每 2s)**:
```json
{"uptime":12345,"bytes_received":1000,"bytes_transmitted":2000,"cpu_usage":12.3,"memory_usage":45.6}
```
面板解密后广播:`{"id":"<nodeId>","type":"info","data":"<解密后payload字符串>"}`(data 是字符串,前端 JSON.parse)。

**命令响应**:
```json
{"type":"AddServiceResponse","success":true,"message":"OK","data":{...},"requestId":"<原样回显>"}
```
- `success` 固定存在;`message` 成功固定 `"OK"`
- `data` 仅 `TcpPingResponse` 携带
- `requestId` **必须原样回显**,面板据此匹配 pending future
- 错误响应:`DecryptError`/`NoDecryptor`/`DecompressError`/`ParseError`/`UnknownCommandResponse`

### 2.2 下行(面板 → agent)

```json
{"type":"AddService","data":{...},"requestId":"<UUID>"}
```
- `requestId` 必填(UUID);面板 `future.get(10s)` 超时 → `{msg:"等待响应超时"}`
- 面板对所有节点命令加密发送

### 2.3 AES-256-GCM 加密(双向对称)

- **密钥**:`key = SHA-256(secret)` 取 32 字节;secret 为空则双方明文
- **加密**:nonce=12 随机字节;`ciphertext = AES-256-GCM-Seal(nonce, plaintext, aad=nil)`;载荷 = `nonce‖ciphertext‖tag`;`data = Base64.StdEncoding(载荷)`
- **包装**:
```json
{"encrypted":true,"data":"<base64>","timestamp":1750000000}
```
- 解密判定:`encrypted==true` 且 data 非空 → 解密;否则明文。`timestamp` 双方都不校验(秒/毫秒均可,Go 存 json.Number)
- 解密失败:agent 回 DecryptError;Java 面板解密失败按明文继续(Go 同)

### 2.4 gzip 压缩(agent 单方面支持,面板从未发送)

```json
{"type":"AddService","compressed":true,"data":"<gzip压缩后的JSON字节>","requestId":"..."}
```
处理顺序:先解密(若 encrypted)→ 解析出 `compressed==true` → gzip 解压 → 解压结果即命令 data → 路由。Go 面板建议保留接收容错,不必发送。

### 2.5 duration 字段

- `forwarder.selector.failTimeout`:面板发字符串 `"600s"`(Java 现状),agent 转纳秒
- `hop.selector.failTimeout`:Java 发数字纳秒 `600000000000`
- Go 面板**两种都接受**(agent 预处理兼容)

## 3. 命令全集

所有命令处理完 agent 都会写 `gost.json` + SIGHUP 热重载;处理前恢复 `panel-observer` 注册。

| 下行 type | 响应 type | data 结构(面板发送) | 用途 |
|---|---|---|---|
| `AddService` | `AddServiceResponse` | **数组** `[ServiceConfig,...]` | 批量创建服务;重名 → "service X already exists";自动注入 observer/`enableStats=true`/`observePeriod=5s` |
| `UpdateService` | `UpdateServiceResponse` | **数组** `[ServiceConfig,...]` | 按 name 替换;不存在 → not found |
| `DeleteService` | `DeleteServiceResponse` | `{"services":["n1","n2"]}` | 按名移除 |
| `PauseService` | `PauseServiceResponse` | `{"services":["n_tcp","n_udp"]}` | 摘除并存入 paused_services.json;全不匹配 → "no matching services found" |
| `ResumeService` | `ResumeServiceResponse` | `{"services":[...]}` | 从 paused_services.json 恢复 |
| `AddChains` | `AddChainsResponse` | `ChainConfig` 对象 | 创建链;重名 → exists |
| `UpdateChains` | `UpdateChainsResponse` | `{"chain":"<name>","data":ChainConfig}` **或**直接 ChainConfig | 按 name 替换 |
| `DeleteChains` | `DeleteChainsResponse` | `{"chain":"<name>"}` **或**字符串 `"<name>"` | 删链 |
| `AddLimiters` | `AddLimitersResponse` | `LimiterConfig` | 创建限速器;重名 → exists |
| `UpdateLimiters` | `UpdateLimitersResponse` | `{"limiter":"<name>","data":LimiterConfig}` **或**直接 LimiterConfig | 替换 |
| `DeleteLimiters` | `DeleteLimitersResponse` | `{"limiter":"<name>"}` **或**字符串 | 删除 |
| `TcpPing` | `TcpPingResponse` | `{"ip":"...","port":8080,"count":4,"timeout":5000}` | TCP 诊断;响应 data 见下 |
| `SetProtocol` | `SetProtocolResponse` | `{"http":0|1,"tls":0|1,"socks":0|1}`(**必须全量;未提供的字段 agent 视为 0 覆盖写回 config.json**,勿只传变化字段) | 协议开关,写回 agent config.json |

**Java 容错规则(必须复刻,幂等)**:
- Add 类命令响应 message 含 `"exists"` → 视为成功
- Delete 类命令响应 message 含 `"not found"` → 视为成功

### 3.1 ChainConfig(AddChains 的 data)

```json
{
  "name": "chains_<tunnelId>",
  "hops": [{
    "name": "hop_<tunnelId>",
    "interface": "<出接口名,仅配置了interfaceName时>",
    "selector": {"strategy": "<strategy>", "maxFails": 1, "failTimeout": 600000000000},
    "nodes": [{
      "name": "node_<inx>",
      "addr": "<ip:port,IPv6加方括号>",
      "connector": {"type": "relay"},
      "dialer": {"type": "<protocol>"}
    }]
  }]
}
```

### 3.2 LimiterConfig(AddLimiters 的 data)

```json
{"name": "<speedLimit.id字符串>", "limits": ["$ <X>MB <X>MB"]}
```

### 3.3 TcpPingResponse(agent 响应 data)

```json
{"ip":"1.2.3.4","port":8080,"success":true,"averageTime":12.34,"packetLoss":0.0,"errorMessage":"...","requestId":"..."}
```
- 面板判 `msg=="OK"` 且解析 data.success/averageTime/packetLoss
- ⚠ Java 侧另有 `"PingResponse"` 特判为死分支(agent 实际响应类型是 `TcpPingResponse`,永不命中),Go 勿复刻
- agent 对 count<=0 默认 4、timeout<=0 默认 5000ms

### 3.4 ServiceConfig 关键字段

**端口转发服务(tcp/udp)**:
```json
{
  "name": "<forwardId>_<userId>_<userTunnelId>_tcp",
  "addr": "<node.tcpListenAddr>:<port>",
  "metadata": {"interface": "<网卡名,仅type=1且有网卡>"},
  "limiter": "<speedLimit.id>",
  "handler": {"type": "tcp", "chain": "chains_<tunnelId>"},
  "listener": {"type": "tcp"},
  "forwarder": {
    "nodes": [{"name": "node_1", "addr": "<remoteAddr>"}, ...],
    "selector": {"strategy": "fifo", "maxFails": 1, "failTimeout": "600s"}
  }
}
```
- 每个转发每台入口节点生成 **tcp、udp 两个服务**;udp listener 附 `"metadata":{"keepAlive":true}`
- `handler.chain` 仅隧道转发(type=2)

**隧道 TLS 服务**(`<tunnelId>_tls`,AddChainService):
```json
{
  "name": "<tunnelId>_tls",
  "addr": "<node.tcpListenAddr>:<port>",
  "handler": {"type": "relay", "chain": "chains_<tunnelId>"},
  "listener": {"type": "<protocol>"}
}
```
- 仅出口节点(chainType==3)设 metadata.interface;仅中转节点(chainType==2)handler 带 chain

## 4. 命名规范(面板逐字生成,agent 按名匹配)

| 对象 | 格式 | 示例 |
|---|---|---|
| 转发服务基础名 | `{forwardId}_{userId}_{userTunnelId}` | `12_3_4`(无权限时第三段 0) |
| 转发服务全名 | 基础名 + `_tcp` / `_udp` | `12_3_4_tcp` |
| 隧道 TLS 服务 | `{tunnelId}_tls` | `5_tls` |
| 链名 | `chains_{tunnelId}` | `chains_5` |
| hop 名 | `hop_{tunnelId}` | `hop_5` |
| 链节点名 | `node_{inx}` | `node_1` |
| 转发节点名 | `node_{num}`(1 起,按 remoteAddr 逗号分隔) | `node_1` |
| 限速器名 | `{speedLimit.id}` | `"7"` |
| 流量上报 n | 完整服务名(含协议后缀) | `12_3_4_tcp` |
| 特殊服务 | `web_api`(不上报、不清理) | `web_api` |

## 5. 流量上报(HTTP)

### 5.1 `/flow/upload`(周期 5s)

- `POST /flow/upload?secret=<secret>`;UA `GOST-Traffic-Reporter/1.0`;超时 5s
- Body(数组,可加密包裹):
```json
[{"n":"12_3_4_tcp","u":1024,"d":2048}]
```
- `n` = 完整服务名;`u` = OutputBytes 增量 → **out_flow**;`d` = InputBytes 增量 → **in_flow**(命名反直觉,勿反)
- 收到响应 `"ok"`(TrimSpace 精确匹配)才清空累积;失败保留重发;无流量服务不出现
- Java 处理链:`n.split("_")` 取前三段 [forwardId, userId, userTunnelId] → 倍率 `d' = trunc(d × trafficRatio) × tunnel.flow`、`u'` 同理(源码:先 BigDecimal 乘 ratio 后 `longValue()` **截断取整,再 ×flow**;勿整体相乘后一次截断)→ 原子累加 forward/user/user_tunnel(userTunnelId=0 跳过)→ 超限检查(用户:in+out ≥ flow×1024³ / 到期 / status≠1;user_tunnel 同理)→ `PauseService`(chain_type='1' 节点,服务名 name_tcp/name_udp)+ forward.status=0

### 5.2 `/flow/config`(10 分钟;当前 agent 为 dead code,但必须实现)

- `POST /flow/config?secret=<secret>`;UA `Config-Reporter/1.0`;超时 10s
- Body:`{"config": <完整gost配置>}`(可加密);service 附运行期 `status:{state, events[], stats{}}`
- 面板孤儿清理:服务名 `_tls` 结尾 → 查 tunnel,不存在删 `<id>_tls`;`_tcp` 结尾 → 查 forward,不存在删 `_tcp`+`_udp`;链名 → 查 tunnel;限流器名 → 查 speed_limit;`web_api` 跳过。⚠ **源码缺陷:仅 `_tls`/`_tcp` 两个分支,`_udp` 结尾的残留服务(如转发删除后遗留的 `12_3_4_udp`)永不清理——Go 侧建议补上 `_udp` 校验**
- 响应 `"ok"`

### 5.3 `/flow/test`

- GET/POST,恒返回 `"test"`(安装脚本健康检查用)

## 6. 兼容性注意(Go 实现要点)

1. 接收先探测 `encrypted` 包装,再探测 `compressed` 包装
2. AddChains/DeleteChains/UpdateChains/AddLimiters 的 data 存在"对象/字符串/包装"三种形态,解析全兼容
3. `selector.failTimeout` 接受字符串 duration 与纳秒数字
4. `/flow/upload` 响应必须精确为 `ok`
5. agent 重连/命令后写本地文件:gost.json、paused_services.json、config.json(协议外行为,勿干预)
6. 加解密失败降级明文;timestamp 不校验
