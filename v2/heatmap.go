package heatmap

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// HeatMap holds the data for the heat map.
type HeatMap struct {
	width  int
	height int
	data   [][]float32
}

// New creates a new HeatMap.
func New(xDim, yDim int, data [][]float32) (*HeatMap, error) {
	if err := validate(xDim, yDim, data); err != nil {
		return nil, err
	}

	return &HeatMap{
		width:  xDim,
		height: yDim,
		data:   data,
	}, nil
}

func validate(xDim, yDim int, data [][]float32) error {
	if xDim < 1 || yDim < 1 {
		return errors.New("dimensions must all be greater or equal to 1")
	}

	if len(data) != xDim {
		return errors.New("data vlolates specified x dimension")
	}

	for _, i := range data {
		if len(i) != yDim {
			return errors.New("data violates specified y dimension")
		}

		for _, j := range i {
			if j < 0 || j > 1 {
				return errors.New("data must contain values that are all between 0 and 1, inclusive")
			}
		}
	}

	return nil
}

// sliceToMax takes a slice an returns a slice of identical length with all values set to 0 except
// for at the indeces of the original slice having the maximum value, which are all set to 1.
func sliceToMax(nums []float32) []float32 {
	r := make([]float32, len(nums))

	maxVal := float32(0)
	for _, num := range nums {
		if num >= maxVal {
			maxVal = num
		}
	}

	// Need to handle possibility of multiple ocurrences of max value
	maxIndex := []int{}
	for i, num := range nums {
		if num == maxVal {
			maxIndex = append(maxIndex, i)
		}
	}

	for _, i := range maxIndex {
		r[i] = 1
	}

	return r
}

// MaxY returns a new HeatMap with the maximum of each column of the original HeatMap set to 1, and 0 everywhere else.
func (h *HeatMap) MaxY() *HeatMap {
	newData := make([][]float32, h.width)

	for i, data := range h.data {
		newData[i] = sliceToMax(data)
	}

	r := &HeatMap{
		width:  h.width,
		height: h.height,
		data:   newData,
	}

	return r
}

func transpose(data [][]float32) [][]float32 {
	r := make([][]float32, len(data[0]))
	for i := range r {
		r[i] = make([]float32, len(data))
	}

	for i, d := range data {
		for j, val := range d {
			r[j][i] = val
		}
	}

	return r
}

// Transpose returns a new HeatMap, with swapped columns and rows of the original.
func (h *HeatMap) Transpose() *HeatMap {
	r := &HeatMap{
		width:  h.height,
		height: h.width,
		data:   transpose(h.data),
	}

	return r
}

// MaxX returns a new HeatMap with the maximum of each row of the original HeatMap set to 1, and 0 everywhere else.
func (h *HeatMap) MaxX() *HeatMap {
	return h.Transpose().MaxY().Transpose()
}

// Average returns an averaged HeatMap from the provided HeatMaps.
func Average(heatmaps ...*HeatMap) (*HeatMap, error) {
	if len(heatmaps) == 0 {
		return nil, errors.New("no heatmaps provided")
	}

	if len(heatmaps) == 1 {
		return heatmaps[0], nil
	}

	// Pairwise compare dimensions
	for i := 1; i < len(heatmaps); i++ {
		x := heatmaps[i]
		y := heatmaps[i-1]

		if x.height != y.height || x.width != y.width {
			return nil, errors.New("heatmap dimensions do not match")
		}
	}

	avg := &HeatMap{
		width:  heatmaps[0].width,
		height: heatmaps[0].height,
	}

	avgData := make([][]float32, avg.width)
	for i := range avgData {
		avgData[i] = make([]float32, avg.height)
	}
	for i := 0; i < avg.width; i++ {
		for j := 0; j < avg.height; j++ {
			avgData[i][j] = avgAtIndex(i, j, heatmaps)
		}
	}

	avg.data = avgData

	return avg, nil
}

func avgAtIndex(i, j int, heatmaps []*HeatMap) float32 {
	var sum float32
	for _, h := range heatmaps {
		sum += h.data[i][j]
	}
	return sum / float32(len(heatmaps))
}

// ToPNG outputs the HeatMap as a png file, using the provided scoringFunctions to translate HeatMap values into colors.
func (h *HeatMap) ToPNG(blockSize int, name string, scoringFunctions func(float32) (color.RGBA, error)) (*image.RGBA, error) {
	if blockSize < 1 {
		return nil, errors.New("heatmap blockSize must be positive")
	}

	if name == "" {
		return nil, errors.New("heatmap filename cannot be blank")
	}

	image := image.NewRGBA(image.Rect(0, 0, blockSize*h.width, blockSize*h.height))

	for i, d := range h.data {
		for j, s := range d {
			if err := h.addBlock(i, j, s, blockSize, image, scoringFunctions); err != nil {
				return nil, err
			}
		}
	}

	f, err := os.Create(name + ".png")
	if err != nil {
		return nil, err
	}

	defer f.Close()

	if err = png.Encode(f, image); err != nil {
		return nil, err
	}

	return image, nil
}

func (h *HeatMap) addBlock(x, y int, score float32, blockSize int, image *image.RGBA, scoringFunction func(float32) (color.RGBA, error)) error {

	rgba, err := scoringFunction(score)
	if err != nil {
		return err
	}

	for x1 := x * blockSize; x1 < (x+1)*blockSize; x1++ {
		// Flip image around the x axis so that the y coordinate increases upwards, not downwards.
		for y1 := (h.height - y) * blockSize; y1 > (h.height-(y+1))*blockSize; y1-- {
			image.Set(x1, y1, rgba)
		}
	}

	return nil
}

// FloatToGreyscale converts float32 inside [0, 1] to greyscale values between [#000000, #FFFFFF].
func FloatToGreyscale(score float32) (color.RGBA, error) {
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

// FloatToRedBlue converts float32 of [0, 1] to red and blue colors.
// Values below 0.5 are increasingly red the closer they get to 0, and
// values above 0.5 are increasingly blue the closer they get to 1.
func FloatToRedBlue(score float32) (color.RGBA, error) {
	if score < 0 {
		return color.RGBA{}, errors.New("score must be greater or equal to 0")
	}

	if score <= float32(0.5) {
		v := score*-510 + 255

		return color.RGBA{
			R: uint8(math.Round(float64(v))),
			G: 0,
			B: 0,
			A: 255,
		}, nil
	}

	v := score*510 - 255

	return color.RGBA{
		R: 0,
		G: 0,
		B: uint8(math.Round(float64(v))),
		A: 255,
	}, nil
}
