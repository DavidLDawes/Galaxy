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
	controlsInit(&w)
	w.Resize(fyne.NewSize(thousand, thousand))
	w.SetContent(Show(w))
	w.ShowAndRun()
}
