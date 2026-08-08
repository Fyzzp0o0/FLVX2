package db

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql
var schemaSQL string

//go:embed data.sql
var dataSQL string

// Init 建立连接池并执行 schema/data 初始化(幂等,对应原 spring.sql.init.mode=always)
func Init(ctx context.Context, host, port, name, user, password string) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_max_conns=5&pool_min_conns=2",
		user, password, host, port, name)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("创建连接池失败: %w", err)
	}
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("连接 PostgreSQL 失败: %w", err)
	}
	// schema.sql / data.sql 均为幂等语句(CREATE TABLE IF NOT EXISTS / ON CONFLICT DO NOTHING / setval)
	if _, err := pool.Exec(ctx, schemaSQL); err != nil {
		pool.Close()
		return nil, fmt.Errorf("执行 schema.sql 失败: %w", err)
	}
	if _, err := pool.Exec(ctx, dataSQL); err != nil {
		pool.Close()
		return nil, fmt.Errorf("执行 data.sql 失败: %w", err)
	}
	// 常用配置缺省补全(前端配置页固定展示;ip 为空时 node/install 会提示用户先设置)
	defaults := []struct{ name, value string }{
		{"app_name", "flux"},
		{"ip", ""},
		{"url", ""},
		{"announcement", ""},
	}
	now := time.Now().UnixMilli()
	for _, d := range defaults {
		if _, err := pool.Exec(ctx,
			`INSERT INTO vite_config (name, value, time) VALUES ($1, $2, $3) ON CONFLICT (name) DO NOTHING`,
			d.name, d.value, now); err != nil {
			pool.Close()
			return nil, fmt.Errorf("初始化配置 %s 失败: %w", d.name, err)
		}
	}
	return pool, nil
}
