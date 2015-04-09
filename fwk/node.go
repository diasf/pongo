package fwk

import (
	"golang.org/x/mobile/f32"
)

type Node struct {
	Parent        *Node
	name          string
	position      Vector
	rotationAxe   Vector
	rotationAngle float32
	children      []NamedDrawable
	modelView     *f32.Mat4
}

func (n *Node) Draw(modelView *f32.Mat4, ratio float64) {
	nodeMV := &f32.Mat4{}
	nodeMV.Translate(modelView, n.position.X, n.position.Y, n.position.Z)
	if n.rotationAngle != 0 {
		nodeMV.Rotate(nodeMV, f32.Radian(n.rotationAngle), &f32.Vec3{n.rotationAxe.X, n.rotationAxe.Y, n.rotationAxe.Z})
	}
	for _, c := range n.children {
		c.Draw(nodeMV, ratio)
	}
}

func (n *Node) OnAttached() {
}

func (n *Node) Move(trans *Vector) {
	n.position.Add(trans)
}

func (n *Node) Rotate(deg float32, up Vector) {
	n.rotationAngle = float32(deg)
	n.rotationAxe = up
}

func (n *Node) GetName() string {
	return n.name
}

func (n *Node) GetPosition() *Vector {
	return &n.position
}

func (n *Node) AddDrawable(d NamedDrawable) (rs *Node) {
	n.children = append(n.children, d)
	d.OnAttached()
	return n
}

func NewNode(parent *Node, name string, position Vector) *Node {
	node := &Node{Parent: parent, name: name, position: position, rotationAngle: 0}
	if parent != nil {
		parent.AddDrawable(node)
	}
	return node
}
