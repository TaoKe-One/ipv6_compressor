package main

import (
	"fyne.io/fyne/v2/app"

	"github.com/TaoKe-One/ipv6_compressor/internal/gui"
)

func main() {
	a := app.NewWithID("com.taokeone.ipv6compressor")

	w := a.NewWindow("IPv6 Compressor v2.0")
	w.SetMaster()

	gui.LoadUI(w)

	w.ShowAndRun()
}
