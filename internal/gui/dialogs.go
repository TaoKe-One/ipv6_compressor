package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/TaoKe-One/ipv6_compressor/pkg/models"
)

// ShowResultDialog 显示处理结果对话框
func ShowResultDialog(result *models.ProcessingResult, parent fyne.Window) {
	if result.Error != nil {
		dialog.ShowError(result.Error, parent)
		return
	}

	content := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("✓ 处理完成!")),
		container.NewGridWithColumns(2,
			widget.NewLabel("处理行数:"), widget.NewLabel(fmt.Sprintf("%d", result.ProcessedRows)),
			widget.NewLabel("修改数量:"), widget.NewLabel(fmt.Sprintf("%d", result.ProcessedCount)),
			widget.NewLabel("耗时:"), widget.NewLabel(fmt.Sprintf("%d ms", result.Duration)),
		),
		widget.NewSeparator(),
		widget.NewLabel(fmt.Sprintf("输出文件: %s", result.OutputFile)),
	)

	d := dialog.NewCustomConfirm("处理结果", "打开文件夹", "关闭", content, func(response bool) {
		if response {
			// TODO: 打开输出文件夹
		}
	}, parent)
	d.Show()
}

// ShowPreviewDialog 显示预览对话框
func ShowPreviewDialog(before, after [][]string, parent fyne.Window) {
	content := container.NewScroll(
		container.NewVBox(
			widget.NewLabel("处理前:"),
			createPreviewTable(before),
			widget.NewSeparator(),
			widget.NewLabel("处理后:"),
			createPreviewTable(after),
		),
	)

	d := dialog.NewCustomConfirm("预览", "关闭", "", content, func(bool) {}, parent)
	d.Resize(fyne.NewSize(600, 400))
	d.Show()
}

// createPreviewTable 创建预览表格
func createPreviewTable(rows [][]string) *fyne.Container {
	if len(rows) == 0 {
		return container.NewVBox(widget.NewLabel("(空)"))
	}

	// 只显示前5行
	maxRows := 5
	if len(rows) < maxRows {
		maxRows = len(rows)
	}

	var items []fyne.CanvasObject
	for i := 0; i < maxRows; i++ {
		rowText := ""
		for j, cell := range rows[i] {
			if j > 0 {
				rowText += " | "
			}
			rowText += cell
		}
		items = append(items, widget.NewLabel(rowText))
	}

	if len(rows) > maxRows {
		items = append(items, widget.NewLabel(fmt.Sprintf("(... 还有 %d 行)", len(rows)-maxRows)))
	}

	return container.NewVBox(items...)
}
