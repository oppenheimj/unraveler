package internal

import (
	"fmt"
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
