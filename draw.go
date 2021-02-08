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
// Zoom slider
const (
	defaultZoom       = uint32(hundred)
	stepString        = tenString
	defaultStep       = uint32(ten)
	defaultStepString = stepString
)

var (
	zoom        = defaultZoom
	zoomSlider  *widget.Slider
	sliderLabel *widget.Label
	stepSelect  *widget.Select
	step        = defaultStep
	stepLabel   *widget.Label
	xSup        *widget.Button
	xSDown      *widget.Button
	xLabel      *widget.Label
	ySup        *widget.Button
	ySDown      *widget.Button
	yLabel      *widget.Label
	zSup        *widget.Button
	zSDown      *widget.Button
	zLabel      *widget.Label
	xCluster    *fyne.Container
	yCluster    *fyne.Container
	zCluster    *fyne.Container

	// Star Screen
	origin = fyne.NewPos(0, 0)
	here   = position{50000.0, 50000.0, 12500.0}

	scale      = scaleFactor
	windowSize = fyne.NewSize(scale, scale)

	zooms = []named{
		{"1", 1, 0},
		{"3", 3, 1},
		{"10", 10, 2},
		{"32", 32, 3},
		{"100", 100, 4},
		{"200", 200, 5},
		{"300", 300, 6},
		{"400", 400, 7},
		{"500", 500, 8},
		{"750", 750, 9},
		{"1000", 1000, 10},
	}
)

func stepNameToValue(name string) uint32 {
	return nameToValue(name, defaultStep)
}

var zoomIndex = uint32(middleZoom)

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
	zoomSlider = &widget.Slider{
		Value:       0,
		Min:         ten,
		Max:         eightHundred,
		Step:        1,
		Orientation: widget.Horizontal,
		OnChanged:   sliderUpdate,
	}
	sliderLabel = widget.NewLabel("Select Zoom Factor\n   in light years\n       100")

	stepSelect = widget.NewSelect(zoomStrings, selectStep)
	stepSelect.Selected = defaultStepString
	stepLabel = widget.NewLabel("Select step size\n   in light years")

	xSup = widget.NewButton("Step X+", xSInc)
	xSDown = widget.NewButton("Step X-", xSDec)
	xLabel = widget.NewLabel(fmt.Sprintf("X position\n%f", here.x))
	xCluster = container.NewVBox(container.NewHBox(xSup, xSDown), xLabel)

	ySup = widget.NewButton("Step Y+", ySInc)
	ySDown = widget.NewButton("Step Y-", ySDec)
	yLabel = widget.NewLabel(fmt.Sprintf("Y position\n%f", here.y))
	yCluster = container.NewVBox(container.NewHBox(ySup, ySDown), yLabel)

	zSup = widget.NewButton("Step Z+", zSInc)
	zSDown = widget.NewButton("Step Z-", zSDec)
	zLabel = widget.NewLabel(fmt.Sprintf("Z position\n%f", here.z))
	zCluster = container.NewVBox(container.NewHBox(zSup, zSDown), zLabel)
}

func sliderUpdate(zoomFactor float64) {
	zoom = uint32(zoomFactor)
	sliderLabel.SetText(fmt.Sprint("Select Zoom Factor\n   in light years\n       ", zoom))
}

func selectStep(selection string) {
	step = stepNameToValue(selection)
	viewPort.Refresh()
}

func getCircle(fromStar star) *canvas.Circle {
	starCircle := canvas.NewCircle(fromStar.brightColor)

	starCircle.FillColor = fromStar.brightColor
	starCircle.StrokeColor = fromStar.brightColor
	starCircle.StrokeWidth = 1
	fudgePixels := float32(1.0)
	if zoomIndex > middleZoom {
		fudgePixels /= fudgeFactor
	}
	if zoomIndex > bigZoom {
		fudgePixels /= fudgeFactor
	}
	starCircle.Resize(fyne.NewSize(float32(fromStar.pixels)*fudgePixels, float32(fromStar.pixels)*fudgePixels))
	starCircle.Move(fyne.NewPos(fromStar.dx, fromStar.dy))

	return starCircle
}

var viewPort = canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: opaque})

func (s *starLayout) Layout(_ []fyne.CanvasObject, size fyne.Size) {
	s.makeStarContainer(viewPort)
}

func (s *starLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return windowSize
}

func (s *starLayout) render() *fyne.Container {
	return s.makeStarContainer(viewPort)
}

var leftHandSide *fyne.Container

func Show(win fyne.Window) fyne.CanvasObject {
	sector := &starLayout{}
	sector.canvas = &canvas.Rectangle{
		FillColor:   color.Black,
		StrokeColor: color.Black,
		StrokeWidth: ten,
	}
	sector.canvas.Move(origin)
	sector.canvas.Resize(windowSize)
	sector.canvas.Refresh()
	starContent := sector.render()
	leftHandSide = container.NewVBox(xCluster, yCluster, zCluster, stepLabel,
		fyne.Widget(stepSelect),
		sliderLabel, zoomSlider)
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
	here.x += float32(step)
	if here.x > xMax {
		here.x = xMax
	}
	xLabel.SetText(fmt.Sprintf("X position\n%f", here.x))
	leftHandSide.Refresh()
	Show(*window)
}

func xSDec() {
	here.x -= float32(step)
	if here.x < xMin {
		here.x = xMin
	}
	xLabel.SetText(fmt.Sprintf("X position\n%f", here.x))
	leftHandSide.Refresh()
	Show(*window)
}

func ySInc() {
	here.y += float32(step)
	if here.y > yMax {
		here.y = yMax
	}
	yLabel.SetText(fmt.Sprintf("Y position\n%f", here.y))
	leftHandSide.Refresh()
	Show(*window)
}

func ySDec() {
	here.y -= float32(step)
	if here.y < yMin {
		here.y = yMin
	}
	yLabel.SetText(fmt.Sprintf("Y position\n%f", here.y))
	leftHandSide.Refresh()
	Show(*window)
}

func zSInc() {
	here.z += float32(step)
	if here.z > zMax {
		here.z = zMax
	}
	zLabel.SetText(fmt.Sprintf("Z position\n%f", here.z))
	leftHandSide.Refresh()
	Show(*window)
}

func zSDec() {
	here.z -= float32(step)
	if here.z < zMin {
		here.z = zMin
	}
	zLabel.SetText(fmt.Sprintf("Z position\n%f", here.z))
	leftHandSide.Refresh()
	Show(*window)
}
