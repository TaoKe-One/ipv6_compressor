package processor

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ExcelProcessor Excel 文件处理器
type ExcelProcessor struct {
	file     *excelize.File
	sheet    string
	rows     [][]string
	filePath string
}

// NewExcelProcessor 创建 Excel 处理器
func NewExcelProcessor(filePath string) (*ExcelProcessor, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开 Excel 文件失败: %w", err)
	}

	// 获取第一个工作表
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		f.Close()
		return nil, fmt.Errorf("Excel 文件中没有工作表")
	}

	// 读取所有行
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("读取 Excel 行失败: %w", err)
	}

	return &ExcelProcessor{
		file:     f,
		sheet:    sheets[0],
		rows:     rows,
		filePath: filePath,
	}, nil
}

// Close 关闭文件
func (p *ExcelProcessor) Close() error {
	return p.file.Close()
}

// GetRows 获取所有行
func (p *ExcelProcessor) GetRows() [][]string {
	return p.rows
}

// GetRowCount 获取行数
func (p *ExcelProcessor) GetRowCount() int {
	return len(p.rows)
}

// GetColumnCount 获取列数
func (p *ExcelProcessor) GetColumnCount() int {
	if len(p.rows) == 0 {
		return 0
	}
	return len(p.rows[0])
}

// GetColumnNames 获取列名（第一行）
func (p *ExcelProcessor) GetColumnNames() []string {
	if len(p.rows) == 0 {
		return nil
	}
	return p.rows[0]
}

// ProcessColumn 处理指定列
// 返回处理后的行数和错误
func (p *ExcelProcessor) ProcessColumn(colIndex int, processFunc func(string) string) (int, error) {
	if len(p.rows) <= 1 {
		return 0, nil
	}

	processed := 0
	// 从第二行开始（跳过表头）
	for rowIdx := 1; rowIdx < len(p.rows); rowIdx++ {
		if colIndex >= len(p.rows[rowIdx]) {
			continue
		}

		original := p.rows[rowIdx][colIndex]
		processedValue := processFunc(original)

		if processedValue != original {
			cell, err := excelize.CoordinatesToCellName(colIndex+1, rowIdx+1)
			if err != nil {
				continue
			}
			p.rows[rowIdx][colIndex] = processedValue
			// 更新 Excel 文件
			p.file.SetCellValue(p.sheet, cell, processedValue)
			processed++
		}
	}

	return processed, nil
}

// ProcessColumnsByName 根据列名处理
func (p *ExcelProcessor) ProcessColumnsByName(colNames []string, processFunc func(string) string) (int, error) {
	if len(p.rows) == 0 {
		return 0, nil
	}

	// 找到列索引
	nameToIndex := make(map[string]int)
	header := p.rows[0]
	for i, name := range header {
		nameToIndex[name] = i
	}

	totalProcessed := 0
	for _, name := range colNames {
		colIndex, exists := nameToIndex[name]
		if !exists {
			continue
		}
		processed, err := p.ProcessColumn(colIndex, processFunc)
		if err != nil {
			return totalProcessed, err
		}
		totalProcessed += processed
	}

	return totalProcessed, nil
}

// Save 保存文件
func (p *ExcelProcessor) Save(outputPath string) error {
	// 生成输出路径
	if outputPath == "" {
		outputPath = p.generateOutputPath()
	}

	return p.file.SaveAs(outputPath)
}

// generateOutputPath 生成输出文件路径
func (p *ExcelProcessor) generateOutputPath() string {
	path := p.filePath
	// 添加 _compressed 后缀
	lastDot := strings.LastIndex(path, ".")
	if lastDot == -1 {
		return path + "_compressed"
	}
	return path[:lastDot] + "_compressed" + path[lastDot:]
}

// GetFilePath 获取文件路径
func (p *ExcelProcessor) GetFilePath() string {
	return p.filePath
}
