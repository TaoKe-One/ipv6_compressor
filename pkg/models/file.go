package models

import (
	"path/filepath"
	"strings"
)

// FileType 文件类型
type FileType int

const (
	FileTypeUnknown FileType = iota
	FileTypeExcel
	FileTypeCSV
)

// String 返回文件类型的字符串表示
func (t FileType) String() string {
	switch t {
	case FileTypeExcel:
		return "Excel"
	case FileTypeCSV:
		return "CSV"
	default:
		return "Unknown"
	}
}

// DetectFileType 检测文件类型
func DetectFileType(filePath string) FileType {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".xlsx", ".xls":
		return FileTypeExcel
	case ".csv":
		return FileTypeCSV
	default:
		return FileTypeUnknown
	}
}

// ProcessingOptions 处理选项
type ProcessingOptions struct {
	SelectedColumns []string  // 选择的列名
	UsePatternMatch bool       // 是否使用正则匹配
	Pattern         string     // 正则表达式模式
	OutputPath      string     // 输出路径
}

// ProcessingResult 处理结果
type ProcessingResult struct {
	Success         bool     // 是否成功
	ProcessedRows   int      // 处理的行数
	ProcessedCount  int      // 实际修改的数量
	OutputFile      string   // 输出文件路径
	Duration        int64    // 处理耗时（毫秒）
	Error           error    // 错误信息
	ModifiedColumns []string // 修改的列名
}

// ColumnData 列数据
type ColumnData struct {
	Name     string
	Index    int
	IPv6Ratio float64
	Selected bool
}

// FileData 文件数据
type FileData struct {
	Path         string      // 文件路径
	Type         FileType    // 文件类型
	RowCount     int         // 行数
	ColumnCount  int         // 列数
	Columns      []ColumnData // 列信息
	PreviewRows  [][]string  // 预览数据（前5行）
}
