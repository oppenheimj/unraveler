package quadtree

import "fmt"

type bounds struct {
	minX float64
	maxX float64
	minY float64
	maxY float64
}

func (b *bounds) update(l location) {
	if l.x < b.minX {
		b.minX = l.x
	}

	if l.x > b.maxX {
		b.maxX = l.x
	}

	if l.y < b.minY {
		b.minY = l.y
	}

	if l.y > b.maxY {
		b.maxY = l.y
	}
}

func findBounds(points []Point) bounds {
	bounds := bounds{
		minX: points[0].location.x,
		maxX: points[0].location.x,
		minY: points[0].location.y,
		maxY: points[0].location.y,
	}

	for _, point := range points {
		bounds.update(point.location)
	}

	fmt.Printf("Bounds are (%.2f, %.2f) and (%.2f, %.2f)\n", bounds.minX, bounds.minY, bounds.maxX, bounds.maxY)

	return bounds
}
