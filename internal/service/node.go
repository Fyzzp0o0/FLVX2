package service

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Fyzzp0o0/FLVX2/internal/model"
)

var ErrNodeNotFound = errors.New("节点不存在")

// NodeService 节点相关(对照 NodeServiceImpl)
type NodeService struct {
	pool *pgxpool.Pool
}

func NewNodeService(pool *pgxpool.Pool) *NodeService { return &NodeService{pool: pool} }

// GetBySecret 按 secret 查节点(WS 握手与流量上报鉴权)
func (s *NodeService) GetBySecret(ctx context.Context, secret string) (*model.Node, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT id, name, secret, server_ip, port, COALESCE(interface_name,''), COALESCE(version,''), http, tls, socks,
		        created_time, COALESCE(updated_time,0), status, tcp_listen_addr, udp_listen_addr
		 FROM node WHERE secret = $1`, secret)
	n := &model.Node{}
	err := row.Scan(&n.ID, &n.Name, &n.Secret, &n.ServerIP, &n.Port, &n.InterfaceName, &n.Version,
		&n.HTTP, &n.TLS, &n.SOCKS, &n.CreatedTime, &n.UpdatedTime, &n.Status, &n.TCPListenAddr, &n.UDPListenAddr)
	if err != nil {
		return nil, ErrNodeNotFound
	}
	return n, nil
}

// GetByID 按 ID 查节点
func (s *NodeService) GetByID(ctx context.Context, id int64) (*model.Node, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT id, name, secret, server_ip, port, COALESCE(interface_name,''), COALESCE(version,''), http, tls, socks,
		        created_time, COALESCE(updated_time,0), status, tcp_listen_addr, udp_listen_addr
		 FROM node WHERE id = $1`, id)
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
