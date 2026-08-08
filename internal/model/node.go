package model

// Node 对应 node 表
type Node struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Secret         string `json:"-"`
	ServerIP       string `json:"serverIp"`
	Port           string `json:"port"` // 端口范围串 "1000-2000,3000"
	InterfaceName  string `json:"interfaceName"`
	Version        string `json:"version"`
	HTTP           int64  `json:"http"`
	TLS            int64  `json:"tls"`
	SOCKS          int64  `json:"socks"`
	CreatedTime    int64  `json:"createdTime"`
	UpdatedTime    int64  `json:"updatedTime"`
	Status         int64  `json:"status"`
	TCPListenAddr  string `json:"tcpListenAddr"`
	UDPListenAddr  string `json:"udpListenAddr"`
}

// Forward 对应 forward 表
type Forward struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"userId"`
	UserName   string `json:"userName"`
	Name       string `json:"name"`
	TunnelID   int64  `json:"tunnelId"`
	RemoteAddr string `json:"remoteAddr"`
	Strategy   string `json:"strategy"`
	InFlow     int64  `json:"inFlow"`
	OutFlow    int64  `json:"outFlow"`
	CreatedTime int64 `json:"createdTime"`
	UpdatedTime int64 `json:"updatedTime"`
	Status     int64  `json:"status"`
	Inx        int64  `json:"inx"`
}

// Tunnel 对应 tunnel 表
type Tunnel struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	TrafficRatio float64 `json:"trafficRatio"`
	Type         int64   `json:"type"`
	Protocol     string  `json:"protocol"`
	Flow         int64   `json:"flow"`
	CreatedTime  int64   `json:"createdTime"`
	UpdatedTime  int64   `json:"updatedTime"`
	Status       int64   `json:"status"`
	InIP         string  `json:"inIp"`
}

// UserTunnel 对应 user_tunnel 表
type UserTunnel struct {
	ID            int64  `json:"id"`
	UserID        int64  `json:"userId"`
	TunnelID      int64  `json:"tunnelId"`
	SpeedID       *int64 `json:"speedId"`
	Num           int64  `json:"num"`
	Flow          int64  `json:"flow"`
	InFlow        int64  `json:"inFlow"`
	OutFlow       int64  `json:"outFlow"`
	FlowResetTime int64  `json:"flowResetTime"`
	ExpTime       int64  `json:"expTime"`
	Status        int64  `json:"status"`
}

// ChainTunnel 对应 chain_tunnel 表(chain_type 为字符串 "1"入口/"2"转发链/"3"出口)
type ChainTunnel struct {
	ID        int64  `json:"id"`
	TunnelID  int64  `json:"tunnelId"`
	ChainType string `json:"chainType"`
	NodeID    int64  `json:"nodeId"`
	Port      *int64 `json:"port"`
	Strategy  string `json:"strategy"`
	Inx       int64  `json:"inx"`
	Protocol  string `json:"protocol"`
}

// SpeedLimit 对应 speed_limit 表
type SpeedLimit struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Speed       int64  `json:"speed"`
	TunnelID    int64  `json:"tunnelId"`
	TunnelName  string `json:"tunnelName"`
	CreatedTime int64  `json:"createdTime"`
	UpdatedTime int64  `json:"updatedTime"`
	Status      int64  `json:"status"`
}

// StatisticsFlow 对应 statistics_flow 表
type StatisticsFlow struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"userId"`
	Flow        int64  `json:"flow"`
	TotalFlow   int64  `json:"totalFlow"`
	Time        string `json:"time"`
	CreatedTime int64  `json:"createdTime"`
}
