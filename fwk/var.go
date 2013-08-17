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

type Color struct {
	R, G, B, A gl.Float
}
