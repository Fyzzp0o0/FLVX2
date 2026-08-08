package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Fyzzp0o0/FLVX2/internal/service"
)

// Handler 聚合各业务 handler
type Handler struct {
	users *service.UserService
	secret string
}

func NewHandler(users *service.UserService, secret string) *Handler {
	return &Handler{users: users, secret: secret}
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
		switch {
		case errors.Is(err, service.ErrCaptchaRequired):
			c.JSON(http.StatusOK, errMsg(err.Error()))
		default:
			c.JSON(http.StatusOK, errMsg(err.Error()))
		}
		return
	}
	token, err := GenerateToken(h.secret, u.ID, u.User, u.RoleID)
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

// CaptchaCheck /api/v1/captcha/check:captcha_enabled=="true" ? 1 : 0
func (h *Handler) CaptchaCheck(c *gin.Context) {
	ctx := context.Background()
	enabled, _ := h.users.GetConfig(ctx, "captcha_enabled")
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

// ---- 未实现端点占位(后续里程碑) ----

// NotImplemented 后续里程碑端点,先返回业务错误避免前端 404 语义不一致
func (h *Handler) NotImplemented(c *gin.Context) {
	path := c.Request.URL.Path
	if strings.HasPrefix(path, "/flow/") {
		c.String(http.StatusOK, "ok")
		return
	}
	c.JSON(http.StatusOK, errMsg("接口尚未实现"))
}
