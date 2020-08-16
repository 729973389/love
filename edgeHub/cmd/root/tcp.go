package root

import (
	"encoding/json"
	"flag"
	"fmt"
	websocket "github.com/gorilla/websocket"
	protobuf "github.com/wuff1996/edgeHub/internal/protobuf"
	"log"
	"net/http"
	"sync"
	"time"
)

var interval = flag.Int("interval", 1, " receive ping longest interval ")

var addr = flag.String("port", ":43211", "http service address")

var upgrader = websocket.Upgrader{} // use default options
var pin = websocket.PingMessage
var pong = websocket.PongMessage

func echo(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	wg.Add(2)
	var connCh = make(chan *websocket.Conn, 1)
	//var pongCh = make (chan int,1)
	//var pingCh=make(chan int,1)
	//pingCh<-1
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	c.SetPingHandler(func(appData string) error {
		c.SetReadDeadline(time.Now().Add(1 * time.Minute))
		log.Println("receiving ping " + appData)
		//<-pingCh
		//<-connCh
		//err := c.WriteMefuncssage(websocket.PongMessage,nil)
		//connCh<-c
		//if err != nil {
		//	log.Println("pong: ",err)
		//	return err
		//}
		return nil
	})
	c.SetPongHandler(func(appData string) error {
		log.Println("receive pong: ",appData)
		//if len(pingCh)>0{
		//	<-pongCh
		//}
		return nil
	})
	defer c.Close()
	connCh <- c
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,websocket.CloseGoingAway,websocket.CloseAbnormalClosure){
				log.Println("read: ",err)
			}
			break
		}
		switch mt {
		case websocket.BinaryMessage:
			connBuf := protobuf.ReadBuf(message)
			jsonBuf, err := json.MarshalIndent(&connBuf, "", " ")
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(string(jsonBuf))
			<-connCh
			err = c.WriteMessage(websocket.BinaryMessage,message)
			connCh<-c
			if err != nil {
				log.Println("writeMessage: ",err)
			}
		}
	}
}

func init() {
	hub := NewHub()
	go hub.Run()
	flag.Parse()

	log.Println("push")
	go func() {
		for  {
			Post()
			time.Sleep(30 * time.Minute)
		}

	}()
	http.HandleFunc("/test", SendMsg)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/hub", func(w http.ResponseWriter, r *http.Request) {
		Servews(hub,w,r)
	})
	fmt.Println("listening", *addr)
	defer log.Fatal(http.ListenAndServe(*addr, nil))

}

func Ping(conn *websocket.Conn) {
	err := conn.WriteMessage(pin, nil)
	log.Println("ping")
	if err != nil {
		log.Println("ping: ",err)
	}
}

//func Pong(conn *websocket.Conn) error {
//	err := conn.WriteMessage(pon, GetTime())
//	if err != nil {
//		log.Println(err)
//	}
//	return err
//}
//
//func GetTime() []byte {
//	t := time.Now().UTC().Unix()
//	s := fmt.Sprint(time.Unix(t, 0).UTC())
//	b, err := syscall.ByteSliceFromString(s)
//	if err != nil {
//		log.Println(err)
//	}
//	return b
//}
//
//func HandlerMessage(c *websocket.Conn,mt int,message []byte,connCh *chan *websocket.Conn){
//	if mt==2 {
//		connBuf :=protobuf.ReadBuf(message)
//		connJSON,err := json.MarshalIndent(connBuf,""," ")
//		if err != nil {
//			log.Println(err)
//		}
//		fmt.Println(string(connJSON))
//		<-*connCh
//		err = c.WriteMessage(websocket.BinaryMessage,message)
//		*connCh<-c
//		if err != nil {
//			log.Println(err)
//		}
//	}
//
//}

//
//go func() {
//	var last,now int64
//	for  {
//		pingGet := <-pingCh
//		if pingCount==0 {
//			pingCount++
//			now=pingGet
//			continue
//		}
//		if pingGet>now {
//			last,now=now,pingGet
//		}
//		if now-last>int64(10*(*interval)){
//			<-connCh
//			pongCh<-now
//			go func() {
//				time.Sleep(5 * time.Second)
//				if len(pongCh)==1{
//					c.Close()
//					return
//				}
//			}()
//			Ping(c)
//			connCh<-c
//		}
//	}
//}()
