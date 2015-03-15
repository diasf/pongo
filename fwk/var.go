package fwk

import "golang.org/x/mobile/f32"

type Drawable interface {
	Draw(modelView *f32.Mat4, ratio float64)
	OnAttached()
}

type Named interface {
	GetName() string
}

type NamedDrawable interface {
	Drawable
	Named
}

type Color struct {
	R, G, B, A float32
}

func (c Color) Slice() []float32 {
	return []float32{c.R, c.G, c.B, c.A}
}
