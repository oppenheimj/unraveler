package internal

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

func (graph *Graph) Unravel(wg *sync.WaitGroup, mt int, c *websocket.Conn) {
	for t := 0; t < 10000; t++ {
		fmt.Println("building quadtree")
		q := ConstructQuadtreeFromGraph(graph)
		fmt.Println("built quadtree")

		for i := range graph.nodes {
			q.computeRepulsion(graph.nodes[i])
			graph.computeAttraction(graph.nodes[i])
		}
		fmt.Println("computed forces")

		avgChange := graph.updateNodes()

		fmt.Println("updated nodes")
		if t%10 == 0 {
			fmt.Println("writing message")
			(*c).WriteMessage(mt, []byte(graph.toString(fmt.Sprint("\"i\":", t, ", \"err\":", avgChange))))
			fmt.Println("wrote message")
		}

	}
}

func (graph *Graph) UnravelOld(wg *sync.WaitGroup, mt int, c *websocket.Conn) {
	fmt.Println("Unraveleing")
	var avgChange float64 = 1
	t := 0
	// blocks := calculateBlocks(4, len(graph.nodes))

	for t < graph.params.maxIters && avgChange > graph.params.minError {
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
