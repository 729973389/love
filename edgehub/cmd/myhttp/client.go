package myhttp

//
//import (
//	"bytes"
//	"encoding/json"
//	"flag"
//	"fmt"
//	"github.com/golang/protobuf/proto"
//	log "github.com/sirupsen/logrus"
//	"github.com/wuff1996/edgeHub/cmd/root"
//	"github.com/wuff1996/edgeHub/internal/protobuf"
//	"net/http"
//)
//
//type HttpClient struct {
//	Schedule *Schedule
//	Hub      *root.Hub
//}
//
////func Run(schedule *Schedule) {
////	c := &Client{
////		schedule,
////		make(map[string]*Client),
////		make(chan string,10),
////	}
////	for _,v:=range c.Schedule.Config{
////		select {
////		case url := <-c.Schedule.Action:
////
////
////		}
////
////	}
////}
//
//var urlFalg = flag.String("url", "192.168.32.150:8081", "set a specific url to connect")
//
//func (c *HttpClient) SendData(b []byte) {
//
//	flag.Parse()
//	sysInfo := &protobuf.InterfaceEdge{}
//	if err := proto.Unmarshal(b, sysInfo); err != nil {
//		log.Error(err)
//	}
//	b, err := json.MarshalIndent(sysInfo, "", " ")
//	if err != nil {
//		log.Println("jsonMarshal: ", err)
//	}
//	fmt.Println(string(b))
//	resp, err := http.Post("http://"+*urlFalg+c.Schedule.SendData, "application/json", bytes.NewReader(b))
//	if err != nil {
//		log.Println("Get: ", err)
//		return
//	}
//	defer resp.Body.Close()
//	respBuf := make([]byte, 256)
//	resp.Body.Read(respBuf)
//	log.Println(string(respBuf))
//
//}
//
//func Serve(s *Schedule, hub *root.Hub) {
//	c := &HttpClient{Schedule: s, Hub: hub}
//
//	for {
//		select {
//		//case url := <-c.Schedule.Action:
//		//switch url {
//		//case c.Schedule.SendData:
//		//	c.SendData()
//		//}
//		case message := <-c.Hub.HttpMessage:
//			c.SendData(message)
//
//			//}
//
//		}
//	}
//}
