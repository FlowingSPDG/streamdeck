//go:build ignore

package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

func main() {
	dir := "imgs"
	writePNG(filepath.Join(dir, "plugin.png"), render(144))
	writePNG(filepath.Join(dir, "plugin@2x.png"), render(288))
	writePNG(filepath.Join(dir, "action.png"), render(72))
	writePNG(filepath.Join(dir, "action@2x.png"), render(144))
	writePNG(filepath.Join(dir, "key.png"), render(72))
	writePNG(filepath.Join(dir, "key@2x.png"), render(144))
}

func writePNG(path string, img image.Image) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func render(size int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	bg := color.RGBA{R: 44, G: 44, B: 44, A: 255}
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	fill(img, 0, 0, size, size, bg)

	s := float64(size) / 72
	type slider struct{ x, knobY int }
	for _, sl := range []slider{{20, 22}, {34, 30}, {48, 18}} {
		x := int(float64(sl.x) * s)
		knobY := int(float64(sl.knobY) * s)
		w := max(1, int(4*s))
		fill(img, x, int(14*s), x+w, int(50*s), white)
		fill(img, x-int(3*s), knobY, x+int(7*s), knobY+int(6*s), white)
	}
	return img
}

func fill(img *image.RGBA, x0, y0, x1, y1 int, c color.Color) {
	b := img.Bounds()
	for y := max(y0, b.Min.Y); y < min(y1, b.Max.Y); y++ {
		for x := max(x0, b.Min.X); x < min(x1, b.Max.X); x++ {
			img.Set(x, y, c)
		}
	}
}
