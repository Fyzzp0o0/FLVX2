package service

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Fyzzp0o0/FLVX2/internal/model"
	"github.com/Fyzzp0o0/FLVX2/internal/ws"
)

// md5Hex 无盐 MD5 十六进制小写(与原 Md5Util.md5 一致)
func md5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", sum)
}

const (
	roleAdmin  = 0
	roleNormal = 1

	statusEnabled  = 1
	statusDisabled = 0

	bytesPerGB = 1024 * 1024 * 1024

	defaultRegisterFlow      = int64(100) // GB
	defaultRegisterNum       = int64(5)
	defaultRegisterExpDays   = int64(30)
	defaultRegisterFlowReset = int64(30) // 天(原代码写为 now+30天 时间戳,见 flow_reset_time 现状说明)
)

var (
	ErrUserNotFound    = errors.New("账号或密码错误")
	ErrUserDisabled    = errors.New("账号被停用")
	ErrUserExists      = errors.New("用户名已存在")
	ErrCaptchaRequired = errors.New("验证码功能尚未启用")
	ErrCannotModifyAdmin = errors.New("请不要作死")
)

// UserService 用户相关业务(对照 UserServiceImpl)
type UserService struct {
	pool *pgxpool.Pool
	cfg  *ViteConfigService
	hub  *ws.Hub
}

func NewUserService(pool *pgxpool.Pool, cfg *ViteConfigService, hub *ws.Hub) *UserService {
	return &UserService{pool: pool, cfg: cfg, hub: hub}
}

// GetByUsername 按用户名查用户("user" 为保留字必须双引号)
func (s *UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT id, "user", pwd, role_id, exp_time, flow, in_flow, out_flow, flow_reset_time, num, created_time, updated_time, status
		 FROM "user" WHERE "user" = $1`, username)
	return scanUser(row)
}

// GetByID 按 ID 查用户
func (s *UserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT id, "user", pwd, role_id, exp_time, flow, in_flow, out_flow, flow_reset_time, num, created_time, updated_time, status
		 FROM "user" WHERE id = $1`, id)
	return scanUser(row)
}

func scanUser(row interface{ Scan(...any) error }) (*model.User, error) {
	u := &model.User{}
	err := row.Scan(&u.ID, &u.User, &u.Pwd, &u.RoleID, &u.ExpTime, &u.Flow, &u.InFlow, &u.OutFlow,
		&u.FlowResetTime, &u.Num, &u.CreatedTime, &u.UpdatedTime, &u.Status)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Login 登录校验(对照 UserServiceImpl.login):返回 user + 是否需强制改密
func (s *UserService) Login(ctx context.Context, username, password, captchaID string) (*model.User, bool, error) {
	// 验证码开关(captcha_enabled=="true" 需二次校验;M3 阶段验证码未实现,暂返回未启用错误)
	captchaEnabled, _ := s.cfg.Get(ctx, "captcha_enabled")
	if captchaEnabled == "true" {
		if captchaID == "" {
			return nil, false, ErrCaptchaRequired
		}
		// TODO(M3 后期): 接入自研滑块验证码 secondaryVerification 校验
		return nil, false, ErrCaptchaRequired
	}
	u, err := s.GetByUsername(ctx, username)
	if err != nil {
		return nil, false, ErrUserNotFound
	}
	if u.Pwd != md5Hex(password) { // 无盐 MD5,已实测与 data.sql 内置哈希一致
		return nil, false, ErrUserNotFound
	}
	if u.Status == statusDisabled {
		return nil, false, ErrUserDisabled
	}
	requireChange := username == "admin_user" || password == "admin_user"
	return u, requireChange, nil
}

// Register 注册(对照 UserServiceImpl.register):默认配额取自 vite_config,可被覆盖
func (s *UserService) Register(ctx context.Context, username, password string) error {
	if err := s.checkUsernameAvailable(ctx, username); err != nil {
		return err
	}
	v1, _ := s.cfg.Get(ctx, "register_default_flow")
	v2, _ := s.cfg.Get(ctx, "register_default_num")
	v3, _ := s.cfg.Get(ctx, "register_default_exp_days")
	v4, _ := s.cfg.Get(ctx, "register_default_flow_reset_days")
	flow, _ := parseConfig(v1, defaultRegisterFlow)
	num, _ := parseConfig(v2, defaultRegisterNum)
	expDays, _ := parseConfig(v3, defaultRegisterExpDays)
	resetDays, _ := parseConfig(v4, defaultRegisterFlowReset)

	now := time.Now().UnixMilli()
	_, err := s.pool.Exec(ctx,
		`INSERT INTO "user" ("user", pwd, role_id, exp_time, flow, in_flow, out_flow, flow_reset_time, num, created_time, updated_time, status)
		 VALUES ($1, $2, $3, $4, $5, 0, 0, $6, $7, $8, $8, 1)`,
		username, md5Hex(password), roleNormal,
		now+expDays*24*3600*1000, flow,
		now+resetDays*24*3600*1000, // 现状:写为毫秒时间戳(见文档 flow_reset_time 决策点)
		num, now)
	if err != nil {
		return err
	}
	return nil
}

// CreateUser 管理员建用户(固定 roleId=1,status=1;原 createUser 固定 status=1,入参 status 忽略)
func (s *UserService) CreateUser(ctx context.Context, username, password string, flow, num, expTime, flowResetTime int64) error {
	if err := s.checkUsernameAvailable(ctx, username); err != nil {
		return err
	}
	now := time.Now().UnixMilli()
	_, err := s.pool.Exec(ctx,
		`INSERT INTO "user" ("user", pwd, role_id, exp_time, flow, in_flow, out_flow, flow_reset_time, num, created_time, updated_time, status)
		 VALUES ($1, $2, $3, $4, $5, 0, 0, $6, $7, $8, $8, 1)`,
		username, md5Hex(password), roleNormal, expTime, flow, flowResetTime, num, now)
	return err
}

// ListUsers 全部普通用户(role_id != 0,不含密码)
func (s *UserService) ListUsers(ctx context.Context) ([]*model.User, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, "user", pwd, role_id, exp_time, flow, in_flow, out_flow, flow_reset_time, num, created_time, updated_time, status
		 FROM "user" WHERE role_id != 0 ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*model.User{}
	for rows.Next() {
		u, err := scanUser(rows)
		if err == nil {
			u.Pwd = ""
			out = append(out, u)
		}
	}
	return out, nil
}

// UpdateUser 管理员更新用户(pwd 为空不更新密码)
func (s *UserService) UpdateUser(ctx context.Context, id int64, username, pwd string, flow, num, expTime, flowResetTime int64) error {
	u, err := s.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}
	if u.RoleID == roleAdmin {
		return ErrCannotModifyAdmin // "请不要作死"
	}
	// 改用户名查重(排除自身)
	var cnt int64
	if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM "user" WHERE "user" = $1 AND id != $2`, username, id).Scan(&cnt); err != nil {
		return err
	}
	if cnt > 0 {
		return ErrUserExists
	}
	now := time.Now().UnixMilli()
	if pwd != "" {
		_, err = s.pool.Exec(ctx,
			`UPDATE "user" SET "user" = $1, pwd = $2, flow = $3, num = $4, exp_time = $5, flow_reset_time = $6, updated_time = $7 WHERE id = $8`,
			username, md5Hex(pwd), flow, num, expTime, flowResetTime, now, id)
	} else {
		_, err = s.pool.Exec(ctx,
			`UPDATE "user" SET "user" = $1, flow = $2, num = $3, exp_time = $4, flow_reset_time = $5, updated_time = $6 WHERE id = $7`,
			username, flow, num, expTime, flowResetTime, now, id)
	}
	return err
}

// DeleteUser 删除用户(级联:删 Gost 服务 → forward/user_tunnel/statistics_flow/user)
func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	u, err := s.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}
	if u.RoleID == roleAdmin {
		return ErrCannotModifyAdmin
	}
	// 先删该用户全部转发的 Gost 服务(入口节点 DeleteService)
	forwards, _ := s.listForwardNames(ctx, id)
	for _, f := range forwards {
		nodeIDs := entryNodeIDs(ctx, s.pool, f.TunnelID)
		for _, nid := range nodeIDs {
			s.hub.SendCommand(nid, "DeleteService", ws.DeleteServicesData{
				Services: []string{f.Name + "_tcp", f.Name + "_udp"},
			})
		}
	}
	_, err = s.pool.Exec(ctx, `DELETE FROM forward WHERE user_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `DELETE FROM user_tunnel WHERE user_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `DELETE FROM statistics_flow WHERE user_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `DELETE FROM "user" WHERE id = $1`, id)
	return err
}

// ResetFlow 流量重置(type==1 清用户,否则清 user_tunnel)
func (s *UserService) ResetFlow(ctx context.Context, id int64, resetType int) error {
	if resetType == 1 {
		_, err := s.pool.Exec(ctx, `UPDATE "user" SET in_flow = 0, out_flow = 0 WHERE id = $1`, id)
		return err
	}
	_, err := s.pool.Exec(ctx, `UPDATE user_tunnel SET in_flow = 0, out_flow = 0 WHERE id = $1`, id)
	return err
}

// UpdatePassword 修改用户名+密码
func (s *UserService) UpdatePassword(ctx context.Context, userID int64, newUsername, currentPassword, newPassword string) error {
	u, err := s.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}
	if u.Pwd != md5Hex(currentPassword) {
		return errors.New("当前密码错误")
	}
	if newUsername != u.User {
		var cnt int64
		if err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM "user" WHERE "user" = $1 AND id != $2`, newUsername, userID).Scan(&cnt); err != nil {
			return err
		}
		if cnt > 0 {
			return ErrUserExists
		}
	}
	_, err = s.pool.Exec(ctx,
		`UPDATE "user" SET "user" = $1, pwd = $2, updated_time = $3 WHERE id = $4`,
		newUsername, md5Hex(newPassword), time.Now().UnixMilli(), userID)
	return err
}

func (s *UserService) checkUsernameAvailable(ctx context.Context, username string) error {
	var cnt int64
	err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM "user" WHERE "user" = $1`, username).Scan(&cnt)
	if err != nil {
		return err
	}
	if cnt > 0 {
		return ErrUserExists
	}
	return nil
}

// forwardName 用户转发的基础名+隧道(级联删除用)
type forwardName struct {
	Name     string
	TunnelID int64
}

func (s *UserService) listForwardNames(ctx context.Context, userID int64) ([]forwardName, error) {
	rows, err := s.pool.Query(ctx, `SELECT name, tunnel_id FROM forward WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []forwardName{}
	for rows.Next() {
		var f forwardName
		if err := rows.Scan(&f.Name, &f.TunnelID); err == nil {
			out = append(out, f)
		}
	}
	return out, nil
}

// entryNodeIDs 隧道入口节点(chain_type='1')
func entryNodeIDs(ctx context.Context, pool *pgxpool.Pool, tunnelID int64) []int64 {
	ids := []int64{}
	rows, err := pool.Query(ctx, `SELECT node_id FROM chain_tunnel WHERE tunnel_id = $1 AND chain_type = '1'`, tunnelID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

func parseConfig(v string, def int64) (int64, error) {
	if v == "" {
		return def, nil
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return def, err
	}
	return n, nil
}

// ---- 套餐(/api/v1/user/package) ----

// UserPackage 用户套餐视图
type UserPackage struct {
	UserInfo          *model.User          `json:"userInfo"`
	TunnelPermissions []*UserTunnelDetail  `json:"tunnelPermissions"`
	Forwards          []*ForwardWithTunnel `json:"forwards"`
	StatisticsFlows   []*model.StatisticsFlow `json:"statisticsFlows"`
}

// Package 当前用户套餐(普通用户/管理员均可)
func (s *UserService) Package(ctx context.Context, userID int64) (*UserPackage, error) {
	u, err := s.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	u.Pwd = ""
	pkg := &UserPackage{UserInfo: u, TunnelPermissions: []*UserTunnelDetail{}, Forwards: []*ForwardWithTunnel{}, StatisticsFlows: []*model.StatisticsFlow{}}

	// 隧道权限
	rows, err := s.pool.Query(ctx,
		`SELECT ut.id, ut.user_id, ut.tunnel_id, ut.flow, ut.num, ut.flow_reset_time, ut.exp_time, ut.speed_id,
		        COALESCE(sl.name,''), COALESCE(sl.speed,0), COALESCE(t.name,''), COALESCE(t.flow,0),
		        ut.in_flow, ut.out_flow, ut.status
		 FROM user_tunnel ut
		 LEFT JOIN tunnel t ON ut.tunnel_id = t.id
		 LEFT JOIN speed_limit sl ON ut.speed_id = sl.id
		 WHERE ut.user_id = $1 ORDER BY ut.id`, userID)
	if err == nil {
		pkg.TunnelPermissions, _ = scanUserTunnelDetails(rows)
		rows.Close()
	}

	// 转发(带隧道 inIp)
	rows, err = s.pool.Query(ctx,
		`SELECT f.id, f.name, f.tunnel_id, COALESCE(t.in_ip,''), COALESCE(t.name,''), f.remote_addr, f.strategy,
		        f.status, f.created_time, f.updated_time, f.user_id, f.user_name, f.in_flow, f.out_flow, f.inx
		 FROM forward f LEFT JOIN tunnel t ON f.tunnel_id = t.id WHERE f.user_id = $1 ORDER BY f.created_time DESC`, userID)
	if err == nil {
		for rows.Next() {
			f := &ForwardWithTunnel{}
			if rows.Scan(&f.ID, &f.Name, &f.TunnelID, &f.InIP, &f.TunnelName, &f.RemoteAddr, &f.Strategy,
				&f.Status, &f.CreatedTime, &f.UpdatedTime, &f.UserID, &f.UserName, &f.InFlow, &f.OutFlow, &f.Inx) == nil {
				// 入口端口
				_ = s.pool.QueryRow(ctx, `SELECT port FROM forward_port WHERE forward_id = $1 ORDER BY id LIMIT 1`, f.ID).Scan(&f.InPort)
				pkg.Forwards = append(pkg.Forwards, f)
			}
		}
		rows.Close()
	}

	// 近 24h 统计(不足补零)
	rows, err = s.pool.Query(ctx,
		`SELECT id, user_id, flow, total_flow, time, created_time FROM statistics_flow WHERE user_id = $1 ORDER BY id DESC LIMIT 24`, userID)
	if err == nil {
		list := []*model.StatisticsFlow{}
		for rows.Next() {
			st := &model.StatisticsFlow{}
			if rows.Scan(&st.ID, &st.UserID, &st.Flow, &st.TotalFlow, &st.Time, &st.CreatedTime) == nil {
				list = append(list, st)
			}
		}
		rows.Close()
		// 倒序(时间从旧到新),不足 24 条补零
		for i := len(list) - 1; i >= 0; i-- {
			pkg.StatisticsFlows = append(pkg.StatisticsFlows, list[i])
		}
		for len(pkg.StatisticsFlows) < 24 {
			pkg.StatisticsFlows = append([]*model.StatisticsFlow{{Flow: 0, TotalFlow: 0}}, pkg.StatisticsFlows...)
		}
	}
	return pkg, nil
}

// QueryUserTunnel 查用户-隧道权限(open_api 用)
func (s *UserService) QueryUserTunnel(ctx context.Context, userID, tunnelID int64) (*model.UserTunnel, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT id, user_id, tunnel_id, speed_id, num, flow, in_flow, out_flow, flow_reset_time, exp_time, status
		 FROM user_tunnel WHERE user_id = $1 AND tunnel_id = $2`, userID, tunnelID)
	ut := &model.UserTunnel{}
	err := row.Scan(&ut.ID, &ut.UserID, &ut.TunnelID, &ut.SpeedID, &ut.Num, &ut.Flow, &ut.InFlow, &ut.OutFlow, &ut.FlowResetTime, &ut.ExpTime, &ut.Status)
	if err != nil {
		return nil, err
	}
	return ut, nil
}
