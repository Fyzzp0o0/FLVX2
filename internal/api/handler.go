package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Fyzzp0o0/FLVX2/internal/jwt"
	"github.com/Fyzzp0o0/FLVX2/internal/service"
	"github.com/Fyzzp0o0/FLVX2/internal/ws"
)

// Handler 聚合各业务 handler
type Handler struct {
	users      *service.UserService
	nodes      *service.NodeService
	tunnels    *service.TunnelService
	forwards   *service.ForwardService
	speedLimits *service.SpeedLimitService
	cfg        *service.ViteConfigService
	flows      *service.FlowService
	hub        *ws.Hub
	secret     string
}

func NewHandler(users *service.UserService, nodes *service.NodeService, tunnels *service.TunnelService,
	forwards *service.ForwardService, speedLimits *service.SpeedLimitService, cfg *service.ViteConfigService,
	flows *service.FlowService, hub *ws.Hub, secret string) *Handler {
	return &Handler{
		users: users, nodes: nodes, tunnels: tunnels, forwards: forwards,
		speedLimits: speedLimits, cfg: cfg, flows: flows, hub: hub, secret: secret,
	}
}

// ---- UserController /api/v1/user ----

type loginReq struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	CaptchaID string `json:"captchaId"`
}

// Login 登录(免鉴权白名单)
func (h *Handler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	u, requireChange, err := h.users.Login(c.Request.Context(), req.Username, req.Password, req.CaptchaID)
	if err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	token, err := jwt.GenerateToken(h.secret, u.ID, u.User, u.RoleID)
	if err != nil {
		c.JSON(http.StatusOK, errCode(-2, "系统异常"))
		return
	}
	c.JSON(http.StatusOK, ok(gin.H{
		"token":                 token,
		"name":                  u.User,
		"role_id":               u.RoleID,
		"requirePasswordChange": requireChange,
	}))
}

type registerReq struct {
	User      string `json:"user" binding:"required,min=3,max=20"`
	Pwd       string `json:"pwd" binding:"required,min=6,max=32"`
	CaptchaID string `json:"captchaId"`
}

// Register 注册(免鉴权白名单)
func (h *Handler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	if err := h.users.Register(c.Request.Context(), req.User, req.Pwd); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

type createUserReq struct {
	User          string `json:"user" binding:"required"`
	Pwd           string `json:"pwd" binding:"required"`
	Flow          int64  `json:"flow" binding:"required"`
	Num           int64  `json:"num" binding:"required"`
	ExpTime       int64  `json:"expTime" binding:"required"`
	FlowResetTime int64  `json:"flowResetTime" binding:"required"`
	Status        *int64 `json:"status"`
}

// UserCreate 管理员建用户
func (h *Handler) UserCreate(c *gin.Context) {
	var req createUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	if err := h.users.CreateUser(c.Request.Context(), req.User, req.Pwd, req.Flow, req.Num, req.ExpTime, req.FlowResetTime); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// UserList 用户列表(管理员)
func (h *Handler) UserList(c *gin.Context) {
	list, err := h.users.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(list))
}

type updateUserReq struct {
	ID            int64  `json:"id" binding:"required"`
	User          string `json:"user" binding:"required"`
	Pwd           string `json:"pwd"`
	Flow          int64  `json:"flow" binding:"required"`
	Num           int64  `json:"num" binding:"required"`
	ExpTime       int64  `json:"expTime" binding:"required"`
	FlowResetTime int64  `json:"flowResetTime" binding:"required"`
	Status        *int64 `json:"status"`
}

// UserUpdate 更新用户(管理员)
func (h *Handler) UserUpdate(c *gin.Context) {
	var req updateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	if err := h.users.UpdateUser(c.Request.Context(), req.ID, req.User, req.Pwd, req.Flow, req.Num, req.ExpTime, req.FlowResetTime); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// UserDelete 删除用户(管理员)
func (h *Handler) UserDelete(c *gin.Context) {
	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	if err := h.users.DeleteUser(c.Request.Context(), req.ID); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// UserReset 重置流量(管理员;type==1 用户,否则 user_tunnel)
func (h *Handler) UserReset(c *gin.Context) {
	var req struct {
		ID   int64 `json:"id" binding:"required"`
		Type int   `json:"type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	if err := h.users.ResetFlow(c.Request.Context(), req.ID, req.Type); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// UserPackage 当前用户套餐
func (h *Handler) UserPackage(c *gin.Context) {
	uid, _, _ := currentUser(c)
	pkg, err := h.users.Package(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(pkg))
}

type changePasswordReq struct {
	NewUsername      string `json:"newUsername" binding:"required"`
	CurrentPassword  string `json:"currentPassword" binding:"required"`
	NewPassword      string `json:"newPassword" binding:"required"`
	ConfirmPassword  string `json:"confirmPassword" binding:"required"`
}

// UserUpdatePassword 修改用户名+密码
func (h *Handler) UserUpdatePassword(c *gin.Context) {
	var req changePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	if req.NewPassword != req.ConfirmPassword {
		c.JSON(http.StatusOK, errMsg("新密码和确认密码不匹配"))
		return
	}
	uid, _, _ := currentUser(c)
	if err := h.users.UpdatePassword(c.Request.Context(), uid, req.NewUsername, req.CurrentPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// CaptchaCheck /api/v1/captcha/check:captcha_enabled=="true" ? 1 : 0
func (h *Handler) CaptchaCheck(c *gin.Context) {
	ctx := context.Background()
	enabled, _ := h.cfg.Get(ctx, "captcha_enabled")
	data := 0
	if enabled == "true" {
		data = 1
	}
	c.JSON(http.StatusOK, ok(data))
}

// ---- FlowController /flow ----

// FlowTest /flow/test 探活,返回纯文本 "test"(安装脚本健康检查用)
func (h *Handler) FlowTest(c *gin.Context) {
	c.String(http.StatusOK, "test")
}

// FlowUpload /flow/upload 节点流量上报(免鉴权,按 secret 认节点;响应必须精确 "ok")
func (h *Handler) FlowUpload(c *gin.Context) {
	secret := c.Query("secret")
	body, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusOK, "ok")
		return
	}
	c.String(http.StatusOK, h.flows.ProcessUpload(c.Request.Context(), secret, body))
}

// FlowConfig /flow/config 节点配置快照上报(免鉴权;孤儿清理在 M4 与 Gost 联动时实现)
func (h *Handler) FlowConfig(c *gin.Context) {
	secret := c.Query("secret")
	body, err := c.GetRawData()
	if err == nil && secret != "" {
		_, _ = h.flows.ParseConfigSnapshot(c.Request.Context(), secret, body)
	}
	c.String(http.StatusOK, "ok")
}

// SystemInfo /system-info WebSocket(type=1 节点 / type=0 管理员)
func (h *Handler) SystemInfo(c *gin.Context) {
	h.hub.HandleWS(c)
}

// ---- 通用 ----

// NotImplemented 未实现端点(当前由 NoRoute 兜底)
func (h *Handler) NotImplemented(c *gin.Context) {
	path := c.Request.URL.Path
	if strings.HasPrefix(path, "/flow/") {
		c.String(http.StatusOK, "ok")
		return
	}
	c.JSON(http.StatusOK, errMsg("接口尚未实现"))
}
