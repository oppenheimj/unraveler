package algorithm

import (
	"fmt"
	"math"
	"sync"
	"github.com/gorilla/websocket"
)

// Unravel runs the force-direction algorithm
func (graph *Graph) Unravel(wg *sync.WaitGroup, mt int, c *websocket.Conn) {
	var avgChange float64 = 1
	t := 0
	blocks := calculateBlocks(1, len(graph.Nodes))

	for t < graph.Document.MaxIters && avgChange > graph.Document.MinError {
		for i := 1; i < len(blocks); i++ {
			wg.Add(1)
			go graph.calculateForces(blocks[i-1], blocks[i], wg)
		}
		wg.Wait()

		avgChange = graph.updateNodes()

		if t % 1 == 0 {
			(*c).WriteMessage(mt, []byte(graph.toString(fmt.Sprint("\"i\":", t, ", \"err\":", avgChange))))
		}
		
		t++
	}

	fmt.Println("Done!", avgChange, t)
}

func calculateBlocks(numThreads int, numNodes int) []int {
	indices := make([]int, numThreads+1)
	npt := math.Pow(float64(numNodes), 2) / float64(2*numThreads)

	for i := numThreads - 1; i > 0; i-- {
		indices[numThreads-i] = numNodes - int(math.Sqrt(float64(i)*npt*2))
	}

	indices[numThreads] = numNodes

	return indices
}

func (graph *Graph) calculateForces(s, f int, wg *sync.WaitGroup) {
	theta := func(n1, n2 *Node) float64 {
		return math.Atan2(n2.Y-n1.Y, n2.X-n1.X)
	}

	for i := s; i < f; i++ {
		ni := &graph.Nodes[i]

		for j := i + 1; j < len(graph.Nodes); j++ {
			nj := &graph.Nodes[j]

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

	for i := range graph.Nodes {
		dx := graph.Document.Kn * graph.Nodes[i].Fx
		dy := graph.Document.Kn * graph.Nodes[i].Fy

		// fmt.Println("node", i, "dx", dx, "dy", dy)

		graph.Nodes[i].X += dx
		graph.Nodes[i].Y += dy

		graph.Nodes[i].Fx = 0
		graph.Nodes[i].Fy = 0

		totalChange += math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	}

	return totalChange / float64(len(graph.Nodes))
}
