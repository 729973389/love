package http

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeHub/internal/protobuf"
	"net/http"
)



type Client struct {
	Schedule *Schedule
}

//func Run(schedule *Schedule) {
//	c := &Client{
//		schedule,
//		make(map[string]*Client),
//		make(chan string,10),
//	}
//	for _,v:=range c.Schedule.Config{
//		select {
//		case url := <-c.Schedule.Action:
//
//
//		}
//
//	}
//}

var urlFalg = flag.String("url", "192.168.32.150:8081", "set a specific url to connect")

func (c *Client) SendData() {

	flag.Parse()
	s := protobuf.GetSystemInfo()
	interfaceEdge := &protobuf.InterfaceEdge{SerialNumber: "1", Data: &s}
	b, err := json.MarshalIndent(interfaceEdge, "", " ")
	if err != nil {
		log.Println("jsonMarshal: ", err)
	}
	fmt.Println(string(b))
	resp, err := http.Post("http://"+*urlFalg+c.Schedule.SendData, "application/json", bytes.NewReader(b))
	if err != nil {
		log.Println("Get: ", err)
		return
	}
	defer resp.Body.Close()
	respBuf := make([]byte, 256)
	resp.Body.Read(respBuf)
	log.Println(string(respBuf))

}

func Serve(s *Schedule) {
	c := &Client{Schedule: s}

	for {
		select {
		case url := <-c.Schedule.Action:
			switch url {
			case c.Schedule.SendData:
				c.SendData()
			}
		}

	}
}
