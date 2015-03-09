package fwk

import (
	"math"

	gl "github.com/chsc/gogl/gl21"
)

type Vector struct {
	X, Y, Z gl.Float
}

func (v *Vector) Add(p *Vector) *Vector {
	v.X, v.Y, v.Z = v.X+p.X, v.Y+p.Y, v.Z+p.Z
	return v
}

func (v Vector) Addition(p *Vector) *Vector {
	return &Vector{X: v.X + p.X, Y: v.Y + p.Y, Z: v.Z + p.Z}
}

func (v *Vector) Substract(p *Vector) *Vector {
	v.X, v.Y, v.Z = v.X-p.X, v.Y-p.Y, v.Z-p.Z
	return v
}

func (v Vector) Substraction(p *Vector) *Vector {
	return &Vector{X: v.X - p.X, Y: v.Y - p.Y, Z: v.Z - p.Z}
}

func (v *Vector) Mult(sc gl.Float) *Vector {
	v.X, v.Y, v.Z = v.X*sc, v.Y*sc, v.Z*sc
	return v
}

func (v *Vector) Multiplication(sc gl.Float) *Vector {
	return &Vector{X: v.X * sc, Y: v.Y * sc, Z: v.Z * sc}
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

func (v Vector) Normalization() *Vector {
	length := gl.Float(v.Length())
	return &Vector{X: v.X / length, Y: v.Y / length, Z: v.Z / length}
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
