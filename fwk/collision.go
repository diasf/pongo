package fwk

import (
	gl "github.com/chsc/gogl/gl21"
)

type CollisionDetector interface {
	Check()
	AddCollidable(c Collidable)
	AddCollisionHandler(h CollisionHandler)
}

type SimpleCollisionDetector struct {
	items    []Collidable
	handlers []CollisionHandler
}

func NewSimpleCollisionDetector() *SimpleCollisionDetector {
	c := &SimpleCollisionDetector{}
	c.handlers = make([]CollisionHandler, 0)
	c.items = make([]Collidable, 0)
	return c
}

func (n *SimpleCollisionDetector) Check() {
	li := len(n.items)
	for i, itemL := range n.items {
		for j := i + 1; j < li; j++ {
			for _, bv := range itemL.GetBoundingVolumes() {
				for _, rbv := range n.items[j].GetBoundingVolumes() {
					if bv.CollidesWith(rbv) {
						n.notifyCollision(newCollisionObject(itemL, bv), newCollisionObject(n.items[j], rbv))
					}
				}
			}
		}
	}
}

func (n *SimpleCollisionDetector) notifyCollision(one, two CollisionObject) {
	for _, h := range n.handlers {
		h.HandleCollision(one, two)
	}
}

func (n *SimpleCollisionDetector) AddCollisionHandler(h CollisionHandler) {
	n.handlers = append(n.handlers, h)
}

func (n *SimpleCollisionDetector) AddCollidable(c Collidable) {
	n.items = append(n.items, c)
}

type BoundingVolume interface {
	CollidesWith(other BoundingVolume) bool
	GetNearestTo(p *Vector) *Vector
}

type CollisionHandler interface {
	HandleCollision(one, two CollisionObject)
}

type Collidable interface {
	Named
	GetBoundingVolumes() []BoundingVolume
}

type CollisionObject struct {
	object         Collidable
	boundingVolume BoundingVolume
}

func (c CollisionObject) GetObject() Collidable {
	return c.object
}

func (c CollisionObject) GetBoundingVolume() BoundingVolume {
	return c.boundingVolume
}

func newCollisionObject(obj Collidable, bv BoundingVolume) CollisionObject {
	return CollisionObject{obj, bv}
}

// Bounding Box -----------------------------------------------

type BoundingBox struct {
	Left, Right, Top, Bottom gl.Float
}

func (v *BoundingBox) CollidesWith(other BoundingVolume) bool {
	if obx, ok := other.(*BoundingBox); ok {
		if (v.Left <= obx.Right && v.Right >= obx.Left) || (v.Right >= obx.Left && v.Left <= obx.Right) {
			if (v.Top >= obx.Bottom && v.Bottom <= obx.Top) || (v.Bottom <= obx.Top && v.Top >= obx.Bottom) {
				return true
			}
		}
	} else {
		if v.inVolume(other.GetNearestTo(&Vector{v.Left, v.Bottom, 0})) ||
			v.inVolume(other.GetNearestTo(&Vector{v.Left, v.Top, 0})) ||
			v.inVolume(other.GetNearestTo(&Vector{v.Right, v.Bottom, 0})) ||
			v.inVolume(other.GetNearestTo(&Vector{v.Right, v.Top, 0})) {
			return true
		}
	}
	return false
}

func (v *BoundingBox) GetNearestTo(p *Vector) *Vector {
	min := Vector{v.Left, v.Bottom, 0}
	max := Vector{v.Right, v.Top, 0}
	nearest := p.Clamp(min, max)
	return &nearest
}

func (v *BoundingBox) inVolume(p *Vector) bool {
	return (v.Left <= p.X && v.Right >= p.X) && (v.Bottom <= p.Y && v.Top >= p.Y)
}

// Bounding Sphere -----------------------------------------------

type BoundingSphere struct {
	radius   float64
	position *Vector
}

func (v *BoundingSphere) CollidesWith(other BoundingVolume) (collides bool) {
	if osp, ok := other.(*BoundingSphere); ok {
		pv := *(v.position)
		pv.Substract(osp.position)
		radii := v.radius + osp.radius
		collides = pv.Dot(&pv) <= (radii * radii)
	} else {
		nv := other.GetNearestTo(v.position)
		nv.Substract(v.position)
		collides = nv.Length() <= v.radius
	}
	return
}

func (v *BoundingSphere) GetNearestTo(p *Vector) *Vector {
	return v.position.Substraction(p).Normalize().Mult(gl.Float(v.radius))
}
