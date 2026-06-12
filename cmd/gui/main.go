package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/TaoKe-One/ipv6_compressor/internal/gui"
)

const (
	windowWidthKey  = "window_width"
	windowHeightKey = "window_height"
	defaultWidth    = 700  // 默认宽度（适配高 DPI）
	defaultHeight   = 500  // 默认高度（适配高 DPI）
)

func main() {
	a := app.NewWithID("com.taokeone.ipv6compressor")

	w := a.NewWindow("IPv6 Compressor v2.1.3")
	w.SetMaster()

	// 从 Preferences 恢复窗口大小，或使用默认值
	prefs := a.Preferences()
	width := float32(defaultWidth)
	height := float32(defaultHeight)

	if w := prefs.Float(windowWidthKey); w > 0 {
		width = float32(w)
	}
	if h := prefs.Float(windowHeightKey); h > 0 {
		height = float32(h)
	}

	// 设置窗口大小
	w.Resize(fyne.NewSize(width, height))

	// 在窗口关闭时保存大小
	w.SetCloseIntercept(func() {
		prefs.SetFloat(windowWidthKey, float64(w.Canvas().Size().Width))
		prefs.SetFloat(windowHeightKey, float64(w.Canvas().Size().Height))
		w.Close()
	})

	gui.LoadUI(w, a)

	w.ShowAndRun()
}
