package game

import (
	gl "github.com/chsc/gogl/gl21"
	"github.com/diasf/pongo/fwk"
	"time"
)

const (
	MOVING_UP = iota
	MOVING_DOWN
	MOVING_STOP
)

type Pad struct {
	node            *fwk.Node
	direction       int
	directionUpdate time.Time
	lockedDirection int
	speed           float32
	width           gl.Float
	widthHalf       gl.Float
	height          gl.Float
	heightHalf      gl.Float
}

func (p *Pad) SetDirection(dir int) {
	p.direction = dir
	p.directionUpdate = time.Now()
}

func (p *Pad) GetDirection() int {
	return p.direction
}

func (p *Pad) GetLastDirectionUpdate() time.Time {
	return p.directionUpdate
}

func (p *Pad) UnLockDirection() {
	p.lockedDirection = -1
}

func (p *Pad) LockDirection(dir int) {
	p.lockedDirection = dir
}

func (p *Pad) IsDirectionLocked() bool {
	return p.lockedDirection >= 0
}

func (p *Pad) IsDirectionLockedOn(dir int) bool {
	return p.lockedDirection == dir
}

func (p *Pad) Move(timeInNano int64) {
	x := gl.Float(p.speed * float32(timeInNano/100000000))
	if p.direction == MOVING_DOWN && !p.IsDirectionLockedOn(MOVING_DOWN) {
		p.node.Move(&fwk.Vector{0, -x, 0})
		p.UnLockDirection()
	} else if p.direction == MOVING_UP && !p.IsDirectionLockedOn(MOVING_UP) {
		p.node.Move(&fwk.Vector{0, x, 0})
		p.UnLockDirection()
	}
}

func (p *Pad) Rotate(deg float32, up fwk.Vector) {
	p.node.Rotate(deg, up)
}

func NewPad(parent *fwk.Node, name string, position fwk.Vector, color fwk.Color, speed float32) *Pad {
	pad := &Pad{}
	pad.width = 10
	pad.widthHalf = pad.width / 2.
	pad.height = 100
	pad.heightHalf = pad.height / 2.
	pad.node = fwk.NewNode(parent, name, position).AddDrawable(&fwk.Rectangle{pad.width, pad.height, color, "Rect"})
	pad.direction = MOVING_STOP
	pad.lockedDirection = -1
	pad.speed = speed
	return pad
}

func (p *Pad) GetBoundingVolumes() []fwk.BoundingVolume {
	pos := p.node.GetPosition()
	return []fwk.BoundingVolume{&fwk.BoundingBox{Left: pos.X - p.widthHalf, Right: pos.X + p.widthHalf, Top: pos.Y + p.heightHalf, Bottom: pos.Y - p.heightHalf}}
}

func (n *Pad) GetName() string {
	return n.node.GetName()
}

func (n *Pad) IsOver(p fwk.Vector) bool {
	return n.GetTop() >= p.Y
}

func (n *Pad) IsUnder(p fwk.Vector) bool {
	return n.GetTop() <= p.Y
}

func (p *Pad) GetTop() gl.Float {
	return p.node.GetPosition().Y + p.heightHalf
}

func (p *Pad) GetBottom() gl.Float {
	return p.node.GetPosition().Y - p.heightHalf
}
