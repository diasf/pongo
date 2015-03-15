package game

import (
	"fmt"
	"time"

	"github.com/diasf/pongo/fwk"
	"golang.org/x/mobile/event"
)

type pongoGame struct {
	fwk.BaseGame
	playerOne    *Pad
	playerTwo    *Pad
	ball         *Ball
	arena        *Arena
	quit         bool
	reactionTime time.Duration
	lastTouch    event.Touch
}

func NewPongoGame(width, height int) fwk.Game {
	pGame := &pongoGame{}
	// initialize base game..
	pGame.BaseGame = fwk.NewBaseGame("PonGo")
	pGame.quit = false
	pGame.reactionTime = time.Duration(100) * time.Millisecond
	// register handlers
	pGame.SetGameSceneBuilder(pGame)
	pGame.SetGameUpdateHandler(pGame)
	pGame.SetTouchEventHandler(pGame)
	pGame.GetCollisionDetector().AddCollisionHandler(pGame)
	fmt.Println("PonGo game created")
	return pGame
}

func (g *pongoGame) Update(duration time.Duration) bool {
	if !g.quit {
		g.ball.Move(duration)
		g.computerPlayerTwo()
		g.playerTwo.Move(duration)
	}
	return !g.quit
}

func (g *pongoGame) computerPlayerTwo() {
	if !time.Now().After(g.ball.directionUpdate.Add(g.reactionTime)) {
		return
	}
	p := g.playerTwo
	myPos := p.node.GetPosition().Y
	if myPos < g.ball.node.GetPosition().Y {
		if !p.IsDirectionLockedOn(MOVING_UP) {
			p.SetDirection(MOVING_UP)
		}
	} else if !p.IsDirectionLockedOn(MOVING_DOWN) {
		p.SetDirection(MOVING_DOWN)
	}
}

func (g *pongoGame) OnTouchEvent(t event.Touch) {
	if t.Type == event.TouchStart {
		g.lastTouch = t
	} else if t.Type == event.TouchMove {
		g.playerOne.MoveY(g.lastTouch.Loc.Y.Px() - t.Loc.Y.Px())
		g.lastTouch = t
	}
}

func (g *pongoGame) BuildGameScene() {
	root := g.GetScene().GetRoot()

	// ball
	g.ball = NewBall(root, "BallNode", fwk.Vector{0., 0., 0.}, fwk.Color{0., 1., 0., 1.}, 3.5)
	g.GetCollisionDetector().AddCollidable(g.ball)
	// player one pad
	g.playerOne = NewPad(root, "Player1Node", fwk.Vector{-185., 0., 0.}, fwk.Color{1., 0., 0., 1.}, 5.)
	g.GetCollisionDetector().AddCollidable(g.playerOne)
	// player two pad
	g.playerTwo = NewPad(root, "Player2Node", fwk.Vector{185., 0., 0.}, fwk.Color{0., 0., 1., 1.}, 5.)
	g.GetCollisionDetector().AddCollidable(g.playerTwo)
	// the ring
	g.arena = NewArena(root, "Arena", 400, 10, fwk.Color{.5, .5, .1, 1})
	g.GetCollisionDetector().AddCollidable(g.arena)
}

func (g *pongoGame) HandleCollision(one, two fwk.CollisionObject) {
	if pad, ok := one.GetObject().(*Pad); ok && two.GetObject() == g.arena {
		if !pad.IsDirectionLocked() {
			pad.SetDirection(MOVING_STOP)
			if g.arena.GetTopBoundingVolume().CollidesWith(one.GetBoundingVolume()) {
				pad.LockDirection(MOVING_UP)
			} else {
				pad.LockDirection(MOVING_DOWN)
			}
		}
	} else if ball, ok := one.GetObject().(*Ball); ok {
		ball.HandleCollision(g.getBallHit(two))
	} else if ball, ok := two.GetObject().(*Ball); ok {
		ball.HandleCollision(g.getBallHit(one))
	}
}

func (g *pongoGame) getBallHit(col fwk.CollisionObject) (x, y int) {
	if pad, ok := col.GetObject().(*Pad); ok {
		if pad == g.playerOne {
			x = 1
		} else {
			x = -1
		}
	} else {
		nearest := col.GetBoundingVolume().GetNearestTo(g.ball.node.GetPosition())
		if nearest.Y > 10 {
			y = -1
		} else {
			y = 1
		}
	}

	return
}
