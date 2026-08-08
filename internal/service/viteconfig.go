package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Fyzzp0o0/FLVX2/internal/model"
)

// ViteConfigService 键值配置(对照 ViteConfigServiceImpl,upsert 语义)
type ViteConfigService struct {
	pool *pgxpool.Pool
}

func NewViteConfigService(pool *pgxpool.Pool) *ViteConfigService {
	return &ViteConfigService{pool: pool}
}

// Get 读单条配置(不存在返回空串)
func (s *ViteConfigService) Get(ctx context.Context, name string) (string, error) {
	var value string
	err := s.pool.QueryRow(ctx, `SELECT value FROM vite_config WHERE name = $1`, name).Scan(&value)
	if err != nil {
		return "", nil
	}
	return value, nil
}

// List 全部配置 → map
func (s *ViteConfigService) List(ctx context.Context) (map[string]string, error) {
	rows, err := s.pool.Query(ctx, `SELECT name, value FROM vite_config`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make(map[string]string)
	for rows.Next() {
		var n, v string
		if err := rows.Scan(&n, &v); err == nil {
			out[n] = v
		}
	}
	return out, nil
}

// GetModel 查单条完整记录(供 /api/v1/config/get 返回)
func (s *ViteConfigService) GetModel(ctx context.Context, name string) (*model.ViteConfig, error) {
	row := s.pool.QueryRow(ctx, `SELECT id, name, value, time FROM vite_config WHERE name = $1`, name)
	vc := &model.ViteConfig{}
	err := row.Scan(&vc.ID, &vc.Name, &vc.Value, &vc.Time)
	if err != nil {
		return nil, err
	}
	return vc, nil
}

// Upsert 存在则更新 value+time,否则插入
func (s *ViteConfigService) Upsert(ctx context.Context, name, value string) error {
	now := time.Now().UnixMilli()
	var cnt int64
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM vite_config WHERE name = $1`, name).Scan(&cnt); err != nil {
		return err
	}
	if cnt > 0 {
		_, err := s.pool.Exec(ctx, `UPDATE vite_config SET value = $2, time = $3 WHERE name = $1`, name, value, now)
		return err
	}
	_, err := s.pool.Exec(ctx, `INSERT INTO vite_config (name, value, time) VALUES ($1, $2, $3)`, name, value, now)
	return err
}
