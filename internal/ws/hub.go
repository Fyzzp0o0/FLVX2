package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Fyzzp0o0/FLVX2/internal/jwt"
	"github.com/Fyzzp0o0/FLVX2/internal/model"
)

// 协议常量(对齐 02-protocol.md)
const (
	cmdTimeout       = 10 * time.Second // 面板下行命令等待响应超时
	readLimit        = 8 << 20          // 8MB 读上限(>1MB 消息写超时 30s 由 agent 侧处理)
	heartbeatMarker  = "memory_usage"   // agent 心跳 payload 特征子串
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true }, // 对应 setAllowedOrigins("*")
}

// client 一条 WS 连接
type client struct {
	conn    *websocket.Conn
	isNode  bool
	nodeID  int64  // type=1
	adminID int64  // type=0
	secret  string // 节点密钥(node 会话用,命令加密)
	writeMu sync.Mutex
}

// NodeStore 节点存取接口(由 service.NodeService 实现;接口定义在消费方避免 import cycle)
type NodeStore interface {
	GetBySecret(ctx context.Context, secret string) (*model.Node, error)
	UpdateOnline(ctx context.Context, id int64, version string, http, tls, socks int64) error
	UpdateOffline(ctx context.Context, id int64) error
}

// GostDto 命令下发结果(兼容 Java GostDto{code,msg,data})
type GostDto struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Hub 管理节点/管理员会话、广播与命令下发
type Hub struct {
	pool   *pgxpool.Pool
	nodes  NodeStore
	jwtKey string
	crypto *Crypto

	mu            sync.Mutex
	nodeSessions  map[int64]*client
	adminSessions map[int64]map[*client]bool

	pendingMu sync.Mutex
	pending   map[string]chan commandResponse
}

type commandResponse struct {
	Type    string          `json:"type"`
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
	ReqID   string          `json:"requestId"`
}

func NewHub(pool *pgxpool.Pool, nodes NodeStore, jwtKey string) *Hub {
	return &Hub{
		pool:          pool,
		nodes:         nodes,
		jwtKey:        jwtKey,
		crypto:        NewCrypto(),
		nodeSessions:  make(map[int64]*client),
		adminSessions: make(map[int64]map[*client]bool),
		pending:       make(map[string]chan commandResponse),
	}
}

// HandleWS /system-info 入口(type=1 节点 secret 鉴权;type=0 管理员 JWT 鉴权)
func (h *Hub) HandleWS(c *gin.Context) {
	q := c.Request.URL.Query()
	connType := q.Get("type")
	secret := q.Get("secret")

	var cli *client
	switch connType {
	case "1":
		node, err := h.nodes.GetBySecret(c.Request.Context(), secret)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest) // 对应握手拒绝 HTTP 400
			return
		}
		httpVal, _ := strconv.ParseInt(q.Get("http"), 10, 64)
		tlsVal, _ := strconv.ParseInt(q.Get("tls"), 10, 64)
		socksVal, _ := strconv.ParseInt(q.Get("socks"), 10, 64)
		version := q.Get("version")
		cli = &client{isNode: true, nodeID: node.ID, secret: secret}
		// 握手后写 node 状态(与 WebSocketServer.afterConnectionEstablished 一致)
		_ = h.nodes.UpdateOnline(context.Background(), node.ID, version, httpVal, tlsVal, socksVal)
	case "0":
		uid, _, _, valid := jwt.ValidateToken(h.jwtKey, secret)
		if !valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		cli = &client{isNode: false, adminID: uid}
	default:
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[ws] 升级失败: %v", err)
		return
	}
	conn.SetReadLimit(readLimit)
	cli.conn = conn

	if cli.isNode {
		h.registerNode(cli)
	} else {
		h.registerAdmin(cli)
	}
	h.readLoop(cli)
}

// registerNode 节点上线:同节点新连接覆盖旧连接并主动 close;广播 status:1
func (h *Hub) registerNode(cli *client) {
	h.mu.Lock()
	if old, ok := h.nodeSessions[cli.nodeID]; ok && old != cli {
		_ = old.conn.Close() // 旧连接被覆盖,关闭后不触发下线广播(判定见 unregisterNode)
	}
	h.nodeSessions[cli.nodeID] = cli
	admins := h.adminSnapshot()
	h.mu.Unlock()

	h.broadcastTo(admins, map[string]any{"id": strconv.FormatInt(cli.nodeID, 10), "type": "status", "data": 1})
	log.Printf("[ws] 节点上线 id=%d", cli.nodeID)
}

func (h *Hub) registerAdmin(cli *client) {
	h.mu.Lock()
	if h.adminSessions[cli.adminID] == nil {
		h.adminSessions[cli.adminID] = make(map[*client]bool)
	}
	h.adminSessions[cli.adminID][cli] = true
	h.mu.Unlock()
}

// unregisterNode 节点下线:仅当该连接仍是活跃会话时才更新 DB 并广播
func (h *Hub) unregisterNode(cli *client) {
	h.mu.Lock()
	active := false
	if cur, ok := h.nodeSessions[cli.nodeID]; ok && cur == cli {
		delete(h.nodeSessions, cli.nodeID)
		active = true
	}
	admins := h.adminSnapshot()
	h.mu.Unlock()
	if active {
		_ = h.nodes.UpdateOffline(context.Background(), cli.nodeID)
		h.broadcastTo(admins, map[string]any{"id": strconv.FormatInt(cli.nodeID, 10), "type": "status", "data": 0})
		log.Printf("[ws] 节点下线 id=%d", cli.nodeID)
	}
}

func (h *Hub) unregisterAdmin(cli *client) {
	h.mu.Lock()
	if set, ok := h.adminSessions[cli.adminID]; ok {
		delete(set, cli)
		if len(set) == 0 {
			delete(h.adminSessions, cli.adminID)
		}
	}
	h.mu.Unlock()
}

func (h *Hub) adminSnapshot() []*client {
	var out []*client
	for _, set := range h.adminSessions {
		for cli := range set {
			out = append(out, cli)
		}
	}
	return out
}

func (h *Hub) broadcastTo(clients []*client, msg any) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	for _, cli := range clients {
		cli.writeRaw(data)
	}
}

// writeRaw 明文写(管理员广播;节点命令走 encryptWrite)
func (c *client) writeRaw(data []byte) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		_ = c.conn.Close()
	}
}

// encryptWrite 加密写(节点会话,密钥=node secret)
func (c *client) encryptWrite(h *Hub, plaintext []byte) error {
	env, err := h.crypto.Encrypt(c.secret, plaintext)
	if err != nil {
		return err
	}
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return c.conn.WriteMessage(websocket.TextMessage, env)
}

// readLoop 读循环:解密 → 心跳识别/广播 → 命令响应匹配
func (h *Hub) readLoop(cli *client) {
	defer func() {
		if cli.isNode {
			h.unregisterNode(cli)
		} else {
			h.unregisterAdmin(cli)
		}
		_ = cli.conn.Close()
	}()
	for {
		_, raw, err := cli.conn.ReadMessage()
		if err != nil {
			return
		}
		// 解密(明文消息原样;解密失败按明文继续,与 Java 端一致)
		payload, _ := h.crypto.Decrypt(cli.secret, raw)
		if cli.isNode {
			h.handleNodeMessage(cli, payload)
		} else {
			h.handleAdminMessage(cli, payload)
		}
	}
}

// handleNodeMessage 节点上行:心跳(含 memory_usage)→ 回 {"type":"call"} + 广播 info;命令响应 → 匹配 pending
func (h *Hub) handleNodeMessage(cli *client, payload []byte) {
	if strings.Contains(string(payload), heartbeatMarker) {
		// 心跳确认(agent 忽略 type=="call" 的消息)
		_ = cli.encryptWrite(h, []byte(`{"type":"call"}`))
		h.mu.Lock()
		admins := h.adminSnapshot()
		h.mu.Unlock()
		h.broadcastTo(admins, map[string]any{
			"id":   strconv.FormatInt(cli.nodeID, 10),
			"type": "info",
			"data": string(payload), // 前端再 JSON.parse
		})
		return
	}
	var resp commandResponse
	if err := json.Unmarshal(payload, &resp); err != nil || !strings.HasSuffix(resp.Type, "Response") {
		return
	}
	if resp.ReqID != "" {
		h.pendingMu.Lock()
		ch, ok := h.pending[resp.ReqID]
		if ok {
			delete(h.pending, resp.ReqID)
		}
		h.pendingMu.Unlock()
		if ok {
			select {
			case ch <- resp:
			default:
			}
		}
	}
}

func (h *Hub) handleAdminMessage(cli *client, payload []byte) {
	// 管理员会话当前无上行协议(前端仅收广播),忽略
	_ = payload
}

// SendCommand 向在线节点下发命令(加密 + requestId 匹配 + 10s 超时)
// 兼容 Java 容错:Add 类 message 含 "exists"、Delete 类含 "not found" 视为成功(幂等)
func (h *Hub) SendCommand(nodeID int64, cmdType string, data any) GostDto {
	h.mu.Lock()
	cli, ok := h.nodeSessions[nodeID]
	h.mu.Unlock()
	if !ok {
		return GostDto{Code: -1, Msg: "节点不在线"}
	}
	reqID := uuid.NewString()
	msg := map[string]any{"type": cmdType, "data": data, "requestId": reqID}
	plain, err := json.Marshal(msg)
	if err != nil {
		return GostDto{Code: -1, Msg: "命令序列化失败"}
	}
	ch := make(chan commandResponse, 1)
	h.pendingMu.Lock()
	h.pending[reqID] = ch
	h.pendingMu.Unlock()

	if err := cli.encryptWrite(h, plain); err != nil {
		h.pendingMu.Lock()
		delete(h.pending, reqID)
		h.pendingMu.Unlock()
		return GostDto{Code: -1, Msg: "发送失败: " + err.Error()}
	}

	select {
	case resp := <-ch:
		if resp.Success && resp.Message == "OK" {
			return GostDto{Code: 0, Msg: "OK", Data: resp.Data}
		}
		msgText := resp.Message
		// 幂等容错:exists/not found 视为成功
		if strings.Contains(msgText, "exists") || strings.Contains(msgText, "not found") {
			return GostDto{Code: 0, Msg: msgText, Data: resp.Data}
		}
		return GostDto{Code: -1, Msg: msgText, Data: resp.Data}
	case <-time.After(cmdTimeout):
		h.pendingMu.Lock()
		delete(h.pending, reqID)
		h.pendingMu.Unlock()
		return GostDto{Code: -1, Msg: "等待响应超时"}
	}
}
