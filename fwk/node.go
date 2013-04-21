package fwk

import (
	gl "github.com/chsc/gogl/gl21"
)

type Node struct {
	name     string
	Position Vector
	children []NamedDrawable
}

func (n *Node) Draw(ratio float64) {
	gl.PushMatrix()
	gl.Translatef(n.Position.X, n.Position.Y, n.Position.Z)
	for _, c := range n.children {
		c.Draw(ratio)
	}
	gl.PopMatrix()
}

func (n *Node) GetName() string {
	return n.name
}

func (n *Node) AddDrawable(d NamedDrawable) (rs *Node) {
	n.children = append(n.children, d)
	return n
}

func NewNode(name string, position Vector) *Node {
	return &Node{name: name, Position: position}
}
