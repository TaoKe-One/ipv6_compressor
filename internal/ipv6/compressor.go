package ipv6

import (
	"net"
	"regexp"
	"strings"
)

// IPv6 正则表达式 - 用于快速识别
var ipv6Pattern = regexp.MustCompile(`(?:(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|` +
	`(?:[0-9a-fA-F]{1,4}:){1,7}:|` +
	`:(?:(?:[0-9a-fA-F]{1,4}:){1,7}|)|` +
	`[0-9a-fA-F]{1,4}::(?:[0-9a-fA-F]{1,4}:){0,5}[0-9a-fA-F]{1,4}|` +
	`::(?:[0-9a-fA-F]{1,4}:){1,7}[0-9a-fA-F]{1,4})`)

// IsIPv6 判断字符串是否为 IPv6 地址
func IsIPv6(value string) bool {
	if value == "" || strings.TrimSpace(value) == "" {
		return false
	}
	return ipv6Pattern.MatchString(strings.TrimSpace(value))
}

// CompressIPv6 将 IPv6 地址转换为 RFC 5952 简写格式
// - 移除前导零
// - 使用 :: 压缩最长的连续零块
func CompressIPv6(ipStr string) string {
	if ipStr == "" {
		return ipStr
	}

	ipStr = strings.TrimSpace(ipStr)
	if !IsIPv6(ipStr) {
		return ipStr
	}

	// 使用 Go 的 net 包自动处理简写（符合 RFC 5952）
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ipStr
	}

	// IPv6 地址压缩格式
	return ip.String()
}
