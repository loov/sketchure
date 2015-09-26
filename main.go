package main

import (
	"flag"
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/loov/sketchure/cleanup"
	"github.com/loov/sketchure/ycbcr"
)

var (
	colored   = flag.Bool("c", false, "try to preserve colors")
	white     = flag.Float64("w", 100, "the highest white value")
	lineWidth = flag.Float64("l", 0.05, "line-width relative to the image width")

	outdir = flag.String("out", ".", "output directory")
	suffix = flag.String("suffix", "~", "suffix to use for processed images")
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

func main() {
	flag.Parse()
	files := flag.Args()

	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "Images not specified.")
		flag.Usage()
		os.Exit(1)
	}

	for _, filename := range files {
		m, err := ycbcr.FromFile(filename)
		check(err)
		process(m)

		name := filepath.Join(*outdir, filepath.Base(filename))
		ext := filepath.Ext(name)

		name = name[:len(name)-len(ext)] + *suffix + ext
		ycbcr.ToFile(name, m)
	}
}
