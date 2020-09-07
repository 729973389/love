package root

import (
	"bytes"
	"context"
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

//http serve that be responsible for PostEdge&PutStatus
func Serve(ctx context.Context, hub *Hub) {
	defer log.Warning("QUIT: ", "HTTP SERVE")
	hc := &HttpClient{Hub: hub}
	for {
		select {
		case message, ok := <-hc.Hub.HttpMessage:
			if !ok {
				log.Warning("HttpRegister: ", "channel closed")
				return
			}
			hc.PostEdge(message)
		case c, ok := <-hc.Hub.HttpUnRegister:
			if !ok {
				log.Warning("httpUnregister: ", "channel closed")
				return
			}
			PutStatus(c.SerialNumber, false)
		case <-ctx.Done():
			log.Warning("HTTP SERVE: ", ctx.Err())
			return


		}
	}
}

//post edge information to the http server
func (c *HttpClient) PostEdge(b []byte) {
	var route = Info.EdgeInfoServer + Info.PostEdge
	sysInfo := &protobuf.InterfaceEdge{}
	if err := proto.Unmarshal(b, sysInfo); err != nil {
		log.Error(err)
		return
	}
	b, err := json.Marshal(sysInfo)
	if err != nil {
		log.Warning("jsonMarshal: ", err)
		return
	}
	log.Info("SEND DATA: ", string(b))
	req, err := http.NewRequest("POST", route, bytes.NewReader(b))
	if err != nil {
		log.Error(errors.Wrap(err, "postEdge"))
		return
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

//put the edge status to the http server
func PutStatus(number string, status bool) error {
	var route = Info.EdgeInfoServer + Info.PutStatus
	httpInfo := &protobuf.HttpOnline{
		SerialNumber: number,
		Online:       status,
	}
	b, err := json.Marshal(httpInfo)
	if err != nil {
		return errors.Wrap(err, "putStatus")
	}
	log.Println(string(b))
	request, err := http.NewRequest("PUT", route, bytes.NewReader(b))
	if err != nil {
		return errors.Wrap(err, "putStatus")
	}
	request.Header.Add("Content-Type", "application/json")
	defer request.Body.Close()
	resp, err := longClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "putStatus")
	}
	defer resp.Body.Close()
	var respBuf = make([]byte, 256)
	resp.Body.Read(respBuf)
	log.Info("putStatus: success: ", string(respBuf))
	return nil

}

//get the specific serialNumber from the http server and check whether the remote token is the same as parameter-list t(token)
func GEtInfo(t string, s string) bool {
	var route = Info.EdgeInfoServer + Info.GetInfo + "?" + fmt.Sprint("serialNumber=", s)
	resp, err := http.Get(route)
	if err != nil {
		log.Error(err)
		return false
	}
	defer resp.Body.Close()
	var b = make([]byte, 1024)
	resp.Body.Read(b)
	log.Println("GET INFO: ", string(b))
	token, err := FindKeyString(string(b), "token")
	if err != nil {
		log.Error(errors.Wrap(err, "http getInfo"))
		return false
	}
	if token != t {
		return false
	}
	return true
}

func PostDeviceInfo(message []byte) {
	var route = Info.DeviceInfoServer + Info.PostDevice
	log.Info("PostDeviceInfo: ", string(message))
	request, err := http.NewRequest("POST", route, strings.NewReader(string(message)))
	if err != nil {
		log.Error(errors.Wrap(err, "post deviceInfo"))
		return
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
