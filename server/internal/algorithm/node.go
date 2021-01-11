package algorithm

import (
	"math"
	"math/rand"
	"sync"
)

// Node is a node
type Node struct {
	X    float64
	Y    float64
	Path [][]float64
	Fx   float64
	Fy   float64
	lock sync.RWMutex
}

// InitializeLocation does what it sounds like
func (node *Node) InitializeLocation() {
	node.X = rand.Float64() * 20
	node.Y = rand.Float64() * 20
}

func (node *Node) getCoords() []float64 {
	return []float64{node.X, node.Y}
}

func (node *Node) distance(otherNode *Node) float64 {
	return math.Sqrt(node.distanceSquared(otherNode))
}

func (node *Node) distanceSquared(otherNode *Node) float64 {
	return math.Pow(otherNode.X-node.X, 2) + math.Pow(otherNode.Y-node.Y, 2)
}
