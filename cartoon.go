package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"time"
	"math/rand"

	// Package image/jpeg is not used explicitly in the code below,
	// but is imported for its initialization side-effect, which allows
	// image.Decode to understand JPEG formatted images. Uncomment these
	// two lines to also understand GIF and PNG images:
	_ "image/gif"
	"image/png"
	_ "image/jpeg"
)

func main() {
  rand.Seed(time.Now().UTC().UnixNano())
	reader, err := os.Open(os.Args[1])
	if err != nil {
	    log.Fatal(err)
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	circles := RandomCirclesApproximation(m, 100000)

	bounds := m.Bounds()
	fmt.Printf(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<svg
  xmlns="http://www.w3.org/2000/svg"
  xmlns:xlink="http://www.w3.org/1999/xlink"
  width="%d"
  height="%d">
`, 3*bounds.Max.X, bounds.Max.Y)

	fmt.Fprintf(os.Stderr, "rendered chisq is %.3g\n", DiscreteObjective(RenderCircles(circles, bounds), m))
	//fmt.Fprintf(os.Stderr, "chisq is %.3g\n", objective(EvaluateCircles(circles), m))
	for _, c := range(circles) {
		fmt.Print(c.SVG())
	}

	file, err := os.Create("foo-rendered.png")
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(file, RenderCircles(circles, bounds))

	fmt.Printf(`
  <image
     y="%d"
     x="%d"
     id="stupid123"
     xlink:href="file:///srv/home/droundy/src/cartoon/foo-rendered.png"
     height="%d"
     width="%d" />
`, 0, 2*bounds.Max.X, bounds.Max.Y, bounds.Max.X)

	fmt.Printf(`
  <image
     y="%d"
     x="%d"
     id="stupid124"
     xlink:href="file:///srv/home/droundy/src/cartoon/%s"
     height="%d"
     width="%d" />
`, 0, bounds.Max.X, os.Args[1], bounds.Max.Y, bounds.Max.X)
	fmt.Println(`</svg>
`)
}
