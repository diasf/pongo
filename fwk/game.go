package fwk

import (
	"errors"
	"fmt"
	gl "github.com/chsc/gogl/gl21"
	"github.com/jteeuwen/glfw"
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
	OnKeyEvent(key, state int)
}

type Game interface {
	Start()
}

type BaseGame struct {
	// width & height of the window
	Width, Height int
	// title to show for the window
	Title string
	// Handler to build the scene
	sceneBuilder GameSceneBuilder
	// Handler to update the game logic
	gameUpdateHandler GameUpdateHandler
	// Handler to handle key events
	keyEventHandler KeyEventHandler
	// the scene object
	scene *Scene
}

func (g *BaseGame) Start() {
	if err := glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "glfw: %s\n", err)
		return
	}
	defer glfw.Terminate()

	glfw.OpenWindowHint(glfw.WindowNoResize, 1)

	if err := glfw.OpenWindow(g.Width, g.Height, 0, 0, 0, 0, 16, 0, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "glfw: %s\n", err)
		return
	}
	defer glfw.CloseWindow()

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle(g.Title)
	if g.keyEventHandler != nil {
		glfw.SetKeyCallback(func(key, state int) { g.keyEventHandler.OnKeyEvent(key, state) })
	}

	if err := gl.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "gl: %s\n", err)
	}

	g.scene = &Scene{Width: gl.Sizei(g.Width), Height: gl.Sizei(g.Height)}
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
	fixedStep := int64(100000000)
	cont := true
	for glfw.WindowParam(glfw.Opened) == 1 && cont {

		delta := timer.Delta().Nanoseconds()

		if delta > 250000000 {
			delta = 250000000
		}
		accu += delta

		for cont && accu >= fixedStep {
			accu -= fixedStep
			if g.gameUpdateHandler != nil {
				cont = g.gameUpdateHandler.Update(fixedStep)
			}
		}

		if cont {
			g.scene.Draw(float64(accu) / float64(fixedStep))
			glfw.SwapBuffers()
		}
	}

	// close window if still open
	if glfw.WindowParam(glfw.Opened) == 1 {
		glfw.CloseWindow()
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
