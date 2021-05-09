package internal

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const KA = 0.0001

type Graph struct {
	Document   *GraphDocument
	nodes      []*node
	edges      map[*node][]*node
	edgesStr   string
	minX       float64
	maxX       float64
	minY       float64
	maxY       float64
	boundsLock sync.Mutex
}

func InitPreferentialAttachment(size int) *Graph {
	graph := Graph{
		minX: math.MaxFloat64,
		maxX: -math.MaxFloat64,
		minY: math.MaxFloat64,
		maxY: -math.MaxFloat64,
	}

	rand.Seed(time.Now().UnixNano())

	graph.Document = &GraphDocument{
		Ka:       0.001,
		Kr:       0.8,
		Kn:       35,
		MaxIters: 3000,
		MinError: 1e-8,
	}

	graph.addNode(size)
	graph.addNode(size)
	graph.addNode(size)

	graph.addEdge(graph.nodes[1], graph.nodes[0])
	graph.addEdge(graph.nodes[0], graph.nodes[1])

	graph.addEdge(graph.nodes[1], graph.nodes[2])
	graph.addEdge(graph.nodes[2], graph.nodes[1])

	var sum int
	var toAttach int

	for i := 3; i < size; i++ {
		sum = 0
		toAttach = rand.Intn(2*(len(graph.nodes)-1)) + 1

		for n := range graph.nodes {
			sum += graph.nodes[n].numEdges

			if sum >= int(toAttach) {
				graph.addNode(size)
				graph.addEdge(graph.nodes[len(graph.nodes)-1], graph.nodes[n])
				graph.addEdge(graph.nodes[n], graph.nodes[len(graph.nodes)-1])
				break
			}
		}
	}

	return &graph
}

func InitCarbonChain(size int) *Graph {
	graph := Graph{}

	graph.Document = &GraphDocument{
		Ka:       0.001,
		Kr:       0.8,
		Kn:       50,
		MaxIters: 2000,
		MinError: 1e-8,
	}

	for i := 0; i < size*3+2; i++ {
		graph.addNode(size)
	}

	for i := 0; i < size; i++ {
		for j := 3*i + 1; j < 3*i+4; j++ {
			graph.addEdge(graph.nodes[i*3], graph.nodes[j])
			graph.addEdge(graph.nodes[j], graph.nodes[i*3])
		}
	}

	graph.addEdge(graph.nodes[0], graph.nodes[size*3+1])
	graph.addEdge(graph.nodes[size*3+1], graph.nodes[0])

	graph.edgesStr = graph.getEdges()

	return &graph
}

func (graph *Graph) addNode(area int) {
	n := &node{}
	n.InitializeLocation(area)

	graph.boundsLock.Lock()
	if n.x <= graph.minX {
		graph.minX = n.x
	}
	if n.x >= graph.maxX {
		graph.maxX = n.x
	}
	if n.y <= graph.minY {
		graph.minY = n.y
	}
	if n.y >= graph.maxY {
		graph.maxY = n.y
	}
	graph.boundsLock.Unlock()

	graph.nodes = append(graph.nodes, n)
}

func (graph *Graph) addEdge(ni, nj *node) {
	if graph.edges == nil {
		graph.edges = make(map[*node][]*node)
	}

	graph.edges[ni] = append(graph.edges[ni], nj)

	ni.numEdges++
}

func (graph *Graph) getAllCoords() [][]float64 {
	var allCoords [][]float64

	for n := range graph.nodes {
		allCoords = append(allCoords, graph.nodes[n].getCoords())
	}

	return allCoords
}

func (graph *Graph) getAllCoordsStr() string {
	coords := "["

	for n := range graph.nodes {
		c := graph.nodes[n].getCoords()
		coords = coords + "[" + strconv.FormatFloat(c[0], 'f', 5, 64) + "," + strconv.FormatFloat(c[1], 'f', 5, 64) + "]"
		if n != len(graph.nodes)-1 {
			coords = coords + ","
		}
	}

	coords += "]"

	return coords
}

func (graph *Graph) getEdges() string {
	// graph.Nodes in here contains wrong objects
	m := make(map[*node]int)

	for i := range graph.nodes {
		m[graph.nodes[i]] = i
	}

	edges := "["

	for i := range graph.nodes {
		e := graph.edges[graph.nodes[i]]

		edges += "["
		if len(e) > 0 {
			for n := range e {
				edges += fmt.Sprint(m[e[n]])

				if n != len(e)-1 {
					edges += ","
				}
			}
		}
		edges += "]"

		if i != len(graph.nodes)-1 {
			edges += ","
		}
	}

	edges += "]" //edges[:len(edges)-1]

	// fmt.Println("AFTER", graph.edges)
	return edges
}

func (graph *Graph) toString(additional string) string {
	str := "{\"edges\":" + graph.getEdges() + ", \"nodes\":" + graph.getAllCoordsStr() + "," + additional + "}"
	return str
}

func (graph *Graph) computeAttraction(n *node) {
	theta := func(n1, n2 *node) float64 {
		return math.Atan2(n2.y-n1.y, n2.x-n1.x)
	}
	for _, neighbor := range graph.edges[n] {
		FijA := KA * n.distance(neighbor) // Ka = 0.001
		thetaij := theta(n, neighbor)

		n.Fx += math.Cos(thetaij) * FijA
		n.Fy += math.Sin(thetaij) * FijA
	}
}
