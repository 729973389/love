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

var (
	Version = "v1"
	Build   = "N/A"
)

const writeTime = 10 * time.Second
const pongTime = 2 * time.Second
const pingTime = (9 * pongTime) / 10
const ping = 15
const loopInformation = 1 * time.Second

var dialer = websocket.Dialer{}

func RunWS(ctx context.Context, hub *Hub) {
	pwd := "runWS"
	defer log.Warning("EXIT RUNWS")
	url := Config.Url
	id := Config.SerialNumber
	token := Config.Token
	c, _, err := dialer.Dial("ws://"+url, nil)
	if err != nil {
		log.Error(errors.Wrap(err, pwd))
		return
	}
	time := GetTime()
	message := protobuf.SetOneOfAuthor(id, token, time, GetHashMac(id, time))
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
		PingPong:     make(chan int, 5),
		SerialNumber: id,
		Token:        token,
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ws.Read()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		ws.Write(ctx)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		ws.LoopInfo(ctx)
	}()
	select {
	case <-ctx.Done():
		log.Warning("closing run")
		return
	}
	wg.Wait()
}

func (w *WS) Write(ctx context.Context) {
	pwd := "write"
	//timer := time.NewTimer(pongTime)
	//var timerCount = 0
	defer func() {
		//timer.Stop()
		w.Conn.Close()
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
				return
			}
			if err := w.Conn.WriteMessage(websocket.BinaryMessage, bm); err != nil {
				log.Error(errors.Wrap(err, "write : deviceInfo"))
				continue
			}
			log.Info(pwd + ": deviceInfo: success")
		case mt := <-w.PingPong:
			switch mt {
			//case websocket.PingMessage:
			//	if err := w.Conn.SetWriteDeadline(time.Now().Add(writeTime)); err != nil {
			//		log.Error("write : setWriteDeadline : ", err)
			//	}
			//	if err := w.Conn.WriteMessage(websocket.PongMessage, nil); err != nil {
			//		log.Error(errors.Wrap(err, "write : pongMessage"))
			//		return
			//	}
			//	log.Println("write : pong : success")
			case ping:
				if err := w.Conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
					log.Error("ping : ", err)
					return
				}
				log.Info(fmt.Sprintf("ping : %s : success", w.Conn.RemoteAddr()))
			}
			//timer.Reset(pongTime)
			//timerCount = 0
		case message, ok := <-w.Send:
			if !ok {
				log.Error("write: channel closed")
				return
			}
			err := w.Conn.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				log.Println("write: message: ", err)
			}
			log.Println("write: message: success")
			continue
		case <-ctx.Done():
			return

			//	timer.Reset(pongTime)
			//	timerCount = 0
			//case <-timer.C:
			//	if timerCount >= 4 {
			//
			//	}
			//	if err := w.Conn.SetWriteDeadline(time.Now().Add(writeTime)); err != nil {
			//		log.Println("set write deadline: ", err)
			//	}
			//	if err := w.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			//		log.Println("ping: ", err)
			//	}
		}
	}
}

func (w *WS) Read() {

	defer func() {
		w.Conn.Close()
		log.Warning("EXIT : READ")
	}()
	//w.Conn.SetPingHandler(func(appData string) error {
	//	log.Println("receive ping")
	//	//if err := w.Conn.SetReadDeadline(time.Now().Add(pongTime)); err != nil {
	//	//	log.Println("set read deadline: ", err)
	//	//}
	//	w.PingPong <- websocket.PingMessage
	//	return nil
	//})
	//w.Conn.SetPongHandler(func(appData string) error {
	//	//if err := w.Conn.SetReadDeadline(time.Now().Add(pongTime)); err != nil {
	//	//	log.Warning(err)
	//	//}
	//	log.Println("receive pong")
	//	return nil
	//})
	for {
		mt, message, err := w.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseTLSHandshake) {
				log.Println("read: message: ", err)
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
				//test
				continue
				//return
			case *protobuf.Message_EdgeInfo:
				edgeInfo := t.EdgeInfo
				b, err := json.Marshal(edgeInfo)
				if err != nil {
					log.WithError(err)
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
		case websocket.CloseMessage:
			log.Warning("CONNECTION WAS CLOSED BY HUB")
			return
		}
	}
}

//send schedule information e.g. systemInfo&ping
func (w *WS) LoopInfo(ctx context.Context) {
	defer log.Warning("EXIT : LOOPINFO")
	log.Info("schedule: start")
	timer1 := time.NewTimer(loopInformation)
	defer timer1.Stop()
	//test
	systemInfo := protobuf.SystemInfo{}
	//systemInfo := protobuf.GetSystemInfo()
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
			//test
			//systemInfo := protobuf.GetSystemInfo()
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
	time := time.Now().UTC()
	return fmt.Sprintf("%d%02d%02dT%02d%02d%02dZ", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second())
	//time, err := time.Now().UTC().MarshalText()
	//var YYYYMMDDHH string
	//if err != nil {
	//	log.Fatal(err)
	//}
	//stime := fmt.Sprintf("%s", time)
	//clear := []string{"-", ":", "T", "Z"}
	//for _, v := range clear {
	//	stime = strings.Replace(stime, v, "", -1)
	//}
	//for i, v := range stime {
	//	if i < 10 {
	//		YYYYMMDDHH += string(v)
	//	}
	//}
	//return YYYYMMDDHH
}
