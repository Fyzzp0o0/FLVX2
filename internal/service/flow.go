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
	items := []FlowItem{}
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
		forwards := []model.Forward{}
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
// 修正 Java 缺陷:_udp 结尾残留服务一并清理
func (s *FlowService) ParseConfigSnapshot(ctx context.Context, secret string, rawBody []byte) (*ConfigSnapshot, error) {
	var nodeID int64
	if err := s.pool.QueryRow(ctx, `SELECT id FROM node WHERE secret = $1`, secret).Scan(&nodeID); err != nil {
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
	if wrapped.Config != nil {
		go s.cleanOrphans(nodeID, wrapped.Config) // 异步(Java @Async)
	}
	return wrapped.Config, nil
}

// cleanOrphans 孤儿配置清理:DB 中找不到对应 forward/tunnel/speed_limit 的即删除
func (s *FlowService) cleanOrphans(nodeID int64, snap *ConfigSnapshot) {
	ctx := context.Background()
	cleaned := map[string]bool{} // base 名去重(_tcp/_udp 各出现一次时只删一次)
	for _, svc := range snap.Services {
		name := svc.Name
		if name == "web_api" {
			continue
		}
		parts := strings.Split(name, "_")
		if len(parts) == 0 {
			continue
		}
		id, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			continue
		}
		switch {
		case strings.HasSuffix(name, "_tls"):
			// <tunnelId>_tls → 查 tunnel
			var cnt int64
			_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM tunnel WHERE id = $1`, id).Scan(&cnt)
			if cnt == 0 {
				s.hub.SendCommand(nodeID, "DeleteService", ws.DeleteServicesData{Services: []string{name}})
			}
		case strings.HasSuffix(name, "_tcp"), strings.HasSuffix(name, "_udp"):
			// <fid>_<uid>_<utid>_tcp|udp → 查 forward;不存在删 tcp+udp(修正 Java 漏删 _udp 缺陷)
			base := strings.TrimSuffix(strings.TrimSuffix(name, "_tcp"), "_udp")
			if cleaned[base] {
				continue
			}
			var cnt int64
			_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM forward WHERE id = $1`, id).Scan(&cnt)
			if cnt == 0 {
				cleaned[base] = true
				s.hub.SendCommand(nodeID, "DeleteService", ws.DeleteServicesData{
					Services: []string{base + "_tcp", base + "_udp"},
				})
			}
		}
	}
	for _, ch := range snap.Chains {
		parts := strings.Split(ch.Name, "_")
		id, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)
		if err != nil {
			continue
		}
		var cnt int64
		_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM tunnel WHERE id = $1`, id).Scan(&cnt)
		if cnt == 0 {
			s.hub.SendCommand(nodeID, "DeleteChains", ws.DeleteChainsData{Chain: ch.Name})
		}
	}
	for _, l := range snap.Limiters {
		id, err := strconv.ParseInt(l.Name, 10, 64)
		if err != nil {
			continue
		}
		var cnt int64
		_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM speed_limit WHERE id = $1`, id).Scan(&cnt)
		if cnt == 0 {
			s.hub.SendCommand(nodeID, "DeleteLimiters", ws.DeleteLimitersData{Limiter: l.Name})
		}
	}
}

// entryNodeIDs 隧道入口节点(chain_type='1')
func (s *FlowService) entryNodeIDs(ctx context.Context, tunnelID int64) []int64 {
	ids := []int64{}
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
