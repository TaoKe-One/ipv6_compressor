package processor

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// CSVProcessor CSV 文件处理器
type CSVProcessor struct {
	rows     [][]string
	filePath string
}

// NewCSVProcessor 创建 CSV 处理器
func NewCSVProcessor(filePath string) (*CSVProcessor, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开 CSV 文件失败: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("读取 CSV 失败: %w", err)
	}

	return &CSVProcessor{
		rows:     rows,
		filePath: filePath,
	}, nil
}

// GetRows 获取所有行
func (p *CSVProcessor) GetRows() [][]string {
	return p.rows
}

// GetRowCount 获取行数
func (p *CSVProcessor) GetRowCount() int {
	return len(p.rows)
}

// GetColumnCount 获取列数
func (p *CSVProcessor) GetColumnCount() int {
	if len(p.rows) == 0 {
		return 0
	}
	return len(p.rows[0])
}

// GetColumnNames 获取列名（第一行）
func (p *CSVProcessor) GetColumnNames() []string {
	if len(p.rows) == 0 {
		return nil
	}
	return p.rows[0]
}

// ProcessColumn 处理指定列
func (p *CSVProcessor) ProcessColumn(colIndex int, processFunc func(string) string) (int, error) {
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
			p.rows[rowIdx][colIndex] = processedValue
			processed++
		}
	}

	return processed, nil
}

// ProcessColumnsByName 根据列名处理
func (p *CSVProcessor) ProcessColumnsByName(colNames []string, processFunc func(string) string) (int, error) {
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
func (p *CSVProcessor) Save(outputPath string) error {
	// 生成输出路径
	if outputPath == "" {
		outputPath = p.generateOutputPath()
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	for _, row := range p.rows {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("写入 CSV 失败: %w", err)
		}
	}

	return nil
}

// generateOutputPath 生成输出文件路径
func (p *CSVProcessor) generateOutputPath() string {
	path := p.filePath
	// 添加 _compressed 后缀
	lastDot := strings.LastIndex(path, ".")
	if lastDot == -1 {
		return path + "_compressed"
	}
	return path[:lastDot] + "_compressed" + path[lastDot:]
}

// GetFilePath 获取文件路径
func (p *CSVProcessor) GetFilePath() string {
	return p.filePath
}

// Close 关闭处理器（CSV 处理器无需关闭资源）
func (p *CSVProcessor) Close() error {
	return nil
}
