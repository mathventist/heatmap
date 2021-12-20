package heatmap

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type heatMap struct {
	image     *image.RGBA
	blockSize int
	baseColor color.RGBA
	width     int
	height    int
}

var (
	Black = color.RGBA{0, 0, 0, 255}
	//Red     = color.RGBA{255, 0, 0, 255}
	//Green   = color.RGBA{0, 255, 0, 255}
	//Blue    = color.RGBA{0, 0, 255, 255}
	//Yellow  = color.RGBA{255, 255, 0, 255}
	//Cyan    = color.RGBA{0, 255, 255, 255}
	//Magenta = color.RGBA{255, 0, 255, 255}
)

func newHeatMap(xDim, yDim int, blockSize int) *heatMap {
	return &heatMap{
		image:     image.NewRGBA(image.Rect(0, 0, blockSize*xDim, blockSize*yDim)),
		blockSize: blockSize,
		baseColor: Black,
		width:     xDim,
		height:    yDim,
	}
}

// AddBlock shades the block at (x, y), colored according to the score
func (h *heatMap) AddBlock(x, y int, score float32) {
	col := floatToRGBA(score, h.baseColor)

	for x1 := x * h.blockSize; x1 < (x+1)*h.blockSize; x1++ {
		// Flip image around the x axis so that the y coordinate increases upwards, not downwards.
		for y1 := (h.height - y) * h.blockSize; y1 > (h.height-(y+1))*h.blockSize; y1-- {
			h.image.Set(x1, y1, col)
		}
	}
}

// DrawHeatMap produces a heatmap image PNG with filenam "name.png"
func DrawHeatMap(data [][]float32, blockSize int, name string) {
	h := newHeatMap(len(data), len(data[0]), blockSize)

	for i, d := range data {
		for j, s := range d {
			h.AddBlock(i, j, s)
		}
	}

	f, err := os.Create(name + ".png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, h.image)
}

// TODO: implement this properly. For now it only operates from black to white.
// Combine a float between 0 and 1, and a base rgba, to create a new rgba, where the
// float is treated as the lightness component of the HSV color model.
func floatToRGBA(f float32, base color.RGBA) color.RGBA {
	v := f * 255

	return color.RGBA{
		R: uint8(math.Round(float64(v))),
		G: uint8(math.Round(float64(v))),
		B: uint8(math.Round(float64(v))),
		A: 255,
	}
	// Step 1: determine the hue from the base color
	//hue := rgbaToHue(base)

	//// Step 2: determine the saturation
	//sat := rgbaToSaturation(base)

	//// Step 3: set the lightness
	//light := f

	//// Step 4: recombine to a new rgba model
	//return HSVToRBGA(hue, sat, light)
}

//func rgbaToHue(c color.RGBA) float32 {
//}
//
//func rgbaToSaturation(c color.RGBA) float32 {
//}
//
//func HSVToRBGA(h, s, l float32) color.RGBA {
//}
