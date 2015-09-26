package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/loov/sketchure/cleanup"
	"github.com/loov/sketchure/cleanup/filter"
	"github.com/loov/sketchure/cleanup/ycbcr"
)

var (
	colored = flag.Bool("colored", false, "try to preserve colors")

	white     = flag.Float64("white", 1, "the highest white value")
	lineWidth = flag.Float64("line", 0.05, "line-width relative to the image width")
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func handle(m *ycbcr.Image) {
	dx := float64(m.Bounds().Dx())
	opts := &cleanup.Options{
		Whiteness: float32(*white),
		LineWidth: int(*lineWidth * dx),
	}

	cleanup.ByBase(m, opts)

	if !*colored {
		filter.Desaturate(m)
	}
}

func ImageToFile(filename string, m image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if filepath.Ext(filename) == ".png" {
		return png.Encode(file, m)
	}
	return jpeg.Encode(file, m, &jpeg.Options{Quality: 80})
}

func ExampleCollage(folder string) {
	images := []*ycbcr.Image{}
	processed := []*ycbcr.Image{}
	files, err := ioutil.ReadDir(folder)
	check(err)

	for _, file := range files {
		fmt.Println("Processing", file.Name())
		m, err := ycbcr.FromFile(filepath.Join(folder, file.Name()))
		check(err)
		images = append(images, m)

		p := m.Clone()
		handle(p)
		processed = append(processed, p)
	}

	const W = 1024
	const H = 768

	collage := image.NewRGBA(image.Rect(0, 0, W*2, H*len(images)))
	for i, in := range images {
		out := processed[i]

		r := image.Rect(0, i*H, W, (i+1)*H)
		draw.Draw(collage, r, in, in.Bounds().Min, draw.Src)

		r = image.Rect(W, i*H, W*2, (i+1)*H)
		draw.Draw(collage, r, out, out.Bounds().Min, draw.Src)
	}

	ImageToFile("output~.jpg", collage)
}

func main() {
	flag.Parse()
	if flag.Arg(0) == "" {
		ExampleCollage("examples")
		return
	}

	m, err := ycbcr.FromFile(flag.Arg(0))
	check(err)
	handle(m)

	ImageToFile("output~.jpg", m)
}
