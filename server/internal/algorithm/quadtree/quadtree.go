package quadtree

import (
	"fmt"
	"image"
    "image/color"
    "image/draw"
    "image/png"
    "os"
	"log"
	"math"
)

type Quadtree struct {
	root *node
	bounds bounds
}

func ConstructQuadtree(points []Point) Quadtree {
	fmt.Printf("Constructing quadtree from %d points\n", len(points))

	bounds := findBounds(points)
	quadtree := Quadtree{bounds: bounds}

	root := node{
		q: &quadtree,
		location: location{x: bounds.minX, y: bounds.minY},
		width: bounds.maxX - bounds.minX,
		height: bounds.maxY - bounds.minY,
	}

	quadtree.root = &root

	for i, point := range points {
		if root.contains(point.location) {
			root.addPoint(&points[i])
		}
	}

	return quadtree
}

// func (q *Quadtree) ComputeForces(p Point) {
// 	theta := 0.1
// 	var recurse func(n node)

// 	recurse = func(n node) {
// 		if n.isLeaf() {
// 			if n.hasPoint() && n.point != &p {
// 				p.addForce(n)
// 			}
// 		} else {
// 			threshold := n.width / math.Sqrt(math.Pow(n.centerOfMass.x-p.location.x, 2)+math.Pow(n.centerOfMass.y-p.location.y, 2))
// 			if threshold < theta {
// 				p.addForce(n)
// 			} else {
// 				for _, child := range n.children {
// 					recurse(child)
// 				}
// 			}
// 		}
// 	}

// 	recurse(*q.root)
// }

func (q *Quadtree) Draw() {
	rectangle := "rectangle.png"

    rectImage := image.NewRGBA(image.Rect(
		int(math.Floor(q.bounds.minX)),
		int(math.Floor(q.bounds.minY)),
		int(math.Floor(q.bounds.maxX)),
		int(math.Floor(q.bounds.maxY)),
	))

	black := color.RGBA{0, 0, 0, 255}
	red := color.RGBA{255, 0, 0, 255}
	// green := color.RGBA{0, 255, 0, 255}
	white := color.RGBA{255, 255, 255, 255}

	draw.Draw(rectImage, rectImage.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)

	var recurse func(n node)
	recurse = func(n node) {
		if !n.isLeaf() {
			line := image.Rect(
				int(math.Floor(n.location.x + n.width / 2)),
				int(math.Floor(n.location.y)),
				int(math.Floor(n.location.x + n.width / 2)) + 1,
				int(math.Floor(n.location.y + n.height)),
			)
			
			draw.Draw(rectImage, line.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

			line = image.Rect(
				int(math.Floor(n.location.x)),
				int(math.Floor(n.location.y + n.height / 2)),
				int(math.Floor(n.location.x + n.width)),
				int(math.Floor(n.location.y + n.height / 2)) + 1,
			)
			
			draw.Draw(rectImage, line.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

			// point := image.Rect(
			// 	int(math.Floor(n.centerOfMass.x)),
			// 	int(math.Floor(n.centerOfMass.y)),
			// 	int(math.Floor(n.centerOfMass.x)) + 1,
			// 	int(math.Floor(n.centerOfMass.y)) + 1,
			// )
			
			// draw.Draw(rectImage, point.Bounds(), &image.Uniform{green}, image.ZP, draw.Src)
			for _, child := range n.children {
				recurse(child)
			}
		} else {
			if n.hasPoint() {
				point := image.Rect(
					int(math.Floor(n.point.location.x)),
					int(math.Floor(n.point.location.y)),
					int(math.Floor(n.point.location.x)) + 1,
					int(math.Floor(n.point.location.y)) + 1,
				)
				
				draw.Draw(rectImage, point.Bounds(), &image.Uniform{red}, image.ZP, draw.Src)
			}
		}
	}

	recurse(*q.root)

    file, err := os.Create(rectangle)
	if err != nil {
		log.Fatalf("failed create file: %s", err)
	}
    png.Encode(file, rectImage)
}