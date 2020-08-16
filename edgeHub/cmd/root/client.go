package root

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/wuff1996/edgeHub/internal/protobuf"
	"log"
	"net/http"
	"time"
)

const (
	writeWaite = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	Hub  *Hub
	Send chan []byte
	Conn *websocket.Conn
}

var UpGrader = websocket.Upgrader{}

//upgrade the http to websocket with client,register every client to hub
func Servews(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := UpGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade: ", err)
		return
	}
	log.Println("connecting:",r.RemoteAddr)
	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte,256)}
	client.Hub.Register <- client
	go client.readPump()
	go client.WritePump()
}
//send all message from this goroutine
func (c *Client) WritePump() {

	timer := time.NewTimer(pingPeriod)
	defer func() {
		timer.Stop()
		c.Hub.UnRegister<-c
	}()
	for  {
		select {
		case _,ok := <-c.Send:
			err :=c.Conn.SetWriteDeadline(time.Now().Add(writeWaite))
			if err != nil {
            log.Println("set write deadline: ",err)
			}
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage,nil)
				return
			}
			c.Conn.WriteMessage(websocket.BinaryMessage,nil)
			timer.Reset(pingPeriod)

			case <-timer.C:
				err := c.Conn.SetWriteDeadline(time.Now().Add(writeWaite))
				if err != nil {
					log.Println("set write deadline: ",err)
				}
				if err := c.Conn.WriteMessage(websocket.PingMessage,nil);err != nil{
					return
				}
		}

	}
}

//read message from client
func (c *Client) readPump() {
	//unregister client and close the websocket connection
	defer func() {
		c.Hub.UnRegister <- c
		c.Hub.UnRegister<-c
	}()
	if err:=c.Conn.SetReadDeadline(time.Now().Add(pongWait));err!=nil{
		log.Println("set read deadline:",err)
	}
	 c.Conn.SetPingHandler(func(appData string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		log.Println("receive ping: ", appData)
		return nil
	})
	for {
		mt, message, err := c.Conn.ReadMessage()
		{
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read: ", err)
				}
				break
			}
			switch mt {
			case websocket.BinaryMessage:
				connBuf := protobuf.ReadBuf(message)
				b,_:=json.MarshalIndent(&connBuf,""," ")
				fmt.Println(string(b))
				c.Send<-message

			}
		}
	}

}


