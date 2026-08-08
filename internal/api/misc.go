package api

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// md5HexForAPI 无盐 MD5(与 service 一致)
func md5HexForAPI(s string) string {
	sum := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", sum)
}

// ---- SpeedLimitController /api/v1/speed-limit(全部管理员) ----

type createSpeedLimitReq struct {
	Name       string `json:"name" binding:"required"`
	Speed      int64  `json:"speed" binding:"required"`
	TunnelID   int64  `json:"tunnelId" binding:"required"`
	TunnelName string `json:"tunnelName" binding:"required"`
}

// SpeedLimitCreate 创建限速规则
func (h *Handler) SpeedLimitCreate(c *gin.Context) {
	var req createSpeedLimitReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	if err := h.speedLimits.Create(c.Request.Context(), req.Name, req.Speed, req.TunnelID, req.TunnelName); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// SpeedLimitList 限速规则列表
func (h *Handler) SpeedLimitList(c *gin.Context) {
	list, err := h.speedLimits.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(list))
}

type updateSpeedLimitReq struct {
	ID    int64  `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Speed int64  `json:"speed" binding:"required"`
}

// SpeedLimitUpdate 更新限速
func (h *Handler) SpeedLimitUpdate(c *gin.Context) {
	var req updateSpeedLimitReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	if err := h.speedLimits.Update(c.Request.Context(), req.ID, req.Name, req.Speed); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// SpeedLimitDelete 删除限速(有引用拒绝)
func (h *Handler) SpeedLimitDelete(c *gin.Context) {
	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	if err := h.speedLimits.Delete(c.Request.Context(), req.ID); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// SpeedLimitTunnels 隧道列表(前端选隧道用)
func (h *Handler) SpeedLimitTunnels(c *gin.Context) {
	list, err := h.speedLimits.Tunnels(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(list))
}

// ---- ViteConfigController /api/v1/config ----

// ConfigList 全部配置(需 JWT;不在白名单)
func (h *Handler) ConfigList(c *gin.Context) {
	m, err := h.cfg.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(m))
}

// ConfigGet 单条配置(免鉴权白名单)
func (h *Handler) ConfigGet(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	vc, err := h.cfg.GetModel(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusOK, errMsg("配置不存在"))
		return
	}
	c.JSON(http.StatusOK, ok(vc))
}

// ConfigUpdate 批量更新(管理员)
func (h *Handler) ConfigUpdate(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	for k, v := range req {
		if err := h.cfg.Upsert(c.Request.Context(), k, v); err != nil {
			c.JSON(http.StatusOK, errMsg(err.Error()))
			return
		}
	}
	c.JSON(http.StatusOK, ok(nil))
}

// ConfigUpdateSingle 单条更新(管理员)
func (h *Handler) ConfigUpdateSingle(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	if err := h.cfg.Upsert(c.Request.Context(), req.Name, req.Value); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// ---- OpenApiController /api/v1/open_api(免鉴权,第三方订阅) ----

// SubStore 订阅信息(原 GET,参数在 query)
func (h *Handler) SubStore(c *gin.Context) {
	user := c.Query("user")
	pwd := c.Query("pwd")
	tunnel := c.DefaultQuery("tunnel", "-1")
	if user == "" || pwd == "" {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	u, err := h.users.GetByUsername(c.Request.Context(), user)
	if err != nil || u.Pwd != md5HexForAPI(pwd) {
		c.JSON(http.StatusOK, errMsg("账号或密码错误"))
		return
	}
	upload, download, total, expire := u.OutFlow, u.InFlow, u.Flow*1024*1024*1024, u.ExpTime/1000
	if tunnel != "-1" {
		tid, err := strconv.ParseInt(tunnel, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, errCode(500, "参数错误"))
			return
		}
		ut, err := h.users.QueryUserTunnel(c.Request.Context(), u.ID, tid)
		if err != nil {
			c.JSON(http.StatusOK, errMsg("隧道权限不存在"))
			return
		}
		upload, download, total, expire = ut.OutFlow, ut.InFlow, ut.Flow*1024*1024*1024, ut.ExpTime/1000
	}
	// 源码 buildSubscriptionHeader(upload=out, download=in) 内部参数位互换 → 实际 upload=<inFlow>; download=<outFlow>
	headerValue := fmt.Sprintf("upload=%d; download=%d; total=%d; expire=%d", download, upload, total, expire)
	c.Header("subscription-userinfo", headerValue)
	c.String(http.StatusOK, headerValue)
}
