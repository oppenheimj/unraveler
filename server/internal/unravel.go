package internal

import (
	"fmt"
	"math"

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

		j.q.computeRepulsion(j.n, j.g.params.Kr, j.g.params.Theta)
		j.g.computeAttraction(j.n)

		// fmt.Println("worker", id, "finished job", j)
		results <- 0
	}
}

func (graph *Graph) Unravel(mt int, c *websocket.Conn) {
	t := 0
	avgChange := 500.0

	for t < graph.params.MaxIters && avgChange > graph.params.MinError {
		q := ConstructQuadtreeFromGraph(graph)

		numJobs := len(graph.nodes)
		numWorkers := graph.params.NumThreads

		jobs := make(chan nodeWorkerData, numJobs)
		results := make(chan int, numJobs)

		for w := 0; w < numWorkers; w++ {
			go nodeWorker(w, jobs, results)
		}

		for i := range graph.nodes {
			jobs <- nodeWorkerData{
				n: graph.nodes[i],
				q: &q,
				g: graph,
			}
		}

		close(jobs)

		for a := 1; a <= numJobs; a++ {
			<-results
		}

		avgChange = graph.updateNodes()

		if t%graph.params.UpdateEvery == 0 {
			(*c).WriteMessage(mt, []byte(graph.toString(fmt.Sprint("\"i\":", t, ", \"err\":", avgChange, ", \"minX\":", graph.minX, ", \"maxX\":", graph.maxX, ", \"minY\":", graph.minY, ", \"maxY\":", graph.maxY))))
		}

		t++

	}
}

func ConstructQuadtreeFromGraph(graph *Graph) treeNode {
	root := treeNode{
		x:      graph.minX,
		y:      graph.minY,
		width:  graph.maxX - graph.minX,
		height: graph.maxY - graph.minY,
	}

	for i, node := range graph.nodes {
		p := point{x: node.x, y: node.y, node: graph.nodes[i]}

		if root.contains(&p) {
			root.addPoint(&p)
		}
	}

	return root
}

func (t *treeNode) computeRepulsion(n *node, kr float64, theta float64) {
	if t.isLeaf() {
		if t.hasPoint() && t.point.node != n {
			t.addRepulsion(n, kr)
		}
	} else {
		threshold := t.width / math.Sqrt(math.Pow(t.comX-n.x, 2)+math.Pow(t.comY-n.y, 2))
		if threshold < theta {
			t.addRepulsion(n, kr)
		} else {
			for _, child := range t.children {
				child.computeRepulsion(n, kr, theta)
			}
		}
	}
}

var g = 1.0

func (t *treeNode) addRepulsion(n *node, kr float64) {
	thetaIJ := math.Atan2(t.comY-n.y, t.comX-n.x) + math.Pi
	distSquared := math.Pow(t.comX-n.x, 2) + math.Pow(t.comY-n.y, 2)

	force := (g * 1 * t.mass * kr) / distSquared // (g * float64(n.numEdges) * t.mass * kr) / distSquared

	n.Fx += force * math.Cos(thetaIJ)
	n.Fy += force * math.Sin(thetaIJ)
}

func (graph *Graph) computeAttraction(n *node) {
	theta := func(n1, n2 *node) float64 {
		return math.Atan2(n2.y-n1.y, n2.x-n1.x)
	}
	for _, neighbor := range graph.edges[n] {
		FijA := graph.params.Ka * n.distance(neighbor) //* math.Log2(float64(n.numEdges)*float64(neighbor.numEdges))
		thetaij := theta(n, neighbor)

		n.Fx += math.Cos(thetaij) * FijA
		n.Fy += math.Sin(thetaij) * FijA
	}
}

func (graph *Graph) updateNodes() float64 {
	var totalChange float64

	bounds := bounds{
		minX: math.MaxFloat64,
		maxX: -math.MaxFloat64,
		minY: math.MaxFloat64,
		maxY: -math.MaxFloat64,
	}

	for i := range graph.nodes {
		dx := graph.params.Kn * graph.nodes[i].Fx
		dy := graph.params.Kn * graph.nodes[i].Fy

		graph.nodes[i].x += math.Copysign(1, dx) * math.Log2(math.Abs(dx)+1)
		graph.nodes[i].y += math.Copysign(1, dy) * math.Log2(math.Abs(dy)+1)

		graph.nodes[i].Fx = 0
		graph.nodes[i].Fy = 0

		totalChange += math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))

		bounds.update(graph.nodes[i])
	}

	graph.minX = bounds.minX
	graph.maxX = bounds.maxX
	graph.minY = bounds.minY
	graph.maxY = bounds.maxY

	return totalChange / float64(len(graph.nodes))
}
