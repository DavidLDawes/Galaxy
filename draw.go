package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
)
// Controls
// Zoom selector
var zoomSelect *widget.Select

var defaultZoom = int32(100)
var defaultZoomString = "100"
var zoom = defaultZoom
var zoomString = defaultZoomString

// Star Screen
var origin = fyne.Position{0,0}
var sectorOrigin = sector{1, 1, 1}

var scale = float32(1000.0)
var lightyearsPerSector = 100
var windowSize = fyne.NewSize(scale, scale)
var scaleFactor = float32(scale)/float32(lightyearsPerSector)

type named struct {
	name string
	value int32
}

var zooms []named = []named{
	{"1", 1}, {"3", 3},
	{"10", 10}, {"32", 32},
	{"100", 100}, {"316", 316},
	{"1000", 1000}, {"3162", 3162},
	{"10000", 10000},

}

func zoomNameToValue(name string) int32 {
	for _, nextZoom := range zooms {
		if name == nextZoom.name {
			return nextZoom.value
		}
	}
	return defaultZoom
}

type starLayout struct {
	starCircles []fyne.CanvasObject
	canvas      fyne.CanvasObject
	stop        bool
}

func controlsInit() {
	zoomStrings := make([]string, len(zooms))
		for i :=0; i < len(zooms); i++ {
			zoomStrings[i] = zooms[i].name

	}
	zoomSelect = widget.NewSelect(zoomStrings, selectZoom)
	zoomSelect.Selected = "100"

}

func selectZoom(selection string) {
	zoom = zoomNameToValue(selection)
	viewPort.Refresh()
}

func getCircle(fromStar star) *canvas.Circle {
	starCircle := canvas.NewCircle(fromStar.brightcolor)

	starCircle.FillColor = fromStar.brightcolor
	starCircle.StrokeColor = fromStar.brightcolor
	starCircle.StrokeWidth = 1
	starCircle.Resize(fyne.NewSize(float32(fromStar.pixels*2), float32(fromStar.pixels*2),))
	starCircle.Move(fyne.NewPos(scaleFactor*fromStar.sx, scaleFactor*fromStar.sy))

	return starCircle
}

var viewPort = canvas.NewRectangle(color.RGBA{0, 0, 0, 255})

func (s *starLayout) Layout(_ []fyne.CanvasObject, size fyne.Size) {
	s.makeStarContainer(viewPort)
}

func (s *starLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return windowSize
}

func (s *starLayout) render() *fyne.Container {
	return s.makeStarContainer(viewPort)
}

func Show(win fyne.Window) fyne.CanvasObject {
	sector := &starLayout{}
	sector.canvas = &canvas.Rectangle{
		FillColor:   color.Black,
		StrokeColor: color.Black,
		StrokeWidth: 10,
	}
	sector.canvas.Move(origin)
	sector.canvas.Resize(windowSize)
	sector.canvas.Refresh()
	starContent := sector.render()
	screen := container.NewHBox(fyne.Widget(zoomSelect), starContent)
	win.SetContent(screen)
	return screen
}

func (s *starLayout) makeStarContainer(rectangle *canvas.Rectangle) (*fyne.Container) {
	s.starCircles = make([]fyne.CanvasObject, 0)
	starContainer := container.NewWithoutLayout()
	starContainer.Resize(windowSize)
	starContainer.Move(origin)
	rectangle.Resize(windowSize)
	rectangle.Move(origin)
	starContainer.Objects = append(starContainer.Objects, rectangle)
	for _, star := range getSectorDetails(sectorOrigin) {
		nextCircle := *getCircle(star)
		s.starCircles = append(s.starCircles, &nextCircle)
		starContainer.Objects = append(starContainer.Objects, &nextCircle)
	}
	s.canvas = starContainer
	starContainer.Layout = s
	starContainer.Show()

	return starContainer
}