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
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeDaemon/internal/protobuf"
	"strings"
	"sync"
	"time"
)

type WS struct {
	Conn         *websocket.Conn
	Send         chan []byte
	PingPong     chan int
	signalCh     chan string
	SerialNumber string
	Token        string
}

var (
	Version = "v1"
	Build   = "N/A"
)

const writeTime = 10 * time.Second
const pongTime = 120 * time.Second
const pingTime = (9 * pongTime) / 10
const ping = 15

var dialer = websocket.Dialer{}

func RunTCP(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	config := GetConfig()
	url := config.Url
	id := config.SerialNumber
	token := config.Token
	c, _, err := dialer.Dial("ws://"+url, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer c.Close()
	time := GetTime()
	message := protobuf.SetOneOfAuthor(id, token, time, GetHashMac(id, time))
	b, err := proto.Marshal(message)
	err = c.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		log.Error(err)
		return
	}
	ws := &WS{
		Conn:         c,
		Send:         make(chan []byte, 1024),
		PingPong:     make(chan int, 5),
		SerialNumber: id,
		signalCh:     make(chan string, 1),
		Token:        token,
	}
	go ws.Write()
	go ws.Read()
	go ws.LoopInfo()

	for {
		select {
		case <-ctx.Done():
			log.Warning("closing run")
			return
		case s := <-ws.signalCh:
			log.Warning("Closing: ", s)
			return
		}
	}
}

func (w *WS) Write() {
	timer := time.NewTimer(pongTime)
	var timerCount = 0
	defer func() {
		timer.Stop()
		w.Conn.Close()
		w.signalCh <- "write"
	}()
	for {
		select {
		case mt := <-w.PingPong:
			switch mt {
			case websocket.PingMessage:
				if err := w.Conn.SetWriteDeadline(time.Now().Add(writeTime)); err != nil {
					log.WithField("set write deadline", err)
				}
				if err := w.Conn.WriteMessage(websocket.PongMessage, nil); err != nil {
					log.WithField("pong", err)
					continue
				}
				log.Println("pong")
			case ping:
				if err := w.Conn.SetWriteDeadline(time.Now().Add(writeTime)); err != nil {
					log.WithField("set write deadline", err)
				}
				if err := w.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.WithField("ping", err)
					continue
				}
				log.Println("ping")
			}
			timer.Reset(pongTime)
			timerCount = 0
		case message, ok := <-w.Send:
			if err := w.Conn.SetWriteDeadline(time.Now().Add(writeTime)); err != nil {
				log.Println("set write deadline: ", err)
			}
			if !ok {
				break
			}
			err := w.Conn.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				log.Println("write message: ", err)
			}
			log.Println("write success")
			timer.Reset(pongTime)
			timerCount = 0
		case <-timer.C:
			if timerCount >= 4 {

			}
			if err := w.Conn.SetWriteDeadline(time.Now().Add(writeTime)); err != nil {
				log.Println("set write deadline: ", err)
			}
			if err := w.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("ping: ", err)
			}
		}
	}
}

func (w *WS) Read() {

	defer func() {
		w.signalCh <- "read"
		w.Conn.Close()
	}()
	if err := w.Conn.SetReadDeadline(time.Now().Add(pongTime)); err != nil {
		log.Println("set read deadline: ", err)
	}
	w.Conn.SetPingHandler(func(appData string) error {
		log.Println("receive ping")
		if err := w.Conn.SetReadDeadline(time.Now().Add(pongTime)); err != nil {
			log.Println("set read deadline: ", err)
		}
		w.PingPong <- websocket.PingMessage
		return nil
	})
	w.Conn.SetPongHandler(func(appData string) error {
		if err := w.Conn.SetReadDeadline(time.Now().Add(pongTime)); err != nil {
			log.Warning(err)
		}
		log.Println("receive pong")

		return nil
	})
	for {
		mt, message, err := w.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseTLSHandshake) {
				log.Println("read: ", err)
			}
			break
		}
		switch mt {
		case websocket.BinaryMessage:
			messageInfo := &protobuf.Message{}
			err := proto.Unmarshal(message, messageInfo)
			if err != nil {
				log.Error("Read: ", err)
			}
			switch t := messageInfo.Switch.(type) {
			case *protobuf.Message_Author:
				break
			case *protobuf.Message_EdgeInfo:
				edgeInfo := t.EdgeInfo
				b, err := json.MarshalIndent(edgeInfo, "", " ")
				if err != nil {
					log.WithError(err)
				}
				fmt.Println(string(b))

			}

		case websocket.CloseMessage:
			log.Error("CONNECTION WAS CLOSED BY HUB")
			return
		}
	}
}

func (w *WS) LoopInfo() {
	go func() {
		for {
			w.PingPong <- ping
			time.Sleep(pingTime)
		}

	}()
	for {
		systemInfo := protobuf.GetSystemInfo()
		edgeInfo := &protobuf.EdgeInfo{
			SerialNumber: w.SerialNumber,
			Data:         &systemInfo,
		}
		message := &protobuf.Message{Switch: &protobuf.Message_EdgeInfo{EdgeInfo: edgeInfo}}
		b, err := proto.Marshal(message)
		if err != nil {
			log.Error("marshal: ", err)
		}
		w.Send <- b
		time.Sleep(30 * time.Minute)
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
	time, err := time.Now().UTC().MarshalText()
	var YYYYMMDDHH string
	if err != nil {
		log.Fatal(err)
	}
	stime := fmt.Sprintf("%s", time)
	clear := []string{"-", ":", "T", "Z"}
	for _, v := range clear {
		stime = strings.Replace(stime, v, "", -1)
	}
	for i, v := range stime {
		if i < 10 {
			YYYYMMDDHH += string(v)
		}
	}
	return YYYYMMDDHH
}
