package root

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
)

//control all message that transformed by each instance.
type Hub struct {
	//Clients holds the map that specific serialNumber to it's *Client,so you can finds specific *Client when you have a
	//serialNumber.
	Clients map[string]*Client
	//deviceMap holds the map that specific serialNumber to it's deviceId.So it can send the edge-device relation when
	//remote edge restart the program.
	DeviceMap map[string][]string
	//Register is a channel transforming the *Client when every single remote edgeDaemon connect the server.
	Register chan *Client
	//Unregister is a channel transforming the *Client when server has to unregister remote edge,it can just transform.
	//specific *Client
	UnRegister chan *Client
	//HttpMessage is a channel that send the edge information to the http Client,so it can post the information to remote
	//server.
	HttpMessage chan []byte
	//HttpUnregister is a channel, when hub unregister the client it will send *Client to httpClient who will change the
	//specific serialNumber status to false to remote http server.
	HttpUnRegister chan *Client
	//deviceId-serialNumber
	DeviceSerialNumberMap map[string]string
}

/*create a hub that controls the client lifecycle and serve as a platform that
transform message between edge client&http client&device websocket client
*/
func NewHub() *Hub {
	return &Hub{
		Clients:               make(map[string]*Client),
		Register:              make(chan *Client),
		UnRegister:            make(chan *Client),
		HttpMessage:           make(chan []byte, 256),
		HttpUnRegister:        make(chan *Client),
		DeviceMap:             make(map[string][]string),
		DeviceSerialNumberMap: make(map[string]string),
	}
}

//select from channel to register and unregister remote client
func (hub *Hub) Run(ctx context.Context) {
	defer log.Warning("EXIT", "RUN HUB")
	for {
		select {
		case client := <-hub.Register:
			if i, ok := hub.Clients[client.SerialNumber]; ok {
				log.Warning("close last client")
				i.Conn.Close()
				delete(hub.Clients, i.SerialNumber)
			}
			log.Println("register: ", client.Conn.RemoteAddr())
			hub.Clients[client.SerialNumber] = client
		case client := <-hub.UnRegister:
			if _, ok := hub.Clients[client.SerialNumber]; ok {
				if client == nil {

				}
				delete(hub.Clients, client.SerialNumber)
				delete(hub.DeviceMap, client.SerialNumber)
				hub.HttpUnRegister <- client
				log.Warning("unregister: ", client.Conn.RemoteAddr())
				client.Conn.Close()
			}
		case <-ctx.Done():
			log.Warning("HUB RUN: ", ctx.Err())
			return
		}
	}
}

//bind specific serialNumber&deviceId to DeviceSerialNumberMap and DeviceMap
func (hub *Hub) Bind(serialNumber, deviceId string) {
	hub.DeviceSerialNumberMap[deviceId] = serialNumber
	hub.DeviceMap[serialNumber] = append(hub.DeviceMap[serialNumber], deviceId)
	log.Info("bindDevice: ", "ok")
}

//unbind specific deviceId in DeviceSerialNumberMap and DeviceMap
func (hub *Hub) UnBind(deviceId string) (string, error) {
	serialNumber, ok := hub.DeviceSerialNumberMap[deviceId]
	if !ok {
		return "", fmt.Errorf("unBind: no serialNumber: %s in deviceSerialNumberMap", serialNumber)
	}
	delete(hub.DeviceSerialNumberMap, deviceId)
	for i, v := range hub.DeviceMap[serialNumber] {
		if v == deviceId {
			copy(hub.DeviceMap[serialNumber][i:], hub.DeviceMap[serialNumber][i+1:])
			hub.DeviceMap[serialNumber] = hub.DeviceMap[serialNumber][:len(hub.DeviceMap[serialNumber])-1]
			log.Info("unbindDevice: ", "ok")
			return serialNumber, nil
		}
	}
	return "", fmt.Errorf("unBind: no deviceId in deviceMap")
}
