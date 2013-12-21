package main

import (
	"os"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"math"
)

type Circle struct {
  X, Y, Radius float64
	R,G,B float64
}
func (c *Circle) ToFloats() []float64 {
	out := make([]float64, 6)
	out[0] = c.X
	out[1] = c.Y
	out[2] = c.Radius
	out[3] = c.R
	out[4] = c.G
	out[5] = c.B
	return out
}
func (c *Circle) FromFloats(in []float64) int {
	c.X      = in[0]
	c.Y			 = in[1]
	c.Radius = in[2]
	c.R			 = in[3]
	c.G			 = in[4]
	c.B			 = in[5]
	return 6
}
func (c *Circle) SVG() string {
	return fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#%02x%02x%02x"/>
`, c.X, c.Y, c.Radius,
		int32(255*c.R), int32(255*c.G), int32(255*c.B))
}
func (c *Circle) String() string {
	return fmt.Sprintf(`(%.1f, %.1f) r=%.1f #%02x%02x%02x`, c.X, c.Y, c.Radius,
		int32(255*c.R), int32(255*c.G), int32(255*c.B))
}
func (c *Circle) Contains(p image.Point) bool {
	dx := float64(p.X) - c.X
	dy := float64(p.Y) - c.Y
	return dx*dx + dy*dy <= c.Radius
}
func (c *Circle) BoundingBox(bounds image.Rectangle) (r image.Rectangle) {
	r.Min.Y = 0
	if int(c.Y - c.Radius) > 0 {
		r.Min.Y = int(c.Y - c.Radius)
	}

	r.Max.Y = bounds.Max.Y
	if int(c.Y + c.Radius) < r.Max.Y {
		r.Max.Y = int(c.Y + c.Radius)
	}

	r.Min.X = 0
	if int(c.X - c.Radius) > 0 {
		r.Min.X = int(c.X - c.Radius)
	}

	r.Max.X = bounds.Max.X
	if int(c.X + c.Radius) < r.Max.X {
		r.Max.X = int(c.X + c.Radius)
	}
	return
}

func (c *Circle) FindMean(i image.Image) (rr, gg, bb, aa float64) {
	bounds := i.Bounds()
	numcounted := int64(0)

	bbox := c.BoundingBox(bounds)
	for y := bbox.Min.Y; y < bbox.Max.Y; y++ {
		for x := bbox.Min.X; x < bbox.Max.X; x++ {
			distsqr := (float64(x)-c.X)*(float64(x)-c.X) + (float64(y)-c.Y)*(float64(y)-c.Y)
			if distsqr <= c.Radius*c.Radius {
				r, g, b, a := i.At(x, y).RGBA()
				rr += float64(r)/65535
				gg += float64(g)/65535
				bb += float64(b)/65535
				aa += float64(a)/65535
				numcounted += 1
			}
		}
	}
	if numcounted > 0 {
		rr /= float64(numcounted)
		gg /= float64(numcounted)
		bb /= float64(numcounted)
		aa /= float64(numcounted)
	} else {
		fmt.Fprintf(os.Stderr, "off the map: %g, %g  vs %d, %d\n", c.X, c.Y, bounds.Max.X, bounds.Max.Y)
		fmt.Fprintf(os.Stderr, "r %g\n", c.Radius)
	}
	return
}

func (c *Circle) SetMean(i image.Image) *Circle {
	c.R, c.G, c.B, _ = c.FindMean(i)
	if c.R < 0 {
		panic("oops on red in SetMean")
	}
	return c
}

func RandomCircle(max image.Point, scale float64) *Circle {
	var c Circle
	c.X = rand.Float64()*float64(max.X)
	c.Y = rand.Float64()*float64(max.Y)
	c.Radius = (math.Abs(rand.NormFloat64())+1)*float64(max.X + max.Y)/1.5/scale
	return &c
}

func RandomCirclesApproximation(m image.Image, numcircles int) (circles Data) {
	bounds := m.Bounds()
	for i := 0; i < numcircles; i++ {
		n := RandomCircle(bounds.Max, math.Sqrt(float64(1000+i))).SetMean(m)
		circles = append(circles, n)
	}
	return
}

func (c *Circle) Render(i *image.NRGBA) {
	bbox := c.BoundingBox(i.Bounds())
	for y := bbox.Min.Y; y < bbox.Max.Y; y++ {
		for x := bbox.Min.X; x < bbox.Max.X; x++ {
			distsqr := (float64(x)-c.X)*(float64(x)-c.X) + (float64(y)-c.Y)*(float64(y)-c.Y)
			if distsqr <= c.Radius*c.Radius {
				i.SetNRGBA(x, y, color.NRGBA{uint8(255*c.R), uint8(255*c.G), uint8(255*c.B), 255})
			}
		}
	}
}
