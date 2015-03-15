package game

import "github.com/diasf/pongo/fwk"

// the arena has two bars
type Arena struct {
	top           *fwk.Node
	bottom        *fwk.Node
	volumes       []fwk.BoundingVolume
	name          string
	ringsize      float32
	halfring      float32
	wallwidth     float32
	halfwallwidth float32
}

func NewArena(parent *fwk.Node, name string, ringsize float32, wallwidth float32, color fwk.Color) *Arena {
	arena := &Arena{}
	arena.name = name
	arena.ringsize = float32(ringsize)
	arena.halfring = arena.ringsize / 2.
	arena.wallwidth = float32(wallwidth)
	arena.halfwallwidth = arena.wallwidth / 2.

	// top
	topPos := fwk.Vector{0., arena.halfring - arena.halfwallwidth, 0.}
	barTop := &fwk.Rectangle{Width: float32(ringsize), Height: float32(wallwidth), Color: color, Name: "BorderRectTop"}
	arena.top = fwk.NewNode(parent, "TopBorderN", topPos).AddDrawable(barTop)
	topBox := &fwk.BoundingBox{Left: topPos.X - arena.halfring, Right: topPos.X + arena.halfring, Top: topPos.Y + arena.halfwallwidth, Bottom: topPos.Y - arena.halfwallwidth}

	// bottom
	bottomPos := fwk.Vector{0., -arena.halfring + arena.halfwallwidth, 0.}
	barBottom := &fwk.Rectangle{Width: float32(ringsize), Height: float32(wallwidth), Color: color, Name: "BorderRectBottom"}
	arena.bottom = fwk.NewNode(parent, "BottomBorderN", bottomPos).AddDrawable(barBottom)
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

func (a *Arena) GetTopBoundingVolume() fwk.BoundingVolume {
	return a.volumes[0]
}

func (a *Arena) GetBottomBoundingVolume() fwk.BoundingVolume {
	return a.volumes[1]
}
