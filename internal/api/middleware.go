package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Fyzzp0o0/FLVX2/internal/jwt"
)

// R 统一响应包装,兼容原 R.java: {code, msg, ts, data}
type R struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Ts   int64       `json:"ts"`
	Data interface{} `json:"data"`
}

func ok(data interface{}) R { return R{Code: 0, Msg: "success", Ts: time.Now().UnixMilli(), Data: data} }

func errMsg(msg string) R { return R{Code: -1, Msg: msg, Ts: time.Now().UnixMilli(), Data: nil} }

func errCode(code int, msg string) R { return R{Code: code, Msg: msg, Ts: time.Now().UnixMilli(), Data: nil} }

// JWT 白名单(对应原 WebMvcConfig 豁免路径)
var jwtWhitelist = []string{
	"/api/v1/open_api/",
	"/api/v1/config/get",
	"/api/v1/user/login",
	"/api/v1/user/register",
	"/api/v1/captcha/",
	"/flow/",
}

func isWhitelisted(path string) bool {
	for _, p := range jwtWhitelist {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

// CORS 中间件,对应原 WebMvcConfig:允许所有来源/头/方法,暴露 Authorization
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

// JWTAuth 拦截 /api/** 除白名单外的请求(裸 token,无 Bearer 前缀)
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if !strings.HasPrefix(path, "/api/") || isWhitelisted(path) {
			c.Next()
			return
		}
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusOK, errCode(401, "未登录或token已过期"))
			c.Abort()
			return
		}
		uid, username, roleID, valid := jwt.ValidateToken(secret, token)
		if !valid {
			c.JSON(http.StatusOK, errCode(401, "无效的token或token已过期"))
			c.Abort()
			return
		}
		c.Set("userId", uid)
		c.Set("username", username)
		c.Set("roleId", roleID)
		c.Next()
	}
}

// RequireAdmin 对应 @RequireRole:role_id == 0(管理员)
func RequireAdmin(c *gin.Context) {
	roleID, exists := c.Get("roleId")
	if !exists {
		c.JSON(http.StatusOK, errCode(401, "未登录或token已过期"))
		c.Abort()
		return
	}
	if roleID.(int64) != 0 {
		c.JSON(http.StatusOK, errCode(403, "权限不足，仅管理员可操作"))
		c.Abort()
		return
	}
	c.Next()
}

// 从上下文取当前登录用户
func currentUser(c *gin.Context) (userID int64, username string, roleID int64) {
	uid, _ := c.Get("userId")
	name, _ := c.Get("username")
	rid, _ := c.Get("roleId")
	userID, _ = uid.(int64)
	username, _ = name.(string)
	roleID, _ = rid.(int64)
	return
}
