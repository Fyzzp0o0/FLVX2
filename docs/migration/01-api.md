# API 清单(Go 重写参考)— 完整端点规格

> 依据 `springboot-backend/src/main/java/com/admin/controller/` 与 DTO/Service 逐类核对。
> 说明:除特别标注外,响应统一包装 `R = {code, msg, ts(毫秒), data}`,`code=0` 成功、`-1` 业务错误、`-2` 未知异常、401/403 也在 body(HTTP 200)。

## 0. 鉴权模型

- **JWT 拦截**:所有 `/api/**`(除白名单),请求头 `Authorization: <裸token>`(无 Bearer 前缀)
- **白名单(免 JWT)**:`/flow/**`、`/api/v1/open_api/**`、`/api/v1/config/get`、`/api/v1/user/login`、`/api/v1/user/register`、`/api/v1/captcha/**`
- **角色**:`@RequireRole` → 要求 token 内 `role_id == 0`(管理员),否则 `R.err(403, "权限不足，仅管理员可操作")`
- **用户数据隔离**:Forward 全部接口、Tunnel `/user/tunnel`、User `/package`、`/updatePassword` 无 @RequireRole,但 service 层按 token 身份过滤——Go 在 handler 层显式实现
- **token claims**:`{sub: userId, user, name, role_id, iat, exp}`,HS256,90 天

## 1. UserController — `/api/v1/user`

| 路由 | 参数 | 鉴权 | 逻辑 | 响应 data |
|---|---|---|---|---|
| POST `/login` | `{username, password, captchaId?}` | 免 | captcha_enabled=true 时二次校验 captchaId;查用户;`md5(password+盐)` 比对;status==0 → 停用;签发 JWT;`requirePasswordChange = (username==\"admin_user\" \|\| password==\"admin_user\")` | `{token, name, role_id, requirePasswordChange}` |
| POST `/register` | `{user(3-20), pwd(6-32), captchaId?}` | 免 | 验证码校验;查重;读默认配额(`register_default_flow`=100GB、`register_default_num`=5、`register_default_exp_days`=30、`register_default_flow_reset_days`=30);建普通用户 roleId=1 | `null` |
| POST `/create` | `{user, pwd, flow≥0, num≥0, expTime, flowResetTime, status?}` | 管理员 | 查重;md5;roleId=1,status=1 | `null` |
| POST `/list` | 无 | 管理员 | 返回 `role_id != 0` 全部用户(密码脱敏) | `User[]` |
| POST `/update` | `{id, user, pwd?, flow, num, expTime, flowResetTime, status?}` | 管理员 | 禁改 roleId=0("请不要作死");查重排除自身;pwd 非空才更新 | `null` |
| POST `/delete` | `{id}` | 管理员 | 禁删管理员;级联删 Gost 服务+forward+user_tunnel+statistics_flow+user | `null` |
| POST `/package` | 无(token 取身份) | JWT | 用户信息+隧道权限列表+转发列表+近 24h 统计 | `UserPackageDto` |
| POST `/updatePassword` | `{newUsername, currentPassword, newPassword, confirmPassword}` | JWT | confirm 一致;当前密码比对;改用户名查重;同时更新用户名+密码 | `null` |
| POST `/reset` | `{id, type}` | 管理员 | type==1 清零 user 流量;否则清零 user_tunnel 流量 | `null` |

## 2. NodeController — `/api/v1/node`(全部管理员)

| 路由 | 参数 | 逻辑 | 响应 data |
|---|---|---|---|
| POST `/create` | `{name, serverIp, port("80,443,1000-2000"), interfaceName?, tcpListenAddr?, udpListenAddr?}` | 校验端口格式;生成 UUID secret;status=0 | `null` |
| POST `/list` | 无 | 按 status 降序,secret 置 null | `Node[]` |
| POST `/update` | `{id, name, serverIp, port, interfaceName?, http?, tls?, socks?, tcpListenAddr?, udpListenAddr?}` | 在线且协议开关变化 → WS 下发 `SetProtocol`(失败即返回);落库 | `null` |
| POST `/delete` | `{id}` | 级联:引用该节点的隧道逐个 deleteTunnel(含 Gost 清理),再删节点 | `null` |
| POST `/install` | `{id}` | 读 vite_config `ip`;返回安装命令字符串 | `String 命令` |

## 3. TunnelController — `/api/v1/tunnel`

| 路由 | 参数 | 鉴权 | 逻辑 | 响应 data |
|---|---|---|---|---|
| POST `/create` | `{name, inNodeId: ChainTunnel[], chainNodes: ChainTunnel[][], outNodeId?, inIp?, type(1端口/2隧道), flow(1单向/2双向), trafficRatio(0-100)}` | 管理员 | 查重;type=2 必须 outNodeId;节点须存在且 status=1 不重复;自动分配端口;保存 tunnel+chain_tunnel;type=2 时各节点建 `chains_<id>` 链与 `<id>_tls` 服务,失败回滚 | `null` |
| POST `/list` | 无 | 管理员 | 全部隧道 | `TunnelDetailDto[]` |
| POST `/update` | `{id, name, flow, inIp?, trafficRatio}` | 管理员 | 仅更新 name/flow/trafficRatio/inIp | `null` |
| POST `/delete` | `{id}` | 管理员 | 先删该隧道全部转发(删 Gost 服务);删 forward/user_tunnel/tunnel;按 chainType 清 `chains_<id>` 与 `<id>_tls`;删 chain_tunnel | `null` |
| POST `/user/assign` | `{userId, tunnelId, flow≥0, num≥0, flowResetTime, expTime, speedId?}` | 管理员 | 同用户+隧道已存在 → "该用户已拥有此隧道权限";建 UserTunnel(status=1) | `null` |
| POST `/user/list` | `{userId}` | 管理员 | 该用户隧道权限列表 | `UserTunnelWithDetailDto[]` |
| POST `/user/remove` | `{id}` | 管理员 | 删该权限下全部转发(删 Gost 服务),再删权限 | `null` |
| POST `/user/update` | `{id, flow, num, flowResetTime, expTime, status, speedId?}` | 管理员 | speedId 变化 → 对每个转发重推 `UpdateService`;**修正原 bug:成功应返回 R.ok()** | `null` |
| POST `/user/tunnel` | 无(token) | JWT | admin → 全部 status=1 隧道;普通用户 → 其授权隧道 | `Tunnel[]` |
| POST `/diagnose` | `{tunnelId}` | 管理员 | 逐跳 `TcpPing`(入口→www.google.com:443 或 入口→跳→出口→外网) | `{tunnelId, tunnelName, tunnelType, results[], timestamp}` |

`ChainTunnel = {id, tunnelId, chainType, nodeId, port, strategy, inx, protocol}`;`chain_type` 字符串 `"1"`=入口、`"2"`=转发链、`"3"`=出口。

## 4. ForwardController — `/api/v1/forward`(全部仅 JWT;普通用户只能操作自己的)

| 路由 | 参数 | 逻辑 | 响应 data |
|---|---|---|---|
| POST `/create` | `{name, tunnelId, remoteAddr(可逗号多目标), strategy?, inPort?}` | 隧道存在且 status=1;普通用户校验:未到期/有权限/权限未到期/流量>0/数量配额;`get_port` 分配(指定 inPort 须所有入口节点都有,否则取交集最小值,再否则各节点各自第一个可用);存 forward+forward_port;入口节点 `AddService`(tcp+udp 两个服务,绑定 limiter),失败回滚 | `null` |
| POST `/list` | 无 | admin 全部/用户自己的 | `ForwardWithTunnelDto[]` |
| POST `/update` | `{id, userId, name, remoteAddr, strategy?, inPort?}` | 属主校验;配额校验;重分配端口(排除自身);`UpdateService` 重推 | `null` |
| POST `/delete` | `{id}` | 属主校验;删 `_tcp`+`_udp` 服务;删 forward_port/forward | `null` |
| POST `/force-delete` | `{id}` | 属主校验;仅删 DB 记录不调 Gost(孤儿由 /flow/config 清理兜底) | `null` |
| POST `/pause` | `{id}` | 校验账号 status;`PauseService`;forward.status=0 | `null` |
| POST `/resume` | `{id}` | 校验隧道 status=1+账号/权限/流量配额;`ResumeService`;status=1 | `null` |
| POST `/diagnose` | `{forwardId}` | 属主校验;按链路逐跳 TcpPing | `{forwardId, forwardName, tunnelType, results[], timestamp}` |
| POST `/update-order` | `{forwards:[{id, inx}]}` | 普通用户仅自己的;批量更新 inx | `null` |

## 5. SpeedLimitController — `/api/v1/speed-limit`(全部管理员)

| 路由 | 参数 | 逻辑 | 响应 data |
|---|---|---|---|
| POST `/create` | `{name, speed≥1(bps), tunnelId, tunnelName}` | speed/8 保留 1 位小数(MB/s);隧道所有 chain_tunnel 节点 `AddLimiters`(名=记录 id);失败回滚 | `null` |
| POST `/list` | 无 | 全部 | `SpeedLimit[]` |
| POST `/update` | `{id, name, speed}` | 各节点 `UpdateLimiters`;落库 | `null` |
| POST `/delete` | `{id}` | 有 user_tunnel 引用 → "该限速规则还有用户在使用 请先取消分配";各节点 `DeleteLimiters`;删记录 | `null` |
| POST `/tunnels` | 无 | 复用 getAllTunnels(前端选隧道) | `TunnelDetailDto[]` |

## 6. ViteConfigController — `/api/v1/config`

| 路由 | 参数 | 鉴权 | 逻辑 | 响应 data |
|---|---|---|---|---|
| POST `/list` | 无 | JWT(不在白名单,需登录) | 全部配置 | `Map<name,value>` |
| POST `/get` | `{name}` | **免** | 单条,不存在 → "配置不存在" | `ViteConfig` |
| POST `/update` | `Map<String,String>` | 管理员 | 逐条 upsert | `null` |
| POST `/update-single` | `{name, value}` | 管理员 | 单条 upsert | `null` |

## 7. CaptchaController — `/api/v1/captcha`(全部免鉴权;非 R 包装)

| 路由 | 参数 | 逻辑 | 响应 |
|---|---|---|---|
| POST `/check` | 无 | 读 captcha_enabled=="true" ? 1 : 0 | `R{data:0\|1}` |
| POST `/generate` | 无 | 读 captcha_type(RANDOM 时随机 SLIDER/WORD_IMAGE_CLICK/ROTATE/CONCAT) | tianai `CaptchaResponse<ImageCaptchaVO>`(含 id/图片数据) |
| POST `/verify` | `{id, data: 轨迹对象}` | matching(id, track);成功返回 `validToken=id` | `{validToken}` |

> 验证码三端点(`check/generate/verify`)路径保留,响应格式前后端同步重写时**自定**(见 05-frontend-ui.md §4),无兼容负担;`captcha_enabled/captcha_type` 配置语义保留。

## 8. OpenApiController — `/api/v1/open_api`(免鉴权,第三方订阅)

| 路由 | 参数 | 逻辑 | 响应 |
|---|---|---|---|
| GET `/sub_store` | Query `user`、`pwd`、`tunnel`(默认-1) | 校验用户+md5(pwd);tunnel==-1 用用户流量,否则查 user_tunnel;设置响应头 `subscription-userinfo: upload=<inFlow>; download=<outFlow>; total=<flow*1024³>; expire=<expTime/1000>`(源码 `buildSubscriptionHeader(outFlow, inFlow,...)` 内参数位互换,故 header 实际为 upload=inFlow/download=outFlow;与传统"upload=上行"语义相反,保持原样) | 成功:200 + headerValue 字符串;失败:`R.err(...)` |

## 9. FlowController — `/flow`(免鉴权,按 secret 认节点;响应纯文本)

| 路由 | 参数 | 逻辑 | 响应 |
|---|---|---|---|
| POST(注解 `@RequestMapping` 任意方法,agent 实际 POST) `/upload` | Query `secret`;Body 数组 `[{n:完整服务名,u,d}]` 或加密包裹 | secret 校验(查 node 表);**无效 secret → 仍返回 "ok" 但丢弃数据**(agent 视为成功清空累积,现状会静默丢流量);解密;`n=="web_api"` 跳过;倍率 `d'=trunc(d×trafficRatio)×tunnel.flow`、`u'` 同理;原子累加 forward/user/user_tunnel 三表;userTunnelId≠0 时超限/到期/停用 → `PauseService` + forward.status=0 | `"ok"`(必须精确,agent 用 TrimSpace=="ok" 判断,多余内容判失败) |
| POST `/config` | Query `secret`;Body `{config: 完整gost配置}` 或加密包裹 | secret 校验;解析 services/chains/limiters;异步孤儿清理(查不到对应 forward/tunnel/speed_limit 的配置 → DeleteService/DeleteChains/DeleteLimiters) | `"ok"` |
| GET/POST(注解 `@RequestMapping` 任意方法) `/test` | 无 | 探活 | `"test"` |

**流量上报字段语义(逐字对齐)**:`u` = agent OutputBytes 增量 → 记 `out_flow`;`d` = InputBytes 增量 → 记 `in_flow`。上报周期 5s,收到 "ok" 才清空累积。

## 10. 错误兜底

- `/error` 404 → 自定义 HTML("你推开了后端的大门，却发现里面只有寂寞。");其余状态码不写响应体(Go 中可简化:SPA fallback 已处理前端 404)

## 11. 定时任务(非 API,必须保留)

| 任务 | 周期 | 逻辑 |
|---|---|---|
| 流量重置 | 每天 00:00:05 | 按 flow_reset_time(0=不重置,1-31=每月 N 号,月末边界处理)清零 user/user_tunnel 的 in/out_flow;**不碰 status** |
| 到期停服 | 每天 00:00:05 | 到期 user(status=1 且 exp_time<now)→ 暂停其全部 forward + status=0;到期 user_tunnel 同理;停服走 `PauseService`(chain_type='1' 节点) |
| 统计快照 | 每小时整点 | 每用户 total_flow=in+out,flow=本次-上次(负则回退为本次值);写 statistics_flow;清理 48h 前 |
| 孤儿清理 | 触发于 /flow/config | 见 §9 |
