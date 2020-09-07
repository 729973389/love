package root

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeHub/internal/protobuf"
	http "net/http"
	"sync"
	"time"
)

//type DeviceInfo struct {
//	DeviceId string      `json:"deviceId"`
//	Data     interface{} `json:"data"`
//}

const (
	writeWaite = 10 * time.Second
	pongWait   = 15 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	Hub          *Hub
	Send         chan []byte
	Conn         *websocket.Conn
	PingPong     chan int
	SerialNumber string
	SendText     chan []byte
}

var UpGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//upgrade the http to websocket with client,register every client to hub
func Servews(ctx context.Context, hub *Hub, w http.ResponseWriter, r *http.Request) {
	defer log.Warning("exit: client: ", r.RemoteAddr)
	ctxChild, cancel := context.WithCancel(ctx)
	conn, err := UpGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade: ", err)
		return
	}
	defer conn.Close()
	var checkCH = make(chan *struct{})
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()
	serialNumber := ""
	go func() {
		serialNumber, err = Check(conn)
		checkCH <- &struct{}{}
	}()
	select {
	case <-timer.C:
		log.Error("check: timeout")
		return
	case <-checkCH:
		break
	}
	if err != nil {
		log.Error(errors.Wrap(err, "servews"))
		return
	}
	timer.Stop()
	close(checkCH)
	if serialNumber == "" {
		return
	}
	log.Println("connecting: ", r.RemoteAddr)
	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 1024), PingPong: make(chan int), SerialNumber: serialNumber, SendText: make(chan []byte, 256)}
	client.Hub.Register <- client
	defer func() {
		client.Hub.UnRegister <- client
		client.Conn.Close()
	}()
	hookWrite := client.HookWrite(ctxChild)
	go func() {
		defer cancel()
		hookWrite()
	}()
	go func() {
		defer cancel()
		//test
		log.Info("reading")
		client.readPump(ctxChild)
	}()
	select {
	case <-ctx.Done():
		log.Warning("close: client: ", ctx.Err())
		return
	case <-ctxChild.Done():
		log.Warning("close: client: ", ctxChild.Err())
		return

	}
}

func (c *Client) HookWrite(ctx context.Context) func() {
	return func() {
		defer log.Warning("EXIT HOOKWRITE")
		ok := c.SendDeviceMap()
		log.Info("SendDeviceMap:", ok)
		c.WritePump(ctx)
	}
}

func (c *Client) SendDeviceMap() bool {
	if len(c.Hub.DeviceMap[c.SerialNumber]) == 0 {
		log.Info("No device")
		return false
	}
	deviceMap := &protobuf.DeviceMap{DeviceId: c.Hub.DeviceMap[c.SerialNumber]}
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
func (c *Client) WritePump(ctx context.Context) {
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
			case websocket.PongMessage:
				timer.Reset(pingPeriod)
			}
		case <-c.SendText:
			err := c.Conn.WriteMessage(websocket.TextMessage, []byte("pong"))
			if err != nil {
				log.Error(errors.Wrap(err, "writePump"))
				return
			}
		case message := <-c.Send:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWaite))
			if err != nil {
				log.Warning("set write deadline: ", err)
			}
			err = c.Conn.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				log.Error(errors.Wrap(err, "client writePump"))
				return
			}
			log.Info(fmt.Sprintf("write remote client: %s: success", c.SerialNumber))
			timer.Reset(pingPeriod)
		case <-timer.C:
			log.Warning("time out")
			return
		case <-ctx.Done():
			log.Info("writePump: ", ctx.Err())
			return
		}
	}
}

//read message from every single client
func (c *Client) readPump(ctx context.Context) {
	defer log.Warning("EXIT READPUMP")
	timer := time.NewTimer(pongWait)
	defer timer.Stop()
	//var once sync.Once
	//unregister client and close the websocket connection
	defer func() {
		log.Warning("closing read: ", c.Conn.RemoteAddr())
		c.Conn.Close()
	}()
	//c.Conn.SetPingHandler(func(appData string) error {
	//	c.Time <- 1
	//	c.PingPong <- websocket.PingMessage
	//	return nil
	//})
	//c.Conn.SetPongHandler(func(appData string) error {
	//	log.Println("receive pong: ", appData)
	//	c.PingPong <- websocket.PongMessage
	//	return nil
	//})
	var timeCH = make(chan *struct{})
	go func() {
		defer c.Conn.Close()
		for {
			select {
			case <-timer.C:
				log.Error("Time out")
				return
				timer.Reset(pongWait)
			case <-ctx.Done():
				log.Warning("read: context: ", ctx.Err())
				return
			}
		}
	}()
	//test
	log.Info("reading inner")
	for {
		mt, message, err := c.Conn.ReadMessage()
		{
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Error("read: ", err)
				}
				log.Error("read: ", err)
				return
			}
			switch mt {
			case websocket.BinaryMessage:
				messageInfo := &protobuf.Message{}
				err := proto.Unmarshal(message, messageInfo)
				if err != nil {
					log.Warning("read: ", err)
					return
				}
				switch messageType := messageInfo.Switch.(type) {
				case *protobuf.Message_Author:
					log.Error("Author massage is not allowed")
					return
				case *protobuf.Message_EdgeInfo:
					edgeInfo := messageType.EdgeInfo
					b, err := proto.Marshal(edgeInfo)
					if err != nil {
						log.Error(err)
						continue
					}
					//test
					log.Info("read: ", string(b))
					c.Send <- b
					c.Hub.HttpMessage <- b
					timeCH <- &struct{}{}
				case *protobuf.Message_DeviceInfo:
					deviceData := messageType.DeviceInfo.Data
					PostDeviceInfo([]byte(deviceData))
				}
			case websocket.TextMessage:
				timeCH <- &struct{}{}
				c.SendText <- message
			}
		}
	}
}

func CheckHmac(author *protobuf.Author) bool {
	var b = []byte(GetConfig().Key)
	hash := hmac.New(sha256.New, b)
	_, err := hash.Write([]byte(author.SerialNumber + author.Time))
	if err != nil {
		log.Error(errors.Wrap(err, "checkHmac"))
		return false
	}
	if hex.EncodeToString(hash.Sum(nil)) == author.Hmac {
		log.Info("checkHmac: success")
		return true
	}
	log.Error("checkHash: false")
	return false
}

func Check(conn *websocket.Conn) (s string, err error) {
	log.Info("checking")
	checkOK := false
	defer func() {
		if checkOK == false {
			log.Error("check: ", checkOK)
			return
		}
		log.Info("check: ", checkOK)
	}()
	mt, message, err := conn.ReadMessage()
	if err != nil {
		return "", errors.Wrap(err, "check")
	}
	if mt != websocket.BinaryMessage {
		return "", fmt.Errorf("check: wrong message type")
	}
	edgeBuf := &protobuf.Message{}
	err = proto.Unmarshal(message, edgeBuf)
	if err != nil {
		return "", errors.Wrap(err, "check")
	}
	switch m := edgeBuf.Switch.(type) {
	case *protobuf.Message_EdgeInfo:
		return "", fmt.Errorf("check: wrong protobuf")
	case *protobuf.Message_Author:
		author := m.Author
		if author.Hmac == "" || author.Token == "" || author.SerialNumber == "" {
			return "", fmt.Errorf("no check information: %v", conn.RemoteAddr())
		}
		if !CheckHmac(author) {
			return "", fmt.Errorf("check: hmac: failed: %v", conn.RemoteAddr())
		}
		if !GEtInfo(author.Token, author.SerialNumber) {
			return "", fmt.Errorf("check: token: failed: %v", conn.RemoteAddr())
		}
		log.Info("Check success:", conn.RemoteAddr())
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch := make(chan error, 1)
			defer close(ch)
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			select {
			case <-ticker.C:
				err = fmt.Errorf("httpRegister: timeout")
				return
			case ch <- PutStatus(author.SerialNumber, true):
				if err = <-ch; err != nil {
					err = errors.Wrap(err, "putStatus")
					return
				}
				checkOK = true
				return
			}
		}()
		wg.Wait()
		return author.SerialNumber, err
	}
	return
}
