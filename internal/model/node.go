package model

// Node 对应 node 表
type Node struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Secret         string `json:"-"`
	ServerIP       string `json:"server_ip"`
	Port           string `json:"port"` // 端口范围串 "1000-2000,3000"
	InterfaceName  string `json:"interface_name"`
	Version        string `json:"version"`
	HTTP           int64  `json:"http"`
	TLS            int64  `json:"tls"`
	SOCKS          int64  `json:"socks"`
	CreatedTime    int64  `json:"created_time"`
	UpdatedTime    int64  `json:"updated_time"`
	Status         int64  `json:"status"`
	TCPListenAddr  string `json:"tcp_listen_addr"`
	UDPListenAddr  string `json:"udp_listen_addr"`
}

// Forward 对应 forward 表
type Forward struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	UserName   string `json:"user_name"`
	Name       string `json:"name"`
	TunnelID   int64  `json:"tunnel_id"`
	RemoteAddr string `json:"remote_addr"`
	Strategy   string `json:"strategy"`
	InFlow     int64  `json:"in_flow"`
	OutFlow    int64  `json:"out_flow"`
	CreatedTime int64 `json:"created_time"`
	UpdatedTime int64 `json:"updated_time"`
	Status     int64  `json:"status"`
	Inx        int64  `json:"inx"`
}

// Tunnel 对应 tunnel 表
type Tunnel struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	TrafficRatio float64 `json:"traffic_ratio"`
	Type         int64   `json:"type"`
	Protocol     string  `json:"protocol"`
	Flow         int64   `json:"flow"`
	CreatedTime  int64   `json:"created_time"`
	UpdatedTime  int64   `json:"updated_time"`
	Status       int64   `json:"status"`
	InIP         string  `json:"in_ip"`
}

// UserTunnel 对应 user_tunnel 表
type UserTunnel struct {
	ID            int64  `json:"id"`
	UserID        int64  `json:"user_id"`
	TunnelID      int64  `json:"tunnel_id"`
	SpeedID       *int64 `json:"speed_id"`
	Num           int64  `json:"num"`
	Flow          int64  `json:"flow"`
	InFlow        int64  `json:"in_flow"`
	OutFlow       int64  `json:"out_flow"`
	FlowResetTime int64  `json:"flow_reset_time"`
	ExpTime       int64  `json:"exp_time"`
	Status        int64  `json:"status"`
}

// ChainTunnel 对应 chain_tunnel 表(chain_type 为字符串 "1"入口/"2"转发链/"3"出口)
type ChainTunnel struct {
	ID        int64  `json:"id"`
	TunnelID  int64  `json:"tunnel_id"`
	ChainType string `json:"chain_type"`
	NodeID    int64  `json:"node_id"`
	Port      *int64 `json:"port"`
	Strategy  string `json:"strategy"`
	Inx       int64  `json:"inx"`
	Protocol  string `json:"protocol"`
}
