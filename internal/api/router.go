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

	// ---- /api/v1 ----
	apiV1 := r.Group("/api/v1")
	{
		user := apiV1.Group("/user")
		{
			user.POST("/login", h.Login)
			user.POST("/register", h.Register)
		}
		captcha := apiV1.Group("/captcha")
		{
			captcha.POST("/check", h.CaptchaCheck)
			// TODO(M3): /generate /verify 自研滑块验证码
		}
		// 后续里程碑端点(节点/隧道/转发/限速/配置/套餐/改密等)由 NoRoute 兜底
	}

	// ---- /flow(agent 上报,免鉴权) ----
	flow := r.Group("/flow")
	{
		flow.Any("/test", h.FlowTest)
		// TODO(M2): /upload /config 流量上报与配置快照
	}

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
