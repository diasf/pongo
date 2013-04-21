package fwk

import (
	gl "github.com/chsc/gogl/gl21"
)

type Scene struct {
	Width, Height gl.Sizei
	root          *Node
}

func (s *Scene) Init() (err error) {
	gl.ClearColor(0., 0., 0., 0.)
	gl.ClearDepth(1)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.Viewport(0, 0, s.Width, s.Height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Frustum(-100, 100, -100, 100, .5, 100.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	s.root = NewNode("Root", Vector{0., 0., -1.})

	return
}

func (*Scene) Destroy() {

}

func (s *Scene) GetRoot() *Node {
	return s.root
}

func (s *Scene) Draw(ratio float64) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	s.root.Draw(ratio)
}
