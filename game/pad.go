package game

import (
	"time"

	"github.com/diasf/pongo/fwk"
	"github.com/diasf/pongo/fwk/tex"
)

const (
	MOVING_UP = iota
	MOVING_DOWN
	MOVING_STOP
)

type Pad struct {
	node            *fwk.Node
	direction       int
	lockedDirection int
	speed           float32
	width           float32
	widthHalf       float32
	height          float32
	heightHalf      float32
}

func (p *Pad) SetDirection(dir int) {
	p.direction = dir
}

func (p *Pad) GetDirection() int {
	return p.direction
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

func (p *Pad) Move(duration time.Duration) {
	x := float32(p.speed * float32(duration.Seconds()))
	if p.direction == MOVING_DOWN && !p.IsDirectionLockedOn(MOVING_DOWN) {
		p.node.Move(&fwk.Vector{0, -x, 0})
		p.UnLockDirection()
	} else if p.direction == MOVING_UP && !p.IsDirectionLockedOn(MOVING_UP) {
		p.node.Move(&fwk.Vector{0, x, 0})
		p.UnLockDirection()
	}
}

func (p *Pad) MoveY(diff float32) {
	p.node.GetPosition().Y += diff
}

func (p *Pad) Rotate(deg float32, up fwk.Vector) {
	p.node.Rotate(deg, up)
}

func NewPad(parent *fwk.Node, name string, position fwk.Vector, speed float32, texturePng string) *Pad {
	pad := &Pad{}
	pad.width = 10
	pad.widthHalf = pad.width / 2.
	pad.height = 100
	pad.heightHalf = pad.height / 2.

	texture := tex.NewTextureFromPNGFile(texturePng)
	texture.SetRepeat()
	texture.SetMagFilterNearest()
	texture.SetMinFilterNearest()

	texCoord := &fwk.RectangleTexCoord{
		BottomLeft:  fwk.Vector2{.0, .0},
		BottomRight: fwk.Vector2{.07, .0},
		TopLeft:     fwk.Vector2{.0, .3},
		TopRight:    fwk.Vector2{.07, .3},
	}

	pad.node = fwk.NewNode(parent, name, position).AddDrawable(&fwk.Rectangle{Width: pad.width, Height: pad.height, Name: "Rect", Texture: texture, TexCoord: texCoord})
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

func (p *Pad) GetTop() float32 {
	return p.node.GetPosition().Y + p.heightHalf
}

func (p *Pad) GetBottom() float32 {
	return p.node.GetPosition().Y - p.heightHalf
}
