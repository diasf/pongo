package fwk

import (
	gl "github.com/chsc/gogl/gl21"
)

type Rectangle struct {
	Width  gl.Float
	Height gl.Float
	Color  Color
	Name   string
}

func (r *Rectangle) Draw(ratio float64) {
	w2 := r.Width / gl.Float(2.)
	h2 := r.Height / gl.Float(2.)
	gl.Color4f(r.Color.R, r.Color.G, r.Color.B, r.Color.A)
	gl.Begin(gl.QUADS)
	gl.Normal3f(0, 0, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-w2, -h2, 0.)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-w2, h2, 0.)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(w2, h2, 0.)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(w2, -h2, 0.)
	gl.End()
}

func (r *Rectangle) GetName() string {
	return r.Name
}
