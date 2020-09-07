package root

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeHub/internal/protobuf"
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

//run device websocket that set the device websocket client and hold the read&write&loopSchedule
func RunDWS(ctx context.Context, hub *Hub) {
	defer log.Warning("EXIT RUNDWS")
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
		defer func() {
			wg.Done()
			cancel()
		}()
		dws.Read()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		dws.Write(ctxChild)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		dws.Loop(ctxChild)
	}()
	wg.Wait()

}

//Write is responsible for write message to websocket client
func (dws *DWS) Write(ctx context.Context) {
	defer func() {
		log.Warning("EXIT DWS WRITE")
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
				continue
			}
			log.Info("send: ", string(message))
		case <-ctx.Done():
			return
		}

	}
}

//Read is responsible for reading from remote websocket server,and when read have a mistake,it will return to finish reading,and make all the websocket function stopped.
func (dws *DWS) Read() {
	defer log.Warning("EXIT: DWS READ")
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
			t, err := FindKeyString(s, "type")
			if err != nil {
				log.Error(errors.Wrap(err, "dws read"))
				continue
			}
			if t == "connected" || t == "refused" {
				continue
			}

			switch deviceId, err := FindKeyString(s, "deviceId"); t {
			case bindDevice:
				if err != nil {
					log.Error(errors.Wrap(err, "dws read: bindDevice"))
					continue
				}
				serialNumber, err := FindKeyString(s, "serialNumber")
				if err != nil {
					log.Error(err)
					continue
				}
				if _, ok := dws.Map[deviceId]; !ok {
					err := dws.DeviceGister(serialNumber, deviceId, "bind")
					if err != nil {
						log.Error(errors.Wrap(err, "dws read: bindDevice"))
						continue
					}
					dws.Map[deviceId] = serialNumber
					dws.Hub.DeviceMap[serialNumber] = append(dws.Hub.DeviceMap[serialNumber], deviceId)
					log.Info("bindDevice: ", "ok")
					continue
				}
				log.Warning(fmt.Sprintf("%s already exits", deviceId))
			case unbindDevice:
				if err != nil {
					log.Error(errors.Wrap(err, "dws read: unbindDevice"))
					continue
				}
				serialNumber, ok := dws.Map[deviceId]
				if !ok {
					log.Error(fmt.Errorf("no such deviceId to %s", serialNumber))
					continue
				}
				err := dws.DeviceGister(serialNumber, deviceId, "unbind")
				if err != nil {
					log.Error(errors.Wrap(err, "dws read: unbind"))
					continue
				}
				delete(dws.Map, deviceId)
				for i, v := range dws.Hub.DeviceMap[serialNumber] {
					if v == deviceId {
						copy(dws.Hub.DeviceMap[serialNumber][i:], dws.Hub.DeviceMap[serialNumber][i+1:])
						dws.Hub.DeviceMap[serialNumber] = dws.Hub.DeviceMap[serialNumber][:len(dws.Hub.DeviceMap[serialNumber])-1]
						log.Info("unbindDevice: ", "ok")
						continue
					}
				}
				log.Warning(fmt.Sprintf("dws read: deviceMap: %s has no %s", serialNumber, deviceId))
			case deleteEdge:
				serialNumber, err := FindKeyString(s, "serialNumber")
				if err != nil {
					log.Error(errors.Wrap(err, "deleteEdge"))
					continue
				}
				c, ok := dws.Hub.Clients[serialNumber]
				if !ok {
					log.Warning("deleteEdge: no such edge")
					continue
				}
				dws.Hub.UnRegister <- c
				log.Info("deleteEdge: ok")
			default:
				log.Warning("No correct information")
			}

		}

	}
}

/*DeviceGister is *DWS function that use specified parameter-list(serialNumber,deviceId,type string) to send the oneof
Message_DeviceGister to edgeDaemon websocket client*/
func (dws *DWS) DeviceGister(s, d, t string) error {
	if c, ok := dws.Hub.Clients[s]; ok {
		deviceGister := &protobuf.DeviceGister{Type: t, DeviceId: d}
		message := &protobuf.Message{Switch: &protobuf.Message_DeviceGister{DeviceGister: deviceGister}}
		b, err := proto.Marshal(message)
		if err != nil {
			return errors.Wrap(err, "deviceGister")
		}
		c.Send <- b
	} else {
		return fmt.Errorf(fmt.Sprintf("deviceGister: no such edge %s", s))
	}
	return nil
}

//Loop is responsible for the schedule that is made ,such as ping remote websocket client
func (dws *DWS) Loop(ctx context.Context) {
	timer := time.NewTimer((DWSTime * 9) / 10)
	defer func() {
		timer.Stop()
		log.Info("EXIT DWS LOOP")
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			dws.Send <- []byte("ping")
			timer.Reset((DWSTime * 9) / 10)

		}
	}

}
