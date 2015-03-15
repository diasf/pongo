package fwk

import (
	"encoding/binary"

	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/gl"
)

type Rectangle struct {
	Width        float32
	Height       float32
	Color        Color
	Name         string
	vertexBuffer gl.Buffer
	indexBuffer  gl.Buffer
	modelView    []float32
	indicesCount int
}

func (r *Rectangle) OnAttached() {
	w2 := r.Width / float32(2.)
	h2 := r.Height / float32(2.)
	col := r.Color.Slice()

	vertices := []float32{}
	vertices = append(vertices, append([]float32{w2, -h2, 0}, col...)...)
	vertices = append(vertices, append([]float32{w2, h2, 0}, col...)...)
	vertices = append(vertices, append([]float32{-w2, h2, 0}, col...)...)
	vertices = append(vertices, append([]float32{-w2, -h2, 0}, col...)...)

	var indices = []byte{
		0, 1, 2,
		2, 3, 0,
	}

	r.indicesCount = len(indices)

	r.vertexBuffer = gl.GenBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.STATIC_DRAW, f32.Bytes(binary.LittleEndian, vertices...))

	r.indexBuffer = gl.GenBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, gl.STATIC_DRAW, indices)

	if renderer, err := GetRenderer(); err == nil {
		renderer.StandardShader.RegisterListener(StandardDrawListener(r.standardShaderExecution))
	}
}

func (r *Rectangle) Draw(modelView *f32.Mat4, ratio float64) {
	r.modelView = conv(modelView)
}

func (r *Rectangle) standardShaderExecution(modelView gl.Uniform, color, position gl.Attrib) {
	gl.UniformMatrix4fv(modelView, r.modelView)

	gl.BindBuffer(gl.ARRAY_BUFFER, r.vertexBuffer)
	gl.EnableVertexAttribArray(position)
	gl.VertexAttribPointer(position, 3, gl.FLOAT, false, 7*4, 0)
	gl.EnableVertexAttribArray(color)
	gl.VertexAttribPointer(color, 4, gl.FLOAT, false, 7*4, 3*4)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.indexBuffer)
	gl.DrawElements(gl.TRIANGLES, gl.UNSIGNED_BYTE, 0, r.indicesCount)

	gl.DisableVertexAttribArray(position)
	gl.DisableVertexAttribArray(color)
}

func (r *Rectangle) GetName() string {
	return r.Name
}
