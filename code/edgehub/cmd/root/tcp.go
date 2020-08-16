package root

import (
	"flag"
	"fmt"
	websocket "github.com/gorilla/websocket"
	"github.com/wuff1996/edgeHub/internal/protobuf"
	"log"
	"net/http"
	"sync"
	"syscall"
	"time"
)

var addr = flag.String("addr", ":43211", "http service address")

var upgrader = websocket.Upgrader{} // use default options
var pin = websocket.PingMessage
var pon = websocket.PongMessage

func echo(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	wg.Add(1)

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	c.SetPingHandler(func(appData string) error {
		log.Println("receiving ping " + appData)
		return Pong(c)
	})
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if mt == websocket.PingMessage {
			err := c.WriteMessage(websocket.PongMessage, nil)
			if err != nil {
				log.Println(err)
			}
		}
		pmsg := protobuf.ReadBuf(message)
		log.Println(pmsg.String())
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break

		}
	}
}

func init() {
	http.HandleFunc("/echo", echo)
	flag.Parse()
	fmt.Println(*addr)
	log.Fatal(http.ListenAndServe(*addr, nil))

}

func Ping(conn *websocket.Conn) {
	err := conn.WriteMessage(pin, nil)
	if err != nil {
		log.Println(err)
	}
}

func Pong(conn *websocket.Conn) error {
	err := conn.WriteMessage(pon, GetTime())
	if err != nil {
		log.Println(err)
	}
	return err
}

func GetTime() []byte {
	t := time.Now().UTC().Unix()
	s := fmt.Sprint(time.Unix(t, 0).UTC())
	b, err := syscall.ByteSliceFromString(s)
	if err != nil {
		log.Println(err)
	}
	return b
}
