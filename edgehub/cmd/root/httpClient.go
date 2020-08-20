package root

import (
	"bytes"
	"encoding/json"
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
		log.Warning("jsonMarshal: ", err)
	}
	log.Println("SEND DATA: ", string(b))
	resp, err := http.Post(GetConfig().Url+GetConfig().SendData, "application/json", bytes.NewReader(b))
	if err != nil {
		log.Warning("Get: ", err)
		return
	}
	if resp==nil {
		log.Error("REMOTE HOST CLOSED")
		return

	}
	defer resp.Body.Close()
	respBuf := make([]byte, 256)
	resp.Body.Read(respBuf)
	log.Println("POST:SEND DATA: ", string(respBuf))
}

func (c *HttpClient) PutStatus(number string, status bool) {
	var rout = GetConfig().Url + GetConfig().PutStatus
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
	defer request.Body.Close()
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Warning(err)
	}
	if resp==nil {
		log.Error("REMOTE CLOSED")
		return
	}
	defer resp.Body.Close()
	var respBuf = make([]byte, 256)
	resp.Body.Read(respBuf)
	log.Println("PUT:PUT STATUS: ", string(respBuf))

}
