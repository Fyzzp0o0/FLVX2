package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Fyzzp0o0/FLVX2/internal/model"
	"github.com/Fyzzp0o0/FLVX2/internal/ws"
)

var ErrForwardNotFound = errors.New("转发不存在")
var ErrNoPermission = errors.New("无权操作该转发")

// ForwardWithTunnel 转发列表项(带隧道信息)
type ForwardWithTunnel struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	InIP        string `json:"inIp"`
	InPort      int64  `json:"inPort"`
	RemoteAddr  string `json:"remoteAddr"`
	Status      int64  `json:"status"`
	CreatedTime int64  `json:"createdTime"`
	UpdatedTime int64  `json:"updatedTime"`
	TunnelName  string `json:"tunnelName"`
	UserName    string `json:"userName"`
	UserID      int64  `json:"userId"`
	TunnelID    int64  `json:"tunnelId"`
	InFlow      int64  `json:"inFlow"`
	OutFlow     int64  `json:"outFlow"`
	Strategy    string `json:"strategy"`
	Inx         int64  `json:"inx"`
}

// ForwardService 转发相关(对照 ForwardServiceImpl)
type ForwardService struct {
	pool    *pgxpool.Pool
	tunnels *TunnelService
	users   *UserService
	nodes   *NodeService
	hub     *ws.Hub
}

func NewForwardService(pool *pgxpool.Pool, tunnels *TunnelService, users *UserService, nodes *NodeService, hub *ws.Hub) *ForwardService {
	return &ForwardService{pool: pool, tunnels: tunnels, users: users, nodes: nodes, hub: hub}
}

// Create 创建转发(校验 + 端口分配 + 存库;AddService 在 M4)
func (s *ForwardService) Create(ctx context.Context, name string, tunnelID int64, remoteAddr, strategy string, inPort int64, userID int64, userName string, roleID int64) error {
	tunnel, err := s.tunnels.GetByID(ctx, tunnelID)
	if err != nil || tunnel.Status != 1 {
		return errors.New("隧道不存在或未启用")
	}
	// 用户隧道权限(userTunnelId=0 表示管理员转发)
	var userTunnelID int64
	var speedID *int64
	if roleID != roleAdmin {
		u, err := s.users.GetByID(ctx, userID)
		if err != nil {
			return errors.New("用户不存在")
		}
		if u.ExpTime > 0 && u.ExpTime <= time.Now().UnixMilli() {
			return errors.New("账号已到期")
		}
		if u.Flow <= 0 {
			return errors.New("用户流量已用完")
		}
		// 该隧道权限
		var utID, utFlow, utExp, utStatus int64
		var utSpeed *int64
		err = s.pool.QueryRow(ctx,
			`SELECT id, flow, exp_time, status, speed_id FROM user_tunnel WHERE user_id = $1 AND tunnel_id = $2`,
			userID, tunnelID).Scan(&utID, &utFlow, &utExp, &utStatus, &utSpeed)
		if err != nil {
			return errors.New("没有该隧道权限")
		}
		if utStatus != 1 || (utExp > 0 && utExp <= time.Now().UnixMilli()) || utFlow <= 0 {
			return errors.New("隧道权限不可用或已到期")
		}
		userTunnelID = utID
		speedID = utSpeed
		// 数量配额
		var cnt1, cnt2 int64
		_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM forward WHERE user_id = $1`, userID).Scan(&cnt1)
		_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM forward WHERE user_id = $1 AND tunnel_id = $2`, userID, tunnelID).Scan(&cnt2)
		if cnt1 >= u.Num || cnt2 >= utFlow {
			return errors.New("转发数量已达上限")
		}
	}
	// 入口节点
	entryNodes := entryNodeIDs(ctx, s.pool, tunnelID)
	if len(entryNodes) == 0 {
		return errors.New("隧道没有入口节点")
	}
	// 端口分配:指定 inPort 须所有入口节点可用;否则取交集最小;否则各节点各自可用
	ports := make([]int, 0, len(entryNodes))
	if inPort > 0 {
		for _, nid := range entryNodes {
			avail, err := s.availablePorts(ctx, nid, 0)
			if err != nil || !contains(avail, int(inPort)) {
				return errors.New("指定端口不可用")
			}
			ports = append(ports, int(inPort))
		}
	} else {
		var intersection []int
		for i, nid := range entryNodes {
			avail, err := s.availablePorts(ctx, nid, 0)
			if err != nil {
				return err
			}
			if i == 0 {
				intersection = avail
			} else {
				intersection = intersect(intersection, avail)
			}
		}
		if len(intersection) > 0 {
			minP := intersection[0]
			for _, p := range intersection {
				if p < minP {
					minP = p
				}
			}
			for range entryNodes {
				ports = append(ports, minP)
			}
		} else {
			for _, nid := range entryNodes {
				avail, err := s.availablePorts(ctx, nid, 0)
				if err != nil {
					return err
				}
				ports = append(ports, avail[0])
			}
		}
	}
	now := time.Now().UnixMilli()
	if strategy == "" {
		strategy = "fifo"
	}
	// 插入 forward 拿 id,再回填服务名
	var forwardID int64
	if err := s.pool.QueryRow(ctx,
		`INSERT INTO forward (user_id, user_name, name, tunnel_id, remote_addr, strategy, in_flow, out_flow, created_time, updated_time, status, inx)
		 VALUES ($1,$2,'',$3,$4,$5,0,0,$6,$6,1,0) RETURNING id`,
		userID, userName, tunnelID, remoteAddr, strategy, now).Scan(&forwardID); err != nil {
		return err
	}
	svcName := fmt.Sprintf("%d_%d_%d", forwardID, userID, userTunnelID)
	if _, err := s.pool.Exec(ctx, `UPDATE forward SET name = $2 WHERE id = $1`, forwardID, svcName); err != nil {
		return err
	}
	for i, nid := range entryNodes {
		if _, err := s.pool.Exec(ctx,
			`INSERT INTO forward_port (forward_id, node_id, port) VALUES ($1,$2,$3)`,
			forwardID, nid, ports[i]); err != nil {
			return err
		}
	}
	// AddService 下发(tcp+udp,含 limiter);失败回滚已建服务与 DB 记录
	tunnelInfo, _ := s.tunnels.GetByID(ctx, tunnelID)
	for i, nid := range entryNodes {
		node, err := s.nodes.GetByID(ctx, nid)
		if err != nil {
			return err
		}
		services := BuildForwardServices(svcName, speedID, node, int64(ports[i]), tunnelInfo, remoteAddr)
		dto := s.hub.SendCommand(nid, "AddService", services)
		if dto.Code != 0 {
			for _, prev := range entryNodes[:i] {
				s.hub.SendCommand(prev, "DeleteService", ws.DeleteServicesData{
					Services: []string{svcName + "_tcp", svcName + "_udp"},
				})
			}
			_ = s.DeleteByID(ctx, forwardID, 0, roleAdmin)
			return errors.New(dto.Msg)
		}
	}
	return nil
}

// List 转发列表(管理员全部/普通用户自己的)
func (s *ForwardService) List(ctx context.Context, userID, roleID int64) ([]*ForwardWithTunnel, error) {
	q := `SELECT f.id, f.name, f.tunnel_id, COALESCE(t.in_ip,''), COALESCE(t.name,''), f.remote_addr, f.strategy,
	              f.status, f.created_time, f.updated_time, f.user_id, f.user_name, f.in_flow, f.out_flow, f.inx
	       FROM forward f LEFT JOIN tunnel t ON f.tunnel_id = t.id`
	if roleID != roleAdmin {
		q += " WHERE f.user_id = " + fmt.Sprintf("%d", userID)
	}
	q += " ORDER BY f.created_time DESC"
	rows, err := s.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*ForwardWithTunnel
	for rows.Next() {
		f := &ForwardWithTunnel{}
		if err := rows.Scan(&f.ID, &f.Name, &f.TunnelID, &f.InIP, &f.TunnelName, &f.RemoteAddr, &f.Strategy,
			&f.Status, &f.CreatedTime, &f.UpdatedTime, &f.UserID, &f.UserName, &f.InFlow, &f.OutFlow, &f.Inx); err != nil {
			continue
		}
		// 入口端口(取该转发在任一入口节点的端口;多入口以逗号拼接)
		var ports []string
		prows, err := s.pool.Query(ctx, `SELECT port FROM forward_port WHERE forward_id = $1 ORDER BY id`, f.ID)
		if err == nil {
			for prows.Next() {
				var p int64
				if prows.Scan(&p) == nil {
					ports = append(ports, fmt.Sprintf("%d", p))
				}
			}
			prows.Close()
		}
		if len(ports) > 0 {
			f.InPort, _ = strconv.ParseInt(ports[0], 10, 64)
		}
		out = append(out, f)
	}
	return out, nil
}

// Delete 删除转发(属主校验 + DeleteService)
func (s *ForwardService) Delete(ctx context.Context, id, userID, roleID int64) error {
	f, err := s.getOwned(ctx, id, userID, roleID)
	if err != nil {
		return err
	}
	nodeIDs := entryNodeIDs(ctx, s.pool, f.TunnelID)
	for _, nid := range nodeIDs {
		s.hub.SendCommand(nid, "DeleteService", ws.DeleteServicesData{
			Services: []string{f.Name + "_tcp", f.Name + "_udp"},
		})
	}
	return s.DeleteByID(ctx, id, userID, roleID)
}

// DeleteByID 仅删数据库记录(级联用;roleID=admin 跳过属主校验)
func (s *ForwardService) DeleteByID(ctx context.Context, id, userID, roleID int64) error {
	if _, err := s.getOwned(ctx, id, userID, roleID); err != nil {
		return err
	}
	if _, err := s.pool.Exec(ctx, `DELETE FROM forward_port WHERE forward_id = $1`, id); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx, `DELETE FROM forward WHERE id = $1`, id)
	return err
}

// ForceDelete 仅删 DB 记录不调 Gost(孤儿由 /flow/config 清理兜底)
func (s *ForwardService) ForceDelete(ctx context.Context, id, userID, roleID int64) error {
	return s.DeleteByID(ctx, id, userID, roleID)
}

// Pause 暂停转发
func (s *ForwardService) Pause(ctx context.Context, id, userID, roleID int64) error {
	f, err := s.getOwned(ctx, id, userID, roleID)
	if err != nil {
		return err
	}
	if roleID != roleAdmin {
		u, err := s.users.GetByID(ctx, userID)
		if err != nil || u.Status != 1 {
			return errors.New("账号不可用")
		}
	}
	nodeIDs := entryNodeIDs(ctx, s.pool, f.TunnelID)
	for _, nid := range nodeIDs {
		s.hub.SendCommand(nid, "PauseService", ws.PauseServicesData{
			Services: []string{f.Name + "_tcp", f.Name + "_udp"},
		})
	}
	_, err = s.pool.Exec(ctx, `UPDATE forward SET status = 0, updated_time = $2 WHERE id = $1`, id, time.Now().UnixMilli())
	return err
}

// Resume 恢复转发(校验隧道/账号/权限/流量)
func (s *ForwardService) Resume(ctx context.Context, id, userID, roleID int64) error {
	f, err := s.getOwned(ctx, id, userID, roleID)
	if err != nil {
		return err
	}
	tunnel, err := s.tunnels.GetByID(ctx, f.TunnelID)
	if err != nil || tunnel.Status != 1 {
		return errors.New("隧道不可用")
	}
	if roleID != roleAdmin {
		u, err := s.users.GetByID(ctx, userID)
		if err != nil || u.Status != 1 {
			return errors.New("账号不可用")
		}
		if u.ExpTime > 0 && u.ExpTime <= time.Now().UnixMilli() {
			return errors.New("账号已到期")
		}
		if u.Flow*bytesPerGB <= u.InFlow+u.OutFlow {
			return errors.New("用户流量已用完")
		}
		var ut model.UserTunnel
		err = s.pool.QueryRow(ctx,
			`SELECT id, flow, exp_time, status FROM user_tunnel WHERE user_id = $1 AND tunnel_id = $2`,
			userID, f.TunnelID).Scan(&ut.ID, &ut.Flow, &ut.ExpTime, &ut.Status)
		if err != nil || ut.Status != 1 {
			return errors.New("隧道权限不可用")
		}
		if ut.ExpTime > 0 && ut.ExpTime <= time.Now().UnixMilli() {
			return errors.New("隧道权限已到期")
		}
		if ut.Flow*bytesPerGB <= ut.InFlow+ut.OutFlow {
			return errors.New("隧道流量已用完")
		}
	}
	nodeIDs := entryNodeIDs(ctx, s.pool, f.TunnelID)
	for _, nid := range nodeIDs {
		s.hub.SendCommand(nid, "ResumeService", ws.PauseServicesData{
			Services: []string{f.Name + "_tcp", f.Name + "_udp"},
		})
	}
	_, err = s.pool.Exec(ctx, `UPDATE forward SET status = 1, updated_time = $2 WHERE id = $1`, id, time.Now().UnixMilli())
	return err
}

// UpdateOrder 批量更新排序(普通用户仅自己的)
func (s *ForwardService) UpdateOrder(ctx context.Context, updates []struct {
	ID  int64 `json:"id"`
	Inx int64 `json:"inx"`
}, userID, roleID int64) error {
	if roleID != roleAdmin {
		var mine int64
		_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM forward WHERE user_id = $1`, userID).Scan(&mine)
		if int64(len(updates)) != mine {
			return errors.New("只能更新自己的转发排序")
		}
		for _, u := range updates {
			var cnt int64
			_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM forward WHERE id = $1 AND user_id = $2`, u.ID, userID).Scan(&cnt)
			if cnt == 0 {
				return ErrNoPermission
			}
		}
	}
	for _, u := range updates {
		if _, err := s.pool.Exec(ctx, `UPDATE forward SET inx = $2, updated_time = $3 WHERE id = $1`, u.ID, u.Inx, time.Now().UnixMilli()); err != nil {
			return err
		}
	}
	return nil
}

// TunnelIDOf 查询转发所属隧道(update 用)
func (s *ForwardService) TunnelIDOf(ctx context.Context, id int64) int64 {
	var tid int64
	_ = s.pool.QueryRow(ctx, `SELECT tunnel_id FROM forward WHERE id = $1`, id).Scan(&tid)
	return tid
}

// listByTunnelAndUser 用户在某隧道下的转发(移除隧道权限用)
func (s *ForwardService) listByTunnelAndUser(ctx context.Context, tunnelID, userID int64) ([]*model.Forward, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, user_id, user_name, name, tunnel_id, remote_addr, strategy, in_flow, out_flow, created_time, updated_time, status, inx
		 FROM forward WHERE tunnel_id = $1 AND user_id = $2`, tunnelID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*model.Forward
	for rows.Next() {
		f := &model.Forward{}
		if err := rows.Scan(&f.ID, &f.UserID, &f.UserName, &f.Name, &f.TunnelID, &f.RemoteAddr, &f.Strategy,
			&f.InFlow, &f.OutFlow, &f.CreatedTime, &f.UpdatedTime, &f.Status, &f.Inx); err == nil {
			out = append(out, f)
		}
	}
	return out, nil
}

// listByTunnel 隧道下全部转发(级联删除用)
func (s *ForwardService) listByTunnel(ctx context.Context, tunnelID int64) ([]*model.Forward, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, user_id, user_name, name, tunnel_id, remote_addr, strategy, in_flow, out_flow, created_time, updated_time, status, inx
		 FROM forward WHERE tunnel_id = $1`, tunnelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*model.Forward
	for rows.Next() {
		f := &model.Forward{}
		if err := rows.Scan(&f.ID, &f.UserID, &f.UserName, &f.Name, &f.TunnelID, &f.RemoteAddr, &f.Strategy,
			&f.InFlow, &f.OutFlow, &f.CreatedTime, &f.UpdatedTime, &f.Status, &f.Inx); err == nil {
			out = append(out, f)
		}
	}
	return out, nil
}

// getOwned 取转发并做属主校验
func (s *ForwardService) getOwned(ctx context.Context, id, userID, roleID int64) (*model.Forward, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT id, user_id, user_name, name, tunnel_id, remote_addr, strategy, in_flow, out_flow, created_time, updated_time, status, inx
		 FROM forward WHERE id = $1`, id)
	f := &model.Forward{}
	if err := row.Scan(&f.ID, &f.UserID, &f.UserName, &f.Name, &f.TunnelID, &f.RemoteAddr, &f.Strategy,
		&f.InFlow, &f.OutFlow, &f.CreatedTime, &f.UpdatedTime, &f.Status, &f.Inx); err != nil {
		return nil, ErrForwardNotFound
	}
	if roleID != roleAdmin && f.UserID != userID {
		return nil, ErrNoPermission
	}
	return f, nil
}

// availablePorts 节点可用端口列表(排除已占用)
func (s *ForwardService) availablePorts(ctx context.Context, nodeID, excludeForwardID int64) ([]int, error) {
	node, err := s.nodes.GetByID(ctx, nodeID)
	if err != nil {
		return nil, errors.New("节点不存在")
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
	var out []int
	for _, p := range parsePorts(node.Port) {
		if !used[p] {
			out = append(out, p)
		}
	}
	return out, nil
}

func contains(list []int, v int) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

func intersect(a, b []int) []int {
	set := map[int]bool{}
	for _, v := range b {
		set[v] = true
	}
	var out []int
	for _, v := range a {
		if set[v] {
			out = append(out, v)
		}
	}
	return out
}



// Update 更新转发(保留记录,UpdateService 重推;对应 Java updateForward)
func (s *ForwardService) Update(ctx context.Context, id, userID, roleID int64, name, remoteAddr, strategy string, inPort int64) error {
	f, err := s.getOwned(ctx, id, userID, roleID)
	if err != nil {
		return err
	}
	tunnel, err := s.tunnels.GetByID(ctx, f.TunnelID)
	if err != nil || tunnel.Status != 1 {
		return errors.New("隧道不存在或未启用")
	}
	entryNodes := entryNodeIDs(ctx, s.pool, f.TunnelID)
	if len(entryNodes) == 0 {
		return errors.New("隧道没有入口节点")
	}
	ports := make([]int, 0, len(entryNodes))
	if inPort > 0 {
		for _, nid := range entryNodes {
			avail, err := s.availablePorts(ctx, nid, id)
			if err != nil || !contains(avail, int(inPort)) {
				return errors.New("指定端口不可用")
			}
			ports = append(ports, int(inPort))
		}
	} else {
		for _, nid := range entryNodes {
			avail, err := s.availablePorts(ctx, nid, id)
			if err != nil {
				return err
			}
			ports = append(ports, avail[0])
		}
	}
	if strategy == "" {
		strategy = "fifo"
	}
	now := time.Now().UnixMilli()
	if _, err := s.pool.Exec(ctx,
		`UPDATE forward SET name=$2, remote_addr=$3, strategy=$4, status=1, updated_time=$5 WHERE id=$1`,
		id, name, remoteAddr, strategy, now); err != nil {
		return err
	}
	if _, err := s.pool.Exec(ctx, `DELETE FROM forward_port WHERE forward_id = $1`, id); err != nil {
		return err
	}
	for i, nid := range entryNodes {
		if _, err := s.pool.Exec(ctx,
			`INSERT INTO forward_port (forward_id, node_id, port) VALUES ($1,$2,$3)`,
			id, nid, ports[i]); err != nil {
			return err
		}
	}
	// UpdateService 重推(含 limiter)
	var speedID *int64
	_ = s.pool.QueryRow(ctx, `SELECT speed_id FROM user_tunnel WHERE user_id = $1 AND tunnel_id = $2`,
		userID, f.TunnelID).Scan(&speedID)
	for i, nid := range entryNodes {
		node, err := s.nodes.GetByID(ctx, nid)
		if err != nil {
			return err
		}
		services := BuildForwardServices(f.Name, speedID, node, int64(ports[i]), tunnel, remoteAddr)
		dto := s.hub.SendCommand(nid, "UpdateService", services)
		if dto.Code != 0 {
			return errors.New(dto.Msg)
		}
	}
	return nil
}
