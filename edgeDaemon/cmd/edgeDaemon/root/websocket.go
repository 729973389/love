package root

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeDaemon/internal/protobuf"
	"sync"
	"time"
)

type WS struct {
	Hub          *Hub
	Conn         *websocket.Conn
	Send         chan []byte
	PingPong     chan int
	SerialNumber string
	Token        string
}

const writeTime = 10 * time.Second
const pongTime = 13 * time.Second
const pingTime = (9 * pongTime) / 10
const ping = 15
const loopInformation = 30 * time.Minute

var dialer = websocket.Dialer{
	HandshakeTimeout: writeTime,
}

func RunWS(ctx context.Context, hub *Hub) {
	ctxchild, cancel := context.WithCancel(ctx)
	defer cancel()
	pwd := "runWS"
	defer log.Warning("EXIT RUNWS")
	url := Config.Url
	id := Config.SerialNumber
	token := Config.Token
	c, resp, err := dialer.Dial("ws://"+url, nil)
	if err != nil {
		log.Error(errors.Wrap(err, pwd))
		return
	}
	defer func() {
		_ = resp.Body.Close()
		_ = c.Close()
	}()
	respByte := make([]byte, 512)
	_, _ = resp.Body.Read(respByte)
	fmt.Println("resp: ", string(respByte))
	c.SetCloseHandler(func(code int, text string) error {
		log.Error(text+": ", code)
		return err
	})
	myTime := GetTime()
	message := protobuf.SetOneOfAuthor(id, token, myTime, GetHashMac(id, myTime))
	b, err := proto.Marshal(message)
	if err != nil {
		log.Error(errors.Wrap(err, pwd))
		return
	}
	err = c.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		log.Error(errors.Wrap(err, pwd))
		return
	}
	ws := &WS{
		Hub:          hub,
		Conn:         c,
		Send:         make(chan []byte, 1024),
		PingPong:     make(chan int, 32),
		SerialNumber: id,
		Token:        token,
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer cancel()
		defer wg.Done()
		ws.Read()
	}()
	wg.Add(1)
	hookWrite := ws.HookWrite(ctx)
	go func() {
		defer cancel()
		defer wg.Done()
		hookWrite()
	}()
	wg.Add(1)
	go func() {
		defer cancel()
		defer wg.Done()
		ws.LoopInfo(ctxchild)
	}()
	wg.Wait()
}

func (w *WS) HookWrite(ctx context.Context) func() {
	return func() {
		defer func() {
			log.Warning("EXIT: HOOKWRITE")
		}()
		if ok := w.SendEdgeProperties(); !ok {
			log.Error("sendEdgeProperties", ok)
			return
		}
		w.Write(ctx)
	}
}

func (w *WS) SendEdgeProperties() bool {
	properties := protobuf.GetEdgeProperties()
	message := &protobuf.Message{Switch: &protobuf.Message_EdgeProperties{EdgeProperties: &properties}}
	b, err := proto.Marshal(message)
	if err != nil {
		log.Error(errors.Wrap(err, "sendEdgeProperties"))
		return false
	}
	if err := w.Conn.WriteMessage(websocket.BinaryMessage, b); err != nil {
		log.Error(errors.Wrap(err, "sendEdgeProperties"))
		return false
	}
	return true
}

func (w *WS) Write(ctx context.Context) {
	pwd := "write"
	defer func() {
		_ = w.Conn.Close()
		log.Warning("EXIT : WRITE")
	}()
	for {
		select {
		case b := <-w.Hub.Up:
			deviceInfo := &protobuf.DeviceInfo{
				DeviceId: "w",
				Data:     string(b),
			}
			message := &protobuf.Message{Switch: &protobuf.Message_DeviceInfo{DeviceInfo: deviceInfo}}
			bm, err := proto.Marshal(message)
			if err != nil {
				log.Error(errors.Wrap(err, "write : deviceInfo"))
				continue
			}
			if err := w.Conn.WriteMessage(websocket.BinaryMessage, bm); err != nil {
				log.Error(errors.Wrap(err, "write : deviceInfo"))
				return
			}
			log.Info(pwd + ": deviceInfo: success")
		case <-w.PingPong:
			p := []byte("ping")
			if err := w.Conn.WriteMessage(websocket.TextMessage, p); err != nil {
				log.Error("ping : ", err)
				return
			}
			log.Info(fmt.Sprintf("ping : %s : success", w.Conn.RemoteAddr()))
		case message, ok := <-w.Send:
			if !ok {
				log.Error("write: channel closed")
				return
			}
			err := w.Conn.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				log.Error("write: message: ", err)
				return
			}
			log.Println("write: message: success")
		case <-ctx.Done():
			return
		}
	}
}

func (w *WS) Read() {

	defer func() {
		_ = w.Conn.Close()
		log.Warning("EXIT: READ")
	}()
	w.Conn.SetPingHandler(func(appData string) error {
		log.Info("pingHandler")
		return nil
	})
	w.Conn.SetPongHandler(func(appData string) error {
		log.Info("pongHandler")
		return nil
	})
	for {
		mt, message, err := w.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseTLSHandshake) {
				log.Println("read: ", err)
			}
			log.Error("read: ", err)
			return
		}
		switch mt {
		case websocket.BinaryMessage:
			messageInfo := &protobuf.Message{}
			err := proto.Unmarshal(message, messageInfo)
			if err != nil {
				log.Error("read: messageInfo: ", err)
				continue
			}
			switch t := messageInfo.Switch.(type) {
			case *protobuf.Message_Author:
				log.Error("read: wrong protobuf")
				continue
			case *protobuf.Message_EdgeInfo:
				edgeInfo := t.EdgeInfo
				b, err := json.Marshal(edgeInfo)
				if err != nil {
					log.Error(errors.Wrap(err, "read: message_edgeInfo"))
					return
				}
				fmt.Println(string(b))
			case *protobuf.Message_DeviceGister:
				deviceGister := t.DeviceGister
				switch deviceGister.Type {
				case "bind":
					w.Hub.Down <- deviceGister
				case "unbind":
					w.Hub.Down <- deviceGister
				}
			case *protobuf.Message_DeviceMap:
				deviceMap := t.DeviceMap
				w.Hub.DeviceMap <- deviceMap.DeviceId

			}
		case websocket.TextMessage:
			fmt.Println(string(message))
			continue
		}
	}
}

//send schedule information e.g. systemInfo&ping

func (w *WS) LoopInfo(ctx context.Context) {
	defer log.Warning("EXIT : LOOPINFO")
	log.Info("schedule: start")
	timer1 := time.NewTimer(loopInformation)
	defer timer1.Stop()
	systemInfo := protobuf.GetSystemInfo()
	edgeInfo := &protobuf.EdgeInfo{
		SerialNumber: w.SerialNumber,
		Data:         &systemInfo,
	}
	message := &protobuf.Message{Switch: &protobuf.Message_EdgeInfo{EdgeInfo: edgeInfo}}
	b, err := proto.Marshal(message)
	if err != nil {
		log.Error(errors.Wrap(err, "marshal proto"))
		return
	}
	w.Send <- b
	timer2 := time.NewTimer(pingTime)
	defer timer2.Stop()
	for {
		select {
		case <-timer1.C:
			systemInfo := protobuf.GetSystemInfo()
			edgeInfo := &protobuf.EdgeInfo{
				SerialNumber: w.SerialNumber,
				Data:         &systemInfo,
			}
			message := &protobuf.Message{Switch: &protobuf.Message_EdgeInfo{EdgeInfo: edgeInfo}}
			b, err := proto.Marshal(message)
			if err != nil {
				log.Error(errors.Wrap(err, "marshal proto"))
				continue
			}
			w.Send <- b
			timer1.Reset(loopInformation)
		case <-ctx.Done():
			return
		case <-timer2.C:
			w.PingPong <- ping
			timer2.Reset(pingTime)
		}
	}
}

func
GetHashMac(id string, time string) string {
	var key = "3141592666"
	var b = []byte(key)
	hash := hmac.New(sha256.New, b)
	_, err := hash.Write([]byte(id + time))
	if err != nil {
		log.Warning("Get hash: ", err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func GetTime() string {
	myTime := time.Now().UTC()
	return fmt.Sprintf(
		"%d%02d%02dT%02d%02d%02dZ",
		myTime.Year(),
		myTime.Month(),
		myTime.Day(),
		myTime.Hour(),
		myTime.Minute(),
		myTime.Second())
}
