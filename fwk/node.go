package fwk

import (
	gl "github.com/chsc/gogl/gl21"
	"math"
)

type Node struct {
	Parent        *Node
	name          string
	position      Vector
	rotationAxe   Vector
	rotationAngle gl.Float
	children      []NamedDrawable
}

func (n *Node) Draw(ratio float64) {
	gl.PushMatrix()
	gl.Translatef(n.position.X, n.position.Y, n.position.Z)
	if n.rotationAngle != 0 {
		gl.Rotatef(n.rotationAngle, n.rotationAxe.X, n.rotationAxe.Y, n.rotationAxe.Z)
	}
	for _, c := range n.children {
		c.Draw(ratio)
	}
	gl.PopMatrix()
}

func (n *Node) Move(trans *Vector) {
	n.position.Add(trans)
}

func (n *Node) Rotate(deg float32, up Vector) {
	println("setting rotation angle to:", deg*(math.Pi/180.), " : ", math.Pi)
	n.rotationAngle = gl.Float(deg)
	n.rotationAxe = up
}

func (n *Node) GetName() string {
	return n.name
}

func (n *Node) AddDrawable(d NamedDrawable) (rs *Node) {
	n.children = append(n.children, d)
	return n
}

func NewNode(parent *Node, name string, position Vector) *Node {
	node := &Node{Parent: parent, name: name, position: position, rotationAngle: 0}
	if parent != nil {
		parent.AddDrawable(node)
	}
	return node
}
