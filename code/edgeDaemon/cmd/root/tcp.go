package root

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sacOO7/gowebsocket"
	"github.com/wuff1996/edgeDaemon/internal/protobuf"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const url = "ws://localhost:43211/echo"
var sigch = make(chan string)

func Connet() {
	var wg sync.WaitGroup
	wg.Add(3)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New(url)
	socket.ConnectionOptions = gowebsocket.ConnectionOptions{
		UseSSL:         false,
		UseCompression: false,
		Subprotocols:   []string{"chat", "superchat"},
	}
	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Fatal("Recieved connect error ", err)
	}
	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server" + url)
	}
	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Recieved message  " + message)
	}
	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		c := protobuf.Readbuf(data)
		j, err := json.MarshalIndent(c,"","   ")
		if err != nil {
			log.Println(err)
		}
		log.Println("Receive message " + string(j))

	}
	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved ping " + data)
	}
	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Receive pong "+ data)
	}
	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
		return
	}
	socket.Connect()
	go func() {
		defer wg.Done()
		systemInfo := protobuf.GetSystemInfo()
		buf := protobuf.Connect{
			Password: "wuff",
			Id:       "1996",
			SystemInfo: &systemInfo,
		}
		for i := 0; i < 10; {
			socket.SendBinary(protobuf.Getbuf(&buf))
			i++
		}
	}()
	go func() {
		defer wg.Done()
		for  {
			SendPing(socket)
			fmt.Println("Ping")
			time.Sleep(10 * time.Second)

		}
	}()
	go func() {
		defer wg.Done()
		for  {
			s := <-sigch
			switch s {
			case "ping":
				socket.Close()
				Connet()
				time.Sleep(10 * time.Second)
				continue
			}

		}

	}()
	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			socket.Close()
			Connet()
			return
		}
	}
}

func SendPing(socket gowebsocket.Socket) {
	err :=socket.Conn.WriteMessage(websocket.PingMessage,GetTime())
	if err != nil {
		log.Println(err)
		sigch<-"ping"
	}

}
func GetTime()[]byte{
	t :=time.Now().UTC().Unix()
	s := fmt.Sprint(time.Unix(t,0).UTC())
	b,err :=syscall.ByteSliceFromString(s)
	if err != nil {
		log.Println(err)
	}
	return b
}


