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

type Color struct {
	R, G, B, A gl.Float
}
