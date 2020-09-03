package root

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeHub/internal/protobuf"
	"net"
	"net/http"
	"strings"
	"time"
)

var longClient *http.Client

const (
	MaxIdleConns        int = 100
	MaxIdleConnsPerHost int = 100
)

func init() {
	longClient = &http.Client{
		Timeout: 90 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				KeepAlive: 60 * time.Second,
				Timeout:   60 * time.Second}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
		},
	}
}

type HttpClient struct {
	Hub *Hub
}

func Serve(hub *Hub) {
	defer log.Warning("Close", "Serve")
	hc := &HttpClient{Hub: hub}
	for {
		select {
		case message, ok := <-hc.Hub.HttpMessage:
			if !ok {
				log.Warning("HttpRegister", "channel closed")
				return
			}
			hc.PostEdge(message)
		case c, ok := <-hc.Hub.HttpUnRegister:
			if !ok {
				log.Warning("HttpUnregister ", "channel closed")
				return
			}
			PutStatus(c.SerialNumber, false)

		}
	}
}

func (c *HttpClient) PostEdge(b []byte) {
	var route = Info.EdgeInfoServer + Info.PostEdge
	sysInfo := &protobuf.InterfaceEdge{}
	if err := proto.Unmarshal(b, sysInfo); err != nil {
		log.Error(err)
	}
	b, err := json.MarshalIndent(sysInfo, "", " ")
	if err != nil {
		log.Warning("jsonMarshal: ", err)
	}
	log.Println("SEND DATA: ", string(b))
	req, err := http.NewRequest("POST", route, bytes.NewReader(b))
	if err != nil {
		log.Error(errors.Wrap(err, "postEdge"))
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := longClient.Do(req)
	//resp, err := http.Post(route, "application/json", bytes.NewReader(b))
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
	var route = Info.EdgeInfoServer + Info.PutStatus
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
	resp, err := longClient.Do(request)
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
	var route = Info.EdgeInfoServer + Info.GetInfo + "?" + fmt.Sprint("serialNumber=", s)
	//log.Println(route) //test
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
			for i2, v2 := range tokens {
				if v2 == ":" {
					token = tokens[i2+1]
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

func PostDeviceInfo(message []byte) {
	var route = Info.DeviceInfoServer + Info.PostDevice
	log.Info("PostDeviceInfo: ", string(message))
	request, err := http.NewRequest("POST", route, strings.NewReader(string(message)))
	if err != nil {
		log.Error(errors.Wrap(err, "post deviceInfo"))
	}
	request.Header.Add("Content-Type", "application/json")
	//resp,err:=http.Post(route,"application/json",bytes.NewReader(message))
	resp, err := longClient.Do(request)
	if err != nil {
		log.Error("PostDeviceInfo: ", err)
		return
	}
	defer resp.Body.Close()
	var b = make([]byte, 512)
	resp.Body.Read(b)
	log.Println("PostDeviceInfo: ", string(b))

}
