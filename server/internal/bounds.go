package internal

type bounds struct {
	minX float64
	maxX float64
	minY float64
	maxY float64
}

func (b *bounds) update(n *node) {
	if n.x < b.minX {
		b.minX = n.x
	}

	if n.x > b.maxX {
		b.maxX = n.x
	}

	if n.y < b.minY {
		b.minY = n.y
	}

	if n.y > b.maxY {
		b.maxY = n.y
	}
}
