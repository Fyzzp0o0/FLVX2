package service

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Fyzzp0o0/FLVX2/internal/ws"
)

// 集成测试:依赖本机开发库 flvx2(DB_USER=flvx2 DB_PASSWORD=flvx2)
func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	pool, err := pgxpool.New(context.Background(),
		"postgres://flvx2:flvx2@127.0.0.1:5432/flvx2?sslmode=disable")
	if err != nil {
		t.Fatalf("连接测试库失败: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func TestResetFlowDaily(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	hub := NewHubForTest(pool)
	jobs := NewJobService(pool, hub)

	now := time.Now()
	// 用户 A:今天重置 + 有流量
	var uidA int64
	if err := pool.QueryRow(ctx,
		`INSERT INTO "user" ("user", pwd, role_id, exp_time, flow, in_flow, out_flow, flow_reset_time, num, created_time, updated_time, status)
		 VALUES ('job_test_a','x',1,$1,10,500,300,$2,5,$3,$3,1) RETURNING id`,
		now.Add(24*time.Hour).UnixMilli(), now.Day(), now.UnixMilli()).Scan(&uidA); err != nil {
		t.Fatalf("插入用户A失败: %v", err)
	}
	// 用户 B:已到期 + status=1 + 一个转发
	var uidB, tidB int64
	if err := pool.QueryRow(ctx,
		`INSERT INTO "user" ("user", pwd, role_id, exp_time, flow, in_flow, out_flow, flow_reset_time, num, created_time, updated_time, status)
		 VALUES ('job_test_b','x',1,$1,10,0,0,0,5,$2,$2,1) RETURNING id`,
		now.Add(-time.Hour).UnixMilli(), now.UnixMilli()).Scan(&uidB); err != nil {
		t.Fatalf("插入用户B失败: %v", err)
	}
	if err := pool.QueryRow(ctx,
		`INSERT INTO tunnel (name, traffic_ratio, type, protocol, flow, created_time, updated_time, status)
		 VALUES ('job_test_t',1,1,'tls',1,$1,$1,1) RETURNING id`, now.UnixMilli()).Scan(&tidB); err != nil {
		t.Fatalf("插入隧道失败: %v", err)
	}
	var fidB int64
	if err := pool.QueryRow(ctx,
		`INSERT INTO forward (user_id, user_name, name, tunnel_id, remote_addr, strategy, in_flow, out_flow, created_time, updated_time, status, inx)
		 VALUES ($1,'job_test_b','9_9_9',$2,'1.1.1.1:1','fifo',0,0,$3,$3,1,0) RETURNING id`,
		uidB, tidB, now.UnixMilli()).Scan(&fidB); err != nil {
		t.Fatalf("插入转发失败: %v", err)
	}
	// 到期 user_tunnel
	var utidB int64
	if err := pool.QueryRow(ctx,
		`INSERT INTO user_tunnel (user_id, tunnel_id, num, flow, in_flow, out_flow, flow_reset_time, exp_time, status)
		 VALUES ($1,$2,5,10,0,0,0,$3,1) RETURNING id`,
		uidB, tidB, now.Add(-time.Hour).UnixMilli()).Scan(&utidB); err != nil {
		t.Fatalf("插入 user_tunnel 失败: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM statistics_flow WHERE user_id IN ($1,$2)`, uidA, uidB)
		_, _ = pool.Exec(ctx, `DELETE FROM forward WHERE id = $1`, fidB)
		_, _ = pool.Exec(ctx, `DELETE FROM user_tunnel WHERE id = $1`, utidB)
		_, _ = pool.Exec(ctx, `DELETE FROM tunnel WHERE id = $1`, tidB)
		_, _ = pool.Exec(ctx, `DELETE FROM "user" WHERE id IN ($1,$2)`, uidA, uidB)
	})

	jobs.ResetFlowDaily(ctx)

	// 验证 A 流量清零
	var inA, outA int64
	if err := pool.QueryRow(ctx, `SELECT in_flow, out_flow FROM "user" WHERE id = $1`, uidA).Scan(&inA, &outA); err != nil {
		t.Fatalf("查询用户A失败: %v", err)
	}
	if inA != 0 || outA != 0 {
		t.Errorf("用户A流量未重置: in=%d out=%d", inA, outA)
	}
	// 验证 B 停服
	var statusB int64
	if err := pool.QueryRow(ctx, `SELECT status FROM "user" WHERE id = $1`, uidB).Scan(&statusB); err != nil {
		t.Fatalf("查询用户B失败: %v", err)
	}
	if statusB != 0 {
		t.Errorf("用户B未停服: status=%d", statusB)
	}
	var fStatus int64
	if err := pool.QueryRow(ctx, `SELECT status FROM forward WHERE id = $1`, fidB).Scan(&fStatus); err != nil {
		t.Fatalf("查询转发失败: %v", err)
	}
	if fStatus != 0 {
		t.Errorf("转发未暂停: status=%d", fStatus)
	}
	var utStatus int64
	if err := pool.QueryRow(ctx, `SELECT status FROM user_tunnel WHERE id = $1`, utidB).Scan(&utStatus); err != nil {
		t.Fatalf("查询 user_tunnel 失败: %v", err)
	}
	if utStatus != 0 {
		t.Errorf("user_tunnel 未停服: status=%d", utStatus)
	}
}

func TestSnapshotHourly(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	hub := NewHubForTest(pool)
	jobs := NewJobService(pool, hub)

	now := time.Now()
	var uid int64
	if err := pool.QueryRow(ctx,
		`INSERT INTO "user" ("user", pwd, role_id, exp_time, flow, in_flow, out_flow, flow_reset_time, num, created_time, updated_time, status)
		 VALUES ('job_test_s','x',1,$1,10,100,200,0,5,$2,$2,1) RETURNING id`,
		now.Add(24*time.Hour).UnixMilli(), now.UnixMilli()).Scan(&uid); err != nil {
		t.Fatalf("插入用户失败: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM statistics_flow WHERE user_id = $1`, uid)
		_, _ = pool.Exec(ctx, `DELETE FROM "user" WHERE id = $1`, uid)
	})

	jobs.SnapshotHourly(ctx)

	var total int64
	if err := pool.QueryRow(ctx, `SELECT total_flow FROM statistics_flow WHERE user_id = $1`, uid).Scan(&total); err != nil {
		t.Fatalf("查询快照失败: %v", err)
	}
	if total != 300 {
		t.Errorf("快照 total_flow 错误: %d(期望 300)", total)
	}
}

// NewHubForTest 测试用 Hub(无节点连接,命令下发会返回"节点不在线",不影响 DB 断言)
func NewHubForTest(pool *pgxpool.Pool) *ws.Hub {
	_ = pool
	return ws.NewHub("test-key")
}
