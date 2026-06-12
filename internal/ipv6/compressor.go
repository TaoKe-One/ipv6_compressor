package ipv6

import (
	"fmt"
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

// ExpandIPv6 将 IPv6 地址展开为完整格式（8组，每组4位十六进制）
// 例如: 2001:db8::1 → 2001:0db8:0000:0000:0000:0000:0000:0001
func ExpandIPv6(ipStr string) string {
	if ipStr == "" {
		return ipStr
	}

	ipStr = strings.TrimSpace(ipStr)
	if !IsIPv6(ipStr) {
		return ipStr
	}

	// 使用 net.ParseIP 解析地址
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ipStr
	}

	// 获取 IPv6 的 16 字节表示
	ip = ip.To16()
	if ip == nil {
		return ipStr
	}

	// 将每 2 字节转换为 4 位十六进制
	var groups []string
	for i := 0; i < 16; i += 2 {
		group := fmt.Sprintf("%02x%02x", ip[i], ip[i+1])
		groups = append(groups, group)
	}

	return strings.Join(groups, ":")
}

// ProcessMode 处理模式
type ProcessMode int

const (
	// ModeCompress 压缩模式（RFC 5952）
	ModeCompress ProcessMode = iota
	// ModeExpand 扩展模式（完整格式）
	ModeExpand
)

// ProcessIPv6 根据模式处理 IPv6 地址
func ProcessIPv6(ipStr string, mode ProcessMode) string {
	switch mode {
	case ModeCompress:
		return CompressIPv6(ipStr)
	case ModeExpand:
		return ExpandIPv6(ipStr)
	default:
		return ipStr
	}
}
