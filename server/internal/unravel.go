package internal

import (
	"fmt"
	"math"
	"sync"

	"github.com/gorilla/websocket"
)

func (graph *Graph) Unravel(wg *sync.WaitGroup, mt int, c *websocket.Conn) {
	for t := 0; t < 10000; t++ {
		q := ConstructQuadtreeFromGraph(graph)

		for i := range graph.nodes {
			q.computeRepulsion(graph.nodes[i])
			graph.computeAttraction(graph.nodes[i])
		}

		avgChange := graph.updateNodes()

		if t%10 == 0 {
			(*c).WriteMessage(mt, []byte(graph.toString(fmt.Sprint("\"i\":", t, ", \"err\":", avgChange))))
		}

	}
}

func (graph *Graph) UnravelOld(wg *sync.WaitGroup, mt int, c *websocket.Conn) {
	fmt.Println("Unraveleing")
	var avgChange float64 = 1
	t := 0
	// blocks := calculateBlocks(4, len(graph.nodes))

	for t < graph.Document.MaxIters && avgChange > graph.Document.MinError {
		// construct quadtree
		q := ConstructQuadtreeFromGraph(graph)
		fmt.Println(q)

		// use threadpool to compute forces
		// update positions

		// for i := 1; i < len(blocks); i++ {
		// 	wg.Add(1)
		// 	go graph.calculateForces(blocks[i-1], blocks[i], wg)
		// }
		// wg.Wait()

		// avgChange = graph.updateNodes()

		// if t%10 == 0 {
		// 	(*c).WriteMessage(mt, []byte(graph.toString(fmt.Sprint("\"i\":", t, ", \"err\":", avgChange))))
		// }

		// TODO: Use mutex for avgChange on graph

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
