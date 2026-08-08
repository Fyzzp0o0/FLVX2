package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Fyzzp0o0/FLVX2/internal/service"
)

// ---- TunnelController /api/v1/tunnel ----

// TunnelCreate 创建隧道
func (h *Handler) TunnelCreate(c *gin.Context) {
	var req service.TunnelCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	if err := h.tunnels.Create(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// TunnelList 隧道列表
func (h *Handler) TunnelList(c *gin.Context) {
	list, err := h.tunnels.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(list))
}

type updateTunnelReq struct {
	ID           int64   `json:"id" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Flow         int64   `json:"flow" binding:"required"`
	InIP         string  `json:"inIp"`
	TrafficRatio float64 `json:"trafficRatio"`
}

// TunnelUpdate 更新隧道
func (h *Handler) TunnelUpdate(c *gin.Context) {
	var req updateTunnelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	if err := h.tunnels.Update(c.Request.Context(), req.ID, req.Name, req.Flow, req.InIP, req.TrafficRatio); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// TunnelDelete 删除隧道
func (h *Handler) TunnelDelete(c *gin.Context) {
	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	if err := h.tunnels.Delete(c.Request.Context(), req.ID, h.forwards); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// ---- 用户隧道授权 ----

type assignUserTunnelReq struct {
	UserID        int64  `json:"userId" binding:"required"`
	TunnelID      int64  `json:"tunnelId" binding:"required"`
	Flow          int64  `json:"flow" binding:"required"`
	Num           int64  `json:"num" binding:"required"`
	FlowResetTime int64  `json:"flowResetTime" binding:"required"`
	ExpTime       int64  `json:"expTime" binding:"required"`
	SpeedID       *int64 `json:"speedId"`
}

// TunnelUserAssign 分配隧道权限
func (h *Handler) TunnelUserAssign(c *gin.Context) {
	var req assignUserTunnelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	if err := h.tunnels.AssignUserTunnel(c.Request.Context(), req.UserID, req.TunnelID, req.Flow, req.Num, req.FlowResetTime, req.ExpTime, req.SpeedID); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// TunnelUserList 某用户隧道权限
func (h *Handler) TunnelUserList(c *gin.Context) {
	var req struct {
		UserID int64 `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	list, err := h.tunnels.ListUserTunnels(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(list))
}

// TunnelUserRemove 移除隧道权限
func (h *Handler) TunnelUserRemove(c *gin.Context) {
	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	if err := h.tunnels.RemoveUserTunnel(c.Request.Context(), req.ID, h.forwards); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

type updateUserTunnelReq struct {
	ID            int64  `json:"id" binding:"required"`
	Flow          int64  `json:"flow" binding:"required"`
	Num           int64  `json:"num" binding:"required"`
	FlowResetTime int64  `json:"flowResetTime" binding:"required"`
	ExpTime       int64  `json:"expTime" binding:"required"`
	Status        int64  `json:"status" binding:"required"`
	SpeedID       *int64 `json:"speedId"`
}

// TunnelUserUpdate 更新隧道权限(修正原 Java bug:成功返回 ok)
func (h *Handler) TunnelUserUpdate(c *gin.Context) {
	var req updateUserTunnelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	if err := h.tunnels.UpdateUserTunnel(c.Request.Context(), req.ID, req.Flow, req.Num, req.FlowResetTime, req.ExpTime, req.Status, req.SpeedID); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// TunnelUserTunnel 当前用户可见隧道(仅 JWT,无 @RequireRole)
func (h *Handler) TunnelUserTunnel(c *gin.Context) {
	uid, _, roleID := currentUser(c)
	list, err := h.tunnels.UserTunnels(c.Request.Context(), uid, roleID)
	if err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(list))
}
