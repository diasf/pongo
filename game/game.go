package game

import (
	"fmt"
	gl "github.com/chsc/gogl/gl21"
	"github.com/jteeuwen/glfw"
	"github.com/pongo/fwk"
)

type pongoGame struct {
	fwk.BaseGame
	playerOne *Pad
	playerTwo *Pad
	arena	  *Arena
	quit      bool
}

func NewPongoGame(width, height int) fwk.Game {
	pGame := &pongoGame{}
	pGame.quit = false
	pGame.Width = width
	pGame.Height = height
	pGame.Title = "PonGo"
	// register handlers
	pGame.SetGameSceneBuilder(pGame)
	pGame.SetGameUpdateHandler(pGame)
	pGame.SetKeyEventHandler(pGame)
	fmt.Println("PonGo game created")
	return pGame
}

func (g *pongoGame) Update(timeInNano int64) bool {
	if !g.quit {
		g.playerOne.Move(timeInNano)
		g.playerTwo.Move(timeInNano)
	}
	return !g.quit
}

func (g *pongoGame) OnKeyEvent(key, state int) {
	if state == glfw.KeyPress && key == glfw.KeyUp {
		g.playerOne.SetDirection(MOVING_UP)
	} else if state == glfw.KeyPress && key == glfw.KeyDown {
		g.playerOne.SetDirection(MOVING_DOWN)
	} else if state == glfw.KeyRelease && (key == glfw.KeyUp || key == glfw.KeyDown) {
		g.playerOne.SetDirection(MOVING_STOP)
	} else if state == glfw.KeyPress && key == glfw.KeyEsc {
		g.quit = true
	}
}

func (g *pongoGame) BuildGameScene() {
	root := g.GetScene().GetRoot()
	height := gl.Float(100.)
	width := gl.Float(10.)

	// player one pad
	g.playerOne = NewPad(root, "Player1Node", fwk.Vector{-190., -height / 2., 0.}, fwk.Color{1., 0., 0., 1.})
	// player two pad
	g.playerTwo = NewPad(root, "Player2Node", fwk.Vector{190. - width, -height / 2., 0.}, fwk.Color{0., 0., 1., 1.})
	// the ring
	g.arena = NewArena(root, 400, 5, fwk.Color{.5, .5, .1, 1})
}
