package game

import (
	"time"

	"github.com/diasf/pongo/fwk"
)

type Ball struct {
	node            *fwk.Node
	direction       fwk.Vector
	directionUpdate time.Time
	speed           float32
	size            float32
}

func (b *Ball) Move(duration time.Duration) {
	b.node.Move(b.direction.Multiplication(float32(b.speed * float32(duration.Seconds()))))
}

func NewBall(parent *fwk.Node, name string, position fwk.Vector, color fwk.Color, speed float32) *Ball {
	ball := &Ball{}
	ball.size = 10
	ball.node = fwk.NewNode(parent, name, position).AddDrawable(&fwk.Rectangle{Width: ball.size, Height: ball.size, Color: color, Name: "Ball"})
	ball.direction = fwk.Vector{-1, 1, 0}
	ball.speed = speed
	go ball.speedIncrement(time.NewTicker(time.Duration(time.Second * 5)))
	return ball
}

func (b *Ball) GetBoundingVolumes() []fwk.BoundingVolume {
	pos := b.node.GetPosition()
	half := b.size / 2
	return []fwk.BoundingVolume{&fwk.BoundingBox{Left: pos.X - half, Right: pos.X + half, Top: pos.Y + half, Bottom: pos.Y - half}}
}

func (b *Ball) GetName() string {
	return b.node.GetName()
}

func (b *Ball) HandleCollision(x, y int) {
	if x != 0 {
		b.direction.X = float32(x)
	}

	if y != 0 {
		b.direction.Y = float32(y)
	}
}

func (b *Ball) speedIncrement(ticker *time.Ticker) {
	for _ = range ticker.C {
		if b.speed < 20 {
			b.speed += 0.5
		} else {
			ticker.Stop()
			return
		}
	}
}

func min(val ...float32) (rs float32) {
	for i, v := range val {
		if i == 0 || v < rs {
			rs = v
		}
	}
	return
}
