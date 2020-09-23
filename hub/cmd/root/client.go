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
	"net/http"
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
	SerialNumber string
	SendText     chan []byte
}

var UpGrader = websocket.Upgrader{
	HandshakeTimeout: writeWaite,
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
	//
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
	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 512), SerialNumber: serialNumber, SendText: make(chan []byte, 1024)}
	client.Hub.Register <- client
	defer PutStatus(client.SerialNumber, nil, false)
	defer delete(client.Hub.Clients, client.SerialNumber)
	hookWrite := client.HookWrite(ctxChild)
	go func() {
		defer cancel()
		hookWrite()
	}()
	go func() {
		defer cancel()
		client.readPump(ctxChild, cancel)
	}()
	for {
		select {
		case <-ctxChild.Done():
			time.Sleep(300)
			log.Warning("close: client: ", ctxChild.Err())
			return
		}
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
	for {
		select {
		case <-c.SendText:
			m := []byte("pong")
			err := c.Conn.WriteMessage(websocket.TextMessage, m)
			if err != nil {
				log.Error(errors.Wrap(err, "writePump"))
				return
			}
			//log.Info("send pong")
			continue
		case message := <-c.Send:
			err := c.Conn.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				log.Error(errors.Wrap(err, "client writePump"))
				return
			}
			log.Info(fmt.Sprintf("write remote client: %s: success", c.SerialNumber))
			continue
		case <-ctx.Done():
			log.Info("writePump: ", ctx.Err())
			return
		}
	}
}

//read message from every single client
func (c *Client) readPump(ctx context.Context, cancel context.CancelFunc) {
	defer log.Warning("EXIT READPUMP")
	timer := time.NewTimer(pongWait)
	defer timer.Stop()
	defer func() {
		log.Warning("closing read: ", c.Conn.RemoteAddr())
	}()
	c.Conn.SetPongHandler(func(appData string) error {
		//log.Info("pongHandler")
		return nil
	})
	c.Conn.SetPingHandler(func(appData string) error {
		//log.Info("pingHandler")
		return nil
	})
	var timeCH = make(chan *struct{})
	go func() {
		for {
			select {
			case <-timer.C:
				log.Error("read: timeout")
				cancel()
				return
			case <-ctx.Done():
				log.Warning("read: context: ", ctx.Err())
				return
			case <-timeCH:
				timer.Reset(pongWait)
			}
		}
	}()
	for {
		timeCH <- &struct{}{}
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
					c.Hub.HttpMessage <- b
					continue
				case *protobuf.Message_DeviceInfo:
					deviceInfo := messageType.DeviceInfo
					deviceData := messageType.DeviceInfo.Data
					uuid := messageType.DeviceInfo.Uuid
					switch deviceInfo.DeviceType {
					case "task":
						PostDeviceInfo([]byte(deviceData))
						continue
					case "response":
						PutUpdateCMD([]byte(deviceData), uuid)
					}
				case *protobuf.Message_EdgeProperties:
					edgeProperties := messageType.EdgeProperties
					if err := PutStatus(c.SerialNumber, edgeProperties, true); err != nil {
						log.Error(errors.Wrap(err, "read: edgeProperties"))
						return
					}
					continue
				}
			case websocket.TextMessage:
				//fmt.Println("read: ping")
				c.SendText <- message
				continue
			}
		}
	}
}

func CheckHmac(author *protobuf.Author) bool {
	var b = []byte(Info.Key)
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
	mt, message, err := conn.ReadMessage()
	if err != nil {
		err = errors.Wrap(err, "check")
		return
	}
	if mt != websocket.BinaryMessage {
		return "", fmt.Errorf("check: wrong message type")
	}
	edgeBuf := &protobuf.Message{}
	err = proto.Unmarshal(message, edgeBuf)
	if err != nil {
		err = errors.Wrap(err, "check")
		return
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
		return author.SerialNumber, nil
	}
	return "", fmt.Errorf("check: unexpected err")
}
