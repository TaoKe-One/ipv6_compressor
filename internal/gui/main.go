package gui

import (
	"fmt"
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
	window        fyne.Window
	fileInfoPanel *FileInfoPanel
	columnSelector *ColumnSelector
	progressBar   *ProgressBar
	statusBar     *StatusBar
	startBtn      *widget.Button
}

// appState 全局应用状态
var appState = &AppState{
	selectedCols: make(map[string]bool),
	statusText:   "准备就绪",
	processMode:  ipv6pkg.ModeCompress, // 默认压缩模式
}

// LoadUI 加载主界面
func LoadUI(w fyne.Window) {
	appState.window = w

	// 创建 UI 组件
	appState.fileInfoPanel = NewFileInfoPanel()
	appState.progressBar = NewProgressBar()
	appState.statusBar = NewStatusBar()
	appState.startBtn = widget.NewButton("开始处理", startProcessing)
	appState.startBtn.Disable()

	// 主内容
	content := createMainContent()

	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 600))
	w.CenterOnScreen()
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
	outputEntry.SetPlaceHolder("默认: 原文件名_compressed.扩展名")
	outputEntry.OnChanged = func(s string) {
		appState.outputPath = s
	}

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
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		defer reader.Close()

		filePath := reader.URI().Path()
		loadFile(filePath)
	}, nil)
}

// showOutputPicker 显示输出文件选择器
func showOutputPicker() {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		defer writer.Close()

		appState.outputPath = writer.URI().Path()
	}, nil)
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
	case models.FileTypeCSV:
		proc, err = processor.NewCSVProcessor(filePath)
	default:
		dialog.ShowError(fmt.Errorf("不支持的文件类型"), appState.window)
		return
	}

	if err != nil {
		dialog.ShowError(err, appState.window)
		return
	}

	appState.currentFile = proc
	appState.fileType = fileType
	appState.statusText = fmt.Sprintf("已加载: %s", filePath)

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
	appState.statusBar.SetText(fmt.Sprintf("已加载 %s，检测到 %d 个 IPv6 列",
		fileType.String(), len(appState.columns)))
}

// detectColumns 检测包含 IPv6 的列
func detectColumns() {
	rows := appState.currentFile.GetRows()

	// 检测列（采样1000行，30%阈值）
	appState.columns = processor.DetectIPv6Columns(rows, 1000, 30)

	// 创建列选择器
	if appState.columnSelector != nil {
		// 移除旧的
	}

	// 默认选中 IPv6 列
	for _, col := range appState.columns {
		if col.IsIPv6 {
			appState.selectedCols[col.Name] = true
		}
	}

	// 更新 UI（这里简化处理，实际需要更复杂的 UI 更新）
	appState.statusBar.SetText(fmt.Sprintf("检测到 %d 个列，%d 个包含 IPv6",
		len(appState.columns), len(processor.FilterIPv6Columns(appState.columns))))
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

	// 更新进度
	appState.progressBar.SetText(fmt.Sprintf("正在处理 %d 列...", len(selectedCols)))

	// 处理每一列
	for _, colName := range selectedCols {
		appState.progressBar.SetText(fmt.Sprintf("正在处理列: %s", colName))

		// 根据模式创建处理函数
		processFunc := func(ip string) string {
			return ipv6pkg.ProcessIPv6(ip, appState.processMode)
		}

		processed, err := appState.currentFile.ProcessColumnsByName(
			[]string{colName},
			processFunc,
		)

		if err != nil {
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
		outputPath = appState.currentFile.GetFilePath() + "_compressed"
	}

	err := appState.currentFile.Save(outputPath)
	if err != nil {
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

	// 更新 UI
	appState.isProcessing = false
	appState.startBtn.Enable()
	appState.progressBar.SetValue(1.0)
	appState.statusBar.SetText(fmt.Sprintf("处理完成！修改了 %d 个地址", rowsProcessed))

	ShowResultDialog(appState.result, appState.window)
}
