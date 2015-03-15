package fwk

import (
	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

type Scene struct {
	Width, Height float32
	root          *Node
	renderer      *Renderer
}

func (s *Scene) Init() (err error) {
	gl.ClearColor(0., 0., 0., 0.)

	if s.renderer, err = GetRenderer(); err != nil {
		return
	}

	s.resetView()

	s.root = NewNode(nil, "Root", Vector{0., 0., -1.})

	return
}

func (*Scene) Destroy() {

}

func (s *Scene) GetRoot() *Node {
	return s.root
}

func (s *Scene) Draw(ratio float64) {
	if s.Width != geom.Width.Px() || s.Height != geom.Height.Px() {
		s.resetView()
	}
	if s.Width == 0 || s.Height == 0 {
		return
	}
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	modelView := &f32.Mat4{}
	modelView.Identity()

	s.root.Draw(modelView, ratio)

	s.renderer.Execute()
}

func (s *Scene) resetView() {
	s.Width = geom.Width.Px()
	s.Height = geom.Height.Px()

	gl.Viewport(0, 0, int(s.Width), int(s.Height))

	s.renderer.StandardShader.SetParams(StandardProjection(ortho(-200, 200, 200, -200, .5, 100)))
}

func ortho(left, right, top, bottom, near, far float32) *f32.Mat4 {
	m := &f32.Mat4{}

	m[0][0] = 2 / (right - left)
	m[1][1] = 2 / (top - bottom)
	m[2][2] = -2 / (far - near)
	m[3][3] = 1

	m[0][3] = -(right + left) / (right - left)
	m[1][3] = -(top + bottom) / (top - bottom)
	m[2][3] = -(far + near) / (far - near)
	return m
}
