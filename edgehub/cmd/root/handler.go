package root

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

func SendMsg(w http.ResponseWriter, r *http.Request) {
	ch := make(chan byte)
	var wg sync.WaitGroup
	wg.Add(2)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panic("Upgrade ", err)
		return
	}
	var sendFlag byte = 1
	go func() {
		defer wg.Done()
		ch <- sendFlag
	}()
	defer c.Close()
	go func() {
		for {
			<-ch
			log.Print("A:")
			SendTest(c)
			ch <- sendFlag

		}
	}()
	go func() {
		for {
			<-ch
			log.Print("B")
			SendTest(c)
			ch <- sendFlag
		}
	}()
	wg.Wait()

}

func SendTest(c *websocket.Conn) {
	b := make([]byte, 0)
	b = append(b, 'l')
	err := c.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		log.Println("SendTest", err)
	}
}
