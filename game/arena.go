package game

import (
	gl "github.com/chsc/gogl/gl21"
	"github.com/diasf/pongo/fwk"
)

// the arena has two bars
type Arena struct {
	top    *fwk.Node
	bottom *fwk.Node
}

func NewArena(parent *fwk.Node, ringsize float32, wallwidth float32, color fwk.Color) *Arena {
	arena := &Arena{}
	halfring := gl.Float(ringsize) / 2.
	bar := &fwk.Rectangle{gl.Float(ringsize), gl.Float(wallwidth), color, "BorderRect"}

	// top
	arena.top = fwk.NewNode(parent, "TopBorderN", fwk.Vector{0., halfring - gl.Float(wallwidth/2.), 0.}).AddDrawable(bar)

	// bottom
	arena.bottom = fwk.NewNode(parent, "BottomBorderN", fwk.Vector{0., -halfring + gl.Float(wallwidth/2.), 0.}).AddDrawable(bar)
	return arena
}

func (a *Arena) GetBoundingVolume() fwk.BoundingVolume {
	return nil
}
