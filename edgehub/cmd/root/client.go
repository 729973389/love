package root

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeHub/internal/protobuf"
	http "net/http"
	"time"
)

const (
	writeWaite = 10 * time.Second
	pongWait   = 150 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	Hub          *Hub
	Send         chan []byte
	Conn         *websocket.Conn
	PingPong     chan int
	SerialNumber string
}

var UpGrader = websocket.Upgrader{}

//upgrade the http to websocket with client,register every client to hub
func Servews(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := UpGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade: ", err)
		return
	}
	number, ok := Check(conn)
	if !ok {
		return
	}
	log.Println("connecting:", r.RemoteAddr)
	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 1024), PingPong: make(chan int), SerialNumber: number}
	client.Hub.Register <- client
	go client.readPump()
	go client.WritePump()
}

//send all message from this goroutine
func (c *Client) WritePump() {
	timerCount := 0
	timer := time.NewTimer(pingPeriod)
	defer func() {
		timer.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case mt := <-c.PingPong:
			switch mt {
			case websocket.PingMessage:
				err := c.Conn.SetWriteDeadline(time.Now().Add(writeWaite))
				if err != nil {
					log.Warning("set write deadline: ", err)
				}
				err = c.Conn.WriteMessage(websocket.PongMessage, nil)
				if err != nil {
					log.Warning(err)
				}
				log.Println("pong: ", c.Conn.RemoteAddr())
				timer.Reset(pingPeriod)
				timerCount = 0
			case websocket.PongMessage:
				timer.Reset(pingPeriod)
				timerCount = 0

			}
		case message, ok := <-c.Send:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWaite))
			if err != nil {
				log.Warning("set write deadline: ", err)
			}
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			c.Conn.WriteMessage(websocket.BinaryMessage, message)
			timer.Reset(pingPeriod)
			timerCount = 0
		case <-timer.C:
			if timerCount >= 4 {
				log.Warning("time out")
				break
			}
			timerCount++
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWaite))
			if err != nil {
				log.Warning("set write deadline: ", err)
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}

//read message from client
func (c *Client) readPump() {
	//var once sync.Once
	//unregister client and close the websocket connection
	defer func() {
		log.Warning("closing read: ", c.Conn.RemoteAddr())
		c.Conn.Close()
	}()
	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Warning("set read deadline:", err)
	}
	c.Conn.SetPingHandler(func(appData string) error {
		err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Warning(err)
		}
		log.Println("receive ping: ", appData)
		c.PingPong <- websocket.PingMessage
		return nil
	})
	c.Conn.SetPongHandler(func(appData string) error {
		err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Warning(err)
		}
		log.Println("receive pong: ", appData)
		c.PingPong <- websocket.PongMessage
		return nil

	})

	for {
		mt, message, err := c.Conn.ReadMessage()
		{
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Error("read: ", err)
				}
				break
			}
			switch mt {
			case websocket.BinaryMessage:
				//edgeBuf := protobuf.ReadEdge(message)
				//once.Do(func() {
				//	if edgeBuf.Token == "" ||edgeBuf.Hmac==""{
				//		c.Hub.UnRegister <- c
				//		return
				//	}
				//	c.SerialNumber = edgeBuf.SerialNumber
				//	c.Token = edgeBuf.Token
				//	c.Hmac = edgeBuf.Hmac
				//	if !c.CheckHmac() {
				//		c.Hub.UnRegister <- c
				//		return
				//	}
				//	c.Hub.HttpRegister <- c
				//})

				c.Send <- message
				c.Hub.HttpMessage <- message
			}
		}
	}
}

func CheckHmac(m string, s string) bool {
	var b = []byte(GetConfig().Key)
	hash := hmac.New(sha256.New, b)

	_, err := hash.Write([]byte(s))
	if err != nil {
		log.Warning("Get hash: ", err)
	}
	if hex.EncodeToString(hash.Sum(nil)) == m {
		log.Println("Check hash success")
		return true
	}
	log.Error("Check hash error")
	return false
}

func Check(conn *websocket.Conn) (string, bool) {
	mt, message, err := conn.ReadMessage()
	if err != nil {
		log.Error(err)
		return "", false
	}
	if mt != websocket.BinaryMessage {
		conn.Close()
		return "", false
	}
	edgeBuf := protobuf.ReadEdge(message)
	if edgeBuf.Hmac == "" || edgeBuf.Token == "" || edgeBuf.SerialNumber == "" {
		log.Error("No check information: ", conn.RemoteAddr())
		conn.Close()
		return "", false
	}
	if !CheckHmac(edgeBuf.Hmac, edgeBuf.SerialNumber) {
		conn.Close()
		return "", false
	}
	if !GEtInfo(edgeBuf.Token, edgeBuf.SerialNumber) {
		log.Error("Close connection: ", conn.RemoteAddr())
		conn.Close()
		return "", false
	}
	PutStatus(edgeBuf.SerialNumber, true)
	return edgeBuf.SerialNumber, true
}
