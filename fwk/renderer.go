package fwk

import (
	"log"

	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/gl"
	"golang.org/x/mobile/gl/glutil"
)

var renderer *Renderer

type Renderer struct {
	StandardShader *StandardShader
}

func (r *Renderer) Execute() {
	r.StandardShader.Execute()
}

func GetRenderer() (*Renderer, error) {
	if renderer != nil {
		return renderer, nil
	}

	if standardShader, err := newStandardShader(); err != nil {
		return nil, err
	} else {
		renderer = &Renderer{StandardShader: standardShader}
		return renderer, nil
	}
}

type StandardShader struct {
	program          gl.Program
	projection       gl.Uniform
	modelView        gl.Uniform
	position         gl.Attrib
	color            gl.Attrib
	projectionMatrix []float32
	listenerCount    int
	listeners        map[int]StandardListener
}

func newStandardShader() (s *StandardShader, err error) {
	s = &StandardShader{}
	if s.program, err = glutil.CreateProgram(vertexShader, fragmentShader); err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}
	s.position = gl.GetAttribLocation(s.program, "Position")
	s.color = gl.GetAttribLocation(s.program, "Color")
	s.projection = gl.GetUniformLocation(s.program, "Projection")
	s.modelView = gl.GetUniformLocation(s.program, "ModelView")
	s.listeners = map[int]StandardListener{}
	return
}

func (s *StandardShader) Execute() {
	if len(s.listeners) == 0 {
		return
	}

	gl.UseProgram(s.program)
	if s.projectionMatrix != nil && len(s.projectionMatrix) > 0 {
		gl.UniformMatrix4fv(s.projection, s.projectionMatrix)
	}

	for _, listener := range s.listeners {
		listener(s)
	}
}

type StandardShaderParam func(std *StandardShader)

func StandardProjection(mat *f32.Mat4) StandardShaderParam {
	return func(std *StandardShader) {
		std.projectionMatrix = conv(mat)
	}
}

func (s *StandardShader) SetParams(param ...StandardShaderParam) {
	if param != nil {
		for _, p := range param {
			p(s)
		}
	}
}

func (s *StandardShader) RegisterListener(listener StandardListener) int {
	id := s.incrListenerId()
	s.listeners[id] = listener
	return id
}

func (s *StandardShader) incrListenerId() int {
	s.listenerCount++
	return s.listenerCount
}

type StandardListener func(std *StandardShader)

func StandardDrawListener(listener func(modelView gl.Uniform, color, position gl.Attrib)) StandardListener {
	return func(std *StandardShader) {
		listener(std.modelView, std.color, std.position)
	}
}

func conv(mat *f32.Mat4) []float32 {
	rs := make([]float32, 16)
	for n, c := range mat[:] {
		for i, r := range c[:] {
			rs[i*4+n] = r
		}
	}
	return rs
}

const vertexShader = `#version 100
attribute vec4 Position;
uniform mat4 Projection;
uniform mat4 ModelView;
attribute vec4 Color;
varying vec4 DestinationColor;
 
void main(void) {
	DestinationColor = Color;
	gl_Position = Projection * ModelView * Position;
}
`

const fragmentShader = `#version 100
precision mediump float;
varying lowp vec4 DestinationColor;
 
void main(void) {
	gl_FragColor = DestinationColor;
}`
