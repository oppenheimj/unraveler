package internal

import (
	"fmt"
	"math"
	"sync"
)

type treeNode struct {
	children  []*treeNode
	x         float64
	y         float64
	width     float64
	height    float64
	numPoints int
	comX      float64
	comY      float64
	mass      float64
	point     *point
}

const theta = 0.9
const KR = 1

func ConstructQuadtreeFromGraph(graph *Graph) treeNode {
	root := treeNode{
		x:      graph.minX,
		y:      graph.minY,
		width:  graph.maxX - graph.minX,
		height: graph.maxY - graph.minY,
	}

	for i, node := range graph.nodes {
		p := point{x: node.x, y: node.y, node: graph.nodes[i]}

		if root.contains(&p) {
			root.addPoint(&p)
		}
	}

	return root
}

func (t *treeNode) computeRepulsion(n *node) {
	if t.isLeaf() {
		if t.hasPoint() && t.point.node != n {
			t.addRepulsion(n)
		}
	} else {
		threshold := t.width / math.Sqrt(math.Pow(t.comX-n.x, 2)+math.Pow(t.comY-n.y, 2))
		if threshold < theta {
			t.addRepulsion(n)
		} else {
			for _, child := range t.children {
				child.computeRepulsion(n)
			}
		}
	}
}

var g = 1.0

func (t *treeNode) addRepulsion(n *node) {
	thetaIJ := math.Atan2(t.comY-n.y, t.comX-n.x) + math.Pi

	dx := t.comX - n.x
	dy := t.comY - n.y
	distance := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	F := (g * 1 * t.mass * KR) / math.Pow(distance, 2) // 1 = node's mass

	// fmt.Println("Adding repulsion from treeNode", t.x, " ", t.y, "to", n.x, n.y, " force ", F*math.Cos(thetaIJ), " ", F*math.Sin(thetaIJ))
	n.Fx += F * math.Cos(thetaIJ) // dx / distance
	n.Fy += F * math.Sin(thetaIJ) // dy / distance
}

func (graph *Graph) calculateForces(s, f int, wg *sync.WaitGroup) {
	theta := func(n1, n2 *node) float64 {
		return math.Atan2(n2.y-n1.y, n2.x-n1.x)
	}

	for i := s; i < f; i++ {
		// ni is our node
		// nj is treeNode
		ni := graph.nodes[i]

		for j := i + 1; j < len(graph.nodes); j++ {
			nj := graph.nodes[j]

			dij2 := ni.distanceSquared(nj)
			FijR := graph.Document.Kr / dij2
			thetaIJ := theta(ni, nj) + math.Pi
			thetaJI := theta(nj, ni) + math.Pi

			ni.lock.Lock()
			ni.Fx += math.Cos(thetaIJ) * FijR
			ni.Fy += math.Sin(thetaIJ) * FijR
			ni.lock.Unlock()

			nj.lock.Lock()
			nj.Fx += math.Cos(thetaJI) * FijR
			nj.Fy += math.Sin(thetaJI) * FijR
			nj.lock.Unlock()
		}

		for _, neighbor := range graph.edges[ni] {
			FijA := graph.Document.Ka * ni.distance(neighbor)
			thetaij := theta(ni, neighbor)

			ni.lock.Lock()
			ni.Fx += math.Cos(thetaij) * FijA
			ni.Fy += math.Sin(thetaij) * FijA
			ni.lock.Unlock()
		}
	}

	wg.Done()
}

func (graph *Graph) updateNodes() float64 {
	var totalChange float64

	bounds := bounds{
		minX: math.MaxFloat64,
		maxX: -math.MaxFloat64,
		minY: math.MaxFloat64,
		maxY: -math.MaxFloat64,
	}

	for i := range graph.nodes {
		dx := graph.Document.Kn * graph.nodes[i].Fx
		dy := graph.Document.Kn * graph.nodes[i].Fy

		// fmt.Println("node", i, "dx", dx, "dy", dy)

		graph.nodes[i].x += dx
		graph.nodes[i].y += dy

		graph.nodes[i].Fx = 0
		graph.nodes[i].Fy = 0

		totalChange += math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))

		bounds.update(graph.nodes[i])
	}

	graph.minX = bounds.minX
	graph.maxX = bounds.maxX
	graph.minY = bounds.minY
	graph.maxY = bounds.maxY

	return totalChange / float64(len(graph.nodes))
}

func (treeNode *treeNode) toString() string {
	return fmt.Sprintf("[(%.2f, %.2f), w: %.2f, h: %.2f]", treeNode.x, treeNode.y, treeNode.width, treeNode.height)
}

func (treeNode *treeNode) contains(p *point) bool {
	inboundsX := treeNode.x <= p.x && p.x <= treeNode.x+treeNode.width
	inboundsY := treeNode.y <= p.y && p.y <= treeNode.y+treeNode.height

	return inboundsX && inboundsY
}

func (t *treeNode) addPoint(p *point) {
	if t.isLeaf() {
		if t.hasPoint() {
			t.convertToInternal()
			t.addPointToChild(p)
		} else {
			t.point = p
			t.comX = p.x
			t.comY = p.y
			t.mass = 1
			t.numPoints++
		}
	} else {
		t.addPointToChild(p)
	}
}

func (t *treeNode) isLeaf() bool {
	return len(t.children) == 0
}

func (t *treeNode) hasPoint() bool {
	return t.point != nil
}

func (t *treeNode) convertToInternal() {
	childLocations := t.childLocations()

	for i := 0; i < 4; i++ {
		t.children = append(t.children, &treeNode{
			x:      childLocations[i].x,
			y:      childLocations[i].y,
			width:  t.width / 2,
			height: t.height / 2,
		})
	}

	tmpPoint := t.point
	t.point = nil
	t.addPoint(tmpPoint)
}

func (t *treeNode) addPointToChild(p *point) {
	newComX := t.comX*float64(t.numPoints) + p.x
	newComY := t.comY*float64(t.numPoints) + p.y
	t.numPoints++
	t.comX = newComX / float64(t.numPoints)
	t.comY = newComY / float64(t.numPoints)

	// placeholder for node's mass
	t.mass++

	i := t.locationToChildrenIndex(p)
	t.children[i].addPoint(p)
}

// [nw, ne, sw, se]
func (treeNode *treeNode) childLocations() []location {
	return []location{
		location{x: treeNode.x, y: treeNode.y + treeNode.height/2},
		location{x: treeNode.x + treeNode.width/2, y: treeNode.y + treeNode.height/2},
		location{x: treeNode.x, y: treeNode.y},
		location{x: treeNode.x + treeNode.width/2, y: treeNode.y},
	}
}

// Internal nodes have an array called children which has four elements
// (since this is a quadtree). This function returns the index of the
// child node that would contain the provided location.
func (t *treeNode) locationToChildrenIndex(p *point) int {
	if p.x > t.x+t.width/2 {
		if p.y > t.y+t.height/2 {
			return 1
		}
		return 3
	}

	if p.y > t.y+t.height/2 {
		return 0
	}

	return 2
}
