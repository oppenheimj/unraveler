package main

import (
	"flag"
	"graph-drawing-microservices/microservices/unraveler/internal/algorithm"
	"log"
	"net/http"
	"runtime"
	"sync"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		// https://github.com/gorilla/websocket/issues/367
		return true
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		runtime.GOMAXPROCS(64)
		g := algorithm.Graph{}
		// g.InitCarbonChainGraph(20)
		g.InitPreferentialAttachment(20)
		g.Unravel(&sync.WaitGroup{}, mt, c)

		// err = c.WriteMessage(mt, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
