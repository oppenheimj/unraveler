package internal

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type nodeWorkerData struct {
	n *node
	q *treeNode
	g *Graph
}

func nodeWorker(id int, jobs <-chan nodeWorkerData, results chan<- int) {
    for j := range jobs {
        // fmt.Println("worker", id, "started  job", j)

		j.q.computeRepulsion(j.n)
		j.g.computeAttraction(j.n)

        // fmt.Println("worker", id, "finished job", j)
        results <- 0
    }
}

func (graph *Graph) Unravel(mt int, c *websocket.Conn) {
	t := 0
	avgChange := 500.0

	for t < graph.params.maxIters && avgChange > graph.params.minError {
		q := ConstructQuadtreeFromGraph(graph)

		numJobs := len(graph.nodes)
		numWorkers := 8
		
		jobs := make(chan nodeWorkerData, numJobs)
		results := make(chan int, numJobs)

		for w := 1; w <= numWorkers; w++ {
			go nodeWorker(w, jobs, results)
		}

		// need worker to do this for every node
		for i := range graph.nodes {
			jobs <- nodeWorkerData{
				n : graph.nodes[i],
				q : &q,
				g : graph,
			}
		}
	
		close(jobs)

		avgChange := graph.updateNodes()

		if t%50 == 0 {
			go (*c).WriteMessage(mt, []byte(graph.toString(fmt.Sprint("\"i\":", t, ", \"err\":", avgChange))))
		}

		t++

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
