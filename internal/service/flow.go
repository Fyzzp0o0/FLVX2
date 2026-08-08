package service

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/Fyzzp0o0/FLVX2/internal/model"
	"github.com/Fyzzp0o0/FLVX2/internal/ws"
)

// FlowService 流量上报闭环(对照 FlowController.processFlowData + checkUserRelatedLimits)
type FlowService struct {
	pool   *pgxpool.Pool
	hub    *ws.Hub
	crypto *ws.Crypto
}

func NewFlowService(pool *pgxpool.Pool, hub *ws.Hub) *FlowService {
	return &FlowService{pool: pool, hub: hub, crypto: ws.NewCrypto()}
}

// FlowItem 上报项:n=完整服务名(含 _tcp/_udp 后缀),u=OutputBytes 增量→out_flow,d=InputBytes 增量→in_flow
type FlowItem struct {
	N string `json:"n"`
	U int64  `json:"u"`
	D int64  `json:"d"`
}

// ProcessUpload 处理 /flow/upload;返回给 agent 的响应体(必须精确 "ok")
func (s *FlowService) ProcessUpload(ctx context.Context, secret string, rawBody []byte) string {
	// ① secret 校验(查 node 表;无效 → 静默丢弃但仍回 ok,保持现状行为)
	if err := s.pool.QueryRow(ctx, `SELECT id FROM node WHERE secret = $1`, secret).Scan(new(int64)); err != nil {
		return "ok"
	}
	// ② 解密(失败按明文继续,与 Java 端一致)
	payload, _ := s.crypto.Decrypt(secret, rawBody)
	var items []FlowItem
	if err := json.Unmarshal(payload, &items); err != nil {
		log.Printf("[flow] upload 解析失败: %v", err)
		return "ok"
	}
	for _, it := range items {
		s.processItem(ctx, it)
	}
	return "ok"
}

func (s *FlowService) processItem(ctx context.Context, it FlowItem) {
	if it.N == "web_api" {
		return // 特殊服务,不计流量
	}
	// ③ 解析服务名 {forwardId}_{userId}_{userTunnelId}_(tcp|udp)
	parts := strings.Split(it.N, "_")
	if len(parts) < 4 {
		return
	}
	forwardID, err1 := strconv.ParseInt(parts[0], 10, 64)
	userID, err2 := strconv.ParseInt(parts[1], 10, 64)
	userTunnelID, err3 := strconv.ParseInt(parts[2], 10, 64)
	if err1 != nil || err2 != nil || err3 != nil {
		return
	}

	// ④ 查 forward + tunnel(倍率)
	var ratio float64
	var tunnelFlow int64
	var forwardTunnelID int64
	err := s.pool.QueryRow(ctx,
		`SELECT f.tunnel_id, t.traffic_ratio, t.flow FROM forward f LEFT JOIN tunnel t ON f.tunnel_id = t.id WHERE f.id = $1`,
		forwardID).Scan(&forwardTunnelID, &ratio, &tunnelFlow)
	if err != nil {
		return // forward/tunnel 不存在 → 丢弃
	}
	// 倍率:d' = trunc(d × trafficRatio) × tunnel.flow(BigDecimal 截断语义)
	newD := decimal.NewFromInt(it.D).Mul(decimal.NewFromFloat(ratio)).IntPart() * tunnelFlow
	newU := decimal.NewFromInt(it.U).Mul(decimal.NewFromFloat(ratio)).IntPart() * tunnelFlow

	// ⑤ 原子累加(单事务,UPDATE 自增天然原子)
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, `UPDATE forward SET in_flow = in_flow + $1, out_flow = out_flow + $2, updated_time = $3 WHERE id = $4`,
		newD, newU, time.Now().UnixMilli(), forwardID); err != nil {
		return
	}
	if _, err := tx.Exec(ctx, `UPDATE "user" SET in_flow = in_flow + $1, out_flow = out_flow + $2, updated_time = $3 WHERE id = $4`,
		newD, newU, time.Now().UnixMilli(), userID); err != nil {
		return
	}
	if userTunnelID != 0 {
		if _, err := tx.Exec(ctx, `UPDATE user_tunnel SET in_flow = in_flow + $1, out_flow = out_flow + $2 WHERE id = $3`,
			newD, newU, userTunnelID); err != nil {
			return
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return
	}

	// ⑥ 超限检查(userTunnelId≠0 时;Java 对应 checkUserRelatedLimits/checkUserTunnelRelatedLimits)
	if userTunnelID != 0 {
		s.checkLimits(ctx, userID, userTunnelID)
	}
}

// checkLimits 用户/隧道超限或到期 → 暂停转发(PauseService + status=0)
func (s *FlowService) checkLimits(ctx context.Context, userID, userTunnelID int64) {
	now := time.Now().UnixMilli()
	// 用户维度(严格小于:in+out > limit 才停)
	var uFlow, uIn, uOut, uExp, uStatus int64
	err := s.pool.QueryRow(ctx, `SELECT flow, in_flow, out_flow, exp_time, status FROM "user" WHERE id = $1`, userID).
		Scan(&uFlow, &uIn, &uOut, &uExp, &uStatus)
	if err != nil {
		return
	}
	userOver := uFlow*bytesPerGB < uIn+uOut || (uExp > 0 && uExp <= now) || uStatus != 1

	// 隧道维度(含等号:>= 即停)
	var utFlow, utIn, utOut, utExp, utStatus, utTunnelID int64
	err = s.pool.QueryRow(ctx, `SELECT flow, in_flow, out_flow, exp_time, status, tunnel_id FROM user_tunnel WHERE id = $1`, userTunnelID).
		Scan(&utFlow, &utIn, &utOut, &utExp, &utStatus, &utTunnelID)
	if err != nil {
		return
	}
	tunnelOver := utIn+utOut >= utFlow*bytesPerGB || (utExp > 0 && utExp <= now) || utStatus != 1

	if userOver || tunnelOver {
		var forwards []model.Forward
		rows, err := s.pool.Query(ctx,
			`SELECT id, user_id, user_name, name, tunnel_id, remote_addr, strategy, in_flow, out_flow, created_time, updated_time, status, inx
			 FROM forward WHERE user_id = $1 AND status = 1`, userID)
		if err == nil {
			for rows.Next() {
				var f model.Forward
				if err := rows.Scan(&f.ID, &f.UserID, &f.UserName, &f.Name, &f.TunnelID, &f.RemoteAddr, &f.Strategy,
					&f.InFlow, &f.OutFlow, &f.CreatedTime, &f.UpdatedTime, &f.Status, &f.Inx); err == nil {
					forwards = append(forwards, f)
				}
			}
			rows.Close()
		}
		for _, f := range forwards {
			// 仅隧道超限时:只暂停该用户在该隧道下的转发
			if tunnelOver && !userOver && f.TunnelID != utTunnelID {
				continue
			}
			// 找到该隧道入口节点(chain_type='1'),逐个 PauseService
			nodeIDs := s.entryNodeIDs(ctx, f.TunnelID)
			for _, nodeID := range nodeIDs {
				s.hub.SendCommand(nodeID, "PauseService", ws.PauseServicesData{
					Services: []string{f.Name + "_tcp", f.Name + "_udp"},
				})
			}
			_, _ = s.pool.Exec(ctx, `UPDATE forward SET status = 0, updated_time = $2 WHERE id = $1`, f.ID, now)
		}
	}
}

// ConfigSnapshot 节点上报的配置快照(02-protocol.md §5.2)
type ConfigSnapshot struct {
	Limiters []struct{ Name string `json:"name"` } `json:"limiters"`
	Chains   []struct{ Name string `json:"name"` } `json:"chains"`
	Services []struct{ Name string `json:"name"` } `json:"services"`
}

// ParseConfigSnapshot 解析 /flow/config(解密+解析);孤儿清理(DeleteService/DeleteChains/DeleteLimiters)
// 在 M4 与 Gost 联动时实现(当前版本仅接收解析,返回 ok)
func (s *FlowService) ParseConfigSnapshot(ctx context.Context, secret string, rawBody []byte) (*ConfigSnapshot, error) {
	if err := s.pool.QueryRow(ctx, `SELECT id FROM node WHERE secret = $1`, secret).Scan(new(int64)); err != nil {
		return nil, err
	}
	payload, err := s.crypto.Decrypt(secret, rawBody)
	if err != nil {
		payload = rawBody
	}
	var wrapped struct {
		Config *ConfigSnapshot `json:"config"`
	}
	if err := json.Unmarshal(payload, &wrapped); err != nil {
		return nil, err
	}
	return wrapped.Config, nil
}

// entryNodeIDs 隧道入口节点(chain_type='1')
func (s *FlowService) entryNodeIDs(ctx context.Context, tunnelID int64) []int64 {
	var ids []int64
	rows, err := s.pool.Query(ctx, `SELECT node_id FROM chain_tunnel WHERE tunnel_id = $1 AND chain_type = '1'`, tunnelID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}
