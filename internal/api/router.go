package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Fyzzp0o0/FLVX2/internal/web"
)

// NewRouter 组装全部路由(6635 前端 / 6636 后端共享同一引擎与中间件)
func NewRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), CORS(), JWTAuth(h.secret))

	apiV1 := r.Group("/api/v1")
	{
		// ---- 用户(login/register/captcha 白名单免 JWT) ----
		user := apiV1.Group("/user")
		{
			user.POST("/login", h.Login)
			user.POST("/register", h.Register)
			user.POST("/create", RequireAdmin, h.UserCreate)
			user.POST("/list", RequireAdmin, h.UserList)
			user.POST("/update", RequireAdmin, h.UserUpdate)
			user.POST("/delete", RequireAdmin, h.UserDelete)
			user.POST("/reset", RequireAdmin, h.UserReset)
			user.POST("/package", h.UserPackage)
			user.POST("/updatePassword", h.UserUpdatePassword)
		}
		captcha := apiV1.Group("/captcha")
		{
			captcha.POST("/check", h.CaptchaCheck)
			// TODO(M3 后期): /generate /verify 自研滑块验证码
		}
		// ---- 节点(全部管理员) ----
		node := apiV1.Group("/node", RequireAdmin)
		{
			node.POST("/create", h.NodeCreate)
			node.POST("/list", h.NodeList)
			node.POST("/update", h.NodeUpdate)
			node.POST("/delete", h.NodeDelete)
			node.POST("/install", h.NodeInstall)
		}
		// ---- 隧道 ----
		tunnel := apiV1.Group("/tunnel")
		{
			tunnel.POST("/create", RequireAdmin, h.TunnelCreate)
			tunnel.POST("/list", RequireAdmin, h.TunnelList)
			tunnel.POST("/update", RequireAdmin, h.TunnelUpdate)
			tunnel.POST("/delete", RequireAdmin, h.TunnelDelete)
			tunnel.POST("/diagnose", RequireAdmin, h.TunnelDiagnose)
			tunnelUser := tunnel.Group("/user")
			{
				tunnelUser.POST("/assign", RequireAdmin, h.TunnelUserAssign)
				tunnelUser.POST("/list", RequireAdmin, h.TunnelUserList)
				tunnelUser.POST("/remove", RequireAdmin, h.TunnelUserRemove)
				tunnelUser.POST("/update", RequireAdmin, h.TunnelUserUpdate)
				tunnelUser.POST("/tunnel", h.TunnelUserTunnel) // 仅 JWT,无 @RequireRole
			}
		}
		// ---- 转发(仅 JWT,角色分支在 service) ----
		forward := apiV1.Group("/forward")
		{
			forward.POST("/create", h.ForwardCreate)
			forward.POST("/list", h.ForwardList)
			forward.POST("/update", h.ForwardUpdate)
			forward.POST("/delete", h.ForwardDelete)
			forward.POST("/force-delete", h.ForwardForceDelete)
			forward.POST("/pause", h.ForwardPause)
			forward.POST("/resume", h.ForwardResume)
			forward.POST("/update-order", h.ForwardUpdateOrder)
			forward.POST("/diagnose", h.ForwardDiagnose)
		}
		// ---- 限速(全部管理员) ----
		speedLimit := apiV1.Group("/speed-limit", RequireAdmin)
		{
			speedLimit.POST("/create", h.SpeedLimitCreate)
			speedLimit.POST("/list", h.SpeedLimitList)
			speedLimit.POST("/update", h.SpeedLimitUpdate)
			speedLimit.POST("/delete", h.SpeedLimitDelete)
			speedLimit.POST("/tunnels", h.SpeedLimitTunnels)
		}
		// ---- 配置 ----
		config := apiV1.Group("/config")
		{
			config.POST("/list", h.ConfigList)           // 需 JWT(不在白名单)
			config.POST("/get", h.ConfigGet)             // 免鉴权白名单
			config.POST("/update", RequireAdmin, h.ConfigUpdate)
			config.POST("/update-single", RequireAdmin, h.ConfigUpdateSingle)
		}
		// ---- open_api(免鉴权,第三方订阅) ----
		openAPI := apiV1.Group("/open_api")
		{
			openAPI.GET("/sub_store", h.SubStore)
		}
	}

	// ---- /flow(agent 上报,免鉴权) ----
	flow := r.Group("/flow")
	{
		flow.Any("/test", h.FlowTest)
		flow.Any("/upload", h.FlowUpload)
		flow.Any("/config", h.FlowConfig)
	}

	// ---- /system-info WebSocket(节点 type=1 / 管理员 type=0) ----
	r.GET("/system-info", h.SystemInfo)

	// ---- 静态托管(仅前端端口 6635 使用,挂同引擎方便同源 API) ----
	registerStatic(r)

	return r
}

// registerStatic 前端静态资源 + SPA fallback(排除 /api /flow /system-info 前缀)
func registerStatic(r *gin.Engine) {
	r.GET("/", web.ServeIndex)
	r.GET("/assets/*filepath", web.ServeAsset)
	r.GET("/favicon.ico", web.ServeIndex)
	// 其余未知路径回退 index.html(等价 nginx try_files;注意排除 API/WS 前缀)
	r.NoRoute(func(c *gin.Context) {
		p := c.Request.URL.Path
		switch {
		case strings.HasPrefix(p, "/flow/"):
			// 未实现的 flow 端点返回 "ok",兼容 agent 上报预期
			c.String(http.StatusOK, "ok")
		case strings.HasPrefix(p, "/api/"), strings.HasPrefix(p, "/system-info"):
			c.JSON(http.StatusNotFound, gin.H{"code": -2, "msg": "not found", "ts": 0, "data": nil})
		default:
			web.ServeIndex(c)
		}
	})
}
