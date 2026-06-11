package processor

import (
	"github.com/TaoKe-One/ipv6-compressor/internal/ipv6"
)

// ColumnInfo 列信息
type ColumnInfo struct {
	Name     string
	Index    int
	IPv6Ratio float64
	IsIPv6   bool
}

// DetectIPv6Columns 自动识别包含 IPv6 地址的列
// sampleSize: 采样检查的行数，提高性能
// threshold: 判定为 IPv6 列的比例阈值（0-1）
func DetectIPv6Columns(rows [][]string, sampleSize, threshold int) []ColumnInfo {
	if len(rows) == 0 {
		return nil
	}

	var columns []ColumnInfo
	numCols := len(rows[0])

	// 确定实际采样数量
	actualSample := sampleSize
	if len(rows) < sampleSize {
		actualSample = len(rows)
	}

	for colIdx := 0; colIdx < numCols; colIdx++ {
		ipv6Count := 0
		validCount := 0

		for rowIdx := 0; rowIdx < actualSample; rowIdx++ {
			if colIdx >= len(rows[rowIdx]) {
				continue
			}

			value := rows[rowIdx][colIdx]
			if value == "" {
				continue
			}

			validCount++
			if ipv6.IsIPv6(value) {
				ipv6Count++
			}
		}

		if validCount == 0 {
			continue
		}

		ratio := float64(ipv6Count) / float64(validCount)
		isIPv6Col := ratio >= float64(threshold)/100.0

		columns = append(columns, ColumnInfo{
			Name:     rows[0][colIdx], // 假设第一行是表头
			Index:    colIdx,
			IPv6Ratio: ratio,
			IsIPv6:   isIPv6Col,
		})
	}

	return columns
}

// FilterIPv6Columns 过滤出 IPv6 列
func FilterIPv6Columns(columns []ColumnInfo) []string {
	var result []string
	for _, col := range columns {
		if col.IsIPv6 {
			result = append(result, col.Name)
		}
	}
	return result
}
