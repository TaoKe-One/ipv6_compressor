package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/TaoKe-One/ipv6-compressor/internal/processor"
	"github.com/TaoKe-One/ipv6-compressor/pkg/models"
)

// ColumnSelector 列选择器组件
type ColumnSelector struct {
	container *fyne.Container
	checks    []*widget.Check
	columns   []processor.ColumnInfo
	onChange  func([]string)
}

// NewColumnSelector 创建列选择器
func NewColumnSelector(columns []processor.ColumnInfo, onChange func([]string)) *ColumnSelector {
	cs := &ColumnSelector{
		columns:  columns,
		onChange: onChange,
	}

	var checks []*widget.Check
	for _, col := range columns {
		col := col
		check := widget.NewCheck(col.Name, func(checked bool) {
			if checked {
				appState.selectedCols[col.Name] = true
			} else {
				delete(appState.selectedCols, col.Name)
			}
			cs.updateSelected()
		})
		if col.IsIPv6 {
			check.Checked = true
			appState.selectedCols[col.Name] = true
		}
		checks = append(checks, check)
	}
	cs.checks = checks

	// 创建滚动容器
	items := make([]fyne.CanvasObject, len(checks))
	for i, check := range checks {
		items[i] = check
	}

	scroll := container.NewScroll(container.NewVBox(items...))
	scroll.SetMinSize(fyne.NewSize(200, 150))

	cs.container = container.NewBorder(
		nil, nil, nil, nil,
		scroll,
	)

	return cs
}

// updateSelected 更新选中状态
func (cs *ColumnSelector) updateSelected() {
	if cs.onChange != nil {
		selected := make([]string, 0, len(appState.selectedCols))
		for name := range appState.selectedCols {
			selected = append(selected, name)
		}
		cs.onChange(selected)
	}
}

// ProgressBar 进度条组件
type ProgressBar struct {
	bar    *widget.ProgressBar
	label  *widget.Label
	value  float64
	text   string
}

// NewProgressBar 创建进度条
func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		bar:   widget.NewProgressBar(),
		label: widget.NewLabel(""),
	}
}

// SetValue 设置进度值
func (pb *ProgressBar) SetValue(value float64) {
	pb.value = value
	pb.bar.SetValue(value)
}

// SetText 设置进度文本
func (pb *ProgressBar) SetText(text string) {
	pb.text = text
	pb.label.SetText(text)
}

// Container 获取容器
func (pb *ProgressBar) Container() *fyne.Container {
	return container.NewVBox(
		pb.label,
		pb.bar,
	)
}

// Reset 重置进度条
func (pb *ProgressBar) Reset() {
	pb.SetValue(0)
	pb.SetText("")
}

// StatusBar 状态栏组件
type StatusBar struct {
	label *widget.Label
}

// NewStatusBar 创建状态栏
func NewStatusBar() *StatusBar {
	return &StatusBar{
		label: widget.NewLabel("准备就绪"),
	}
}

// SetText 设置状态文本
func (sb *StatusBar) SetText(text string) {
	sb.label.SetText(text)
	appState.statusText = text
}

// Container 获取容器
func (sb *StatusBar) Container() *fyne.Container {
	return container.NewPadded(sb.label)
}

// FileInfoPanel 文件信息面板
type FileInfoPanel struct {
	container *fyne.Container
	fileName  *widget.Label
	fileType  *widget.Label
	rowCount  *widget.Label
	colCount  *widget.Label
}

// NewFileInfoPanel 创建文件信息面板
func NewFileInfoPanel() *FileInfoPanel {
	fip := &FileInfoPanel{
		fileName: widget.NewLabel("-"),
		fileType: widget.NewLabel("-"),
		rowCount: widget.NewLabel("-"),
		colCount: widget.NewLabel("-"),
	}

	fip.container = container.NewVBox(
		widget.NewLabel("文件信息:"),
		container.NewGridWithColumns(2,
			widget.NewLabel("文件名:"), fip.fileName,
			widget.NewLabel("类型:"), fip.fileType,
			widget.NewLabel("行数:"), fip.rowCount,
			widget.NewLabel("列数:"), fip.colCount,
		),
	)

	return fip
}

// Update 更新文件信息
func (fip *FileInfoPanel) Update(filePath string, fileType models.FileType, rowCount, colCount int) {
	fip.fileName.SetText(filePath)
	fip.fileType.SetText(fileType.String())
	fip.rowCount.SetText(fmt.Sprintf("%d", rowCount))
	fip.colCount.SetText(fmt.Sprintf("%d", colCount))
}

// Clear 清空信息
func (fip *FileInfoPanel) Clear() {
	fip.fileName.SetText("-")
	fip.fileType.SetText("-")
	fip.rowCount.SetText("-")
	fip.colCount.SetText("-")
}

// Container 获取容器
func (fip *FileInfoPanel) Container() *fyne.Container {
	return fip.container
}
