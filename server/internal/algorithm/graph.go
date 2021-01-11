package algorithm

import (
	"strconv"
	"sync"
	"fmt"

	"graph-drawing-microservices/microservices/unraveler/internal/adapters"
	"go.mongodb.org/mongo-driver/bson"
)

// Graph is a graph
type Graph struct {
	Document *adapters.GraphDocument
	Nodes    []Node
	edges    map[*Node][]*Node
	edgesStr string
	lock     sync.RWMutex
}

func (graph *Graph) InitCarbonChainGraph() {
	graph.Document = &adapters.GraphDocument{
		Ka: 0.001,
		Kr: 0.8,
		Kn: 50,
		MaxIters: 500,
		MinError: 1e-8,
	}

	length := 4

	for i := 0; i < length*3+2; i++ {
		graph.addNode()
	}

	for i := 0; i < length; i++ {
		for j := 3*i+1; j < 3*i+4; j++ {
			graph.addEdge(&graph.Nodes[i*3], &graph.Nodes[j])
			graph.addEdge(&graph.Nodes[j], &graph.Nodes[i*3])
		}
	}

	graph.addEdge(&graph.Nodes[0], &graph.Nodes[length*3+1])
	graph.addEdge(&graph.Nodes[length*3+1], &graph.Nodes[0])

	graph.edgesStr = graph.getEdges()
}

// InitGraphFromDocument initializes a Graph from a GraphDocument
func (graph *Graph) InitGraphFromDocument(gd *adapters.GraphDocument) {
	graph.Document = gd

	for i := 0; i < len(gd.E); i++ {
		graph.addNode()
	}

	for index := range graph.Nodes {
		for _, neighbor := range gd.E[strconv.Itoa(index)] {
			graph.addEdge(&graph.Nodes[index], &graph.Nodes[neighbor])
		}
	}
}

// SaveResult saves coords to mongo
func (graph *Graph) SaveResult() {
	filter := bson.D{{
		Key:   "_id",
		Value: graph.Document.ID,
	}}
	update := bson.D{{
		Key: "$set",
		Value: bson.D{{
			Key:   "coords",
			Value: graph.getAllCoords(),
		}},
	}}

	adapters.UpdateGraph(filter, update)
}

func (graph *Graph) addNode() {
	graph.lock.Lock()
	n := Node{}
	n.InitializeLocation()
	graph.Nodes = append(graph.Nodes, n)
	graph.lock.Unlock()
}

func (graph *Graph) addEdge(ni, nj *Node) {
	graph.lock.Lock()
	if graph.edges == nil {
		graph.edges = make(map[*Node][]*Node)
	}

	graph.edges[ni] = append(graph.edges[ni], nj)
	// graph.edges[nj] = append(graph.edges[nj], ni)
	graph.lock.Unlock()
}

func (graph *Graph) getAllCoords() [][]float64 {
	var allCoords [][]float64

	for n := range graph.Nodes {
		allCoords = append(allCoords, graph.Nodes[n].getCoords())
	}

	return allCoords
}

func (graph *Graph) getAllCoordsStr() string {
	coords := "["

	for n := range graph.Nodes {
		c := graph.Nodes[n].getCoords()
		coords = coords + "[" + strconv.FormatFloat(c[0], 'f', 5, 64) + "," + strconv.FormatFloat(c[1], 'f', 5, 64) + "]"
		if n != len(graph.Nodes)-1 {
			coords = coords + ","
		}
	}

	coords += "]"

	return coords
}

func (graph *Graph) getEdges() string {
	m := make(map[*Node]int)

	for i := range graph.Nodes {
		m[&graph.Nodes[i]] = i
	}

	edges := "["

	for i := range graph.Nodes {
		e := graph.edges[&(graph.Nodes[i])]

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

		if i != len(graph.Nodes)-1 {
			edges += ","
		}
	}



	edges += "]" //edges[:len(edges)-1]

	return edges
}

func (graph *Graph) toString(additional string) string {
	return "{\"edges\":" + graph.getEdges() + ", \"nodes\":" + graph.getAllCoordsStr() + "," + additional + "}"
}