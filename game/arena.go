package game

import (
	gl "github.com/chsc/gogl/gl21"
	"github.com/pongo/fwk"
)

// the arena has two bars
type Arena struct {
	top      *fwk.Node
	bottom   *fwk.Node
}

func NewArena(parent * fwk.Node, ringsize float32, wallwidth float32, color fwk.Color) *Arena {
	arena := &Arena{}
	halfring := gl.Float(ringsize) / 2.
	bar := &fwk.Rectangle{gl.Float(ringsize), gl.Float(wallwidth), color, "BorderRect"}

	// top
	arena.top = fwk.NewNode("TopBorderN", fwk.Vector{-halfring, halfring - gl.Float(wallwidth), 0.}).AddDrawable(bar)
	parent.AddDrawable(arena.top)

	// bottom
	arena.bottom = fwk.NewNode("BottomBorderN", fwk.Vector{-halfring, -halfring, 0.}).AddDrawable(bar)
	parent.AddDrawable(arena.bottom)
	return arena
}
