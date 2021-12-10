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

func hewHeatMap(xDim, yDim int, blockSize int) *heatMap {
	return &heatMap{
		image:     image.NewRGBA(image.Rect(0, 0, blockSize*xDim, blockSize*yDim)),
		blockSize: blockSize,
		baseColor: Black,
	}
}

// AddBlock shades the block at (x,y} according to the score
func (h *heatMap) AddBlock(x, y int, score float32) {
	col := floatToRGBA(score, h.baseColor)

	// TODO: figure out which way the image is orientated.
	for x1 := x * h.blockSize; x1 < (x+1)*h.blockSize; x1++ {
		for y1 := y * h.blockSize; y1 < (y+1)*h.blockSize; y1++ {
			h.image.Set(x1, y1, col)
		}
	}
}

func DrawHeatMap(data [][]float32, blockSize int, name string) {
	h := hewHeatMap(len(data), len(data[0]), blockSize)

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
