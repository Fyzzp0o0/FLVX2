package service

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Fyzzp0o0/FLVX2/internal/ws"
)

// DiagnosisResult 单跳诊断结果(对齐 Java DiagnosisResult)
type DiagnosisResult struct {
	NodeID       int64   `json:"nodeId"`
	NodeName     string  `json:"nodeName"`
	TargetIP     string  `json:"targetIp"`
	TargetPort   int64   `json:"targetPort"`
	Description  string  `json:"description"`
	Success      bool    `json:"success"`
	Message      string  `json:"message"`
	AverageTime  float64 `json:"averageTime"`
	PacketLoss   float64 `json:"packetLoss"`
	Timestamp    int64   `json:"timestamp"`
	FromChainType string `json:"fromChainType"`
	FromInx      int     `json:"fromInx"`
	ToChainType  string  `json:"toChainType"`
	ToInx        int     `json:"toInx"`
}

// tcpPing 对指定节点执行 TcpPing
func tcpPing(ctx context.Context, hub *ws.Hub, nodeID int64, ip string, port int64) *DiagnosisResult {
	res := &DiagnosisResult{NodeID: nodeID, TargetIP: ip, TargetPort: port, Timestamp: time.Now().UnixMilli()}
	dto := hub.SendCommand(nodeID, "TcpPing", ws.TcpPingData{IP: ip, Port: int(port), Count: 4, Timeout: 5000})
	if dto.Code != 0 {
		res.Message = dto.Msg
		return res
	}
	raw, _ := json.Marshal(dto.Data)
	var resp struct {
		Success      bool    `json:"success"`
		AverageTime  float64 `json:"averageTime"`
		PacketLoss   float64 `json:"packetLoss"`
		ErrorMessage string  `json:"errorMessage"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		res.Message = "响应解析失败"
		return res
	}
	res.Success = resp.Success
	res.AverageTime = resp.AverageTime
	res.PacketLoss = resp.PacketLoss
	if resp.ErrorMessage != "" {
		res.Message = resp.ErrorMessage
	}
	return res
}

// linkNodes 链路节点序列:入口 → 转发链(按 inx) → 出口
type linkNode struct {
	NodeID     int64
	NodeName   string
	ServerIP   string
	ChainType  string
	Inx        int
	Port       *int64
}

func (s *TunnelService) linkNodes(ctx context.Context, tunnelID int64) ([]linkNode, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT ct.node_id, ct.chain_type, ct.inx, ct.port, n.name, n.server_ip
		 FROM chain_tunnel ct JOIN node n ON ct.node_id = n.id
		 WHERE ct.tunnel_id = $1`, tunnelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries, mids, outs []linkNode
	for rows.Next() {
		var ln linkNode
		var inx *int
		if err := rows.Scan(&ln.NodeID, &ln.ChainType, &inx, &ln.Port, &ln.NodeName, &ln.ServerIP); err != nil {
			continue
		}
		if inx != nil {
			ln.Inx = *inx
		}
		switch ln.ChainType {
		case "1":
			entries = append(entries, ln)
		case "2":
			mids = append(mids, ln)
		case "3":
			outs = append(outs, ln)
		}
	}
	sort.Slice(mids, func(i, j int) bool { return mids[i].Inx < mids[j].Inx })
	seq := append(entries, mids...)
	return append(seq, outs...), nil
}

// DiagnoseForward 转发诊断:按链路逐跳 TcpPing,最后一段指向转发目标
func (s *ForwardService) Diagnose(ctx context.Context, forwardID, userID, roleID int64) ([]*DiagnosisResult, error) {
	f, err := s.getOwned(ctx, forwardID, userID, roleID)
	if err != nil {
		return nil, err
	}
	links, err := s.tunnels.linkNodes(ctx, f.TunnelID)
	if err != nil || len(links) == 0 {
		return nil, errors.New("链路信息缺失")
	}
	var results []*DiagnosisResult
	targets := splitAddrs(f.RemoteAddr)
	// 逐跳:node[i] → node[i+1];最后一跳节点 → 每个转发目标
	for i, ln := range links {
		if i+1 < len(links) {
			next := links[i+1]
			port := int64(0)
			if next.Port != nil {
				port = *next.Port
			}
			res := tcpPing(ctx, s.hub, ln.NodeID, next.ServerIP, port)
			res.NodeName = ln.NodeName
			res.Description = ln.NodeName + " → " + next.NodeName
			res.FromChainType = ln.ChainType
			res.FromInx = ln.Inx
			res.ToChainType = next.ChainType
			res.ToInx = next.Inx
			results = append(results, res)
			continue
		}
		for _, target := range targets {
			ip, port := splitHostPort(target)
			res := tcpPing(ctx, s.hub, ln.NodeID, ip, port)
			res.NodeName = ln.NodeName
			res.Description = ln.NodeName + " → " + target
			res.FromChainType = ln.ChainType
			res.FromInx = ln.Inx
			results = append(results, res)
		}
	}
	return results, nil
}

// Diagnose 隧道诊断:端口转发测 www.google.com:443;隧道转发逐跳后出口测外网
func (s *TunnelService) Diagnose(ctx context.Context, tunnelID int64) ([]*DiagnosisResult, error) {
	if _, err := s.GetByID(ctx, tunnelID); err != nil {
		return nil, err
	}
	links, err := s.linkNodes(ctx, tunnelID)
	if err != nil || len(links) == 0 {
		return nil, errors.New("链路信息缺失")
	}
	var results []*DiagnosisResult
	for i, ln := range links {
		if i+1 < len(links) {
			next := links[i+1]
			port := int64(0)
			if next.Port != nil {
				port = *next.Port
			}
			res := tcpPing(ctx, s.hub, ln.NodeID, next.ServerIP, port)
			res.NodeName = ln.NodeName
			res.Description = ln.NodeName + " → " + next.NodeName
			res.FromChainType = ln.ChainType
			res.FromInx = ln.Inx
			res.ToChainType = next.ChainType
			res.ToInx = next.Inx
			results = append(results, res)
			continue
		}
		// 末跳 → 外网(端口转发/隧道统一测 www.google.com:443)
		res := tcpPing(ctx, s.hub, ln.NodeID, "www.google.com", 443)
		res.NodeName = ln.NodeName
		res.Description = ln.NodeName + " → www.google.com:443"
		res.FromChainType = ln.ChainType
		res.FromInx = ln.Inx
		results = append(results, res)
	}
	return results, nil
}

// splitAddrs 逗号分隔目标
func splitAddrs(remoteAddr string) []string {
	var out []string
	for _, a := range strings.Split(remoteAddr, ",") {
		if t := strings.TrimSpace(a); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// splitHostPort 解析 ip:port(IPv6 带括号)
func splitHostPort(addr string) (string, int64) {
	addr = strings.TrimSpace(addr)
	if strings.HasPrefix(addr, "[") {
		if idx := strings.LastIndex(addr, "]:"); idx > 0 {
			port, _ := strconv.ParseInt(addr[idx+2:], 10, 64)
			return addr[1:idx], port
		}
		return strings.Trim(addr, "[]"), 0
	}
	if idx := strings.LastIndex(addr, ":"); idx > 0 {
		port, _ := strconv.ParseInt(addr[idx+1:], 10, 64)
		return addr[:idx], port
	}
	return addr, 0
}
