package game

import (
	gl "github.com/chsc/gogl/gl21"
	"github.com/pongo/fwk"
)

const (
	MOVING_UP = iota
	MOVING_DOWN
	MOVING_STOP
)

type Pad struct {
	node      *fwk.Node
	direction int
	speed     float32
}

func (p *Pad) SetDirection(dir int) {
	p.direction = dir
}

func (p *Pad) Move(timeInNano int64) {
	x := gl.Float(p.speed * float32(timeInNano/100000000))
	if p.direction == MOVING_DOWN {
		p.node.Position.Y -= x
	} else if p.direction == MOVING_UP {
		p.node.Position.Y += x
	}
}

func NewPad(parent *fwk.Node, name string, position fwk.Vector, color fwk.Color) *Pad {
	pad := &Pad{}
	pad.node = fwk.NewNode(name, position).AddDrawable(&fwk.Rectangle{10, 100, color, "Rect"})
	pad.direction = MOVING_STOP
	pad.speed = 2.
	parent.AddDrawable(pad.node)
	return pad
}
