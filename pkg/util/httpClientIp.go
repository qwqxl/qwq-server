package util

import (
	"net"
	"net/http"
)

import (
	"strings"
)

func GetClientIP(req *http.Request) string {
	// 优先级 1: 标准 Forwarded 头部 (RFC 7239)
	if forwarded := req.Header.Get("Forwarded"); forwarded != "" {
		if ip := parseForwardedHeader(forwarded); ip != "" {
			return normalizeIP(ip)
		}
	}

	// 优先级 2: 常见代理头部 (按可信度排序)
	headers := []string{
		"X-Forwarded-For",  // 需注意可能包含多个IP
		"CF-Connecting-IP", // Cloudflare
		"Fastly-Client-IP", // Fastly
		"True-Client-IP",   // Akamai 和 Cloudflare
		"X-Real-IP",
	}

	for _, header := range headers {
		value := req.Header.Get(header)
		if value == "" {
			continue
		}

		// 处理逗号分隔的IP列表 (常见于 X-Forwarded-For)
		ips := strings.Split(value, ",")
		for i := range ips {
			ip := strings.TrimSpace(ips[i])

			// 跳过已知的无效IP
			if ip == "" || isPrivateIP(ip) {
				continue
			}

			// 验证IP格式
			if parsed := net.ParseIP(ip); parsed != nil {
				return normalizeIP(ip)
			}
		}
	}

	// 优先级 3: 直接连接地址
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		// 处理无端口的情况
		ip = req.RemoteAddr
	}
	return normalizeIP(ip)
}

// 解析 Forwarded 头部 (e.g. "for=192.0.2.60;proto=http;by=203.0.113.43")
func parseForwardedHeader(header string) string {
	for _, part := range strings.Split(header, ";") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "for=") {
			ip := strings.TrimPrefix(part, "for=")

			// 移除引号和 IPv6 方括号
			ip = strings.Trim(ip, `"[]`)

			// 处理带端口的格式 (e.g. "_gazonk:1234")
			if host, _, err := net.SplitHostPort(ip); err == nil {
				ip = host
			}

			return ip
		}
	}
	return ""
}

// IP标准化处理
func normalizeIP(ip string) string {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return "" // 无效IP
	}

	// IPv4兼容格式处理
	if v4 := parsed.To4(); v4 != nil {
		return v4.String()
	}

	// IPv6环回地址标准化
	if parsed.IsLoopback() {
		return "127.0.0.1"
	}

	return parsed.String()
}

// 检测私有IP地址 (RFC 1918)
func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	return ip.IsLoopback() ||
		ip.IsPrivate() || // Go 1.17+ 新增
		ip.IsUnspecified() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast()
}
