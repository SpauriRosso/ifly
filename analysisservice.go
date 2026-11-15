package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	_ "image/png"
	"math"
	"strings"
)

// StructuralAnalysis contains measured visual features; it makes no semantic or emotional claims.
type StructuralAnalysis struct {
	Brightness     float64 `json:"brightness"`
	Contrast       float64 `json:"contrast"`
	EdgeDensity    float64 `json:"edgeDensity"`
	ColorDiversity float64 `json:"colorDiversity"`
	Complexity     float64 `json:"complexity"`
	Width          int     `json:"width"`
	Height         int     `json:"height"`
}

type AnalysisService struct{}

// Analyze decodes a small PNG data URL and returns deterministic visual measurements or an error.
func (s *AnalysisService) Analyze(dataURL string) (StructuralAnalysis, error) {
	if len(dataURL) > 2_000_000 || !strings.HasPrefix(dataURL, "data:image/png;base64,") {
		return StructuralAnalysis{}, errors.New("invalid analysis image")
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(dataURL, "data:image/png;base64,"))
	if err != nil {
		return StructuralAnalysis{}, err
	}
	config, format, err := image.DecodeConfig(bytes.NewReader(raw))
	if err != nil || format != "png" || config.Width < 1 || config.Height < 1 || config.Width > 256 || config.Height > 256 {
		return StructuralAnalysis{}, errors.New("invalid analysis dimensions")
	}
	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return StructuralAnalysis{}, err
	}
	return measureImage(img), nil
}

// measureImage samples pixels, adjacent luminance changes, and quantized colors to describe image structure.
func measureImage(img image.Image) StructuralAnalysis {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	luminance := make([]float64, w*h)
	colors := make(map[uint16]struct{})
	var sum, squared, edges float64
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			alpha := float64(a) / 65535
			red := alpha*float64(r)/65535 + (1 - alpha)
			green := alpha*float64(g)/65535 + (1 - alpha)
			blue := alpha*float64(b)/65535 + (1 - alpha)
			value := .2126*red + .7152*green + .0722*blue
			index := y*w + x
			luminance[index] = value
			sum += value
			squared += value * value
			colors[uint16(red*7)<<6|uint16(green*7)<<3|uint16(blue*7)] = struct{}{}
			if x > 0 && math.Abs(value-luminance[index-1]) > .12 {
				edges++
			}
			if y > 0 && math.Abs(value-luminance[index-w]) > .12 {
				edges++
			}
		}
	}
	count := float64(w * h)
	mean := sum / count
	contrast := math.Sqrt(math.Max(0, squared/count-mean*mean)) * 2
	edgeDensity := edges / math.Max(1, float64((w-1)*h+w*(h-1)))
	diversity := math.Min(1, float64(len(colors))/128)
	complexity := math.Min(1, .4*edgeDensity*3+.3*contrast+.3*diversity)
	return StructuralAnalysis{mean, math.Min(1, contrast), edgeDensity, diversity, complexity, w, h}
}
