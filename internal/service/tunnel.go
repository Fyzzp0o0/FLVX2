package service

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Fyzzp0o0/FLVX2/internal/model"
	"github.com/Fyzzp0o0/FLVX2/internal/ws"
)

var (
	ErrTunnelNotFound = errors.New("隧道不存在")
	ErrTunnelExists   = errors.New("隧道名称已存在")
)

// ChainTunnelIn 隧道链节点入参(chain_type:1入口/2转发链/3出口)
type ChainTunnelIn struct {
	ID       int64   `json:"id"`
	TunnelID int64   `json:"tunnelId"`
	ChainType int    `json:"chainType"`
	NodeID   int64   `json:"nodeId"`
	Port     *int64  `json:"port"`
	Strategy string  `json:"strategy"`
	Inx      int     `json:"inx"`
	Protocol string  `json:"protocol"`
}

// TunnelCreate 创建隧道入参
type TunnelCreate struct {
	Name         string          `json:"name" binding:"required"`
	InNodeID     []ChainTunnelIn `json:"inNodeId"`
	ChainNodes   [][]ChainTunnelIn `json:"chainNodes"`
	OutNodeID    []ChainTunnelIn `json:"outNodeId"`
	InIP         string          `json:"inIp"`
	Type         int             `json:"type"` // 1端口转发/2隧道转发
	Flow         int             `json:"flow"` // 1单向/2双向
	TrafficRatio float64         `json:"trafficRatio"`
}

// TunnelDetail 隧道详情(含节点分组)
type TunnelDetail struct {
	ID           int64            `json:"id"`
	Name         string           `json:"name"`
	Type         int64            `json:"type"`
	Flow         int64            `json:"flow"`
	TrafficRatio float64          `json:"trafficRatio"`
	Status       int64            `json:"status"`
	CreatedTime  int64            `json:"createdTime"`
	UpdatedTime  int64            `json:"updatedTime"`
	InIP         string           `json:"inIp"`
	InNodeID     []ChainTunnelIn  `json:"inNodeId"`
	ChainNodes   [][]ChainTunnelIn `json:"chainNodes"`
	OutNodeID    []ChainTunnelIn  `json:"outNodeId"`
}

// TunnelService 隧道相关(对照 TunnelServiceImpl)
type TunnelService struct {
	pool    *pgxpool.Pool
	nodes   *NodeService
	hub     *ws.Hub
}

func NewTunnelService(pool *pgxpool.Pool, nodes *NodeService, hub *ws.Hub) *TunnelService {
	return &TunnelService{pool: pool, nodes: nodes, hub: hub}
}

// GetByID 查隧道
func (s *TunnelService) GetByID(ctx context.Context, id int64) (*model.Tunnel, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT id, name, traffic_ratio, type, protocol, flow, created_time, updated_time, status, COALESCE(in_ip,'')
		 FROM tunnel WHERE id = $1`, id)
	t := &model.Tunnel{}
	err := row.Scan(&t.ID, &t.Name, &t.TrafficRatio, &t.Type, &t.Protocol, &t.Flow, &t.CreatedTime, &t.UpdatedTime, &t.Status, &t.InIP)
	if err != nil {
		return nil, ErrTunnelNotFound
	}
	return t, nil
}

// Create 创建隧道
func (s *TunnelService) Create(ctx context.Context, dto TunnelCreate) error {
	var cnt int64
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM tunnel WHERE name = $1`, dto.Name).Scan(&cnt); err != nil {
		return err
	}
	if cnt > 0 {
		return ErrTunnelExists
	}
	if dto.Type == 2 && len(dto.OutNodeID) == 0 {
		return errors.New("隧道转发必须填写出口节点")
	}
	// 校验节点存在且启用、不重复
	allNodes := map[int64]bool{}
	for _, ct := range dto.InNodeID {
		allNodes[ct.NodeID] = true
	}
	for _, group := range dto.ChainNodes {
		for _, ct := range group {
			allNodes[ct.NodeID] = true
		}
	}
	for _, ct := range dto.OutNodeID {
		allNodes[ct.NodeID] = true
	}
	for nid := range allNodes {
		n, err := s.nodes.GetByID(ctx, nid)
		if err != nil || n.Status != 1 {
			return errors.New("节点不存在或未启用")
		}
	}
	// 入口节点 IP
	inIP := dto.InIP
	if inIP == "" {
		var ips []string
		for _, ct := range dto.InNodeID {
			n, err := s.nodes.GetByID(ctx, ct.NodeID)
			if err == nil {
				ips = append(ips, n.ServerIP)
			}
		}
		inIP = strings.Join(ips, ",")
	}
	now := time.Now().UnixMilli()
	var tunnelID int64
	err := s.pool.QueryRow(ctx,
		`INSERT INTO tunnel (name, traffic_ratio, type, protocol, flow, created_time, updated_time, status, in_ip)
		 VALUES ($1,$2,$3,'tls',$4,$5,$5,1,$6) RETURNING id`,
		dto.Name, dto.TrafficRatio, dto.Type, dto.Flow, now, inIP).Scan(&tunnelID)
	if err != nil {
		return err
	}
	// 组装 chain_tunnel 并分配端口(入口无端口,转发链/出口自动分配)
	var chains []ChainTunnelIn
	for _, ct := range dto.InNodeID {
		ct.ChainType = 1
		ct.TunnelID = tunnelID
		chains = append(chains, ct)
	}
	inx := 0
	for _, group := range dto.ChainNodes {
		for _, ct := range group {
			ct.ChainType = 2
			ct.TunnelID = tunnelID
			ct.Inx = inx
			if ct.Port == nil {
				p, err := s.GetNodePort(ctx, ct.NodeID, 0)
				if err != nil {
					return err
				}
				pp := int64(p)
				ct.Port = &pp
			}
			chains = append(chains, ct)
		}
		inx++
	}
	for _, ct := range dto.OutNodeID {
		ct.ChainType = 3
		ct.TunnelID = tunnelID
		if ct.Port == nil {
			p, err := s.GetNodePort(ctx, ct.NodeID, 0)
			if err != nil {
				return err
			}
			pp := int64(p)
			ct.Port = &pp
		}
		chains = append(chains, ct)
	}
	for _, ct := range chains {
		if _, err := s.pool.Exec(ctx,
			`INSERT INTO chain_tunnel (tunnel_id, chain_type, node_id, port, strategy, inx, protocol)
			 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			ct.TunnelID, strconv.Itoa(ct.ChainType), ct.NodeID, ct.Port, ct.Strategy, ct.Inx, ct.Protocol); err != nil {
			return err
		}
	}
	// type==2 隧道:下发链与 TLS 服务(对应 TunnelServiceImpl.createTunnel 的 Gost 联动)
	if dto.Type == 2 {
		if err := s.deployChainServices(ctx, tunnelID, dto); err != nil {
			// 回滚:清理已下发与 DB 记录
			_, _ = s.pool.Exec(ctx, `DELETE FROM chain_tunnel WHERE tunnel_id = $1`, tunnelID)
			_, _ = s.pool.Exec(ctx, `DELETE FROM tunnel WHERE id = $1`, tunnelID)
			return err
		}
	}
	return nil
}

// deployChainServices 下发 chains_<id> 链与 <id>_tls 服务(Java TunnelServiceImpl:145-232 逻辑)
func (s *TunnelService) deployChainServices(ctx context.Context, tunnelID int64, dto TunnelCreate) error {
	// 节点详情缓存
	nodeMap := map[int64]*model.Node{}
	getNode := func(id int64) (*model.Node, error) {
		if n, ok := nodeMap[id]; ok {
			return n, nil
		}
		n, err := s.nodes.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		nodeMap[id] = n
		return n, nil
	}
	hopNodes := func(list []ChainTunnelIn) []ChainHopNode {
		var out []ChainHopNode
		for _, ct := range list {
			n, _ := getNode(ct.NodeID)
			if n == nil || ct.Port == nil {
				continue
			}
			out = append(out, ChainHopNode{Inx: ct.Inx, Addr: n.ServerIP + ":" + strconv.FormatInt(*ct.Port, 10), Protocol: ct.Protocol})
		}
		return out
	}
	addChain := func(nodeID int64, strategy string, next []ChainHopNode) error {
		if len(next) == 0 {
			return nil
		}
		n, _ := getNode(nodeID)
		iface := ""
		if n != nil {
			iface = n.InterfaceName
		}
		data := BuildChainData(tunnelID, strategy, iface, next)
		dto := s.hub.SendCommand(nodeID, "AddChains", data)
		if dto.Code != 0 {
			return errors.New(dto.Msg)
		}
		return nil
	}
	addTLSService := func(nodeID int64, ct ChainTunnelIn) error {
		n, err := getNode(ct.NodeID)
		if err != nil || ct.Port == nil {
			return errors.New("节点信息缺失")
		}
		hostNode, _ := getNode(nodeID)
		iface := ""
		if hostNode != nil {
			iface = hostNode.InterfaceName
		}
		svc := BuildChainService(tunnelID, ct.ChainType, ct.Protocol, n, *ct.Port, iface)
		dto := s.hub.SendCommand(nodeID, "AddService", []map[string]any{svc})
		if dto.Code != 0 {
			return errors.New(dto.Msg)
		}
		return nil
	}

	// 入口节点:链指向第一跳或出口
	outHops := hopNodes(dto.OutNodeID)
	for _, in := range dto.InNodeID {
		next := outHops
		if len(dto.ChainNodes) > 0 {
			next = hopNodes(dto.ChainNodes[0])
		}
		strategy := "fifo"
		if len(dto.ChainNodes) > 0 && len(dto.ChainNodes[0]) > 0 {
			strategy = dto.ChainNodes[0][0].Strategy
		}
		if err := addChain(in.NodeID, strategy, next); err != nil {
			return err
		}
	}
	// 转发链节点:链指向下一跳或出口 + 自身 tls 服务
	for i, group := range dto.ChainNodes {
		next := outHops
		if i+1 < len(dto.ChainNodes) {
			next = hopNodes(dto.ChainNodes[i+1])
		}
		for _, ct := range group {
			strategy := ct.Strategy
			if strategy == "" {
				strategy = "fifo"
			}
			if err := addChain(ct.NodeID, strategy, next); err != nil {
				return err
			}
			ct.ChainType = 2
			if err := addTLSService(ct.NodeID, ct); err != nil {
				return err
			}
		}
	}
	// 出口节点:tls 服务
	for _, out := range dto.OutNodeID {
		out.ChainType = 3
		if err := addTLSService(out.NodeID, out); err != nil {
			return err
		}
	}
	return nil
}

// List 全部隧道详情
func (s *TunnelService) List(ctx context.Context) ([]*TunnelDetail, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, name, traffic_ratio, type, protocol, flow, created_time, updated_time, status, COALESCE(in_ip,'') FROM tunnel ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*TunnelDetail
	for rows.Next() {
		t := &model.Tunnel{}
		if err := rows.Scan(&t.ID, &t.Name, &t.TrafficRatio, &t.Type, &t.Protocol, &t.Flow, &t.CreatedTime, &t.UpdatedTime, &t.Status, &t.InIP); err != nil {
			continue
		}
		d := &TunnelDetail{
			ID: t.ID, Name: t.Name, Type: t.Type, Flow: t.Flow, TrafficRatio: t.TrafficRatio,
			Status: t.Status, CreatedTime: t.CreatedTime, UpdatedTime: t.UpdatedTime, InIP: t.InIP,
			InNodeID: []ChainTunnelIn{}, ChainNodes: [][]ChainTunnelIn{}, OutNodeID: []ChainTunnelIn{},
		}
		s.fillChainNodes(ctx, d)
		out = append(out, d)
	}
	return out, nil
}

// fillChainNodes 填充隧道节点分组(chain_type 1/2/3;chainNodes 按 inx 分组)
func (s *TunnelService) fillChainNodes(ctx context.Context, d *TunnelDetail) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, tunnel_id, chain_type, node_id, port, COALESCE(strategy,''), COALESCE(inx,0), COALESCE(protocol,'')
		 FROM chain_tunnel WHERE tunnel_id = $1 ORDER BY id`, d.ID)
	if err != nil {
		return
	}
	defer rows.Close()
	chainByInx := map[int][]ChainTunnelIn{}
	for rows.Next() {
		var ct ChainTunnelIn
		var chainType string
		if err := rows.Scan(&ct.ID, &ct.TunnelID, &chainType, &ct.NodeID, &ct.Port, &ct.Strategy, &ct.Inx, &ct.Protocol); err != nil {
			continue
		}
		ct.ChainType, _ = strconv.Atoi(chainType)
		switch ct.ChainType {
		case 1:
			d.InNodeID = append(d.InNodeID, ct)
		case 2:
			chainByInx[ct.Inx] = append(chainByInx[ct.Inx], ct)
		case 3:
			d.OutNodeID = append(d.OutNodeID, ct)
		}
	}
	keys := make([]int, 0, len(chainByInx))
	for k := range chainByInx {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		d.ChainNodes = append(d.ChainNodes, chainByInx[k])
	}
}

// Update 更新隧道(name/flow/trafficRatio/inIp)
func (s *TunnelService) Update(ctx context.Context, id int64, name string, flow int64, inIP string, trafficRatio float64) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return ErrTunnelNotFound
	}
	if inIP == "" {
		// 按入口节点 serverIp 重拼
		rows, err := s.pool.Query(ctx,
			`SELECT n.server_ip FROM chain_tunnel ct JOIN node n ON ct.node_id = n.id WHERE ct.tunnel_id = $1 AND ct.chain_type = '1'`, id)
		if err == nil {
			var ips []string
			for rows.Next() {
				var ip string
				if rows.Scan(&ip) == nil {
					ips = append(ips, ip)
				}
			}
			rows.Close()
			inIP = strings.Join(ips, ",")
		}
	}
	_, err := s.pool.Exec(ctx,
		`UPDATE tunnel SET name=$2, flow=$3, traffic_ratio=$4, in_ip=$5, updated_time=$6 WHERE id=$1`,
		id, name, flow, trafficRatio, inIP, time.Now().UnixMilli())
	return err
}

// Delete 删除隧道(级联:先删该隧道转发(含 Gost 服务)→ user_tunnel/tunnel → 清链与服务 → chain_tunnel)
func (s *TunnelService) Delete(ctx context.Context, id int64, forwards *ForwardService) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return ErrTunnelNotFound
	}
	// 该隧道全部转发
	fs, err := forwards.listByTunnel(ctx, id)
	if err != nil {
		return err
	}
	for _, f := range fs {
		if err := forwards.DeleteByID(ctx, f.ID, 0, roleAdmin); err != nil {
			return err
		}
	}
	if _, err := s.pool.Exec(ctx, `DELETE FROM user_tunnel WHERE tunnel_id = $1`, id); err != nil {
		return err
	}
	if _, err := s.pool.Exec(ctx, `DELETE FROM tunnel WHERE id = $1`, id); err != nil {
		return err
	}
	// 清节点上的链与服务(按 chain_type 分组节点)
	rows, err := s.pool.Query(ctx,
		`SELECT node_id, chain_type FROM chain_tunnel WHERE tunnel_id = $1`, id)
	if err == nil {
		nodeSet := map[int64]bool{}
		for rows.Next() {
			var nid int64
			var ct string
			if rows.Scan(&nid, &ct) == nil {
				nodeSet[nid] = true
			}
		}
		rows.Close()
		for nid := range nodeSet {
			s.hub.SendCommand(nid, "DeleteChains", ws.DeleteChainsData{Chain: "chains_" + strconv.FormatInt(id, 10)})
			s.hub.SendCommand(nid, "DeleteService", ws.DeleteServicesData{Services: []string{strconv.FormatInt(id, 10) + "_tls"}})
		}
	}
	_, err = s.pool.Exec(ctx, `DELETE FROM chain_tunnel WHERE tunnel_id = $1`, id)
	return err
}

// GetNodePort 分配节点可用端口:parsePorts(node.port) - 已占用(chain_tunnel + forward_port,排除自身转发),取最小
func (s *TunnelService) GetNodePort(ctx context.Context, nodeID, excludeForwardID int64) (int, error) {
	node, err := s.nodes.GetByID(ctx, nodeID)
	if err != nil {
		return 0, errors.New("节点不存在")
	}
	used := map[int]bool{}
	rows, err := s.pool.Query(ctx, `SELECT port FROM chain_tunnel WHERE node_id = $1 AND port IS NOT NULL`, nodeID)
	if err == nil {
		for rows.Next() {
			var p int64
			if rows.Scan(&p) == nil {
				used[int(p)] = true
			}
		}
		rows.Close()
	}
	rows, err = s.pool.Query(ctx, `SELECT port FROM forward_port WHERE node_id = $1 AND forward_id != $2`, nodeID, excludeForwardID)
	if err == nil {
		for rows.Next() {
			var p int64
			if rows.Scan(&p) == nil {
				used[int(p)] = true
			}
		}
		rows.Close()
	}
	for _, p := range parsePorts(node.Port) {
		if !used[p] {
			return p, nil
		}
	}
	return 0, errors.New("节点可用端口不足")
}

// parsePorts 解析端口串 "1000-2000,3000" → 升序去重
func parsePorts(port string) []int {
	var out []int
	seen := map[int]bool{}
	for _, seg := range strings.Split(port, ",") {
		parts := strings.SplitN(seg, "-", 2)
		start, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		end := start
		if len(parts) == 2 {
			end, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
		}
		for p := start; p <= end && p <= 65535; p++ {
			if !seen[p] {
				seen[p] = true
				out = append(out, p)
			}
		}
	}
	sort.Ints(out)
	return out
}

// ---- 用户隧道授权(对照 TunnelController /user/* 与 UserTunnelServiceImpl) ----

// UserTunnelDetail 用户隧道权限详情(联查)
type UserTunnelDetail struct {
	ID            int64   `json:"id"`
	UserID        int64   `json:"userId"`
	TunnelID      int64   `json:"tunnelId"`
	Flow          int64   `json:"flow"`
	Num           int64   `json:"num"`
	FlowResetTime int64   `json:"flowResetTime"`
	ExpTime       int64   `json:"expTime"`
	SpeedID       *int64  `json:"speedId"`
	SpeedLimitName string  `json:"speedLimitName"`
	Speed         int64   `json:"speed"`
	TunnelName    string  `json:"tunnelName"`
	TunnelFlow    int64   `json:"tunnelFlow"`
	InFlow        int64   `json:"inFlow"`
	OutFlow       int64   `json:"outFlow"`
	Status        int64   `json:"status"`
}

// AssignUserTunnel 分配隧道权限(重复检查)
func (s *TunnelService) AssignUserTunnel(ctx context.Context, userID, tunnelID, flow, num, flowResetTime, expTime int64, speedID *int64) error {
	var cnt int64
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM user_tunnel WHERE user_id = $1 AND tunnel_id = $2`, userID, tunnelID).Scan(&cnt); err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("该用户已拥有此隧道权限")
	}
	_, err := s.pool.Exec(ctx,
		`INSERT INTO user_tunnel (user_id, tunnel_id, speed_id, num, flow, in_flow, out_flow, flow_reset_time, exp_time, status)
		 VALUES ($1,$2,$3,$4,$5,0,0,$6,$7,1)`,
		userID, tunnelID, speedID, num, flow, flowResetTime, expTime)
	return err
}

// ListUserTunnels 某用户全部隧道权限
func (s *TunnelService) ListUserTunnels(ctx context.Context, userID int64) ([]*UserTunnelDetail, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT ut.id, ut.user_id, ut.tunnel_id, ut.flow, ut.num, ut.flow_reset_time, ut.exp_time, ut.speed_id,
		        COALESCE(sl.name,''), COALESCE(sl.speed,0), COALESCE(t.name,''), COALESCE(t.flow,0),
		        ut.in_flow, ut.out_flow, ut.status
		 FROM user_tunnel ut
		 LEFT JOIN tunnel t ON ut.tunnel_id = t.id
		 LEFT JOIN speed_limit sl ON ut.speed_id = sl.id
		 WHERE ut.user_id = $1 ORDER BY ut.id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanUserTunnelDetails(rows)
}

// RemoveUserTunnel 移除权限(先删该用户该隧道的全部转发)
func (s *TunnelService) RemoveUserTunnel(ctx context.Context, id int64, forwards *ForwardService) error {
	var userID, tunnelID int64
	if err := s.pool.QueryRow(ctx, `SELECT user_id, tunnel_id FROM user_tunnel WHERE id = $1`, id).Scan(&userID, &tunnelID); err != nil {
		return ErrTunnelNotFound
	}
	fs, err := forwards.listByTunnelAndUser(ctx, tunnelID, userID)
	if err != nil {
		return err
	}
	for _, f := range fs {
		if err := forwards.DeleteByID(ctx, f.ID, 0, roleAdmin); err != nil {
			return err
		}
	}
	_, err = s.pool.Exec(ctx, `DELETE FROM user_tunnel WHERE id = $1`, id)
	return err
}

// UpdateUserTunnel 更新权限(speedId 变化 → UpdateService 重推,在 M4 联动)
func (s *TunnelService) UpdateUserTunnel(ctx context.Context, id int64, flow, num, flowResetTime, expTime, status int64, speedID *int64) error {
	var oldSpeed *int64
	if err := s.pool.QueryRow(ctx, `SELECT speed_id FROM user_tunnel WHERE id = $1`, id).Scan(&oldSpeed); err != nil {
		return ErrTunnelNotFound
	}
	// 修正原 Java bug:更新成功应返回成功(原代码固定返回 R.err)
	if _, err := s.pool.Exec(ctx,
		`UPDATE user_tunnel SET flow=$2, num=$3, flow_reset_time=$4, exp_time=$5, status=$6, speed_id=$7 WHERE id=$1`,
		id, flow, num, flowResetTime, expTime, status, speedID); err != nil {
		return err
	}
	// TODO(M4): speedId 变化时对该用户+隧道的每个转发重推 UpdateService 刷新限速
	_ = oldSpeed
	return nil
}

// UserTunnels 当前用户可见隧道(admin 全部 status=1;普通用户仅授权的)
func (s *TunnelService) UserTunnels(ctx context.Context, userID, roleID int64) ([]*model.Tunnel, error) {
	var rows pgx.Rows
	var err error
	if roleID == roleAdmin {
		rows, err = s.pool.Query(ctx, `SELECT id, name, traffic_ratio, type, protocol, flow, created_time, updated_time, status, COALESCE(in_ip,'') FROM tunnel WHERE status = 1 ORDER BY id`)
	} else {
		rows, err = s.pool.Query(ctx,
			`SELECT t.id, t.name, t.traffic_ratio, t.type, t.protocol, t.flow, t.created_time, t.updated_time, t.status, COALESCE(t.in_ip,'')
			 FROM tunnel t JOIN user_tunnel ut ON t.id = ut.tunnel_id
			 WHERE t.status = 1 AND ut.user_id = $1 AND ut.status = 1 ORDER BY t.id`, userID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*model.Tunnel
	for rows.Next() {
		t := &model.Tunnel{}
		if err := rows.Scan(&t.ID, &t.Name, &t.TrafficRatio, &t.Type, &t.Protocol, &t.Flow, &t.CreatedTime, &t.UpdatedTime, &t.Status, &t.InIP); err == nil {
			out = append(out, t)
		}
	}
	return out, nil
}

func scanUserTunnelDetails(rows interface {
	Next() bool
	Scan(...any) error
}) ([]*UserTunnelDetail, error) {
	var out []*UserTunnelDetail
	for rows.Next() {
		d := &UserTunnelDetail{}
		if err := rows.Scan(&d.ID, &d.UserID, &d.TunnelID, &d.Flow, &d.Num, &d.FlowResetTime, &d.ExpTime, &d.SpeedID,
			&d.SpeedLimitName, &d.Speed, &d.TunnelName, &d.TunnelFlow, &d.InFlow, &d.OutFlow, &d.Status); err == nil {
			out = append(out, d)
		}
	}
	return out, nil
}
