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
	c.handlers = make([]CollisionHandler, 5)
	c.items = make([]Collidable, 5)
	return c
}

func (n *SimpleCollisionDetector) Check() {
	li := len(n.items)
	for i := 0; i < li; i++ {
		for j := i + 1; j < li; j++ {
			if n.items[i].GetBoundingVolume().CollidesWith(n.items[j].GetBoundingVolume()) {
				n.notifyCollision(n.items[i], n.items[j])
			}
		}
	}
}

func (n *SimpleCollisionDetector) notifyCollision(one, two Collidable) {
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
	HandleCollision(one, two Collidable)
}

type Collidable interface {
	GetBoundingVolume() BoundingVolume
}

// Bounding Box -----------------------------------------------

type BoundingBox struct {
	left, right, top, bottom gl.Float
}

func (v *BoundingBox) CollidesWith(other BoundingVolume) bool {
	if obx, ok := other.(*BoundingBox); ok {
		if (v.left <= obx.right && v.right >= obx.left) || (v.right >= obx.left && v.left <= obx.right) {
			if (v.top >= obx.bottom && v.bottom <= obx.top) || (v.bottom <= obx.top && v.top >= obx.bottom) {
				return true
			}
		}
	} else {
		if v.inVolume(other.GetNearestTo(&Vector{v.left, v.bottom, 0})) || v.inVolume(other.GetNearestTo(&Vector{v.left, v.top, 0})) || v.inVolume(other.GetNearestTo(&Vector{v.left, v.top, 0})) || v.inVolume(other.GetNearestTo(&Vector{v.left, v.top, 0})) {
			return true
		}
	}
	return false
}

func (v *BoundingBox) GetNearestTo(p *Vector) *Vector {
	return nil
}

func (v *BoundingBox) inVolume(p *Vector) bool {
	return (v.left <= p.X && v.right >= p.X) && (v.bottom <= p.Y && v.top >= p.Y)
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
