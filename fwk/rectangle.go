package fwk

import (
	"encoding/binary"
	"fmt"
	"log"

	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/gl"

	"github.com/diasf/pongo/fwk/tex"
)

type Rectangle struct {
	Width          float32
	Height         float32
	Color          Color
	Texture        *tex.Texture
	TexCoord       *RectangleTexCoord
	Name           string
	vertexBuffer   gl.Buffer
	indexBuffer    gl.Buffer
	modelView      []float32
	indicesCount   int
	verticesStride int
}

type RectangleTexCoord struct {
	TopLeft     Vector2
	TopRight    Vector2
	BottomLeft  Vector2
	BottomRight Vector2
}

func (r *Rectangle) OnAttached() {
	w2 := r.Width / float32(2.)
	h2 := r.Height / float32(2.)

	texCoord := r.buildTexCoord()

	vertices := []float32{}
	vertices = append(vertices, w2, -h2, 0, texCoord.BottomRight.X, texCoord.BottomRight.Y)
	vertices = append(vertices, w2, h2, 0, texCoord.TopRight.X, texCoord.TopRight.Y)
	vertices = append(vertices, -w2, h2, 0, texCoord.TopLeft.X, texCoord.TopLeft.Y)
	vertices = append(vertices, -w2, -h2, 0, texCoord.BottomLeft.X, texCoord.BottomLeft.Y)
	r.verticesStride = 5

	var indices = []byte{
		0, 1, 2,
		2, 3, 0,
	}

	r.indicesCount = len(indices)

	r.vertexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, f32.Bytes(binary.LittleEndian, vertices...), gl.STATIC_DRAW)

	r.indexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indices, gl.STATIC_DRAW)

	if r.Texture != nil {
		if err := r.Texture.Upload(); err != nil {
			log.Fatalln(err)
			r.Texture = nil
		}
	}

	if renderer, err := GetRenderer(); err == nil {
		renderer.StandardShader.RegisterListener(StandardDrawListener(r.standardShaderExecution))
	}
}

func (r *Rectangle) buildTexCoord() *RectangleTexCoord {
	if r.TexCoord != nil {
		return r.TexCoord
	}
	return &RectangleTexCoord{
		BottomLeft:  Vector2{0, 0},
		BottomRight: Vector2{0, 1},
		TopLeft:     Vector2{1, 0},
		TopRight:    Vector2{1, 1},
	}
}

func (r *Rectangle) Draw(modelView *f32.Mat4, ratio float64) {
	r.modelView = conv(modelView)
}

func (r *Rectangle) standardShaderExecution(modelView, withTexture, texture gl.Uniform, color, texCoord, position gl.Attrib) {
	gl.UniformMatrix4fv(modelView, r.modelView)

	gl.BindBuffer(gl.ARRAY_BUFFER, r.vertexBuffer)
	gl.EnableVertexAttribArray(position)
	gl.VertexAttribPointer(position, 3, gl.FLOAT, false, r.verticesStride*4, 0)

	if r.Texture != nil {
		gl.Uniform1i(withTexture, 1)
		gl.EnableVertexAttribArray(texCoord)
		gl.VertexAttribPointer(texCoord, 2, gl.FLOAT, false, r.verticesStride*4, 3*4)

		gl.ActiveTexture(gl.TEXTURE0)
		r.Texture.Bind()
		gl.Uniform1i(texture, 0)
	} else {
		gl.Uniform1i(withTexture, 0)
	}

	gl.VertexAttrib4fv(color, r.Color.Slice())

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.indexBuffer)
	gl.DrawElements(gl.TRIANGLES, r.indicesCount, gl.UNSIGNED_BYTE, 0)

	gl.DisableVertexAttribArray(position)
	gl.DisableVertexAttribArray(color)
}

func (r *Rectangle) GetName() string {
	return r.Name
}

func errDrain() string {
	var errs []gl.Enum
	for {
		e := gl.GetError()
		if e == 0 {
			break
		}
		errs = append(errs, e)
	}
	if len(errs) > 0 {
		return fmt.Sprintf(" error: %v", errs)
	}
	return ""
}
