package root

//
//import (
//	"crypto/tls"
//	"encoding/json"
//	"flag"
//	"github.com/gorilla/websocket"
//	"github.com/wuff1996/edgeDaemon/internal/protobuf"
//	"log"
//	"os"
//	"os/signal"
//	"sync"
//	"time"
//)
//
//var urlFalg = flag.String("url", "localhost:43211/hub", "set a specific url to connect")
//
//func init() {
//	flag.Parse()
//}
//func Connect() {
//
//	var wg sync.WaitGroup
//	wg.Add(2)
//	var connCh = make(chan *websocket.Conn, 1)
//	var pingCh = make(chan int, 1)
//	interrupt := make(chan os.Signal, 1)
//	signal.Notify(interrupt, os.Interrupt)
//
//	MyDialer := &websocket.Dialer{
//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//	}
//	conn, _, err := MyDialer.Dial("ws://"+*urlFalg, nil)
//	if err != nil {
//		log.Fatal("Dial  :", err)
//	}
//	conn.SetPongHandler(func(appData string) error {
//		log.Println("Receive pong , ", appData)
//		return nil
//	})
//	conn.SetPingHandler(func(appData string) error {
//		log.Println("Receive ping , ", appData)
//		if err != nil {
//			log.Println(err)
//			return err
//		}
//		pingCh <- 1
//		return nil
//	})
//	defer conn.Close()
//	connCh <- conn
//	go func() {
//		for ; ; <-pingCh {
//
//			<-connCh
//			SendPing(conn)
//			connCh <- conn
//
//		}
//	}()
//	go func() {
//		for {
//			select {
//			case <-interrupt:
//				log.Println("interrupt")
//				os.Exit(0)
//				return
//			}
//		}
//	}()
//	go func() {
//		for {
//			systemInfo := protobuf.GetSystemInfo()
//			buf := protobuf.Connect{
//				Id:         "xd",
//				Password:   "xc",
//				SystemInfo: &systemInfo,
//			}
//			<-connCh
//			err := conn.WriteMessage(websocket.BinaryMessage, protobuf.GetBuf(&buf))
//			connCh <- conn
//			if err != nil {
//				log.Println(err)
//			}
//			time.Sleep(3 * time.Minute)
//		}
//	}()
//	go func() {
//		for {
//			<-connCh
//			SendPing(conn)
//			connCh <- conn
//			time.Sleep(50 * time.Second)
//		}
//	}()
//	for {
//		mt, message, err := conn.ReadMessage()
//		if err != nil {
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
//				log.Println("err : ", err)
//
//			}
//			break
//		}
//		switch mt {
//		case websocket.BinaryMessage:
//			connBuf := protobuf.ReadBuf(message)
//			jsonConn, err := json.MarshalIndent(&connBuf, "", " ")
//			if err != nil {
//				log.Println(err)
//			}
//			log.Println(string(jsonConn))
//		}
//	}
//}
//
//func SendPing(conn *websocket.Conn) {
//	log.Println("sendPing")
//	err := conn.WriteMessage(websocket.PingMessage, nil)
//	if err != nil {
//		log.Println("sendPing", err)
//	}
//}
//
////func GetTime() []byte {
////	t := time.Now().UTC().Unix()
////	s := fmt.Sprint(time.Unix(t, 0).UTC())
////	b, err := syscall.ByteSliceFromString(s)
////	if err != nil {
////		log.Println(err)
////	}
////	return b
////}
