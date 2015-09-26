// +build ignore

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"

	"github.com/loov/sketchure/cleanup"
	"github.com/loov/sketchure/ycbcr"
)

var (
	colored   = flag.Bool("c", false, "try to preserve colors")
	white     = flag.Float64("w", 100, "the highest white value")
	lineWidth = flag.Float64("l", 0.05, "line-width relative to the image width")
	output    = flag.String("o", "collage.jpg", "result image file")

	width  = flag.Int("width", 640, "single image width")
	height = flag.Int("height", 480, "single image height")
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func process(m *image.YCbCr) {
	dx := float64(m.Bounds().Dx())
	opts := &cleanup.Options{
		Whiteness:  *white,
		LineWidth:  int(*lineWidth * dx),
		Desaturate: !*colored,
	}

	cleanup.ByBase(m, opts)
}

func fitted(into image.Rectangle, size image.Point) image.Rectangle {
	ratio := float64(size.X) / float64(size.Y)

	if size.X > into.Dx() {
		size.X = into.Dx()
		size.Y = int(float64(into.Dx()) / ratio)
	}

	if size.Y > into.Dy() {
		size.X = int(float64(into.Dy()) * ratio)
		size.Y = into.Dy()
	}

	r := image.Rectangle{
		Min: into.Min,
		Max: into.Min.Add(size),
	}
	r = r.Add(image.Point{into.Dx() / 2, into.Dy() / 2}).
		Sub(image.Point{size.X / 2, size.Y / 2})

	return r
}

func drawcell(collage draw.Image, src image.Image, x, y, w, h int) {
	sz := src.Bounds().Size()
	bounds := image.Rect(
		x*w, y*h,
		(x+1)*w, (y+1)*h,
	)

	dr := fitted(bounds, sz)
	draw.CatmullRom.Scale(collage, dr, src, src.Bounds(), draw.Over, nil)
}

func main() {
	flag.Parse()

	dir := flag.Arg(0)
	if dir == "" {
		fmt.Fprintf(os.Stderr, "Images not specified.")
		flag.Usage()
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(dir)
	check(err)

	w, h := *width, *height
	collage := image.NewRGBA(image.Rect(0, 0, w*2, h*len(files)))
	draw.Draw(collage, collage.Bounds(), image.NewUniform(color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}), image.Point{}, draw.Src)

	for i, file := range files {
		fmt.Println("Processing ", file.Name())
		m, err := ycbcr.FromFile(filepath.Join(dir, file.Name()))
		check(err)

		drawcell(collage, m, 0, i, w, h)
		process(m)
		drawcell(collage, m, 1, i, w, h)
	}

	check(ycbcr.ToFile(*output, collage))
}
