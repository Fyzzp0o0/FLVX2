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

	defaultRegisterFlow        = int64(100) // GB
	defaultRegisterNum         = int64(5)
	defaultRegisterExpDays     = int64(30)
	defaultRegisterFlowReset   = int64(30) // 天(原代码写为 now+30天 时间戳,见 flow_reset_time 现状说明)
	adminExpTimeFarFuture      = int64(2727251700000)
	adminFlowUnlimited         = int64(99999)
)

var (
	ErrUserNotFound      = errors.New("账号或密码错误")
	ErrUserDisabled      = errors.New("账号被停用")
	ErrUserExists        = errors.New("用户名已存在")
	ErrCaptchaRequired   = errors.New("验证码功能尚未启用")
)

// UserService 用户相关业务(对照 UserServiceImpl)
type UserService struct {
	pool *pgxpool.Pool
}

func NewUserService(pool *pgxpool.Pool) *UserService { return &UserService{pool: pool} }

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

// GetConfig 读 vite_config 单条配置(name 不存在返回空串)
func (s *UserService) GetConfig(ctx context.Context, name string) (string, error) {
	var value string
	err := s.pool.QueryRow(ctx, `SELECT value FROM vite_config WHERE name = $1`, name).Scan(&value)
	if err != nil {
		return "", nil // 不存在视为空
	}
	return value, nil
}

// Login 登录校验(对照 UserServiceImpl.login):返回 user + 是否需强制改密
func (s *UserService) Login(ctx context.Context, username, password, captchaID string) (*model.User, bool, error) {
	// 验证码开关(captcha_enabled=="true" 需二次校验;M1 阶段验证码未实现,暂返回未启用错误)
	captchaEnabled, _ := s.GetConfig(ctx, "captcha_enabled")
	if captchaEnabled == "true" {
		if captchaID == "" {
			return nil, false, ErrCaptchaRequired
		}
		// TODO(M3): 接入自研滑块验证码 secondaryVerification 校验
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
	v1, _ := s.GetConfig(ctx, "register_default_flow")
	v2, _ := s.GetConfig(ctx, "register_default_num")
	v3, _ := s.GetConfig(ctx, "register_default_exp_days")
	v4, _ := s.GetConfig(ctx, "register_default_flow_reset_days")
	flow, _ := parseGBConfig(v1, defaultRegisterFlow)
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

func parseGBConfig(v string, def int64) (int64, error) { return parseConfig(v, def) }
