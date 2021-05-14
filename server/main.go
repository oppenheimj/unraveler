package main

import (
	"encoding/json"
	"flag"
	"graph-drawing-microservices/microservices/unraveler/internal"
	"log"
	"net/http"
	"runtime"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func unpackFrontendPayload(message []byte) internal.Params {
	var params internal.Params

	if err := json.Unmarshal(message, &params); err != nil {
		panic(err)
	}

	// fmt.Printf("%+v\n", params)

	return params
}

func echo(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		// https://github.com/gorilla/websocket/issues/367
		return true
	}

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer connection.Close()

	for {
		runtime.GOMAXPROCS(64)

		messageType, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("recv: %s", message)

		g := internal.InitPreferentialAttachment(unpackFrontendPayload(message))
		// g := internal.InitCarbonChain(200)

		g.Unravel(messageType, connection)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
