package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/themes"

	"github.com/TaoKe-One/ipv6-compressor/internal/gui"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(themes.DefaultTheme())

	w := a.NewWindow("IPv6 Compressor v2.0")
	w.SetMaster()

	gui.LoadUI(w)

	w.ShowAndRun()
}
