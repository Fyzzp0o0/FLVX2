package ws

// 12 种命令的 data 构造,字段与 02-protocol.md §3 逐字对齐(agent 按名匹配)

// PauseServicesData PauseService/ResumeService 的 data
type PauseServicesData struct {
	Services []string `json:"services"`
}

// DeleteServicesData DeleteService 的 data
type DeleteServicesData struct {
	Services []string `json:"services"`
}

// DeleteChainsData DeleteChains 的 data(对象形态;agent 兼容直接字符串)
type DeleteChainsData struct {
	Chain string `json:"chain"`
}

// DeleteLimitersData DeleteLimiters 的 data
type DeleteLimitersData struct {
	Limiter string `json:"limiter"`
}

// UpdateLimitersData UpdateLimiters 的 data
type UpdateLimitersData struct {
	Limiter string      `json:"limiter"`
	Data    LimiterConfig `json:"data"`
}

// LimiterConfig AddLimiters 的 data:{"name":"<id>","limits":["$ XMB XMB"]}
type LimiterConfig struct {
	Name   string   `json:"name"`
	Limits []string `json:"limits"`
}

// SetProtocolData SetProtocol 的 data(必须全量;未提供字段 agent 视为 0 覆盖)
type SetProtocolData struct {
	HTTP  int `json:"http"`
	TLS   int `json:"tls"`
	SOCKS int `json:"socks"`
}

// TcpPingData TcpPing 的 data
type TcpPingData struct {
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Count   int    `json:"count"`
	Timeout int    `json:"timeout"`
}
