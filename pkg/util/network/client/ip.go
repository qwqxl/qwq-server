package client

import (
	"math/big"
	"net"
	"net/http"
	"strings"
	"sync"
)

var (
	privateNetworks []*net.IPNet
	ipInitOnce      sync.Once
)

func init() {
	initPrivateNetworks()
}

func initPrivateNetworks() {
	_, ipv4Unicast, _ := net.ParseCIDR("0.0.0.0/8")
	_, rfc1918_1, _ := net.ParseCIDR("10.0.0.0/8")
	_, rfc1918_2, _ := net.ParseCIDR("172.16.0.0/12")
	_, rfc1918_3, _ := net.ParseCIDR("192.168.0.0/16")
	_, carrierNAT, _ := net.ParseCIDR("100.64.0.0/10")
	_, linkLocal, _ := net.ParseCIDR("169.254.0.0/16")
	_, loopback, _ := net.ParseCIDR("127.0.0.0/8")
	_, ipv6Local, _ := net.ParseCIDR("fe80::/10")
	_, ipv6UniqueLocal, _ := net.ParseCIDR("fc00::/7")

	privateNetworks = []*net.IPNet{
		ipv4Unicast, rfc1918_1, rfc1918_2, rfc1918_3,
		carrierNAT, linkLocal, loopback,
		ipv6Local, ipv6UniqueLocal,
	}
}

// GetClientIP 获取客户端IP
func GetClientIP(r *http.Request) string {
	// 优先级 1: 标准 Forwarded 头部 (RFC 7239)
	if forwarded := r.Header.Get("Forwarded"); forwarded != "" {
		if ip := parseForwardedHeader(forwarded); ip != "" {
			return ip
		}
	}

	// 优先级 2: 代理专用头部
	proxyHeaders := []struct {
		name    string
		multi   bool
		trusted bool
	}{
		{"X-Client-IP", false, true},
		{"CF-Connecting-IP", false, true},
		{"Fastly-Client-IP", false, true},
		{"True-Client-IP", false, true},
		{"Fly-Client-IP", false, true},
		{"X-Forwarded-For", true, false},
		{"X-Real-IP", false, false},
		{"X-Cluster-Client-IP", false, false},
	}

	for _, header := range proxyHeaders {
		value := r.Header.Get(header.name)
		if value == "" {
			continue
		}

		ips := strings.Split(value, ",")
		for i := range ips {
			idx := i
			if header.multi {
				idx = len(ips) - 1 - i
			}

			ipStr := strings.TrimSpace(ips[idx])
			ip := validateAndNormalizeIP(ipStr)

			if ip == "" {
				continue
			}

			if header.trusted || !isPrivateIP(ip) {
				return ip
			}
		}
	}

	// 优先级 3: 直接连接地址
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	return validateAndNormalizeIP(ip)
}

// 其他辅助函数...

// 补充

// parseForwardedHeader 解析Forwarded头部
func parseForwardedHeader(header string) string {
	// 实现完整的RFC 7239解析
	parts := strings.Split(header, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "for=") {
			ip := strings.TrimPrefix(part, "for=")
			ip = strings.Trim(ip, `"[] `)

			// 移除IPv6的方括号
			if strings.HasPrefix(ip, "[") && strings.HasSuffix(ip, "]") {
				ip = ip[1 : len(ip)-1]
			}

			// 处理带端口的情况
			if host, _, err := net.SplitHostPort(ip); err == nil {
				ip = host
			}

			return ip
		}
	}
	return ""
}

// validateAndNormalizeIP 验证和标准化IP
func validateAndNormalizeIP(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ""
	}

	// 处理IPv4映射的IPv6地址
	if ip.To4() != nil {
		return ip.String()
	}

	// IPv6标准化
	return "[" + ip.String() + "]"
}

// isPrivateIP 检查是否为私有IP
func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	for _, network := range privateNetworks {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

// IsIPv6 检查是否为IPv6地址
func IsIPv6(ip string) bool {
	return strings.Contains(ip, ":")
}

// IPToInt IP转整数
func IPToInt(ip net.IP) *big.Int {
	i := big.NewInt(0)
	i.SetBytes(ip)
	return i
}
