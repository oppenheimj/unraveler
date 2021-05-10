package internal

import "fmt"

type point struct {
	x    float64
	y    float64
	node *node
}

func (p *point) toString() string {
	return fmt.Sprintf("(%.2f, %.2f)", p.x, p.y)
}
