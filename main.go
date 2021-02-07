package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("Galaxy")
	controlsInit()
	w.Resize(fyne.NewSize(1000, 1000))
	w.SetContent(Show(w))
	w.ShowAndRun()
}
