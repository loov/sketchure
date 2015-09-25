package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/loov/sketchure/cleanup"
	"github.com/loov/sketchure/cleanup/filter"
)

var (
	colored = flag.Bool("colored", false, "try to preserve colors")

	white     = flag.Float64("white", 100, "the highest white value")
	lineWidth = flag.Float64("line", 0.05, "line-width relative to the image width")
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func handle(m *image.YCbCr) {
	dx := float64(m.Bounds().Dx())
	opts := &cleanup.Options{
		Whiteness: *white,
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

func YCbCrFromFile(filename string) (*image.YCbCr, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	if ycbcr, ok := src.(*image.YCbCr); ok {
		return ycbcr, nil
	}

	return ImageToYCbCr(src), nil
}

func ImageToYCbCr(src image.Image) *image.YCbCr {
	r := src.Bounds()
	dst := image.NewYCbCr(src.Bounds(), image.YCbCrSubsampleRatio444)
	for x := r.Min.X; x < r.Max.X; x++ {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			c := src.At(x, y)
			r, g, b, _ := c.RGBA()
			yy, cb, cr := color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))

			i := dst.YOffset(x, y)
			dst.Y[i] = yy
			dst.Cb[i] = cb
			dst.Cr[i] = cr
		}
	}

	return dst
}

func YCbCrClone(src *image.YCbCr) *image.YCbCr {
	dst := *src
	dst.Y = make([]uint8, len(src.Y))
	copy(dst.Y, src.Y)
	dst.Cb = make([]uint8, len(src.Cb))
	copy(dst.Cb, src.Cb)
	dst.Cr = make([]uint8, len(src.Cr))
	copy(dst.Cr, src.Cr)
	return &dst
}

func ExampleCollage(folder string) {
	images := []*image.YCbCr{}
	processed := []*image.YCbCr{}
	files, err := ioutil.ReadDir(folder)
	check(err)

	for _, file := range files {
		fmt.Println("Processing", file.Name())
		m, err := YCbCrFromFile(filepath.Join(folder, file.Name()))
		check(err)
		images = append(images, m)

		p := YCbCrClone(m)
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

	m, err := YCbCrFromFile(flag.Arg(0))
	check(err)
	handle(m)

	ImageToFile("output~.jpg", m)
}
