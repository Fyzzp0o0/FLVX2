package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Fyzzp0o0/FLVX2/internal/model"
	"github.com/Fyzzp0o0/FLVX2/internal/ws"
)

// SpeedLimitService 限速规则(对照 SpeedLimitServiceImpl;命令联动 Add/Update/DeleteLimiters)
type SpeedLimitService struct {
	pool    *pgxpool.Pool
	tunnels *TunnelService
	hub     *ws.Hub
}

func NewSpeedLimitService(pool *pgxpool.Pool, tunnels *TunnelService, hub *ws.Hub) *SpeedLimitService {
	return &SpeedLimitService{pool: pool, tunnels: tunnels, hub: hub}
}

// Create 创建限速并下发 AddLimiters(隧道全部节点);失败回滚已下发与记录
func (s *SpeedLimitService) Create(ctx context.Context, name string, speed int64, tunnelID int64, tunnelName string) error {
	if _, err := s.tunnels.GetByID(ctx, tunnelID); err != nil {
		return errors.New("隧道不存在")
	}
	now := time.Now().UnixMilli()
	var limitID int64
	if err := s.pool.QueryRow(ctx,
		`INSERT INTO speed_limit (name, speed, tunnel_id, tunnel_name, created_time, updated_time, status)
		 VALUES ($1,$2,$3,$4,$5,$5,1) RETURNING id`,
		name, speed, tunnelID, tunnelName, now).Scan(&limitID); err != nil {
		return err
	}
	speedMBPS := ConvertBitsToMBPS(speed)
	nodes := s.tunnelNodeIDs(ctx, tunnelID)
	success := []int64{}
	for _, nid := range nodes {
		dto := s.hub.SendCommand(nid, "AddLimiters", BuildLimiterData(limitID, speedMBPS))
		if dto.Code != 0 {
			for _, okID := range success {
				s.hub.SendCommand(okID, "DeleteLimiters", ws.DeleteLimitersData{Limiter: itoa(limitID)})
			}
			_, _ = s.pool.Exec(ctx, `DELETE FROM speed_limit WHERE id = $1`, limitID)
			return errors.New(dto.Msg)
		}
		success = append(success, nid)
	}
	return nil
}

// List 全部限速规则
func (s *SpeedLimitService) List(ctx context.Context) ([]*model.SpeedLimit, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, name, speed, tunnel_id, tunnel_name, created_time, COALESCE(updated_time,0), status FROM speed_limit ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*model.SpeedLimit{}
	for rows.Next() {
		sl := &model.SpeedLimit{}
		if err := rows.Scan(&sl.ID, &sl.Name, &sl.Speed, &sl.TunnelID, &sl.TunnelName, &sl.CreatedTime, &sl.UpdatedTime, &sl.Status); err == nil {
			out = append(out, sl)
		}
	}
	return out, nil
}

// Update 更新限速并下发 UpdateLimiters(全部节点)
func (s *SpeedLimitService) Update(ctx context.Context, id int64, name string, speed int64) error {
	var tunnelID int64
	if err := s.pool.QueryRow(ctx, `SELECT tunnel_id FROM speed_limit WHERE id = $1`, id).Scan(&tunnelID); err != nil {
		return errors.New("限速规则不存在")
	}
	if _, err := s.pool.Exec(ctx, `UPDATE speed_limit SET name=$2, speed=$3, updated_time=$4 WHERE id=$1`,
		id, name, speed, time.Now().UnixMilli()); err != nil {
		return err
	}
	speedMBPS := ConvertBitsToMBPS(speed)
	data := ws.UpdateLimitersData{Limiter: itoa(id), Data: ws.LimiterConfig{Name: itoa(id), Limits: []string{"$ " + speedMBPS + "MB " + speedMBPS + "MB"}}}
	for _, nid := range s.tunnelNodeIDs(ctx, tunnelID) {
		if dto := s.hub.SendCommand(nid, "UpdateLimiters", data); dto.Code != 0 {
			return errors.New(dto.Msg)
		}
	}
	return nil
}

// Delete 删除限速(有 user_tunnel 引用则拒绝;下发 DeleteLimiters)
func (s *SpeedLimitService) Delete(ctx context.Context, id int64) error {
	var tunnelID int64
	if err := s.pool.QueryRow(ctx, `SELECT tunnel_id FROM speed_limit WHERE id = $1`, id).Scan(&tunnelID); err != nil {
		return errors.New("限速规则不存在")
	}
	var cnt int64
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM user_tunnel WHERE speed_id = $1`, id).Scan(&cnt); err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("该限速规则还有用户在使用 请先取消分配")
	}
	if _, err := s.pool.Exec(ctx, `DELETE FROM speed_limit WHERE id = $1`, id); err != nil {
		return err
	}
	for _, nid := range s.tunnelNodeIDs(ctx, tunnelID) {
		s.hub.SendCommand(nid, "DeleteLimiters", ws.DeleteLimitersData{Limiter: itoa(id)})
	}
	return nil
}

// Tunnels 供前端选隧道(复用 TunnelService.List)
func (s *SpeedLimitService) Tunnels(ctx context.Context) ([]*TunnelDetail, error) {
	return s.tunnels.List(ctx)
}

// tunnelNodeIDs 隧道全部节点(chain_tunnel 去重)
func (s *SpeedLimitService) tunnelNodeIDs(ctx context.Context, tunnelID int64) []int64 {
	rows, err := s.pool.Query(ctx, `SELECT DISTINCT node_id FROM chain_tunnel WHERE tunnel_id = $1`, tunnelID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	ids := []int64{}
	for rows.Next() {
		var id int64
		if rows.Scan(&id) == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

func itoa(v int64) string {
	return strconv.FormatInt(v, 10)
}
