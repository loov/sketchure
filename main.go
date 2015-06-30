package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/loov/sketch-capture/filter"
	"github.com/loov/sketch-capture/lab"

	_ "image/jpeg"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func ImageToFile(filename string, m image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, m)
}

func ExampleCollage(folder string) {
	images := []*lab.Image{}
	processed := []*lab.Image{}
	files, err := ioutil.ReadDir(folder)
	check(err)

	for _, file := range files {
		fmt.Println("Processing", file.Name())
		m, err := lab.ImageFromFile(filepath.Join(folder, file.Name()))
		check(err)
		images = append(images, m)

		p := m.Clone()
		filter.Normalize(p, 110, 15)
		// filter.Desaturate(p)
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

	ImageToFile("output~.png", collage)
}

func main() {
	ExampleCollage("examples")
}
