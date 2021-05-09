package internal

import (
	"math"
	"math/rand"
	"sync"
)

type node struct {
	x        float64
	y        float64
	Path     [][]float64
	Fx       float64
	Fy       float64
	numEdges int
	lock     sync.RWMutex
}

func (n *node) InitializeLocation(area int) {
	n.x = rand.Float64() * float64(area)
	n.y = rand.Float64() * float64(area)
}

func (n *node) getCoords() []float64 {
	return []float64{n.x, n.y}
}

func (n *node) distance(otherNode *node) float64 {
	return math.Sqrt(n.distanceSquared(otherNode))
}

func (n *node) distanceSquared(otherNode *node) float64 {
	return math.Pow(otherNode.x-n.x, 2) + math.Pow(otherNode.y-n.y, 2)
}