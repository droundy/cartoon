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

	bounds := m.Bounds()
	fmt.Printf(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<svg
  xmlns="http://www.w3.org/2000/svg"
  xmlns:xlink="http://www.w3.org/1999/xlink"
  width="%d"
  height="%d">
`, 3*bounds.Max.X, bounds.Max.Y)

	circles := Data{RandomCircle(m.Bounds().Max, 1).SetMean(m)}
	simplex := CreateSimplex(circles, m)
	for nn := 0; nn < 4; nn++ {
		//circles = append(simplex.Data, RandomCirclesApproximation(m, 1)...)
		circles = append(simplex.Data, RandomCircle(m.Bounds().Max, 20).SetMean(m))
		fmt.Fprintf(os.Stderr, "I now have %d circles\n", len(circles))
		simplex = CreateSimplex(circles, m)
		fmt.Fprintf(os.Stderr, "first guess chisq is %.4v\n", simplex.Badnesses[0])
		for nn:=0; nn<10; nn++ {
			for ii := 0; ii < len(simplex.X[0]); ii++ {
				simplex.Improve()
			}
			fmt.Fprintf(os.Stderr, "chisq is %.4v latest circle:\n%v",
				simplex.Badnesses[simplex.best], simplex.Data)
		}
	}

	rendered := image.NewNRGBA(m.Bounds())
	circles.Render(rendered)
	fmt.Fprintf(os.Stderr, "chisq is %.3g\n", DiscreteObjective(rendered, m))
	fmt.Print(simplex.SVG())

	file, err := os.Create("foo-rendered.png")
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(file, rendered)

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
