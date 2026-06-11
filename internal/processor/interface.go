package processor

// FileProcessor 文件处理器接口
type FileProcessor interface {
	GetRows() [][]string
	GetRowCount() int
	GetColumnCount() int
	GetColumnNames() []string
	ProcessColumn(colIndex int, processFunc func(string) string) (int, error)
	ProcessColumnsByName(colNames []string, processFunc func(string) string) (int, error)
	Save(outputPath string) error
	GetFilePath() string
	Close() error
}
