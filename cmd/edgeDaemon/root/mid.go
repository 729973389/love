package root

import (
	"context"
	"crypto/tls"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeDaemon/internal/protobuf"
	"strings"
	"sync"
)

type Client struct {
	Hub    *Hub
	Client mqtt.Client
	Send   chan []byte
	Map    map[string]bool
}

var c *Client

const subscribe = "easyfetch/device/properties/"
const publish = "easyfetch/device/cmd/"
const echo = "easyfetch/device/response/"

func RunMQTT(ctx context.Context, hub *Hub) {
	ctxChild, cancel := context.WithCancel(ctx)
	defer log.Warning("EXIT RUNMQTT")
	ops := mqtt.NewClientOptions().SetClientID(Config.SerialNumber).AddBroker("tcp://localhost:7883").
		SetUsername(Config.SerialNumber).
		SetPassword(Config.Token).
		SetTLSConfig(&tls.Config{InsecureSkipVerify: true}).
		SetAutoReconnect(true).
		SetProtocolVersion(4).
		SetCleanSession(true)
	client := mqtt.NewClient(ops)
	c = &Client{Hub: hub, Client: client, Send: make(chan []byte, 1024), Map: make(map[string]bool)}
	if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
		log.Error(errors.Wrap(token.Error(), "connect broker"))
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			cancel()
			wg.Done()
		}()
		c.Receive(ctxChild)
	}()
	//go c.PublishCommand()
	select {
	case <-ctxChild.Done():
		break
	}
	wg.Wait()
}

func (c *Client) Receive(ctx context.Context) {
	defer func() {
		log.Warning("EXIT RECEIVE")
		c.Client.Disconnect(200)
	}()
	for {
		select {
		case deviceId := <-c.Hub.DeviceMap:
			log.Info("register deviceId")
			fmt.Println(deviceId)
			for _, v := range deviceId {
				c.Map[v] = true
				if token := c.Client.Subscribe("easyfetch/device/properties/"+v, 0, messageHandler); token.Wait() && token.Error() != nil {
					log.Warning(errors.Wrap(token.Error(), fmt.Sprintf("register deviceId: %s", v)))
				}
			}

		case deviceGister := <-c.Hub.Down:
			log.Info("receive: ", deviceGister.String())
			deviceId := deviceGister.DeviceId
			if deviceId == "" {
				continue
			}
			switch t := deviceGister.Type; t {
			case "bind":
				log.Info("bind")
				if _, ok := c.Map[deviceId]; ok {
					log.Warning(errors.Wrap(fmt.Errorf(deviceId), "already exit"))
					continue
				}
				c.Map[deviceId] = true
				log.Info("bind: ", deviceId, ": success")
				c.Client.Subscribe(subscribe+deviceId, 0, messageHandler)
			case "unbind":
				log.Info("unbind")
				if _, ok := c.Map[deviceId]; !ok {
					log.Warning("unbind: no deviceId: ", deviceId)
					continue
				}
				c.Client.Unsubscribe(subscribe + deviceId)
				delete(c.Map, deviceId)
			case "controlDevice":
				log.Info("controlDevice")
				if _, ok := c.Map[deviceId]; !ok {
					log.Warning("controlDevice: no deviceId: ", deviceId)
					continue
				}
				uuid, err := FindKeyString(deviceGister.Data, "uuid")
				if err != nil {
					log.Error(errors.Wrap(err, "controlDevice"))
					continue
				}
				uuid = "/" + uuid
				if token := c.Client.Publish(publish+deviceId+uuid, 1, false, deviceGister.Data); token.Wait() && token.Error() != nil {
					log.Error(errors.Wrap(token.Error(), "controlDevice"))
					continue
				}
				c.Client.Subscribe(echo+deviceId+"/#", 1, echoHandler)
				log.Info("controlDevice: ", publish+deviceId+uuid, " :success")
			}
		case <-ctx.Done():
			return
		}
	}
}

type Ack struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var echoHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	topic := message.Topic()
	count := 0
	uuid := ""
	for _, v := range []rune(topic) {
		if string(v) == "/" {
			count++
		}
		if count == 4 && string(v) != "/" {
			uuid += string(v)
		}
	}
	if count != 4 {
		log.Error(fmt.Sprintf("mqtt echoHandler: wrong topic: %s", message.Topic()))
		return
	}
	deviceInfo := &protobuf.DeviceInfo{DeviceType: "response", Data: string(message.Payload()), Uuid: uuid}
	m := &protobuf.Message{Switch: &protobuf.Message_DeviceInfo{DeviceInfo: deviceInfo}}
	c.Hub.Up <- m
	log.Println(string(message.Payload()))
}

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	deviceInfo := &protobuf.DeviceInfo{DeviceType: "task", Data: string(message.Payload())}
	m := &protobuf.Message{Switch: &protobuf.Message_DeviceInfo{DeviceInfo: deviceInfo}}
	c.Hub.Up <- m
	log.Println(string(message.Payload()))
}

//FindKeyString finds string value when given a specified parameter-list(json,key string)
func FindKeyString(s string, key string) (string, error) {
	if !strings.Contains(s, ",") {
		if strings.Contains(s, "\""+key+"\"") {
			t := strings.Split(s, "\"")
			for i, t2 := range t {
				if strings.Contains(t2, ":") {
					if t[i-1] == key {
						return t[i+1], nil
					}
				}
			}
		}
	}
	line := strings.Split(s, ",")
	for _, v := range line {
		if strings.Contains(v, "\""+key+"\":") {
			t := strings.Split(v, "\"")
			for i, t2 := range t {
				if strings.Contains(t2, ":") {
					if t[i-1] == key {
						return t[i+1], nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("FindKeyString: can't find %s from %s", key, s)
}

//func (c *Client) PublishCommand() {
//	for {
//		select {
//		case command := <-c.Hub.Command:
//			deviceId, err := FindKeyString(string(command), "deviceId")
//			if err != nil {
//				log.Error(errors.Wrap(err, "publishCommand"))
//			}
//			if token := c.Client.Publish("easyfetch/device/command"+deviceId, 1, false, command); token.Wait() && token.Error() != nil {
//				log.Warning(token.Error())
//				continue
//			}
//			c.Client.Subscribe("easyfetch/device/"+deviceId+"/command/response", 0, messageHandler)
//		}
//	}
//}
