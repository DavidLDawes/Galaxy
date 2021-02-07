package main

import (
	"encoding/binary"
	"image/color"
	"math/rand"

	"github.com/spaolacci/murmur3"
)

var (
	bright    = uint8(255)
	mixbright = uint8(bright * 165 / 255)
	tween     = uint8(228)
	med       = uint8(196)
	mixmed    = uint8(med * 165 / 255)
	dim       = uint8(128)
	mixdim    = uint8(dim * 165 / 255)
	submix    = uint8(99)
)

var timesBigger = 10.0

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
	x           float32
	y           float32
	z           float32
	sx          float32
	sy          float32
	sz          float32
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

var classO classdetails = classdetails{
	class:       "O",
	brightcolor: color.RGBA{dim, dim, bright, 255},
	medcolor:    color.RGBA{dim / 2, dim / 2, med, 255},
	dimcolor:    color.RGBA{dim / 4, dim / 4, dim, 255},
	odds:        .0006,
	fudge:       .000402,
	minmass:     16.00001,
	deltamass:   243.2,
	minradii:    6,
	deltaradii:  17.3,
	minlum:      30000,
	deltalum:    147000.2,
	pixels:      11,
}

var classB classdetails = classdetails{
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

var classA classdetails = classdetails{
	class:       "A",
	brightcolor: color.RGBA{bright, bright, bright, 255},
	medcolor:    color.RGBA{med, med, med, 255},
	dimcolor:    color.RGBA{dim, dim, dim, 255},
	odds:        .0006,
	fudge:       .0018,
	minmass:     1.4,
	deltamass:   .7,
	minradii:    1.4,
	deltaradii:  .4,
	minlum:      5,
	deltalum:    20,
	pixels:      6,
}

var classF classdetails = classdetails{
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

var classG classdetails = classdetails{
	class:       "G",
	brightcolor: color.RGBA{tween, tween, 0, 255},
	medcolor:    color.RGBA{med, med, 0, 255},
	dimcolor:    color.RGBA{dim, dim, 0, 255},
	odds:        .076,
	fudge:       .01502,
	minmass:     .8,
	deltamass:   .24,
	minradii:    .96,
	deltaradii:  .19,
	minlum:      .6,
	deltalum:    .9,
	pixels:      4,
}

var classK classdetails = classdetails{
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

var classM classdetails = classdetails{
	class:       "M",
	brightcolor: color.RGBA{bright, 0, 0, 255},
	medcolor:    color.RGBA{med, 0, 0, 255},
	dimcolor:    color.RGBA{dim, 0, 0, 255},
	odds:        .7645,
	fudge:       .14,
	minmass:     1.04,
	deltamass:   .36,
	minradii:    1.15,
	deltaradii:  .25,
	minlum:      1.5,
	deltalum:    3.5,
	pixels:      2,
}

var starDetailsByClass = [7]classdetails{classO, classB, classA, classF, classG, classK, classM}

var classByZoom = [8]int{7, 7, 7, 6, 5, 3, 2, 1}

var hashCharByIteration = [7]int{0, 2, 8, 16, 24, 12, 23}

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
)

func getStarDetails(classDetails classdetails, sector sector, random *rand.Rand) []star {
	stars := make([]star, 0)
	loopsize := int32(1000 * (classDetails.odds - classDetails.fudge + 2*classDetails.fudge*random.Float32()))
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

func getSectorDetails(fromSector sector) []star {
	result := make([]star, 0)
	random := getHash(fromSector)
	for _, starDetails := range starDetailsByClass {
		nextClass := getStarDetails(starDetails, fromSector, random)
		for _, classStar := range nextClass {
			result = append(result, classStar)
		}
	}
	return result
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
