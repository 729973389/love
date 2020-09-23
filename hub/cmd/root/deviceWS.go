package root

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeHub/internal/protobuf"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

type DWS struct {
	Conn *websocket.Conn
	Hub  *Hub
	Send chan []byte
}

//remote json type of bindDevice
const bindDevice = "bind_device"

//remote json type of unbindDevice
const unbindDevice = "unbind_device"

//remote json type of deleteEdge
const deleteEdge = "delete_edge"

//demo
const controlDevice = "device_cmd"

//pong time
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
	defer func() {
		_ = c.Close()
		_ = resp.Body.Close()
	}()
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errors.Wrap(err, "DWS: runDWS"))
		return
	}
	_ = resp.Body.Close()
	log.Info(string(respByte))
	dws := &DWS{Hub: hub, Conn: c, Send: make(chan []byte, 1024)}
	wg.Add(1)
	go func() {
		defer func() {
			cancel()
			wg.Done()
		}()
		dws.Read()
	}()
	wg.Add(1)
	go func() {
		defer func() {
			cancel()
			wg.Done()
		}()
		dws.Write(ctxChild)
	}()
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		dws.Loop(ctxChild)
	}()
	wg.Wait()
}

//Write is responsible for write message to websocket client
func (dws *DWS) Write(ctx context.Context) {
	defer func() {
		log.Warning("EXIT DWS WRITE")
		_ = dws.Conn.Close()
	}()
	for {
		select {
		case message := <-dws.Send:
			if err := dws.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Error(errors.Wrap(err, "dws write"))
				continue
			}
			if string(message) == "ping" {
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
	defer func() {
		_ = dws.Conn.Close()
		log.Warning("EXIT: DWS READ")
	}()
	for {
		mt, message, err := dws.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseTLSHandshake) {
				log.Error(errors.Wrap(err, "dws read"))
			}
			log.Info(fmt.Sprintf("messageType: %d: %v", mt, err))
			return
		}
		switch mt {
		case websocket.TextMessage:
			s := string(message)
			if s == "ping" || s == "pong" {
				continue
			}
			//test start
			fmt.Printf("dws textMessage:  %s\n",s)
			//test end
			data, err := FindKey(s, "data")
			if err != nil {
				log.Warning(errors.Wrap(err, "controlDevice"))
				continue
			}
			fs := strings.Replace(s, data, "", 1)
			t, err := FindKeyString(fs, "type")
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
				dws.Hub.Bind(serialNumber, deviceId)
				err = dws.DeviceGister(serialNumber, deviceId, "bind", "")
				if err != nil {
					log.Warning(errors.Wrap(err, "dws read: bindDevice"))
					continue
				}
				continue
			case unbindDevice:
				if err != nil {
					log.Error(errors.Wrap(err, "dws read: unbindDevice"))
					continue
				}
				serialNumber, err := dws.Hub.UnBind(deviceId)
				if err != nil {
					log.Error(errors.Wrap(err, "dws read: unbindDevice"))
					continue
				}
				err = dws.DeviceGister(serialNumber, deviceId, "unbind", "")
				if err != nil {
					log.Error(errors.Wrap(err, "dws read: unbind"))
					continue
				}
				continue
			case deleteEdge:
				serialNumber, err := FindKeyString(s, "serialNumber")
				if err != nil {
					log.Error(errors.Wrap(err, "deleteEdge"))
					continue
				}
				c, ok := dws.Hub.Clients[serialNumber]
				if !ok {
					log.Warning("deleteEdge: no such edge")
					delete(dws.Hub.DeviceMap, serialNumber)
				} else {
					dws.Hub.UnRegister <- c
				}
				log.Info("deleteEdge: ok")
			case controlDevice:
				if err != nil {
					log.Error(errors.Wrap(err, "dws read: controlDevice"))
					continue
				}
				serialNumber, ok := dws.Hub.DeviceSerialNumberMap[deviceId]
				if !ok {
					log.Error(fmt.Errorf("no such deviceId to %s", serialNumber))
					continue
				}
				err := dws.DeviceGister(serialNumber, deviceId, "controlDevice", data)
				if err != nil {
					log.Error(errors.Wrap(err, "dws read: controlDevice"))
					continue
				}
			default:
				log.Warning("No correct information")
			}
		}
	}
}

/*DeviceGister is *DWS function that use specified parameter-list(serialNumber,deviceId,type,data string) to send the oneof
Message_DeviceGister to edgeDaemon websocket client*/
func (dws *DWS) DeviceGister(s, d, t, data string) error {
	if c, ok := dws.Hub.Clients[s]; ok {
		deviceGister := &protobuf.DeviceGister{Type: t, DeviceId: d, Data: data}
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
