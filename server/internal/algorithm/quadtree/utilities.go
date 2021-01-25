package quadtree

import (
	"math/rand"
	"time"
)

func GenerateRandomPoints(n int, width float64) []Point {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomLocation := func() location {
		return location{
			x: r.Float64() * width,
			y: r.Float64() * width,
		}
	}

	points := []Point{}

	for i := 0; i < n; i++ {
		points = append(points, Point{
			location: randomLocation(),
			mass: 1.0,
		})
	}

	return points
}
