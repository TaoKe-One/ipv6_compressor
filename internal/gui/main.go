package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	ipv6pkg "github.com/TaoKe-One/ipv6_compressor/internal/ipv6"
	"github.com/TaoKe-One/ipv6_compressor/internal/processor"
	"github.com/TaoKe-One/ipv6_compressor/pkg/models"
)

// AppState 应用状态
type AppState struct {
	currentFile   processor.FileProcessor
	fileType      models.FileType
	columns       []processor.ColumnInfo
	selectedCols  map[string]bool
	outputPath    string
	isProcessing  bool
	result        *models.ProcessingResult
	statusText    string
	processMode   ipv6pkg.ProcessMode

	// UI 引用
	app           fyne.App
	window        fyne.Window
	fileInfoPanel *FileInfoPanel
	columnSelector *ColumnSelector
	progressBar   *ProgressBar
	statusBar     *StatusBar
	startBtn      *widget.Button
	outputEntry   *widget.Entry
	lastDirPath   string  // 记住上次使用的目录
}

// appState 全局应用状态
var appState = &AppState{
	selectedCols: make(map[string]bool),
	statusText:   "准备就绪",
	processMode:  ipv6pkg.ModeCompress, // 默认压缩模式
}

// LoadUI 加载主界面
func LoadUI(w fyne.Window, a fyne.App) {
	appState.window = w
	appState.app = a

	// 创建 UI 组件
	appState.fileInfoPanel = NewFileInfoPanel()
	appState.progressBar = NewProgressBar()
	appState.statusBar = NewStatusBar()
	appState.startBtn = widget.NewButton("开始处理", startProcessing)
	appState.startBtn.Disable()

	// 主内容
	content := createMainContent()

	w.SetContent(content)
	w.CenterOnScreen()

	// 设置文件拖拽处理
	w.SetOnDropped(func(pos fyne.Position, dropped []fyne.URI) {
		if len(dropped) > 0 {
			filePath := dropped[0].Path()
			loadFile(filePath)
		}
	})
}

// createMainContent 创建主内容
func createMainContent() *fyne.Container {
	// 拖拽区域
	dropArea := createDropArea()

	// 选择文件按钮
	selectFileBtn := widget.NewButton("选择文件", showFilePicker)

	// 处理模式选择
	modeSelect := widget.NewRadioGroup([]string{"压缩 (RFC 5952)", "扩展 (完整格式)"}, func(selected string) {
		if selected == "压缩 (RFC 5952)" {
			appState.processMode = ipv6pkg.ModeCompress
		} else if selected == "扩展 (完整格式)" {
			appState.processMode = ipv6pkg.ModeExpand
		}
	})
	modeSelect.SetSelected("压缩 (RFC 5952)")

	// 左侧面板
	leftPanel := container.NewVBox(
		dropArea,
		selectFileBtn,
		widget.NewSeparator(),
		appState.fileInfoPanel.Container(),
		widget.NewSeparator(),
		container.NewVBox(
			widget.NewLabel("处理模式:"),
			modeSelect,
		),
		widget.NewSeparator(),
		container.NewVBox(
			widget.NewLabel("检测到的列 (IPv6):"),
		),
	)

	// 输出路径区域
	outputEntry := widget.NewEntry()
	outputEntry.SetPlaceHolder("默认: 原文件同目录/原文件名_压缩版.扩展名")
	outputEntry.OnChanged = func(s string) {
		appState.outputPath = s
	}
	// 保存引用以便后续更新
	appState.outputEntry = outputEntry

	outputGroup := container.NewBorder(
		nil, nil,
		widget.NewLabel("输出路径:"), widget.NewButton("浏览", showOutputPicker),
		outputEntry,
	)

	// 右侧面板
	rightPanel := container.NewVBox(
		appState.progressBar.Container(),
		appState.startBtn,
		outputGroup,
		appState.statusBar.Container(),
	)

	// 主布局
	mainContainer := container.NewBorder(
		nil, nil, nil, nil,
		container.NewHSplit(leftPanel, rightPanel),
	)

	return mainContainer
}

// createDropArea 创建拖拽区域
func createDropArea() *fyne.Container {
	card := widget.NewCard(
		"",
		"拖拽 Excel/CSV 文件到这里，或点击下方按钮选择文件",
		widget.NewLabel(""),
	)
	return container.NewPadded(card)
}

// showFilePicker 显示文件选择器
func showFilePicker() {
	// 创建自定义文件选择器（可调整大小）
	picker := NewFilePicker(appState.window, appState.app, func(filePath string) {
		loadFile(filePath)
		// 记住本次使用的目录
		appState.lastDirPath = filepath.Dir(filePath)
	})

	// 设置文件过滤器 - 只显示支持的文件类型
	picker.SetFilter(func(name string) bool {
		ext := filepath.Ext(name)
		ext = strings.ToLower(ext)
		for _, validExt := range []string{".xlsx", ".xls", ".csv"} {
			if ext == validExt {
				return true
			}
		}
		return false
	})

	picker.Show()
}

// FileExtensionFilter 文件扩展名过滤器
type FileExtensionFilter struct {
	extensions []string
}

// NewFileExtensionFilter 创建文件扩展名过滤器
func NewFileExtensionFilter(extensions []string) *FileExtensionFilter {
	return &FileExtensionFilter{extensions: extensions}
}

// Matches 实现 storage.FileFilter 接口
func (f *FileExtensionFilter) Matches(uri fyne.URI) bool {
	path := uri.Path()
	ext := filepath.Ext(path)
	ext = strings.ToLower(ext)

	for _, validExt := range f.extensions {
		if ext == validExt {
			return true
		}
	}
	// 如果没有扩展名，也允许选择（可能是目录或其他情况）
	if ext == "" {
		return true
	}
	return false
}

// showOutputPicker 显示输出文件选择器
func showOutputPicker() {
	// 创建自定义文件夹选择器（可调整大小）
	picker := NewFilePicker(appState.window, appState.app, func(filePath string) {
		// 用户选择了文件，直接使用
		appState.outputPath = filePath
		if appState.outputEntry != nil {
			appState.outputEntry.SetText(appState.outputPath)
		}
	})

	// 设置文件过滤器 - 允许所有文件
	picker.SetFilter(func(name string) bool {
		return true
	})

	// 修改标题为"选择输出文件"
	picker.window.SetTitle("选择输出文件")

	picker.Show()
}

// generateOutputPath 生成默认输出路径（与输入文件同目录）
func generateOutputPath(inputPath string, fileType models.FileType) string {
	dir := filepath.Dir(inputPath)
	filename := filepath.Base(inputPath)
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	// 根据处理模式生成后缀
	suffix := "_压缩版"
	if appState.processMode == ipv6pkg.ModeExpand {
		suffix = "_扩展版"
	}

	// 检查文件是否已存在，如果存在则添加序号
	outputName := nameWithoutExt + suffix + ext
	outputPath := filepath.Join(dir, outputName)

	counter := 1
	for fileExists(outputPath) {
		outputName = fmt.Sprintf("%s%s_%d%s", nameWithoutExt, suffix, counter, ext)
		outputPath = filepath.Join(dir, outputName)
		counter++
	}

	return outputPath
}

// fileExists 检查文件是否存在
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// loadFile 加载文件
func loadFile(filePath string) {
	// 重置状态
	if appState.currentFile != nil {
		appState.currentFile.Close()
	}

	// 根据文件类型创建处理器
	fileType := models.DetectFileType(filePath)
	var proc processor.FileProcessor
	var err error

	switch fileType {
	case models.FileTypeExcel:
		proc, err = processor.NewExcelProcessor(filePath)
		// 如果 Excel 打开失败，尝试作为 CSV 处理
		if err != nil {
			appState.statusBar.SetText(fmt.Sprintf("Excel 打开失败，尝试作为 CSV 处理..."))
			proc, err = processor.NewCSVProcessor(filePath)
			if err != nil {
				dialog.ShowError(fmt.Errorf("无法打开文件（既不是有效的 Excel 也不是 CSV）: %w", err), appState.window)
				return
			}
			fileType = models.FileTypeCSV
		}
	case models.FileTypeCSV:
		proc, err = processor.NewCSVProcessor(filePath)
	default:
		// 未知文件类型，尝试 CSV
		proc, err = processor.NewCSVProcessor(filePath)
		if err != nil {
			dialog.ShowError(fmt.Errorf("不支持的文件类型: %w", err), appState.window)
			return
		}
		fileType = models.FileTypeCSV
	}

	if err != nil {
		dialog.ShowError(err, appState.window)
		return
	}

	appState.currentFile = proc
	appState.fileType = fileType
	appState.statusText = fmt.Sprintf("已加载: %s", filePath)

	// 设置默认输出路径（同目录下）
	defaultOutputPath := generateOutputPath(filePath, fileType)
	appState.outputPath = defaultOutputPath
	if appState.outputEntry != nil {
		appState.outputEntry.SetText(defaultOutputPath)
	}

	// 更新文件信息
	appState.fileInfoPanel.Update(
		filePath,
		fileType,
		proc.GetRowCount(),
		proc.GetColumnCount(),
	)

	// 检测 IPv6 列
	detectColumns()

	// 启用开始按钮
	appState.startBtn.Enable()

	// 更新状态
	appState.statusBar.SetText(fmt.Sprintf("已加载 %s (%d 行, %d 列)，选中 %d 个列进行处理",
		fileType.String(), proc.GetRowCount(), proc.GetColumnCount(), len(appState.selectedCols)))
}

// detectColumns 检测包含 IPv6 的列
func detectColumns() {
	rows := appState.currentFile.GetRows()

	// 检测列（采样1000行，降低阈值到5%，更敏感）
	appState.columns = processor.DetectIPv6Columns(rows, 1000, 5)

	// 创建列选择器
	if appState.columnSelector != nil {
		// 移除旧的
	}

	// 自动选中包含 IPv6 的列
	for _, col := range appState.columns {
		if col.IsIPv6 {
			appState.selectedCols[col.Name] = true
		}
	}

	// 如果没有检测到明显的 IPv6 列，自动选择所有列进行处理
	// 这样可以处理包含 IPv6 地址但比例较低的列
	if len(appState.selectedCols) == 0 && len(appState.columns) > 0 {
		for _, col := range appState.columns {
			appState.selectedCols[col.Name] = true
		}
	}

	// 更新 UI
	appState.statusBar.SetText(fmt.Sprintf("检测到 %d 个列，自动选中 %d 个列进行处理",
		len(appState.columns), len(appState.selectedCols)))
}

// startProcessing 开始处理
func startProcessing() {
	if appState.currentFile == nil {
		dialog.ShowError(fmt.Errorf("请先选择文件"), appState.window)
		return
	}

	if len(appState.selectedCols) == 0 {
		dialog.ShowError(fmt.Errorf("请至少选择一列"), appState.window)
		return
	}

	appState.isProcessing = true
	appState.startBtn.Disable()
	appState.progressBar.Reset()

	// 获取选中的列名
	selectedCols := make([]string, 0, len(appState.selectedCols))
	for name := range appState.selectedCols {
		selectedCols = append(selectedCols, name)
	}

	// 启动处理协程
	go processFile(selectedCols)
}

// processFile 处理文件
func processFile(selectedCols []string) {
	startTime := time.Now()
	totalRows := appState.currentFile.GetRowCount()
	rowsProcessed := 0

	// 处理每一列
	for _, colName := range selectedCols {
		appState.progressBar.SetText(fmt.Sprintf("正在处理列: %s (%d/%d)", colName, len(appState.selectedCols), len(selectedCols)))

		// 根据模式创建处理函数
		processFunc := func(ip string) string {
			return ipv6pkg.ProcessIPv6(ip, appState.processMode)
		}

		processed, err := appState.currentFile.ProcessColumnsByName(
			[]string{colName},
			processFunc,
		)

		if err != nil {
			appState.window.Canvas().Refresh(appState.progressBar.Container())
			dialog.ShowError(err, appState.window)
			appState.isProcessing = false
			appState.startBtn.Enable()
			return
		}

		rowsProcessed += processed

		// 更新进度
		progress := float64(rowsProcessed) / float64(len(selectedCols)*totalRows)
		appState.progressBar.SetValue(progress)
	}

	// 保存文件
	appState.progressBar.SetText("正在保存...")
	outputPath := appState.outputPath
	if outputPath == "" {
		outputPath = generateOutputPath(appState.currentFile.GetFilePath(), appState.fileType)
	}

	err := appState.currentFile.Save(outputPath)
	if err != nil {
		appState.window.Canvas().Refresh(appState.progressBar.Container())
		dialog.ShowError(err, appState.window)
		appState.isProcessing = false
		appState.startBtn.Enable()
		return
	}

	// 处理完成
	duration := time.Since(startTime).Milliseconds()

	appState.result = &models.ProcessingResult{
		Success:        true,
		ProcessedRows:  totalRows,
		ProcessedCount: rowsProcessed,
		OutputFile:     outputPath,
		Duration:       duration,
		ModifiedColumns: selectedCols,
	}

	// 更新 UI - 从主线程更新
	appState.isProcessing = false
	appState.startBtn.Enable()
	appState.progressBar.SetValue(1.0)
	appState.statusBar.SetText(fmt.Sprintf("处理完成！修改了 %d 个地址，耗时 %d ms", rowsProcessed, duration))

	ShowResultDialog(appState.result, appState.window)
}
