package fwk

import (
	gl "github.com/chsc/gogl/gl21"
)

type Drawable interface {
	Draw(ratio float64)
}

type Named interface {
	GetName() string
}

type NamedDrawable interface {
	Drawable
	Named
}

type Vector struct {
	X, Y, Z gl.Float
}

func (v *Vector) Add(p *Vector) {
	v.X, v.Y, v.Z = v.X + p.X, v.Y + p.Y, v.Z + p.Z
}

type Color struct {
	R, G, B, A gl.Float
}
