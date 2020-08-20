package root

import (
	"flag"
	"fmt"
	websocket "github.com/gorilla/websocket"
	"log"
	"net/http"
)



var addr = flag.String("port", ":43211", "http service address")

// use default options
var upgrader = websocket.Upgrader{}
var pin = websocket.PingMessage
var pong = websocket.PongMessage

func Run() {
	hub := NewHub()
	go hub.Run()
	go Serve(hub)
	flag.Parse()
	http.HandleFunc("/hub", func(w http.ResponseWriter, r *http.Request) {
		Servews(hub, w, r)
	})
	fmt.Println("listening", *addr)
	defer log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", GetConfig().Socket), nil))

}

func Ping(conn *websocket.Conn) {
	err := conn.WriteMessage(pin, nil)
	log.Println("ping")
	if err != nil {
		log.Println("ping: ", err)
	}
}
