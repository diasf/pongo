package fwk

import "math"

type Vector2 struct {
	X, Y float32
}

type Vector struct {
	X, Y, Z float32
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

func (v *Vector) Mult(sc float32) *Vector {
	v.X, v.Y, v.Z = v.X*sc, v.Y*sc, v.Z*sc
	return v
}

func (v *Vector) Multiplication(sc float32) *Vector {
	return &Vector{X: v.X * sc, Y: v.Y * sc, Z: v.Z * sc}
}

func (v *Vector) Dot(p *Vector) float32 {
	return float32(v.X*p.X + v.Y*p.Y + v.Z*p.Z)
}

func (v *Vector) Length() float32 {
	return float32(math.Sqrt(float64(v.Dot(v))))
}

func (v *Vector) Normalize() *Vector {
	length := float32(v.Length())
	v.X, v.Y, v.Z = v.X/length, v.Y/length, v.Z/length
	return v
}

func (v Vector) Normalization() *Vector {
	length := float32(v.Length())
	return &Vector{X: v.X / length, Y: v.Y / length, Z: v.Z / length}
}

func (v Vector) Clamp(min, max Vector) (p Vector) {
	p.X = clamp(v.X, min.X, max.X)
	p.Y = clamp(v.Y, min.Y, max.Y)
	p.Z = clamp(v.Z, min.Z, max.Z)
	return p
}

func clamp(a, min, max float32) float32 {
	if a < min {
		return min
	} else if a > max {
		return max
	}
	return a
}
