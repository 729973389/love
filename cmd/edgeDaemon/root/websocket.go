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
	"github.com/wuff1996/edgeDaemon/api/container"
	"github.com/wuff1996/edgeDaemon/internal/protobuf"
	"time"
)

type WS struct {
	Hub           *Hub
	Conn          *websocket.Conn
	Send          chan []byte
	PingPong      chan int
	SerialNumber  string
	Token         string
	Docker        map[string]*container.Docker
	DockerChannel chan *protobuf.Docker
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
		Hub:           hub,
		Conn:          c,
		Send:          make(chan []byte, 1024),
		PingPong:      make(chan int, 32),
		SerialNumber:  id,
		Token:         token,
		Docker:        make(map[string]*container.Docker),
		DockerChannel: make(chan *protobuf.Docker, 32),
	}
	go func() {
		defer cancel()
		ws.Read()
	}()
	hookWrite := ws.HookWrite(ctx)
	go func() {
		defer cancel()
		hookWrite()
	}()
	go func() {
		defer cancel()
		ws.LoopInfo(ctxchild)
	}()
	go func() {
		ws.DockerCTL(ctxchild)
	}()
	for {
		select {
		case <-ctxchild.Done():
			time.Sleep(300)
			log.Warning("close: client: ", ctxchild.Err())
			return
		}
	}
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
			bm, err := proto.Marshal(b)
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
				log.Error("ping: ", err)
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
				case "controlDevice":
					w.Hub.Down <- deviceGister
				}
			case *protobuf.Message_DeviceMap:
				deviceMap := t.DeviceMap
				w.Hub.DeviceMap <- deviceMap.DeviceId
			case *protobuf.Message_Docker:
				docker := t.Docker
				w.DockerChannel <- docker
			}
		case websocket.TextMessage:
			fmt.Println(string(message))
			continue
		}
	}
}

//DockerCTL hold all the docker goroutine
func (w *WS) DockerCTL(ctx context.Context) {
	defer log.Warning("EXIT DOCKERCTL")
	for {
		select {
		case <-ctx.Done():
			return
		case protoDocker := <-w.DockerChannel:
			if protoDocker.Container.Name == "" {
				continue
			}
			if containerDocker := w.Docker[protoDocker.Container.Name]; containerDocker != nil {
				containerDocker.CH <- protoDocker
				continue
			}
			containerDocker := &container.Docker{Docker: protoDocker, Resp: make(chan string, 256), CH: make(chan *protobuf.Docker)}
			w.Docker[protoDocker.Container.Name] = containerDocker
			go func() {
				ctxChild, cancel := context.WithCancel(ctx)
				go func() {
					defer cancel()
					DockerRun(ctxChild, containerDocker, w)
				}()
				go func() {
					defer cancel()
					DockerReceive(ctxChild, containerDocker, w)
				}()
			}()
			containerDocker.CH <- protoDocker
		}
	}

}

//DockerReceive receives message from docker,and can communicate to other channels.
func DockerReceive(ctx context.Context, d *container.Docker, w *WS) {
	for {
		select {
		case <-ctx.Done():
			return
		case response, ok := <-d.Resp:
			if !ok {
				return
			}
			//receive response
			log.Info(response)
		}
	}

}

//DockerRun is responsible for
func DockerRun(ctx context.Context, d *container.Docker, w *WS) {
	defer delete(w.Docker, d.Container.Name)
	for {
		select {
		case <-ctx.Done():
			log.Warning("dockerRun: exit: ", d.Container.Name)
			return
		case protoDocker := <-d.CH:
			switch protoDocker.Type {
			case "pull":
				if err := d.Pull(); err != nil {
					log.Error(errors.Wrap(err, "dockerRun"))
					continue
				}
				if err := d.Run(); err != nil {
					log.Error(errors.Wrap(err, "dockerRun"))
					continue
				}
				log.Info("dockerRun: pull: success")
			case "update":
				d.Docker = protoDocker
				if err := d.Update(); err != nil {
					log.Error(errors.Wrap(err, "dockerRun"))
					continue
				}
				log.Info("dockerRun: update: success")
			case "delete":
				if err := d.Remove(); err != nil {
					log.Error(errors.Wrap(err, "dockerRun"))
					delete(w.Docker, protoDocker.Container.Name)
					log.Error("dockerRun: delete: failed")
					return
				}
				delete(w.Docker, protoDocker.Container.Name)
				log.Info("dockerRun: delete: success")
				return
			default:
				log.Error("type is not right")
				continue
			}
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

func GetHashMac(id string, time string) string {
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
