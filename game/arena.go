package game

import (
	"fmt"
	gl "github.com/chsc/gogl/gl21"
	"github.com/diasf/pongo/fwk"
)

// the arena has two bars
type Arena struct {
	top           *fwk.Node
	bottom        *fwk.Node
	volumes       []fwk.BoundingVolume
	name          string
	ringsize      gl.Float
	halfring      gl.Float
	wallwidth     gl.Float
	halfwallwidth gl.Float
}

func NewArena(parent *fwk.Node, name string, ringsize float32, wallwidth float32, color fwk.Color) *Arena {
	arena := &Arena{}
	arena.name = name
	arena.ringsize = gl.Float(ringsize)
	arena.halfring = arena.ringsize / 2.
	arena.wallwidth = gl.Float(wallwidth)
	arena.halfwallwidth = arena.wallwidth / 2.
	bar := &fwk.Rectangle{gl.Float(ringsize), gl.Float(wallwidth), color, "BorderRect"}

	// top
	topPos := fwk.Vector{0., arena.halfring - arena.halfwallwidth, 0.}
	arena.top = fwk.NewNode(parent, "TopBorderN", topPos).AddDrawable(bar)
	topBox := &fwk.BoundingBox{Left: topPos.X - arena.halfring, Right: topPos.X + arena.halfring, Top: topPos.Y + arena.halfwallwidth, Bottom: topPos.Y - arena.halfwallwidth}
	fmt.Println("Arena TopBox:", topBox)

	// bottom
	bottomPos := fwk.Vector{0., -arena.halfring + arena.halfwallwidth, 0.}
	arena.bottom = fwk.NewNode(parent, "BottomBorderN", bottomPos).AddDrawable(bar)
	bottomBox := &fwk.BoundingBox{Left: bottomPos.X - arena.halfring, Right: bottomPos.X + arena.halfring, Top: bottomPos.Y + arena.halfwallwidth, Bottom: bottomPos.Y - arena.halfwallwidth}

	arena.volumes = []fwk.BoundingVolume{topBox, bottomBox}
	return arena
}

func (a *Arena) GetBoundingVolumes() []fwk.BoundingVolume {
	return a.volumes
}

func (a *Arena) GetName() string {
	return a.name
}
