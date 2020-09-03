package root

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeHub/internal/protobuf"
	"strings"
	"sync"
	"time"
)

type DWS struct {
	Conn *websocket.Conn
	Hub  *Hub
	Send chan []byte
	Map  map[string]string
}

const writeTimeout = 10 * time.Second
const bindDevice = "bind_device"
const unbindDevice = "unbind_device"
const deleteEdge = "delete_edge"
const DWSTime = 15 * time.Second

var dialer = websocket.Dialer{}

func RunDWS(ctx context.Context, hub *Hub) {
	ctxChild, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	c, resp, err := dialer.Dial(Info.DeviceControlServer, nil)
	if err != nil {
		log.Error(errors.Wrap(err, "DWS dial"))
		return
	}
	defer c.Close()
	defer resp.Body.Close()
	respByte := make([]byte, 512)
	resp.Body.Read(respByte)
	resp.Body.Close()
	log.Info(string(respByte))
	dws := &DWS{Hub: hub, Conn: c, Send: make(chan []byte, 1024), Map: make(map[string]string)}
	wg.Add(1)
	go func() {
		defer wg.Done()
		dws.Read()
	}()
	go dws.Write(ctxChild)
	go dws.Loop(ctxChild)
	wg.Wait()
	cancel()

}

func (dws *DWS) Write(ctx context.Context) {
	defer func() {
		log.Warning("EXIT dws write")
		dws.Conn.Close()
	}()
	for {
		select {
		case message := <-dws.Send:
			if err := dws.Conn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
				log.Warning(errors.Wrap(err, "dws write"))
			}
			if err := dws.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Error(errors.Wrap(err, "dws write"))
				return
			}
			log.Info("send", string(message))
		case <-ctx.Done():
			return
		}

	}
}

func (dws *DWS) Read() {
	defer log.Warning("EXIT dws read")
	defer dws.Conn.Close()
	for {
		//if err := dws.Conn.SetReadDeadline(time.Now().Add(DWSTime)); err != nil {
		//	log.Error(errors.Wrap(err, "dws read"))
		//	continue
		//}
		mt, message, err := dws.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseTLSHandshake) {
				log.Error(errors.Wrap(err, "dws read"))
			}
			return
		}
		switch mt {
		case websocket.TextMessage:
			s := string(message)
			log.Info(s)
			if s == "ping" || s == "pong" {
				continue
			}
			t, _ := FindKeyString(s, "type")
			if t == "connected" || t == "refused" {
				continue
			}
			deviceId, err := FindKeyString(s, "deviceId")
			if err != nil {
				log.Warning(err)
				continue
			}
			switch t {
			case bindDevice:
				log.Info("bindDevice")
				serialNumber, err := FindKeyString(s, "serialNumber")
				if err != nil {
					log.Error(err)
					return
				}
				if _, ok := dws.Map[deviceId]; !ok {
					err := dws.DeviceGister(serialNumber, deviceId, "bind")
					if err == nil {
						dws.Map[deviceId] = serialNumber
					}
				} else {
					log.Warning("deviceId already exists")
				}
			case unbindDevice:
				log.Info("unbindDevice")
				serialNumber, ok := dws.Map[deviceId]
				if !ok {
					log.Error(fmt.Errorf("no such deviceId to %s", serialNumber))
					continue
				}
				err := dws.DeviceGister(serialNumber, deviceId, "unbind")
				if err == nil {
					delete(dws.Map, deviceId)
				}
			case deleteEdge:
				log.Info("deleteEdge")
				serialNumber, err := FindKeyString(s, "serialNumber")
				if err != nil {
					log.Error(errors.Wrap(err, "deleteEdge"))
				}
				c, ok := dws.Hub.Clients[serialNumber]
				if !ok {
					log.Warning("deleteEdge: no such edge")
					continue
				}
				dws.Hub.UnRegister <- c
			default:
				log.Warning("No correct information")

			}

		}

	}
}

func (dws *DWS) DeviceGister(s, d, t string) error {
	if c, ok := dws.Hub.Clients[s]; ok {
		deviceGister := &protobuf.DeviceGister{Type: t, DeviceId: d}
		message := &protobuf.Message{Switch: &protobuf.Message_DeviceGister{DeviceGister: deviceGister}}
		b, err := proto.Marshal(message)
		if err != nil {
			log.Error(errors.Wrap(err, t))
			return errors.Wrap(err, t)
		}
		c.Send <- b
	} else {
		log.Info(dws.Hub.Clients)
		log.Error("no such edge")
		return fmt.Errorf("no such edge")
	}
	return nil
}

func FindKeyString(s string, key string) (string, error) {
	line := strings.Split(s, ",")
	for _, v := range line {
		if strings.Contains(v, "\""+key+"\":") {
			t := strings.Split(v, "\"")
			for i, s := range t {
				if strings.Contains(s, ":") {
					if t[i-1] == key {
						return t[i+1], nil
					}
				}
			}
		}
	}
	return "", errors.Wrap(fmt.Errorf("can't find %s", key), "findKeyString")

}

func (dws *DWS) Loop(ctx context.Context) {
	defer func() {
		log.Info("EXIT DWS Loop")
	}()
	go func() {
		for {
			select {
			case dws.Send <- []byte("ping"):
				time.Sleep((DWSTime * 9) / 10)
			case <-ctx.Done():
				return
			}
		}

	}()

}
