package game

import (
	"fmt"
	"github.com/diasf/pongo/fwk"
	glfw "github.com/go-gl/glfw3"
)

type pongoGame struct {
	fwk.BaseGame
	playerOne *Pad
	playerTwo *Pad
	arena     *Arena
	quit      bool
}

func NewPongoGame(width, height int) fwk.Game {
	pGame := &pongoGame{}
	// initialize base game..
	pGame.BaseGame = fwk.NewBaseGame(width, height, "PonGo")
	pGame.quit = false
	// register handlers
	pGame.SetGameSceneBuilder(pGame)
	pGame.SetGameUpdateHandler(pGame)
	pGame.SetKeyEventHandler(pGame)
	pGame.GetCollisionDetector().AddCollisionHandler(pGame)
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

func (g *pongoGame) OnKeyEvent(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press && key == glfw.KeyUp {
		g.playerOne.SetDirection(MOVING_UP)
	} else if action == glfw.Press && key == glfw.KeyDown {
		g.playerOne.SetDirection(MOVING_DOWN)
	} else if action == glfw.Release && (key == glfw.KeyUp || key == glfw.KeyDown) {
		g.playerOne.SetDirection(MOVING_STOP)
	} else if action == glfw.Press && key == glfw.KeyEscape {
		g.quit = true
	}
}

func (g *pongoGame) BuildGameScene() {
	root := g.GetScene().GetRoot()

	// player one pad
	g.playerOne = NewPad(root, "Player1Node", fwk.Vector{-185., 0., 0.}, fwk.Color{1., 0., 0., 1.})
	g.GetCollisionDetector().AddCollidable(g.playerOne)
	// player two pad
	g.playerTwo = NewPad(root, "Player2Node", fwk.Vector{185., 0., 0.}, fwk.Color{0., 0., 1., 1.})
	g.GetCollisionDetector().AddCollidable(g.playerTwo)
	// the ring
	g.arena = NewArena(root, 400, 5, fwk.Color{.5, .5, .1, 1})
	g.GetCollisionDetector().AddCollidable(g.arena)
}

func (g *pongoGame) HandleCollision(one, two fwk.Collidable) {
	fmt.Println("Collision between", one, " and ", two)
}
