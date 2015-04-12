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
	withTexture      gl.Uniform
	texture          gl.Uniform
	position         gl.Attrib
	color            gl.Attrib
	texCoord         gl.Attrib
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
	s.position = gl.GetAttribLocation(s.program, "a_position")
	s.color = gl.GetAttribLocation(s.program, "a_color")
	s.texCoord = gl.GetAttribLocation(s.program, "a_texCoord")

	s.projection = gl.GetUniformLocation(s.program, "u_projection")
	s.modelView = gl.GetUniformLocation(s.program, "u_modelView")
	s.withTexture = gl.GetUniformLocation(s.program, "u_withTexture")
	s.texture = gl.GetUniformLocation(s.program, "u_texture")

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

func StandardDrawListener(listener func(modelView, withTexture, texture gl.Uniform, color, texCoord, position gl.Attrib)) StandardListener {
	return func(std *StandardShader) {
		listener(std.modelView, std.withTexture, std.texture, std.color, std.texCoord, std.position)
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
uniform mat4 u_projection;
uniform mat4 u_modelView;

attribute vec4 a_position;
attribute vec4 a_color;
attribute vec2 a_texCoord;

varying vec4 v_destinationColor;
varying vec2 v_texCoord;
 
void main(void) {
	v_destinationColor = a_color;
	v_texCoord = a_texCoord;
	gl_Position = u_projection * u_modelView * a_position;
}
`

const fragmentShader = `#version 100
precision mediump float;

varying lowp vec4 v_destinationColor;
varying vec2 v_texCoord;

uniform bool u_withTexture;
uniform sampler2D u_texture;
 
void main(void) {
	if(u_withTexture) {
		//gl_FragColor = texture2D(u_texture, v_texCoord);
		gl_FragColor = v_destinationColor;
	} else {
		gl_FragColor = v_destinationColor;
	}
}`
