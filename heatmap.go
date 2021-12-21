package heatmap

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type heatMap struct {
	image     *image.RGBA
	blockSize int
	width     int
	height    int
}

func newHeatMap(xDim, yDim int, blockSize int) *heatMap {
	return &heatMap{
		image:     image.NewRGBA(image.Rect(0, 0, blockSize*xDim, blockSize*yDim)),
		blockSize: blockSize,
		width:     xDim,
		height:    yDim,
	}
}

func (h *heatMap) addBlock(x, y int, score float32) error {
	if x < 0 || x > h.width-1 || y < 0 || y > h.height-1 {
		errorString := fmt.Sprintf("heatmap coordinate (%v, %v) is out of bounds", x, y)
		return errors.New(errorString)
	}

	if score < 0 {
		return errors.New("score must be greater or equal to 0")
	}

	col, err := floatToGreyscale(score)
	if err != nil {
		return err
	}

	for x1 := x * h.blockSize; x1 < (x+1)*h.blockSize; x1++ {
		// Flip image around the x axis so that the y coordinate increases upwards, not downwards.
		for y1 := (h.height - y) * h.blockSize; y1 > (h.height-(y+1))*h.blockSize; y1-- {
			h.image.Set(x1, y1, col)
		}
	}

	return nil
}

// DrawHeatMap produces a heatmap image PNG with filenam "name.png"
func DrawHeatMap(data [][]float32, blockSize int, name string) error {
	if blockSize < 1 {
		return errors.New("heatmap blocksize must be positive")
	}

	if name == "" {
		return errors.New("heatmap filename cannot be blank")
	}

	if len(data) == 0 || len(data[0]) == 0 {
		return errors.New("heatmap is missing data in one or more dimensions")
	}

	h := newHeatMap(len(data), len(data[0]), blockSize)

	for i, d := range data {
		for j, s := range d {
			if err := h.addBlock(i, j, s); err != nil {
				return err
			}
		}
	}

	f, err := os.Create(name + ".png")
	if err != nil {
		return err
	}

	defer f.Close()

	if err = png.Encode(f, h.image); err != nil {
		return err
	}

	return nil
}

// Convert float [0, 1] to greyscale, [#000000, #FFFFFF]
func floatToGreyscale(score float32) (color.RGBA, error) {
	if score < 0 {
		return color.RGBA{}, errors.New("score must be greater or equal to 0")
	}

	v := score * 255

	return color.RGBA{
		R: uint8(math.Round(float64(v))),
		G: uint8(math.Round(float64(v))),
		B: uint8(math.Round(float64(v))),
		A: 255,
	}, nil
}
