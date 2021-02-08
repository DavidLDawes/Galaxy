package main

import (
	"encoding/binary"
	"image/color"
	"math"
	"math/rand"

	"github.com/spaolacci/murmur3"
)

const (
	xMin = 0
	xMax = 100000
	yMin = 0
	yMax = 100000
	zMin = 0
	zMax = 25000
	xAdj = 100
	yAdj = 100
	zAdj = 100
)

type classDetails struct {
	class       string
	brightColor color.RGBA
	medColor    color.RGBA
	dimColor    color.RGBA
	odds        float32
	fudge       float32
	minMass     float32
	deltaMass   float32
	minRadii    float32
	deltaRadii  float32
	minLum      float32
	deltaLum    float32
	pixels      int32
}

type star struct {
	class       string
	brightColor color.RGBA
	medColor    color.RGBA
	dimColor    color.RGBA
	pixels      int32
	mass        float32
	radii       float32
	luminance   float32
	// 3D position
	x float32
	y float32
	z float32
	// sector location, 0 <= sx, sy, sz < 100
	sx float32
	sy float32
	sz float32
	// display location, 0 <= dx, dy, dz < 1000
	dx float32
	dy float32
	dz float32
}

type sector struct {
	x uint32
	y uint32
	z uint32
}

type position struct {
	x float32
	y float32
	z float32
}

var (
	bright = uint8(math.MaxUint8)
	tween  = uint8(sevenEighths)
	med    = uint8(threeQuarters)
	dim    = uint8(half)

	classO = classDetails{
		class:       "O",
		brightColor: color.RGBA{R: 0, G: 0, B: bright, A: opaque},
		medColor:    color.RGBA{R: 0, G: 0, B: bright, A: opaque},
		dimColor:    color.RGBA{R: 0, G: 0, B: med, A: opaque},
		odds:        .0000003,
		fudge:       .0000000402,
		minMass:     16.00001,
		deltaMass:   243.2,
		minRadii:    6,
		deltaRadii:  17.3,
		minLum:      30000,
		deltaLum:    147000.2,
		pixels:      11,
	}

	classB = classDetails{
		class:       "B",
		brightColor: color.RGBA{R: dim, G: dim, B: bright, A: opaque},
		medColor:    color.RGBA{R: dim / fudgeFactor, G: dim / fudgeFactor, B: med, A: opaque},
		dimColor: color.RGBA{
			R: dim / (fudgeFactor * fudgeFactor),
			G: dim / (fudgeFactor * fudgeFactor), B: dim, A: opaque,
		},
		odds:       .0013,
		fudge:      .0003,
		minMass:    2.1,
		deltaMass:  13.9,
		minRadii:   1.8,
		deltaRadii: 4.8,
		minLum:     25,
		deltaLum:   29975,
		pixels:     8,
	}

	classA = classDetails{
		class:       "A",
		brightColor: color.RGBA{R: bright, G: bright, B: bright, A: opaque},
		medColor:    color.RGBA{R: med, G: med, B: med, A: opaque},
		dimColor:    color.RGBA{R: dim, G: dim, B: dim, A: opaque},
		odds:        .006,
		fudge:       .0018,
		minMass:     1.4,
		deltaMass:   .7,
		minRadii:    1.4,
		deltaRadii:  .4,
		minLum:      5,
		deltaLum:    20,
		pixels:      6,
	}

	classF = classDetails{
		class:       "F",
		brightColor: color.RGBA{R: bright, G: bright, B: tween, A: opaque},
		medColor:    color.RGBA{R: tween, G: tween, B: dim, A: opaque},
		dimColor:    color.RGBA{R: med, G: med, B: dim / fudgeFactor, A: opaque},
		odds:        .03,
		fudge:       .012,
		minMass:     1.04,
		deltaMass:   .36,
		minRadii:    1.15,
		deltaRadii:  .25,
		minLum:      1.5,
		deltaLum:    3.5,
		pixels:      5,
	}

	classG = classDetails{
		class:       "G",
		brightColor: color.RGBA{R: tween, G: tween, B: 0, A: opaque},
		medColor:    color.RGBA{R: med, G: med, B: 0, A: opaque},
		dimColor:    color.RGBA{R: dim, G: dim, B: 0, A: opaque},
		odds:        .076,
		fudge:       .01102,
		minMass:     .8,
		deltaMass:   .24,
		minRadii:    .96,
		deltaRadii:  .19,
		minLum:      .6,
		deltaLum:    .9,
		pixels:      4,
	}

	classK = classDetails{
		class:       "K",
		brightColor: color.RGBA{R: 0xFE, G: 0xD8, B: 0xB1, A: opaque},
		medColor:    color.RGBA{R: 3 * (0xFE / 4), G: 3 * (0xD8 / 4), B: 3 * (0xB1 / 4), A: opaque},
		dimColor:    color.RGBA{R: 0xFE / 2, G: uint8(0xD8) / fudgeFactor, B: uint8(0xB1) / fudgeFactor, A: opaque},
		odds:        .121,
		fudge:       .042,
		minMass:     .45,
		deltaMass:   .35,
		minRadii:    .7,
		deltaRadii:  .26,
		minLum:      .08,
		deltaLum:    .52,
		pixels:      3,
	}

	classM = classDetails{
		class:       "M",
		brightColor: color.RGBA{R: bright, G: 0, B: 0, A: opaque},
		medColor:    color.RGBA{R: med, G: 0, B: 0, A: opaque},
		dimColor:    color.RGBA{R: dim, G: 0, B: 0, A: opaque},
		odds:        .7645,
		fudge:       .04,
		minMass:     1.04,
		deltaMass:   .36,
		minRadii:    1.15,
		deltaRadii:  .25,
		minLum:      1.5,
		deltaLum:    3.5,
		pixels:      2,
	}

	starDetailsByClass = [7]classDetails{classO, classB, classA, classF, classG, classK, classM}
	classByZoom        = [11]int{7, 7, 7, 7, 7, 7, 6, 5, 4, 3, 2}
)

func getStarDetails(classDetails classDetails, sector sector, random1m *rand.Rand) []star {
	stars := make([]star, 0)
	loopSize := int32(423.728813559 * (classDetails.odds - classDetails.fudge + 2*classDetails.fudge*random1m.Float32()))
	for i := 0; i < int(loopSize); i++ {
		nextStar := star{}
		random1 := random1m.Float32()
		nextStar.sx = random1m.Float32() * xAdj
		nextStar.sy = random1m.Float32() * yAdj
		nextStar.sz = random1m.Float32() * zAdj
		nextStar.x = float32(sector.x)*xAdj + nextStar.sx
		nextStar.y = float32(sector.y)*yAdj + nextStar.sy
		nextStar.z = float32(sector.z)*zAdj + nextStar.sz
		nextStar.class = classDetails.class
		nextStar.brightColor = classDetails.brightColor
		nextStar.dimColor = classDetails.dimColor
		nextStar.mass = classDetails.minMass + classDetails.deltaMass*(1+random1)
		nextStar.radii = classDetails.minRadii + random1*classDetails.deltaRadii
		nextStar.luminance = classDetails.minLum + random1*classDetails.deltaLum
		nextStar.pixels = classDetails.pixels
		stars = append(stars, nextStar)
	}

	return stars
}

func getSectorDetails(fromSector sector) []star {
	result := make([]star, 0)
	random1m := getHash(fromSector)
	classCount := 0
	for _, starDetails := range starDetailsByClass {
		nextClass := getStarDetails(starDetails, fromSector, random1m)
		result = append(result, nextClass...)
		classCount++
		if classCount > classByZoom[zoomIndex] {
			break
		}
	}

	return result
}

type saveDetails struct {
	savePosition position
	saveZoom     uint32
	saveStars    []star
	init         bool
}

var lastDetails = saveDetails{
	position{50000, 50000, 12500},
	100, make([]star, 0), false,
}

func getGalaxyDetails() []star {
	if lastDetails.init {
		if here.x == lastDetails.savePosition.x &&
			here.y == lastDetails.savePosition.y &&
			here.z == lastDetails.savePosition.z &&
			zoom == lastDetails.saveZoom {
			return lastDetails.saveStars
		}
	}
	xMin := here.x
	yMin := here.y
	zMin := here.z
	xMax := here.x + float32(zoom)
	yMax := here.y + float32(zoom)
	zMax := here.z + float32(zoom)

	resultStars := make([]star, 0)

	var extraX, extraY uint32
	if math.Mod(float64(xMin), float64(scale)) != 0 {
		extraX = 1
	} else {
		extraX = 0
	}
	if math.Mod(float64(yMin), float64(scale)) != 0 {
		extraY = 1
	} else {
		extraY = 0
	}
	for xi := uint32(0); 100*xi < uint32(xMax-xMin)+extraX; xi++ {
		for yi := uint32(0); 100*yi < uint32(yMax-yMin)+extraY; yi++ {
			for zi := uint32(0); 100*zi < uint32(zMax-zMin); zi++ {
				for _, star := range getSectorDetails(getSectorFromPosition(
					position{
						here.x + 100.0*float32(xi),
						here.y + 100.0*float32(yi),
						here.z + 100.0*float32(zi),
					})) {
					if !(star.x < xMin) && !(star.x > xMax) &&
						!(star.y < yMin) && !(star.y > yMax) &&
						!(star.z < zMin) && !(star.z > zMax) {
						star.dx = (star.x - xMin) * 1000 / (xMax - xMin)
						star.dy = (star.y - yMin) * 1000 / (yMax - yMin)
						star.dz = (star.z - zMin) * 1000 / (zMax - zMin)
						resultStars = append(resultStars, star)
					}
				}
			}
		}
	}
	lastDetails.savePosition.x = here.x
	lastDetails.savePosition.y = here.y
	lastDetails.savePosition.z = here.z
	lastDetails.saveZoom = zoom
	lastDetails.saveStars = resultStars
	lastDetails.init = true

	return resultStars
}

func getSectorFromPosition(now position) sector {
	return sector{uint32(now.x / 100), uint32(now.y / 100), uint32(now.z / 100)}
}

func getHash(aSector sector) *rand.Rand {
	id := murmur3.New64()
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, aSector.x)
	_, err := id.Write(buf)
	if err != nil {
		print("Failed to hash part 1")
	}

	binary.LittleEndian.PutUint32(buf, aSector.y)
	_, err = id.Write(buf)
	if err != nil {
		print("Failed to hash part 2")
	}

	binary.LittleEndian.PutUint32(buf, aSector.z)
	_, err = id.Write(buf)
	if err != nil {
		print("Failed to hash part 3")
	}

	return rand.New(rand.NewSource(int64(id.Sum64())))
}
