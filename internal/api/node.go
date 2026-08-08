package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---- NodeController /api/v1/node(全部管理员) ----

type createNodeReq struct {
	Name          string `json:"name" binding:"required"`
	ServerIP      string `json:"serverIp" binding:"required"`
	Port          string `json:"port" binding:"required"`
	InterfaceName string `json:"interfaceName"`
	TCPListenAddr string `json:"tcpListenAddr"`
	UDPListenAddr string `json:"udpListenAddr"`
}

// NodeCreate 创建节点
func (h *Handler) NodeCreate(c *gin.Context) {
	var req createNodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	node, err := h.nodes.Create(c.Request.Context(), req.Name, req.ServerIP, req.Port, req.InterfaceName, req.TCPListenAddr, req.UDPListenAddr)
	if err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	node.Secret = "" // 不回传密钥
	c.JSON(http.StatusOK, ok(node))
}

// NodeList 节点列表
func (h *Handler) NodeList(c *gin.Context) {
	list, err := h.nodes.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	for _, n := range list {
		n.Secret = "" // 密钥置 null
	}
	c.JSON(http.StatusOK, ok(list))
}

type updateNodeReq struct {
	ID            int64  `json:"id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	ServerIP      string `json:"serverIp" binding:"required"`
	Port          string `json:"port" binding:"required"`
	InterfaceName string `json:"interfaceName"`
	HTTP          *int64 `json:"http"`
	TLS           *int64 `json:"tls"`
	SOCKS         *int64 `json:"socks"`
	TCPListenAddr string `json:"tcpListenAddr"`
	UDPListenAddr string `json:"udpListenAddr"`
}

// NodeUpdate 更新节点(在线且协议变化 → SetProtocol 全量下发)
func (h *Handler) NodeUpdate(c *gin.Context) {
	var req updateNodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	httpV, tlsV, socksV := int64(0), int64(0), int64(0)
	if req.HTTP != nil {
		httpV = *req.HTTP
	}
	if req.TLS != nil {
		tlsV = *req.TLS
	}
	if req.SOCKS != nil {
		socksV = *req.SOCKS
	}
	if err := h.nodes.Update(c.Request.Context(), req.ID, req.Name, req.ServerIP, req.Port, req.InterfaceName,
		httpV, tlsV, socksV, req.TCPListenAddr, req.UDPListenAddr); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// NodeDelete 删除节点
func (h *Handler) NodeDelete(c *gin.Context) {
	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	if err := h.nodes.Delete(c.Request.Context(), req.ID, h.tunnels, h.forwards); err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(nil))
}

// NodeInstall 返回节点安装命令
func (h *Handler) NodeInstall(c *gin.Context) {
	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusOK, errCode(500, "参数错误"))
		return
	}
	cmd, err := h.nodes.InstallCmd(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusOK, errMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ok(cmd))
}
