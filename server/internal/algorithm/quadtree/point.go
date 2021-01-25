package quadtree

import "fmt"

type Point struct {
	location location
	mass float64
}

func (p *Point) toString() string {
	return fmt.Sprintf("(%.2f, %.2f)", p.location.x, p.location.y)
}

func (p *Point) addForce(n node) {}
