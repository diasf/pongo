package fwk

import (
	gl "github.com/chsc/gogl/gl21"
	"math"
)

type Vector struct {
	X, Y, Z gl.Float
}

func (v *Vector) Add(p *Vector) *Vector {
	v.X, v.Y, v.Z = v.X+p.X, v.Y+p.Y, v.Z+p.Z
	return v
}

func (v Vector) Addition(p *Vector) (addition *Vector) {
	addition.X, addition.Y, addition.Z = v.X+p.X, v.Y+p.Y, v.Z+p.Z
	return
}

func (v *Vector) Substract(p *Vector) *Vector {
	v.X, v.Y, v.Z = v.X-p.X, v.Y-p.Y, v.Z-p.Z
	return v
}

func (v Vector) Substraction(p *Vector) (substraction *Vector) {
	substraction.X, substraction.Y, substraction.Z = v.X-p.X, v.Y-p.Y, v.Z-p.Z
	return
}

func (v *Vector) Mult(sc gl.Float) *Vector {
	v.X, v.Y, v.Z = v.X*sc, v.Y*sc, v.Z*sc
	return v
}

func (v *Vector) Multiplication(sc gl.Float) (multiplication *Vector) {
	multiplication.X, multiplication.Y, multiplication.Z = v.X*sc, v.Y*sc, v.Z*sc
	return multiplication
}

func (v *Vector) Dot(p *Vector) float64 {
	return float64(v.X*p.X + v.Y*p.Y + v.Z*p.Z)
}

func (v *Vector) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v *Vector) Normalize() *Vector {
	length := gl.Float(v.Length())
	v.X, v.Y, v.Z = v.X/length, v.Y/length, v.Z/length
	return v
}

func (v Vector) Normalization() (normalized *Vector) {
	length := gl.Float(v.Length())
	normalized.X, normalized.Y, normalized.Z = v.X/length, v.Y/length, v.Z/length
	return
}

func (v Vector) Clamp(min, max Vector) (p Vector) {
	p.X = clamp(v.X, min.X, max.X)
	p.Y = clamp(v.Y, min.Y, max.Y)
	p.Z = clamp(v.Z, min.Z, max.Z)
	return p
}

func clamp(a, min, max gl.Float) gl.Float {
	if a < min {
		return min
	} else if a > max {
		return max
	}
	return a
}
