package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type named struct {
	name  string
	value uint32
	index uint32
}

// Controls
// Zoom selector
const (
	zoomString        = "100"
	defaultZoom       = uint32(100)
	defaultZoomString = zoomString
	stepString        = "10"
	defaultStep       = uint32(10)
	defaultStepString = stepString
)

var (
	zoomSelect *widget.Select
	zoom    = defaultZoom
	zoomLabel *widget.Label
	stepSelect *widget.Select
	step    = defaultStep
	stepLabel *widget.Label
	xSup       *widget.Button
	xSdown     *widget.Button
	xLabel    *widget.Label
	ySup       *widget.Button
	ySdown     *widget.Button
	yLabel    *widget.Label
	zup       *widget.Button
	zdown     *widget.Button
	zSup       *widget.Button
	zSdown     *widget.Button
	zLabel    *widget.Label
	xCluster  *fyne.Container
	yCluster  *fyne.Container
	zCluster  *fyne.Container

	// Star Screen
	origin        = fyne.Position{0, 0}
	sectorOrigin  = sector{0, 0, 0}
	currentSector = sector{50000, 50000, 12500}
	here          = position{50000.0, 50000.0, 12500.0}

	scale               = float32(1000.0)
	lightyearsPerSector = uint32(100)
	windowSize          = fyne.NewSize(scale, scale)
	scaleFactor         = float32(scale) / float32(lightyearsPerSector)

	zooms []named = []named{
		{"1", 1, 0},
		{"3", 3, 1},
		{"10", 10, 2},
		{"32", 32, 3},
		{"100", 100, 4},
		{"316", 316, 5},
		{"1000", 1000, 6},
		{"3162", 3162, 7},
	}
)


func zoomNameToValue(name string) uint32 {
	return nameToValue(name, defaultZoom)
}


func stepNameToValue(name string) uint32 {
	return nameToValue(name, uint32(defaultStep))
}

var zoomIndex = uint32(4)

func nameToValue(name string, defaultValue uint32) uint32 {
	for _, nextZoom := range zooms {
		if name == nextZoom.name {
			zoomIndex = nextZoom.index
			return nextZoom.value
		}
	}
	zoomIndex = 4
	return defaultValue
}

type starLayout struct {
	starCircles []fyne.CanvasObject
	canvas      fyne.CanvasObject
	stop        bool
}

var window *fyne.Window

func controlsInit(inWindow *fyne.Window) {
	window = inWindow
	zoomStrings := make([]string, len(zooms))
	for i := 0; i < len(zooms); i++ {
		zoomStrings[i] = zooms[i].name
	}
	zoomSelect = widget.NewSelect(zoomStrings, selectZoom)
	zoomSelect.Selected = defaultZoomString
	zoomLabel = widget.NewLabel("Select window size\n   in light years")

	stepSelect = widget.NewSelect(zoomStrings, selectStep)
	stepSelect.Selected = defaultStepString
	stepLabel = widget.NewLabel("Select step size\n   in light years")

	xSup = widget.NewButton("Step X+", xSInc)
	xSdown = widget.NewButton("Step X-", xSDec)
	xLabel = widget.NewLabel(fmt.Sprintf("X position\n%f", here.x))
	xCluster = container.NewVBox(container.NewHBox(xSup, xSdown), xLabel)

	ySup = widget.NewButton("Step Y+", ySInc)
	ySdown = widget.NewButton("Step Y-", ySDec)
	yLabel = widget.NewLabel(fmt.Sprintf("Y position\n%f", here.y))
	yCluster = container.NewVBox(container.NewHBox(ySup, ySdown), yLabel)

	zSup = widget.NewButton("Step Z+", zSInc)
	zSdown = widget.NewButton("Step Z-", zSDec)
	zLabel = widget.NewLabel(fmt.Sprintf("Z position\n%f", here.z))
	zCluster = container.NewVBox(container.NewHBox(zSup, zSdown), zLabel)

	setPosition()
}

func selectZoom(selection string) {
	zoom = zoomNameToValue(selection)
	Show(*window)
}

func selectStep(selection string) {
	step = stepNameToValue(selection)
	viewPort.Refresh()
}

func getCircle(fromStar star) *canvas.Circle {
	starCircle := canvas.NewCircle(fromStar.brightcolor)

	starCircle.FillColor = fromStar.brightcolor
	starCircle.StrokeColor = fromStar.brightcolor
	starCircle.StrokeWidth = 1
	fudgePixels := float32(1.0)
	if zoomIndex > 4 {
		fudgePixels = fudgePixels / 2.0
	}
	if zoomIndex > 5 {
		fudgePixels = fudgePixels / 2.0
	}
	if zoomIndex > 6 {
		fudgePixels = fudgePixels / 2.0
	}
	starCircle.Resize(fyne.NewSize(float32(fromStar.pixels)*fudgePixels, float32(fromStar.pixels)*fudgePixels))
	starCircle.Move(fyne.NewPos(fromStar.dx, fromStar.dy))

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
	leftHandSide := container.NewVBox(xCluster, yCluster, zCluster, stepLabel, fyne.Widget(stepSelect), zoomLabel, fyne.Widget(zoomSelect))
	screen := container.NewHBox(leftHandSide, starContent)
	win.SetContent(screen)
	return screen
}

func (s *starLayout) makeStarContainer(rectangle *canvas.Rectangle) *fyne.Container {
	s.starCircles = make([]fyne.CanvasObject, 0)
	starContainer := container.NewWithoutLayout()
	starContainer.Resize(windowSize)
	starContainer.Move(origin)
	rectangle.Resize(windowSize)
	rectangle.Move(origin)
	starContainer.Objects = append(starContainer.Objects, rectangle)
	for _, star := range getGalaxyDetails() {
		nextCircle := *getCircle(star)
		s.starCircles = append(s.starCircles, &nextCircle)
		starContainer.Objects = append(starContainer.Objects, &nextCircle)
	}
	s.canvas = starContainer
	starContainer.Layout = s
	starContainer.Show()

	return starContainer
}

func xSInc() {
	here.x = here.x + float32(step)
	if here.x > xmax {
		here.x = xmax
	}
	setPosition()
	xLabel.SetText(fmt.Sprintf("X position\n%f", here.x))
	Show(*window)
}

func xSDec() {
	here.x = here.x - float32(step)
	if here.x < 0 {
		here.x = 0
	}
	setPosition()
	xLabel.SetText(fmt.Sprintf("X position\n%f", here.x))
	Show(*window)
}

func ySInc() {
	here.y = here.y + float32(step)
	if here.y > ymax {
		here.y = ymax
	}
	setPosition()
	yLabel.SetText(fmt.Sprintf("Y position\n%f", here.x))
	Show(*window)
}

func ySDec() {
	here.y = here.y - float32(step)
	if here.y < 0 {
		here.y = 0
	}
	setPosition()
	yLabel.SetText(fmt.Sprintf("Y position\n%f", here.x))
	Show(*window)
}

func zSInc() {
	here.z = here.z + float32(step)
	if here.z > zmax {
		here.z = zmax
	}
	setPosition()
	zLabel.SetText(fmt.Sprintf("Z position\n%f", here.x))
	Show(*window)
}

func zSDec() {
	here.z = here.z - float32(step)
	if here.z < 0 {
		here.z = 0
	}
	setPosition()
	zLabel.SetText(fmt.Sprintf("Z position\n%f", here.x))
	Show(*window)
}

func setPosition() {
	currentSector = sector{
		lightyearsPerSector * uint32(here.x/float32(lightyearsPerSector)),
		lightyearsPerSector * uint32(here.y/float32(lightyearsPerSector)),
		lightyearsPerSector * uint32(here.z/float32(lightyearsPerSector)),
	}
}
