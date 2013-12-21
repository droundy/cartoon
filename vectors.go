package main

import (
	"image"
)

type Datum interface {
	ToFloats() []float64
	FromFloats([]float64) int
	SVG() string
	String() string
	Render(*image.NRGBA)
}

type Data []Datum

func (ds Data) ToFloats() (out []float64) {
	for _,d := range(ds) {
		out = append(out, d.ToFloats()...)
	}
	return
}
func (ds Data) FromFloats(in []float64) (num int) {
	for i := range(ds) {
		n := ds[i].FromFloats(in)
		in = in[n:]
		num += n
	}
	return
}
func (ds Data) Render(im *image.NRGBA) {
	for i := range(ds) {
		ds[i].Render(im)
	}
}
func (ds Data) SVG() (out string) {
	for _,d := range(ds) {
		out += d.SVG()
		out += "\n"
	}
	return
}
func (ds Data) String() (out string) {
	out = "\t"
	for _,d := range(ds) {
		out += d.String()
		out += "\n\t"
	}
	return out[:len(out)-1]
}

type Simplex struct {
	Data
	I image.Image
	X [][]float64
	Badnesses []float64

	Xo []float64
	bXo float64

	best int
}

func CreateSimplex(d Data, im image.Image) (o Simplex) {
	o.Data = d
	o.I = im
	v0 := d.ToFloats()
	nvert := len(v0)+1
	o.X = make([][]float64, nvert)
	o.X[0] = v0
	for i:=1; i<nvert; i++ {
		o.X[i] = append(o.X[i], v0...)
		if o.X[i][i-1] < 1.0 && o.X[i][i-1] > 0 {
			o.X[i][i-1] *= 0.9
		} else {
			o.X[i][i-1] *= 0.8
		}
	}
	o.Badnesses = make([]float64, nvert)
	for i := range(o.X) {
		o.Badnesses[i] = o.Badness(o.X[i])
	}
	return
}

func (o *Simplex) Badness(x []float64) float64 {
	o.FromFloats(x)
	return o.Objective(o.I)
}

func (o *Simplex) Improve() {
	worst := 0
	secondworst := 0
	for i := range(o.Badnesses) {
		if o.Badnesses[i] >= o.Badnesses[worst] {
			secondworst = worst
			worst = i
		} else if o.Badnesses[i] < o.Badnesses[o.best] {
			o.best = i
		} else if o.Badnesses[i] > o.Badnesses[secondworst] && o.Badnesses[i] < o.Badnesses[worst] {
			secondworst = i
		}
	}

	// Wikipedia step 2
	o.Xo = make([]float64, len(o.X[0]))
	for iv,v := range(o.X) {
		if iv != worst {
			for i := range(o.Xo) {
				o.Xo[i] += v[i]
			}
		}
	}
	for i := range(o.Xo) {
		o.Xo[i] /= float64(len(o.X)-1)
	}
	o.bXo = o.Badness(o.Xo)

	// Wikipedia step 3, reflection
	Xr := make([]float64, len(o.X[0]))
	for i := range(Xr) {
		Xr[i] = 2*o.Xo[i] - o.X[worst][i]
	}

	bXr := o.Badness(Xr)
	if bXr < o.Badnesses[secondworst] && bXr > o.Badnesses[o.best] {
		o.Badnesses[worst] = bXr
		o.X[worst] = Xr
		return
	}

	// Wikipedia step 4, expansion
	if bXr < o.Badnesses[o.best] {
		Xe := make([]float64, len(o.X[0]))
		for i := range(Xe) {
			Xe[i] = 3*o.Xo[i] - 2*o.X[worst][i]
		}
		bXe := o.Badness(Xe)

		if bXe < bXr {
			o.Badnesses[worst] = bXe
			o.X[worst] = Xe
		} else {
			o.Badnesses[worst] = bXr
			o.X[worst] = Xr
		}
		return
	}

	// Wikipedia step 5, contraction
	Xc := make([]float64, len(o.X[0]))
	for i := range(Xc) {
		Xc[i] = 0.5*(o.Xo[i] + o.X[worst][i])
	}
	bXc := o.Badness(Xc)

	if bXc < o.Badnesses[worst] {
		o.Badnesses[worst] = bXc
		o.X[worst] = Xc
	}

	// Wikipedia step 6, reduction
	for n := range(o.X) {
		if n != o.best {
			for i := range(o.X[n]) {
				o.X[n][i] = 0.5*(o.X[n][i] + o.X[o.best][i])
			}
		}
	}
	o.FromFloats(o.X[o.best])
}
