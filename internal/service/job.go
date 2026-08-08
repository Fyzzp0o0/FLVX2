package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Fyzzp0o0/FLVX2/internal/ws"
)

// JobService 定时任务(对照 ResetFlowAsync / StatisticsFlowAsync)
type JobService struct {
	pool *pgxpool.Pool
	hub  *ws.Hub
}

func NewJobService(pool *pgxpool.Pool, hub *ws.Hub) *JobService {
	return &JobService{pool: pool, hub: hub}
}

// ResetFlowDaily 每天 00:00:05:流量重置 + 到期停服
func (s *JobService) ResetFlowDaily(ctx context.Context) {
	now := time.Now()
	currentDay := now.Day()
	lastDay := daysInMonth(now.Year(), now.Month())
	log.Printf("[job] 流量重置开始: 当月第%d天/最后%d天", currentDay, lastDay)

	// 1) 重置用户流量(月末边界:flow_reset_time>lastDay 也在月末重置)
	cond := resetCondition(currentDay, lastDay)
	rows, err := s.pool.Query(ctx, `SELECT id FROM "user" WHERE flow_reset_time != 0 AND `+cond)
	if err == nil {
		for rows.Next() {
			var id int64
			if rows.Scan(&id) == nil {
				_, _ = s.pool.Exec(ctx, `UPDATE "user" SET in_flow = 0, out_flow = 0 WHERE id = $1`, id)
			}
		}
		rows.Close()
	}
	// 2) 重置用户隧道流量
	rows, err = s.pool.Query(ctx, `SELECT id FROM user_tunnel WHERE flow_reset_time != 0 AND `+cond)
	if err == nil {
		for rows.Next() {
			var id int64
			if rows.Scan(&id) == nil {
				_, _ = s.pool.Exec(ctx, `UPDATE user_tunnel SET in_flow = 0, out_flow = 0 WHERE id = $1`, id)
			}
		}
		rows.Close()
	}
	log.Printf("[job] 流量重置完成")

	// 3) 到期用户停服(user.status=0 + 暂停其转发)
	s.expireUsers(ctx, now)
	// 4) 到期隧道停服(user_tunnel.status=0 + 暂停对应转发)
	s.expireUserTunnels(ctx, now)
}

// expireUsers 到期用户:暂停其全部启用转发 + status=0
func (s *JobService) expireUsers(ctx context.Context, now time.Time) {
	rows, err := s.pool.Query(ctx,
		`SELECT id FROM "user" WHERE role_id != 0 AND status = 1 AND exp_time IS NOT NULL AND exp_time < $1`,
		now.UnixMilli())
	if err != nil {
		return
	}
	var ids []int64
	for rows.Next() {
		var id int64
		if rows.Scan(&id) == nil {
			ids = append(ids, id)
		}
	}
	rows.Close()
	for _, uid := range ids {
		s.pauseUserForwards(ctx, uid, 0)
		_, _ = s.pool.Exec(ctx, `UPDATE "user" SET status = 0 WHERE id = $1`, uid)
		log.Printf("[job] 用户到期停服 id=%d", uid)
	}
}

// expireUserTunnels 到期隧道权限:暂停该用户该隧道转发 + status=0
func (s *JobService) expireUserTunnels(ctx context.Context, now time.Time) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, user_id, tunnel_id FROM user_tunnel WHERE status = 1 AND exp_time IS NOT NULL AND exp_time < $1`,
		now.UnixMilli())
	if err != nil {
		return
	}
	type ut struct{ id, uid, tid int64 }
	var list []ut
	for rows.Next() {
		var u ut
		if rows.Scan(&u.id, &u.uid, &u.tid) == nil {
			list = append(list, u)
		}
	}
	rows.Close()
	for _, u := range list {
		s.pauseUserForwards(ctx, u.uid, u.tid)
		_, _ = s.pool.Exec(ctx, `UPDATE user_tunnel SET status = 0 WHERE id = $1`, u.id)
		log.Printf("[job] 隧道权限到期停服 id=%d", u.id)
	}
}

// pauseUserForwards 暂停用户转发(可选限定隧道);PauseService 下发 + status=0
func (s *JobService) pauseUserForwards(ctx context.Context, userID, tunnelID int64) {
	query := `SELECT id, user_id, user_name, name, tunnel_id FROM forward WHERE user_id = $1 AND status = 1`
	args := []any{userID}
	if tunnelID > 0 {
		query += ` AND tunnel_id = $2`
		args = append(args, tunnelID)
	}
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return
	}
	var forwards []struct {
		ID       int64
		Name     string
		TunnelID int64
	}
	for rows.Next() {
		var f struct {
			ID       int64
			Name     string
			TunnelID int64
		}
		if rows.Scan(&f.ID, new(int64), new(string), &f.Name, &f.TunnelID) == nil {
			forwards = append(forwards, f)
		}
	}
	rows.Close()
	for _, f := range forwards {
		for _, nid := range entryNodeIDs(ctx, s.pool, f.TunnelID) {
			s.hub.SendCommand(nid, "PauseService", ws.PauseServicesData{
				Services: []string{f.Name + "_tcp", f.Name + "_udp"},
			})
		}
		_, _ = s.pool.Exec(ctx, `UPDATE forward SET status = 0 WHERE id = $1`, f.ID)
	}
}

// SnapshotHourly 每小时整点:清理 48h 前数据 + 每用户流量快照
func (s *JobService) SnapshotHourly(ctx context.Context) {
	now := time.Now()
	hourStr := now.Format("15:04")
	cutoff := now.UnixMilli() - 48*60*60*1000

	_, _ = s.pool.Exec(ctx, `DELETE FROM statistics_flow WHERE created_time < $1`, cutoff)

	rows, err := s.pool.Query(ctx, `SELECT id, in_flow, out_flow FROM "user"`)
	if err != nil {
		return
	}
	type u struct {
		id          int64
		total       int64
	}
	var users []u
	for rows.Next() {
		var uu u
		var inF, outF int64
		if rows.Scan(&uu.id, &inF, &outF) == nil {
			uu.total = inF + outF
			users = append(users, uu)
		}
	}
	rows.Close()
	for _, user := range users {
		// 上一条快照
		var lastTotal int64
		_ = s.pool.QueryRow(ctx,
			`SELECT total_flow FROM statistics_flow WHERE user_id = $1 ORDER BY id DESC LIMIT 1`, user.id).Scan(&lastTotal)
		increment := user.total - lastTotal
		if increment < 0 {
			increment = user.total // 负数回退为当前累计(Java 行为)
		}
		_, _ = s.pool.Exec(ctx,
			`INSERT INTO statistics_flow (user_id, flow, total_flow, time, created_time) VALUES ($1,$2,$3,$4,$5)`,
			user.id, increment, user.total, hourStr, now.UnixMilli())
	}
	log.Printf("[job] 统计快照完成: %d 用户, %s", len(users), hourStr)
}

// resetCondition 流量重置 SQL 条件(月末边界)
func resetCondition(currentDay, lastDay int) string {
	if currentDay == lastDay {
		return fmt.Sprintf("(flow_reset_time = %d OR flow_reset_time > %d)", currentDay, lastDay)
	}
	return fmt.Sprintf("flow_reset_time = %d", currentDay)
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
