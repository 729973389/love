package root

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeDaemon/internal/protobuf"
	"os"
	"os/signal"
	"time"
)

var url = GetConfig().Url
var token = GetConfig().Token
var id = GetConfig().SerialNumber

type WS struct {
	Conn     *websocket.Conn
	Send     chan []byte
	PingPong chan int
	signalCh chan os.Signal
}

const writeTime = 10 * time.Second
const pongTime = 120 * time.Second
const pingTime = (9 * pongTime) / 10
const ping = 15

var dialer = websocket.Dialer{}

func RunTCP() {
	c, _, err := dialer.Dial("ws://"+url, nil)
	if err != nil {
		log.Panic(err)
	}
	w := &WS{
		Conn:     c,
		Send:     make(chan []byte, 256),
		PingPong: make(chan int, 5),
		signalCh: make(chan os.Signal, 1),
	}
	signal.Notify(w.signalCh)
	go w.Write()
	go w.Read()
	go w.LoopInfo()
	go func() {
		for {
			select {
			case s := <-w.signalCh:
				switch s {
				case os.Interrupt:
					log.Error("interrupt")
					os.Exit(0)
				case os.Kill:
					log.Error("EXIT")
					os.Exit(0)

				}
			}

		}
	}()

}

func (w *WS) Write() {
	timer := time.NewTimer(pongTime)
	var timerCount = 0
	defer func() {
		timer.Stop()
		w.Conn.Close()
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
				w.signalCh <- os.Interrupt
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
		w.signalCh <- os.Kill
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
			connectInfo := protobuf.ReadBuf(message)
			b, err := json.MarshalIndent(&connectInfo, "", " ")
			if err != nil {
				log.WithError(err)
			}
			fmt.Println(string(b))
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
		//test
		log.Println(url, token, id)
		//
		systemInfo := protobuf.GetSystemInfo()
		edgeInfo := &protobuf.EdgeInfo{
			SerialNumber: id,
			Data:         &systemInfo,
		}
		message := protobuf.GetBufEdgeInfo(edgeInfo)
		w.Send <- message
		time.Sleep(30 * time.Minute)
	}

}
