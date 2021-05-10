package internal

import (
	"fmt"
)

// FailOnError reduces boilerplate
func FailOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// func GenerateRandomPoints(n int, width float64) []point {
// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	randomLocation := func() location {
// 		return location{
// 			x: r.Float64() * width,
// 			y: r.Float64() * width,
// 		}
// 	}

// 	points := []point{}

// 	for i := 0; i < n; i++ {
// 		points = append(points, point{

// 			location: randomLocation(),
// 			mass:     1.0,
// 		})
// 	}

// 	return points
// }
