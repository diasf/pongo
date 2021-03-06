package fwk

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"time"

	gl "github.com/chsc/gogl/gl21"
	glfw "github.com/go-gl/glfw3"
)

type GameSceneBuilder interface {
	BuildGameScene()
}

type GameUpdateHandler interface {
	Update(duration time.Duration) bool
}

type KeyEventHandler interface {
	OnKeyEvent(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
}

type Game interface {
	Start()
}

type BaseGame struct {
	width, height     int               // width & height of the window
	title             string            // title to show for the window
	fixedStep         time.Duration     // fixed time at which the GameUpdateHandler will be invoked
	maxDelta          time.Duration     // maximum delta between updates
	sceneBuilder      GameSceneBuilder  // Handler to build the scene
	gameUpdateHandler GameUpdateHandler // Handler to update the game logic
	keyEventHandler   KeyEventHandler   // Handler to handle key events
	collisionDetector CollisionDetector // Used for checking collisions
	scene             *Scene            // the scene object
	window            *glfw.Window
}

func NewBaseGame(width, height int, title string) BaseGame {
	bg := BaseGame{
		width:             width,
		height:            height,
		title:             title,
		fixedStep:         time.Duration(100 * time.Millisecond),
		maxDelta:          time.Duration(250 * time.Millisecond),
		sceneBuilder:      nil,
		gameUpdateHandler: nil,
		keyEventHandler:   nil,
		collisionDetector: NewSimpleCollisionDetector(),
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

	var err error
	if g.window, err = glfw.CreateWindow(g.width, g.height, g.title, nil, nil); err != nil {
		fmt.Fprintf(os.Stderr, "glfw: %s\n", err)
		return
	}

	g.window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if g.keyEventHandler != nil {
		g.window.SetKeyCallback(func(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
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

	g.gameLoop()
}

func (g *BaseGame) gameLoop() {
	timer := NewTimer()
	accu := time.Duration(0)
	cont := true
	ratio := float64(0)
	for !g.window.ShouldClose() && cont {

		delta := timer.Delta()

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
			ratio = accu.Seconds() / g.fixedStep.Seconds()
			g.GetCollisionDetector().Check()
			g.scene.Draw(ratio)
			g.window.SwapBuffers()
			glfw.PollEvents()
		}
	}
}

func (g *BaseGame) GetScene() *Scene {
	return g.scene
}

func (g *BaseGame) GetCollisionDetector() CollisionDetector {
	return g.collisionDetector
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
