package root

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeHub/internal/protobuf"
	"net/http"
)

type HttpClient struct {
	Hub *Hub
}

func Serve(hub *Hub) {
	hc := &HttpClient{Hub: hub}

	for {
		select {
		case message := <-hc.Hub.HttpMessage:
			hc.SendData(message)
		case c := <-hc.Hub.HttpRegister:
			hc.PutStatus(c.SerialNumber, true)
		case c := <-hc.Hub.HttpUnRegister:
			hc.PutStatus(c.SerialNumber, false)

		}
	}
}

func (c *HttpClient) SendData(b []byte) {

	sysInfo := &protobuf.InterfaceEdge{}
	if err := proto.Unmarshal(b, sysInfo); err != nil {
		log.Error(err)
	}
	b, err := json.MarshalIndent(sysInfo, "", " ")
	if err != nil {
		log.Println("jsonMarshal: ", err)
	}
	fmt.Println(string(b))
	resp, err := http.Post(GetConfig().Url+GetConfig().SendData, "application/json", bytes.NewReader(b))
	if err != nil {
		log.Println("Get: ", err)
		return
	}
	defer resp.Body.Close()
	respBuf := make([]byte, 256)
	resp.Body.Read(respBuf)
	log.Println(string(respBuf))
}

func (c *HttpClient) PutStatus(number string, status bool) {
	var rout = GetConfig().Url + GetConfig().PutStatus
	flag.Parse()
	httpInfo := &protobuf.HttpOnline{
		SerialNumber: number,
		Online:       status,
	}
	b, err := json.MarshalIndent(httpInfo, "", " ")
	if err != nil {
		log.Warning(err)
	}
	log.Println(string(b))
	request, err := http.NewRequest("PUT", rout, bytes.NewReader(b))
	request.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Warning(err)
	}
	defer resp.Body.Close()
	var respBuf = make([]byte, 256)
	resp.Body.Read(respBuf)
	log.Println(string(respBuf))

}
