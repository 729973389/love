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
	"io"
	"io/ioutil"
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
			PutStatus(c.SerialNumber, nil, false)
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
	req.Header.Add("token", Info.EdgeInfoServerToken)
	resp, err := longClient.Do(req)
	//resp, err := http.Post(route, "application/json", bytes.NewReader(b))
	if err != nil {
		log.Error("Get: ", err)
		return
	}
	defer resp.Body.Close()
	respBuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errors.Wrap(err, "postEdge"))
		return
	}
	log.Println("POST:SEND DATA: ", string(respBuf))
}

type PutProperties struct {
	SerialNumber string                   `json:"serialNumber"`
	Property     *protobuf.EdgeProperties `json:"property,omitempty"`
	Online       bool                     `json:"online"`
}

//put the edge status to the http server
func PutStatus(number string, data *protobuf.EdgeProperties, status bool) error {
	var route = Info.EdgeInfoServer + Info.PutStatus
	properties := &PutProperties{
		SerialNumber: number,
		Property:     data,
		Online:       status,
	}
	b, err := json.Marshal(properties)
	if err != nil {
		return errors.Wrap(err, "putStatus")
	}
	fmt.Println(string(b))
	req, err := http.NewRequest("PUT", route, bytes.NewReader(b))
	if err != nil {
		return errors.Wrap(err, "putStatus")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", Info.EdgeInfoServerToken)
	defer req.Body.Close()
	resp, err := longClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "putStatus")
	}
	defer resp.Body.Close()
	respBuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		return errors.Wrap(err, "putStatus")
	}
	log.Info("putStatus: success: ", string(respBuf))
	return nil

}

//get the specific serialNumber from the http server and check whether the remote token is the same as parameter-list
//t(token) s(serialNumber)
func GEtInfo(t string, s string) bool {
	var route = Info.EdgeInfoServer + Info.GetInfo + "?" + fmt.Sprint("serialNumber=", s)
	req, err := http.NewRequest("GET", route, io.ReadCloser(nil))
	if err != nil {
		log.Error(errors.Wrap(err, "getInfo"))
		return false
	}
	req.Header.Add("token", Info.EdgeInfoServerToken)
	resp, err := longClient.Do(req)
	if err != nil {
		log.Error(err)
		return false
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errors.Wrap(err, "getInfo"))
		return false
	}
	log.Println("GET INFO: ", string(b))
	token, err := FindKeyString(string(b), "token")
	if err != nil {
		log.Error(errors.Wrap(err, "getInfo"))
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
	req, err := http.NewRequest("POST", route, strings.NewReader(string(message)))
	if err != nil {
		log.Error(errors.Wrap(err, "postDeviceInfo"))
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", Info.EdgeInfoServerToken)
	defer req.Body.Close()
	resp, err := longClient.Do(req)
	if err != nil {
		log.Error("PostDeviceInfo: ", err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errors.Wrap(err, "postDeviceInfo"))
	}
	log.Println("PostDeviceInfo: ", string(b))

}

type ResponseCMD struct {
	Uuid string `json:"uuid"`
	Ack  Ack    `json:"ack"`
}

type Ack struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func PutUpdateCMD(data []byte, uuid string) {
	var route = Info.EdgeInfoServer + Info.PutUpdateCMD
	ack := &Ack{}
	if err := json.Unmarshal(data, ack); err != nil {
		log.Error(errors.Wrap(err, "putUpdateCMD"))
		return
	}
	responseCMD := &ResponseCMD{Uuid: uuid, Ack: *ack}
	message, err := json.Marshal(responseCMD)
	if err != nil {
		log.Error(errors.Wrap(err, "putUpdateCMD"))
		return
	}
	log.Info("PutUpdateCMD: ", string(message))
	req, err := http.NewRequest("PUT", route, bytes.NewReader(message))
	if err != nil {
		log.Error(errors.Wrap(err, "putUpdateCMD"))
		return
	}
	defer req.Body.Close()
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", Info.EdgeInfoServerToken)
	resp, err := longClient.Do(req)
	if err != nil {
		log.Error(errors.Wrap(err, "putUpdateCMD"))
		return
	}
	if resp.Body == nil {
		log.Warning("putUpdateCMD: response body is empty")
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Error(fmt.Errorf("putUpdateCMD err: %d", resp.StatusCode))
	}
	readerByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errors.Wrap(err, "putUpdateCMD"))
		return
	}
	log.Info("putUpdateCMD: ", string(readerByte))
}
