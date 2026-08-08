package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---- ForwardController /api/v1/forward(仅 JWT,普通用户仅自己的) ----

type createForwardReq struct {
	Name       string `json:"name" binding:"required"`
	TunnelID   int64  `json:"tunnelId" binding:"required"`
	RemoteAddr string `json:"remoteAddr" binding:"required"`
	Strategy   string `json:"strategy"`
	InPort     int64  `json:"inPort"`
}

// ForwardCreate 创建转发
func (h *Handler) ForwardCreate(c *gin.Context) {
	var req createForwardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	uid, name, roleID := currentUser(c)
	if err := h.forwards.Create(c.Request.Context(), req.Name, req.TunnelID, req.RemoteAddr, req.Strategy, req.InPort, uid, name, roleID); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// ForwardList 转发列表(角色分支)
func (h *Handler) ForwardList(c *gin.Context) {
	uid, _, roleID := currentUser(c)
	list, err := h.forwards.List(c.Request.Context(), uid, roleID)
	if err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(list))
}

type updateForwardReq struct {
	ID         int64  `json:"id" binding:"required"`
	UserID     int64  `json:"userId" binding:"required"`
	Name       string `json:"name" binding:"required"`
	RemoteAddr string `json:"remoteAddr" binding:"required"`
	Strategy   string `json:"strategy"`
	InPort     int64  `json:"inPort"`
}

// ForwardUpdate 更新转发(重分配端口 + 重推服务,M4 完整)
func (h *Handler) ForwardUpdate(c *gin.Context) {
	var req updateForwardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	uid, name, roleID := currentUser(c)
	// M3 简化:删除重建(端口/地址变化;属主校验在 DeleteByID);M4 将改为 UpdateService 语义
	tunnelID := h.forwards.TunnelIDOf(c.Request.Context(), req.ID)
	if tunnelID == 0 {
		c.JSON(http.StatusOK, errMsg("转发不存在"))
		return
	}
	if err := h.forwards.DeleteByID(c.Request.Context(), req.ID, uid, roleID); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	if err := h.forwards.Create(c.Request.Context(), req.Name, tunnelID, req.RemoteAddr, req.Strategy, req.InPort, uid, name, roleID); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// ForwardDelete 删除转发
func (h *Handler) ForwardDelete(c *gin.Context) {
	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	uid, _, roleID := currentUser(c)
	if err := h.forwards.Delete(c.Request.Context(), req.ID, uid, roleID); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// ForwardForceDelete 强制删除(仅删 DB,残留服务由孤儿清理兜底)
func (h *Handler) ForwardForceDelete(c *gin.Context) {
	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	uid, _, roleID := currentUser(c)
	if err := h.forwards.ForceDelete(c.Request.Context(), req.ID, uid, roleID); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// ForwardPause 暂停转发
func (h *Handler) ForwardPause(c *gin.Context) {
	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	uid, _, roleID := currentUser(c)
	if err := h.forwards.Pause(c.Request.Context(), req.ID, uid, roleID); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// ForwardResume 恢复转发
func (h *Handler) ForwardResume(c *gin.Context) {
	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	uid, _, roleID := currentUser(c)
	if err := h.forwards.Resume(c.Request.Context(), req.ID, uid, roleID); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// ForwardUpdateOrder 批量更新排序
func (h *Handler) ForwardUpdateOrder(c *gin.Context) {
	var req struct {
		Forwards []struct {
			ID  int64 `json:"id"`
			Inx int64 `json:"inx"`
		} `json:"forwards" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	uid, _, roleID := currentUser(c)
	if err := h.forwards.UpdateOrder(c.Request.Context(), req.Forwards, uid, roleID); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}
