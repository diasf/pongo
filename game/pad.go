package game

import (
	"fmt"
	gl "github.com/chsc/gogl/gl21"
	"github.com/diasf/pongo/fwk"
)

const (
	MOVING_UP = iota
	MOVING_DOWN
	MOVING_STOP
)

type Pad struct {
	node       *fwk.Node
	direction  int
	speed      float32
	width      gl.Float
	widthHalf  gl.Float
	height     gl.Float
	heightHalf gl.Float
}

func (p *Pad) SetDirection(dir int) {
	p.direction = dir
}

func (p *Pad) Move(timeInNano int64) {
	x := gl.Float(p.speed * float32(timeInNano/100000000))
	if p.direction == MOVING_DOWN {
		p.node.Move(&fwk.Vector{0, -x, 0})
	} else if p.direction == MOVING_UP {
		p.node.Move(&fwk.Vector{0, x, 0})
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
	pad.speed = speed

	box, _ := pad.GetBoundingVolumes()[0].(*fwk.BoundingBox)
	fmt.Println(name, ": ", box)

	return pad
}

func (p *Pad) GetBoundingVolumes() []fwk.BoundingVolume {
	pos := p.node.GetPosition()
	return []fwk.BoundingVolume{&fwk.BoundingBox{Left: pos.X - p.widthHalf, Right: pos.X + p.widthHalf, Top: pos.Y + p.heightHalf, Bottom: pos.Y - p.heightHalf}}
}

func (n *Pad) GetName() string {
	return n.node.GetName()
}
