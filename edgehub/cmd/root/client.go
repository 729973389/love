package root

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeHub/internal/protobuf"
	http "net/http"
	"time"
)

//type DeviceInfo struct {
//	DeviceId string      `json:"deviceId"`
//	Data     interface{} `json:"data"`
//}

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
	Time         chan int
}

var UpGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//upgrade the http to websocket with client,register every client to hub
func Servews(ctx context.Context, hub *Hub, w http.ResponseWriter, r *http.Request) {
	defer log.Info("serve quit:", r.RemoteAddr)
	log.Info("serve start:", r.RemoteAddr)
	conn, err := UpGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade: ", err)
		return
	}
	var checkCH = make(chan string, 1)
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()
	go func() {
		checkCH <- Check(conn)
	}()
	var number string
	select {
	case <-timer.C:
		log.Error("check: timeout")
		return
	case number = <-checkCH:
		if number == "" {
			return
		}
		break
	}
	log.Println("connecting:", r.RemoteAddr)
	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 1024), PingPong: make(chan int), SerialNumber: number, Time: make(chan int)}
	client.Hub.Register <- client
	hookWrite := client.HookWrite()
	go hookWrite()
	go client.WritePump()
	go client.readPump()
	for {
		select {
		case <-ctx.Done():
			client.Conn.Close()
			close(client.PingPong)
			log.Warning("Close", "Servews")
			return
		}
	}
}

func (c *Client) HookWrite() func() {
	return func() {
		ok := c.SendDeviceMap()
		log.Info("SendDeviceMap:", ok)
		c.WritePump()
	}
}

func (c *Client) SendDeviceMap() bool {
	if len(c.Hub.deviceMap[c.SerialNumber]) == 0 {
		log.Info("No device")
		return false
	}
	deviceMap := &protobuf.DeviceMap{DeviceId: c.Hub.deviceMap[c.SerialNumber]}
	message := &protobuf.Message{Switch: &protobuf.Message_DeviceMap{DeviceMap: deviceMap}}
	b, err := proto.Marshal(message)
	if err != nil {
		log.Error(errors.Wrap(err, "sendDeviceMap"))
		return false
	}
	if err := c.Conn.WriteMessage(websocket.BinaryMessage, b); err != nil {
		log.Error(errors.Wrap(err, "sendDeviceMap"))
		return false
	}
	return true
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
				//err := c.Conn.SetWriteDeadline(time.Now().Add(writeWaite))
				//if err != nil {
				//	log.Warning("set write deadline: ", err)
				//}
				//err = c.Conn.WriteMessage(websocket.PongMessage, nil)
				//if err != nil {
				//	log.Warning(err)
				//}
				//log.Println("pong: ", c.Conn.RemoteAddr())
				timer.Reset(pingPeriod)
				timerCount = 0
			case websocket.PongMessage:
				timer.Reset(pingPeriod)
				timerCount = 0
			}
		case message := <-c.Send:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWaite))
			if err != nil {
				log.Warning("set write deadline: ", err)
			}
			c.Conn.WriteMessage(websocket.BinaryMessage, message)
			timer.Reset(pingPeriod)
			timerCount = 0
		case <-timer.C:
			if timerCount >= 2 {
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
	timer := time.NewTimer(pongWait)
	defer timer.Stop()
	//var once sync.Once
	//unregister client and close the websocket connection
	defer func() {
		log.Warning("closing read: ", c.Conn.RemoteAddr())
		c.Conn.Close()
	}()
	c.Conn.SetPingHandler(func(appData string) error {
		//err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		//if err != nil {
		//	log.Warning(err)
		//}
		//log.Println("receive ping: ", appData)
		c.Time <- 1
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
	go func() {
		defer c.Conn.Close()
		for {
			select {
			case <-timer.C:
				log.Error("Time out")
			case <-c.Time:
				timer.Reset(pongWait)
			}
		}
	}()
	for {
		//if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		//	log.Warning("set read deadline:", err)
		//}

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
				messageInfo := &protobuf.Message{}
				err := proto.Unmarshal(message, messageInfo)
				if err != nil {
					log.Warning("Read: ", err)
					break
				}
				switch messageType := messageInfo.Switch.(type) {
				case *protobuf.Message_Author:
					log.Error("Author massage is not allowed")
					break
				case *protobuf.Message_EdgeInfo:
					edgeInfo := messageType.EdgeInfo
					b, err := proto.Marshal(edgeInfo)
					if err != nil {
						log.Error(err)
					}
					c.Send <- message
					c.Hub.HttpMessage <- b
					c.Time <- 1
				case *protobuf.Message_DeviceInfo:
					deviceData := messageType.DeviceInfo.Data
					//device := DeviceInfo{DeviceId: deviceInfo.DeviceId, Data: deviceInfo.Data}
					//b, err := json.Marshal(device)
					//if err != nil {
					//	log.Error(errors.Wrap(err, "deviceInfo"))
					//}
					PostDeviceInfo([]byte(deviceData))
				}

			}
		}
	}
}

func CheckHmac(author *protobuf.Author) bool {
	var b = []byte(GetConfig().Key)
	hash := hmac.New(sha256.New, b)

	_, err := hash.Write([]byte(author.SerialNumber + author.Time))
	if err != nil {
		log.Warning(errors.Wrap(err, "checkHmac"))
		return false
	}
	if hex.EncodeToString(hash.Sum(nil)) == author.Hmac {
		log.Println("Check hash success")
		return true
	}
	log.Error("Check hash error")
	return false
}

func Check(conn *websocket.Conn) string {
	log.Info("Checking")
	mt, message, err := conn.ReadMessage()
	if err != nil {
		log.Error(err)
		return ""
	}
	if mt != websocket.BinaryMessage {
		conn.Close()
		return ""
	}
	edgeBuf := &protobuf.Message{}
	err = proto.Unmarshal(message, edgeBuf)
	if err != nil {
		log.Error("Message: ", err)
		return ""
	}
	switch m := edgeBuf.Switch.(type) {
	case *protobuf.Message_EdgeInfo:
		conn.Close()
		return ""
	case *protobuf.Message_Author:
		author := m.Author
		if author.Hmac == "" || author.Token == "" || author.SerialNumber == "" {
			log.Error("No check information: ", conn.RemoteAddr())
			conn.Close()
			return ""
		}
		if !CheckHmac(author) {
			conn.Close()
			return ""
		}
		if !GEtInfo(author.Token, author.SerialNumber) {
			log.Error("Check token failed: ", conn.RemoteAddr())
			conn.Close()
			return ""
		}
		log.Info("Check success:", conn.RemoteAddr())
		PutStatus(author.SerialNumber, true)
		return author.SerialNumber
	}
	conn.Close()
	return ""

}
