package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

var white = color.RGBA{255, 253, 246, 255}
var black = color.RGBA{30, 30, 32, 255}

// canvas creates a light drawing surface for one benchmark sketch.
func canvas() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 512, 512))
	for y := 0; y < 512; y++ {
		for x := 0; x < 512; x++ {
			img.Set(x, y, white)
		}
	}
	return img
}

// line draws a thick freehand-like segment between two points.
func line(img *image.RGBA, x0, y0, x1, y1, width int, ink color.RGBA) {
	steps := int(math.Hypot(float64(x1-x0), float64(y1-y0))) * 2
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := int(float64(x0)*(1-t) + float64(x1)*t)
		y := int(float64(y0)*(1-t) + float64(y1)*t)
		circle(img, x, y, width/2, ink)
	}
}

// circle fills a disk used for heads, foliage, sun, and paint marks.
func circle(img *image.RGBA, cx, cy, radius int, ink color.RGBA) {
	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			if x >= 0 && y >= 0 && x < 512 && y < 512 && (x-cx)*(x-cx)+(y-cy)*(y-cy) <= radius*radius {
				img.Set(x, y, ink)
			}
		}
	}
}

// rect fills a rectangular region on the benchmark sketch.
func rect(img *image.RGBA, x0, y0, x1, y1 int, ink color.RGBA) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			img.Set(x, y, ink)
		}
	}
}

// save writes one named PNG image into the benchmark directory.
func save(dir, name string, img *image.RGBA) {
	f, err := os.Create(filepath.Join(dir, name+".png"))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

// main draws five distinct simple scenes that reveal both recognition and mood errors.
func main() {
	dir := os.Args[1]
	green := color.RGBA{43, 107, 58, 255}
	blue := color.RGBA{93, 169, 214, 255}
	yellow := color.RGBA{244, 197, 53, 255}
	red := color.RGBA{183, 48, 55, 255}
	landscape := canvas()
	circle(landscape, 405, 90, 45, yellow)
	line(landscape, 15, 350, 130, 270, 8, blue)
	line(landscape, 130, 270, 255, 350, 8, blue)
	line(landscape, 255, 350, 360, 280, 8, blue)
	line(landscape, 360, 280, 500, 365, 8, blue)
	for _, x := range []int{80, 180, 330, 440} {
		line(landscape, x, 410, x, 300, 12, black)
		circle(landscape, x, 300, 38, green)
	}
	line(landscape, 0, 410, 511, 410, 8, green)
	save(dir, "landscape", landscape)

	abstract := canvas()
	circle(abstract, 150, 140, 80, red)
	circle(abstract, 350, 165, 70, blue)
	rect(abstract, 130, 310, 395, 410, yellow)
	line(abstract, 70, 430, 410, 70, 25, black)
	save(dir, "abstract", abstract)

	person := canvas()
	circle(person, 255, 105, 45, black)
	circle(person, 255, 105, 39, white)
	line(person, 255, 150, 255, 335, 12, black)
	line(person, 255, 205, 150, 260, 12, black)
	line(person, 255, 205, 360, 260, 12, black)
	line(person, 255, 335, 165, 455, 12, black)
	line(person, 255, 335, 345, 455, 12, black)
	save(dir, "person", person)

	happy := canvas()
	circle(happy, 415, 90, 48, yellow)
	rect(happy, 100, 260, 395, 455, color.RGBA{236, 192, 124, 255})
	line(happy, 80, 265, 250, 120, 14, red)
	line(happy, 250, 120, 415, 265, 14, red)
	rect(happy, 225, 345, 280, 455, color.RGBA{105, 68, 48, 255})
	circle(happy, 100, 110, 31, red)
	circle(happy, 155, 75, 31, blue)
	line(happy, 100, 140, 115, 245, 3, black)
	line(happy, 155, 105, 115, 245, 3, black)
	save(dir, "happy_house", happy)

	threat := canvas()
	rect(threat, 0, 0, 512, 512, color.RGBA{72, 68, 76, 255})
	circle(threat, 145, 190, 35, black)
	line(threat, 145, 225, 145, 390, 16, black)
	line(threat, 145, 275, 260, 250, 15, black)
	line(threat, 145, 390, 70, 490, 14, black)
	line(threat, 145, 390, 220, 490, 14, black)
	line(threat, 260, 250, 345, 200, 7, color.RGBA{220, 220, 225, 255})
	line(threat, 345, 200, 360, 170, 7, color.RGBA{220, 220, 225, 255})
	circle(threat, 410, 335, 32, black)
	line(threat, 410, 365, 410, 460, 14, black)
	circle(threat, 345, 80, 42, red)
	save(dir, "threat", threat)
}
