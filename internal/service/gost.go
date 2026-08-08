package service

import (
	"strconv"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/Fyzzp0o0/FLVX2/internal/model"
)

// Gost 命令 data 构造,与 GostUtil.java 逐字段对齐(02-protocol.md §3)

// BuildForwardServices 端口转发服务(tcp+udp 两个)
func BuildForwardServices(svcName string, limiter *int64, node *model.Node, port int64, tunnel *model.Tunnel, remoteAddr string) []map[string]any {
	protocols := []string{"tcp", "udp"}
	services := make([]map[string]any, 0, 2)
	for _, proto := range protocols {
		service := map[string]any{
			"name": svcName + "_" + proto,
		}
		addr := node.TCPListenAddr + ":" + strconv.FormatInt(port, 10)
		if proto == "udp" {
			addr = node.UDPListenAddr + ":" + strconv.FormatInt(port, 10)
		}
		service["addr"] = addr
		// 只在端口转发(type==1)时设置 interface(隧道转发时 interface 在转发链节点上)
		if tunnel.Type == 1 && node.InterfaceName != "" {
			service["metadata"] = map[string]any{"interface": node.InterfaceName}
		}
		if limiter != nil {
			service["limiter"] = strconv.FormatInt(*limiter, 10)
		}
		handler := map[string]any{"type": proto}
		if tunnel.Type == 2 {
			handler["chain"] = "chains_" + strconv.FormatInt(tunnel.ID, 10)
		}
		service["handler"] = handler
		service["listener"] = buildListener(proto)
		service["forwarder"] = buildForwarder(remoteAddr, defaultStrategy)
		services = append(services, service)
	}
	return services
}

// BuildChainService 隧道 TLS 服务(<tunnelId>_tls;chainType 2=中转带 chain,3=出口带 interface)
func BuildChainService(tunnelID int64, chainType int, protocol string, node *model.Node, port int64, ifaceName string) map[string]any {
	service := map[string]any{
		"name": strconv.FormatInt(tunnelID, 10) + "_tls",
		"addr": node.TCPListenAddr + ":" + strconv.FormatInt(port, 10),
	}
	// 只为出口节点(chainType==3)设置 interface
	if chainType == 3 && ifaceName != "" {
		service["metadata"] = map[string]any{"interface": ifaceName}
	}
	handler := map[string]any{"type": "relay"}
	if chainType == 2 {
		handler["chain"] = "chains_" + strconv.FormatInt(tunnelID, 10)
	}
	service["handler"] = handler
	service["listener"] = map[string]any{"type": protocol}
	return service
}

// BuildChainData 链配置(name=chains_<tunnelId>;nodes=下一跳节点列表)
func BuildChainData(tunnelID int64, strategy, hopInterface string, nextHops []ChainHopNode) map[string]any {
	nodes := make([]map[string]any, 0, len(nextHops))
	for _, n := range nextHops {
		nodes = append(nodes, map[string]any{
			"name":      "node_" + strconv.Itoa(n.Inx),
			"addr":      processServerAddress(n.Addr),
			"connector": map[string]any{"type": "relay"},
			"dialer":    map[string]any{"type": n.Protocol},
		})
	}
	hop := map[string]any{
		"name": "hop_" + strconv.FormatInt(tunnelID, 10),
		"selector": map[string]any{
			"strategy":    strategy,
			"maxFails":    1,
			"failTimeout": int64(600000000000), // 600s,纳秒数字(Java 原样)
		},
		"nodes": nodes,
	}
	if hopInterface != "" {
		hop["interface"] = hopInterface
	}
	return map[string]any{
		"name": "chains_" + strconv.FormatInt(tunnelID, 10),
		"hops": []any{hop},
	}
}

// ChainHopNode 链的下一跳节点
type ChainHopNode struct {
	Inx      int
	Addr     string // ip:port
	Protocol string
}

// BuildLimiterData 限速器配置:{"name":"<id>","limits":["$ XMB XMB"]}
func BuildLimiterData(name int64, speedMBPS string) map[string]any {
	return map[string]any{
		"name":   strconv.FormatInt(name, 10),
		"limits": []string{"$ " + speedMBPS + "MB " + speedMBPS + "MB"},
	}
}

// ConvertBitsToMBPS speed/8.0 保留 1 位小数(HALF_UP),与 Java convertBitsToMBps 一致
func ConvertBitsToMBPS(speed int64) string {
	v := decimal.NewFromInt(speed).Div(decimal.NewFromInt(8))
	return v.Round(1).String()
}

// processServerAddress IPv6 加方括号(Java StrUtil/正则逻辑)
func processServerAddress(serverAddr string) string {
	if serverAddr == "" {
		return serverAddr
	}
	// 已是 [ipv6]:port 或 host:port 则原样
	if strings.HasPrefix(serverAddr, "[") {
		return serverAddr
	}
	// 冒号数 >1 视为 IPv6(无括号)
	if strings.Count(serverAddr, ":") > 1 {
		// 形如 2001:db8::1:8080,最后一段为端口
		idx := strings.LastIndex(serverAddr, ":")
		return "[" + serverAddr[:idx] + "]" + serverAddr[idx:]
	}
	return serverAddr
}

func buildListener(protocol string) map[string]any {
	if protocol == "udp" {
		return map[string]any{"type": "udp", "metadata": map[string]any{"keepAlive": true}}
	}
	return map[string]any{"type": protocol}
}

// buildForwarder remoteAddr 逗号分隔多目标 + selector
func buildForwarder(remoteAddr, strategy string) map[string]any {
	nodes := make([]map[string]any, 0)
	for i, addr := range strings.Split(remoteAddr, ",") {
		nodes = append(nodes, map[string]any{
			"name": "node_" + strconv.Itoa(i+1),
			"addr": strings.TrimSpace(addr),
		})
	}
	if strategy == "" {
		strategy = defaultStrategy
	}
	return map[string]any{
		"nodes": nodes,
		"selector": map[string]any{
			"strategy":    strategy,
			"maxFails":    1,
			"failTimeout": "600s",
		},
	}
}

const defaultStrategy = "fifo"
