package main

import (
	"image"
	//"fmt"
	//"os"
)

func DiscreteObjective(a, b image.Image) (chisq float64) {
	bounds := a.Bounds()
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			ra, ga, ba, _ := a.At(x, y).RGBA()
			rb, gb, bb, _ := b.At(x, y).RGBA()
			dr := float64(int32(ra-rb))/65535
			db := float64(int32(ba-bb))/65535
			dg := float64(int32(ga-gb))/65535
			//fmt.Fprintf(os.Stderr, "dg %.3g with ga %d gb %d\n", dg, ga, gb)
			chisq += dr*dr + db*db + dg*dg
		}
	}
	return chisq/float64(bounds.Max.Y)/float64(bounds.Max.X)
}

func (d Data) Objective(i image.Image) (chisq float64) {
	r := image.NewNRGBA(i.Bounds())
	d.Render(r)
	return DiscreteObjective(r, i)
}
