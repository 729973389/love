package root

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeHub/internal/protobuf"
	"net/http"
	"strings"
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
		case c := <-hc.Hub.HttpUnRegister:
			PutStatus(c.SerialNumber, false)

		}
	}
}

func (c *HttpClient) SendData(b []byte) {
	var route = GetConfig().Url + GetConfig().SendData
	sysInfo := &protobuf.InterfaceEdge{}
	if err := proto.Unmarshal(b, sysInfo); err != nil {
		log.Error(err)
	}
	b, err := json.MarshalIndent(sysInfo, "", " ")
	if err != nil {
		log.Warning("jsonMarshal: ", err)
	}
	log.Println("SEND DATA: ", string(b))
	resp, err := http.Post(route, "application/json", bytes.NewReader(b))
	if err != nil {
		log.Error("Get: ", err)
		return
	}
	defer resp.Body.Close()
	respBuf := make([]byte, 256)
	resp.Body.Read(respBuf)
	log.Println("POST:SEND DATA: ", string(respBuf))
}

func PutStatus(number string, status bool) {
	var route = GetConfig().Url + GetConfig().PutStatus
	httpInfo := &protobuf.HttpOnline{
		SerialNumber: number,
		Online:       status,
	}
	b, err := json.MarshalIndent(httpInfo, "", " ")
	if err != nil {
		log.Warning(err)
	}
	log.Println(string(b))
	request, err := http.NewRequest("PUT", route, bytes.NewReader(b))
	request.Header.Add("Content-Type", "application/json")
	defer request.Body.Close()
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
	var respBuf = make([]byte, 256)
	resp.Body.Read(respBuf)
	log.Println("PUT:PUT STATUS: ", string(respBuf))

}

func GEtInfo(t string, s string) bool {
	var route = GetConfig().Url + GetConfig().GetInfo + "?" + fmt.Sprint("serialNumber=", s)
	log.Println(route) //test
	var token string
	resq, err := http.Get(route)
	if err != nil {
		log.Error(err)
		return false
	}
	defer resq.Body.Close()
	var b = make([]byte, 512)
	resq.Body.Read(b)
	log.Println("GET INFO: ", string(b))
	ft1 := strings.Split(string(b), ",")
	for _, v := range ft1 {
		if strings.Contains(v, "\"token\":") {
			tokens := strings.Split(v, "\"")
			for _, v2 := range tokens {
				if strings.ContainsAny(v2, "0123456789abcd") {
					token = v2
					break
				}
			}
		}
	}
	if token == "" {
		return false
	}
	if token == t {
		return true
	}
	return false
}
