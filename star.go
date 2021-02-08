package main

import (
	"encoding/binary"
	"image/color"
	"math"
	"math/rand"

	"github.com/spaolacci/murmur3"
)

const (
	xmin = 0
	xmax = 100000
	ymin = 0
	ymax = 100000
	zmin = 0
	zmax = 25000
	xadj = 100
	yadj = 100
	zadj = 100
	starsPerMegaCubicLY = float32(4023.728813559)
)

type classdetails struct {
	class       string
	brightcolor color.RGBA
	medcolor    color.RGBA
	dimcolor    color.RGBA
	odds        float32
	fudge       float32
	minmass     float32
	deltamass   float32
	minradii    float32
	deltaradii  float32
	minlum      float32
	deltalum    float32
	pixels      int32
}

type star struct {
	class       string
	brightcolor color.RGBA
	medcolor    color.RGBA
	dimcolor    color.RGBA
	pixels      int32
	mass        float32
	radii       float32
	lumens      float32
	// 3D postion
	x           float32
	y           float32
	z           float32
	// sector location, 0 <= sx, sy, sz < 100
	sx          float32
	sy          float32
	sz          float32
	// display location, 0 <= dx, dy, dz < 1000
	dx          float32
	dy          float32
	dz          float32
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
	bright    = uint8(255)
	mixbright = uint8(bright * 165 / 255)
	tween     = uint8(228)
	med       = uint8(196)
	mixmed    = uint8(med * 165 / 255)
	dim       = uint8(128)
	mixdim    = uint8(dim * 165 / 255)
	submix    = uint8(99)

	classO classdetails = classdetails{
		class:       "O",
		brightcolor: color.RGBA{dim, dim, bright, 255},
		medcolor:    color.RGBA{dim / 2, dim / 2, med, 255},
		dimcolor:    color.RGBA{dim / 4, dim / 4, dim, 255},
		odds:        .0000003,
		fudge:       .0000000402,
		minmass:     16.00001,
		deltamass:   243.2,
		minradii:    6,
		deltaradii:  17.3,
		minlum:      30000,
		deltalum:    147000.2,
		pixels:      11,
	}


	classB classdetails = classdetails{
		class:       "B",
		brightcolor: color.RGBA{dim / 2, dim / 2, bright, 255},
		medcolor:    color.RGBA{dim / 4, dim / 4, med, 255},
		dimcolor:    color.RGBA{dim / 8, dim / 8, dim, 255},
		odds:        .0013,
		fudge:       .0003,
		minmass:     2.1,
		deltamass:   13.9,
		minradii:    1.8,
		deltaradii:  4.8,
		minlum:      25,
		deltalum:    29975,
		pixels:      8,
	}

	classA classdetails = classdetails{
		class:       "A",
		brightcolor: color.RGBA{bright, bright, bright, 255},
		medcolor:    color.RGBA{med, med, med, 255},
		dimcolor:    color.RGBA{dim, dim, dim, 255},
		odds:        .006,
		fudge:       .0018,
		minmass:     1.4,
		deltamass:   .7,
		minradii:    1.4,
		deltaradii:  .4,
		minlum:      5,
		deltalum:    20,
		pixels:      6,
	}

	classF classdetails = classdetails{
		class:       "F",
		brightcolor: color.RGBA{bright, bright, 0, 255},
		medcolor:    color.RGBA{tween, tween, 0, 255},
		dimcolor:    color.RGBA{med, med, 0, 255},
		odds:        .03,
		fudge:       .012,
		minmass:     1.04,
		deltamass:   .36,
		minradii:    1.15,
		deltaradii:  .25,
		minlum:      1.5,
		deltalum:    3.5,
		pixels:      5,
	}

	classG classdetails = classdetails{
		class:       "G",
		brightcolor: color.RGBA{tween, tween, 0, 255},
		medcolor:    color.RGBA{med, med, 0, 255},
		dimcolor:    color.RGBA{dim, dim, 0, 255},
		odds:        .076,
		fudge:       .01102,
		minmass:     .8,
		deltamass:   .24,
		minradii:    .96,
		deltaradii:  .19,
		minlum:      .6,
		deltalum:    .9,
		pixels:      4,
	}

	classK classdetails = classdetails{
		class:       "K",
		brightcolor: color.RGBA{bright, submix, 0, 255},
		medcolor:    color.RGBA{med, uint8(submix * med / 255), 0, 255},
		dimcolor:    color.RGBA{dim, uint8(submix * dim / 255), 0, 255},
		odds:        .121,
		fudge:       .042,
		minmass:     .45,
		deltamass:   .35,
		minradii:    .7,
		deltaradii:  .26,
		minlum:      .08,
		deltalum:    .52,
		pixels:      3,
	}

	classM classdetails = classdetails{
		class:       "M",
		brightcolor: color.RGBA{bright, 0, 0, 255},
		medcolor:    color.RGBA{med, 0, 0, 255},
		dimcolor:    color.RGBA{dim, 0, 0, 255},
		odds:        .7645,
		fudge:       .04,
		minmass:     1.04,
		deltamass:   .36,
		minradii:    1.15,
		deltaradii:  .25,
		minlum:      1.5,
		deltalum:    3.5,
		pixels:      2,
	}

	starDetailsByClass = [7]classdetails{classO, classB, classA, classF, classG, classK, classM}
	classByZoom = [8]int{7, 7, 7, 7, 7, 3, 2, 2}
)

func getStarDetails(classDetails classdetails, sector sector, random *rand.Rand) []star {
	stars := make([]star, 0)
	loopsize := int32(423.728813559 * (classDetails.odds - classDetails.fudge + 2*classDetails.fudge*random.Float32()))
	for i := 0; i < int(loopsize); i++ {
		nextstar := star{}
		rando := random.Float32()
		nextstar.sx = random.Float32() * xadj
		nextstar.sy = random.Float32() * yadj
		nextstar.sz = random.Float32() * zadj
		nextstar.x = float32(sector.x)*xadj + nextstar.sx
		nextstar.y = float32(sector.y)*yadj + nextstar.sy
		nextstar.z = float32(sector.z)*zadj + nextstar.sz
		nextstar.class = classDetails.class
		nextstar.brightcolor = classDetails.brightcolor
		nextstar.dimcolor = classDetails.dimcolor
		nextstar.mass = classDetails.minmass + classDetails.deltamass*(1+rando)
		nextstar.radii = classDetails.minradii + rando*classDetails.deltaradii
		nextstar.lumens = classDetails.minlum + rando*classDetails.deltalum
		nextstar.pixels = classDetails.pixels
		stars = append(stars, nextstar)
	}
	return stars
}

func getSectorDetails(fromSector sector, xtrim float32) []star {
	result := make([]star, 0)
	random := getHash(fromSector)
	classCount := 0
	for _, starDetails := range starDetailsByClass {
		nextClass := getStarDetails(starDetails, fromSector, random)
		for _, classStar := range nextClass {
			result = append(result, classStar)
		}
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

var lastDetails saveDetails = saveDetails{
	position{50000, 50000, 12500},
	100, make([]star, 0), false,
}

func getGalaxyDetails() ([]star) {
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

	resultStars := make([]star , 0)

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
	for xi := uint32(0); 100*xi < uint32(xMax - xMin) + extraX; xi++ {
		xtrim := float32(100*xi)
		for yi := uint32(0); 100*yi < uint32(yMax - yMin) + extraY; yi++ {
			for zi := uint32(0); 100*zi < uint32(zMax - zMin); zi++ {
				for _, star := range getSectorDetails(getSectorFromPosition(
					position{here.x + 100.0*float32(xi),
						here.y + 100.0*float32(yi),
						here.z + 100.0*float32(zi)}),
						xtrim) {
					if !(star.x < xMin) && !(star.x > xMax) &&
						!(star.y < yMin) && !(star.y > yMax) &&
						!(star.z < zMin) && !(star.z > zMax) {
						star.dx = (star.x - xMin)*1000/(xMax - xMin)
						star.dy = (star.y - yMin)*1000/(yMax - yMin)
						star.dz = (star.z - zMin)*1000/(zMax - zMin)
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
	return sector{uint32(now.x/100), uint32(now.y/100), uint32(now.z/100)}
}

func getHash(aSector sector) *rand.Rand {
	id := murmur3.New64()
	var buf []byte = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf[:], aSector.x)
	id.Write(buf)
	binary.LittleEndian.PutUint32(buf[:], aSector.y)
	id.Write(buf)
	binary.LittleEndian.PutUint32(buf[:], aSector.z)
	id.Write(buf)

	return rand.New(rand.NewSource(int64(id.Sum64()))) //nolint:gosec
}
