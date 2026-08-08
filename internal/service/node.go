package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Fyzzp0o0/FLVX2/internal/model"
	"github.com/Fyzzp0o0/FLVX2/internal/ws"
)

var ErrNodeNotFound = errors.New("节点不存在")
var ErrNodeInUse = errors.New("节点被隧道引用,无法删除")

var portRangeRe = regexp.MustCompile(`^[0-9]+(-[0-9]+)?(,[0-9]+(-[0-9]+)?)*$`)

// NodeService 节点相关(对照 NodeServiceImpl)
type NodeService struct {
	pool *pgxpool.Pool
	cfg  *ViteConfigService
	hub  *ws.Hub
}

func NewNodeService(pool *pgxpool.Pool, cfg *ViteConfigService, hub *ws.Hub) *NodeService {
	return &NodeService{pool: pool, cfg: cfg, hub: hub}
}

// GetBySecret 按 secret 查节点(WS 握手与流量上报鉴权)
func (s *NodeService) GetBySecret(ctx context.Context, secret string) (*model.Node, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT id, name, secret, server_ip, port, COALESCE(interface_name,''), COALESCE(version,''), http, tls, socks,
		        created_time, COALESCE(updated_time,0), status, tcp_listen_addr, udp_listen_addr
		 FROM node WHERE secret = $1`, secret)
	return scanNode(row)
}

// GetByID 按 ID 查节点
func (s *NodeService) GetByID(ctx context.Context, id int64) (*model.Node, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT id, name, secret, server_ip, port, COALESCE(interface_name,''), COALESCE(version,''), http, tls, socks,
		        created_time, COALESCE(updated_time,0), status, tcp_listen_addr, udp_listen_addr
		 FROM node WHERE id = $1`, id)
	return scanNode(row)
}

func scanNode(row interface{ Scan(...any) error }) (*model.Node, error) {
	n := &model.Node{}
	err := row.Scan(&n.ID, &n.Name, &n.Secret, &n.ServerIP, &n.Port, &n.InterfaceName, &n.Version,
		&n.HTTP, &n.TLS, &n.SOCKS, &n.CreatedTime, &n.UpdatedTime, &n.Status, &n.TCPListenAddr, &n.UDPListenAddr)
	if err != nil {
		return nil, ErrNodeNotFound
	}
	return n, nil
}

// UpdateOnline 节点上线:status=1 + version/http/tls/socks
func (s *NodeService) UpdateOnline(ctx context.Context, id int64, version string, http, tls, socks int64) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE node SET status = 1, version = $2, http = $3, tls = $4, socks = $5, updated_time = $6 WHERE id = $1`,
		id, version, http, tls, socks, time.Now().UnixMilli())
	return err
}

// UpdateOffline 节点下线:仅当 status=1 时置 0
func (s *NodeService) UpdateOffline(ctx context.Context, id int64) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE node SET status = 0, updated_time = $2 WHERE id = $1 AND status = 1`,
		id, time.Now().UnixMilli())
	return err
}

// Create 创建节点(校验端口格式;生成 UUID secret;status=0)
func (s *NodeService) Create(ctx context.Context, name, serverIP, port, interfaceName, tcpListenAddr, udpListenAddr string) (*model.Node, error) {
	if err := validatePortRange(port); err != nil {
		return nil, err
	}
	if tcpListenAddr == "" {
		tcpListenAddr = "0.0.0.0"
	}
	if udpListenAddr == "" {
		udpListenAddr = "0.0.0.0"
	}
	now := time.Now().UnixMilli()
	node := &model.Node{
		Name: name, Secret: uuid.NewString(), ServerIP: serverIP, Port: port,
		InterfaceName: interfaceName, TCPListenAddr: tcpListenAddr, UDPListenAddr: udpListenAddr,
		CreatedTime: now, UpdatedTime: now, Status: 0,
	}
	err := s.pool.QueryRow(ctx,
		`INSERT INTO node (name, secret, server_ip, port, interface_name, tcp_listen_addr, udp_listen_addr, created_time, updated_time, status)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$8,0) RETURNING id`,
		node.Name, node.Secret, node.ServerIP, node.Port, node.InterfaceName, node.TCPListenAddr, node.UDPListenAddr, now,
	).Scan(&node.ID)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// List 全部节点按 status 降序(secret 由 handler 置 null)
func (s *NodeService) List(ctx context.Context) ([]*model.Node, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, name, secret, server_ip, port, COALESCE(interface_name,''), COALESCE(version,''), http, tls, socks,
		        created_time, COALESCE(updated_time,0), status, tcp_listen_addr, udp_listen_addr
		 FROM node ORDER BY status DESC, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*model.Node{}
	for rows.Next() {
		n, err := scanNode(rows)
		if err == nil {
			out = append(out, n)
		}
	}
	return out, nil
}

// Update 更新节点;在线且协议开关变化 → SetProtocol 全量下发
func (s *NodeService) Update(ctx context.Context, id int64, name, serverIP, port, interfaceName string,
	http, tls, socks int64, tcpListenAddr, udpListenAddr string) error {
	node, err := s.GetByID(ctx, id)
	if err != nil {
		return ErrNodeNotFound
	}
	if err := validatePortRange(port); err != nil {
		return err
	}
	if node.Status == 1 && (node.HTTP != http || node.TLS != tls || node.SOCKS != socks) {
		// SetProtocol 必须全量(未提供字段 agent 视为 0 覆盖)
		dto := s.hub.SendCommand(id, "SetProtocol", ws.SetProtocolData{HTTP: int(http), TLS: int(tls), SOCKS: int(socks)})
		if dto.Code != 0 {
			return errors.New(dto.Msg)
		}
	}
	_, err = s.pool.Exec(ctx,
		`UPDATE node SET name=$2, server_ip=$3, port=$4, interface_name=$5, http=$6, tls=$7, socks=$8,
		        tcp_listen_addr=$9, udp_listen_addr=$10, updated_time=$11 WHERE id=$1`,
		id, name, serverIP, port, interfaceName, http, tls, socks, tcpListenAddr, udpListenAddr, time.Now().UnixMilli())
	return err
}

// Delete 删除节点(引用该节点的隧道先逐个删除)
func (s *NodeService) Delete(ctx context.Context, id int64, tunnels *TunnelService, forwards *ForwardService) error {
	// 检查引用:chain_tunnel 按 tunnel_id 分组
	rows, err := s.pool.Query(ctx, `SELECT DISTINCT tunnel_id FROM chain_tunnel WHERE node_id = $1`, id)
	if err != nil {
		return err
	}
	tunnelIDs := []int64{}
	for rows.Next() {
		var tid int64
		if err := rows.Scan(&tid); err == nil {
			tunnelIDs = append(tunnelIDs, tid)
		}
	}
	rows.Close()
	for _, tid := range tunnelIDs {
		if err := tunnels.Delete(ctx, tid, forwards); err != nil {
			return err
		}
	}
	_, err = s.pool.Exec(ctx, `DELETE FROM node WHERE id = $1`, id)
	return err
}

// InstallCmd 返回节点安装命令(读 vite_config ip)
func (s *NodeService) InstallCmd(ctx context.Context, id int64) (string, error) {
	node, err := s.GetByID(ctx, id)
	if err != nil {
		return "", ErrNodeNotFound
	}
	ip, _ := s.cfg.Get(ctx, "ip")
	if ip == "" {
		return "", errors.New("请先前往网站配置中设置ip")
	}
	return fmt.Sprintf("curl -L https://github.com/bqlpfy/flux-panel/releases/download/2.0.7-beta/install.sh -o ./install.sh && chmod +x ./install.sh && ./install.sh -a %s -s %s", ip, node.Secret), nil
}

// validatePortRange 校验端口范围串 "80,443,1000-2000"
func validatePortRange(port string) error {
	if !portRangeRe.MatchString(port) {
		return errors.New("端口格式错误")
	}
	for _, seg := range strings.Split(port, ",") {
		parts := strings.SplitN(seg, "-", 2)
		start, err := strconv.Atoi(parts[0])
		if err != nil || start < 0 || start > 65535 {
			return errors.New("端口格式错误")
		}
		if len(parts) == 2 {
			end, err := strconv.Atoi(parts[1])
			if err != nil || end < start || end > 65535 {
				return errors.New("端口格式错误")
			}
		}
	}
	return nil
}
