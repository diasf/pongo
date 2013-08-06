package fwk

import (
	"errors"
	"fmt"
	gl "github.com/chsc/gogl/gl21"
	glfw "github.com/go-gl/glfw3"
	"image"
	"image/png"
	"io"
	"os"
)

type GameSceneBuilder interface {
	BuildGameScene()
}

type GameUpdateHandler interface {
	Update(timeInNano int64) bool
}

type KeyEventHandler interface {
	OnKeyEvent(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
}

type Game interface {
	Start()
}

type BaseGame struct {
	// width & height of the window
	width, height int
	// title to show for the window
	title string
	// fixed time at which the GameUpdateHandler will be invoked
	fixedStep int64
	// maximum delta between updates
	maxDelta int64
	// Handler to build the scene
	sceneBuilder GameSceneBuilder
	// Handler to update the game logic
	gameUpdateHandler GameUpdateHandler
	// Handler to handle key events
	keyEventHandler KeyEventHandler
	// the scene object
	scene *Scene
}

func NewBaseGame(width, height int, title string) BaseGame {
	bg := BaseGame{
		width:             width,
		height:            height,
		title:             title,
		fixedStep:         100000000,
		maxDelta:          250000000,
		sceneBuilder:      nil,
		gameUpdateHandler: nil,
		keyEventHandler:   nil,
		scene:             nil,
	}
	return bg
}

func (g *BaseGame) Start() {
	if !glfw.Init() {
		fmt.Fprintf(os.Stderr, "Error initializing glfw\n")
		return
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, 0)
	glfw.WindowHint(glfw.Resizable, 0)

	var window * glfw.Window
	var err error
	if window, err = glfw.CreateWindow(g.width, g.height, g.title, nil, nil); err != nil {
		fmt.Fprintf(os.Stderr, "glfw: %s\n", err)
		return
	}

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if g.keyEventHandler != nil {
		window.SetKeyCallback(func(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			g.keyEventHandler.OnKeyEvent(key, scancode, action, mods)
		})
	}

	if err := gl.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "gl: %s\n", err)
	}

	g.scene = &Scene{Width: gl.Sizei(g.width), Height: gl.Sizei(g.height)}
	if err := g.scene.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "init: %s\n", err)
		return
	}
	defer g.scene.Destroy()

	if g.sceneBuilder != nil {
		g.sceneBuilder.BuildGameScene()
	}

	timer := NewTimer()
	accu := int64(0)
	cont := true
	for !window.ShouldClose() && cont {

		delta := timer.Delta().Nanoseconds()

		if delta > g.maxDelta {
			delta = g.maxDelta
		}
		accu += delta

		for cont && accu >= g.fixedStep {
			accu -= g.fixedStep
			if g.gameUpdateHandler != nil {
				cont = g.gameUpdateHandler.Update(g.fixedStep)
			}
		}

		if cont {
			g.scene.Draw(float64(accu) / float64(g.fixedStep))
			window.SwapBuffers()
			glfw.PollEvents()
		}
	}
}

func (g *BaseGame) GetScene() *Scene {
	return g.scene
}

func (g *BaseGame) SetGameSceneBuilder(handler GameSceneBuilder) {
	g.sceneBuilder = handler
}

// Handler to update the game logic
func (g *BaseGame) SetGameUpdateHandler(handler GameUpdateHandler) {
	g.gameUpdateHandler = handler
}

// Handler to handle key events
func (g *BaseGame) SetKeyEventHandler(handler KeyEventHandler) {
	g.keyEventHandler = handler
}

func createTexture(r io.Reader) (textureId gl.Uint, err error) {
	img, err := png.Decode(r)
	if err != nil {
		return 0, err
	}

	rgbaImg, ok := img.(*image.NRGBA)
	if !ok {
		return 0, errors.New("texture must be an NRGBA image")
	}

	gl.GenTextures(1, &textureId)
	gl.BindTexture(gl.TEXTURE_2D, textureId)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	// flip image: first pixel is lower left corner
	imgWidth, imgHeight := img.Bounds().Dx(), img.Bounds().Dy()
	data := make([]byte, imgWidth*imgHeight*4)
	lineLen := imgWidth * 4
	dest := len(data) - lineLen
	for src := 0; src < len(rgbaImg.Pix); src += rgbaImg.Stride {
		copy(data[dest:dest+lineLen], rgbaImg.Pix[src:src+rgbaImg.Stride])
		dest -= lineLen
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, 4, gl.Sizei(imgWidth), gl.Sizei(imgHeight), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Pointer(&data[0]))

	return textureId, nil
}
