package quadtree

import "fmt"

type node struct {
	q *Quadtree
	children []node
	location location
	width    float64
	height	float64
	numPoints int
	centerOfMass location
	mass float64
	point *Point
}

func (n *node) toString() string {
	return fmt.Sprintf("[(%.2f, %.2f), w: %.2f, h: %.2f]", n.location.x, n.location.y, n.width, n.height)
}

func (n *node) contains(l location) bool {
	return n.location.x <= l.x && l.x <= n.location.x + n.width && n.location.y <= l.y && l.y <= n.location.y + n.height
}

func (n *node) addPoint(p *Point) {
	fmt.Printf("Adding point %s to node %s\n", p.toString(), n.toString())

	if n.isLeaf() {
		if n.hasPoint() {
			n.convertToInternal()
			n.addPointToChild(p)
		} else {
			n.point = p
			n.centerOfMass = p.location
			n.mass = p.mass
			n.numPoints++
		}
	} else {
		n.addPointToChild(p)
	}
}

func (n *node) isLeaf() bool {
	return len(n.children) == 0
}

func (n *node) hasPoint() bool {
	return n.point != nil
}

func (n *node) convertToInternal() {
	fmt.Printf("Converting node %s to internal\n", n.toString())
	childLocations := childLocations(n.location, n.width, n.height)

	for i := 0; i < 4; i++ {
		n.children = append(n.children, node{
			location: childLocations[i],
			width: n.width / 2,
			height: n.height / 2,
		})
	}

	pTmp := n.point
	n.point = nil
	n.addPoint(pTmp)
}

func (n *node) addPointToChild(p *Point) {
	newComX := n.centerOfMass.x * float64(n.numPoints) + p.location.x
	newComY := n.centerOfMass.y * float64(n.numPoints) + p.location.y
	n.numPoints++
	n.centerOfMass.x = newComX / float64(n.numPoints)
	n.centerOfMass.y = newComY / float64(n.numPoints)

	n.mass += p.mass

	i := n.locationToChildrenIndex(p.location)
	n.children[i].addPoint(p)
}

// [nw, ne, sw, se]
func childLocations(l location, w float64, h float64) []location {
	return []location{
		location{x: l.x, y: l.y + h/2},
		location{x: l.x + w/2, y: l.y + h/2},
		location{x: l.x, y: l.y},
		location{x: l.x + w/2, y: l.y},
	}
}

// Internal nodes have an array called children which has four elements
// (since this is a quadtree). This function returns the index of the
// child node that would contain the provided location.
func (n *node) locationToChildrenIndex(l location) int {
	if l.x > n.location.x + n.width/2 {
		if l.y > n.location.y+n.height/2 {
			return 1
		}
		return 3
	}

	if l.y > n.location.y + n.height/2 {
		return 0
	}

	return 2
}
