package gui

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// FilePicker 自定义文件选择器窗口
type FilePicker struct {
	parentWindow fyne.Window
	window       fyne.Window
	onSelected   func(string)
	filter       func(string) bool
	currentDir   string
	selectedFile string
	fileList     *widget.List
	pathLabel    *widget.Label
	dirs         []string
	files        []string
	app          fyne.App
	title        string
}

// NewFilePicker 创建文件选择器
func NewFilePicker(parent fyne.Window, app fyne.App, onSelected func(string)) *FilePicker {
	return &FilePicker{
		parentWindow: parent,
		app:          app,
		onSelected:   onSelected,
		filter:       func(s string) bool { return true },
	}
}

// SetFilter 设置文件过滤器
func (fp *FilePicker) SetFilter(filter func(string) bool) {
	fp.filter = filter
}

// SetTitle 设置窗口标题
func (fp *FilePicker) SetTitle(title string) {
	fp.title = title
}

// Show 显示文件选择器
func (fp *FilePicker) Show() {
	// 创建新窗口（可调整大小）
	title := fp.title
	if title == "" {
		title = "选择文件"
	}
	fp.window = fp.app.NewWindow(title)
	fp.window.Resize(fyne.NewSize(700, 500))
	fp.window.CenterOnScreen()

	// 当前目录标签
	fp.pathLabel = widget.NewLabel("当前目录:")

	// 文件列表
	fp.fileList = widget.NewList(
		func() int {
			fp.refreshFiles()
			return len(fp.dirs) + len(fp.files)
		},
		func() fyne.CanvasObject {
			icon := widget.NewIcon(theme.FolderIcon())
			label := widget.NewLabel("Template")
			return container.NewHBox(icon, label)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			cont := obj.(*fyne.Container)
			icon := cont.Objects[0].(*widget.Icon)
			label := cont.Objects[1].(*widget.Label)

			if id < len(fp.dirs) {
				// 显示目录
				icon.SetResource(theme.FolderIcon())
				label.SetText(fp.dirs[id] + "/")
			} else {
				// 显示文件
				fileIdx := id - len(fp.dirs)
				if fileIdx < len(fp.files) {
					icon.SetResource(theme.FileIcon())
					label.SetText(fp.files[fileIdx])
				}
			}
		},
	)

	// 点击选择文件
	fp.fileList.OnSelected = func(id widget.ListItemID) {
		if id < len(fp.dirs) {
			// 点击目录，进入
			fp.enterDir(fp.dirs[id])
			fp.fileList.UnselectAll() // 取消选中，允许再次点击
		} else {
			fileIdx := id - len(fp.dirs)
			if fileIdx < len(fp.files) {
				// 点击文件，确认选择
				fp.selectFile(fp.files[fileIdx])
			}
		}
	}

	// 按钮区域
	upBtn := widget.NewButton("上级目录", func() {
		fp.goUp()
	})

	homeBtn := widget.NewButton("主目录", func() {
		fp.goHome()
	})

	cancelBtn := widget.NewButton("取消", func() {
		fp.window.Close()
	})

	// 主布局
	buttonBox := container.NewHBox(upBtn, homeBtn)
	content := container.NewBorder(
		container.NewVBox(fp.pathLabel, buttonBox),
		container.NewHBox(layout.NewSpacer(), cancelBtn),
		nil, nil,
		fp.fileList,
	)

	fp.window.SetContent(content)

	// 初始化到用户主目录或上次目录
	startDir := os.Getenv("HOME")
	if startDir == "" {
		startDir = os.Getenv("USERPROFILE")
	}
	if startDir == "" {
		startDir = "."
	}
	fp.setCurrentDir(startDir)

	fp.window.Show()
}

// setCurrentDir 设置当前目录
func (fp *FilePicker) setCurrentDir(dir string) {
	fp.currentDir = dir
	fp.pathLabel.SetText("当前目录: " + dir)
	if fp.fileList != nil {
		fp.fileList.Refresh()
		fp.fileList.UnselectAll()
	}
}

// refreshFiles 刷新文件列表
func (fp *FilePicker) refreshFiles() {
	fp.dirs = nil
	fp.files = nil

	entries, err := os.ReadDir(fp.currentDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			if name != "." {
				fp.dirs = append(fp.dirs, name)
			}
		} else {
			if fp.filter(name) {
				fp.files = append(fp.files, name)
			}
		}
	}
}

// enterDir 进入目录
func (fp *FilePicker) enterDir(name string) {
	newDir := filepath.Join(fp.currentDir, name)
	fp.setCurrentDir(newDir)
}

// goUp 返回上级目录
func (fp *FilePicker) goUp() {
	parentDir := filepath.Dir(fp.currentDir)
	if parentDir != fp.currentDir {
		fp.setCurrentDir(parentDir)
	}
}

// goHome 跳转到主目录
func (fp *FilePicker) goHome() {
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home != "" {
		fp.setCurrentDir(home)
	}
}

// selectFile 选择文件
func (fp *FilePicker) selectFile(name string) {
	filePath := filepath.Join(fp.currentDir, name)
	fp.window.Close()
	if fp.onSelected != nil {
		fp.onSelected(filePath)
	}
}
