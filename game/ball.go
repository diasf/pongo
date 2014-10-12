package game

import (
	gl "github.com/chsc/gogl/gl21"
	"github.com/diasf/pongo/fwk"
	"time"
)

type Ball struct {
	node            *fwk.Node
	direction       fwk.Vector
	directionUpdate time.Time
	speed           float32
	size            gl.Float
}

func (b *Ball) Move(timeInNano int64) {
	x := gl.Float(b.speed * float32(timeInNano/100000000))
	b.node.Move(b.direction.Mult(x))
}

func NewBall(parent *fwk.Node, name string, position fwk.Vector, color fwk.Color, speed float32) *Ball {
	ball := &Ball{}
	ball.size = 10
	ball.node = fwk.NewNode(parent, name, position).AddDrawable(&fwk.Rectangle{ball.size, ball.size, color, "Ball"})
	ball.direction = fwk.Vector{-1, 1, 0}
	ball.speed = speed
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

func (b *Ball) HandleCollision(this fwk.CollisionObject, other fwk.CollisionObject) {
	pos := b.node.GetPosition()
	half := b.size / 2
	nearest := other.GetBoundingVolume().GetNearestTo(pos)
	left := nearest.X - pos.X - half
	right := nearest.X - pos.X + half
	top := nearest.Y - pos.Y + half
	bottom := nearest.Y - pos.Y - half
	m := min(left, right, top, bottom)
	if m == left || m == right {
		b.direction.X *= -1
	} else {
		b.direction.Y *= -1
	}
}

func min(val ...gl.Float) (rs gl.Float) {
	for i, v := range val {
		if i == 0 || v < rs {
			rs = v
		}
	}
	return
}
