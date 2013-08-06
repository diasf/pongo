package fwk

import (
)

type SimpleCollision struct {
	items []Collidable
	handlers []CollisionHandler
}

type CollisionHandler interface {
	HandleCollision(one, two Collidable)
}

type Collidable interface {
	CollidesWith(other Collidable) bool
}

func (n * SimpleCollision) Check() {
	li := len(n.items)
	for i:=0; i < li; i++ {
		for j:=i+1; j < li; j++ {
			if n.items[i].CollidesWith(n.items[j]) {
				n.notifyCollision(n.items[i], n.items[j])
			}
		}
	}
}

func (n * SimpleCollision) notifyCollision(one, two Collidable) {
	for _,h := range n.handlers {
		h.HandleCollision(one, two)
	}
}
