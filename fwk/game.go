package fwk

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

type GameSceneBuilder interface {
	BuildGameScene()
}

type GameUpdateHandler interface {
	Update(duration time.Duration) bool
}

type TouchEventHandler interface {
	OnTouchEvent(t event.Touch)
}

type Game interface {
	Start()
}

type BaseGame struct {
	width, height     float32           // width & height of the window
	fixedStep         time.Duration     // fixed time at which the GameUpdateHandler will be invoked
	maxDelta          time.Duration     // maximum delta between updates
	sceneBuilder      GameSceneBuilder  // Handler to build the scene
	gameUpdateHandler GameUpdateHandler // Handler to update the game logic
	touchEventHandler TouchEventHandler // Handler to handle touch events
	collisionDetector CollisionDetector // Used for checking collisions
	scene             *Scene            // the scene object
}

func NewBaseGame(title string) BaseGame {
	bg := BaseGame{
		fixedStep:         time.Duration(100 * time.Millisecond),
		maxDelta:          time.Duration(250 * time.Millisecond),
		sceneBuilder:      nil,
		gameUpdateHandler: nil,
		touchEventHandler: nil,
		collisionDetector: NewSimpleCollisionDetector(),
		scene:             nil,
	}
	return bg
}

func (g *BaseGame) Start() {
	app.Run(app.Callbacks{
		Start: g.onStart,
		Stop:  g.onStop,
		Draw:  g.onDraw(),
		Touch: g.onTouch,
	})
}

func (g *BaseGame) onStart() {
	g.width = geom.Width.Px()
	g.height = geom.Height.Px()

	g.scene = &Scene{Width: g.width, Height: g.height}
	if err := g.scene.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "init: %s\n", err)
		return
	}

	if g.sceneBuilder != nil {
		g.sceneBuilder.BuildGameScene()
	}
}

func (g *BaseGame) onStop() {
	if g.scene != nil {
		g.scene.Destroy()
	}
}

func (g *BaseGame) onTouch(t event.Touch) {
	if g.touchEventHandler != nil {
		g.touchEventHandler.OnTouchEvent(t)
	}
}

func (g *BaseGame) onDraw() func() {
	timer := NewTimer()
	accu := time.Duration(0)
	cont := true
	ratio := float64(0)
	return func() {
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

// Handler to handle touch events
func (g *BaseGame) SetTouchEventHandler(handler TouchEventHandler) {
	g.touchEventHandler = handler
}

func createTextureFromPng(r io.Reader) (textureId gl.Texture, err error) {
	img, err := png.Decode(r)
	if err != nil {
		return
	}

	rgbaImg, ok := img.(*image.NRGBA)
	if !ok {
		return textureId, errors.New("texture must be an NRGBA image")
	}

	textureId = gl.GenTexture()
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
	gl.TexImage2D(gl.TEXTURE_2D, 0, imgWidth, imgHeight, gl.RGBA, gl.UNSIGNED_BYTE, data)
	return
}
