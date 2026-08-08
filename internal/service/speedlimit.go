package service

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Fyzzp0o0/FLVX2/internal/model"
)

// SpeedLimitService 限速规则(对照 SpeedLimitServiceImpl)
type SpeedLimitService struct {
	pool    *pgxpool.Pool
	tunnels *TunnelService
}

func NewSpeedLimitService(pool *pgxpool.Pool, tunnels *TunnelService) *SpeedLimitService {
	return &SpeedLimitService{pool: pool, tunnels: tunnels}
}

// Create 创建限速规则(speed bps → MB/s 保留 1 位小数;AddLimiters 在 M4)
func (s *SpeedLimitService) Create(ctx context.Context, name string, speed int64, tunnelID int64, tunnelName string) error {
	if _, err := s.tunnels.GetByID(ctx, tunnelID); err != nil {
		return errors.New("隧道不存在")
	}
	now := time.Now().UnixMilli()
	_, err := s.pool.Exec(ctx,
		`INSERT INTO speed_limit (name, speed, tunnel_id, tunnel_name, created_time, updated_time, status)
		 VALUES ($1,$2,$3,$4,$5,$5,1)`,
		name, speed, tunnelID, tunnelName, now)
	// TODO(M4): 隧道所有 chain_tunnel 节点 AddLimiters(名=记录 id)
	return err
}

// List 全部限速规则
func (s *SpeedLimitService) List(ctx context.Context) ([]*model.SpeedLimit, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, name, speed, tunnel_id, tunnel_name, created_time, COALESCE(updated_time,0), status FROM speed_limit ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*model.SpeedLimit
	for rows.Next() {
		sl := &model.SpeedLimit{}
		if err := rows.Scan(&sl.ID, &sl.Name, &sl.Speed, &sl.TunnelID, &sl.TunnelName, &sl.CreatedTime, &sl.UpdatedTime, &sl.Status); err == nil {
			out = append(out, sl)
		}
	}
	return out, nil
}

// Update 更新限速(UpdateLimiters 在 M4)
func (s *SpeedLimitService) Update(ctx context.Context, id int64, name string, speed int64) error {
	var cnt int64
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM speed_limit WHERE id = $1`, id).Scan(&cnt); err != nil {
		return err
	}
	if cnt == 0 {
		return errors.New("限速规则不存在")
	}
	_, err := s.pool.Exec(ctx, `UPDATE speed_limit SET name=$2, speed=$3, updated_time=$4 WHERE id=$1`,
		id, name, speed, time.Now().UnixMilli())
	return err
}

// Delete 删除限速(有 user_tunnel 引用则拒绝;DeleteLimiters 在 M4)
func (s *SpeedLimitService) Delete(ctx context.Context, id int64) error {
	var cnt int64
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM user_tunnel WHERE speed_id = $1`, id).Scan(&cnt); err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("该限速规则还有用户在使用 请先取消分配")
	}
	_, err := s.pool.Exec(ctx, `DELETE FROM speed_limit WHERE id = $1`, id)
	return err
}

// Tunnels 供前端选隧道(复用 TunnelService.List)
func (s *SpeedLimitService) Tunnels(ctx context.Context) ([]*TunnelDetail, error) {
	return s.tunnels.List(ctx)
}
